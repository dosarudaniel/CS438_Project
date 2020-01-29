package main

import (
	"context"
	"flag"
	"fmt"
	chordService "github.com/dosarudaniel/CS438_Project/services/chord_service"
	clientService "github.com/dosarudaniel/CS438_Project/services/client_service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log = logrus.New()

func main() {
	peersterAddress := flag.String("PeersterAddress", "", "Peerster address to connect to")
	command := flag.String("command", "", "Command to be sent to Peerster: download/upload/findSuccessor")
	file := flag.String("file", "", "file name at owner")
	ID := flag.String("ID", "", "Download: File owner's ID / FindSuccessor: ID for which the IP is requested ")
	nameToStore := flag.String("nameToStore", "", "Name used to store the downloaded file")
	query := flag.String("query", "", "Search query (required for search command)")
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
		/*
			How to use:
			client/client -PeersterAddress 127.0.0.1:5000 -command upload -file="hello world"
		*/
		log.Info("Sending an upload request to " + *peersterAddress)

		if *file == "" { // required field
			log.Fatal("Upload: No file name given. Specify which file do you upload.")
		}

		conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			log.WithField("err", err).Fatal("Failed to dial")
		}
		defer conn.Close()

		client := clientService.NewClientServiceClient(conn)
		_, err = client.UploadFile(context.Background(), &clientService.Filename{
			Filename: *file,
		})

		if err != nil {
			log.WithFields(logrus.Fields{
				"filename": *file,
				"err":      err,
			}).Warn("uploading a file failed...")
		}

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

	case "search":
		/*
			How to use:
			client/client -PeersterAddress 127.0.0.1:5000 -command search -query="hello"
		*/
		if *query == "" {
			log.Fatal("-query is required but not given")
		}

		log.Info("Search query is being processed...")

		conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			log.WithField("err", err).Fatal("Failed to dial")
		}
		defer conn.Close()

		log.Info("hi")

		client := clientService.NewClientServiceClient(conn)

		fmt.Println(client.SearchFile(context.Background(), &clientService.Query{
			Query: *query,
		}))

	default:
		log.Fatal(fmt.Sprintf("No correct command given, try one of the following download/upload/findSuccessor"))
	}

}
