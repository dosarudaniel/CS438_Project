package chord

import (
	"encoding/hex"
	"math/big"
)

func idToInt(id string) uint64 {
	y, err := hex.DecodeString(id)
	if err != nil {
		return 0
	}
	z := new(big.Int)
	z.SetBytes(y)
	return z.Uint64()
}

func getIthFingerID(id string, ithFinger, numOfBitsInID int) (string, error) {
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
