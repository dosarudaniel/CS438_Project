// For core Chord-related functionality
package chord

import (
	"context"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"github.com/golang/protobuf/ptypes/empty"
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

	// to search the local finger table for the highest predecessor of id
	ClosestPrecedingNode(ID) Node

	// to learn about newly joined nodes
	StabilizeDaemon()

	// to make sure finger table entries are correct (and up to date)
	FixFingersDaemon()

	// to clear the node's predecessor pointer if the predecessor has failed
	CheckPredecessorsDaemon()
}

// Chord node can be created using `NewChordNode`
// However, it must be followed
// either by creating own Chord network using `create`
// or by joining an existing network using `join`
// Both `create` and `join` must be followed by a call to `launchFingerTableDaemons`
type ChordNode struct {
	// constant to keep own IP and ID
	// Note: it must not be changed after initial definition in NewChordNode()
	node Node

	predecessor    nodeWithMux
	successorsList successorsListWithMux

	fingerTable fingerTableWithMux

	// to keep ChordClients, which allows avoiding creation of connections to other nodes every time
	// we want to communicate with them: connections are memoized upon creation
	stubsPool stubsPoolWithMux

	// for gRPC server functionality
	chordServer *grpc.Server

	// for gRPC file share server functionality
	fileShareServer *grpc.Server
}

// NewChordNode is a constructor for ChordNode struct
func NewChordNode(listener net.Listener) (*ChordNode, error) {
	chordNode := &ChordNode{}

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
		make([]*Node, 0),
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
	RegisterChordServer(chordNode.chordServer, chordNode)
	RegisterFileShareServiceServer(chordNode.chordServer, chordNode)
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
	chordNode.successorsList.list = []*Node{&chordNode.node}
	chordNode.successorsList.Unlock()
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
	chordNode.predecessor.node = nil
	chordNode.predecessor.Unlock()

	n0Stub, err := chordNode.getStubFor(ipAddr(n0.Ip))
	if err != nil {
		return err
	}
	succ, err := n0Stub.FindSuccessor(context.Background(), &ID{Id: chordNode.node.Id})
	if err != nil {
		return err
	}

	chordNode.fingerTable.Lock()
	chordNode.fingerTable.table = []Node{*succ}
	chordNode.fingerTable.Unlock()

	succStub, err := chordNode.getStubFor(ipAddr(succ.Ip))
	if err != nil {
		return err
	}
	nodesPtr, err := succStub.GetSuccessorsList(context.Background(), &empty.Empty{})
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

func (chordNode *ChordNode) String() string {
	outputString := "Node " + chordNode.node.Id + " ip: " + chordNode.node.Ip + "\n"

	outputString += "\t Predecesor: "
	chordNode.predecessor.RLock()
	if chordNode.predecessor.node != nil {
		outputString += chordNode.predecessor.node.Id + " ip: " + chordNode.predecessor.node.Ip + "\n"
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
	for i, node := range chordNode.fingerTable.table {
		outputString += fmt.Sprintf("\t\t [%d] %s %s\n", i, node.Id, node.Ip)
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
