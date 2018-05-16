package web

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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
//  "<application>/account"
func (jsc *AccountBalanceController) Show(c *gin.Context) {
	store := jsc.App.Store
	txm := store.TxManager

	if account, err := store.KeyStore.GetAccount(); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if ethBalance, err := txm.GetEthBalance(account.Address); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if linkBalance, err := txm.GetLinkBalance(account.Address, common.HexToAddress(store.Config.LinkContractAddress)); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		ab := presenters.AccountBalance{
			Address:     account.Address.Hex(),
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
		}
		if json, err := jsonapi.Marshal(ab); err != nil {
			c.JSON(500, gin.H{
				"errors": []string{fmt.Errorf("failed to marshal account using jsonapi: %+v", err).Error()},
			})
		} else {
			c.Data(200, MediaType, json)
		}
	}
}
