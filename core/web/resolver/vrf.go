package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
)

type VRFKeyResolver struct {
	key vrfkey.KeyV2
}

func NewVRFKeyResolver(key vrfkey.KeyV2) VRFKeyResolver {
	return VRFKeyResolver{key: key}
}

// ID returns the ID of the VRF key, which is the public key.
func (k VRFKeyResolver) ID() graphql.ID {
	return graphql.ID(k.key.ID())
}

// Compressed returns the compressed version of the public key.
func (k VRFKeyResolver) Compressed() string {
	return k.key.PublicKey.String()
}

func (k VRFKeyResolver) Uncompressed() string {
	// It's highly unlikely that this call will return an error.
	// If it does, we'd likely have issues all throughout the application.
	// However, it's still good practice to handle the error that is returned
	// rather than completely ignoring it.
	uncompressed, err := k.key.PublicKey.StringUncompressed()
	if err != nil {
		uncompressed = "error: unable to uncompress public key"
	}
	return uncompressed
}

// Hash returns the hash of the VRF public key.
func (k VRFKeyResolver) Hash() string {
	var hashStr string

	// It's highly unlikely that this call will return an error.
	// If it does, we'd likely have issues all throughout the application.
	// However, it's still good practice to handle the error that is returned
	// rather than completely ignoring it.
	hash, err := k.key.PublicKey.Hash()
	if err != nil {
		hashStr = "error: unable to get public key hash"
	} else {
		hashStr = hash.String()
	}
	return hashStr
}

type VRFKeySuccessResolver struct {
	key vrfkey.KeyV2
}

func NewVRFKeySuccessResolver(key vrfkey.KeyV2) *VRFKeySuccessResolver {
	return &VRFKeySuccessResolver{key: key}
}

func (r *VRFKeySuccessResolver) Key() VRFKeyResolver {
	return NewVRFKeyResolver(r.key)
}

type VRFKeyPayloadResolver struct {
	key vrfkey.KeyV2
	NotFoundErrorUnionType
}

func NewVRFKeyPayloadResolver(key vrfkey.KeyV2, err error) *VRFKeyPayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			return errors.Is(errors.Cause(err), keystore.ErrMissingVRFKey)
		}}
	}

	return &VRFKeyPayloadResolver{key: key, NotFoundErrorUnionType: e}
}

func (r *VRFKeyPayloadResolver) ToVRFKeySuccess() (*VRFKeySuccessResolver, bool) {
	if r.err == nil {
		return NewVRFKeySuccessResolver(r.key), true
	}
	return nil, false
}

type VRFKeysPayloadResolver struct {
	keys []vrfkey.KeyV2
}

func NewVRFKeysPayloadResolver(keys []vrfkey.KeyV2) *VRFKeysPayloadResolver {
	return &VRFKeysPayloadResolver{keys: keys}
}

func (r *VRFKeysPayloadResolver) Results() []VRFKeyResolver {
	var results []VRFKeyResolver
	for _, k := range r.keys {
		results = append(results, NewVRFKeyResolver(k))
	}
	return results
}

type CreateVRFKeyPayloadResolver struct {
	key vrfkey.KeyV2
}

func NewCreateVRFKeyPayloadResolver(key vrfkey.KeyV2) *CreateVRFKeyPayloadResolver {
	return &CreateVRFKeyPayloadResolver{key: key}
}

func (r *CreateVRFKeyPayloadResolver) Key() VRFKeyResolver {
	return NewVRFKeyResolver(r.key)
}

type DeleteVRFKeySuccessResolver struct {
	key vrfkey.KeyV2
}

func NewDeleteVRFKeySuccessResolver(key vrfkey.KeyV2) *DeleteVRFKeySuccessResolver {
	return &DeleteVRFKeySuccessResolver{key: key}
}

func (r *DeleteVRFKeySuccessResolver) Key() VRFKeyResolver {
	return NewVRFKeyResolver(r.key)
}

type DeleteVRFKeyPayloadResolver struct {
	key vrfkey.KeyV2
	NotFoundErrorUnionType
}

func NewDeleteVRFKeyPayloadResolver(key vrfkey.KeyV2, err error) *DeleteVRFKeyPayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			return errors.Is(errors.Cause(err), keystore.ErrMissingVRFKey)
		}}
	}

	return &DeleteVRFKeyPayloadResolver{key: key, NotFoundErrorUnionType: e}
}

func (r *DeleteVRFKeyPayloadResolver) ToDeleteVRFKeySuccess() (*DeleteVRFKeySuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteVRFKeySuccessResolver(r.key), true
	}
	return nil, false
}
