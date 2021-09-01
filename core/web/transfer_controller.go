package web

import (
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

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

	chain, err := getChain(c, tc.App.GetChainSet(), tr.EVMChainID.String())
	switch err {
	case ErrInvalidChainID, ErrMultipleChains, ErrMissingChainID:
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	db := tc.App.GetStore().DB

	etx, err := bulletprooftxmanager.SendEther(db, chain.ID(), tr.FromAddress, tr.DestinationAddress, tr.Amount, chain.Config().EvmGasLimitTransfer())
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %v", err))
		return
	}

	jsonAPIResponse(c, presenters.NewEthTxResource(etx), "eth_tx")
}
