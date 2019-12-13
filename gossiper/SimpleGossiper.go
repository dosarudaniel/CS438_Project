package gossiper

import (
	"fmt"
	"github.com/dosarudaniel/CS438_Project/types"
	"github.com/dosarudaniel/CS438_Project/communication"
	"github.com/dosarudaniel/CS438_Project/logger"
)

// Gossiper represents the main logic of CS438_Project
// It handles incoming messages
type SimpleGossiper struct {
	name        string
	address     string
	broadcaster *communication.PeerPool

	gossipInput <-chan types.GossipPacket
	clientInput <-chan types.Message
	log         *logger.Logger
}

// NewSimpleGossiper creates a new gossiper logic
func NewSimpleGossiper(name string, address string, peerPool *communication.PeerPool,
	gossipInput <-chan types.GossipPacket,
	clientInput <-chan types.Message,
	log *logger.Logger) *SimpleGossiper {
	return &SimpleGossiper{
		name,
		address,
		peerPool,
		gossipInput,
		clientInput,
		log,
	}
}

// Run starts the gossiper process. It is a blocking method
func (g *SimpleGossiper) Run() {
	g.log.Info("Gossiper running")
	for {
		select {
		case gossip := <-g.gossipInput:
			g.handleGossip(gossip)
		case cmd := <-g.clientInput:
			g.handleClientCmd(cmd)
		}
	}
}

// handleGossip handles incoming gossip messages
func (g *SimpleGossiper) handleGossip(gossip types.GossipPacket) {
	fmt.Println(gossip.Simple.String())
	fromPeer := gossip.Simple.RelayPeerAddr
	err := g.broadcaster.Insert(fromPeer)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable to add peer %v: %v", fromPeer, err))
	}
	fmt.Println(g.broadcaster.String())
	gossip.Simple.RelayPeerAddr = g.address
	err = g.broadcaster.Broadcast(&gossip, fromPeer)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable to brodcast message %v: %v", gossip, err))
	}
}

// handleClientCmd handles incoming client commands
func (g *SimpleGossiper) handleClientCmd(cmd types.Message) {
	fmt.Println(cmd.String())
	fmt.Println(g.broadcaster.String())
	simpleMessage := types.SimpleMessage{
		OriginalName:  g.name,
		RelayPeerAddr: g.address,
		Contents:      cmd.Text,
	}
	gossip := types.GossipPacket{Simple: &simpleMessage}
	err := g.broadcaster.Broadcast(&gossip, "")
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable to brodcast message %v: %v", gossip, err))
	}
}
