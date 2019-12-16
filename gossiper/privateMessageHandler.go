package gossiper

import (
	"fmt"
	"github.com/dosarudaniel/CS438_Project/types"
)

func (g *Gossiper) handlePrivateMessage(sender string, pm *types.PrivateMessage) {
	if pm.Destination == g.name {
		g.log.Info(fmt.Sprintf("Got a private message for me: %v", pm))
		fmt.Printf("PRIVATE origin %v hop-limit %v contents %v\n", pm.Origin, pm.HopLimit, pm.Text)
		g.privateMsgRepository.Insert(pm)
	} else {
		g.log.Debug(fmt.Sprintf("Got a private message, forwarding it (before hop decrement): %v", pm))
		err := g.forwardPacket(types.GossipPacket{Private: pm}, pm.Destination)
		if err != nil {
			g.log.Warn(fmt.Sprintf("Unable to forward private message: %v", err))
		}
	}
}

func (g *Gossiper) sendPrivateMessage(message types.Message) {
	pm := types.PrivateMessage{
		Origin:      g.name,
		ID:          0,
		Text:        message.Text,
		Destination: *message.Destination,
		HopLimit:    10,
	}

	err := g.forwardPacket(types.GossipPacket{Private: &pm}, pm.Destination)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable to send private message %v: %v", pm, err))
		return
	}

	g.privateMsgRepository.Insert(&pm)
}
