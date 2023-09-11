package recommendations

import (
	"FDBackend/cypherQueries"
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

type Pairing struct {
	Name     string   `json:"name"`
	Strength int      `json:"strength"`
	Labels   []string `json:"labels"`
}

type RecommendationsResponse struct {
	Flavor          string    `json:"flavor"`
	Recommendations []Pairing `json:"recommendations"`
}

func NewHandler(driver neo4j.DriverWithContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flavor := r.URL.Query().Get("flavor")
		if flavor == "" {
			http.Error(w, "Flavor is required", http.StatusBadRequest)
			return
		}

		query, err := cypherQueries.GetRecommendationsQuery("GetRecommendationsQuery")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		recommendations, err := getRecommendations(flavor, driver, query)
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

			name, ok1 := record.Get("recommendation")
			strength, ok2 := record.Get("value")
			labels, ok3 := record.Get("labels")

			if ok1 && ok2 && ok3 && name != nil && name != flavor {
					if nameStr, ok := name.(string); ok {
							if strengthVal, ok := strength.(int64); ok {
									if labelsVal, ok := labels.([]interface{}); ok {
											strLabels := make([]string, len(labelsVal))
											for i, label := range labelsVal {
													strLabels[i] = label.(string)
											}
											recommendations = append(recommendations, Pairing{Name: nameStr, Strength: int(strengthVal), Labels: strLabels})
									}
							}
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
