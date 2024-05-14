package sessions

import (
	"fmt"
	"net/mail"
	"time"

	pkgerrors "github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// User holds the credentials for API user.
type User struct {
	Email             string
	HashedPassword    string
	Role              UserRole
	CreatedAt         time.Time
	TokenKey          null.String
	TokenSalt         null.String
	TokenHashedSecret null.String
	UpdatedAt         time.Time
}

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleEdit  UserRole = "edit"
	UserRoleRun   UserRole = "run"
	UserRoleView  UserRole = "view"
)

// https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length
const (
	MaxBcryptPasswordLength = 50
)

// NewUser creates a new user by hashing the passed plainPwd with bcrypt.
func NewUser(email string, plainPwd string, role UserRole) (User, error) {
	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	pwd, err := ValidateAndHashPassword(plainPwd)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:          email,
		HashedPassword: pwd,
		Role:           role,
	}, nil
}

// ValidateEmail is the single point of logic for user email validations
func ValidateEmail(email string) error {
	if len(email) == 0 {
		return pkgerrors.New("Must enter an email")
	}
	_, err := mail.ParseAddress(email)
	return err
}

// ValidateAndHashPassword is the single point of logic for user password validations
func ValidateAndHashPassword(plainPwd string) (string, error) {
	if err := utils.VerifyPasswordComplexity(plainPwd); err != nil {
		return "", pkgerrors.Wrapf(err, "password insufficiently complex:\n%s", utils.PasswordComplexityRequirements)
	}
	if len(plainPwd) > MaxBcryptPasswordLength {
		return "", pkgerrors.Errorf("must enter a password less than %v characters", MaxBcryptPasswordLength)
	}

	pwd, err := utils.HashPassword(plainPwd)
	if err != nil {
		return "", err
	}

	return pwd, nil
}

// GetUserRole is the single point of logic for mapping role string to UserRole
func GetUserRole(role string) (UserRole, error) {
	if role == string(UserRoleAdmin) {
		return UserRoleAdmin, nil
	}
	if role == string(UserRoleEdit) {
		return UserRoleEdit, nil
	}
	if role == string(UserRoleRun) {
		return UserRoleRun, nil
	}
	if role == string(UserRoleView) {
		return UserRoleView, nil
	}

	errStr := fmt.Sprintf(
		"Invalid role: %s. Allowed roles: '%s', '%s', '%s', '%s'.",
		role,
		UserRoleAdmin,
		UserRoleEdit,
		UserRoleRun,
		UserRoleView,
	)
	return UserRole(""), pkgerrors.New(errStr)
}
