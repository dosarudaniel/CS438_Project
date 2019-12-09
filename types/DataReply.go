package types

import "fmt"

type DataReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   []byte
	Data        []byte
}

func (dr *DataReply) String() string {
	return fmt.Sprintf("DataReply{Origin %v sent Hash %v for %v Data.len=%v. Hop-limit: %v}",
		dr.Origin, dr.HashValue, dr.Destination, len(dr.Data), dr.HopLimit)
}
