package recommendations

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

type Pairing struct {
	Name     string `json:"name"`
	Strength int    `json:"strength"`
}

type RecommendationsResponse struct {
	Flavor          string   `json:"flavor"`
	Recommendations []Pairing `json:"recommendations"`
}

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

func getRecommendations(flavor string, driver neo4j.DriverWithContext) ([]Pairing, error) {
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var recommendations []Pairing

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	query := `
	MATCH (i1)-[r:pairs_with]->(i2)
	WHERE i1.name = $flavor OR properties(r).Value = $flavor
	RETURN i2.name AS recommendation, r.Value AS strength
	`
	params := map[string]interface{}{"flavor": flavor}

	result, err := tx.Run(ctx, query, params)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	for result.Next(ctx) {
		record := result.Record()
		name, ok1 := record.Get("recommendation")
		strength, ok2 := record.Get("strength")
		if ok1 && ok2 {
				recommendations = append(recommendations, Pairing{Name: name.(string), Strength: strength.(int)})
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
