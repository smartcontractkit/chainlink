package models

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"golang.org/x/crypto/sha3"
)

// ExternalInitiatorRequest is the incoming record used to create an ExternalInitiator.
type ExternalInitiatorRequest struct {
	Name string `json:"name"`
	URL  WebURL `json:"url"`
}

// ExternalInitiator represents a user that can initiate runs remotely
type ExternalInitiator struct {
	*gorm.Model
	Name           string `gorm:"not null,unique"`
	URL            WebURL `gorm:"not null"`
	AccessKey      string `gorm:"not null"`
	Salt           string `gorm:"not null"`
	HashedSecret   string `gorm:"not null"`
	OutgoingSecret string `gorm:"not null"`
	OutgoingToken  string `gorm:"not null"`
}

// NewExternalInitiator generates an ExternalInitiator from an
// ExternalInitiatorAuthentication, hashing the password for storage
func NewExternalInitiator(
	eia *ExternalInitiatorAuthentication,
	eir *ExternalInitiatorRequest,
) (*ExternalInitiator, error) {
	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := HashedSecret(eia, salt)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing secret for external initiator")
	}

	return &ExternalInitiator{
		Name:           strings.ToLower(eir.Name),
		URL:            eir.URL,
		AccessKey:      eia.AccessKey,
		HashedSecret:   hashedSecret,
		Salt:           salt,
		OutgoingToken:  utils.NewSecret(utils.DefaultSecretSize),
		OutgoingSecret: utils.NewSecret(utils.DefaultSecretSize),
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
		Secret:    utils.NewSecret(utils.DefaultSecretSize),
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
