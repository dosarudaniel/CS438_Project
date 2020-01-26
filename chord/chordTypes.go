package chord

import (
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"sync"
)

/*
	We place here types that are directly related to ChordNode implementation
*/

type nodeID string
type ipAddr string

type fingerTable []*Node

type fingerTableWithMux struct {
	table fingerTable
	sync.RWMutex
}

type nodeWithMux struct {
	nodePtr *Node
	sync.RWMutex
}

type ipToStubMap map[ipAddr]ChordClient

type stubsPoolWithMux struct {
	pool ipToStubMap
	sync.RWMutex
}

type ChordConfig struct {
	NumOfBitsInID int // number of bits in ID
}
