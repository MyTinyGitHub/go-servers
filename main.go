package main

import "net/http"

func serveHTTPHealthz(res http.ResponseWriter, req *http.Request) {
  res.WriteHeader(200)
  res.Header().Add("Content-Type", "text/plain")
  res.Write([]byte("OK"))
}


func main() {
  mux := http.NewServeMux()
  server := &http.Server{
    Handler: mux,
    Addr: ":8080",
  }

  mux.Handle("/app/", http.FileServer(http.Dir(".")))
  mux.HandleFunc("/healthz", serveHTTPHealthz)
  
  server.ListenAndServe()
}
