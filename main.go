package main

import (
	"fmt"
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
  fmt.Fprintf(res, "Hits: %v", value)
}

func (cfg *apiConfig) resetMetrics(res http.ResponseWriter, req *http.Request) {
  res.WriteHeader(200)
  res.Header().Add("Content-Type", "text/plain")
  cfg.fileserverHits.Store(0)
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
  mux.HandleFunc("GET /healthz", serveHTTPHealthz)
  mux.HandleFunc("GET /metrics", conf.serveMetrics)
  mux.HandleFunc("POST /reset", conf.resetMetrics)
  
  server.ListenAndServe()
}
