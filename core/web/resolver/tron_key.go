package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
)

type TronKeyResolver struct {
	key tronkey.Key
}

func NewTronKey(key tronkey.Key) *TronKeyResolver {
	return &TronKeyResolver{key: key}
}

func NewTronKeys(keys []tronkey.Key) []*TronKeyResolver {
	resolvers := make([]*TronKeyResolver, 0, len(keys))

	for _, k := range keys {
		resolvers = append(resolvers, NewTronKey(k))
	}

	return resolvers
}

func (r *TronKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.ID())
}

// -- GetTronKeys Query --

type TronKeysPayloadResolver struct {
	keys []tronkey.Key
}

func NewTronKeysPayload(keys []tronkey.Key) *TronKeysPayloadResolver {
	return &TronKeysPayloadResolver{keys: keys}
}

func (r *TronKeysPayloadResolver) Results() []*TronKeyResolver {
	return NewTronKeys(r.keys)
}
