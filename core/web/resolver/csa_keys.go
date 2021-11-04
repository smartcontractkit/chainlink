package resolver

import (
	"errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
)

// CSAKeyResolver resolves the single CSA Key object
type CSAKeyResolver struct {
	key csakey.KeyV2
}

func NewCSAKey(key csakey.KeyV2) *CSAKeyResolver {
	return &CSAKeyResolver{key}
}

// Version resolves the CSA Key version number.
func (r *CSAKeyResolver) Version() int32 {
	return int32(r.key.Version)
}

// PubKey resolves the CSA Key public key string.
func (r *CSAKeyResolver) PubKey() string {
	return r.key.PublicKeyString()
}

// -- CSAKeys Query --

type CSAKeysPayloadResolver struct {
	keys []csakey.KeyV2
}

func NewCSAKeysResolver(keys []csakey.KeyV2) *CSAKeysPayloadResolver {
	return &CSAKeysPayloadResolver{keys}
}

func (r *CSAKeysPayloadResolver) Results() []*CSAKeyResolver {
	return NewCSAKeys(r.keys)
}

func NewCSAKeys(keys []csakey.KeyV2) []*CSAKeyResolver {
	var resolvers []*CSAKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewCSAKey(k))
	}

	return resolvers
}

// -- CreateCSAKey Mutation --

type CreateCSAKeyPayloadResolver struct {
	key *csakey.KeyV2
	err error
}

func NewCreateCSAKeyPayload(key *csakey.KeyV2, err error) *CreateCSAKeyPayloadResolver {
	return &CreateCSAKeyPayloadResolver{key, err}
}

func (r *CreateCSAKeyPayloadResolver) ToCreateCSAKeySuccess() (*CreateCSAKeySuccessResolver, bool) {
	if r.key != nil {
		return NewCreateCSAKeySuccessResolver(r.key), true
	}

	return nil, false
}

func (r *CreateCSAKeyPayloadResolver) ToCSAKeyExistsError() (*CSAKeyExistsErrorResolver, bool) {
	if r.err != nil && errors.Is(r.err, keystore.ErrCSAKeyExists) {
		return NewCSAKeyExistsError(r.err.Error()), true
	}

	return nil, false
}

type CreateCSAKeySuccessResolver struct {
	key *csakey.KeyV2
}

func NewCreateCSAKeySuccessResolver(key *csakey.KeyV2) *CreateCSAKeySuccessResolver {
	return &CreateCSAKeySuccessResolver{key}
}

func (r *CreateCSAKeySuccessResolver) CSAKey() *CSAKeyResolver {
	return NewCSAKey(*r.key)
}

type CSAKeyExistsErrorResolver struct {
	message string
}

func NewCSAKeyExistsError(message string) *CSAKeyExistsErrorResolver {
	return &CSAKeyExistsErrorResolver{
		message: message,
	}
}

func (r *CSAKeyExistsErrorResolver) Message() string {
	return r.message
}

func (r *CSAKeyExistsErrorResolver) Code() ErrorCode {
	return ErrorCodeUnprocessable
}
