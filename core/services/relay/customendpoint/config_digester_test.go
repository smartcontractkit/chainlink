package customendpoint

import (
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
)

func TestConfigDigester(t *testing.T) {
	initialDigester := OffchainConfigDigester{
		EndpointName:   "dydx",
		EndpointTarget: "staging",
		PayloadType:    "ETHUSD",
	}
	assert.Equal(t, initialDigester.ConfigDigestPrefix(), ConfigDigestPrefixCustomEndpoint)
	initialDigest, _ := initialDigester.ConfigDigest(types.ContractConfig{})

	// Changing EndpointName changes the digest, but not the prefix.
	modifiedDigester := initialDigester
	modifiedDigester.EndpointName = "modified"
	modifiedDigest, _ := modifiedDigester.ConfigDigest(types.ContractConfig{})
	assert.Equal(t, modifiedDigester.ConfigDigestPrefix(), ConfigDigestPrefixCustomEndpoint)
	assert.NotEqual(t, initialDigest, modifiedDigest)

	// Changing EndpointTarget changes the digest, but not the prefix.
	modifiedDigester = initialDigester
	modifiedDigester.EndpointTarget = "modified"
	modifiedDigest, _ = modifiedDigester.ConfigDigest(types.ContractConfig{})
	assert.Equal(t, modifiedDigester.ConfigDigestPrefix(), ConfigDigestPrefixCustomEndpoint)
	assert.NotEqual(t, initialDigest, modifiedDigest)

	// Changing PayloadType changes the digest, but not the prefix.
	modifiedDigester = initialDigester
	modifiedDigester.PayloadType = "modified"
	modifiedDigest, _ = modifiedDigester.ConfigDigest(types.ContractConfig{})
	assert.Equal(t, modifiedDigester.ConfigDigestPrefix(), ConfigDigestPrefixCustomEndpoint)
	assert.NotEqual(t, initialDigest, modifiedDigest)
}
