package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
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
	from, err := retrieveFromAddress(tr.FromAddress, store)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	tx, err := store.TxManager.CreateTxWithEth(from, tr.DestinationAddress, tr.Amount)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %v", err))
		return
	}

	jsonAPIResponse(c, presenters.NewTx(tx), "transaction")
}

func retrieveFromAddress(from common.Address, store *store.Store) (common.Address, error) {
	if from != utils.ZeroAddress {
		return from, nil
	}
	ma := store.TxManager.NextActiveAccount()
	if ma == nil {
		return common.Address{}, errors.New("must activate an account before creating a transaction")
	}

	return ma.Address, nil
}
