package chord

import (
	"context"
	"errors"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"time"
)

/*
 StabilizeDaemon() is called periodically.
 It verifies n’s immediate successor, and tells the successor about n.
 Algorithm:
 	n.stabilize()
  		x = successor.predecessor;
	 	if (x is_in (n; successor))
	   		successor = x;
	 	successor.notify(n);

*/
func (chordNode *ChordNode) StabilizeDaemon() {
	succ, doesExist := chordNode.getSuccessor()
	if !doesExist {
		return
	}

	// TODO what if successor is down
	x, err := chordNode.stubGetPredecessor(ipAddr(succ.Ip), context.Background())
	if err != nil || x == nil {
		return
	}

	if chordNode.node.Id < x.Id && x.Id < succ.Id {
		chordNode.setSuccessor(x)
	}

	_ = chordNode.stubNotify(ipAddr(x.Ip), context.Background(), &chordNode.node)
}

/*
 FixFingersDaemon() is called periodically, refreshes finger table entries.
 `next` stores the index of the next finger to fix.

 Algorithm:
 n.fix_fingers()
 	next = next + 1
 	if next >= m
 		next = 0
    finger[next] = find_successor(n + 2^next)
*/
func (chordNode *ChordNode) FixFingersDaemon(intervalMs int) {
	next := -1
	m := chordNode.config.NumOfBitsInID
	n := chordNode.node.Id

	select {
	case <-time.Tick(time.Duration(intervalMs) * time.Second):
		var err error

		next++
		if next >= int(m) {
			next = 0
		}

		// n + 2^next
		id, err := addIntToID(1<<next, n)
		if err != nil {
			return
		}

		succ, err := chordNode.FindSuccessor(context.Background(), &ID{Id: id})
		if err != nil {
			return
		}

		chordNode.fingerTable.Lock()
		chordNode.fingerTable.table[next] = succ
		chordNode.fingerTable.Unlock()
	}
}

/*
 CheckPredecessorDaemon checks whether the node's predecessor has failed
 Algorithm
 n.check_predecessor()
 	if (predecessor has failed) <- in our case responds to a FindSuccessor rpc call within 3 seconds
		predecessor = nil;
*/
func (chordNode *ChordNode) CheckPredecessorDaemon() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //TODO the timeout should be flexible
	defer cancel()

	pred, doesExist := chordNode.getPredecessor()
	if !doesExist {
		return
	}

	_, err = chordNode.stubFindSuccessor(ipAddr(pred.Ip), ctx, &ID{Id: chordNode.node.Ip})
	if errors.Is(err, context.DeadlineExceeded) {
		chordNode.setPredecessor(nil)
	}
}
