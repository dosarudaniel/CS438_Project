package types

import "fmt"

type FingerEntry struct {
	// ID hash of (n + 2^i) % 2^m
	// n is the current node
	// i the index in the finger table i can one of the following values: 0,1,2,3,...,m-1
	// m is the number of bits
	id   			string // hexstring of size m bits, eg. "3fde56" of size m=32
	gossiperAddress string

}

type FingerTable []FingerEntry


func (fe FingerEntry) String() string {
	return fmt.Sprintf("FingerEntry: Content{%v - %v}", fe.id, fe.gossiperAddress)
}

func (ft FingerTable) String() string {
	outputString := ""
	// for entry in FingerTable {
	// 	outputString += entry.id + " <-> " + entry.gossiperAddress
	// }

	return outputString
}
