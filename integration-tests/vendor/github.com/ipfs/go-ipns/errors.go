package ipns

import (
	"errors"
)

// ErrExpiredRecord should be returned when an ipns record is
// invalid due to being too old
var ErrExpiredRecord = errors.New("expired record")

// ErrUnrecognizedValidity is returned when an IpnsRecord has an
// unknown validity type.
var ErrUnrecognizedValidity = errors.New("unrecognized validity type")

// ErrInvalidPath should be returned when an ipns record path
// is not in a valid format
var ErrInvalidPath = errors.New("record path invalid")

// ErrSignature should be returned when an ipns record fails
// signature verification
var ErrSignature = errors.New("record signature verification failed")

// ErrKeyFormat should be returned when an ipns record key is
// incorrectly formatted (not a peer ID)
var ErrKeyFormat = errors.New("record key could not be parsed into peer ID")

// ErrPublicKeyNotFound should be returned when the public key
// corresponding to the ipns record path cannot be retrieved
// from the peer store
var ErrPublicKeyNotFound = errors.New("public key not found in peer store")

// ErrPublicKeyMismatch should be returned when the public key embedded in the
// record doesn't match the expected public key.
var ErrPublicKeyMismatch = errors.New("public key in record did not match expected pubkey")

// ErrBadRecord should be returned when an ipns record cannot be unmarshalled
var ErrBadRecord = errors.New("record could not be unmarshalled")
