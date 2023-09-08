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
)

func echoEnvHandler(w http.ResponseWriter, r *http.Request) {
	neo4jURL := os.Getenv("NEO4J_URL")
	fmt.Fprintf(w, "NEO4J_URL: %s", neo4jURL)
}

func enforceHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "http" || proto == "" && r.TLS == nil {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/echo/env", echoEnvHandler)
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

	if isProduction {
		http.ListenAndServe(":"+port, corsMiddleware(enforceHTTPS(router)))
	} else {
		http.ListenAndServe(":"+port, corsMiddleware(router))
	}
}
