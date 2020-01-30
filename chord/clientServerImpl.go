package chord

import (
	"context"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	. "github.com/dosarudaniel/CS438_Project/services/client_service"
	"github.com/golang/protobuf/ptypes/empty"
	txtdist "github.com/masatana/go-textdistance"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
)

// RPC used by the CLI and Web client for communication with one Node from the Chord ring
func (chordNode *ChordNode) RequestFile(ctx context.Context, fileMetadata *FileMetadata) (*Response, error) {
	// Having an ID, first get the Node and its IP
	ownerNode, err := chordNode.FindSuccessor(ctx, &ID{Id: fileMetadata.OwnersID})
	if err != nil {
		return &Response{Text: "Could not find owner's IP for id " + fileMetadata.OwnersID}, err
	}
	// Create a connection with the owner of the file and Download it
	err = chordNode.RequestFileFromIP(fileMetadata.FilenameAtOwner, fileMetadata.NameToStore, ownerNode.Ip)
	if err != nil {
		return &Response{Text: "Could not download file from IP: " + ownerNode.Ip}, err
	}

	return &Response{Text: "Success! File downloaded at _download/" + fileMetadata.NameToStore}, nil
}

func (chordNode *ChordNode) UploadFile(ctx context.Context, msgFilenamePtr *Filename) (*empty.Empty, error) {
	if msgFilenamePtr == nil {
		return nil, nilError("message Filename")
	}

	filename := msgFilenamePtr.Filename

	keywords := strings.Split(cleanFilename(filename), " ")

	for _, keyword := range keywords {
		err := chordNode.PutInDHT(keyword, filename) //TODO impl proper error handling
		if err != nil {
			log.WithFields(logrus.Fields{
				"filename": filename,
				"keyword":  keyword,
				"err":      err,
			}).Warn("putting keyword-fileRecord failed")
		}
	}
	log.Info("Uploading a file has finished")

	return &empty.Empty{}, nil
}

func (chordNode *ChordNode) SearchFile(ctx context.Context, msgQueryPtr *Query) (*FileRecords, error) {
	if msgQueryPtr == nil {
		return nil, nilError("message Query")
	}

	query := msgQueryPtr.Query

	query = cleanFilename(query)

	keywords := strings.Split(query, " ")

	type uniqueFileRecord struct {
		Filename string
		OwnerIP  string
	}

	unionOfFileRecords := make(map[uniqueFileRecord]*FileRecord)
	for _, keyword := range keywords {
		fileRecords, err := chordNode.FindInDHT(keyword)
		if err != nil {
			log.WithFields(logrus.Fields{
				"query":   query,
				"keyword": keyword,
				"err":     err,
			}).Warn("couldn't search for a keyword...")
			continue
		}

		for _, fileRecord := range fileRecords {
			if fileRecord == nil {
				continue
			}
			unionOfFileRecords[uniqueFileRecord{fileRecord.Filename, fileRecord.OwnerIp}] = fileRecord
		}
	}

	allFileRecords := make([]*FileRecord, 0)
	for _, fileRecord := range unionOfFileRecords {
		if fileRecord != nil {
			allFileRecords = append(allFileRecords, fileRecord)
		}
	}

	sort.Slice(allFileRecords, func(i, j int) bool {
		txtdist.LevenshteinDistance(query, allFileRecords[i].Filename)
		return txtdist.LevenshteinDistance(query, allFileRecords[i].Filename) < txtdist.LevenshteinDistance(query, allFileRecords[j].Filename)
	})

	return &FileRecords{
		FileRecords: allFileRecords,
	}, nil
}

// RPC used by the CLI and Web client for communication with one Node from the Chord ring
func (chordNode *ChordNode) FindSuccessorClient(ctx context.Context, id *Identifier) (*Response, error) {
	var node *Node
	var err error
	responseIp := "nil"

	if chordNode.node.Id == id.Id { // Client asked about this node's IP
		responseIp = chordNode.node.Ip
	} else {
		// Having an ID, get the Node and its IP
		node, err = chordNode.FindSuccessor(ctx, &ID{Id: id.Id})
		if err != nil {
			return &Response{Text: "Could not find IP for id given.", Info: responseIp}, err
		}
		responseIp = node.Ip
	}
	return &Response{Text: "Success! IP found for id given:", Info: responseIp}, nil
}

func (chordNode *ChordNode) KeyToID(ctx context.Context, msgKeyPtr *Key) (*ID, error) {
	if msgKeyPtr == nil {
		return nil, nilError("message Key")
	}

	key := msgKeyPtr.Keyword

	id, err := chordNode.hashString(key)
	if err != nil {
		return nil, err
	}

	return &ID{
		Id: id,
	}, nil
}
