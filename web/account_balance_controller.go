package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// AccountBalanceController returns information for the active account
type AccountBalanceController struct {
	App *services.ChainlinkApplication
}

// Show returns the address, plus it's ETH & LINK balance
// Example:
//  "<application>/account_balance"
func (jsc *AccountBalanceController) Show(c *gin.Context) {
	store := jsc.App.Store
	txm := store.TxManager

	if account, err := store.KeyStore.GetAccount(); err != nil {
		publicError(c, 400, err)
	} else if ethBalance, err := txm.GetEthBalance(account.Address); err != nil {
		c.AbortWithError(500, err)
	} else if linkBalance, err := txm.GetLinkBalance(account.Address); err != nil {
		c.AbortWithError(500, err)
	} else {
		ab := presenters.AccountBalance{
			Address:     account.Address.Hex(),
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
		}
		if json, err := jsonapi.Marshal(ab); err != nil {
			c.AbortWithError(500, fmt.Errorf("failed to marshal account using jsonapi: %+v", err))
		} else {
			c.Data(200, MediaType, json)
		}
	}
}
