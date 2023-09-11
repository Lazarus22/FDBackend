// GetRecommendationsQuery
MATCH (i1)-[r:pairs_with]->(i2) 
WHERE i1.name = $flavor OR i2.name = $flavor 
RETURN i2.name as recommendation, r.value as value

