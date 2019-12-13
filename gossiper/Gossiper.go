package gossiper

import (
	"time"
	"github.com/dosarudaniel/CS438_Project/logger"

	"github.com/dosarudaniel/CS438_Project/types"
	"github.com/dosarudaniel/CS438_Project/messaging"
	"github.com/dosarudaniel/CS438_Project/communication"
	"github.com/dosarudaniel/CS438_Project/file_sharing"
	"github.com/dosarudaniel/CS438_Project/ui"
)

// Gossiper represents the main logic of CS438_Project
// It handles incoming messages
type Gossiper struct {
	settings types.Settings
	name     string
	address  string

	peerPool              *communication.PeerPool
	rumorRepository       *messaging.RumorRepository
	privateMsgRepository  *messaging.PrivateMessageRepository
	pendingRumors         *messaging.PeerStatusRepository
	dsdv                  *communication.DSDV
	sharedFilesRepository *file_sharing.SharedFilesRepository
	downloadManager       *file_sharing.DownloadManager

	gossipInput      <-chan types.GossipPacket
	clientInput      <-chan types.Message
	dataRequestInput <-chan *types.DataRequest
	antiEntropy      <-chan time.Time
	routeTimer       <-chan time.Time

	uiInterface ui.BackendInterface
	log         *logger.Logger
}

// New creates a new gossiper logic
func New(name string, address string, peerPool *communication.PeerPool, antiEntropy int, routeTimer int,
	gossipInput <-chan types.GossipPacket,
	clientInput <-chan types.Message,
	uiInterface ui.BackendInterface,
	log *logger.Logger) (*Gossiper, error) {

	rumorRepository := messaging.NewRumorRepository(log)
	rumorRepository.AddOrigin(name)

	var antiEntropyTicker <-chan time.Time = nil
	if antiEntropy > 0 {
		antiEntropyTicker = time.NewTicker(time.Duration(antiEntropy) * time.Second).C
	}

	var rumorTicker <-chan time.Time = nil
	if routeTimer > 0 {
		rumorTicker = time.NewTicker(time.Duration(routeTimer) * time.Second).C
	}

	sharedPath, err := file_sharing.ExecutableRelativePath("_SharedFiles")
	if err != nil {
		return nil, err
	}

	sharedFilesRepository, err := file_sharing.NewSharedFilesRepository(sharedPath, log)
	if err != nil {
		return nil, err
	}

	downloadPath, err := file_sharing.ExecutableRelativePath("_Downloads")
	if err != nil {
		return nil, err
	}

	downloadManager, dataRequestInput, err := file_sharing.NewDownloadManager(downloadPath, name, log)
	if err != nil {
		return nil, err
	}

	return &Gossiper{
		types.Settings{
			GossipAddr:  address,
			Name:        name,
			AntiEntropy: antiEntropy,
			RouteTimer:  routeTimer,
		},
		name,
		address,
		peerPool,
		rumorRepository,
		messaging.NewPrivateMessageRepository(),
		messaging.NewPeerStatusRepository(),
		communication.NewDSDV(),
		sharedFilesRepository,
		downloadManager,
		gossipInput,
		clientInput,
		dataRequestInput,
		antiEntropyTicker,
		rumorTicker,
		uiInterface,
		log,
	}, nil
}

// Run starts the gossiper process. It is a blocking method
func (g *Gossiper) Run() {
	g.log.Info("Gossiper running")

	// Start-up Route Rumor
	g.handleRouteTimer()
	for {
		select {
		case gossip := <-g.gossipInput:
			g.handleGossipPacket(gossip)
		case cmd := <-g.clientInput:
			g.handleClientCmd(cmd)
		case _ = <-g.antiEntropy:
			g.handleAntiEntropy()
		case _ = <-g.routeTimer:
			g.handleRouteTimer()
		case uiRequest := <-g.uiInterface.RequestChannel:
			g.handleUIRequest(uiRequest)
		case dataRequest := <-g.dataRequestInput:
			g.handleDataRequest(dataRequest)
		}
	}
}

// handleGossipPacket handles incoming gossip packets
func (g *Gossiper) handleGossipPacket(gossip types.GossipPacket) {
	addr := *gossip.GossipAddr
	if addr != g.address {
		err := g.peerPool.Insert(addr)
		if err != nil {
			g.log.Warn(err.Error())
		}
	}

	if gossip.Rumor != nil {
		g.handleRumor(addr, gossip.Rumor)
	} else if gossip.Status != nil {
		g.handleStatus(addr, gossip.Status)
	} else if gossip.Private != nil {
		g.handlePrivateMessage(addr, gossip.Private)
	} else if gossip.DataRequest != nil {
		g.handleDataRequest(gossip.DataRequest)
	} else if gossip.DataReply != nil {
		g.handleDataReply(gossip.DataReply)
	} else if gossip.SearchRequest != nil {
		g.handleSearchRequest(gossip.SearchRequest)
	} else if gossip.SearchReply != nil {
		g.handleSearchReply(gossip.SearchReply)
	} else {
		g.log.Warn("Erroneous GossipPacket discarded")
	}
}
