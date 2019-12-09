package messaging

import (
	"sync"
	"github.com/2_alt_hw2/Peerster/types"
)

// PeerStatusRepository stores pending communications with other peers
type PeerStatusRepository struct {
	//            targetAddr origin     seqNum
	pendingRumors sync.Map //Concurrent version of map[string]map[string]map[uint32]interface{}
}

// NewPeerStatusRepository creates a new PendingRumorRepository
func NewPeerStatusRepository() *PeerStatusRepository {
	return &PeerStatusRepository{}
}

// Insert inserts a new pending rumor in the repository
// Duplicates have no effect
func (p *PeerStatusRepository) Insert(targetAddr, origin string, ID uint32) {
	targetPendingRumors, found := p.pendingRumors.Load(targetAddr)
	if !found {
		var tmp sync.Map
		targetPendingRumors, _ = p.pendingRumors.LoadOrStore(targetAddr, &tmp)
	}

	originPendingRumors, found := (targetPendingRumors.(*sync.Map)).Load(origin)
	if !found {
		var tmp sync.Map
		originPendingRumors, _ = (targetPendingRumors.(*sync.Map)).LoadOrStore(origin, &tmp)
	}

	// We store nil to indicate that the key ID exists
	(originPendingRumors.(*sync.Map)).Store(ID, nil)
}

// Ack checks all pending communications with targetAddr concerning rumors from origin up to given ID
// It returns a list of all matching ID's
func (p *PeerStatusRepository) Ack(targetAddr string, status []types.PeerStatus) []types.PeerStatus {
	ackedIDs := make([]types.PeerStatus, 0)

	targetPendingRumors, found := p.pendingRumors.Load(targetAddr)
	if found {
		for _, ack := range status {
			originPendingRumors, found := (targetPendingRumors.(*sync.Map)).Load(ack.Identifier)
			if found {
				(originPendingRumors.(*sync.Map)).Range(func(key, _ interface{}) bool {
					pendingID := key.(uint32)
					if pendingID <= ack.NextID {
						ackedIDs = append(ackedIDs, types.PeerStatus{
							Identifier: ack.Identifier,
							NextID:     pendingID,
						})
						(originPendingRumors.(*sync.Map)).Delete(pendingID)
					}
					return true
				})
			}
		}
	}

	return ackedIDs
}

// Delete deletes one entry from the repository
// No effect on nonexistent entries
func (p *PeerStatusRepository) Delete(targetAddr string, origin string, ID uint32) bool {
	targetPendingRumors, found := p.pendingRumors.Load(targetAddr)
	if found {
		originPendingRumors, found := (targetPendingRumors.(*sync.Map)).Load(origin)
		if found {
			_, found := (originPendingRumors.(*sync.Map)).Load(ID)
			if found {
				(originPendingRumors.(*sync.Map)).Delete(ID)
				return true
			}
		}
	}
	return false
}

// AllPending returns the list of all pending rumors
func (p *PeerStatusRepository) AllPending() map[string][]types.PeerStatus {
	pending := make(map[string][]types.PeerStatus, 0)

	p.pendingRumors.Range(func(targetAddr, pendingRumors interface{}) bool {
		(pendingRumors.(*sync.Map)).Range(func(origin, originPendingRumors interface{}) bool {
			(originPendingRumors.(*sync.Map)).Range(func(pendingID, _ interface{}) bool {
				pending[targetAddr.(string)] = append(pending[targetAddr.(string)], types.PeerStatus{
					Identifier: origin.(string),
					NextID:     pendingID.(uint32),
				})
				return true
			})
			return true
		})
		return true
	})

	return pending
}

// Pending returns the list of pending rumors for a peer
func (p *PeerStatusRepository) Pending(targetAddr string) []types.PeerStatus {
	pending := make([]types.PeerStatus, 0)

	targetPendingRumors, found := p.pendingRumors.Load(targetAddr)
	if found {
		(targetPendingRumors.(*sync.Map)).Range(func(origin, originPendingRumors interface{}) bool {
			(originPendingRumors.(*sync.Map)).Range(func(pendingID, _ interface{}) bool {
				pending = append(pending, types.PeerStatus{
					Identifier: origin.(string),
					NextID:     pendingID.(uint32),
				})
				return true
			})
			return true
		})
	}

	return pending
}
