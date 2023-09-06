// GetRecommendationsQuery
MATCH (i1)-[r:pairs_with]->(i2)
WHERE i1.name = 'chicken' OR r.value = 'chicken'
RETURN i2.name AS recommendation, r.value AS strength