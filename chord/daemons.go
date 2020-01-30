package chord

import (
	"context"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
	"time"
)

func (chordNode *ChordNode) RunAtInterval(daemon func(*ChordNode), intervalSeconds time.Duration) {
	daemon(chordNode)
	for {
		select {
		case <-time.Tick(intervalSeconds):
			daemon(chordNode)
		}
	}
}

/*
 StabilizeDaemon() verifies n’s immediate successor, and tells the successor about n.
 Algorithm:
 	n.stabilize()
  		x = successor.predecessor;
	 	if (x is_in (n; successor))
	   		successor = x;
	 	successor.notify(n);

*/
func StabilizeDaemon(chordNode *ChordNode) {
	succ, doesExist := chordNode.getSuccessor()
	if !doesExist {
		return
	}

	// TODO what if successor is down
	x, err := chordNode.stubGetPredecessor(ipAddr(succ.Ip), context.Background())
	if err == nil && x != nil {
		if isBetweenTwoNodesExclusive(chordNode.node.Id, x.Id, succ.Id) || chordNode.node.Id == succ.Id {
			chordNode.setSuccessor(x)
		}
	}

	if succ.Id != chordNode.node.Id {
		err = chordNode.stubNotify(ipAddr(succ.Ip), context.Background(), &chordNode.node)
	}
	log.Debug(chordNode)
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
func FixFingersDaemon(chordNode *ChordNode) func(*ChordNode) {
	next := -1
	m := chordNode.config.NumOfBitsInID
	n := chordNode.node.Id

	return func(chordNode *ChordNode) {
		var err error

		next = (next + 1) % m

		// n + 2^next % 2^m
		id, err := getFingerStart(n, next, m)
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

		log.Debug(chordNode)
	}
}

/*
 CheckPredecessorDaemon checks whether the node's predecessor has failed
 Algorithm
 n.check_predecessor()
 	if (predecessor has failed) <- in our case responds to a FindSuccessor rpc call within 3 seconds
		predecessor = nil;
 FIXME: implement correctly
*/
func CheckPredecessorDaemon(chordNode *ChordNode) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //TODO the timeout should be flexible
	defer cancel()

	pred, doesExist := chordNode.getPredecessor()
	if !doesExist {
		return
	}

	// transfer keys from the predecessor that have ID higher than
	// predecessor's ID
	err = chordNode.stubTransferKeys(ipAddr(pred.Ip), ctx, pred.Id, chordNode.node)

	log.WithField("err", err).Info("check predecessor daemon failed")
}
