package types

import "fmt"

// RumorMessage contains the actual text of a user message to be gossipped.
type RumorMessage struct {
	Origin string
	ID     uint32
	Text   string
}

// String returns the string representation of a RumorMessage
func (r *RumorMessage) String(from string) string {
	return fmt.Sprintf("RUMOR origin %v from %v ID %v contents %v", r.Origin, from, r.ID, r.Text)
}
