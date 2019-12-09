package communication

import (
	"fmt"
//========= REDACTED =========
	"net"
	"sync"
	"github.com/2_alt_hw2/Peerster/logger"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/dedis/protobuf"

//========= REDACTED =========
//========= REDACTED =========
)

// Peer represents a peer in the pool
type Peer struct {
	peerAddr string
	peerCh   chan<- []byte
	wg       *sync.WaitGroup
	log      *logger.Logger
}

// NewPeer creates a new Peer object
func NewPeer(peerAddr string, writer *UDPWriter, wg *sync.WaitGroup, log *logger.Logger) (*Peer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", peerAddr)
	if err != nil {
		return nil, err
	}

	peerChan := make(chan []byte, 100)
	p := Peer{
		peerAddr,
		peerChan,
		wg,
		log,
	}
	p.startRoutine(peerChan, udpAddr, writer)

	return &p, nil
}

// startSender starts a goroutine sending incoming packets from channel `peerChan` to `peer` through UDP
func (p *Peer) startRoutine(recvPeerCh <-chan []byte, udpAddr *net.UDPAddr, writer *UDPWriter) {
	go func() {
		p.wg.Add(1)
		p.log.Debug(fmt.Sprintf("Connection to %v established", p.peerAddr))

		for packet := range recvPeerCh {
			p.log.Trace(fmt.Sprintf("Sending packet to %v", p.peerAddr))
			_, err := writer.WriteTo(packet, udpAddr)
			if err != nil {
				p.log.Warn(fmt.Sprintf("Error while sending message to %v: %v", p.peerAddr, err))
			}
		}

		p.log.Debug(fmt.Sprintf("Connection to %v closed", p.peerAddr))
		p.wg.Done()
	}()
}

// SendPacket sends a packet to the peer
func (p *Peer) SendPacket(packet *types.GossipPacket) error {
	p.log.Trace(fmt.Sprintf("Sending packet: %v", packet))
	packetBytes, err := protobuf.Encode(packet)
	if err != nil {
		return fmt.Errorf("error while encoding packet: %v", err)
	}
	p.peerCh <- packetBytes
	return nil
}

// SendBytes sends bytes to the peer
func (p *Peer) SendBytes(packetBytes []byte) {
	p.peerCh <- packetBytes
}

// Close closes the connection to the peer
func (p *Peer) Close() {
	close(p.peerCh)
}

// String returns the string representation of the peer
func (p *Peer) String() string {
	return p.peerAddr
}
