package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

func (cfg *apiConfig) addUser(res http.ResponseWriter, req *http.Request) {
	type email struct {
		Email string `json:"email"`
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError("Something went wrong", http.StatusBadRequest, res)
		return
	}

	var userEmail email
	err = json.Unmarshal(body, &userEmail)
	if err != nil {
		respondWithError("Unable to unmarshal data", http.StatusBadRequest, res)
		return
	}

	value := sql.NullString{
		String: userEmail.Email,
		Valid:  true,
	}

	user, err := cfg.dabaseQueries.CreateUser(req.Context(), value)
	if err != nil {
		respondWithError("Error creating user"+err.Error(), http.StatusBadRequest, res)
		return
	}

	type userType struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	userData := userType{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
		Email:     user.Email.String,
	}

	res.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(userData)
	if err != nil {
		respondWithError("Error Marshalling created user", http.StatusBadRequest, res)
		return
	}

	res.Write(data)
}
