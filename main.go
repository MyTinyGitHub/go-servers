package main

import (
	"database/sql"
	"encoding/json"
	"go-servers/internal/database"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dabaseQueries  *database.Queries
	JwtSecret      string
	PolkaKey       string
}

func serveHTTPHealthz(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "text/plain")
	res.Write([]byte("OK"))
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	var conf apiConfig
	conf.dabaseQueries = dbQueries
	conf.JwtSecret = jwtSecret
	conf.PolkaKey = polkaKey

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.HandleFunc("/app/", conf.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", serveHTTPHealthz)
	mux.HandleFunc("GET /admin/metrics", conf.serveMetrics)
	mux.HandleFunc("POST /admin/reset", conf.resetMetrics)

	mux.HandleFunc("POST /api/polka/webhooks", conf.polkaWebhook)
	mux.HandleFunc("POST /api/users", conf.addUser)
	mux.HandleFunc("PUT /api/users", conf.updateUser)

	mux.HandleFunc("POST /api/login", conf.login)
	mux.HandleFunc("POST /api/refresh", conf.refresh)
	mux.HandleFunc("POST /api/revoke", conf.revoke)

	mux.HandleFunc("POST /api/chirps", conf.addChirp)
	mux.HandleFunc("GET /api/chirps", conf.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", conf.getChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", conf.deleteChirp)

	server.ListenAndServe()
}

type error_value struct {
	Error string `json:"error"`
}

func respondWithError(message string, httpStatus int, res http.ResponseWriter) {
	res.WriteHeader(httpStatus)
	res.Header().Add("Content-Type", "text/plain")
	dat, _ := json.Marshal(error_value{Error: message})
	res.Write(dat)
}
