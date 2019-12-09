package types

import "fmt"

// PrivateMessage represents a message with a specific destination
type PrivateMessage struct {
	Origin      string
	ID          uint32
	Text        string
	Destination string
	HopLimit    uint32
}

// ToRumor return the rumor equivalent of the private message
func (p PrivateMessage) ToRumor() RumorMessage {
	return RumorMessage{
		Origin: p.Origin,
		ID:     p.ID,
		Text:   p.Text,
	}
}

func (p PrivateMessage) String() string {
	return fmt.Sprintf("Private message: origin %v hop-limit %v contents %v", p.Origin, p.HopLimit, p.Text)
}
