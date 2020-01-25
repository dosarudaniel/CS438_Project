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

func (ChordNode * ChordNode) RequestFileFromID(filename string, id string) {

	// TODO: Step 1: get IP from ID using Chord
	serverAddr := "127.0.0.123:5000"

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewFileShareServiceClient(conn)

	fileInfo := FileInfo{Filename: filename}

	Download(client, &fileInfo)  // Or `go download(client, &fileInfo)`
}


//
func Download(client FileShareServiceClient, fileInfo *FileInfo) {
	log.Printf("Downloading %v", fileInfo)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.TransferFile(ctx, fileInfo)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Download(_) = _, %v", client, err)
		}
		log.Println(chunk)
	}
}





