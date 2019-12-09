package types

import "fmt"

type SearchResult struct {
	FileName string
	MetafileHash []byte
	ChunkMap []uint64
	ChunkCount uint64
}

func (sResult *SearchResult) String() string {
	return fmt.Sprintf("SearchResult{FileName %v MetafileHash.len %v ChunkMap.len: %v ChunkCount: %v}",
		sResult.FileName, len(sResult.MetafileHash), len(sResult.ChunkMap), sResult.ChunkCount)
}
