package chord

import (
	"context"
	"errors"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"github.com/golang/protobuf/ptypes/empty"
)

// GetPredecessor (RPC) returns a pointer to the predecessor node
func (chordNode *ChordNode) GetPredecessor(ctx context.Context, e *empty.Empty) (*Node, error) {
	chordNode.predecessor.RLock()
	defer chordNode.predecessor.RUnlock()
	return chordNode.predecessor.nodePtr, nil
}

// FIXME implement functions below
func (chordNode *ChordNode) FindSuccessor(ctx context.Context, in *ID) (*Node, error) {
	return &Node{Ip: "localhost:5000", Id: "f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb"}, nil
}

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
	predecessorID := chordNode.predecessor.nodePtr.Id
	n0ID := n0.Id
	nID := chordNode.node.Id
	if chordNode.predecessor.nodePtr == nil || (predecessorID < n0ID && n0ID < nID) {
		chordNode.predecessor.nodePtr = n0
	}
	return &empty.Empty{}, nil
}
