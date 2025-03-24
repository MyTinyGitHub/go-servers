package auth

import (
  mRand "math/rand"
	"encoding/hex"
	"strconv"
)

func MakeRefreshToken() (string, error) {
	value := mRand.Int()

	res := hex.EncodeToString([]byte(strconv.Itoa(value)))

	return res, nil
}
