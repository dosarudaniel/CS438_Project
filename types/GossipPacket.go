package types

import "fmt"

// GossipPacket represents the only type of Peerster packet for the gossiper
type GossipPacket struct {
	Simple      	*SimpleMessage
	Rumor       	*RumorMessage
	Status      	*StatusPacket
	Private     	*PrivateMessage
	DataRequest 	*DataRequest
	DataReply   	*DataReply
	GossipAddr  	*string   			// This field caused me problems from HW3
	SearchRequest 	*SearchRequest
	SearchReply 	*SearchReply
}

// DecrementHop removes one from the hop limit and returns true if the packet is still valid, false otherwise
// Does nothing for packets without hop value
func (g GossipPacket) DecrementHop() bool {
	if g.Private != nil {
		if g.Private.HopLimit > 0 {
			g.Private.HopLimit -= 1
		} else {
			return false
		}
	} else if g.DataRequest != nil {
		if g.DataRequest.HopLimit > 0 {
			g.DataRequest.HopLimit -= 1
		} else {
			return false
		}
	} else if g.DataReply != nil {
		if g.DataReply.HopLimit > 0 {
			g.DataReply.HopLimit -= 1
		} else {
			return false
		}
	} else if g.SearchReply != nil {
		if g.SearchReply.HopLimit > 0 {
			g.SearchReply.HopLimit -= 1
		} else {
			return false
		}
	}

	return true
}

// DecrementBudget removes one from the budget limit and returns true if the packet is still valid, false otherwise
// Does nothing for packets without budget value
func (g GossipPacket) DecrementBudget() bool {
	 if g.SearchRequest != nil {
		if g.SearchRequest.Budget > 0 {
			g.SearchRequest.Budget -= 1
		} else {
			return false
		}
	}

	return true
}

// String returns the string representation of a Gossip packet
func (g GossipPacket) String() string {
	var inner string
	if g.Simple != nil {
		inner = g.Simple.String()
	} else if g.Rumor != nil {
		if g.GossipAddr == nil {
			inner = g.Rumor.String("nil")
		} else {
			inner = g.Rumor.String(*g.GossipAddr)
		}
	} else if g.Status != nil {
		inner = g.Status.String()
	} else if g.Private != nil {
		inner = g.Private.String()
	} else if g.DataRequest != nil {
		inner = g.DataRequest.String()
	} else if g.DataReply != nil {
		inner = g.DataReply.String()
	} else if g.SearchReply != nil {
		inner = g.SearchReply.String()
	} else if g.SearchRequest != nil {
		inner = g.SearchRequest.String()
	} else {
		inner = "____Nothing____"
	}
	return fmt.Sprintf("GossipPacket{%v}", inner)
}
