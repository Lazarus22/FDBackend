package recommendations

import (
	"FDBackend/cypherQueries"
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"strings"
	"time"
)

type Pairing struct {
	Name             string `json:"name"`
	Strength         int    `json:"strength"`
	RelationshipType string `json:"relationshipType"`
	NodeType         string `json:"nodeType"`
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

		// Initialize to a non-nil zero-length slice if recommendations is nil
		if recommendations == nil {
			recommendations = make([]Pairing, 0)
		}

		lowercaseFlavor := strings.ToLower(flavor)

		response := RecommendationsResponse{
			Flavor:          lowercaseFlavor,
			Recommendations: recommendations,
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(http.TimeFormat)) // Set Expires header
		json.NewEncoder(w).Encode(response)
	}
}

func AutoCompleteHandler(driver neo4j.DriverWithContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := r.URL.Query().Get("prefix")
		if prefix == "" {
			http.Error(w, "Prefix is required", http.StatusBadRequest)
			return
		}

		query, err := cypherQueries.GetRecommendationsQuery("GetAutocompleteSuggestions")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		suggestions, err := getSuggestions(prefix, driver, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Initialize to a non-nil zero-length slice if suggestions is nil
		if suggestions == nil {
			suggestions = make([]string, 0)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(suggestions)
	}
}

func getSuggestions(prefix string, driver neo4j.DriverWithContext, query string) ([]string, error) {
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	var suggestions []string

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{"prefix": prefix}

	result, err := tx.Run(ctx, query, params)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	for result.Next(ctx) {
		record := result.Record()
		suggestion, _ := record.Get("suggestion")
		if suggestionStr, ok := suggestion.(string); ok {
			suggestions = append(suggestions, suggestionStr)
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

	if suggestions == nil {
		suggestions = make([]string, 0)
	}

	return suggestions, nil
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
		strength, _ := record.Get("value")
		relationshipType, _ := record.Get("relationshipType")
		nodeType, _ := record.Get("nodeType")

		if nameStr, ok := name.(string); ok {
			if strengthVal, ok := strength.(int64); ok {
				if relationshipTypeStr, ok := relationshipType.(string); ok {
					if nodeTypeStr, ok := nodeType.(string); ok {
						recommendations = append(recommendations, Pairing{
							Name:             nameStr,
							Strength:         int(strengthVal),
							RelationshipType: relationshipTypeStr,
							NodeType:         nodeTypeStr,
						})
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

	if recommendations == nil {
		recommendations = make([]Pairing, 0)
	}

	return recommendations, nil
}
