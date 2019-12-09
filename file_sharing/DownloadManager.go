package file_sharing

import (
	"fmt"
//========= REDACTED =========
//========= REDACTED =========

	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/logger"
	"os"
	"path/filepath"
	"time"
)

type labeledChunk struct {
	label int
	data  Chunk
}

type download struct {
	File
	metaFile      []Hash
	size          uint64
	missingChunks uint32
	chunks        map[Hash]*labeledChunk
}

type DownloadManager struct {
	name                   string
	downloadFilesDirectory string
	dataReqChannel         chan<- *types.DataRequest

	filesByName     map[string]*download
	wantedMetaFiles map[Hash][]*download
	wantedChunks    map[Hash][]*download

	log *logger.Logger
}

type pendingState int

const (
	notPending      pendingState = iota
	metafilePending pendingState = iota
	chunkPending    pendingState = iota
)

func NewDownloadManager(
	directory string,
	name string,
	log *logger.Logger) (*DownloadManager, <-chan *types.DataRequest, error) {
	log.Debug(fmt.Sprintf("DownloadManager init. Download folder: %v", directory))
	err := createDirectory(directory)
	if err != nil {
		return nil, nil, err
	}

	dataRequestChannel := make(chan *types.DataRequest, 1024)
	return &DownloadManager{
		name:                   name,
		downloadFilesDirectory: filepath.FromSlash(directory + "/"),
		dataReqChannel:         dataRequestChannel,

		filesByName:     make(map[string]*download),
		wantedMetaFiles: make(map[Hash][]*download),
		wantedChunks:    make(map[Hash][]*download),

		log: log,
	}, dataRequestChannel, nil
}

func (dm *DownloadManager) Download(file File) error {
	_, found := dm.filesByName[file.Name]
	if found || fileExists(dm.downloadFilesDirectory+file.Name) {
		return fmt.Errorf("file already exists with name %v", file.Name)
	}

	hash, err := HexToHash(file.Hash)
	if err != nil {
		return err
	}

	downloadFile := &download{
		File:          file,
		metaFile:      nil,
		size:          0,
		missingChunks: 0,
		chunks:        nil,
	}

	dm.insertNewDownload(downloadFile, hash)

	// TODO Check for known chunks for direct-delivery

	dm.sendNewRequest(downloadFile, hash)
	return nil
}

func (dm *DownloadManager) insertNewDownload(downloadFile *download, hash Hash) {
	dm.filesByName[downloadFile.Name] = downloadFile

	files, found := dm.wantedMetaFiles[hash]
	if !found {
		files = make([]*download, 0)
	}
	files = append(files, downloadFile)
	dm.wantedMetaFiles[hash] = files
}

func (dm *DownloadManager) Deliver(origin string, hash Hash, chunk Chunk) error {
	files, found := dm.wantedMetaFiles[hash]
	if found {
		// If empty response, stop wanting this chunk from this origin
		if len(chunk) == 0 {
			tmp := files[:0]
			for _, file := range files {
				if file.Origin != origin {
					tmp = append(tmp, file)
				}
			}
			dm.wantedMetaFiles[hash] = tmp
			return nil
		}

		err := chunk.isValidMetaFile()
		if err != nil {
			return err
		}

		for _, file := range files {
			dm.deliverMetaFile(file, chunk)
		}
		delete(dm.wantedMetaFiles, hash)
	}

	files, found = dm.wantedChunks[hash]
	if found {
		// If empty response, stop wanting this chunk from this origin
		if len(chunk) == 0 {
			tmp := files[:0]
			for _, file := range files {
				if file.Origin != origin {
					tmp = append(tmp, file)
				}
			}
			dm.wantedChunks[hash] = tmp
			return nil
		}

		err := chunk.isValidChunk()
		if err != nil {
			return err
		}

		for _, file := range files {
			labeledChunk := file.chunks[hash]
			labeledChunk.data = chunk
			file.size += uint64(len(chunk))
			file.missingChunks -= 1

			if file.missingChunks == 0 {
				go dm.reconstructFile(file)
			}
		}
		delete(dm.wantedChunks, hash)
	} else {
		return fmt.Errorf("no match for chuck with hash %v", hash)
	}

	return nil
}

func (dm *DownloadManager) reconstructFile(downloadedFile *download) {
	filePath := dm.downloadFilesDirectory + downloadedFile.Name
	file, err := os.Create(filePath)
	if err != nil {
		dm.log.Warn(fmt.Sprintf("Unable to reconstruct file %v: %v", downloadedFile.Name, err))
	}

	for i, hash := range downloadedFile.metaFile {
		_, err := file.Write(downloadedFile.chunks[hash].data)
		if err != nil {
			dm.log.Warn(fmt.Sprintf("Unable to reconstruct file %v. Chunk #%v: %v", downloadedFile.Name, i, err))
		}
	}
	fmt.Printf("RECONSTRUCTED file %v\n", downloadedFile.Name)
}

func (dm *DownloadManager) deliverMetaFile(file *download, data Chunk) {
	nbHash := len(data) / 32
	file.missingChunks = uint32(nbHash)
	file.metaFile = make([]Hash, nbHash)
	file.chunks = make(map[Hash]*labeledChunk, nbHash)
	dm.log.Trace(fmt.Sprintf("Metafile: %v", data))
	dm.log.Trace(fmt.Sprintf("Prepare to request the %v chunks of file %v", nbHash, file.Name))
	for i := 0; i < nbHash; i += 1 {
		var chunkHash Hash
		copy(chunkHash[:], data[i*32:(i+1)*32])
		file.metaFile[i] = chunkHash
		file.chunks[chunkHash] = &labeledChunk{i + 1, nil}
		files, found := dm.wantedChunks[chunkHash]
		if !found {
			files = make([]*download, 0)
		}
		dm.wantedChunks[chunkHash] = append(files, file)
		dm.sendNewRequest(file, chunkHash)
	}
}

func (dm *DownloadManager) sendNewRequest(file *download, hash Hash) {
	request := &types.DataRequest{
		Origin:      dm.name,
		Destination: file.Origin,
		HopLimit:    10,
		HashValue:   hash[:],
	}
	dm.sendRequest(file, hash, request)
}

func (dm *DownloadManager) sendRequest(file *download, hash Hash, request *types.DataRequest) {
	pendingState := dm.isPending(hash, request.Destination)
	switch pendingState {
	case metafilePending:
		fmt.Printf("DOWNLOADING metafile of %v from %v\n", file.Name, file.Origin)
	case chunkPending:
		fmt.Printf("DOWNLOADING %v chunk %v from %v\n", file.Name, file.chunks[hash].label, file.Origin)

	}

	if pendingState != notPending {
		dm.log.Trace(fmt.Sprintf("DownloadManager sent data request %v to gossiper for forwarding", request))
		dm.dataReqChannel <- request

		// Timeout
		time.AfterFunc(time.Duration(10)*time.Second, func() {
			dm.log.Trace(fmt.Sprintf("Timeout for the data request %v", request))
			dm.sendNewRequest(file, hash)
		})
	}
}

func (dm *DownloadManager) isPending(hash Hash, origin string) pendingState {
	files, found := dm.wantedMetaFiles[hash]
	if found {
		for _, file := range files {
			if file.Origin == origin {
				return metafilePending
			}
		}
	}

	files, found = dm.wantedChunks[hash]
	if found {
		for _, file := range files {
			if file.Origin == origin {
				return chunkPending
			}
		}
	}

	return notPending
}
