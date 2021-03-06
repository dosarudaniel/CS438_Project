package chord

import (
	"crypto"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// In this file, we put all helper functions and types used by Chord package
// "helper" can be defined as something that is not worth own file/package and has tiny role in the chord functionality
// DANGEROUS: assumes chordNode.config.NumOfBitsInID is set => must not be called before chordNode.config is set
func (chordNode *ChordNode) hashString(s string) (string, error) {
	numOfBitsInID := chordNode.config.NumOfBitsInID
	sha256 := crypto.SHA256.New()
	_, err := sha256.Write([]byte(s))
	if err != nil {
		return "", err
	}
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash)[64-(numOfBitsInID/4):], nil
}

func (chordNode *ChordNode) logDebugPretty() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	log.Debug(chordNode.String())
}

// cleanFilename removes frequent punctuation marks and " and " and " of "
func cleanFilename(filename string) string {
	cleanFilename := filename
	for _, char := range []string{".", ",", "?", "!", " of ", " and ", " a ", " the "} {
		cleanFilename = strings.ReplaceAll(cleanFilename, char, " ")
	}
	cleanFilename = strings.ToLower(cleanFilename)
	return cleanFilename
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
