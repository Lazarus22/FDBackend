package recommendations

import (
	"FDBackend/cypherQueries"
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"fmt"
)

type Pairing struct {
	Name     string `json:"name"`
	Strength int64  `json:"strength"`
}

type RecommendationsResponse struct {
	Flavor          string    `json:"flavor"`
	Recommendations []Pairing `json:"recommendations"`
}

func NewHandler(driver neo4j.DriverWithContext) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			flavor := r.URL.Query().Get("flavor")

			fmt.Printf("Received flavor: %s\n", flavor)

			// Create a session
			session := driver.NewSession(neo4j.SessionConfig{})
			defer session.Close()

			// Fetching the Cypher query from your map
			query, err := cypherQueries.GetRecommendationsQuery("GetRecommendationsQuery")
			if err != nil {
				fmt.Printf("Error fetching query: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		var recommendations []Pairing
		result, err := session.Run(query, map[string]interface{}{"flavor": flavor})

		if err != nil {
			fmt.Printf("Error running query: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for result.Next() {
			record := result.Record()
			name, _ := record.Get("recommendation")
			strengthMap, _ := record.Get("strength")

			if strengthData, ok := strengthMap.(map[string]interface{}); ok {
				strength := strengthData["low"].(int64)

				if nameStr, ok := name.(string); ok {
					recommendations = append(recommendations, Pairing{Name: nameStr, Strength: strength})
				}
			}
		}

		if err = result.Err(); err != nil {
			fmt.Printf("Error iterating through query results: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Printf("Final recommendations: %v\n", recommendations)

		response := RecommendationsResponse{
			Flavor:          flavor,
			Recommendations: recommendations,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}


func getRecommendations(flavor string, driver neo4j.DriverWithContext, query string) ([]Pairing, error) {
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var recommendations []Pairing

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{"flavor": flavor}

	result, err := tx.Run(ctx, query, params)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	for result.Next(ctx) {
    record := result.Record()
    name, _ := record.Get("recommendation")
    strengthMap, _ := record.Get("strength")

    if strengthData, ok := strengthMap.(map[string]interface{}); ok {
        low, _ := strengthData["low"].(int64)
        high, _ := strengthData["high"].(int64)
        strength := low + (high << 32)
        
        if nameStr, ok := name.(string); ok {
            recommendations = append(recommendations, Pairing{Name: nameStr, Strength: strength})
        }
    }
}


	if err = result.Err(); err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return recommendations, nil
}
