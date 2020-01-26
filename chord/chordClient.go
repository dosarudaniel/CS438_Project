// For RPC related methods of ChordNode
package chord

import (
	"context"
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

func (chordNode *ChordNode) stubFindSuccessor(ip ipAddr, ctx context.Context, id *ID) (*Node, error) {
	stub, err := chordNode.getStubFor(ip)
	if err != nil {
		return nil, err
	}

	return stub.FindSuccessor(ctx, id)
}

func (chordNode *ChordNode) stubNotify(ip ipAddr, ctx context.Context, node *Node) error {
	stub, err := chordNode.getStubFor(ip)
	if err != nil {
		return err
	}

	_, err = stub.Notify(ctx, node)

	return err
}

func (chordNode *ChordNode) stubGetPredecessor(ip ipAddr, ctx context.Context) (*Node, error) {
	stub, err := chordNode.getStubFor(ip)
	if err != nil {
		return nil, err
	}

	return stub.GetPredecessor(ctx, &empty.Empty{})
}
