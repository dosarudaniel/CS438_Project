package chord

import (
	"context"
	"errors"
	"fmt"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
)

// GetPredecessor (RPC) returns a pointer to the predecessor node
func (chordNode *ChordNode) GetPredecessor(ctx context.Context, _ *empty.Empty) (*Node, error) {
	pred, doesExist := chordNode.getPredecessor()
	if !doesExist {
		return nil, &nilPredecessor{}
	}
	return &pred, nil
}

// FIXME CheckPredecessor relies on use of ctx.Timeout
/*
 FindSuccessor finds the successor of id
 n.find_successor(id)
 	if id is_in (n, successor]
 		return successor;
 	else
 		n' = closest_preceding_node(id);
 		return n'.find_successor(id);
*/
func (chordNode *ChordNode) FindSuccessor(ctx context.Context, messageIDPtr *ID) (*Node, error) {
	if messageIDPtr == nil {
		return nil, errors.New("id must not be nil")
	}

	id := messageIDPtr.Id
	n := chordNode.node
	succ, doesExist := chordNode.getSuccessor()
	if !doesExist {
		return nil, errors.New("successor does not exist")
	}

	if isBetweenTwoNodesRightInclusive(n.Id, id, succ.Id) {
		return &succ, nil
	} else {
		n0 := chordNode.ClosestPrecedingFinger(nodeID(id))
		if n0.Id == n.Id {
			return &n, nil
		}
		return chordNode.stubFindSuccessor(ipAddr(n0.Ip), context.Background(), &ID{Id: id})
	}
}

// Notify(n0) checks whether n0 needs to be my predecessor
// Algorithm:
// n.notify(n0)
//	 if (predecessor is nil or n0 is_in (predecessor; n))
//	 	predecessor = n0;
func (chordNode *ChordNode) Notify(ctx context.Context, n0Ptr *Node) (*empty.Empty, error) {
	emptyPtr := &empty.Empty{}

	switch {
	case n0Ptr == nil:
		return emptyPtr, errors.New("trying to notify nil node")
	case n0Ptr.Id == chordNode.node.Id:
		return emptyPtr, errors.New(fmt.Sprintf("node %s is notifying itself", chordNode.node.Id))
	}

	pred, doesPredExist := chordNode.getPredecessor()
	if !doesPredExist || isBetweenTwoNodesExclusive(pred.Id, n0Ptr.Id, chordNode.node.Id) {
		chordNode.setPredecessor(n0Ptr)
		if doesPredExist {
			_, err := chordNode.TransferKeys(context.Background(), &TransferKeysRequest{
				FromId: pred.Id,
				ToNode: n0Ptr,
			})
			if err != nil {
				log.WithFields(logrus.Fields{
					"fromID": pred.Id,
					"toID":   n0Ptr.Id,
					"err":    err,
				}).Warn("TransferKeys failed")
			}
		}
	}

	return emptyPtr, nil
}

func (chordNode *ChordNode) Put(ctx context.Context, keyValPtr *FileRecordWithKeyword) (*empty.Empty, error) {
	emptyPtr := &empty.Empty{}

	if keyValPtr == nil {
		return emptyPtr, nilError("FileRecordWithKeyword")
	}

	err := chordNode.hashTable.PutOrAppendOne(keyValPtr.Keyword, keyValPtr.Val)

	return emptyPtr, err
}

// Get implements RPC method Get (Key) returns (Val);
// It returns a list of documents that were stored under this keyword
func (chordNode *ChordNode) Get(ctx context.Context, messageKeyPtr *Key) (*Val, error) {
	if messageKeyPtr == nil {
		return nil, nilError("Key")
	}

	chordNode.hashTable.RLock()
	defer chordNode.hashTable.RUnlock() //TODO can using array of pointers create a data race here?

	fileRecords, doesExist := chordNode.hashTable.table[messageKeyPtr.Keyword]

	if doesExist {
		return &Val{FileRecords: fileRecords}, nil
	} else {
		return nil, errors.New(
			fmt.Sprintf("such key does not exist at node IP %s ID %s", chordNode.node.Ip, chordNode.node.Ip))
	}
}

func (chordNode *ChordNode) TransferKeys(ctx context.Context, req *TransferKeysRequest) (*empty.Empty, error) {
	emptyPtr := &empty.Empty{}

	if req == nil {
		return emptyPtr, nilError("TransferKeyRequest")
	}

	fromID := req.FromId
	toNode := req.ToNode

	if toNode.Id == chordNode.node.Id {
		return emptyPtr, errors.New("trying to transfer keys from to same node")
	}

	chordNode.hashTable.Lock()
	defer chordNode.hashTable.Unlock()

	keysToRemove := make([]string, 0)
	for key, val := range chordNode.hashTable.table {
		hashedKey, err := chordNode.hashString(key)
		if err != nil {
			return emptyPtr, err
		}
		// Check that the hashed key lies in the correct range before putting the value in our predecessor
		if isBetweenTwoNodesRightInclusive(fromID, hashedKey, toNode.Id) {
			// TODO potential improvement: enable PutManyFileRecords rpc to put several file records at once
			for _, fileRecordPtr := range val {
				if fileRecordPtr == nil {
					continue
				}
				if err := chordNode.stubPut(ipAddr(toNode.Ip), context.Background(),
					key, fileRecordPtr.Filename, fileRecordPtr.OwnerIp); err != nil {
					return emptyPtr, err
				}
			}

			keysToRemove = append(keysToRemove, key)
		}
	}

	for _, key := range keysToRemove {
		delete(chordNode.hashTable.table, key)
	}

	return emptyPtr, nil
}
