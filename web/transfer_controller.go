package web

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// TransfersController can send LINK tokens to another address
type TransfersController struct {
	App services.Application
}

// Create sends ETH from the Chainlink's account to a specified address.
//
// Example: "<application>/withdrawals"
func (tc *TransfersController) Create(c *gin.Context) {
	var tr models.SendEtherRequest

	store := tc.App.GetStore()

	if err := c.ShouldBindJSON(&tr); err != nil {
		publicError(c, 400, err)
	} else if from, err := retrieveFromAddress(tr.FromAddress, store); err != nil {
		publicError(c, 400, err)
	} else if tx, err := store.TxManager.CreateTxWithEth(from, tr.DestinationAddress, tr.Amount); err != nil {
		publicError(c, 400, fmt.Errorf("Transaction failed: %v", err))
	} else {
		c.JSON(200, tx)
	}
}

func retrieveFromAddress(from common.Address, store *store.Store) (common.Address, error) {
	if from != utils.ZeroAddress {
		return from, nil
	}
	ma := store.TxManager.NextActiveAccount()
	if ma == nil {
		return common.Address{}, errors.New("Must activate an account before creating a transaction")
	}

	return ma.Address, nil
}
