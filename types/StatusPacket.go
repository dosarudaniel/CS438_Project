package types

import (
	"fmt"
//========= REDACTED =========
"github.com/2_alt_hw2/Peerster/logger"
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
