package chord

import (
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
)

// RPC implementation
func (chordNode *ChordNode) TransferFile(fileInfo *FileInfo, stream FileShareService_TransferFileServer) error {
	// Find the file with filename == fileInfo.Filename

	//

	// Create an array of chunks named fileChunks
	//fileChunks := make([][]byte, 5)
	//for _, chunk := range fileChunks {
	//	fileChunk := FileChunk{Content: chunk}
	//	if err := stream.Send(&fileChunk); err != nil {
	//		return err
	//	}
	//
	//}
	fmt.Println("TransferFile was called")
	chunk1 := []byte("Hello, ")
	fileChunk1 := FileChunk{Content: chunk1}
	chunk2 := []byte("World!")
	fileChunk2 := FileChunk{Content: chunk2}

	if err := stream.Send(&fileChunk1); err != nil {
		return err
	}

	if err := stream.Send(&fileChunk2); err != nil {
		return err
	}

	return nil
}