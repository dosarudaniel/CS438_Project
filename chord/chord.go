// For core Chord-related functionality
package chord

import (
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// IChordNode interface defines all of the functions that Chord will have (which I know at this moment)
type IChordNode interface {
	ChordServer
	Create()
	Join(Node)
	ClosestPrecedingNode(ID) Node
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

	stubsPool stubsPoolWithMux

	chordServer *grpc.Server
}

func NewChordNode(listener net.Listener) (ChordNode, error) {
	chordNode := ChordNode{}

	ip := listener.Addr().String()
	id, err := hashString(ip)
	if err != nil {
		return chordNode, err
	}
	chordNode.node = Node{
		Id: id,
		Ip: ip,
	}
	chordNode.predecessor = nodeWithMux{
		nil,
		sync.RWMutex{},
	}
	chordNode.successorsList = successorsListWithMux{
		make([]Node, 0),
		sync.RWMutex{},
	}
	chordNode.fingerTable = fingerTableWithMux{
		make(fingerTable, 0),
		sync.RWMutex{},
	}
	chordNode.stubsPool = stubsPoolWithMux{
		make(ipToStubMap),
		sync.RWMutex{},
	}

	chordNode.chordServer = grpc.NewServer()
	RegisterChordServer(chordNode.chordServer, &chordNode)
	go chordNode.chordServer.Serve(listener)

	return chordNode, nil
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
