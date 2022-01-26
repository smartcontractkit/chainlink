package presenters

import (
	"github.com/duo-labs/webauthn/protocol"
)

// RegistrationSettings represents an enrollment settings object
type RegistrationSettings struct {
	JAID
	Settings protocol.CredentialCreation `json:"settings"`
}

// GetName implements the api2go EntityNamer interface
func (r RegistrationSettings) GetName() string {
	return "registrationsettings"
}

// NewRegistrationSettings creates a new structure to enroll a new hardware
// key for authentication
func NewRegistrationSettings(settings protocol.CredentialCreation) *RegistrationSettings {
	return &RegistrationSettings{
		JAID:     NewJAID("registration_settings"),
		Settings: settings,
	}
}
