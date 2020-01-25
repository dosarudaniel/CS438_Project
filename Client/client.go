package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dosarudaniel/CS438_Project/logger"
	clientService "github.com/dosarudaniel/CS438_Project/services/client_service"
	"google.golang.org/grpc"
	"time"
)

// printFeature gets the feature for the given point.
func requestFile(client clientService.ClientServiceClient, fileMetadata *clientService.FileMetadata, log *logger.Logger) error {
	log.Info(fmt.Sprint("Request file %v from %v", fileMetadata.FilenameAtOwner, fileMetadata.OwnersID))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.RequestFile(ctx, fileMetadata)
	if err != nil {
		log.Fatal(fmt.Sprint("%v.RequestFile(_) = _, %v: ", client, err))
		return err
	}

	log.Info(response.Text)
	return nil
}

func main() {
	peersterAddress := flag.String("PeersterAddress", "127.0.0.1:5000", "Peerster address to connect to")
	file := flag.String("file", "file1.txt", "file name at owner")
	ownersID := flag.String("ownersID", "f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb", "File owner's ID")
	nameToStore := flag.String("nameToStore", *file, "Name used to store the download file")
	// TODO: Add upload a file command
	verbose := flag.Bool("v", false, "verbose mode")

	// TODO: Add safety check's

	fileMetadata := &clientService.FileMetadata{FilenameAtOwner: *file, OwnersID: *ownersID, NameToStore: *nameToStore}

	flag.Parse()

	log := logger.DefaultLogger()
	if *verbose {
		log.Level = logger.DebugLevel
	}

	conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(),  grpc.WithInsecure(),)
	if err != nil {
		log.Fatal(fmt.Sprint("Fail to dial: %v", err))
	}
	defer conn.Close()


	client := clientService.NewClientServiceClient(conn)

	err = requestFile(client, fileMetadata, log)
	if err != nil {
		log.Fatal(fmt.Sprint("Fail to requestFile: %v", err))
	}

}
