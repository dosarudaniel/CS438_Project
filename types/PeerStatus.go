package types

// PeerStatus summarizes the set of messages the sending peer has seen so far
type PeerStatus struct {
	Identifier string
	NextID     uint32
}
