// GetRecommendationsQuery
MATCH (i1)-[r:pairs_with]->(i2)
WHERE i1.name = $flavor OR r.property = $flavor
RETURN i2.name AS recommendation, r.value AS strength
// GetAllFlavorsQuery
MATCH (n:Flavor)
RETURN n.name AS flavorName