package main

import (
	"encoding/json"
	"fmt"
	"go-servers/internal/auth"
	"go-servers/internal/database"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type chirp struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
}

type chirpData struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

func (cfg *apiConfig) getChirps(res http.ResponseWriter, req *http.Request) {
	chirps, _ := cfg.dabaseQueries.GetAllChirps(req.Context())

	var chirpss []chirp

	for _, chirpep := range chirps {
		chirpss = append(chirpss, chirp{
			Id:        chirpep.ID.String(),
			CreatedAt: chirpep.CreatedAt.String(),
			UpdatedAt: chirpep.UpdatedAt.String(),
			Body:      chirpep.Body,
			UserId:    chirpep.UserID.String(),
		})
	}

	data, _ := json.Marshal(chirpss)

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func (cfg *apiConfig) deleteChirp(res http.ResponseWriter, req *http.Request) {
	uu_id, _ := uuid.Parse(req.PathValue("chirpId"))
	token, err := auth.GetBearerToken(&req.Header)
	if err != nil {
		msg := fmt.Sprintf("unable to get access token: %v", err)
		respondWithError(msg, http.StatusUnauthorized, res)
		return
	}

	userId, err := auth.ValidateJWT(token, "TOP")
	if err != nil {
		msg := fmt.Sprintf("unable to validate the token: %v", err)
		respondWithError(msg, http.StatusForbidden, res)
		return
	}

	ok, err := cfg.dabaseQueries.UserHasChirp(req.Context(), database.UserHasChirpParams{
		ID:     uu_id,
		UserID: userId,
	})

	if err != nil || !ok {
		msg := fmt.Sprintf("chirp not fonud, or missing autharization: %v", err)
		respondWithError(msg, http.StatusForbidden, res)
		return
	}

	err = cfg.dabaseQueries.DeleteChirpOfUser(req.Context(), database.DeleteChirpOfUserParams{
		ID:     uu_id,
		UserID: userId,
	})

	if err != nil {
		msg := fmt.Sprintf("unable to delete a chirp: %v ", err)
		respondWithError(msg, http.StatusNotFound, res)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) getChirpById(res http.ResponseWriter, req *http.Request) {
	uu_id, _ := uuid.Parse(req.PathValue("chirpId"))
	chirpik, err := cfg.dabaseQueries.GetChirpById(req.Context(), uu_id)
	if err != nil {
		msg := fmt.Sprintf("unable to find a chirp: %v ", err)
		respondWithError(msg, http.StatusNotFound, res)
		return
	}
	chirping := chirp{
		Id:        chirpik.ID.String(),
		CreatedAt: chirpik.CreatedAt.String(),
		UpdatedAt: chirpik.UpdatedAt.String(),
		Body:      chirpik.Body,
		UserId:    chirpik.UserID.String(),
	}

	data, _ := json.Marshal(chirping)

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func (cfg *apiConfig) addChirp(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError("Something went wrong", http.StatusBadRequest, res)
		return
	}

	var chirpData chirpData
	err = json.Unmarshal(body, &chirpData)
	if err != nil {
		respondWithError("Unable to unmarshal data", http.StatusBadRequest, res)
		return
	}

	bearer, err := auth.GetBearerToken(&req.Header)
	if err != nil {
		respondWithError("Bearer token not provided", http.StatusUnauthorized, res)
		return
	}

	userId, err := auth.ValidateJWT(bearer, "TOP")
	if err != nil {
		respondWithError("Bearer token not provided", http.StatusUnauthorized, res)
		return
	}

	if ok := validateChirp(chirpData.Body, res); !ok {
		return
	}

	dbChirp, err := cfg.dabaseQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		UserID: userId,
		Body:   chirpData.Body,
	})

	if err != nil {
		respondWithError("Error creating chirp: "+err.Error(), http.StatusBadRequest, res)
		return
	}

	userData := chirp{
		Id:        dbChirp.ID.String(),
		CreatedAt: dbChirp.CreatedAt.String(),
		UpdatedAt: dbChirp.UpdatedAt.String(),
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID.String(),
	}

	res.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(userData)
	if err != nil {
		respondWithError("Error Marshalling created user", http.StatusBadRequest, res)
		return
	}

	res.Write(data)
}

func validateChirp(message string, res http.ResponseWriter) bool {

	if len(message) > 140 {
		respondWithError("Chirp is too long", http.StatusBadRequest, res)
		return false
	}

	return true
}
