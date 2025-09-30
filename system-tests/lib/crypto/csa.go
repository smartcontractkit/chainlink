package crypto

import "strings"

type CSAKey struct {
	Key string
}

// CleansedKey returns the key without the "csa_" prefix
func (c *CSAKey) CleansedKey() string {
	return strings.TrimPrefix(c.Key, "csa_")
}
