package chord

import (
	"context"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

// CLIENT part :  A node which request a file from another node

// This function should be called whenever a node receives a  request from its local client (used for interaction)
func (ChordNode * ChordNode) RequestFileFromID(filename string, id string) {

	// TODO: Step 1: get IP from ID using Chord
	serverAddr := "127.0.0.1:5000"

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewFileShareServiceClient(conn)

	fileInfo := FileInfo{Filename: filename}

	Download(client, &fileInfo)  // Or `go download(client, &fileInfo)`
}



func Download(client FileShareServiceClient, fileInfo *FileInfo) {
	log.Printf("Downloading %v", fileInfo)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	stream, err := client.TransferFile(ctx, fileInfo)
	if err != nil {
		log.Fatalf("%v.Download(_) = _, %v", client, err)
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Download(_) = _, %v", client, err)
		}
		log.Println("Received one chunk: " + string(chunk.Content))
	}
}





