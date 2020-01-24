package chord

import (
	"context"
	"errors"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
)

// GetPredecessor (RPC) returns a pointer to the predecessor node
func (chordNode *ChordNode) GetPredecessor(ctx context.Context, e *empty.Empty) (*Node, error) {
	pred, doesExist := chordNode.getPredecessor()
	if !doesExist {
		return nil, errors.New("predecessor is nil")
	}
	return &pred, nil
}

// FIXME CheckPredecessor relies on use of ctx.Timeout
/*
 FindSuccessor finds the successor of id
 n.find_successor(id)
 	if id is_in (n, successor]
 		return successor;
 	else
 		n' = closest_preceding_node(id);
 		return n'.find_successor(id);
*/
func (chordNode *ChordNode) FindSuccessor(ctx context.Context, messageIDPtr *ID) (*Node, error) {
	log.Println("running FindSuccessor")
	if messageIDPtr == nil {
		return nil, errors.New("id must not be nil")
	}

	id := messageIDPtr.Id
	log.Printf("id : %d", id)
	n := chordNode.node
	succ, doesExist := chordNode.getSuccessor()
	if !doesExist {
		return nil, errors.New("successor does not exist")
	}

	if n.Id < id && id <= succ.Id {
		return &succ, nil
	} else {
		n0 := chordNode.ClosestPrecedingNode(nodeID(id))
		if n0.Id == n.Id {
			return &n, nil
		}
		return chordNode.stubFindSuccessor(ipAddr(n0.Ip), context.Background(), &ID{Id: id})
	}
}

// Notify(n0) checks whether n0 needs to be my predecessor
// Algorithm:
// n.notify(n0)
//	 if (predecessor is nil or n0 is_in (predecessor; n))
//	 predecessor = n0;
func (chordNode *ChordNode) Notify(ctx context.Context, n0 *Node) (*empty.Empty, error) {
	emptyPtr := &empty.Empty{}
	log.Println("running notify")
	if n0 == nil {
		return emptyPtr, errors.New("trying to notify nil node")
	}
	log.Printf("running notify with n0 = %v", n0)
	pred, doesExist := chordNode.getPredecessor()
	if !doesExist || (pred.Id < n0.Id && n0.Id < chordNode.node.Id) {
		chordNode.setPredecessor(n0)
	}

	return emptyPtr, nil
}
