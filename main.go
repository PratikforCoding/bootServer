package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
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
	w.Header().Set("Content-Type", "text/plain: charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", hits)
	
}
func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	apiCfg := &apiConfig{}
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/app/assets", assetHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)

	server := &http.Server{
		Addr: ":8080",
		Handler: corsMux,
	}

	log.Println("Server is running on port : 8080...")
	log.Fatal(server.ListenAndServe())
}