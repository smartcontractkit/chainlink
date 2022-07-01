package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/sessions"
)

// UserResource represents a User JSONAPI resource.
type UserResource struct {
	JAID
	Email             string            `json:"email"`
	Role              sessions.UserRole `json:"role"`
	HasActiveApiToken string            `json:"hasActiveApiToken"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r UserResource) GetName() string {
	return "users"
}

// NewUserResource constructs a new UserResource.
//
// A User does not have an ID primary key, so we must use the email
func NewUserResource(u sessions.User) *UserResource {
	hasToken := "false"
	if u.TokenKey.Valid {
		hasToken = "true"
	}
	return &UserResource{
		JAID:              NewJAID(u.Email),
		Email:             u.Email,
		Role:              sessions.UserRole(u.Role),
		HasActiveApiToken: hasToken,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

func NewUserResources(users []sessions.User) []UserResource {
	us := []UserResource{}
	for _, user := range users {
		us = append(us, *NewUserResource(user))
	}
	return us
}
