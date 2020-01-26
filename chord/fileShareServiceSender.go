package chord

import (
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"os"
)


func fileExist(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


// RPC implementation
func (chordNode *ChordNode) TransferFile(fileInfo *FileInfo, stream FileShareService_TransferFileServer) error {

	// Create an array of chunks named fileChunks
	//fileChunks := make([][]byte, 5)
	//for _, chunk := range fileChunks {
	//	fileChunk := FileChunk{Content: chunk}
	//	if err := stream.Send(&fileChunk); err != nil {
	//		return err
	//	}
	//
	//}

	filePath := "_Upload/" + fileInfo.Filename

	if !fileExist(filePath) {
		fmt.Println(fmt.Sprintf("File [%v] does not exist", filePath))
		return nil //errors.New(fmt.Sprintf("File [%v] does not exist", fileInfo.Filename))
	} else {
		fmt.Println(fmt.Sprintf("File [%v] does exist", filePath))
	}

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