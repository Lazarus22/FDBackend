// GetRecommendationsQuery
MATCH (n)-[r:pairs_with]-(m)
WHERE n.name = $flavor AND n.name <> m.name
RETURN m.name as recommendation, r.value as value, labels(m) as labels
