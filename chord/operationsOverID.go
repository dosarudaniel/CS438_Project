package chord

import (
	"encoding/hex"
	"math/big"
)

// idToBigIntString returns an integer represenation of an ID,
// where ID is a hex string, meaning each char is between 0-9 or a-f
func idToBigIntString(id string) string {
	y, err := hex.DecodeString(id)
	if err != nil {
		return ""
	}
	z := new(big.Int)
	z.SetBytes(y)
	return z.String()
}

func getFingerStart(id string, ithFinger, numOfBitsInID int) (string, error) {
	y, err := hex.DecodeString(id)
	if err != nil {
		return "", err
	}
	z := new(big.Int)
	z.SetBytes(y)
	z.Add(big.NewInt(int64(1<<ithFinger)), z) //TODO check that converting ithFinger to int64 is safe
	z = z.Mod(z, big.NewInt(int64(1<<numOfBitsInID)))
	return hex.EncodeToString(z.Bytes()), nil
}

/*
  For the Chord ring:
	m is_in (l, r) is equivalent to
		when r == l
			empty interval => false
		when l < r
			l < m && m < r => true
		when l > r
			l < m || m < r => true
		else => false
*/
func isBetweenTwoNodesExclusive(leftmostID, middleID, rightmostID string) bool {
	l := leftmostID
	m := middleID
	r := rightmostID

	switch {
	case l == r: // because interval is exclusive, for which l == r means, it's essentially empty
		return false
	case r > l && (l < m && m < r):
		return true
	case r < l && (l < m || m < r):
		return true
	default:
		return false
	}
}

/*
  For the Chord ring:
	m is_in (l, r] is equivalent to
		when r == l
			m == r => true
		when l < r
			l < m && m <= r => true
		when l > r
			l < m || m <= r => true
		else => false
*/
func isBetweenTwoNodesRightInclusive(leftmostID, middleID, rightmostID string) bool {
	l := leftmostID
	m := middleID
	r := rightmostID

	switch {
	case l == r && m == r: // because interval is exclusive, for which l == r means, it's essentially empty
		return true
	case r > l && (l < m && m <= r):
		return true
	case r < l && (l < m || m <= r):
		return true
	default:
		return false
	}
}
