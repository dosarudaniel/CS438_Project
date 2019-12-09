package messaging

import (
	"github.com/2_alt_hw2/Peerster/types"
)

// PrivateMessageRepository represent a simple volatile database of rumors
type PrivateMessageRepository struct {
	list []*types.PrivateMessage
}

// NewPrivateMessageRepository creates a new privateMessage repository
func NewPrivateMessageRepository() *PrivateMessageRepository {
	return &PrivateMessageRepository{
		make([]*types.PrivateMessage, 0),
	}
}

// Insert adds a new privateMessage to the repository
func (r *PrivateMessageRepository) Insert(pm *types.PrivateMessage) {
	r.list = append(r.list, pm)
}

// List returns the list of all private after a given offset
func (r *PrivateMessageRepository) List(offset uint32) []*types.PrivateMessage {
	if offset >= uint32(len(r.list)) {
		return []*types.PrivateMessage{}
	}
	return r.list[offset:]
}
