package main

import (
	"context"
	"flag"
	"fmt"
	chordService "github.com/dosarudaniel/CS438_Project/services/chord_service"
	clientService "github.com/dosarudaniel/CS438_Project/services/client_service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

var log = logrus.New()

func main() {
	peersterAddress := flag.String("PeersterAddress", "", "Peerster address to connect to")
	command := flag.String("command", "", "Command to be sent to Peerster: download/upload/findSuccessor")
	file := flag.String("file", "", "file name at owner")
	ID := flag.String("ID", "", "Download: File owner's ID / FindSuccessor: ID for which the IP is requested ")
	nameToStore := flag.String("nameToStore", "", "Name used to store the downloaded file")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.TraceLevel)
	}

	// Safety checks
	if *peersterAddress == "" {
		log.Fatal("No PeersterAddress given")
	}

	switch *command {
	case "download":
		// Safety checks
		if *file == "" || *ID == "" { // required fields
			log.Fatal("Download: No file name given. Specify which file do you request.")
		}

		if *nameToStore == "" { // optional field
			log.Info("Download: No nameToStore given. Using " + *file + " name to store it.")
			*nameToStore = *file
		}

		log.Info("Sending a download request to " + *peersterAddress)

		fileMetadata := &clientService.FileMetadata{FilenameAtOwner: *file, OwnersID: *ID, NameToStore: *nameToStore}

		conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Fail to dial: %v", err)
		}
		defer conn.Close()

		client := clientService.NewClientServiceClient(conn)

		err = requestFile(client, fileMetadata)
		if err != nil {
			fmt.Printf("Fail to requestFile: %v", err)
		}

	case "upload":
		log.Info("Sending an upload request to " + *peersterAddress)

		if *file == "" { // required field
			log.Fatal("Upload: No file name given. Specify which file do you upload.")
		}
		// TODO

	case "findSuccessor":
		log.Info("Sending a findSuccessor request to " + *peersterAddress)

		if *ID == "" { // required
			log.Fatal("FindSuccessor: No ID given. Specify an ID to find the corresponding IP.")
		}

		conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			log.WithField("err", err).Fatal("Failed to dial")
		}
		defer conn.Close()

		client := clientService.NewClientServiceClient(conn)
		id := &chordService.ID{Id: *ID}
		response, err := findSuccessorClient(client, id)
		if err != nil {
			log.WithField("err", err).Fatal("Fail to findSuccessorClient")
		}
		fmt.Println(response.Text + response.Info) // Print the IP

	default:
		log.Fatal(fmt.Sprintf("No correct command given, try one of the following download/upload/findSuccessor"))
	}

}

// RequestFile RPC caller function
func requestFile(client clientService.ClientServiceClient, fileMetadata *clientService.FileMetadata) error {
	log.WithFields(logrus.Fields{
		"requested file":  fileMetadata.FilenameAtOwner,
		"file owner's ID": fileMetadata.OwnersID,
	}).Info(fmt.Sprintf("Requesting a file..."))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // FIXME ctx is not used in rpc, so doesn't work imo
	defer cancel()

	response, err := client.RequestFile(ctx, fileMetadata)
	if err != nil {
		fmt.Printf("%v.RequestFile(_) = _, %v: ", client, err)
		return err
	}

	log.Info(response.Text)
	return nil
}

// FindSuccessorClient RPC caller function
func findSuccessorClient(client clientService.ClientServiceClient, ID *chordService.ID) (clientService.Response, error) {
	log.WithField("ID", ID.Id).Info("Finding a successor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.FindSuccessorClient(ctx, &clientService.Identifier{Id: ID.Id})
	if err != nil {
		fmt.Printf("%v.FindSuccessorClient(_) = _, %v: ", client, err)
		return *response, err
	}

	log.Info(response.Text + response.Info)

	return *response, nil
}
