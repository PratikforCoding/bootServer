package main

import (
	"log"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/app/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/logo.png")
	})
	mux.HandleFunc("/healthz", healthzHandler)

	server := &http.Server{
		Addr: ":8080",
		Handler: corsMux,
	}

	log.Println("Server is running on port : 8080...")
	log.Fatal(server.ListenAndServe())
}