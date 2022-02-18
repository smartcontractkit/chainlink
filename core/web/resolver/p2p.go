package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
)

type P2PKeyResolver struct {
	key p2pkey.KeyV2
}

func NewP2PKey(key p2pkey.KeyV2) P2PKeyResolver {
	return P2PKeyResolver{key: key}
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

// -- P2PKeys Query --

type P2PKeysPayloadResolver struct {
	keys []p2pkey.KeyV2
}

func NewP2PKeysPayload(keys []p2pkey.KeyV2) *P2PKeysPayloadResolver {
	return &P2PKeysPayloadResolver{keys: keys}
}

func (r *P2PKeysPayloadResolver) Results() []P2PKeyResolver {
	var results []P2PKeyResolver
	for _, k := range r.keys {
		results = append(results, NewP2PKey(k))
	}
	return results
}

// -- CreateP2PKey Mutation --

type CreateP2PKeySuccessResolver struct {
	key p2pkey.KeyV2
}

func NewCreateP2PKeySuccess(key p2pkey.KeyV2) *CreateP2PKeySuccessResolver {
	return &CreateP2PKeySuccessResolver{key: key}
}

func (r *CreateP2PKeySuccessResolver) P2PKey() P2PKeyResolver {
	return NewP2PKey(r.key)
}

type CreateP2PKeyPayloadResolver struct {
	p2pKey p2pkey.KeyV2
}

func NewCreateP2PKeyPayload(key p2pkey.KeyV2) *CreateP2PKeyPayloadResolver {
	return &CreateP2PKeyPayloadResolver{p2pKey: key}
}

func (r *CreateP2PKeyPayloadResolver) P2PKey() P2PKeyResolver {
	return NewP2PKey(r.p2pKey)
}

func (r *CreateP2PKeyPayloadResolver) ToCreateP2PKeySuccess() (*CreateP2PKeySuccessResolver, bool) {
	return NewCreateP2PKeySuccess(r.p2pKey), true
}

// -- DeleteP2PKey Mutation --

type DeleteP2PKeySuccessResolver struct {
	p2pKey p2pkey.KeyV2
}

func NewDeleteP2PKeySuccess(p2pKey p2pkey.KeyV2) *DeleteP2PKeySuccessResolver {
	return &DeleteP2PKeySuccessResolver{p2pKey: p2pKey}
}

func (r *DeleteP2PKeySuccessResolver) P2PKey() P2PKeyResolver {
	return NewP2PKey(r.p2pKey)
}

type DeleteP2PKeyPayloadResolver struct {
	p2pKey p2pkey.KeyV2
	NotFoundErrorUnionType
}

func NewDeleteP2PKeyPayload(p2pKey p2pkey.KeyV2, err error) *DeleteP2PKeyPayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			return errors.As(err, &keystore.KeyNotFoundError{})
		}}
	}

	return &DeleteP2PKeyPayloadResolver{p2pKey: p2pKey, NotFoundErrorUnionType: e}
}

func (r *DeleteP2PKeyPayloadResolver) ToDeleteP2PKeySuccess() (*DeleteP2PKeySuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteP2PKeySuccess(r.p2pKey), true
	}
	return nil, false
}
