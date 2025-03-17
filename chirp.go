package main

import (
	"encoding/json"
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
      Id: chirpep.ID.UUID.String(),
      CreatedAt: chirpep.CreatedAt.String(),
      UpdatedAt: chirpep.UpdatedAt.String(),
      Body: chirpep.Body,
      UserId: chirpep.UserID.UUID.String(),
    })
  }

	data, _ := json.Marshal(chirpss)

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func (cfg *apiConfig) getChirpById(res http.ResponseWriter, req *http.Request) {
  uu_id, _ := uuid.Parse(req.PathValue("chirpId"))
	chirpik, _ := cfg.dabaseQueries.GetChirpById(req.Context(), uuid.NullUUID{
    UUID: uu_id,
    Valid: true,
  })

  chirping := chirp{
      Id: chirpik.ID.UUID.String(),
      CreatedAt: chirpik.CreatedAt.String(),
      UpdatedAt: chirpik.UpdatedAt.String(),
      Body: chirpik.Body,
      UserId: chirpik.UserID.UUID.String(),
  }
  

	data, _ := json.Marshal(chirping)

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func (cfg *apiConfig) addChirp(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError("Something went wrong",http.StatusBadRequest, res)
		return
	}

	var chirpData chirpData
	err = json.Unmarshal(body, &chirpData)
	if err != nil {
		respondWithError("Unable to unmarshal data", http.StatusBadRequest, res)
		return
	}

	if ok := validateChirp(chirpData.Body, res); !ok {
		return
	}

	uu_id, _ := uuid.Parse(chirpData.UserId)

	userId := uuid.NullUUID{
		UUID:  uu_id,
		Valid: true,
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
		Id:        dbChirp.ID.UUID.String(),
		CreatedAt: dbChirp.CreatedAt.String(),
		UpdatedAt: dbChirp.UpdatedAt.String(),
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID.UUID.String(),
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
