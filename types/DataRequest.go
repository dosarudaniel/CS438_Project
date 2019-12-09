package types

import "fmt"

type DataRequest struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   []byte
}

func (dr *DataRequest) String() string {
	return fmt.Sprintf("DataRequest{Origin %v wants Hash %v from %v. Hop-limit: %v}",
		dr.Origin, dr.HashValue, dr.Destination, dr.HopLimit)
}
