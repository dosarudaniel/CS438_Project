package file_sharing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Hash [32]byte

func ChunkToHash(chunk Chunk) Hash {
	return sha256.Sum256(chunk)
}

func (hash Hash) ToHex() string {
	return hex.EncodeToString(hash[:])
}

func BytesToHash(bytes []byte) (Hash, error) {
	var hash Hash
	if len(bytes) != 32 {
		return hash, fmt.Errorf("invalid hash format (not 256 bits)")
	}
	copy(hash[:], bytes)
	return hash, nil
}

func HexToHash(str string) (Hash, error) {
	var hash Hash
	if len(str) != 64 {
		return hash, fmt.Errorf("invalid hash format (not 256 bits)")
	}
	decoded, err := hex.DecodeString(str)
	if err != nil {
		return hash, err
	}
	copy(hash[:], decoded)
	return hash, nil
}
