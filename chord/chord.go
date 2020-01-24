// For core Chord-related functionality
package chord

import (
	"context"
	"errors"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
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

	predecessor    nodeWithMux
	successorsList successorsListWithMux

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
		make([]*Node, chordNode.config.LenOfSuccList),
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

	return chordNode, nil
}

// Create creates a Chord ring
func (chordNode *ChordNode) Create() {
	// create a new Chord ring.
	// n.create()
	//   predecessor = nil;
	//	 successorsList[0] = n;
	//   finger[0] = n
	chordNode.predecessor.Lock()
	chordNode.predecessor.nodePtr = nil
	chordNode.predecessor.Unlock()

	chordNode.successorsList.Lock()
	chordNode.successorsList.list[0] = &chordNode.node
	chordNode.successorsList.Unlock()

	chordNode.fingerTable.Lock()
	chordNode.fingerTable.table[0] = &chordNode.node
	chordNode.fingerTable.Unlock()
}

// Join lets the ChordNode to join a Chord ring containing node n0.
// Algorithm:
// n.join(n0)
//   predecessor = nil;
//   succ := n0.find_successor(n);
//	 fingerTable[0] = succ
//   succList := succ.getSuccessorsList();
//   successorsList = succ :: succList[:len(succList) - 1]
func (chordNode *ChordNode) Join(n0 Node) error {
	var err error

	chordNode.predecessor.Lock()
	chordNode.predecessor.nodePtr = nil
	chordNode.predecessor.Unlock()

	succ, err := chordNode.stubFindSuccessor(ipAddr(n0.Ip), context.Background(), &ID{Id: chordNode.node.Id})
	if err != nil {
		return err
	}

	chordNode.fingerTable.Lock()
	chordNode.fingerTable.table[0] = succ
	chordNode.fingerTable.Unlock()

	nodesPtr, err := chordNode.stubGetSuccessorsList(ipAddr(succ.Ip), context.Background())
	switch {
	case err != nil:
		return err
	case nodesPtr == nil:
		return fmt.Errorf("successors list from %s is nil", succ.Ip)
	}
	succList := nodesPtr.NodeArray

	chordNode.successorsList.Lock()
	defer chordNode.successorsList.Unlock()
	chordNode.successorsList.list = append([]*Node{succ}, succList[:len(succList)-1]...)

	return nil
}

// CheckPredecessorDaemon checks whether the node's predecessor has failed
// n.check_predecessor()
//	 if (predecessor has failed) <- in our case responds to a FindSuccessor rpc call within 3 seconds
//	   predecessor = nil;
func (chordNode *ChordNode) CheckPredecessorDaemon() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	chordNode.predecessor.RLock()
	predPtr := chordNode.predecessor.nodePtr
	chordNode.predecessor.RUnlock()
	if predPtr == nil {
		return
	}
	_, err := chordNode.stubFindSuccessor(ipAddr(predPtr.Ip), ctx, &ID{Id: chordNode.node.Ip})
	if errors.Is(err, context.DeadlineExceeded) {
		chordNode.predecessor.Lock()
		chordNode.predecessor.nodePtr = nil
		chordNode.predecessor.Unlock()
	}
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

	outputString += "\t Successors list: \n"
	chordNode.successorsList.RLock()
	for _, successor := range chordNode.successorsList.list {
		if successor != nil {
			outputString += "\t\t" + successor.Id + " " + successor.Ip + "\n"
		} else {
			outputString += "\t\t" + "nil" + "\n"
		}
	}
	chordNode.successorsList.RUnlock()

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
