package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	clsession "github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/utils"
	webauth "github.com/smartcontractkit/chainlink/core/web/auth"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// UserController manages the current Session's User.
type UserController struct {
	App chainlink.Application
}

// UpdatePasswordRequest defines the request to set a new password for the
// current session's User.
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Index lists all API users
func (c *UserController) Index(ctx *gin.Context) {
	users, err := c.App.SessionORM().ListUsers()
	if err != nil {
		c.App.GetLogger().Errorf("Unable to list users", "err", err)
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(ctx, presenters.NewUserResources(users), "users")
}

// Create creates a new API user with provided context arugments.
func (c *UserController) Create(ctx *gin.Context) {
	type newUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var request newUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	userRole, err := clsession.GetUserRole(request.Role)
	if err != nil {
		jsonAPIError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := clsession.NewUser(request.Email, request.Password, userRole)
	if err != nil {
		jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("error creating API user: %s", err))
		return
	}
	if err = c.App.SessionORM().CreateUser(&user); err != nil {
		// If this is a duplicate key error (code 23505), return a nicer error message
		if err, ok := err.(*pgconn.PgError); ok {
			fmt.Println(err.Code)
			if err.Code == "23505" {
				jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("user with email %s already exists", request.Email))
				return
			}
		}
		c.App.GetLogger().Errorf("Error creating new API user", "err", err)
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("error creating API user"))
		return
	}

	jsonAPIResponse(ctx, presenters.NewUserResource(user), "user")
}

// Update changes sets email, password, or role fields of a specified API user.
func (c *UserController) Update(ctx *gin.Context) {
	type updateUserRequest struct {
		Email       string `json:"email"`
		NewEmail    string `json:"newEmail"`
		NewPassword string `json:"newPassword"`
		NewRole     string `json:"newRole"`
	}

	var request updateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	// First, attempt to load specified user by email
	// This user object will have its fields patched for a final update call
	user, err := c.App.SessionORM().FindUser(request.Email)
	if err != nil {
		jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("specified user not found: %s", request.Email))
		return
	}

	// Don't allow current admin user to edit self
	sessionUser, ok := webauth.GetAuthenticatedUser(ctx)
	if !ok {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user from context"))
		return
	}
	if sessionUser.Email == request.Email {
		jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("can not change state or permissions of current admin user"))
		return
	}

	// Changes to the user may require the entries in the sessions table to be removed
	purgeSessions := false

	// Patch validated email is newEmail is specified
	if request.NewEmail != "" {
		if err := clsession.ValidateEmail(request.NewEmail); err != nil {
			jsonAPIError(ctx, http.StatusBadRequest, err)
			return
		}
		user.Email = request.NewEmail

		// Email changed, purge associated sessions at final step
		purgeSessions = true
	}

	// Patch validated password is newPassword is specified
	if request.NewPassword != "" {
		pwd, err := clsession.ValidateAndHashPassword(request.NewPassword)
		if err != nil {
			jsonAPIError(ctx, http.StatusBadRequest, err)
			return
		}
		user.HashedPassword = pwd

		// Password changed, purge associated sessions at final step
		purgeSessions = true
	}

	// Patch validated role is newRole is specified
	if request.NewRole != "" {
		userRole, err := clsession.GetUserRole(request.NewRole)
		if err != nil {
			jsonAPIError(ctx, http.StatusBadRequest, err)
			return
		}
		user.Role = userRole
	}

	if purgeSessions {
		err := c.App.SessionORM().DeleteSessionByUser(request.Email)
		if err != nil {
			c.App.GetLogger().Errorf("Failed to purge user sessions via DeleteSessionByUser", "err", err)
			jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("error updating API user"))
			return
		}
	}

	if err = c.App.SessionORM().UpdateUser(request.Email, &user); err != nil {
		// If this is a duplicate key error (code 23505), return a nicer error message
		if err, ok := err.(*pgconn.PgError); ok {
			fmt.Println(err.Code)
			if err.Code == "23505" {
				jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("user with email %s already exists", request.Email))
				return
			}
		}
		c.App.GetLogger().Errorf("Error updating new API user", "err", err)
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("error updating API user"))
		return
	}

	jsonAPIResponse(ctx, presenters.NewUserResource(user), "user")
}

// Delete deletes an API user and any sessions by email
func (c *UserController) Delete(ctx *gin.Context) {
	email := ctx.Param("email")

	// Attempt find user by email
	_, err := c.App.SessionORM().FindUser(email)
	if err != nil {
		jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("specified user not found: %s", email))
		return
	}

	// Don't allow current admin user to delete self
	sessionUser, ok := webauth.GetAuthenticatedUser(ctx)
	if !ok {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user from context"))
		return
	}
	if sessionUser.Email == email {
		jsonAPIError(ctx, http.StatusBadRequest, pkgerrors.Errorf("can not delete currently logged in admin user"))
		return
	}

	if err = c.App.SessionORM().DeleteUser(email); err != nil {
		c.App.GetLogger().Errorf("Error deleting API user", "err", err, "email", email)
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("error deleting API user"))
		return
	}

	jsonAPIResponse(ctx, presenters.NewUserResource(clsession.User{Email: email}), "user")
}

// UpdatePassword changes the password for the current User.
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var request UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	sessionUser, ok := webauth.GetAuthenticatedUser(ctx)
	if !ok {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user from context"))
		return
	}
	user, err := c.App.SessionORM().FindUser(sessionUser.Email)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user record: %+v", err))
		return
	}
	if !utils.CheckPasswordHash(request.OldPassword, user.HashedPassword) {
		jsonAPIError(ctx, http.StatusConflict, errors.New("old password does not match"))
		return
	}
	if err := c.updateUserPassword(ctx, &user, request.NewPassword); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(ctx, presenters.NewUserResource(user), "user")
}

// NewAPIToken generates a new API token for a user overwriting any pre-existing one set.
func (c *UserController) NewAPIToken(ctx *gin.Context) {
	var request clsession.ChangeAuthTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	sessionUser, ok := webauth.GetAuthenticatedUser(ctx)
	if !ok {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user from context"))
		return
	}
	user, err := c.App.SessionORM().FindUser(sessionUser.Email)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user record: %+v", err))
		return
	}
	if !utils.CheckPasswordHash(request.Password, user.HashedPassword) {
		jsonAPIError(ctx, http.StatusUnauthorized, errors.New("incorrect password"))
		return
	}
	newToken := auth.NewToken()
	if err := c.App.SessionORM().SetAuthToken(&user, newToken); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(ctx, newToken, "auth_token", http.StatusCreated)
}

// DeleteAPIToken deletes and disables a user's API token.
func (c *UserController) DeleteAPIToken(ctx *gin.Context) {
	var request clsession.ChangeAuthTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	sessionUser, ok := webauth.GetAuthenticatedUser(ctx)
	if !ok {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user from context"))
		return
	}
	user, err := c.App.SessionORM().FindUser(sessionUser.Email)
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, pkgerrors.Errorf("failed to obtain current user record: %+v", err))
		return
	}
	if !utils.CheckPasswordHash(request.Password, user.HashedPassword) {
		jsonAPIError(ctx, http.StatusUnauthorized, errors.New("incorrect password"))
		return
	}
	if err := c.App.SessionORM().DeleteAuthToken(&user); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}
	{
		jsonAPIResponseWithStatus(ctx, nil, "auth_token", http.StatusNoContent)
	}
}

func getCurrentSessionID(ctx *gin.Context) (string, error) {
	session := sessions.Default(ctx)
	sessionID, ok := session.Get(webauth.SessionIDKey).(string)
	if !ok {
		return "", errors.New("unable to get current session ID")
	}
	return sessionID, nil
}

func (c *UserController) updateUserPassword(ctx *gin.Context, user *clsession.User, newPassword string) error {
	sessionID, err := getCurrentSessionID(ctx)
	if err != nil {
		return err
	}
	orm := c.App.SessionORM()
	if err := orm.ClearNonCurrentSessions(sessionID); err != nil {
		return pkgerrors.Errorf("failed to clear non current user sessions: %+v", err)
	}
	if err := orm.SetPassword(user, newPassword); err != nil {
		return pkgerrors.Errorf("failed to update current user password: %+v", err)
	}
	return nil
}
