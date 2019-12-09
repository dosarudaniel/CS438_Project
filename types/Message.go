package types

import (
	"encoding/hex"
	"fmt"
)

// Message represents a client text message without metadata
type Message struct {
	Text        string
	Destination *string
	File        *string
	Request     *[]byte
}

func (s Message) String() string {
	if s.Text != "" {
		if s.Destination != nil {
			return fmt.Sprintf("CLIENT MESSAGE %v dest %v", s.Text, *s.Destination)
		} else {
			return fmt.Sprintf("CLIENT MESSAGE %v", s.Text)
		}
	} else if s.File != nil {
		if s.Request != nil && s.Destination != nil {
			return fmt.Sprintf("File download request \"%v\" (%v) at %v", *s.File, hex.EncodeToString(*s.Request), *s.Destination)
		} else if s.Request == nil && s.Destination == nil {
			return fmt.Sprintf("File sharing request \"%v\"", *s.File)
		}
	}

	return fmt.Sprintf("Invalid message format: {\n"+
		"\tText: %v\n"+
		"\tDestination: %v\n"+
		"\tFile: %v\n"+
		"\tRequest: %v\n"+
		"}", s.Text, s.Destination, s.File, s.Request)
}
