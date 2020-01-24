package chord

import (
	"context"
	"errors"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"github.com/golang/protobuf/ptypes/empty"
)

// GetPredecessor (RPC) returns a pointer to the predecessor node
func (chordNode *ChordNode) GetPredecessor(ctx context.Context, e *empty.Empty) (*Node, error) {
	pred, doesExist := chordNode.getPredecessor()
	if !doesExist {
		return nil, errors.New("predecessor is nil")
	}
	return &pred, nil
}

// FIXME implement functions below
// FIXME CheckPredecessor relies on use of ctx.Timeout
func (chordNode *ChordNode) FindSuccessor(ctx context.Context, in *ID) (*Node, error) {
	return &Node{Ip: "localhost:5000", Id: "f98eeff24e2fced1a1336182a3e8775326262914cc4087066d9346431795ccdb"}, nil
}

// Notify(n0) checks whether n0 needs to be my predecessor
// Algorithm:
// n.notify(n0)
//	 if (predecessor is nil or n0 is_in (predecessor; n))
//	 predecessor = n0;
func (chordNode *ChordNode) Notify(ctx context.Context, n0 *Node) (*empty.Empty, error) {
	emptyPtr := &empty.Empty{}

	if n0 == nil {
		return emptyPtr, errors.New("trying to notify nil node")
	}
	pred, doesExist := chordNode.getPredecessor()
	if !doesExist || (pred.Id < n0.Id && n0.Id < chordNode.node.Id) {
		chordNode.setPredecessor(n0)
	}

	return emptyPtr, nil
}
