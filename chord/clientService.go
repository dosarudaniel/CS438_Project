package chord


import (
	"context"
	. "github.com/dosarudaniel/CS438_Project/services/client_service"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
)

// RPC implementation
func (chordNode *ChordNode) RequestFile(ctx context.Context, fileMetadata *FileMetadata) (*Response, error) {
	ownerNode, err := chordNode.FindSuccessor(ctx, &ID{Id: fileMetadata.OwnersID}) // FIXME
	if err != nil {
		return &Response{Text: "Could not find owner's IP for id " + fileMetadata.OwnersID}, err
	}

	err = chordNode.RequestFileFromIp(fileMetadata.FilenameAtOwner, fileMetadata.NameToStore, ownerNode.Ip)
	if err != nil {
		return &Response{Text: "Could not download file from IP: " + ownerNode.Ip}, err
	}

	return &Response{Text: "Success! File downloaded at _Download/" + fileMetadata.NameToStore}, nil
}