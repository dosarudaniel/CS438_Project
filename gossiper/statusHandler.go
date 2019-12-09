package gossiper

import (
	"fmt"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/logger"
	"math/rand"
)

// handleStatus handles incoming status messages
func (g *Gossiper) handleStatus(sender string, status *types.StatusPacket) {
	g.log.Debug(fmt.Sprintf("STATUS from %v %v\n", sender, status.String()))

	senderPeer, err := g.peerPool.Find(sender)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable to find sender %v in repository: %v", sender, err))
		return
	}

	// Searching for unknown origins in sender's status
	for _, statMsg := range status.Want {
		g.rumorRepository.AddOrigin(statMsg.Identifier)
	}

	myStatus := g.rumorRepository.GetStatus()
	rumor, peerHasInterestingRumors := g.rumorRepository.FirstUnknown(status.Want)
	if rumor != nil {
		g.log.Info(fmt.Sprintf("\tThere are rumors to send to %v. Sending: %v\tOWN-STATUS %v", sender, rumor, myStatus))
		g.sendRumor(senderPeer, rumor)
		return
	}

	// Searching for rumor to ask
	if peerHasInterestingRumors {
		g.log.Info(fmt.Sprintf("\t%v is up to date but has interesting rumors for us. Sending our status:", sender))
		g.log.Debug(fmt.Sprintf("\tSending our status: %v", myStatus))
		err := senderPeer.SendPacket(&types.GossipPacket{Status: &types.StatusPacket{Want: myStatus}})
		if err != nil {
			g.log.Debug(fmt.Sprintf("\tStatus(ASK for rumor) not sent: %v", err))
		}
	} else {
		g.log.Debug(fmt.Sprintf("IN SYNC WITH %v\n", sender))

		if g.log.Level <= logger.DebugLevel {
			pendingRumors := g.pendingRumors.Pending(sender)
			g.log.Debug(fmt.Sprintf("\tPending rumors with %v (%v): %v", sender, len(pendingRumors), pendingRumors))
		}
		ackedRumors := g.pendingRumors.Ack(sender, status.Want)
		g.log.Debug(fmt.Sprintf("\tAcked %v rumors: %v", len(ackedRumors), ackedRumors))

		if g.log.Level <= logger.DebugLevel {
			pendingRumors := g.pendingRumors.Pending(sender)
			g.log.Debug(fmt.Sprintf("\tPending rumors(%v): %v", len(pendingRumors), pendingRumors))
		}

		for _, ackedRumor := range ackedRumors {
			if (rand.Int() % 2) == 0 {
				randomPeer := g.peerPool.GetRandom()
				g.log.Info(fmt.Sprintf("FLIPPED COIN sending rumor to %v", randomPeer))
				rumor, err := g.rumorRepository.Find(ackedRumor.Identifier, ackedRumor.NextID)
				if err == nil {
					g.sendRumor(randomPeer, rumor)
				} else {
					g.log.Warn(fmt.Sprintf("Unable to find acked rumor in repository: %v", err))
				}
			} else {
				g.log.Info(fmt.Sprintf("Coin flipped ; got tails: Not sending %v", ackedRumor))
			}
		}
	}

	g.log.Debug(fmt.Sprintf(g.peerPool.String()))
}

// handleAntiEntropy handles Anti Entropy ticks and send a Status message to a random peer
func (g *Gossiper) handleAntiEntropy() {
	randomPeer := g.peerPool.GetRandom()
	if randomPeer == nil {
		g.log.Info("No anti-entropy with empty peer pool")
		return
	}
	g.log.Info(fmt.Sprintf("ANTI-ENTROPY with %v", randomPeer))
	status := g.rumorRepository.GetStatus()
	err := randomPeer.SendPacket(&types.GossipPacket{Status: &types.StatusPacket{Want: status}})
	if err != nil {
		g.log.Debug(fmt.Sprintf("\tStatus(anti entropy) not sent: %v", err))
	}
}
