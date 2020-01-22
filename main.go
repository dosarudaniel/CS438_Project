package main

import (
	"flag"
	"fmt"
	"github.com/dosarudaniel/CS438_Project/logger"
	"math/rand"
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

	switch {
	case *shouldCreateDHT && !*shouldJoinExistingDHT:
		// TODO: Create a new Chord DHT network with one node
		log.Info(fmt.Sprint("Creating a new Chord DHT network ..."))

	case *shouldJoinExistingDHT && !*shouldCreateDHT:
		if *existingNodeId == "" || *existingNodeIp == "" {
			log.Fatal("With flag 'join' you also need to provide arguments 'existingNodeIp' and 'existingNodeId'")
			os.Exit(-1)
		}
		if *existingNodeIp == *peersterAddr {
			log.Fatal("'existingNodeIp' should not be the same as peersterAddr")
		}
		// TODO: This Node joins to a Chord DHT network using the existing Node with existingNodeIp
		log.Info(fmt.Sprint("Joining to id: ", *existingNodeId))
		log.Info(fmt.Sprint("Joining to existing node with IP: ", *existingNodeIp))

	default:
		log.Fatal(fmt.Sprintf("One of the following flags should be true: 'create' or 'join'"))
		os.Exit(-1)
	}
}
