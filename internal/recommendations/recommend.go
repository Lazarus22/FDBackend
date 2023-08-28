package recommendations

import (
	"encoding/json"
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

// RecommendationsResponse represents the structure of the response
type RecommendationsResponse struct {
	Flavor         string   `json:"flavor"`
	Recommendations []string `json:"recommendations"`
}

// Handler function to handle HTTP requests for recommendations
func Handler(w http.ResponseWriter, r *http.Request) {
	flavor := r.URL.Query().Get("flavor")
	if flavor == "" {
		http.Error(w, "Flavor is required", http.StatusBadRequest)
		return
	}

	recommendations, err := getRecommendations(flavor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := RecommendationsResponse{
			Flavor:         flavor,
			Recommendations: recommendations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getRecommendations(flavor string) ([]string, error) {
	ctx := context.Background()
	uri := "neo4j+s://e29d74fa.databases.neo4j.io"
	username := "neo4j" // Your Neo4j username
	password := "Allagash22" // Your Neo4j password

	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// Updated Cypher query using correct labels and relationship type
	query := "MATCH (i1:Ingredient {name: $flavor})-[:pairs_with]->(i2:Ingredient) RETURN i2.name AS recommendation"

	// Run query using a managed transaction
	people, err := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (interface{}, error) {
			result, err := tx.Run(ctx, query, map[string]interface{}{
				"flavor": flavor,
			})
			if err != nil {
				return nil, err
			}
			records, err := result.Collect(ctx)
			if err != nil {
				return nil, err
			}
			return records, nil
		})
	if err != nil {
		return nil, err
	}

	// Collect recommendations
	var recommendations []string
	for _, record := range people.([]*neo4j.Record) {
		recommendations = append(recommendations, record.AsMap()["recommendation"].(string))
	}

	return recommendations, nil
}
