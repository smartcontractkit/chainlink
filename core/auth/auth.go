package auth

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"hash"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// ErrorAuthFailed is a generic authentication failed - but not because of
	// some system failure on our behalf (i.e. HTTP 5xx), more detail is not
	// given
	ErrorAuthFailed = pkgerrors.New("Authentication failed")
)

// Token is used for API authentication.
type Token struct {
	AccessKey string `json:"accessKey"`
	Secret    string `json:"secret"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (ta *Token) GetID() string {
	return ta.AccessKey
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (ta *Token) GetName() string {
	return "auth_tokens"
}

// SetID returns the ID of this structure for jsonapi serialization.
func (ta *Token) SetID(id string) error {
	ta.AccessKey = id
	return nil
}

// NewToken returns a new Authentication Token.
func NewToken() *Token {
	return &Token{
		AccessKey: utils.NewBytes32ID(),
		Secret:    utils.NewSecret(utils.DefaultSecretSize),
	}
}

func hashInput(ta *Token, salt string) []byte {
	return fmt.Appendf(nil, "v0-%s-%s-%s", ta.AccessKey, ta.Secret, salt)
}

// HashedSecret generates a hashed password for an external initiator
// authentication
func HashedSecret(ta *Token, salt string) (string, error) {
	hasher := hash.Hash(sha3.New256())
	_, err := hasher.Write(hashInput(ta, salt))
	if err != nil {
		return "", pkgerrors.Wrap(err, "error writing external initiator authentication to hasher")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
