package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers *http.Header) (string, error) {
	h := headers.Get("Authorization")
	if !strings.Contains(h, "Bearer") {
		return "", fmt.Errorf("header does not contain authorization")
	}

	h = strings.Replace(h, "Bearer", "", 1)
	h = strings.Trim(h, " ")

	return h, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: &jwt.NumericDate{
			Time: time.Now().UTC(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().UTC().Add(time.Hour),
		},
		Subject: userID.String(),
	}).SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to parse claim: %v", err)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to extract subject claim: %v", err)
	}

	return uuid.Parse(sub)
}
