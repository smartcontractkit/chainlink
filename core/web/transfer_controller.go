package web

import (
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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

	chain, err := getChain(tc.App.GetChainSet(), tr.EVMChainID.String())
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

	if tr.FromAddress == utils.ZeroAddress {
		jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("invalid withdrawal source address address: %v", tr.FromAddress))
		return
	}

	if !tr.AllowHigherAmounts {
		balance := chain.BalanceMonitor().GetEthBalance(tr.FromAddress)

		// ETH balance is less than the sent amount
		if balance == nil || balance.Cmp(&tr.Amount) == -1 {
			jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("balance is too low for this transaction to be executed: %v", balance))
			return
		}
	}

	db := tc.App.GetSqlxDB()
	q := pg.NewQ(db, tc.App.GetLogger(), tc.App.GetConfig())
	etx, err := chain.TxManager().SendEther(q, chain.ID(), tr.FromAddress, tr.DestinationAddress, tr.Amount, chain.Config().EvmGasLimitTransfer())
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %v", err))
		return
	}

	jsonAPIResponse(c, presenters.NewEthTxResource(etx), "eth_tx")
}
