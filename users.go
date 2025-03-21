package main

import (
	"encoding/json"
	"go-servers/internal/auth"
	"go-servers/internal/database"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) addUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
    Password string `json:"password"`
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError("Something went wrong", http.StatusBadRequest, res)
		return
	}

	var input parameters
	err = json.Unmarshal(body, &input)
	if err != nil {
		respondWithError("Unable to unmarshal data", http.StatusBadRequest, res)
		return
	}

  hashedPassword, _ := hashPassword(input.Password)

	value := database.CreateUserParams {
    Email: input.Email,
    HashedPassword: hashedPassword,
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
		Email:     user.Email,
	}

	res.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(userData)
	if err != nil {
		respondWithError("Error Marshalling created user", http.StatusBadRequest, res)
		return
	}

	res.Write(data)
}

func (cfg *apiConfig) login(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
    Password string `json:"password"`
    ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError("Something went wrong", http.StatusBadRequest, res)
		return
	}

	var input parameters
	err = json.Unmarshal(body, &input)
	if err != nil {
		respondWithError("Unable to unmarshal data", http.StatusBadRequest, res)
		return
	}

	user, err := cfg.dabaseQueries.GetUserByEmail(req.Context(), input.Email)
	if err != nil {
		respondWithError("Error retrieving user"+err.Error(), http.StatusBadRequest, res)
		return
	}

  if ok := checkPasswordHash(input.Password, user.HashedPassword); ok != nil {
    respondWithError("Cannot authenticate user", http.StatusUnauthorized, res)
    return
  }


  if input.ExpiresInSeconds == 0 || input.ExpiresInSeconds > 60 {
    input.ExpiresInSeconds = 60
  }

  token, err := auth.MakeJWT(user.ID, "TOP", time.Duration(input.ExpiresInSeconds) * time.Second)
  if err != nil {
    respondWithError("Unable to create jwt token", http.StatusInternalServerError, res)
  }

	type userType struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
    Token string `json:"token"`
	}

	userData := userType{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
		Email:     user.Email,
    Token: token, 
	}

	res.WriteHeader(http.StatusOK)
	data, err := json.Marshal(userData)
	if err != nil {
		respondWithError("Error Marshalling created user", http.StatusBadRequest, res)
		return
	}

	res.Write(data)
  
}

func checkPasswordHash(password string, hash string) error {
  return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
}

func hashPassword(password string) (string, error) {
  value, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  return string(value), err
}
