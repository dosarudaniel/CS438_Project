package chord

// StabilizeDaemon() is called periodically.
// It verifies n’s immediate successor, and tells the successor about n.
// Algorithm:
// n.stabilize()
//   x = successor.predecessor;
//	 if (x is_in (n; successor))
//	 successor = x;
//	 successor.notify(n);
func (chordNode *ChordNode) StabilizeDaemon() {
	// TODO Implement
}

// TODO implement the rest of daemons
