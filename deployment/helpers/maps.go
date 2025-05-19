package helpers

// AddValueToNestedMap adds a value to a map nested within another map.
func AddValueToNestedMap[K1 comparable, K2 comparable, V any](mapping map[K1]map[K2]V, key1 K1, key2 K2, value V) map[K1]map[K2]V {
	if mapping == nil {
		mapping = make(map[K1]map[K2]V)
	}
	if mapping[key1] == nil {
		mapping[key1] = make(map[K2]V)
		mapping[key1][key2] = value
		return mapping
	}
	mapping[key1][key2] = value
	return mapping
}
