package main

import (
	"context"
	"fmt"
	chordService "github.com/dosarudaniel/CS438_Project/services/chord_service"
	clientService "github.com/dosarudaniel/CS438_Project/services/client_service"
	"github.com/sirupsen/logrus"
	"time"
)

// RequestFile RPC caller function
func requestFile(client clientService.ClientServiceClient, fileMetadata *clientService.FileMetadata) (clientService.Response, error) {
	log.WithFields(logrus.Fields{
		"requested file":  fileMetadata.FilenameAtOwner,
		"file owner's ID": fileMetadata.OwnersID,
	}).Info(fmt.Sprintf("Requesting a file..."))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //
	defer cancel()

	response, err := client.RequestFile(ctx, fileMetadata)
	if err != nil {
		fmt.Printf("%v.RequestFile(_) = _, %v: ", client, err)
		return *response, err
	}

	return *response, nil
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
