package chord

import "sync"
import . "github.com/dosarudaniel/CS438_Project/services/chord_service"

// In this file, we put all helper functions and types used by Chord package
// "helper" can be defined as something that is not worth own file/package and has tiny role in the chord functionality

type nodeWithMux struct {
	node *Node
	sync.RWMutex
}

type successorsListWithMux struct {
	successorsList []Node
	sync.RWMutex
}
