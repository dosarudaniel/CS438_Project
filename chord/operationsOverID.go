package chord

import (
	"encoding/hex"
	"math/big"
)

func addIntToID(x int, id string) (string, error) {
	y, err := hex.DecodeString(id)
	if err != nil {
		return "", err
	}
	z := new(big.Int)
	z.SetBytes(y)
	z.Add(big.NewInt(int64(x)), z) //TODO check that converting x to int64 is safe
	return hex.EncodeToString(z.Bytes()), nil
}
