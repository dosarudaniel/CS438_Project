package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/dosarudaniel/CS438_Project/logger"
	"github.com/dosarudaniel/CS438_Project/communication"
	"github.com/dosarudaniel/CS438_Project/ui/http"
	"github.com/dosarudaniel/CS438_Project/gossiper"
)

func main() {
	uiPort := flag.String("UIPort", "8080", "port for the UI client")
	guiPort := flag.String("GUIPort", "8080", "port for the graphical UI client")
	gossipAddr := flag.String("gossipAddr", "127.0.0.1:5000", "ip:port for the gossiper")
	name := flag.String("name", "", "name of the gossiper")
	peers := flag.String("peers", "", "command separated list of peers of the form ip:port")
	simple := flag.Bool("simple", false, "run gossiper in simple broadcast mode")
	debug := flag.Bool("v", false, "verbosity of the program")
	trace := flag.Bool("vv", false, "more verbosity of the program")
	antiEntropy := flag.Int("antiEntropy", 10, "use the given timeout in seconds for anti-entropy")
	routeTimer := flag.Int("rtimer", 0, "Timeout in seconds to send route rumors. 0 (default) means disable sending route rumors.")

	flag.Parse()

	log := logger.DefaultLogger()

	if *trace {
		log.Level = logger.TraceLevel
	} else if *debug {
		log.Level = logger.DebugLevel
	}

	// Setup the random seed
	rand.Seed(time.Now().UnixNano())

	// Check name
	if *name == "" {
		*name = "Anon_" + strconv.Itoa(1000+rand.Intn(10000))
	}

	log.Info(fmt.Sprint("UIPort: ", *uiPort))
	log.Info(fmt.Sprint("gossipAddr: ", *gossipAddr))
	log.Info(fmt.Sprint("name: ", *name))
	log.Info(fmt.Sprint("peers: ", *peers))
	log.Info(fmt.Sprint("simple: ", *simple))
	log.Info(fmt.Sprint("antiEntropy: ", *antiEntropy))
	log.Info(fmt.Sprint("routeTimer: ", *routeTimer))

	// Setup client communication
	clientInput, clientPacketHandler := communication.NewClientPacketDecoder(log)
	_, err := communication.NewUDPListener("127.0.0.1:"+*uiPort, clientPacketHandler, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to start client handler: %v", err))
		os.Exit(1)
	}

	// Setup peer communication
	peerList := strings.Split(*peers, ",")
	if len(peerList) == 1 && peerList[0] == "" {
		peerList = nil
	}
	for i, peer := range peerList {
		if peer == *gossipAddr {
			size := len(peerList)
			peerList[i] = peerList[size-1]
			peerList = peerList[:size-1]
			break
		}
	}
	gossipInput, gossipPacketHandler := communication.NewGossipPacketDecoder(log)
	gossipWriter, err := communication.NewUDPListener(*gossipAddr, gossipPacketHandler, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to start gossip handler: %v", err))
		os.Exit(1)
	}
	peerPool := communication.NewPeerPool(peerList, gossipWriter, log)

	// Setup UI
	uiInterface, err := http.NewHttpUI("127.0.0.1:"+*guiPort, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to start http server: %v", err))
		os.Exit(1)
	}

	// Start gossiper
	if *simple {
		g := gossiper.NewSimpleGossiper(
			*name,
			*gossipAddr,
			peerPool,
			gossipInput,
			clientInput,
			log)
		g.Run()
	} else {
		g, err := gossiper.New(
			*name,
			*gossipAddr,
			peerPool,
			*antiEntropy,
			*routeTimer,
			gossipInput,
			clientInput,
			uiInterface,
			log)

		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}

		g.Run()
	}
}
