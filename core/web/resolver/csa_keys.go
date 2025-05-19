package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
)

// CSAKeyResolver resolves the single CSA Key object
type CSAKeyResolver struct {
	key csakey.KeyV2
}

func NewCSAKey(key csakey.KeyV2) *CSAKeyResolver {
	return &CSAKeyResolver{key: key}
}

// ID resolves the CSA Key public key as the id.
func (r *CSAKeyResolver) ID() graphql.ID {
	return graphql.ID(r.key.ID())
}

// PubKey resolves the CSA Key public key string.
func (r *CSAKeyResolver) PublicKey() string {
	return "csa_" + r.key.PublicKeyString()
}

// Version resolves the CSA Key version number.
func (r *CSAKeyResolver) Version() int32 {
	return int32(r.key.Version)
}

// -- CSAKeys Query --

type CSAKeysPayloadResolver struct {
	keys []csakey.KeyV2
}

func NewCSAKeysResolver(keys []csakey.KeyV2) *CSAKeysPayloadResolver {
	return &CSAKeysPayloadResolver{keys: keys}
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
	return &CreateCSAKeyPayloadResolver{key: key, err: err}
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
	return &CreateCSAKeySuccessResolver{key: key}
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

type DeleteCSAKeySuccessResolver struct {
	key csakey.KeyV2
}

func NewDeleteCSAKeySuccess(key csakey.KeyV2) *DeleteCSAKeySuccessResolver {
	return &DeleteCSAKeySuccessResolver{key: key}
}

func (r *DeleteCSAKeySuccessResolver) CSAKey() *CSAKeyResolver {
	return NewCSAKey(r.key)
}

type DeleteCSAKeyPayloadResolver struct {
	key csakey.KeyV2
	NotFoundErrorUnionType
}

func NewDeleteCSAKeyPayload(key csakey.KeyV2, err error) *DeleteCSAKeyPayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			return errors.As(err, &keystore.KeyNotFoundError{})
		}}
	}

	return &DeleteCSAKeyPayloadResolver{key: key, NotFoundErrorUnionType: e}
}

func (r *DeleteCSAKeyPayloadResolver) ToDeleteCSAKeySuccess() (*DeleteCSAKeySuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteCSAKeySuccess(r.key), true
	}
	return nil, false
}
