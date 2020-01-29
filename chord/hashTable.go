/*
	Important notice:
	Read access must be RLocked and RUnlocked. HashTableWithMux only provides thread-safe write access.
*/

package chord

import (
	chord "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"sync"
)

type hashTable map[string][]*chord.FileRecord

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

// PutOrReplacePair puts or replace a key-value pair
func (tableWithMux *HashTableWithMux) PutOrReplacePair(keyword string, fileRecords []*chord.FileRecord) {
	tableWithMux.Lock()
	defer tableWithMux.Unlock()

	tableWithMux.table[keyword] = fileRecords
}

// PutOrAppendOne appends an existing value in key-value pair or creates a new pair with single-element value
func (tableWithMux *HashTableWithMux) PutOrAppendOne(keyword string, fileRecord *chord.FileRecord) error {
	if fileRecord == nil {
		return nilError("fileRecord")
	}

	tableWithMux.Lock()
	defer tableWithMux.Unlock()

	fileRecords, doExist := tableWithMux.table[keyword]
	if doExist {
		tableWithMux.table[keyword] = append(fileRecords, fileRecord)
	} else {
		tableWithMux.table[keyword] = []*chord.FileRecord{fileRecord}
	}

	return nil
}

func (tableWithMux *HashTableWithMux) NumOfPairs() int {
	tableWithMux.RLock()
	defer tableWithMux.RUnlock()

	return len(tableWithMux.table)
}
