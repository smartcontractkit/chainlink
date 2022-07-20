package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/sessions"
)

// UserResource represents a User JSONAPI resource.
type UserResource struct {
	JAID
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// GetName implements the api2go EntityNamer interface
func (r UserResource) GetName() string {
	return "users"
}

// NewUserResource constructs a new UserResource.
//
// A User does not have an ID primary key, so we must use the email
func NewUserResource(u sessions.User) *UserResource {
	return &UserResource{
		JAID:      NewJAID(u.Email),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
