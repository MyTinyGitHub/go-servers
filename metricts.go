package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	}
}

func (cfg *apiConfig) serveMetrics(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "text/plain")

	value := cfg.fileserverHits.Load()

	page := `
    <html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
      </body>
    </html>
  `

	fmt.Fprintf(res, page, value)
}

func (cfg *apiConfig) resetMetrics(res http.ResponseWriter, req *http.Request) {
  platform := os.Getenv("PLATFORM")

  if platform != "dev" {
    res.WriteHeader(http.StatusForbidden)
    res.Header().Add("Content-Type", "text/plain")
    return
  }

  cfg.dabaseQueries.DeleteUsers(req.Context())
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "text/plain")
}
