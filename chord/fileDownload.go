package chord

import (
	"context"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

// CLIENT part :  A node which request a file from another node

// RequestFileFromIP should be called whenever a node receives a request from its local client (used for interaction)
func (chordNode *ChordNode) RequestFileFromIP(filename string, nameToStore string, ownersIp string) error {

	conn, err := grpc.Dial(ownersIp, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewFileShareServiceClient(conn)

	fileInfo := FileInfo{Filename: filename}

	err = Download(client, &fileInfo, nameToStore) // Or `go download(client, &fileInfo)`
	if err != nil {
		log.Fatalf("%v.Download() failed, err = %v", client, err)
		return err
	}

	return nil
}

// comment?
func Download(client FileShareServiceClient, fileInfo *FileInfo, nameToStore string) error {
	log.Printf("Downloading %v", fileInfo)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.TransferFile(ctx, fileInfo)
	if err != nil {
		log.Fatalf("%v.Download(_) = _, %v", client, err)
		return err
	}
	f, err := os.Create("_Download/" + nameToStore)
	if err != nil {
		log.Fatalf("%v.Download(_): Could not create file %v", client, nameToStore)
		return err
	}
	defer f.Close()
	for {
		// TODO Save the chunks into nameToStore file
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Download(_) = _, %v", client, err)
			return err
		}
		n, err := f.Write(chunk.Content)
		if err != nil {
			log.Fatalf("%v.Download(_): Could not write %v bytes into file %v", client, n, nameToStore)
			return err
		}
		log.Println("Received one chunk: " + string(chunk.Content))
	}

	return nil
}
