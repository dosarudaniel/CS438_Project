package types

import (
	"fmt"
	"github.com/dosarudaniel/CS438_Project/logger"
	"strings"
)

// StatusPacket lists the messages wanted by a peer
type StatusPacket struct {
	Want []PeerStatus
}

func (s StatusPacket) String() string {
	if len(s.Want) == 0 {
		return ""
	}
	var b strings.Builder
	for _, status := range s.Want {
		_, err := fmt.Fprintf(&b, " peer %v nextID %v", status.Identifier, status.NextID)
		if err != nil {
			logger.DefaultLogger().Warn(fmt.Sprintf("Unable to construct status packet: %v", err))
		}
	}
	return b.String()[1:]
}
