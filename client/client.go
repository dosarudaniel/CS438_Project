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
	log.Info(fmt.Sprintf("Request file %v from %v", fileMetadata.FilenameAtOwner, fileMetadata.OwnersID))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.RequestFile(ctx, fileMetadata)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v.RequestFile(_) = _, %v: ", client, err))
		return err
	}

	log.Info(response.Text)
	return nil
}

func main() {

	peersterAddress := flag.String("PeersterAddress", "", "Peerster address to connect to")
	file := flag.String("file", "", "file name at owner")
	ownersID := flag.String("ownersID", "", "File owner's ID")
	nameToStore := flag.String("nameToStore", "", "Name used to store the download file")
	// TODO: Add in a future PR: upload a file, add findIP(ID) maybe
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	log := logger.DefaultLogger()
	if *verbose {
		log.Level = logger.DebugLevel
	}

	// Safety checks
	if *peersterAddress == "" {
		fmt.Println("No PeersterAddress given")
	}
	// TODO Add more safety checks in a future PR

	fileMetadata := &clientService.FileMetadata{FilenameAtOwner: *file, OwnersID: *ownersID, NameToStore: *nameToStore}

	conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		fmt.Println(fmt.Sprint("Fail to dial: %v", err))
	}
	defer conn.Close()

	client := clientService.NewClientServiceClient(conn)

	err = requestFile(client, fileMetadata, log)
	if err != nil {
		fmt.Println(fmt.Sprint("Fail to requestFile: %v", err))
	}
}
