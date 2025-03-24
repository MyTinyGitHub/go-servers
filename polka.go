package main

import (
	"encoding/json"
	"fmt"
	"go-servers/internal/auth"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhook(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

  token, err := auth.GetAuthenticationToken(&req.Header, "ApiKey")
  if err != nil || token != cfg.PolkaKey {
		respondWithError("Unable to get token from header", http.StatusUnauthorized, res)
		return
  }

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError("Unable to parse the body", http.StatusInternalServerError, res)
		return
	}

	var input parameters
	err = json.Unmarshal(body, &input)
	if err != nil {
		respondWithError("Unable to unmarshal the body", http.StatusInternalServerError, res)
		return
	}

	if input.Event != "user.upgraded" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(input.Data.UserId)
	if err != nil {
		respondWithError("Unable to parse the userid", http.StatusInternalServerError, res)
		return
	}

	err = cfg.dabaseQueries.SetChirpToRed(req.Context(), userId)
	if err != nil {
		msg := fmt.Sprintf("User not found: %v", err)
		respondWithError(msg, http.StatusNotFound, res)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
