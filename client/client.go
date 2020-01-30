package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	chordService "github.com/dosarudaniel/CS438_Project/services/chord_service"
	clientService "github.com/dosarudaniel/CS438_Project/services/client_service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"strings"
	"time"
)

var log = logrus.New()

func main() {
	peersterAddress := flag.String("PeersterAddress", "", "Peerster address to connect to")
	command := flag.String("command", "", "Command to be sent to Peerster: download/upload/findSuccessor")
	file := flag.String("file", "", "file name at owner")
	ID := flag.String("ID", "", "Download: File owner's ID / FindSuccessor: ID for which the IP is requested ")
	nameToStore := flag.String("nameToStore", "", "Name used to store the downloaded file")
	query := flag.String("query", "", "Search query (required for search command)")
	isWithDownload := flag.Bool("withDownload", false, "Used with search command if you want to download one of the found results")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
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
			log.Debug("Download: No nameToStore given. Using " + *file + " name to store it.")
			*nameToStore = *file
		}

		log.Debug("Sending a download request to " + *peersterAddress)

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
		log.Debug("Sending an upload request to " + *peersterAddress)

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
		log.Debug("Sending a findSuccessor request to " + *peersterAddress)

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

		start := time.Now()
		response, err := findSuccessorClient(client, id)
		if err != nil {
			log.WithField("err", err).Fatal("Fail to findSuccessorClient")
		}
		elapsed := time.Since(start)
		if *verbose == true {
			fmt.Println("Time spent with request:" +  *ID + ",", elapsed.Microseconds())
		}
		fmt.Println(response.Text + response.Info) // Print the IP

	case "search":
		/*
			How to use:
			client/client -PeersterAddress 127.0.0.1:5000 -command search -query="hello"
			client/client -PeersterAddress 127.0.0.1:5000 -command search -query="hello" -withDownload
			client/client -PeersterAddress 127.0.0.1:5000 -command search -query="hello" -withDownload -file newfilename
		*/
		if *query == "" {
			log.Fatal("-query is required but not given")
		}

		fmt.Println("Search query is being processed...")

		conn, err := grpc.Dial(*peersterAddress, grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			log.WithField("err", err).Fatal("Failed to dial")
		}
		defer conn.Close()

		client := clientService.NewClientServiceClient(conn)

		msgFileRecords, err := client.SearchFile(context.Background(), &clientService.Query{
			Query: *query,
		})

		if err != nil {
			log.Error(err)
			return
		}

		fileRecords := msgFileRecords.FileRecords

		for i, fileRecord := range fileRecords {
			if fileRecord != nil {
				fmt.Printf("\t[%d] %s at %s\n", i, fileRecord.Filename, fileRecord.OwnerIp)
			}
		}

		if !*isWithDownload {
			return
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter index of a file you want to download, e.g., 0: ")
		indexString, _ := reader.ReadString('\n')
		index, err := strconv.Atoi(strings.TrimRight(indexString, "\n"))
		if err != nil {
			log.Error(err)
			fmt.Println("Should've entered an integer. Crashing...")
			os.Exit(1)
		}

		ownersID, err := client.KeyToID(context.Background(), &chordService.Key{
			Keyword: fileRecords[index].OwnerIp,
		})
		if err != nil || ownersID == nil {
			fmt.Println("Problem with getting ID from IP. Crashing...")
			os.Exit(1)
		}

		if *file == "" {
			*file = fileRecords[index].Filename
		}

		fileMetadata := &clientService.FileMetadata{
			FilenameAtOwner: fileRecords[index].Filename,
			OwnersID:        ownersID.Id,
			NameToStore:     *file,
		}

		err = requestFile(client, fileMetadata)
		if err != nil {
			fmt.Printf("Fail to requestFile: %v", err)
		}

	default:
		log.Fatal(fmt.Sprintf("No correct command given, try one of the following download/upload/findSuccessor"))
	}
}
