package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserController manages the current Session's User User.
type UserController struct {
	App chainlink.Application
}

// UpdatePassword changes the password for the current User.
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var request models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := c.App.GetStore().FindUser()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
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

	jsonAPIResponse(ctx, presenters.UserPresenter{User: &user}, "user")
}

// NewAPIToken generates a new API token for a user overwriting any pre-existing one set.
func (c *UserController) NewAPIToken(ctx *gin.Context) {
	var request models.ChangeAuthTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := c.App.GetStore().FindUser()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
		return
	}
	if !utils.CheckPasswordHash(request.Password, user.HashedPassword) {
		jsonAPIError(ctx, http.StatusUnauthorized, errors.New("incorrect password"))
		return
	}
	newToken, err := user.GenerateAuthToken()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}
	if err := c.App.GetStore().SaveUser(&user); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(ctx, newToken, "auth_token", http.StatusCreated)
}

// DeleteAPIToken deletes and disables a user's API token.
func (c *UserController) DeleteAPIToken(ctx *gin.Context) {
	var request models.ChangeAuthTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := c.App.GetStore().FindUser()
	if err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
		return
	}
	if !utils.CheckPasswordHash(request.Password, user.HashedPassword) {
		jsonAPIError(ctx, http.StatusUnauthorized, errors.New("incorrect password"))
		return
	}
	if user.DeleteAuthToken(); false {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}
	if err := c.App.GetStore().SaveUser(&user); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		return
	}
	{
		jsonAPIResponseWithStatus(ctx, nil, "auth_token", http.StatusNoContent)
	}
}

// AccountBalances returns the account balances of ETH & LINK.
// Example:
//  "<application>/user/balances"
func (c *UserController) AccountBalances(ctx *gin.Context) {
	store := c.App.GetStore()
	accounts := store.KeyStore.Accounts()
	balances := []presenters.AccountBalance{}
	for _, a := range accounts {
		pa := getAccountBalanceFor(ctx, store, a)
		if ctx.IsAborted() {
			return
		}
		balances = append(balances, pa)
	}

	jsonAPIResponse(ctx, balances, "balances")
}

func (c *UserController) getCurrentSessionID(ctx *gin.Context) (string, error) {
	session := sessions.Default(ctx)
	sessionID, ok := session.Get(SessionIDKey).(string)
	if !ok {
		return "", errors.New("unable to get current session ID")
	}
	return sessionID, nil
}

func (c *UserController) saveNewPassword(user *models.User, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.HashedPassword = hashedPassword
	return c.App.GetStore().SaveUser(user)
}

func (c *UserController) updateUserPassword(ctx *gin.Context, user *models.User, newPassword string) error {
	sessionID, err := c.getCurrentSessionID(ctx)
	if err != nil {
		return err
	}
	if err := c.App.GetStore().ClearNonCurrentSessions(sessionID); err != nil {
		return fmt.Errorf("failed to clear non current user sessions: %+v", err)
	}
	if err := c.saveNewPassword(user, newPassword); err != nil {
		return fmt.Errorf("failed to update current user password: %+v", err)
	}
	return nil
}

func getAccountBalanceFor(ctx *gin.Context, store *store.Store, account accounts.Account) presenters.AccountBalance {
	txm := store.TxManager
	ethBalance, err := store.EthClient.BalanceAt(context.TODO(), account.Address, nil)
	if err != nil {
		err = fmt.Errorf("error calling getEthBalance on Ethereum node: %v", err)
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		ctx.Abort()
		return presenters.AccountBalance{}
	}

	linkBalance, err := txm.GetLINKBalance(account.Address)
	if err != nil {
		err = fmt.Errorf("error calling getLINKBalance on Ethereum node: %v", err)
		jsonAPIError(ctx, http.StatusInternalServerError, err)
		ctx.Abort()
		return presenters.AccountBalance{}
	}

	return presenters.AccountBalance{
		Address:     account.Address.Hex(),
		EthBalance:  (*assets.Eth)(ethBalance),
		LinkBalance: linkBalance,
	}
}
