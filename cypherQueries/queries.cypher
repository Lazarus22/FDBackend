// GetRecommendationsQuery
MATCH (i1)-[r:pairs_with]-(i2) 
WHERE (i1.name=$flavor OR i2.name=$flavor) AND i1 <> i2
RETURN i1.name, labels(i1), i2.name, labels(i2), r.value as strength


