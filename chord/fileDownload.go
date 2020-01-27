package chord

import (
	"context"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"google.golang.org/grpc"
	"io"
	"os"
	"time"
)

// RequestFileFromIP is called whenever a node receives a request from its local client (CLI or Web)
func (chordNode *ChordNode) RequestFileFromIP(filename string, nameToStore string, ownersIp string) error {
	// Create a connection between the owner of the file and the node which request the file
	conn, err := grpc.Dial(ownersIp, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println(fmt.Sprintf("fail to dial: %v", err))
	}
	defer conn.Close()
	client := NewFileShareServiceClient(conn)

	fileInfo := FileInfo{Filename: filename}

	err = Download(client, &fileInfo, nameToStore) // Or `go download(client, &fileInfo)`
	if err != nil {
		fmt.Println(fmt.Sprintf("%v.Download() failed, err = %v", client, err))
		return err
	}

	return nil
}

// Download function uses the TransferFile RPC to receive the file chunks which will be stored under _download/nameToStore
func Download(client FileShareServiceClient, fileInfo *FileInfo, nameToStore string) error {
	log.Printf("Downloading %v", fileInfo) // TODO for logging PR: Log.Info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Request a file from another node by using TransferFile RPC
	stream, err := client.TransferFile(ctx, fileInfo)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v.Download(_) = _, %v", client, err))
		return err
	}
	f, err := os.Create("_download/" + nameToStore)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v.Download(_): Could not create file %v", client, nameToStore))
		return err
	}
	defer f.Close()

	for {
		// Received the file chunks
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("%v.Download(_) = _, %v", client, err))
			return err
		}
		// write in the _download/nameToStore file each chunk received
		n, err := f.Write(chunk.Content)
		if err != nil {
			fmt.Println(fmt.Sprintf("%v.Download(_): Could not write %v bytes into file %v", client, n, nameToStore))
			return err
		}
		log.Println("Received one chunk: " + string(chunk.Content)) // TODO for logging PR: Log.Info
	}

	return nil
}
