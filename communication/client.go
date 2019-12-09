package communication

//========= REDACTED =========
import (
"github.com/2_alt_hw2/Peerster/types"
)

type ClientMessageType int

const (
	None           ClientMessageType = iota
	Rumor          ClientMessageType = iota
	PrivateMessage ClientMessageType = iota
	FileSharing    ClientMessageType = iota
	FileRequest    ClientMessageType = iota
)

func CheckClientMessageTypeFromStrings(msg string, destination string, file string, request string) ClientMessageType {
	if msg != "" {
		if file != "" || request != "" {
			return None
		} else if destination != "" {
			return PrivateMessage
		} else {
			return Rumor
		}
	} else if file != "" {
		if request == "" && destination == "" {
			return FileSharing
		} else if request != "" && destination != "" {
			return FileRequest
		}
	}

	return None
}

func CheckClientMessageTypeFromStruct(message types.Message) ClientMessageType {
	if message.Text != "" {
		if message.File != nil || message.Request != nil {
			return None
		} else if message.Destination != nil && *message.Destination != "" {
			return PrivateMessage
		} else {
			return Rumor
		}
	} else if message.File != nil {
		if message.Request == nil && (message.Destination == nil || *message.Destination == "") {
			return FileSharing
		} else if message.Request != nil && message.Destination != nil && *message.Destination != "" {
			return FileRequest
		}
	}

	return None
}
