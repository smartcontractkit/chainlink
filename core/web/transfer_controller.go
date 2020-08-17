package web

import (
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/gin-gonic/gin"
)

// TransfersController can send LINK tokens to another address
type TransfersController struct {
	App chainlink.Application
}

// Create sends ETH from the Chainlink's account to a specified address.
//
// Example: "<application>/withdrawals"
func (tc *TransfersController) Create(c *gin.Context) {
	var tr models.SendEtherRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	store := tc.App.GetStore()

	if store.Config.EnableBulletproofTxManager() {
		etx, err := bulletprooftxmanager.SendEther(store, tr.FromAddress, tr.DestinationAddress, tr.Amount)
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %v", err))
			return
		}

		jsonAPIResponse(c, etx, "eth_tx")
	} else {
		tx, err := store.TxManager.CreateTxWithEth(tr.FromAddress, tr.DestinationAddress, &tr.Amount)
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %v", err))
			return
		}

		jsonAPIResponse(c, presenters.NewTx(tx), "transaction")
	}
}
