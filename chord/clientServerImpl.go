package chord

import (
	"context"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	. "github.com/dosarudaniel/CS438_Project/services/client_service"
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


// RPC used by the CLI and Web client for communication with one Node from the Chord ring
func (chordNode *ChordNode) FindSuccessorClient(ctx context.Context, id *Identifier) (*Response, error) {
	var ownerNode *Node
	var err error
	responseIp := "nil"

	if chordNode.node.Id == id.Id { // Client asked about this node's IP
		responseIp = chordNode.node.Ip
		fmt.Println("heeeeeeeeeeeeeeeeeere")
	} else {
		// Having an ID, get the Node and its IP
		ownerNode, err = chordNode.FindSuccessor(ctx, &ID{Id: id.Id})
		if err != nil {
			return &Response{Text: "Could not find IP for id given.", Info: responseIp}, err
		}
		responseIp = ownerNode.Ip
	}

	return &Response{Text: "Success! IP found for id given:", Info: responseIp}, nil
}