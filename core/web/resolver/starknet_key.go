package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

type StarkNetKeyResolver struct {
	key starkkey.Key
}

func NewStarkNetKey(key starkkey.Key) *StarkNetKeyResolver {
	return &StarkNetKeyResolver{key: key}
}

func NewStarkNetKeys(keys []starkkey.Key) []*StarkNetKeyResolver {
	var resolvers []*StarkNetKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewStarkNetKey(k))
	}

	return resolvers
}

func (r *StarkNetKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.StarkKeyStr())
}

// -- GetStarkNetKeys Query --

type StarkNetKeysPayloadResolver struct {
	keys []starkkey.Key
}

func NewStarkNetKeysPayload(keys []starkkey.Key) *StarkNetKeysPayloadResolver {
	return &StarkNetKeysPayloadResolver{keys: keys}
}

func (r *StarkNetKeysPayloadResolver) Results() []*StarkNetKeyResolver {
	return NewStarkNetKeys(r.keys)
}
