package auth

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

func MakeRefreshToken() (string, error) {
	value, err := rand.Read([]byte("Test"))
	if err != nil {
		return "", err
	}

	res := hex.EncodeToString([]byte(strconv.Itoa(value)))

	return res, nil
}
