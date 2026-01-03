package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/haedarrfd/simple-rss-aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// Configuration struct for API handlers
type apiConfig struct {
	DB *database.Queries `json:"db" toml:"db" yaml:"db" xml:"db"`
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}
	dbURL := os.Getenv("DB_URL")
	
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	// Open a connection to postgresql
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	// Allowing to interact with the database through pre-defined methods
	dbQueries := database.New(conn)
	apiCfg := apiConfig{
		DB: dbQueries,
	}

	// Start the feed scraping process
	go startScraping(dbQueries, 10, time.Hour)

	router := chi.NewRouter()

	// Set up CORS (Cross-Origin Resource Sharing) middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Create sub-path route, then the route is handled by handler function
	v1Router := chi.NewRouter()
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	// Mount sub-router in the main router
	router.Mount("/v1", v1Router)

	// Creates a new HTTP server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v\n", port)

	// Start the server
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
