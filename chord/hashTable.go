package chord

import (
	"sync"
)

type hashTable map[string]ipAddr

type HashTableWithMux struct {
	table hashTable
	sync.RWMutex
}

func NewHashTable() HashTableWithMux {
	return HashTableWithMux{
		table:   make(hashTable),
		RWMutex: sync.RWMutex{},
	}
}

// Get returns the IP address, which the key is mapped to, and doesExist
func (tableWithMux *HashTableWithMux) Get(key string) (ipAddr, bool) {
	tableWithMux.RLock()
	defer tableWithMux.RUnlock()

	addr, ok := tableWithMux.table[key]
	if ok {
		return addr, true
	} else {
		return "", false
	}
}

func (tableWithMux *HashTableWithMux) Put(key string, addr ipAddr) {
	tableWithMux.Lock()
	defer tableWithMux.Unlock()

	tableWithMux.table[key] = addr
}

// Retrieve deletes key-value pair from the hash table and returns the deleted value and if such pair existed
// returns ipAddr, didSuchPairExist
func (tableWithMux *HashTableWithMux) Retrieve(key string) (ipAddr, bool) {
	tableWithMux.Lock()
	defer tableWithMux.Unlock()

	addr, present := tableWithMux.table[key]
	if present {
		delete(tableWithMux.table, key)
		return addr, true
	} else {
		return "", false
	}
}

func (tableWithMux *HashTableWithMux) NumOfPairs() int {
	tableWithMux.RLock()
	defer tableWithMux.RUnlock()

	return len(tableWithMux.table)
}
