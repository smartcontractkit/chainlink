package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
)

type SuiKeyResolver struct {
	key suikey.Key
}

func NewSuiKey(key suikey.Key) *SuiKeyResolver {
	return &SuiKeyResolver{key: key}
}

func NewSuiKeys(keys []suikey.Key) []*SuiKeyResolver {
	resolvers := []*SuiKeyResolver{}

	for _, k := range keys {
		resolvers = append(resolvers, NewSuiKey(k))
	}

	return resolvers
}

func (r *SuiKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.PublicKeyStr())
}

func (r *SuiKeyResolver) Account() string {
	return r.key.Account()
}

// -- GetSuiKeys Query --

type SuiKeysPayloadResolver struct {
	keys []suikey.Key
}

func NewSuiKeysPayload(keys []suikey.Key) *SuiKeysPayloadResolver {
	return &SuiKeysPayloadResolver{keys: keys}
}

func (r *SuiKeysPayloadResolver) Results() []*SuiKeyResolver {
	return NewSuiKeys(r.keys)
}
