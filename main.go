package main

import (
	"log"
	"os"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/PratikforCoding/chirpy.git/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB *database.DB
	jwt_secret string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apicfg := apiConfig {
		fileserverHits: 0,
		DB: db,
		jwt_secret: jwtSecret,
	}

	router := chi.NewRouter()
	fsHandler := apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	router.Handle("/app/", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apicfg.handlerReset)
	apiRouter.Post("/chirps", apicfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apicfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{id}", apicfg.handlerGetChirpbyId)
	apiRouter.Post("/users", apicfg.handlerUsersCreate)
	apiRouter.Post("/login", apicfg.handlerLogin)
	apiRouter.Put("/users", apicfg.handlerUsersUpdate)

	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrcis", apicfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr: ":8080",
		Handler: corsMux,
	}

	log.Println("Server is running on port : 8080...")
	log.Fatal(server.ListenAndServe())
}