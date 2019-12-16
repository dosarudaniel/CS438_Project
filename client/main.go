package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"github.com/dosarudaniel/CS438_Project/logger"
	"github.com/dosarudaniel/CS438_Project/communication"
	"github.com/dosarudaniel/CS438_Project/types"
	"github.com/dosarudaniel/CS438_Project/file_sharing"
	"github.com/dedis/protobuf"
)

func main() {
	uiPort := flag.String("UIPort", "8080", "port for the UI client")
	msg := flag.String("msg", "", "message to be sent")
	destination := flag.String("dest", "", "destination for the private message; can be omitted")
	file := flag.String("file", "", "file to be indexed by the gossiper")
	request := flag.String("request", "", "request a chunk or metafile of this hash")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	log := logger.DefaultLogger()
	if *verbose {
		log.Level = logger.DebugLevel
	}

	message := prepareMessage(msg, destination, file, request, log)

	err := sendMessage(message, *uiPort)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("Sent message to gossiper at 127.0.0.1:%v :\n\t%v", *uiPort, message))
}

func prepareMessage(msg *string, destination *string, file *string, request *string, log *logger.Logger) types.Message {
	message := types.Message{
		Text: *msg,
	}
	switch communication.CheckClientMessageTypeFromStrings(*msg, *destination, *file, *request) {
	case communication.Rumor:
		// Nothing to do

	case communication.PrivateMessage:
		message.Destination = destination

	case communication.FileSharing:
		message.File = file

	case communication.FileRequest:
		log.Debug(fmt.Sprintf("hash: %v", *request))
		hash, err := file_sharing.HexToHash(*request)
		if err != nil {
			log.Fatal(fmt.Sprintf("Unable to decode hex hash: %v", err))
			fmt.Println("ERROR (Unable to decode hex hash)")
			os.Exit(1)
		}
		message.File = file
		message.Destination = destination
		request := hash[:]
		message.Request = &request
	default:
		log.Fatal("ERROR (Bad argument combination)")
		fmt.Println("ERROR (Bad argument combination)")
		os.Exit(1)
	}
	return message
}

func sendMessage(message types.Message, uiPort string) error {
	packetBytes, err := protobuf.Encode(&message)
	if err != nil {
		return fmt.Errorf("error while encoding the message: %v", err)
	}

	conn, err := net.Dial("udp4", "127.0.0.1:"+uiPort)
	if err != nil {
		return fmt.Errorf("error while resolving gossiper address")
	}

	_, err = conn.Write(packetBytes)
	if err != nil {
		return fmt.Errorf("error while sending message: %v", err)
	}

	return nil
}
