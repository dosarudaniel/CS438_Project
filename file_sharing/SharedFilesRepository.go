package file_sharing

import (
	"fmt"
	"github.com/2_alt_hw2/Peerster/logger"
	"io/ioutil"
	"path/filepath"
)

type SharedFilesRepository struct {
	sharedList           []File
	metaHashes           map[string]Hash
	chunks               map[Hash]Chunk
	sharedFilesDirectory string
	log                  *logger.Logger
}

func NewSharedFilesRepository(directory string, log *logger.Logger) (*SharedFilesRepository, error) {
	log.Debug(fmt.Sprintf("SharedFilesRepository init. SharedFiles folder: %v", directory))
	err := createDirectory(directory)
	if err != nil {
		return nil, err
	}

	return &SharedFilesRepository{
		sharedList:           make([]File, 0),
		metaHashes:           make(map[string]Hash, 0),
		chunks:               make(map[Hash]Chunk, 0),
		sharedFilesDirectory: filepath.FromSlash(directory + "/"),
		log:                  log,
	}, nil
}

func (s *SharedFilesRepository) Share(filename string) error {
	_, alreadyShared := s.metaHashes[filename]
	if alreadyShared {
		return fmt.Errorf("file %v is already shared", filename)
	}

	// Reading chunks
	s.log.Debug(fmt.Sprintf("Indexing file %v for sharing", filename))
	chunks, err := fileToChunks(s.sharedFilesDirectory + filename)
	if err != nil {
		return err
	}

	// Computing and saving hashes
	s.log.Trace(fmt.Sprintf("Saving the %v chunks of %v", len(chunks), filename))
	metaFile := make(Chunk, 32*len(chunks))
	for i, chunk := range chunks {
		hash := ChunkToHash(chunk)
		offset := i * 32
		copy(metaFile[offset:offset+32], hash[:])
		s.log.Trace(fmt.Sprintf("\tChunk %v, metafile-offset %v: %v", i, offset, metaFile[offset:offset+32]))
		s.chunks[hash] = chunk
	}

	// Computing meta hash and saving meta file
	metaHash := ChunkToHash(metaFile)
	s.chunks[metaHash] = metaFile
	s.log.Trace(fmt.Sprintf("Metafhash of %v: %v", filename, metaHash))
	s.log.Trace(fmt.Sprintf("Metafile: %v", metaFile))

	// Registering the file as shared
	s.sharedList = append(s.sharedList, File{
		Name: filename,
		Hash: metaHash.ToHex(),
	})
	s.metaHashes[filename] = metaHash

	return nil
}

func (s *SharedFilesRepository) ListShared() []File {
	return s.sharedList
}

func (s *SharedFilesRepository) ListSharable() ([]string, error) {
	files, err := ioutil.ReadDir(s.sharedFilesDirectory)
	if err != nil {
		return nil, err
	}

	sharableFiles := make([]string, 0)
	for _, file := range files {
		_, alreadyShared := s.metaHashes[file.Name()]
		if !alreadyShared {
			sharableFiles = append(sharableFiles, file.Name())
		}
	}

	return sharableFiles, nil
}

func (s *SharedFilesRepository) GetChunk(hash Hash) Chunk {
	chunk, found := s.chunks[hash]
	if !found {
		chunk = []byte{}
	}
	return chunk
}
