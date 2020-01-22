// For RPC related methods of ChordNode
package chord

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)
import . "github.com/dosarudaniel/CS438_Project/services/chord_service"

// getStubFor either returns the stub for the known stub
// or creates a new one and memoizes it
func (chordNode *ChordNode) getStubFor(ip ipAddr) (ChordClient, error) {
	chordNode.stubsPool.RLock()
	stub, isPresent := chordNode.stubsPool.pool[ip]
	chordNode.stubsPool.RUnlock()
	if isPresent {
		return stub, nil
	}

	// dial, memoize, return
	conn, err := grpc.Dial(string(ip), grpc.WithBlock())
	if err != nil {
		return stub, err
	}
	newStub := NewChordClient(conn)
	chordNode.stubsPool.Lock()
	chordNode.stubsPool.pool[ip] = newStub
	chordNode.stubsPool.Unlock()
	return newStub, nil
}

// FIXME implement functions below
func (chordNode *ChordNode) FindSuccessor(ctx context.Context, in *ID) (*Nodes, error) {
	return nil, nil
}
func (chordNode *ChordNode) Notify(ctx context.Context, in *Node) (*empty.Empty, error) {
	return nil, nil
}
