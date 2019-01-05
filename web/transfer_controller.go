package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
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
		return
	} else if tx, err := store.TxManager.CreateTx(tr.DestinationAddress, []byte{}); err != nil {
		publicError(c, 400, fmt.Errorf("Transaction failed: %v", err))
	} else {
		c.JSON(200, tx)
	}
}
