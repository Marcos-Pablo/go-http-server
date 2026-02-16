package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Marcos-Pablo/go-http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db             *sql.DB
	queries        *database.Queries
	fileserverHits atomic.Int32
	platform       string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	apiCfg := loadAPIConfig()
	defer apiCfg.db.Close()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func loadAPIConfig() *apiConfig {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("couldn't open connection to the database: %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		db:             db,
		queries:        dbQueries,
		fileserverHits: atomic.Int32{},
		platform:       platform,
	}

	return apiCfg
}
