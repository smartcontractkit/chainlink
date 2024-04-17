package telemetry

import (
	"fmt"
)

var keyString = fmt.Sprintf("%064b", 0)

// getDummyKeyString returns a dummy key string
// satisfies the wsrpc key length constraints
func GetDummyKeyString() string {
	return keyString
}

// getDummyKeyString returns a dummy key string with the specified prefix
// satisfies the wsrpc key length constraints
func GetDummyKeyStringWithPrefix(prefix string) string {
	combo := prefix + GetDummyKeyString()
	return combo[:64]
}
