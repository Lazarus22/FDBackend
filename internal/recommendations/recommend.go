package recommendations

import (
	"FDBackend/cypherQueries"
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
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
		ctx := context.Background()

		flavor := r.URL.Query().Get("flavor")

		session := driver.NewSession(ctx, neo4j.SessionConfig{})
		defer session.Close(ctx)

		query, err := cypherQueries.GetRecommendationsQuery("GetAllFlavorsQuery")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		params := map[string]interface{}{"flavor": flavor}
		result, err := session.Run(ctx, query, params)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var recommendations []Pairing
		for result.Next(ctx) {
			record := result.Record()
			name, ok := record.Get("flavorName")
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if nameStr, ok := name.(string); ok {
				recommendations = append(recommendations, Pairing{Name: nameStr, Strength: 0}) // Here Strength is set to 0 as a placeholder.
			}
		}

		if err = result.Err(); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resp := RecommendationsResponse{
			Flavor:          flavor,
			Recommendations: recommendations,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

