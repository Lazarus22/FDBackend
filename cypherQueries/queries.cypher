// GetAllFlavorsQuery1
MATCH (f:Flavor {name: $flavor})-[r:PAIRS_WITH]->(recommendation:Flavor)
RETURN recommendation.name AS recommendation, r.strength AS strength