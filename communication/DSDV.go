package communication

import "fmt"

// DSDV represents a destination-sequenced distance vector
type DSDV struct {
	routingTable map[string]*routingEntry
}

type routingEntry struct {
	nextHop   string
	seqNumber uint32
}

// NewDSDV creates and return a new DSDV
func NewDSDV() *DSDV {
	return &DSDV{
		make(map[string]*routingEntry),
	}
}

// Upsert updates the entry for origin if the seqNumber is bigger than the stored one
// New entries are stored unconditionally
func (dsdv *DSDV) Upsert(origin string, seqNumber uint32, peerAddr string) bool {
	entry, found := dsdv.routingTable[origin]
	if !found {
		dsdv.routingTable[origin] = &routingEntry{
			nextHop:   peerAddr,
			seqNumber: seqNumber,
		}
		return true
	} else if entry.seqNumber < seqNumber {
		dsdv.routingTable[origin].seqNumber = seqNumber
		return true
	}

	return false
}

// RouteTo returns the nextHop for a given origin if it is known or an error
func (dsdv *DSDV) RouteTo(origin string) (string, error) {
	entry, found := dsdv.routingTable[origin]
	if found {
		return entry.nextHop, nil
	}
	return "", fmt.Errorf("route to %v not found", origin)
}
