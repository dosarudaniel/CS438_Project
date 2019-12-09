package gossiper

import (
	"fmt"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/communication"
	"github.com/2_alt_hw2/Peerster/logger"
	"time"
)

// handleRumor handles incoming rumor messages
func (g *Gossiper) handleRumor(sender string, rumor *types.RumorMessage) {
	fmt.Println(rumor.String(sender))
	if !g.rumorRepository.Contains(rumor.Origin, rumor.ID) {
		g.rumorRepository.Insert(rumor)
		if g.dsdv.Upsert(rumor.Origin, rumor.ID, sender) && rumor.Text != "" {
			fmt.Printf("DSDV %v %v\n", rumor.Origin, sender)
		}
		randomPeer := g.peerPool.GetRandom()
		g.sendRumor(randomPeer, rumor)
	}

	senderPeer, err := g.peerPool.Find(sender)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Sender peer of received rumor not found: %v", err))
	}
	err = senderPeer.SendPacket(&types.GossipPacket{Status: &types.StatusPacket{Want: g.rumorRepository.GetStatus()}})
	if err != nil {
		g.log.Debug(fmt.Sprintf("\tStatus(ACK) packet not sent: %v", err))
	}
}

// handleRouteTimer handles Rumor timer ticks to announce our route to other peers
func (g *Gossiper) handleRouteTimer() {
	randomPeer := g.peerPool.GetRandom()
	if randomPeer == nil {
		g.log.Info("No route announcement with empty peer pool")
		return
	}
	g.log.Info(fmt.Sprintf("Route announcement to %v", randomPeer))
	g.sendNewRumor(randomPeer, "")
}

func (g *Gossiper) sendRumor(target *communication.Peer, rumor *types.RumorMessage) {
	if target == nil {
		return
	}

	g.pendingRumors.Insert(target.String(), rumor.Origin, rumor.ID)
	if rumor.Text != "" {
		fmt.Printf("MONGERING with %v\n", target)
	}
	g.log.Info(fmt.Sprintf("\tSending rumor %v", rumor))
	if g.log.Level <= logger.TraceLevel {
		pendingRumors := g.pendingRumors.AllPending()
		g.log.Trace(fmt.Sprintf("\tPending rumors(%v): %v", len(pendingRumors), pendingRumors))
	}
	err := target.SendPacket(&types.GossipPacket{Rumor: rumor})
	if err != nil {
		g.log.Debug(fmt.Sprintf("\tRumor not sent: %v", err))
	}

	// Timeout
	time.AfterFunc(time.Duration(10)*time.Second, func() {
		if g.pendingRumors.Delete(target.String(), rumor.Origin, rumor.ID) {
			g.log.Trace(fmt.Sprintf("Timeout for the message %v to %v", rumor, target))
			g.sendRumor(g.peerPool.GetRandom(), rumor)
		}
	})
}

func (g *Gossiper) sendNewRumor(target *communication.Peer, message string) {
	rumor := types.RumorMessage{
		Origin: g.name,
		ID:     g.rumorRepository.NextRumor(g.name),
		Text:   message,
	}
	g.rumorRepository.Insert(&rumor)
	g.sendRumor(target, &rumor)
}
