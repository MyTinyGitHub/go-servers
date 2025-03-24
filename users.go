package main

import (
	"encoding/json"
	"fmt"
	"go-servers/internal/auth"
	"go-servers/internal/database"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) updateUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(&req.Header)
	if err != nil {
		msg := fmt.Sprintf("unable to get access token: %v", err)
		respondWithError(msg, http.StatusUnauthorized, res)
		return
	}

	userId, err := auth.ValidateJWT(token, "TOP")
	if err != nil {
		msg := fmt.Sprintf("unable to validate the token: %v", err)
		respondWithError(msg, http.StatusUnauthorized, res)
		return
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

	user, err := cfg.dabaseQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          input.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		msg := fmt.Sprintf("unable to update user: %v", err)
		respondWithError(msg, http.StatusInternalServerError, res)
		return
	}

	type userType struct {
		Id          string `json:"id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	userData := userType{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt.Time.String(),
		UpdatedAt:   user.UpdatedAt.Time.String(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	data, err := json.Marshal(userData)
	if err != nil {
		respondWithError("Error Marshalling created user", http.StatusBadRequest, res)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func (cfg *apiConfig) addUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
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

	value := database.CreateUserParams{
		Email:          input.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.dabaseQueries.CreateUser(req.Context(), value)
	if err != nil {
		respondWithError("Error creating user"+err.Error(), http.StatusBadRequest, res)
		return
	}

	type userType struct {
		Id          string `json:"id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	userData := userType{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt.Time.String(),
		UpdatedAt:   user.UpdatedAt.Time.String(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
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
		Email    string `json:"email"`
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

	user, err := cfg.dabaseQueries.GetUserByEmail(req.Context(), input.Email)
	if err != nil {
		respondWithError("Error retrieving user"+err.Error(), http.StatusBadRequest, res)
		return
	}

	if ok := checkPasswordHash(input.Password, user.HashedPassword); ok != nil {
		respondWithError("Cannot authenticate user", http.StatusUnauthorized, res)
		return
	}

	token, err := auth.MakeJWT(user.ID, "TOP")
	if err != nil {
		respondWithError("Unable to create jwt token", http.StatusInternalServerError, res)
		return
	}

	type userType struct {
		Id           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}

	rtoken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError("Unable to create refresh token", http.StatusInternalServerError, res)
		return
	}

	_, err = cfg.dabaseQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:  rtoken,
		UserID: user.ID,
	})

	if err != nil {
		msg := fmt.Sprintf("unable to create refresh token: %v", err)
		respondWithError(msg, http.StatusInternalServerError, res)
		return
	}

	userData := userType{
		Id:           user.ID.String(),
		CreatedAt:    user.CreatedAt.Time.String(),
		UpdatedAt:    user.UpdatedAt.Time.String(),
		Email:        user.Email,
		Token:        token,
		RefreshToken: rtoken,
		IsChirpyRed:  user.IsChirpyRed,
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
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func hashPassword(password string) (string, error) {
	value, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(value), err
}
