package gossiper

import (
	"fmt"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/file_sharing"
	"github.com/2_alt_hw2/Peerster/ui"
)

func (g *Gossiper) handleUIRequest(r ui.Request) {
	g.uiInterface.ResponseChannel <- r.Process(g)
}

// HandlePullPeersRequest implements ui.HandlePullPeersRequest
func (g *Gossiper) HandlePullPeersRequest() []string {
	g.log.Trace("UI asked for peer list")
	return g.peerPool.List()
}

// HandlePushPeerRequest implements ui.HandlePushPeerRequest
func (g *Gossiper) HandlePushPeerRequest(peer string) {
	g.log.Debug(fmt.Sprintf("UI asked to push peer: %v", peer))
	err := g.peerPool.Insert(peer)
	if err != nil {
		g.log.Debug(fmt.Sprintf("Unable to add peer: %v", err))
	}
}

// HandlePullContactsRequest implements ui.HandlePullContactsRequest
func (g *Gossiper) HandlePullContactsRequest() []string {
	g.log.Trace("UI asked for contact list")
	return g.rumorRepository.KnownOrigins(0)
}

// HandlePushContactRequest implements ui.HandlePushContactRequest
func (g *Gossiper) HandlePushContactRequest(contact string) {
	g.log.Debug(fmt.Sprintf("UI asked to push contact: %v", contact))
	g.rumorRepository.AddOrigin(contact)
}

// HandlePullRumorsRequest implements ui.HandlePullRumorsRequest
func (g *Gossiper) HandlePullRumorsRequest(offset uint32) []*types.RumorMessage {
	g.log.Trace(fmt.Sprintf("UI asked for rumor list from offset %v", offset))
	return g.rumorRepository.List(offset)
}

// HandlePushRumorRequest implements ui.HandlePushRumorRequest
func (g *Gossiper) HandlePushRumorRequest(message *types.Message) {
	g.log.Debug(fmt.Sprintf("UI asked to push rumor: %v", message))
	g.handleClientCmd(*message)
}

// HandlePullSettingsRequest implements ui.HandlePullPeersRequest
func (g *Gossiper) HandlePullSettingsRequest() types.Settings {
	g.log.Trace("UI asked for settings")
	return g.settings
}

// HandlePushSettingsRequest implements ui.HandlePushPeerRequest
func (g *Gossiper) HandlePushSettingsRequest(settings types.Settings) {
	g.log.Debug(fmt.Sprintf("UI asked to change settings to %v", settings))
	g.log.Warn("Settings update: Only name update is implemented yet")
	g.name = settings.Name
}

// HandlePullPrivateMessageRequest implements ui.HandlePullPrivateMessageRequest
func (g *Gossiper) HandlePullPrivateMessageRequest(offset uint32) []*types.PrivateMessage {
	g.log.Trace(fmt.Sprintf("UI asked for private message list from offset %v", offset))
	return g.privateMsgRepository.List(offset)
}

// HandlePushPrivateMessageRequest implements ui.HandlePushPrivateMessageRequest
func (g *Gossiper) HandlePushPrivateMessageRequest(message *types.Message) {
	g.log.Trace(fmt.Sprintf("UI asked to push private message %v", message))
	g.handleClientCmd(*message)
}

// HandlePullSharedFilesRequest implements ui.HandlePullSharedFilesRequest
func (g *Gossiper) HandlePullSharedFilesRequest() []file_sharing.File {
	g.log.Trace("UI asked for shared files list")
	return g.sharedFilesRepository.ListShared()
}

// HandlePullSharableFilesRequest implements ui.HandlePullSharableFilesRequest
func (g *Gossiper) HandlePullSharableFilesRequest() []string {
	g.log.Trace("UI asked for sharable files list")
	files, err := g.sharedFilesRepository.ListSharable()
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable to get sharable files: %v", err))
	}
	return files
}

// HandlePushSharedFileRequest implements ui.HandlePushSharedFileRequest
func (g *Gossiper) HandlePushSharedFileRequest(filename string) {
	g.log.Trace(fmt.Sprintf("UI asked to share file %v", filename))
	err := g.sharedFilesRepository.Share(filename)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable share file %v: %v", filename, err))
	}
}

// HandlePushDownloadFileRequest implements ui.HandlePushDownloadFileRequest
func (g *Gossiper) HandlePushDownloadFileRequest(file file_sharing.File) {
	g.log.Trace(fmt.Sprintf("UI asked to download file %v", file))
	err := g.downloadManager.Download(file)
	if err != nil {
		g.log.Warn(fmt.Sprintf("Unable download file %v: %v", file.Name, err))
	}
}
