package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
)

type P2PKeyResolver struct {
	key p2pkey.KeyV2
}

func NewP2PKeyResolver(key p2pkey.KeyV2) P2PKeyResolver {
	return P2PKeyResolver{key}
}

func (k P2PKeyResolver) ID() graphql.ID {
	return graphql.ID(k.key.ID())
}

func (k P2PKeyResolver) PeerID() string {
	return k.key.PeerID().String()
}

func (k P2PKeyResolver) PublicKey() string {
	return k.key.PublicKeyHex()
}

type P2PKeysPayloadResolver struct {
	keys []p2pkey.KeyV2
}

func NewP2PKeysPayloadResolver(keys []p2pkey.KeyV2) *P2PKeysPayloadResolver {
	return &P2PKeysPayloadResolver{keys}
}

func (r *P2PKeysPayloadResolver) Results() []P2PKeyResolver {
	results := []P2PKeyResolver{}
	for _, k := range r.keys {
		results = append(results, NewP2PKeyResolver(k))
	}
	return results
}

type CreateP2PKeyPayloadResolver struct {
	key p2pkey.KeyV2
}

func NewCreateP2PKeyPayloadResolver(key p2pkey.KeyV2) *CreateP2PKeyPayloadResolver {
	return &CreateP2PKeyPayloadResolver{key}
}

func (r *CreateP2PKeyPayloadResolver) Key() P2PKeyResolver {
	return NewP2PKeyResolver(r.key)
}

type DeleteP2PKeySuccessResolver struct {
	key p2pkey.KeyV2
}

func NewDeleteP2PKeySuccessResolver(key p2pkey.KeyV2) *DeleteP2PKeySuccessResolver {
	return &DeleteP2PKeySuccessResolver{key}
}

func (r *DeleteP2PKeySuccessResolver) Key() P2PKeyResolver {
	return NewP2PKeyResolver(r.key)
}

type DeleteP2PKeyPayloadResolver struct {
	key p2pkey.KeyV2
	err error
}

func NewDeleteP2PKeyPayloadResolver(key p2pkey.KeyV2, err error) *DeleteP2PKeyPayloadResolver {
	return &DeleteP2PKeyPayloadResolver{key, err}
}

func (r *DeleteP2PKeyPayloadResolver) ToDeleteP2PKeySuccess() (*DeleteP2PKeySuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteP2PKeySuccessResolver(r.key), true
	}
	return nil, false
}

func (r *DeleteP2PKeyPayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil {
		return NewNotFoundError(r.err.Error()), true
	}
	return nil, false
}
