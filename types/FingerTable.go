package types

import "fmt"

type FingerEntry struct {
	Id   []byte
	GossiperAddress string
}

type FingerTable []*FingerEntry

func (ft *FingerTable) String() string {
	return fmt.Sprintf("FingerTable Content{TODO}")
}
