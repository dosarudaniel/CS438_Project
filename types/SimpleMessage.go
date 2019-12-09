package types

import "fmt"

// SimpleMessage represents a Peerster message in the simple protocol
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

func (s *SimpleMessage) String() string {
	return fmt.Sprintf("SIMPLE MESSAGE origin %v from %v contents %v",
		s.OriginalName,
		s.RelayPeerAddr,
		s.Contents)
}
