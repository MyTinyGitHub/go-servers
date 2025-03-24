package auth

import (
	"encoding/hex"
	mRand "math/rand"
	"strconv"
)

func MakeRefreshToken() (string, error) {
	value := mRand.Int()

	res := hex.EncodeToString([]byte(strconv.Itoa(value)))

	return res, nil
}
