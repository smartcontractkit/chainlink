package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

// UserController manages the current Session's User User.
type UserController struct {
	App *services.ChainlinkApplication
}

func (c *UserController) getCurrentSessionID(ctx *gin.Context) (string, error) {
	session := sessions.Default(ctx)
	sessionID, ok := session.Get(SessionIDKey).(string)
	if !ok {
		return "", errors.New("unable to get current session ID")
	}
	return sessionID, nil
}

func (c *UserController) clearNonCurrentSessions(sessionID string) error {
	var sessions []models.Session
	err := c.App.Store.Select(q.Not(q.Eq("ID", sessionID))).Find(&sessions)
	if err != nil && err != storm.ErrNotFound {
		return err
	}

	for _, s := range sessions {
		err := c.App.Store.DeleteStruct(&s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *UserController) saveNewPassword(user *models.User, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.HashedPassword = hashedPassword
	return c.App.Store.Save(user)
}

func (c *UserController) updateUserPassword(ctx *gin.Context, user *models.User, newPassword string) error {
	if sessionID, err := c.getCurrentSessionID(ctx); err != nil {
		return err
	} else if err := c.clearNonCurrentSessions(sessionID); err != nil {
		return fmt.Errorf("failed to clear non current user sessions: %+v", err)
	} else if err := c.saveNewPassword(user, newPassword); err != nil {
		return fmt.Errorf("failed to update current user password: %+v", err)
	}
	return nil
}

// UpdatePassword changes the password for the current User.
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var request models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		publicError(ctx, http.StatusUnprocessableEntity, err)
	} else if user, err := c.App.Store.FindUser(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
	} else if !utils.CheckPasswordHash(request.OldPassword, user.HashedPassword) {
		publicError(ctx, http.StatusConflict, errors.New("Old password does not match"))
	} else if err := c.updateUserPassword(ctx, &user, request.NewPassword); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else if json, err := jsonapi.Marshal(presenters.UserPresenter{User: &user}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to marshal password reset response using jsonapi: %+v", err))
	} else {
		ctx.Data(http.StatusOK, MediaType, json)
	}
}

// AccountBalances returns the account balances of ETH & LINK.
// Example:
//  "<application>/user/balances"
func (c *UserController) AccountBalances(ctx *gin.Context) {
	store := c.App.Store
	txm := store.TxManager

	if account, err := store.KeyStore.GetAccount(); err != nil {
		publicError(ctx, 400, err)
	} else if ethBalance, err := txm.GetEthBalance(account.Address); err != nil {
		ctx.AbortWithError(500, err)
	} else if linkBalance, err := txm.GetLinkBalance(account.Address); err != nil {
		ctx.AbortWithError(500, err)
	} else {
		ab := presenters.AccountBalance{
			Address:     account.Address.Hex(),
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
		}
		if json, err := jsonapi.Marshal(ab); err != nil {
			ctx.AbortWithError(500, fmt.Errorf("failed to marshal account using jsonapi: %+v", err))
		} else {
			ctx.Data(200, MediaType, json)
		}
	}
}
