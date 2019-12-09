package types

import "fmt"

type SearchRequest struct {
	Origin string
	Budget uint64
	Keywords []string
}

func (sRequest *SearchRequest) String() string {
	return fmt.Sprintf("SearchRequest{Origin %v Budget %v Keywords.len: %v}",
		sRequest.Origin, sRequest.Budget, len(sRequest.Keywords))
}
