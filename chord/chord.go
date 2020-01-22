package chord

import (
	"context"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"google.golang.org/grpc"
	"log"
)

// IChordNode defines all of the functions that Chord will have (which I know at this moment)
type IChordNode interface {
	ChordServer
	Create()
	Join(Node)
	ClosestPrecedingNode(ID) Node
	Notify(Node)
	StabilizeDaemon()
	FixFingersDaemon()
	CheckPredecessorsDaemon()
}

// Chord node can be created using `NewChordNode`
// However, it must be followed
// either by creating own Chord network using `create`
// or by joining an existing network using `join`
// Both `create` and `join` must be followed by a call to `launchFingerTableDaemons`
type ChordNode struct {
	node Node

	predecessor    nodeWithMux
	successorsList successorsListWithMux

	fingerTable fingerTableWithMux

	chordServer *grpc.Server
}

// Create creates a Chord ring
func (chordNode *ChordNode) Create() {
	// create a new Chord ring.
	// n.create()
	//   predecessor = nil;
	//	 successor = [n];
	chordNode.predecessor.Lock()
	chordNode.predecessor.node = nil
	chordNode.predecessor.Unlock()

	chordNode.successorsList.Lock()
	chordNode.successorsList.list = []Node{chordNode.node}
	chordNode.successorsList.Unlock()
}
