package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAuthenticationToken(headers *http.Header, token string) (string, error) {
	h := headers.Get("Authorization")
	if !strings.Contains(h, token) {
		return "", fmt.Errorf("header does not contain authorization")
	}

	h = strings.Replace(h, token, "", 1)
	h = strings.Trim(h, " ")

	return h, nil
}
