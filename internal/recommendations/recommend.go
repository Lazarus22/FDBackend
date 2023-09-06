package recommendations

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

// RecommendationsResponse represents the structure of the response
type RecommendationsResponse struct {
	Flavor          string   `json:"flavor"`
	Recommendations []string `json:"recommendations"`
}

// NewHandler returns a new HTTP handler function for recommendations.
func NewHandler(driver neo4j.DriverWithContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flavor := r.URL.Query().Get("flavor")
		if flavor == "" {
			http.Error(w, "Flavor is required", http.StatusBadRequest)
			return
		}

		recommendations, err := getRecommendations(flavor, driver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := RecommendationsResponse{
			Flavor:          flavor,
			Recommendations: recommendations,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getRecommendations(flavor string, driver neo4j.DriverWithContext) ([]string, error) {
	ctx := context.Background()

	// Add context as the first argument
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var recommendations []string

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	query := `
	MATCH (i1:Ingredient)-[r:pairs_with]->(i2:Ingredient)
	WHERE i1.name = $flavor OR r.Property = $flavor
	RETURN i2.name AS recommendation, r.Value AS strength
	`
	params := map[string]interface{}{"flavor": flavor}

	result, err := tx.Run(ctx, query, params)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	// Add context as an argument
	for result.Next(ctx) {
		record := result.Record()
		value, ok := record.Get("recommendation")
		if ok {
			recommendations = append(recommendations, value.(string))
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
