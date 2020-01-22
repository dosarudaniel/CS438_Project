package chord

import (
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"sync"
)

type fingerTable []Node

type fingerTableWithMux struct {
	table fingerTable
	sync.RWMutex
}
