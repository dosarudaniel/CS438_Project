// For core Chord-related functionality
package chord

import (
	"context"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// IChordNode interface defines all of the functions that Chord will have (which I know at this moment)
type IChordNode interface {
	// interface created by the chord_service proto file for RPC methods
	ChordServer

	// to create own Chord ring (network)
	Create()

	// to join the Chord ring (network) knowing a single node already in the ring
	Join(Node) error

	// to search the local finger table for the highest predecessor of nodeID
	ClosestPrecedingNode(nodeID) Node

	// to learn about newly joined nodes
	StabilizeDaemon()

	// to make sure finger table entries are correct (and up to date)
	FixFingersDaemon()

	// to clear the node's predecessor pointer if the predecessor has failed
	CheckPredecessorDaemon()
}

// Chord node can be created using `NewChordNode`
// However, it must be followed
// either by creating own Chord network using `create`
// or by joining an existing network using `join`
// Both `create` and `join` must be followed by a call to `launchFingerTableDaemons`
type ChordNode struct {
	// constant; must not be changed after first initialization
	config ChordConfig

	// constant to keep own IP and ID
	// Note: it must not be changed after initialization in NewChordNode()
	node Node

	predecessor nodeWithMux

	fingerTable fingerTableWithMux

	// to keep ChordClients, which allows avoiding creation of connections to other nodes every time
	// we want to communicate with them: connections are memoized upon creation
	stubsPool stubsPoolWithMux

	// for gRPC server functionality
	chordServer *grpc.Server
}

// NewChordNode is a constructor for ChordNode struct
func NewChordNode(listener net.Listener, config ChordConfig) (*ChordNode, error) {
	chordNode := &ChordNode{}

	chordNode.config = config

	ip := listener.Addr().String()
	id, err := hashString(ip, int(config.NumOfBitsInID))
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
	chordNode.fingerTable = fingerTableWithMux{
		make(fingerTable, chordNode.config.NumOfBitsInID),
		sync.RWMutex{},
	}
	chordNode.stubsPool = stubsPoolWithMux{
		make(ipToStubMap),
		sync.RWMutex{},
	}

	chordNode.chordServer = grpc.NewServer()
	RegisterChordServer(chordNode.chordServer, chordNode)
	go chordNode.chordServer.Serve(listener)

	// TODO replace by a constant or config.fixFingerInterval
	go chordNode.RunAtInterval(StabilizeDaemon, 5)
	go chordNode.RunAtInterval(FixFingersDaemon(chordNode), 5)
	go chordNode.RunAtInterval(CheckPredecessorDaemon, 5)

	return chordNode, nil
}

func (chordNode *ChordNode) setSuccessor(nodePtr *Node) {
	chordNode.fingerTable.Lock()
	defer chordNode.fingerTable.Unlock()
	chordNode.fingerTable.table[0] = nodePtr
}

// returns Node, doesExist
func (chordNode *ChordNode) getSuccessor() (Node, bool) {
	chordNode.fingerTable.RLock()
	defer chordNode.fingerTable.RUnlock()
	if chordNode.fingerTable.table[0] == nil {
		return Node{}, false
	} else {
		return *chordNode.fingerTable.table[0], true
	}
}

func (chordNode *ChordNode) setPredecessor(nodePtr *Node) {
	chordNode.predecessor.Lock()
	defer chordNode.predecessor.Unlock()
	chordNode.predecessor.nodePtr = nodePtr
}

// returns Node, doesExist
func (chordNode *ChordNode) getPredecessor() (Node, bool) {
	chordNode.predecessor.RLock()
	defer chordNode.predecessor.RUnlock()
	if chordNode.predecessor.nodePtr == nil {
		return Node{}, false
	} else {
		return *chordNode.predecessor.nodePtr, true
	}
}

// Create creates a Chord ring
func (chordNode *ChordNode) Create() {
	// create a new Chord ring.
	// n.create()
	//   predecessor = nil;
	//   successor = n;
	chordNode.setPredecessor(nil)
	chordNode.setSuccessor(&chordNode.node)
}

// Join lets the ChordNode to join a Chord ring containing node n0.
// Algorithm:
// n.join(n0)
//	 predecessor = nil;
//	 successor = n0.find_successor(n);
func (chordNode *ChordNode) Join(n0 Node) error {
	var err error

	chordNode.setPredecessor(nil)

	succ, err := chordNode.stubFindSuccessor(ipAddr(n0.Ip), context.Background(), &ID{Id: chordNode.node.Id})
	fmt.Printf("Im joining succ: %v", succ)
	if err != nil {
		return err
	}

	chordNode.setSuccessor(succ)

	return nil
}

func (chordNode *ChordNode) ClosestPrecedingNode(nodeID nodeID) Node {
	id := string(nodeID)
	n := chordNode.node.Id
	chordNode.fingerTable.RLock()
	defer chordNode.fingerTable.RUnlock()
	fingerTable := chordNode.fingerTable.table
	for i := len(chordNode.fingerTable.table) - 1; i >= 0; i-- {
		finger := fingerTable[i]
		if finger == nil {
			continue
		}
		if n < finger.Id && finger.Id < id {
			return *finger
		}
	}
	return chordNode.node
}

func (chordNode *ChordNode) String() string {
	outputString := "Node " + chordNode.node.Id + " ip: " + chordNode.node.Ip + "\n"

	outputString += "\t Predecesor: "
	chordNode.predecessor.RLock()
	if chordNode.predecessor.nodePtr != nil {
		outputString += chordNode.predecessor.nodePtr.Id + " ip: " + chordNode.predecessor.nodePtr.Ip + "\n"
	} else {
		outputString += "nil\n"
	}
	chordNode.predecessor.RUnlock()

	outputString += "\t Finger table: \n"
	chordNode.fingerTable.RLock()
	for i, nodePtr := range chordNode.fingerTable.table {
		if nodePtr == nil {
			outputString += fmt.Sprintf("\t\t [%d] nil\n", i)
		} else {
			outputString += fmt.Sprintf("\t\t [%d] %s %s\n", i, nodePtr.Id, nodePtr.Ip)
		}
	}
	chordNode.fingerTable.RUnlock()

	outputString += "\t Connections: \n"
	chordNode.stubsPool.RLock()
	for ip := range chordNode.stubsPool.pool {
		outputString += "\t\t" + string(ip) + "\n"
	}
	chordNode.stubsPool.RUnlock()

	return outputString
}
