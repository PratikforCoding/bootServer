package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}



func validateHandler(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type returnVal struct {
		Valid bool `json:"valid"`
	}
	type newChirp struct {
		CleanBod string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w , http.StatusInternalServerError, "could't decode parameters")
		return 
	}
	if len(params.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	} 
	/*responseWithJson(w, http.StatusOK, returnVal{
		Valid: true,
	})*/
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	dat := strings.ToLower(params.Body)
	dat2 := strings.Split(params.Body, " ")
	chirpSlice := strings.Split(dat, " ")
	for i, word := range chirpSlice {
		for _, word1 := range profane {
			if word == word1 {
				chirpSlice[i] = "****"
			} 
		}
	}
	for i, word := range chirpSlice {
		if word != "****" && word != dat2[i] {
			chirpSlice[i] = dat2[i]
		}
	}
	new := strings.Join(chirpSlice, " ")
	responseWithJson(w, http.StatusOK, newChirp{
		CleanBod: new,
	})

}


func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499{
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	responseWithJson(w, code, errorResponse{
		Error: msg,
	})
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return 
	}
	w.WriteHeader(code)
	w.Write(dat)
}
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

func assetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := `<pre>
	<a href="logo.png">logo.png</a>
	</pre>`

	_, err := w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Failed to write resoponse!", http.StatusInternalServerError)
		return
	}
}
func (apiCfg *apiConfig)metricsHandler(w http.ResponseWriter, r *http.Request) {
	
	hits := apiCfg.fileserverHits
	response := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	
}
func main() {
	r := chi.NewRouter()
	rt := chi.NewRouter()
	corsMux := middlewareCors(r)
	rt.Use(middlewareCors)
	apiCfg := &apiConfig{}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)
	r.Get("/app/assets", assetHandler)
	rt.Get("/healthz", healthzHandler)
	rt.Get("/metrics", apiCfg.metricsHandler)
	rt.Post("/validate_chirp", validateHandler)
	r.Mount("/api", rt)

	server := &http.Server{
		Addr: ":8080",
		Handler: corsMux,
	}

	log.Println("Server is running on port : 8080...")
	log.Fatal(server.ListenAndServe())
}