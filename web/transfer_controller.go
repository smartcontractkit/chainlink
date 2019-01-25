package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
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
	} else if tr.FromAddress != utils.ZeroAddress {
		if tx, err := store.TxManager.CreateTxWithEth(tr.DestinationAddress, tr.Amount, tr.FromAddress); err != nil {
			publicError(c, 400, fmt.Errorf("Transaction failed: %v", err))
		} else {
			c.JSON(200, tx)
		}
	} else if tx, err := store.TxManager.CreateTxWithEth(tr.DestinationAddress, tr.Amount); err != nil {
		publicError(c, 400, fmt.Errorf("Transaction failed: %v", err))
	} else {
		c.JSON(200, tx)
	}
}
