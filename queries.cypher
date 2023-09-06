// GetRecommendationsQuery
MATCH (i1)-[r:pairs_with]->(i2)
WHERE i1.name = $flavor OR properties(r).Value = $flavor
RETURN i2.name AS recommendation, r.Value AS strength
