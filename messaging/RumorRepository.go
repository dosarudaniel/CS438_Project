package messaging

import (
	"fmt"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/logger"
)

// RumorRepository represent a simple volatile database of rumors
type RumorRepository struct {
	rumors                map[string]map[uint32]*types.RumorMessage
	origins               []string
	list                  []*types.RumorMessage
	status                map[string]uint32
	peerStatusBuffer      []types.PeerStatus
	peerStatusBufferDirty bool
	log                   *logger.Logger
}

// NewRumorRepository creates a new rumor repository
func NewRumorRepository(log *logger.Logger) *RumorRepository {
	return &RumorRepository{
		make(map[string]map[uint32]*types.RumorMessage, 0),
		make([]string, 0),
		make([]*types.RumorMessage, 0),
		make(map[string]uint32, 0),
		make([]types.PeerStatus, 0),
		true,
		log,
	}
}

// AddOrigin inserts a new empty sub-repository for a given peer
func (r *RumorRepository) AddOrigin(origin string) {
	if origin != "" {
		_, ok := r.rumors[origin]
		if !ok {
			r.rumors[origin] = make(map[uint32]*types.RumorMessage)
			r.origins = append(r.origins, origin)
			r.updateStatus(origin, 0)
		}
	}
}

// Insert adds a new rumor to the repository. Nothing happen for know rumors
func (r *RumorRepository) Insert(rumor *types.RumorMessage) {
	rumorsOfOrigin, ok := r.rumors[rumor.Origin]
	if ok {
		_, ok := rumorsOfOrigin[rumor.ID]
		if !ok {
			rumorsOfOrigin[rumor.ID] = rumor
		}
	} else {
		r.origins = append(r.origins, rumor.Origin)
		rumorsOfOrigin := make(map[uint32]*types.RumorMessage)
		rumorsOfOrigin[rumor.ID] = rumor
		r.rumors[rumor.Origin] = rumorsOfOrigin
	}
	if rumor.Text != "" {
		r.list = append(r.list, rumor)
	}
	r.updateStatus(rumor.Origin, rumor.ID)
}

// Contains indicates if rumor for a given origin and id is in the repository
func (r *RumorRepository) Contains(origin string, id uint32) bool {
	rumorsOfOrigin, ok := r.rumors[origin]
	if ok {
		_, ok = rumorsOfOrigin[id]
		return ok
	}
	return false
}

// Find searches for a known rumor for a given origin and id
func (r *RumorRepository) Find(origin string, id uint32) (*types.RumorMessage, error) {
	rumorsOfOrigin, ok := r.rumors[origin]
	if !ok {
		return nil, fmt.Errorf("origin(%v) of rumor not found", origin)
	}

	rumor, ok := rumorsOfOrigin[id]
	if !ok {
		return nil, fmt.Errorf("rumor(%v) not found for known origin(%v)", id, origin)
	}
	return rumor, nil
}

// List returns the list of all rumors after a given offset
func (r *RumorRepository) List(offset uint32) []*types.RumorMessage {
	if offset >= uint32(len(r.list)) {
		return []*types.RumorMessage{}
	}
	return r.list[offset:]
}

// List returns the list of all rumors after a given offset
func (r *RumorRepository) KnownOrigins(offset uint32) []string {
	if offset >= uint32(len(r.origins)) {
		return []string{}
	}
	return r.origins[offset:]
}

// GetStatus returns the list of wanted rumor by peerID
func (r *RumorRepository) GetStatus() []types.PeerStatus {
	if r.peerStatusBufferDirty {
		r.peerStatusBuffer = make([]types.PeerStatus, 0)
		for peerID, nextRumorID := range r.status {
			s := types.PeerStatus{Identifier: peerID, NextID: nextRumorID}
			r.peerStatusBuffer = append(r.peerStatusBuffer, s)
		}
		r.peerStatusBufferDirty = false
	}
	return r.peerStatusBuffer
}

// FirstUnknown compares a list of PeerStatus with the repository content and returns:
// * The first matching Rumor if one is found (and the boolean represents nothing)
// * Or a nil rumor and the boolean represents interesting rumor(s) in the status
func (r *RumorRepository) FirstUnknown(status []types.PeerStatus) (*types.RumorMessage, bool) {
	peerHasInterestingRumors := false
	for _, statMsg := range status {
		rumor, err := r.Find(statMsg.Identifier, statMsg.NextID)
		if err == nil {
			return rumor, false
		} else if !peerHasInterestingRumors && r.NextRumor(statMsg.Identifier) < statMsg.NextID {
			// The peer has a at least one rumor we want
			peerHasInterestingRumors = true
		}
	}
	return nil, peerHasInterestingRumors
}

// NextRumor returns the next wanted rumor for a peer name
func (r *RumorRepository) NextRumor(origin string) uint32 {
	return r.status[origin]
}

func (r *RumorRepository) updateStatus(origin string, id uint32) {
	r.log.Info(fmt.Sprintf("Updating status for: peer %v id %v", origin, id))

	currentWantedID, found := r.status[origin]
	if !found {
		r.log.Trace(fmt.Sprintf("\tOrigin not found. Init wanted to 1"))
		currentWantedID = 1
		r.status[origin] = currentWantedID
		r.peerStatusBufferDirty = true
	} else {
		r.log.Trace(fmt.Sprintf("\tOrigin found. Current wanted is: %v", currentWantedID))
	}

	if id == currentWantedID {
		currentWantedID++
		r.peerStatusBufferDirty = true
		for r.Contains(origin, currentWantedID) {
			currentWantedID++
		}
		r.status[origin] = currentWantedID
	}

	r.log.Debug(fmt.Sprintf("\tNew status: %v", r.GetStatus()))
}
