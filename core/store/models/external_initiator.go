package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"golang.org/x/crypto/sha3"
)

// ExternalInitiator represents a user that can initiate runs remotely
type ExternalInitiator struct {
	*gorm.Model
	AccessKey    string
	Salt         string
	HashedSecret string
}

// NewExternalInitiator generates an ExternalInitiator from an
// ExternalInitiatorAuthentication, hashing the password for storage
func NewExternalInitiator(eia *ExternalInitiatorAuthentication) (*ExternalInitiator, error) {
	salt := NewSecret()
	hashedSecret, err := HashedSecret(eia, salt)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing secret for external initiator")
	}

	return &ExternalInitiator{
		AccessKey:    eia.AccessKey,
		HashedSecret: hashedSecret,
		Salt:         salt,
	}, nil
}

// AuthenticateExternalInitiator compares an auth against an initiator and
// returns true if the password hashes match
func AuthenticateExternalInitiator(eia *ExternalInitiatorAuthentication, ea *ExternalInitiator) (bool, error) {
	hashedSecret, err := HashedSecret(eia, ea.Salt)
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare([]byte(hashedSecret), []byte(ea.HashedSecret)) == 1, nil
}

// NewExternalInitiatorAuthentication returns a new
// ExternalInitiatorAuthentication with a freshly generated access key and
// secret, this is intended to be supplied to the user and saved, as it cannot
// be regenerated in the future.
func NewExternalInitiatorAuthentication() *ExternalInitiatorAuthentication {
	return &ExternalInitiatorAuthentication{
		AccessKey: utils.NewBytes32ID(),
		Secret:    NewSecret(),
	}
}

func hashInput(eia *ExternalInitiatorAuthentication, salt string) []byte {
	return []byte(fmt.Sprintf("v0-%s-%s-%s", eia.AccessKey, eia.Secret, salt))
}

// HashedSecret generates a hashed password for an external initiator
// authentication
func HashedSecret(eia *ExternalInitiatorAuthentication, salt string) (string, error) {
	hasher := sha3.New256()
	_, err := hasher.Write(hashInput(eia, salt))
	if err != nil {
		return "", errors.Wrap(err, "error writing external initiator authentication to hasher")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// ExternalInitiatorAuthentication represents the credentials needed to
// authenticate as an external initiator
type ExternalInitiatorAuthentication struct {
	AccessKey string
	Secret    string
}

// GetID returns the ID of this structure for jsonapi serialization.
func (eia *ExternalInitiatorAuthentication) GetID() string {
	return eia.AccessKey
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (eia *ExternalInitiatorAuthentication) GetName() string {
	return "external_initiators"
}

// SetID returns the ID of this structure for jsonapi serialization.
func (eia *ExternalInitiatorAuthentication) SetID(id string) error {
	eia.AccessKey = id
	return nil
}

var characters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numChars *big.Int

func init() {
	numChars = big.NewInt(int64(len(characters)))
}

// NewSecret returns a new secret for use for authenticating external initiators
func NewSecret() string {
	b := make([]rune, 64)
	for i := range b {
		cidx, err := rand.Int(rand.Reader, numChars)
		if err != nil {
			panic(errors.Wrapf(err, "while generating secret for external initiator"))
		}
		b[i] = characters[cidx.Int64()]
	}
	return string(b)
}
