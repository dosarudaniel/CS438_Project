package gossiper

import (
	"fmt"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/communication"
	"github.com/2_alt_hw2/Peerster/file_sharing"
)

// handleClientCmd handles incoming client commands
func (g *Gossiper) handleClientCmd(message types.Message) {
	fmt.Println(message.String())
	g.log.Debug(fmt.Sprintf(g.peerPool.String()))

	switch communication.CheckClientMessageTypeFromStruct(message) {
	case communication.Rumor:
		g.sendNewRumor(g.peerPool.GetRandom(), message.Text)
	case communication.PrivateMessage:
		g.sendPrivateMessage(message)
	case communication.FileSharing:
		err := g.sharedFilesRepository.Share(*message.File)
		if err != nil {
			g.log.Warn(fmt.Sprintf("Unable to index file %v for sharing: %v", message.File, err))
		}
	case communication.FileRequest:
		hash, err := file_sharing.BytesToHash(*message.Request)
		if err != nil {
			g.log.Warn(fmt.Sprintf("Unable to download file %v from %v: %v", message.File, message.Destination, err))
		}

		file := file_sharing.File{
			Name:   *message.File,
			Origin: *message.Destination,
			Hash:   hash.ToHex(),
		}
		err = g.downloadManager.Download(file)
		if err != nil {
			g.log.Warn(fmt.Sprintf("Unable to download file %v from %v: %v", message.File, message.Destination, err))
		}
	}
}
