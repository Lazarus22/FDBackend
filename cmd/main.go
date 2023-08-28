package main

import (
	"github.com/gorilla/handlers"
	"net/http"
	"FlavorDB/backend/internal/recommendations"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/recommend", recommendations.Handler)

	// Apply CORS middleware
	corsMiddleware := handlers.CORS(handlers.AllowedOrigins([]string{"http://localhost:3000"}))
	http.ListenAndServe(":8080", corsMiddleware(router))
}