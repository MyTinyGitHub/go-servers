package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
  fileserverHits atomic.Int32
} 


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    cfg.fileserverHits.Add(1)
    next.ServeHTTP(res, req)
  }
}

func serveHTTPHealthz(res http.ResponseWriter, req *http.Request) {
  res.WriteHeader(200)
  res.Header().Add("Content-Type", "text/plain")
  res.Write([]byte("OK"))
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
  res.WriteHeader(200)
  res.Header().Add("Content-Type", "text/plain")
  cfg.fileserverHits.Store(0)
}

func validateChirp(res http.ResponseWriter, req *http.Request) {
  body, error := io.ReadAll(req.Body)
  if error != nil {
    res.WriteHeader(http.StatusBadRequest)
    res.Header().Add("Content-Type", "text/plain")
    fmt.Fprintf(res, `
      {
        "error": "Something went wrong"
      }`)
    return 
  }

  stringBody := string(body)
  if len(stringBody) > 140 {
    res.WriteHeader(http.StatusBadRequest)
    res.Header().Add("Content-Type", "text/plain")
    fmt.Fprintf(res, `
      {
        "error": "Chirp is too long"
      }`)
    return 
  }

    res.WriteHeader(http.StatusOK)
    res.Header().Add("Content-Type", "text/plain")
    fmt.Fprintf(res, `
      {
        "valid": true
      }
      `)

}

func main() {
  mux := http.NewServeMux()

  var conf apiConfig 

  server := &http.Server{
    Handler: mux,
    Addr: ":8080",
  }

  handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
  
  mux.HandleFunc("/app/", conf.middlewareMetricsInc(handler))
  mux.HandleFunc("GET /api/healthz", serveHTTPHealthz)
  mux.HandleFunc("GET /admin/metrics", conf.serveMetrics)
  mux.HandleFunc("POST /admin/reset", conf.resetMetrics)
  mux.HandleFunc("POST /api/validate_chirp", validateChirp)
  
  server.ListenAndServe()
}
