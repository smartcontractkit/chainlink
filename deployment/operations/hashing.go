package operations

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"sync"
)

func constructUniqueHashFrom(hashCache *sync.Map, def Definition, input any) (string, error) {
	// convert def and input into string to be used as cachekey
	// as input can be any type, some types cannot be used as cache argument for sync.map
	// converting to string ensure argument is always valid

	key, err := json.Marshal(def)
	if err != nil {
		return "", err
	}

	// convert to a canonical representation that's stable regardless of field order
	// If the keys are not sorted, the hash generated can be different in maps.
	inputBytes, err := canonicalizeJSON(input)
	if err != nil {
		return "", err
	}

	key = append(key, inputBytes...)
	if cached, ok := hashCache.Load(string(key)); ok {
		return cached.(string), nil
	}

	hash := sha256.Sum256(key)
	result := hex.EncodeToString(hash[:])

	hashCache.Store(string(key), result)
	return result, nil
}

// canonicalizeJSON converts value to a canonical JSON representation with sorted keys
func canonicalizeJSON(value any) ([]byte, error) {
	// First marshal to standard JSON
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a structure that can be canonically marshaled
	var data any
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	// Marshal with sorted keys
	return json.Marshal(canonicalize(data))
}

// canonicalize recursively processes data structures for stable serialization
func canonicalize(data any) any {
	switch v := data.(type) {
	case map[string]any:
		// For maps, create a sorted representation
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		m := make(map[string]any, len(v))
		for _, k := range keys {
			m[k] = canonicalize(v[k]) // Recursively process values
		}
		return m

	case []any:
		// For arrays, recursively process each element
		for i, val := range v {
			v[i] = canonicalize(val)
		}
	}
	return data
}
