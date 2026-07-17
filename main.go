package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/arjunsingh14/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	queries *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect to datbase: %s", err)
	}
	dbQueries := database.New(db)

	filePathRoot := "."
	port := "8080"
	apiConfig := &apiConfig{queries: dbQueries}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", apiConfig.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))

	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetMetrics)

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)


	server := &http.Server{
		Handler: mux, 
		Addr: ":" + port,
	}

	server.ListenAndServe()

}
