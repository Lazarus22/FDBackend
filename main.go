package main

import (
	"FDBackend/internal/recommendations"
	"FDBackend/cypherQueries"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"os"
	forcesslheroku "github.com/jonahgeorge/force-ssl-heroku"
)

func main() {
	// Initialize Queries
	err := cypherQueries.InitializeQueries()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}
	fmt.Println("Listening on port:", port)

	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	ctx := context.Background()

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	router := http.NewServeMux()
	router.HandleFunc("/recommendations", recommendations.NewHandler(driver))

	corsMiddleware := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	isProduction := os.Getenv("ENV") == "PRODUCTION"

	// Wrap your router with the force-ssl-heroku middleware
	secureRouter := forcesslheroku.ForceSsl(router)

	if isProduction {
		http.ListenAndServe(":"+port, corsMiddleware(secureRouter))
	} else {
		http.ListenAndServe(":"+port, corsMiddleware(router))
	}
}
