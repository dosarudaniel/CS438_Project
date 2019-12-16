package gossiper

import (
	"fmt"
	"github.com/dosarudaniel/CS438_Project/types"
)

func (g *Gossiper) forwardPacket(packet types.GossipPacket, destination string) error {
	g.log.Trace(fmt.Sprintf("Forwarding packet (before decrement): %v", packet))

	if !packet.DecrementHop() {
		return fmt.Errorf("hop limit reached")
	}

	nextPeerAddr, err := g.dsdv.RouteTo(destination)
	if err != nil {
		return err
	}
	peer, err := g.peerPool.Find(nextPeerAddr)
	if err != nil {
		return err
	}
	err = peer.SendPacket(&packet)
	if err != nil {
		return err
	}
	return nil
}
