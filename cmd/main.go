package main

import (
	"os"
	"net/http"
	"github.com/gorilla/handlers"
	"FDBackend/internal/recommendations"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // localhost
	}
	router := http.NewServeMux()
	router.HandleFunc("/recommend", recommendations.Handler)

	corsMiddleware := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))
	http.ListenAndServe(":"+port, corsMiddleware(router))
}
