package gossiper

import (
	//"fmt"
	"github.com/dosarudaniel/CS438_Project/types"
	//"github.com/dosarudaniel/CS438_Project/file_sharing"
)

func (g *Gossiper) handleSearchRequest(request *types.SearchRequest) {
	// if request.Destination == g.name {
	// 	g.log.Info(fmt.Sprintf("Got a data request for me: %v", request))
	// 	reply, err := g.processDataRequest(request)
	// 	if err != nil {
	// 		g.log.Warn(fmt.Sprintf("Unable to process data request %v: %v", request, err))
	// 		return
	// 	}
	// 	g.log.Trace(fmt.Sprintf("\tReply data: %v", reply.Data))
	// 	err = g.forwardPacket(types.GossipPacket{DataReply: reply}, reply.Destination)
	// 	if err != nil {
	// 		g.log.Warn(fmt.Sprintf("Unable to forward data request: %v", err))
	// 		return
	// 	}
	// } else {
	// 	g.log.Info(fmt.Sprintf("Got a data request, forwarding it (pre-hop decrement): %v", request))
	// 	err := g.forwardPacket(types.GossipPacket{DataRequest: request}, request.Destination)
	// 	if err != nil {
	// 		g.log.Warn(fmt.Sprintf("Unable to forward data request: %v", err))
	// 		return
	// 	}
	// }
}

func (g *Gossiper) handleSearchReply(reply *types.SearchReply) {
	// if reply.Destination == g.name {
	// 	g.log.Info(fmt.Sprintf("Got a data reply for me: %v", reply))
	// 	err := g.processDataReply(reply)
	// 	if err != nil {
	// 		g.log.Warn(fmt.Sprintf("Unable to process data reply %v: %v", reply, err))
	// 	}
	// } else {
	// 	g.log.Info(fmt.Sprintf("Got a data reply, forwarding it (pre-hop decrement): %v", reply))
	// 	err := g.forwardPacket(types.GossipPacket{DataReply: reply}, reply.Destination)
	// 	if err != nil {
	// 		g.log.Warn(fmt.Sprintf("Unable to forward data reply: %v", err))
	// 	}
	// }
}

// func (g *Gossiper) processDataRequest(request *types.DataRequest) (*types.DataReply, error) {
// 	hash, err := file_sharing.BytesToHash(request.HashValue)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	chunk := g.sharedFilesRepository.GetChunk(hash)
//
// 	return &types.DataReply{
// 		Origin:      g.name,
// 		Destination: request.Origin,
// 		HopLimit:    10,
// 		HashValue:   request.HashValue,
// 		Data:        chunk,
// 	}, nil
// }
//
// func (g *Gossiper) processDataReply(reply *types.DataReply) error {
// 	hash, err := file_sharing.BytesToHash(reply.HashValue)
// 	if err != nil {
// 		return err
// 	}
// 	return g.downloadManager.Deliver(reply.Origin, hash, reply.Data)
// }
