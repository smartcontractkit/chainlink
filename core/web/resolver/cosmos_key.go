package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

type CosmosKeyResolver struct {
	key cosmoskey.Key
}

func NewCosmosKey(key cosmoskey.Key) *CosmosKeyResolver {
	return &CosmosKeyResolver{key: key}
}

func NewCosmosKeys(keys []cosmoskey.Key) []*CosmosKeyResolver {
	var resolvers []*CosmosKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewCosmosKey(k))
	}

	return resolvers
}

func (r *CosmosKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.PublicKeyStr())
}

// -- GetCosmosKeys Query --

type CosmosKeysPayloadResolver struct {
	keys []cosmoskey.Key
}

func NewCosmosKeysPayload(keys []cosmoskey.Key) *CosmosKeysPayloadResolver {
	return &CosmosKeysPayloadResolver{keys: keys}
}

func (r *CosmosKeysPayloadResolver) Results() []*CosmosKeyResolver {
	return NewCosmosKeys(r.keys)
}
