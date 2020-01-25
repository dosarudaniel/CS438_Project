package chord

import (
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
)

// RPC implementation
func (chordNode *ChordNode) TransferFile(fileInfo *FileInfo, stream FileShareService_TransferFileServer) error {
	// Find the file with filename == fileInfo.Filename

	//

	// Create an array of chunks named fileChunks
	fileChunks := make([][]byte, 5)
	for _, chunk := range fileChunks {
		fileChunk := FileChunk{Content: chunk}
		if err := stream.Send(&fileChunk); err != nil {
			return err
		}

	}
	return nil
}