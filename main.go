package main

import (
	"flag"
	"fmt"
	"github.com/dosarudaniel/CS438_Project/chord"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"github.com/dosarudaniel/CS438_Project/web"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
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
	m := flag.Int("m", 8, "Number of bits in one node's id; max = 256, min = 4 (multiple of 4)")
	guiIPAddr := flag.String("guiIPAddr", "", "ip:port for running the web GUI")
	fixFingerInterval := flag.Int("fixFingerInterval", 1, "Number of seconds between two runs of FixFingers Daemon")
	stabilizeInterval := flag.Int("stabilizeInterval", 1, "Number of seconds between two runs of Stabilize Daemon")
	checkPredecessorInterval := flag.Int("checkPredecessorInterval", 1, "Number of seconds between two runs of CheckPredecessor Daemon")

	flag.Parse()

	log := logrus.New()
	if *trace {
		log.SetLevel(logrus.DebugLevel)
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

	listener, err := net.Listen("tcp", *peersterAddr)
	if err != nil {
		log.Fatal(fmt.Sprintf("listening to %s failed: %v", *peersterAddr, err))
	}

	chordNode, err := chord.NewChordNode(listener, chord.ChordConfig{
		NumOfBitsInID:            *m,
		ChunkSize:                1024,
		StabilizeInterval:        time.Duration(*stabilizeInterval) * time.Second,
		FixFingersInterval:       time.Duration(*fixFingerInterval) * time.Second,
		CheckPredecessorInterval: time.Duration(*checkPredecessorInterval) * time.Second,
	}, *trace)
	if err != nil || chordNode == nil {
		log.Fatal("creating new Chord node failed")
		return // log.Fatal terminates the program; this return is for IDE not to complain that chordNode may be non-nil
	}

	switch {
	case *shouldCreateDHT && !*shouldJoinExistingDHT:
		log.Info("creating a new Chord ring...")

		chordNode.Create()
		fmt.Println(chordNode)

	case *shouldJoinExistingDHT && !*shouldCreateDHT:
		if *existingNodeId == "" || *existingNodeIp == "" {
			log.Fatal("With flag 'join' you also need to provide arguments 'existingNodeIp' and 'existingNodeId'")
		}
		if *existingNodeIp == *peersterAddr {
			log.Fatal("'existingNodeIp' should not be the same as peersterAddr")
		}
		log.Info(fmt.Sprintf("Joining to id: %v ", *existingNodeId))
		log.Info(fmt.Sprintf("Joining to existing node with IP: %v", *existingNodeIp))

		fmt.Println(chordNode)

		err := chordNode.Join(Node{Id: *existingNodeId, Ip: *existingNodeIp})
		if err != nil {
			log.Fatal(fmt.Sprint(err))
		}

	default:
		log.Fatal(fmt.Sprintf("One of the following flags should be true: 'create' or 'join'"))
	}

	if *guiIPAddr != "" {
		go web.RunServer(*guiIPAddr, chordNode)
	}

	// runs infinitely
	select {}
}
