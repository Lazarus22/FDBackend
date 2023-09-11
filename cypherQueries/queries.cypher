// GetRecommendationsQuery
MATCH (n)-[r]-(m)
WHERE n.name = $flavor AND n.name <> m.name AND type(r) IN [
  'pairs_with', 
  'in_season', 
  'has_function', 
  'related_to', 
  'key_ingredient', 
  'has_taste', 
  'has_volume', 
  'has_weight', 
  'uses_technique'
]
RETURN m.name as recommendation, r.value as value, labels(m) as labels, type(r) as relationshipType, head(labels(m)) as nodeType
