package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
)

type TONKeyResolver struct {
	key tonkey.Key
}

func NewTONKey(key tonkey.Key) *TONKeyResolver {
	return &TONKeyResolver{key: key}
}

func NewTONKeys(keys []tonkey.Key) []*TONKeyResolver {
	resolvers := make([]*TONKeyResolver, 0, len(keys))

	for _, k := range keys {
		resolvers = append(resolvers, NewTONKey(k))
	}

	return resolvers
}

func (r *TONKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.ID())
}

func (r *TONKeyResolver) AddressBase64() string {
	return r.key.AddressBase64()
}

func (r *TONKeyResolver) RawAddress() string {
	return r.key.RawAddress()
}

// -- GetTONKeys Query --

type TONKeysPayloadResolver struct {
	keys []tonkey.Key
}

func NewTONKeysPayload(keys []tonkey.Key) *TONKeysPayloadResolver {
	return &TONKeysPayloadResolver{keys: keys}
}

func (r *TONKeysPayloadResolver) Results() []*TONKeyResolver {
	return NewTONKeys(r.keys)
}
