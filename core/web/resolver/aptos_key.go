package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
)

type AptosKeyResolver struct {
	key aptoskey.Key
}

func NewAptosKey(key aptoskey.Key) *AptosKeyResolver {
	return &AptosKeyResolver{key: key}
}

func NewAptosKeys(keys []aptoskey.Key) []*AptosKeyResolver {
	var resolvers []*AptosKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewAptosKey(k))
	}

	return resolvers
}

func (r *AptosKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.PublicKeyStr())
}

func (r *AptosKeyResolver) Account() string {
	return r.key.Account()
}

// -- GetAptosKeys Query --

type AptosKeysPayloadResolver struct {
	keys []aptoskey.Key
}

func NewAptosKeysPayload(keys []aptoskey.Key) *AptosKeysPayloadResolver {
	return &AptosKeysPayloadResolver{keys: keys}
}

func (r *AptosKeysPayloadResolver) Results() []*AptosKeyResolver {
	return NewAptosKeys(r.keys)
}
