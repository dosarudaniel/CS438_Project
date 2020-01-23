package chord

import (
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"github.com/golang/protobuf/ptypes/empty"
	"context"
)

// StabilizeDaemon() is called periodically.
// It verifies n’s immediate successor, and tells the successor about n.
// Algorithm:
// n.stabilize()
//   x = successor.predecessor;
//	 if (x is_in (n; successor))
//	 successor = x;
//	 successor.notify(n);
func (chordNode *ChordNode) StabilizeDaemon() error {

	chordNode.successorsList.RLock()
	successor := chordNode.successorsList.list[0]
	chordNode.successorsList.RUnlock()

	succStub, err := chordNode.getStubFor(ipAddr(successor.Ip))
	if err != nil {
		return err
	}

	x, err := succStub.GetPredecessor(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}

	if x.Id > chordNode.node.Id && x.Id < successor.Id {
		chordNode.successorsList.Lock()
		// Remove last element from successors, add first the new successor
		chordNode.successorsList.list = append([]*Node{successor}, chordNode.successorsList.list[:len(chordNode.successorsList.list)-1]...)
		chordNode.successorsList.Unlock()
	}

	_, err = succStub.Notify(context.Background(), &chordNode.node)
	if err != nil {
		return err
	}

	return nil
}

// called periodically. refreshes finger table entries.
// next stores the index of the next finger to fix.
// Algorithm:
//n.fix_fingers()
//	next = next + 1;
//	if (next > m)
//		next = 1;
//	finger[next] = find_successor(n + 2^(next-1) );
func (chordNode *ChordNode) FixFingersDaemon() error {
	nextCopy := 0
	chordNode.next.Lock()
	chordNode.next.value += 1
	if chordNode.next.value > chordNode.node.M {
		chordNode.next.value = 1
	}
	nextCopy = int(chordNode.next.value - 1)
	chordNode.next.Unlock()

	chordNode.successorsList.RLock()
	successorIp := chordNode.successorsList.list[0].Ip
	chordNode.successorsList.RUnlock()
	succStub, err := chordNode.getStubFor(ipAddr(successorIp))
	if err != nil {
		return err
	}

	succ, err := succStub.FindSuccessor(context.Background(), &ID{Id: chordNode.node.Id  + 1 << nextCopy})
	if err != nil {
		return err
	}

	chordNode.fingerTable.Lock()
	// nextCopy - 1 because we start index counting from 0
	chordNode.fingerTable.table[nextCopy - 1] = *succ
	chordNode.fingerTable.Unlock()

	return nil
}

// to clear the node's predecessor pointer if the predecessor has failed
func (chordNode *ChordNode) CheckPredecessorsDaemon() {

}
