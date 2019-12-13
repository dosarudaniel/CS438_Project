package communication

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/dosarudaniel/CS438_Project/logger"
	"github.com/dosarudaniel/CS438_Project/types"
	"github.com/dedis/protobuf"
)

// PeerPool hides the unexported type to keep to properties private
type PeerPool struct {
	peers     map[string]*Peer
	peerSlice []*Peer // For efficient random access
	writer    *UDPWriter
	strRepr   *strings.Builder
	waitGroup *sync.WaitGroup
	log       *logger.Logger
}

// NewPeerPool creates a new broadcaster object
func NewPeerPool(peerAddrList []string, writer *UDPWriter, log *logger.Logger) *PeerPool {
	var strRepr strings.Builder
	strRepr.WriteString("PEERS ")

	p := PeerPool{
		make(map[string]*Peer, len(peerAddrList)),
		make([]*Peer, 0, len(peerAddrList)),
		writer,
		&strRepr,
		new(sync.WaitGroup),
		log,
	}

	for _, peerAddr := range peerAddrList {
		err := p.Insert(peerAddr)
		if err != nil {
			log.Warn(fmt.Sprintf("Unable to add peer %v: %v", peerAddr, err))
		}
	}

	return &p
}

// Insert method creates a new peer connection in the broadcaster
func (p *PeerPool) Insert(peerAddr string) error {
	_, found := p.peers[peerAddr]
	if found {
		return nil
	}

	peer, err := NewPeer(peerAddr, p.writer, p.waitGroup, p.log)
	if err != nil {
		return fmt.Errorf("unable to insert peer: %v", err)
	}

	p.peers[peerAddr] = peer
	p.peerSlice = append(p.peerSlice, peer)
	p.strRepr.WriteString(peerAddr + ",")
	return nil
}

// Contains indicates if peer is in the pool
func (p *PeerPool) Contains(peer string) bool {
	_, ok := p.peers[peer]
	return ok
}

// Find searches for a channel for a given peer
func (p *PeerPool) Find(peer string) (*Peer, error) {
	ch, ok := p.peers[peer]
	if !ok {
		return nil, fmt.Errorf("peer (%v) not found", peer)
	}
	return ch, nil
}

// List return the list of all peer name
func (p *PeerPool) List() []string {
	list := make([]string, 0, len(p.peerSlice))

	for _, peer := range p.peerSlice {
		list = append(list, peer.String())
	}

	return list
}

// GetRandom returns a random peer in the pool or nil if the pool is empty
func (p *PeerPool) GetRandom() *Peer {
	nbPeer := len(p.peerSlice)
	if nbPeer == 0 {
		p.log.Info("GetRandom peer from empty pool")
		return nil
	}
	return p.peerSlice[rand.Int()%nbPeer]
}

// Broadcast sends a packet to all connected peers
func (p *PeerPool) Broadcast(packet *types.GossipPacket, peerException string) error {
	packetBytes, err := protobuf.Encode(packet)
	if err != nil {
		return fmt.Errorf("error while encoding packet: %v", err)
	}

	p.log.Info(fmt.Sprintf("Broadcasting %v", packet))
	for peerAddr, peer := range p.peers {
		if peerAddr != peerException {
			p.log.Trace(fmt.Sprintf("\tSending to channel of %v", peerAddr))
			peer.SendBytes(packetBytes)
		}
	}
	return nil
}

// Close closes all channels to peer senders, therefore stopping the associated goroutines
func (p *PeerPool) Close() {
	for _, peer := range p.peers {
		peer.Close()
	}
	p.waitGroup.Wait()
}

// String return the string representation of a broadcaster
func (p *PeerPool) String() string {
	return strings.TrimSuffix(p.strRepr.String(), ",") // Cut the last coma
}
