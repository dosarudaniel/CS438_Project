package ui

import (
	"sync"
)

// Request represents UI request to the backend
type Request message

// Response represents UI response from the backend
type Response message

// Requester represents an interface to make asynchronous requests to backend
type Requester struct {
	nextSeqNumber        uint32
	nextSeqNumberMutex   *sync.Mutex
	cReq                 chan Request
	cRes                 chan Response
	pendingRequests      map[uint32]interface{}
	pendingRequestsMutex *sync.Mutex
}

// BackendInterface encapsulate two channel for the backend communication
type BackendInterface struct {
	RequestChannel  <-chan Request
	ResponseChannel chan<- Response
}

type messageType uint32

type message struct {
	seqNumber   uint32
	messageType messageType
	parameter   interface{}
}

// NewRequester creates a new interface to make UI requests to the backend
func NewRequester() *Requester {
	var nextSeqNumberMutex sync.Mutex
	var pendingRequestsMutex sync.Mutex
	r := Requester{
		0,
		&nextSeqNumberMutex,
		make(chan Request),
		make(chan Response),
		make(map[uint32]interface{}, 0),
		&pendingRequestsMutex,
	}
	go func() {
		for msg := range r.cRes {
			pendingRequestsMutex.Lock()
			callback := r.pendingRequests[msg.seqNumber]
			msg.handleResponse(callback)
			delete(r.pendingRequests, msg.seqNumber)
			pendingRequestsMutex.Unlock()
		}
	}()
	return &r
}

// BackendInterface is a getter for the backend interface
func (r *Requester) BackendInterface() BackendInterface {
	return BackendInterface{
		r.cReq,
		r.cRes,
	}
}

func (r *Requester) makeRequest(msgType messageType, parameter interface{}, callback interface{}) {
	r.nextSeqNumberMutex.Lock()
	seqNumber := r.nextSeqNumber
	r.nextSeqNumber++
	r.nextSeqNumberMutex.Unlock()

	r.pendingRequestsMutex.Lock()
	r.pendingRequests[seqNumber] = callback
	r.pendingRequestsMutex.Unlock()

	r.cReq <- Request{
		seqNumber:   seqNumber,
		messageType: msgType,
		parameter:   parameter,
	}
}

// Process processes the request with the given request handler
func (r *Request) Process(handler RequestHandler) Response {
	return Response{
		seqNumber:   r.seqNumber,
		messageType: r.messageType,
		parameter:   r.handleRequest(handler),
	}
}
