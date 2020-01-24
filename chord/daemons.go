package chord

import (
	"context"
	"errors"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"time"
)

// StabilizeDaemon() is called periodically.
// It verifies n’s immediate successor, and tells the successor about n.
// Algorithm:
// n.stabilize()
//   x = successor.predecessor;
//	 if (x is_in (n; successor))
//	   successor = x;
//	 successor.notify(n);
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

// CheckPredecessorDaemon checks whether the node's predecessor has failed
// n.check_predecessor()
//	 if (predecessor has failed) <- in our case responds to a FindSuccessor rpc call within 3 seconds
//	   predecessor = nil;
func (chordNode *ChordNode) CheckPredecessorDaemon() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
