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
