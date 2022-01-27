package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
)

type SolanaKeyResolver struct {
	key solkey.Key
}

func NewSolanaKey(key solkey.Key) *SolanaKeyResolver {
	return &SolanaKeyResolver{key: key}
}

func NewSolanaKeys(keys []solkey.Key) []*SolanaKeyResolver {
	var resolvers []*SolanaKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewSolanaKey(k))
	}

	return resolvers
}

func (r *SolanaKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.PublicKeyStr())
}

// -- GetSolanaKeys Query --

type SolanaKeysPayloadResolver struct {
	keys []solkey.Key
}

func NewSolanaKeysPayload(keys []solkey.Key) *SolanaKeysPayloadResolver {
	return &SolanaKeysPayloadResolver{keys: keys}
}

func (r *SolanaKeysPayloadResolver) Results() []*SolanaKeyResolver {
	return NewSolanaKeys(r.keys)
}
