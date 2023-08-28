package main

import (
	"os"
	"net/http"
	"github.com/gorilla/handlers"
	"FDBackend/internal/recommendations"
	"fmt"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
			port = "8080" // default port
	}
	fmt.Println("Listening on port:", port)
	router := http.NewServeMux()
	router.HandleFunc("/recommend", recommendations.Handler)

	corsMiddleware := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))
	http.ListenAndServe(":"+port, corsMiddleware(router))
}
