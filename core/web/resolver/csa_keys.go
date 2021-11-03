package resolver

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
)

// CSAKeyResolver resolves the single CSA Key object
type CSAKeyResolver struct {
	key csakey.KeyV2
}

func NewCSAKey(key csakey.KeyV2) *CSAKeyResolver {
	return &CSAKeyResolver{key}
}

// Version resolves the CSA Key version number.
func (r *CSAKeyResolver) Version() int32 {
	return int32(r.key.Version)
}

// PubKey resolves the CSA Key public key string.
func (r *CSAKeyResolver) PubKey() string {
	return r.key.PublicKeyString()
}

// -- CSAKeys Query --

type CSAKeysPayloadResolver struct {
	keys []csakey.KeyV2
}

func NewCSAKeysResolver(keys []csakey.KeyV2) *CSAKeysPayloadResolver {
	return &CSAKeysPayloadResolver{keys}
}

func (r *CSAKeysPayloadResolver) Results() []*CSAKeyResolver {
	return NewCSAKeys(r.keys)
}

func NewCSAKeys(keys []csakey.KeyV2) []*CSAKeyResolver {
	var resolvers []*CSAKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewCSAKey(k))
	}

	return resolvers
}
