package auth

import (
	"encoding/hex"
	"fmt"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

var (
	// ErrorAuthFailed is a generic authentication failed - but not because of
	// some system failure on our behalf (i.e. HTTP 5xx), more detail is not
	// given
	ErrorAuthFailed = errors.New("Authentication failed")
)

// Token is used for API authentication.
type Token struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// GetID returns the ID of the Token struct for jsonapi serialization.
func (ta *Token) GetID() string {
	return ta.AccessKey
}

// GetName returns the pluralized "type" of Token struct for jsonapi serialization.
func (ta *Token) GetName() string {
	return "auth_tokens"
}

// SetID returns the ID of Token struct for jsonapi serialization.
func (ta *Token) SetID(id string) error {
	ta.AccessKey = id
	return nil
}

// NewToken returns a new Authentication Token.
func NewToken() *Token {
	return &Token{
		AccessKey: utils.NewBytes32ID(),
		SecretKey: utils.NewSecret(utils.DefaultSecretSize),
	}
}

// HashInput gets both access key and secret key with additional salt for HashedSecret function
func HashInput(ta *Token, salt string) []byte {
	return []byte(fmt.Sprintf("v0-%s-%s-%s", ta.AccessKey, ta.SecretKey, salt))
}

// HashedSecret generates a hashed password for an external initiator
// authentication
func HashedSecret(ta *Token, salt string) (string, error) {
	hasher := sha3.New256()
	_, err := hasher.Write(HashInput(ta, salt))
	if err != nil {
		return "", errors.Wrap(err, "error writing external initiator authentication to hasher")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}