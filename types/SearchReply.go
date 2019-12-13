package types

import "fmt"

type SearchReply struct {
	Origin string
	Destination string
	HopLimit uint32
	Results []*SearchResult
}

func (sReply *SearchReply) String() string {
	return fmt.Sprintf("SearchReply{Origin %v Destination %v HopLimit: %v, Results.len = %v}",
		sReply.Origin sReply.Destination, sReply.HopLimit, len(sReply.Results))
}
