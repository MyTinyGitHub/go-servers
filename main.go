package main

import (
	"database/sql"
	"go-servers/internal/database"
	"net/http"
  "encoding/json"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dabaseQueries  *database.Queries
}

func serveHTTPHealthz(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "text/plain")
	res.Write([]byte("OK"))
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	var conf apiConfig
	conf.dabaseQueries = dbQueries

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.HandleFunc("/app/", conf.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", serveHTTPHealthz)
	mux.HandleFunc("GET /admin/metrics", conf.serveMetrics)
	mux.HandleFunc("POST /admin/reset", conf.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/users", conf.createUser)

	server.ListenAndServe()
}

type error_value struct {
  Error string `json:"error"`
}

func badRequest(message string, res http.ResponseWriter) {
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Add("Content-Type", "text/plain")
		dat, _ := json.Marshal(error_value{Error: message})
		res.Write(dat)
}
