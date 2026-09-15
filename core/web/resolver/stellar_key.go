package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
)

type StellarKeyResolver struct {
	key stellarkey.Key
}

func NewStellarKey(key stellarkey.Key) *StellarKeyResolver {
	return &StellarKeyResolver{key: key}
}

func NewStellarKeys(keys []stellarkey.Key) []*StellarKeyResolver {
	resolvers := make([]*StellarKeyResolver, 0, len(keys))

	for _, k := range keys {
		resolvers = append(resolvers, NewStellarKey(k))
	}

	return resolvers
}

func (r *StellarKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.ID())
}

// Account is the Stellar account address in StrKey form (G…), which is what a
// feeds-manager chain config stores as its accountAddr.
func (r *StellarKeyResolver) Account() string {
	return r.key.Account()
}

// -- GetStellarKeys Query --

type StellarKeysPayloadResolver struct {
	keys []stellarkey.Key
}

func NewStellarKeysPayload(keys []stellarkey.Key) *StellarKeysPayloadResolver {
	return &StellarKeysPayloadResolver{keys: keys}
}

func (r *StellarKeysPayloadResolver) Results() []*StellarKeyResolver {
	return NewStellarKeys(r.keys)
}
