package main

import (
	"flag"
	"fmt"
	"github.com/dosarudaniel/CS438_Project/chord"
	"github.com/dosarudaniel/CS438_Project/logger"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	peersterAddr := flag.String("ipAddr", "127.0.0.1:5000", "ip:port for the Peerster")
	name := flag.String("name", "", "name of the Peerster")
	shouldCreateDHT := flag.Bool("create", false, "Pass this flag to create a new Chord ring")
	shouldJoinExistingDHT := flag.Bool("join", false, "Pass this flag to join to an existing Chord ring")
	existingNodeId := flag.String("existingNodeId", "", "The id to which this node should join")
	existingNodeIp := flag.String("existingNodeIp", "", "ip:port for the existing Peerster in the Chord ring to join")
	trace := flag.Bool("v", false, "more verbosity of the program")
	m := flag.Int("m", 32, "Number of bits in one node's id")
	r := flag.Int("r", 2, "Number of nodes in the successor list")

	flag.Parse()

	log := logger.DefaultLogger()
	if *trace {
		log.Level = logger.TraceLevel
	}

	// Setup the random seed
	rand.Seed(time.Now().UnixNano())

	// Assign name
	if *name == "" {
		*name = "Anon_" + strconv.Itoa(1000+rand.Intn(10000))
	}

	log.Info(fmt.Sprint("Peerster IP Address: ", *peersterAddr))
	log.Info(fmt.Sprint("Peerster Name: ", *name))
	log.Info(fmt.Sprint("Number of bits in one node's id: ", *m))
	log.Info(fmt.Sprint("Number of nodes in the successor list: ", *r))

	listener, err := net.Listen("tcp", *peersterAddr)
	if err != nil {
		log.Fatal(fmt.Sprintf("listening to %s failed: %v", *peersterAddr, err))
	}

	chordNode, err := chord.NewChordNode(listener)
	if err != nil {
		log.Fatal("creating new Chord node failed")
		os.Exit(-1)
	}

	switch {
	case *shouldCreateDHT && !*shouldJoinExistingDHT:
		log.Info("creating a new Chord ring...")

		chordNode.Create()
		fmt.Println(chordNode)

	case *shouldJoinExistingDHT && !*shouldCreateDHT:
		if *existingNodeId == "" || *existingNodeIp == "" {
			log.Fatal("With flag 'join' you also need to provide arguments 'existingNodeIp' and 'existingNodeId'")
			os.Exit(-1)
		}
		if *existingNodeIp == *peersterAddr {
			log.Fatal("'existingNodeIp' should not be the same as peersterAddr")
		}
		log.Info(fmt.Sprint("Joining to id: ", *existingNodeId))
		log.Info(fmt.Sprint("Joining to existing node with IP: ", *existingNodeIp))

		err := chordNode.Join(Node{Id: *existingNodeId, Ip: *existingNodeIp})
		if err != nil {
			log.Fatal(fmt.Sprint(err))
			os.Exit(-1)
		}

		fmt.Println(chordNode)

	default:
		log.Fatal(fmt.Sprintf("One of the following flags should be true: 'create' or 'join'"))
		os.Exit(-1)
	}

	select {}
}
