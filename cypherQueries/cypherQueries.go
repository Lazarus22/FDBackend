package cypherQueries

import (
	"errors"
	"os"
	"strings"
)

var QueryMap map[string]string

func init() {
	err := InitializeQueries()
	if err != nil {
			panic(err)
	}
}

func InitializeQueries() error {
	QueryMap = make(map[string]string)

	// Read the file
	fileContent, err := os.ReadFile("./queries.cypher") 
	if err != nil {
		return err
	}

	// Convert to string and split by semicolon
	queries := strings.Split(string(fileContent), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		lines := strings.Split(query, "\n")
		comment := strings.TrimSpace(lines[0])
		queryBody := strings.Join(lines[1:], "\n")

		QueryMap[comment] = queryBody
	}

	return nil
}

func GetRecommendationsQuery(key string) (string, error) {
	if query, exists := QueryMap[key]; exists {
		return query, nil
	}
	return "", errors.New("query not found")
}
