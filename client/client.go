package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dosarudaniel/CS438_Project/logger"
	clientService "github.com/dosarudaniel/CS438_Project/services/client_service"
	"google.golang.org/grpc"
	"os"
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
	command := flag.String("command", "" , "Command to be sent to Peerster: download/upload/findSuccessor")
	file := flag.String("file", "", "file name at owner")
	ownersID := flag.String("ownersID", "", "File owner's ID")
	nameToStore := flag.String("nameToStore", "", "Name used to store the download file")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	log := logger.DefaultLogger()
	if *verbose {
		log.Level = logger.DebugLevel
	}

	// Safety checks
	if *peersterAddress == "" {
		log.Fatal("No PeersterAddress given")
		os.Exit(-1)
	}

	switch {
	case *command == "download":
		// Safety checks
		if *file == "" { 			// mandatory
			log.Fatal("Download: No file name given. Specify which file do you request.")
		}
		if *ownersID == "" { 		// mandatory
			log.Fatal("Download: No ID given. Specify the file owner's ID.")
		}
		if *nameToStore == "" {  	// optional
			log.Info("Download: No nameToStore given. Using " + *file + " name to store it.")
			*nameToStore = *file
		}

		log.Info("Sending a download request to " + *peersterAddress)

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

	case *command == "upload":
		log.Info("Sending an upload request to " + *peersterAddress)

	case *command == "findSuccessor":
		log.Info("Sending a findSuccessor request to " + *peersterAddress)

	default:
		log.Fatal(fmt.Sprintf("No correct command given, try one of the following download/upload/findSuccessor"))
		os.Exit(-1)
	}

	// TODO Add more safety checks in a future PR

}
