package chord

import (
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"io"
	"os"
)

// Returns true if a file exist at path fileName
func fileExist(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// RPC implementation
func (chordNode *ChordNode) TransferFile(fileInfo *FileInfo, stream FileShareService_TransferFileServer) error {

	fmt.Println("TransferFile was called")

	filePath := "_upload/" + fileInfo.Filename

	if !fileExist(filePath) {
		fmt.Println(fmt.Sprintf("File [%v] does not exist", filePath))
		return nil // Do not stop the server
	} else {
		fmt.Println(fmt.Sprintf("File [%v] does exist", filePath))
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil // Do not stop the server
	}
	defer file.Close()

	buffer := make([]byte, chordNode.config.ChunkSize)

	for {
		bytesRead, err := file.Read(buffer)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break // end of the file
		}

		fmt.Println("Bytes read from the file to string: ", string(buffer[:bytesRead]))

		fileChunk := FileChunk{Content: buffer[:bytesRead]}
		if err := stream.Send(&fileChunk); err != nil { // send the chunk to the client (a node who request this file)
			return err
		}
	}

	return nil
}
