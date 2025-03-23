package main

import (
	"encoding/json"
	"fmt"
	"go-servers/internal/auth"
	"net/http"
)

func (cfg *apiConfig) refresh(res http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(&req.Header)
	if err != nil {
		respondWithError("unable to get bearer token", http.StatusUnauthorized, res)
		return
	}

	dbToken, err := cfg.dabaseQueries.GetRefreshTokenByToken(req.Context(), token)
	if err != nil {
		msg := fmt.Sprintf("refresh token does not exists or is expired: %v", err)
		respondWithError(msg, http.StatusUnauthorized, res)
		return
	}

	err = cfg.dabaseQueries.UpdateExpiresAtForRevoked(req.Context(), dbToken.Token)
	if err != nil {
		msg := fmt.Sprintf("refresh token does not exists or is expired: %v", err)
		respondWithError(msg, http.StatusUnauthorized, res)
		return
	}

	user, err := cfg.dabaseQueries.GetUserByRefreshToken(req.Context(), dbToken.Token)
	if err != nil {
		respondWithError("unable to get user from token", http.StatusInternalServerError, res)
		return
	}

	newJwtToken, err := auth.MakeJWT(user.ID, "TOP")
	if err != nil {
		respondWithError("unable to create JWT token", http.StatusInternalServerError, res)
		return
	}

	rToken := struct {
		Token string `json:"token"`
	}{
		Token: newJwtToken,
	}

	response, _ := json.Marshal(rToken)

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (cfg *apiConfig) revoke(res http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(&req.Header)
	if err != nil {
		respondWithError("unable to get bearer token", http.StatusUnauthorized, res)
		return
	}

	cfg.dabaseQueries.RevokeTokenByToken(req.Context(), token)

	res.WriteHeader(http.StatusNoContent)
}
