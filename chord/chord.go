// For core Chord-related functionality
package chord

import (
	"context"
	"errors"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	. "github.com/dosarudaniel/CS438_Project/services/client_service"
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var log = logrus.New()

// IChordNode interface defines all of the functions that Chord will have (which I know at this moment)
type IChordNode interface {
	// interface created by the chord_service proto file for RPC methods
	ChordServer

	// to create own Chord ring (network)
	Create()

	// to join the Chord ring (network) knowing a single node already in the ring
	Join(Node) error

	// to search the local finger table for the highest predecessor of nodeID
	ClosestPrecedingFinger(nodeID) Node

	// to learn about newly joined nodes
	StabilizeDaemon()

	// to make sure finger table entries are correct (and up to date)
	FixFingersDaemon()

	// to clear the node's predecessor pointer if the predecessor has failed
	CheckPredecessorDaemon()
}

// ChordNode can be created using `NewChordNode`
// However, it must be followed
// either by creating own Chord network using `create`
// or by joining an existing network using `join`
type ChordNode struct {
	// constant; must not be changed after first initialization
	config ChordConfig

	// constant to keep own IP and ID
	// Note: it must not be changed after initialization in NewChordNode()
	node Node

	predecessor nodeWithMux

	fingerTable fingerTableWithMux

	hashTable HashTableWithMux

	// to keep ChordClients, which allows avoiding creation of connections to other nodes every time
	// we want to communicate with them: connections are memoized upon creation
	stubsPool stubsPoolWithMux

	// for gRPC server functionality
	grpcServer *grpc.Server
}

// NewChordNode is a constructor for ChordNode struct
func NewChordNode(listener net.Listener, config ChordConfig, verbose bool) (*ChordNode, error) {
	if verbose {
		log.SetLevel(logrus.TraceLevel)
	}

	chordNode := &ChordNode{}

	chordNode.config = config

	ip := listener.Addr().String()
	id, err := chordNode.hashString(ip)
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
	chordNode.hashTable = NewHashTable()
	chordNode.stubsPool = stubsPoolWithMux{
		make(ipToStubMap),
		sync.RWMutex{},
	}

	chordNode.grpcServer = grpc.NewServer()
	RegisterChordServer(chordNode.grpcServer, chordNode)
	RegisterFileShareServiceServer(chordNode.grpcServer, chordNode)
	RegisterClientServiceServer(chordNode.grpcServer, chordNode)
	go chordNode.grpcServer.Serve(listener)

	// TODO replace by a constant or config.fixFingerInterval
	go chordNode.RunAtInterval(StabilizeDaemon, 1)
	go chordNode.RunAtInterval(FixFingersDaemon(chordNode), 1)
	go chordNode.RunAtInterval(CheckPredecessorDaemon, 1)

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
	switch {
	case err != nil:
		return err
	case succ == nil:
		return &nilSuccessor{}
	case chordNode.node.Id == succ.Id:
		return errors.New("node with such ID already exists in the Chord ring")
	}

	chordNode.setSuccessor(succ)

	return nil
}

func (chordNode *ChordNode) ClosestPrecedingFinger(nodeID nodeID) Node {
	id := string(nodeID)
	chordNode.fingerTable.RLock()
	defer chordNode.fingerTable.RUnlock()
	fingerTable := chordNode.fingerTable.table
	for i := len(chordNode.fingerTable.table) - 1; i >= 0; i-- {
		finger := fingerTable[i]
		if finger == nil {
			continue
		}
		if isBetweenTwoNodesExclusive(chordNode.node.Id, finger.Id, id) {
			return *finger
		}
	}
	return chordNode.node
}

// lookup takes in a key, returns the ip address of the node that should store that chord pair
// return err is failure
func (chordNode *ChordNode) lookup(key string) (ipAddr, error) {
	var err error

	hashedKey, err := chordNode.hashString(key)
	if err != nil {
		return "", err
	}

	succ, err := chordNode.FindSuccessor(context.Background(), &ID{Id: hashedKey})
	if err != nil {
		return "", err
	}

	return ipAddr(succ.Ip), nil
}

// PutInDHT stores the key in the Chord ring (some other node who's responsible for that key)
func (chordNode *ChordNode) PutInDHT(keyword, filename string) error {
	ip, err := chordNode.lookup(keyword)
	if err != nil {
		return err
	}

	if string(ip) == chordNode.node.Ip {
		return chordNode.hashTable.PutOrAppendOne(keyword, &FileRecord{
			Filename: filename,
			OwnerIp:  chordNode.node.Ip,
		})
	}

	return chordNode.stubPut(ip, context.Background(), keyword, filename, chordNode.node.Ip)
}

// FindInDHT finds the key from the Chord ring, i.e., gets it from other node
func (chordNode *ChordNode) FindInDHT(keyword string) ([]*FileRecord, error) {
	emptyFileRecords := make([]*FileRecord, 0)

	ip, err := chordNode.lookup(keyword) // ip of the node to store keyword
	if err != nil {
		return emptyFileRecords, err
	}

	if string(ip) == chordNode.node.Ip {
		chordNode.hashTable.RLock()
		fileRecords, doesExist := chordNode.hashTable.table[keyword]
		chordNode.hashTable.RUnlock()

		if doesExist {
			return fileRecords, nil
		} else {
			return emptyFileRecords, errors.New("keyword does not exist")
		}
	}

	val, err := chordNode.stubGet(ip, context.Background(), keyword)
	if err != nil {
		return emptyFileRecords, err
	}

	return val.FileRecords, nil
}

func (chordNode *ChordNode) String() string {
	outputString := "\nNODE INFO: "
	outputString += fmt.Sprintf("ID (string %s) (big.int %s) IP %s\n",
		chordNode.node.Id, idToBigIntString(chordNode.node.Id), chordNode.node.Ip)

	outputString += "\t Predecesor: "
	chordNode.predecessor.RLock()
	if chordNode.predecessor.nodePtr != nil {
		outputString +=
			fmt.Sprintf("ID (string %s) (big.int %s) IP %s\n",
				chordNode.predecessor.nodePtr.Id,
				idToBigIntString(chordNode.predecessor.nodePtr.Id),
				chordNode.predecessor.nodePtr.Ip)
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
			fingerStart, _ := getFingerStart(chordNode.node.Id, i, chordNode.config.NumOfBitsInID)
			outputString += fmt.Sprintf("\t\t [%d] [start: %s] (string %s) (big.int %s) %s\n",
				i, idToBigIntString(fingerStart), nodePtr.Id, idToBigIntString(nodePtr.Id), nodePtr.Ip)
		}
	}
	chordNode.fingerTable.RUnlock()

	outputString += "\t Inverted index tree: \n"
	chordNode.hashTable.RLock()
	for key, fileRecords := range chordNode.hashTable.table {
		hashKey, _ := chordNode.hashString(key)
		outputString += fmt.Sprintf("\t\t (key %s) (hashKey string %s) (hashKey big.int %s)\n",
			key, hashKey, idToBigIntString(hashKey))
		for _, fileRecordPtr := range fileRecords {
			if fileRecordPtr == nil {
				outputString += "\t\t\t<nil>\n"
			} else {
				outputString += fmt.Sprintf("\t\t\t\"%s\" \tin %s\n",
					fileRecordPtr.Filename, fileRecordPtr.OwnerIp)
			}
		}
	}
	chordNode.hashTable.RUnlock()

	outputString += "\t Connections: \n"
	chordNode.stubsPool.RLock()
	for ip := range chordNode.stubsPool.pool {
		outputString += "\t\t" + string(ip) + "\n"
	}
	chordNode.stubsPool.RUnlock()

	return outputString
}
