package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
) 


func (cfg *apiConfig) createUser(res http.ResponseWriter, req *http.Request) {
	type error_value struct {
		Error string `json:"error"`
	}

  type email struct {
    Email string `json:"email"`
  }

	body, err := io.ReadAll(req.Body)
	if err != nil {
    badRequest("Something went wrong", res)
		return
	}

  var userEmail email
  err = json.Unmarshal(body, &userEmail)
	if err != nil {
    badRequest("Unable to unmarshal data", res)
		return
	}


  value := sql.NullString{
    String:userEmail.Email,
    Valid: true,
  }

  user, err := cfg.dabaseQueries.CreateUser(req.Context(), value)
	if err != nil {
    badRequest("Error creating user" + err.Error(), res)
		return
	}

  type userType struct {
    Id string `json:"id"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    Email string `json:"email"`
  }

  userData := userType{
    Id: user.ID.UUID.String(),
    CreatedAt: user.CreatedAt.Time.String(),
    UpdatedAt: user.UpdatedAt.Time.String(),
    Email: user.Email.String,
  }

  res.WriteHeader(http.StatusCreated)
  data, err := json.Marshal(userData)
	if err != nil {
    badRequest("Error Marshalling created user", res)
		return
	}

  res.Write(data)
}
