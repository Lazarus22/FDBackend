package main

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"os"
	"FDBackend/internal/recommendations" // Adjust the path to match your project structure
)

func main() {
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
	http.ListenAndServe(":"+port, corsMiddleware(router))
}
