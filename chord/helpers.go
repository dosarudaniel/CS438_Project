package chord

import (
	"crypto"
	"encoding/hex"
)

// In this file, we put all helper functions and types used by Chord package
// "helper" can be defined as something that is not worth own file/package and has tiny role in the chord functionality

func hashString(s string) (string, error) {
	sha256 := crypto.SHA256.New()
	_, err := sha256.Write([]byte(s))
	if err != nil {
		return "", err
	}
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash), nil
}
