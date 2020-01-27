package chord

import (
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
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
	fmt.Println("ID===", hex.EncodeToString(hash))
	return hex.EncodeToString(hash)[64-(numOfBitsInID/4):], nil
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
func isBetweenTwoNodesExclusive(leftmostNode, nodeBetween, rightmostNode Node) bool {
	l := leftmostNode.Id
	m := nodeBetween.Id
	r := rightmostNode.Id

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
func isBetweenTwoNodesRightInclusive(leftmostNode, nodeBetween, rightmostNode Node) bool {
	l := leftmostNode.Id
	m := nodeBetween.Id
	r := rightmostNode.Id

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

func nilError(s string) error {
	return errors.New(s + " is nil")
}

type nilSuccessor struct{}

func (m *nilSuccessor) Error() string {
	return "successor is nil"
}

type nilPredecessor struct{}

func (m *nilPredecessor) Error() string {
	return "predecessor is nil"
}

type nilNode struct{}

func (m *nilNode) Error() string {
	return "nodeIsNil"
}
