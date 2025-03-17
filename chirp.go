package main

import (
  "net/http"
  "io"
  "encoding/json"
  "strings"
)

func validateChirp(res http.ResponseWriter, req *http.Request) {
	type error_value struct {
		Error string `json:"error"`
	}

	body, error := io.ReadAll(req.Body)
	if error != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Add("Content-Type", "text/plain")
		dat, _ := json.Marshal(error_value{Error: "Something went wrong"})
		res.Write(dat)
		return
	}

	type requestBody struct {
		Body string `json:"body"`
	}

	var rBody requestBody
	json.Unmarshal(body, &rBody)

	if len(rBody.Body) > 140 {
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Add("Content-Type", "text/plain")
		dat, _ := json.Marshal(error_value{Error: "Chirp is too long"})
		res.Write(dat)
		return
	}

	type resBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	stringBody := strings.ReplaceAll(rBody.Body, "kerfuffle", "****")
	stringBody = strings.ReplaceAll(stringBody, "Kerfuffle", "****")
	stringBody = strings.ReplaceAll(stringBody, "sharbert", "****")
	stringBody = strings.ReplaceAll(stringBody, "Sharbert", "****")
	stringBody = strings.ReplaceAll(stringBody, "fornax", "****")
	stringBody = strings.ReplaceAll(stringBody, "Fornax", "****")

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "text/plain")
	dat, _ := json.Marshal(resBody{CleanedBody: stringBody})
	res.Write(dat)
}

