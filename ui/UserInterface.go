package ui

import (
	"github.com/dosarudaniel/CS438_Project/file_sharing"
	"github.com/dosarudaniel/CS438_Project/types"
)

type pullPeersCallback func([]string)
type pushPeerCallback func()
type pullContactsCallback func([]string)
type pushContactCallback func()
type pullRumorsCallback func([]*types.RumorMessage)
type pushRumorCallback func()
type pullSettingsCallback func(types.Settings)
type pushSettingsCallback func()
type pullPrivateMessageCallback func([]*types.PrivateMessage)
type pushPrivateMessageCallback func()
type pullSharedFilesCallback func([]file_sharing.File)
type pullSharableFilesCallback func([]string)
type pushSharedFileCallback func()
type pushDownloadFileCallback func()

const (
	pullPeers          messageType = iota
	pushPeer           messageType = iota
	pullContacts       messageType = iota
	pushContact        messageType = iota
	pullRumors         messageType = iota
	pushRumor          messageType = iota
	pullSettings       messageType = iota
	pushSettings       messageType = iota
	pullPrivateMessage messageType = iota
	pushPrivateMessage messageType = iota
	pullSharedFiles    messageType = iota
	pullSharableFiles  messageType = iota
	pushSharedFile     messageType = iota
	pushDownloadFile   messageType = iota
)

// RequestPullPeers creates a new asynchronous request to fetch the peer list
func (r *Requester) RequestPullPeers(callback pullPeersCallback) {
	r.makeRequest(pullPeers, nil, callback)
}

// RequestPushPeer creates a new asynchronous request to push a new peer
func (r *Requester) RequestPushPeer(peer string, callback pushPeerCallback) {
	r.makeRequest(pushPeer, peer, callback)
}

// RequestPullContacts creates a new asynchronous request to fetch the contact list
func (r *Requester) RequestPullContacts(callback pullContactsCallback) {
	r.makeRequest(pullContacts, nil, callback)
}

// RequestPushContact creates a new asynchronous request to push a new contact
func (r *Requester) RequestPushContact(contact string, callback pushContactCallback) {
	r.makeRequest(pushContact, contact, callback)
}

// RequestPullRumors creates a new asynchronous request to fetch rumors
func (r *Requester) RequestPullRumors(lastID uint32, callback pullRumorsCallback) {
	r.makeRequest(pullRumors, lastID, callback)
}

// RequestPushRumor creates a new asynchronous request to push a new rumor
func (r *Requester) RequestPushRumor(rumor types.Message, callback pushRumorCallback) {
	r.makeRequest(pushRumor, &rumor, callback)
}

// RequestPullSettings creates a new asynchronous request to fetch settings
func (r *Requester) RequestPullSettings(callback pullSettingsCallback) {
	r.makeRequest(pullSettings, nil, callback)
}

// RequestPushSettings creates a new asynchronous request to push settings
func (r *Requester) RequestPushSettings(settings types.Settings, callback pushSettingsCallback) {
	r.makeRequest(pushSettings, settings, callback)
}

// RequestPullPrivateMessage creates a new asynchronous request to fetch settings
func (r *Requester) RequestPullPrivateMessage(lastID uint32, callback pullPrivateMessageCallback) {
	r.makeRequest(pullPrivateMessage, lastID, callback)
}

// RequestPushPrivateMessage creates a new asynchronous request to push settings
func (r *Requester) RequestPushPrivateMessage(pm types.Message, callback pushPrivateMessageCallback) {
	r.makeRequest(pushPrivateMessage, &pm, callback)
}

// RequestPullSharedFiles creates a new asynchronous request to fetch shared files
func (r *Requester) RequestPullSharedFiles(callback pullSharedFilesCallback) {
	r.makeRequest(pullSharedFiles, nil, callback)
}

// RequestPullSharableFiles creates a new asynchronous request to fetch sharable files
func (r *Requester) RequestPullSharableFiles(callback pullSharableFilesCallback) {
	r.makeRequest(pullSharableFiles, nil, callback)
}

// RequestPushSharedFile creates a new asynchronous request to push a shared file
func (r *Requester) RequestPushSharedFile(filename string, callback pushSharedFileCallback) {
	r.makeRequest(pushSharedFile, filename, callback)
}

// RequestPushDownloadFile creates a new asynchronous request to push a download file
func (r *Requester) RequestPushDownloadFile(file file_sharing.File, callback pushDownloadFileCallback) {
	r.makeRequest(pushDownloadFile, file, callback)
}

// RequestHandler describes an interface that can process UI requests
type RequestHandler interface {
	HandlePullPeersRequest() []string
	HandlePushPeerRequest(string)
	HandlePullContactsRequest() []string
	HandlePushContactRequest(string)
	HandlePullRumorsRequest(uint32) []*types.RumorMessage
	HandlePushRumorRequest(*types.Message)
	HandlePullSettingsRequest() types.Settings
	HandlePushSettingsRequest(types.Settings)
	HandlePullPrivateMessageRequest(uint32) []*types.PrivateMessage
	HandlePushPrivateMessageRequest(*types.Message)
	HandlePullSharedFilesRequest() []file_sharing.File
	HandlePullSharableFilesRequest() []string
	HandlePushSharedFileRequest(string)
	HandlePushDownloadFileRequest(file_sharing.File)
}

func (r *Request) handleRequest(handler RequestHandler) interface{} {
	switch r.messageType {
	case pullPeers:
		return handler.HandlePullPeersRequest()
	case pushPeer:
		handler.HandlePushPeerRequest(r.parameter.(string))
	case pullContacts:
		return handler.HandlePullContactsRequest()
	case pushContact:
		handler.HandlePushContactRequest(r.parameter.(string))
	case pullRumors:
		return handler.HandlePullRumorsRequest(r.parameter.(uint32))
	case pushRumor:
		handler.HandlePushRumorRequest(r.parameter.(*types.Message))
	case pullSettings:
		return handler.HandlePullSettingsRequest()
	case pushSettings:
		handler.HandlePushSettingsRequest(r.parameter.(types.Settings))
	case pullPrivateMessage:
		return handler.HandlePullPrivateMessageRequest(r.parameter.(uint32))
	case pushPrivateMessage:
		handler.HandlePushPrivateMessageRequest(r.parameter.(*types.Message))
	case pullSharedFiles:
		return handler.HandlePullSharedFilesRequest()
	case pullSharableFiles:
		return handler.HandlePullSharableFilesRequest()
	case pushSharedFile:
		handler.HandlePushSharedFileRequest(r.parameter.(string))
	case pushDownloadFile:
		handler.HandlePushDownloadFileRequest(r.parameter.(file_sharing.File))
	}
	return nil
}

func (msg *Response) handleResponse(callback interface{}) {
	switch msg.messageType {
	case pullPeers:
		callback.(pullPeersCallback)(msg.parameter.([]string))
	case pushPeer:
		callback.(pushPeerCallback)()
	case pullContacts:
		callback.(pullContactsCallback)(msg.parameter.([]string))
	case pushContact:
		callback.(pushContactCallback)()
	case pullRumors:
		callback.(pullRumorsCallback)(msg.parameter.([]*types.RumorMessage))
	case pushRumor:
		callback.(pushRumorCallback)()
	case pullSettings:
		callback.(pullSettingsCallback)(msg.parameter.(types.Settings))
	case pushSettings:
		callback.(pushSettingsCallback)()
	case pullPrivateMessage:
		callback.(pullPrivateMessageCallback)(msg.parameter.([]*types.PrivateMessage))
	case pushPrivateMessage:
		callback.(pushPrivateMessageCallback)()
	case pullSharedFiles:
		callback.(pullSharedFilesCallback)(msg.parameter.([]file_sharing.File))
	case pullSharableFiles:
		callback.(pullSharableFilesCallback)(msg.parameter.([]string))
	case pushSharedFile:
		callback.(pushSharedFileCallback)()
	case pushDownloadFile:
		callback.(pushDownloadFileCallback)()
	}
}
