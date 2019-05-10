package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// UserController manages the current Session's User User.
type UserController struct {
	App services.Application
}

// UpdatePassword changes the password for the current User.
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var request models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		jsonAPIError(ctx, http.StatusUnprocessableEntity, err)
	} else if user, err := c.App.GetStore().FindUser(); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to obtain current user record: %+v", err))
	} else if !utils.CheckPasswordHash(request.OldPassword, user.HashedPassword) {
		jsonAPIError(ctx, http.StatusConflict, errors.New("Old password does not match"))
	} else if err := c.updateUserPassword(ctx, &user, request.NewPassword); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(ctx, presenters.UserPresenter{User: &user}, "user")
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
	if sessionID, err := c.getCurrentSessionID(ctx); err != nil {
		return err
	} else if err := c.App.GetStore().ClearNonCurrentSessions(sessionID); err != nil {
		return fmt.Errorf("failed to clear non current user sessions: %+v", err)
	} else if err := c.saveNewPassword(user, newPassword); err != nil {
		return fmt.Errorf("failed to update current user password: %+v", err)
	}
	return nil
}

func getAccountBalanceFor(ctx *gin.Context, store *store.Store, account accounts.Account) presenters.AccountBalance {
	txm := store.TxManager
	if ethBalance, err := txm.GetEthBalance(account.Address); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
	} else if linkBalance, err := txm.GetLINKBalance(account.Address); err != nil {
		jsonAPIError(ctx, http.StatusInternalServerError, err)
	} else {
		return presenters.AccountBalance{
			Address:     account.Address.Hex(),
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
		}
	}
	return presenters.AccountBalance{}
}
