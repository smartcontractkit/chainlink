package models

import (
	"errors"
	"regexp"
	"time"

	"chainlink/core/utils"
)

// User holds the credentials for API user.
type User struct {
	Email          string    `json:"email" gorm:"primary_key"`
	HashedPassword string    `json:"hashedPassword"`
	CreatedAt      time.Time `json:"createdAt" gorm:"index"`
	TokenKey       string    `json:"tokenKey"`
	TokenSecret    string    `json:"tokenSecret"`
}

// https://davidcel.is/posts/stop-validating-email-addresses-with-regex/
var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// NewUser creates a new user by hashing the passed plainPwd with bcrypt.
func NewUser(email, plainPwd string) (User, error) {
	if len(email) == 0 {
		return User{}, errors.New("Must enter an email")
	}

	if !emailRegexp.MatchString(email) {
		return User{}, errors.New("Invalid email format")
	}

	if len(plainPwd) < 8 || len(plainPwd) > 1028 {
		return User{}, errors.New("Must enter a password with 8 - 1028 characters")
	}

	pwd, err := utils.HashPassword(plainPwd)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:          email,
		HashedPassword: pwd,
	}, nil
}

// SessionRequest encapsulates the fields needed to generate a new SessionID,
// including the hashed password.
type SessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Session holds the unique id for the authenticated session.
type Session struct {
	ID        string    `json:"id" gorm:"primary_key"`
	LastUsed  time.Time `json:"lastUsed" gorm:"index"`
	CreatedAt time.Time `json:"createdAt" gorm:"index"`
}

// NewSession returns a session instance with ID set to a random ID and
// LastUsed to to now.
func NewSession() Session {
	return Session{
		ID:       utils.NewBytes32ID(),
		LastUsed: time.Now(),
	}
}

// ChangePasswordRequest sets a new password for the current Session's User.
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
