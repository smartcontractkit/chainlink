package models

import (
	"math/rand"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/utils"
)

// ExternalInitiator represents a user that can initiate runs remotely
type ExternalInitiator struct {
	*gorm.Model
	AccessKey    string
	HashedSecret string
}

// NewExternalInitiator generates an ExternalInitiator from an
// ExternalInitiatorAuthentication, hashing the password for storage
func NewExternalInitiator(eia *ExternalInitiatorAuthentication) (*ExternalInitiator, error) {
	hashedSecret, err := utils.HashPassword(eia.Secret)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing secret for external initiator")
	}

	return &ExternalInitiator{
		AccessKey:    eia.AccessKey,
		HashedSecret: hashedSecret,
	}, nil
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

// NewSecret returns a new secret for use for authenticating external initiators
func NewSecret() string {
	var characters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 64)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}
