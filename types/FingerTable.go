package types

import "fmt"

type FingerEntry struct {
	Id   			[]byte
	GossiperAddress string
}

type FingerTable []FingerEntry


func (fe FingerEntry) String() string {
	return fmt.Sprintf("FingerEntry: Content{%v - %v}", fe.Id, fe.GossiperAddress)
}

func (ft FingerTable) String() string {
	outputString := ""
	// for entry in *FingerTable {
	// 	outputString +=
	// }

	return outputString //fmt.Sprintf("FingerTable Content{TODO}")
}
