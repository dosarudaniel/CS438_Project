package chord

import (
	"crypto"
	"encoding/hex"
	. "github.com/dosarudaniel/CS438_Project/services/chord_service"
)

// In this file, we put all helper functions and types used by Chord package
// "helper" can be defined as something that is not worth own file/package and has tiny role in the chord functionality

func hashString(s string, numOfBitsInID int) (string, error) {
	sha256 := crypto.SHA256.New()
	_, err := sha256.Write([]byte(s))
	if err != nil {
		return "", err
	}
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash[len(hash)-(numOfBitsInID/8):]), nil
}

/*
  For the Chord ring:
	m is_in (l, r) means l < m & m < r when r > l
	m is_in (l, r) means l < m         when r < l
*/
func isBetweenTwoNodesExclusive(leftmostNode, nodeBetween, rightmostNode Node) bool {
	l := leftmostNode.Id
	m := nodeBetween.Id
	r := rightmostNode.Id

	switch {
	case l == r: // because interval is exclusive, for which l == r means, it's essentially empty
		return false
	case r > l && l < m && m < r:
		return true
	case r < l && l < m:
		return true
	default:
		return false
	}
}
