// For RPC related methods of ChordNode
package chord

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)
import . "github.com/dosarudaniel/CS438_Project/services/chord_service"

// getStubFor either returns the stub for the known stub (ChordClient) for gRPC
// or creates a new one and memoizes it
func (chordNode *ChordNode) getStubFor(ip ipAddr) (ChordClient, error) {
	chordNode.stubsPool.RLock()
	stub, isPresent := chordNode.stubsPool.pool[ip]
	chordNode.stubsPool.RUnlock()
	if isPresent {
		return stub, nil
	}

	// dial, memoize, return
	conn, err := grpc.Dial(string(ip), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return stub, err
	}
	newStub := NewChordClient(conn)
	chordNode.stubsPool.Lock()
	chordNode.stubsPool.pool[ip] = newStub
	chordNode.stubsPool.Unlock()
	return newStub, nil
}

// GetPredecessor (RPC) returns a pointer to the predecessor node
func (chordNode *ChordNode) GetPredecessor(ctx context.Context, e *empty.Empty) (*Node, error) {
	chordNode.predecessor.RLock()
	defer chordNode.predecessor.RUnlock()
	return chordNode.predecessor.node, nil
}

// FIXME implement functions below
func (chordNode *ChordNode) FindSuccessor(ctx context.Context, in *ID) (*Node, error) {
	return &Node{Ip: "localhost:5000", Id: "f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb"}, nil
}

// FIXME implement functions below
func (chordNode *ChordNode) GetSuccessorsList(ctx context.Context, in *empty.Empty) (*Nodes, error) {
	chordNode.successorsList.RLock()
	defer chordNode.successorsList.RUnlock()
	return &Nodes{
		NodeArray: chordNode.successorsList.list,
	}, nil
}

// Notify(n0) checks whether n0 needs to be my predecessor
// Algorithm:
// n.notify(n0)
//	 if (predecessor is nil or n0 is_in (predecessor; n))
//	 predecessor = n0;
func (chordNode *ChordNode) Notify(ctx context.Context, n0 *Node) (*empty.Empty, error) {
	if n0 == nil {
		return &empty.Empty{}, errors.New("trying to notify nil node")
	}
	chordNode.predecessor.Lock()
	defer chordNode.predecessor.Unlock()
	predecessorID := chordNode.predecessor.node.Id
	n0ID := n0.Id
	nID := chordNode.node.Id
	if chordNode.predecessor.node == nil || (predecessorID < n0ID && n0ID < nID) {
		chordNode.predecessor.node = n0
	}
	return &empty.Empty{}, nil
}
