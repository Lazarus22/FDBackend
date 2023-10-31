package cypherQueries

import (
	"errors"
	"os"
	"strings"
	"fmt"
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
	fileContent, err := os.ReadFile("./cypherQueries/queries.cypher")
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
			comment := strings.TrimSpace(strings.TrimPrefix(lines[0], "//"))  // This line changed
			queryBody := strings.Join(lines[1:], "\n")

			QueryMap[comment] = queryBody
			fmt.Printf("Loaded query with key: %s\n", comment) // Logging for debugging
	}

	return nil
}

func GetRecommendationsQuery(key string) (string, error) {
	lowerKey := strings.ToLower(key)
	for k, v := range QueryMap {
			if strings.ToLower(k) == lowerKey {
					return v, nil
			}
	}
	return "", errors.New("query not found")
}


