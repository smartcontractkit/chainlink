package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// WithdrawalsController can send LINK tokens to another address
type WithdrawalsController struct {
	App services.Application
}

var naz = assets.NewLink(1)

// Create sends LINK from the configured oracle contract to the given address
// See models.WithdrawalRequest for expected inputs. ContractAddress field is
// optional.
//
// Example: "<application>/withdrawals"
func (abc *WithdrawalsController) Create(c *gin.Context) {
	store := abc.App.GetStore()
	txm := store.TxManager
	wr := models.WithdrawalRequest{}

	if err := c.ShouldBindJSON(&wr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	addressWasNotSpecifiedInRequest := utils.ZeroAddress

	if wr.Amount.Cmp(naz) < 0 {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf(
			"Must withdraw at least %v LINK", naz.String()))
	} else if wr.DestinationAddress == addressWasNotSpecifiedInRequest {
		jsonAPIError(c, http.StatusBadRequest, errors.New("Invalid withdrawal address"))
	} else if linkBalance, err := txm.ContractLINKBalance(wr); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if linkBalance.Cmp(wr.Amount) < 0 {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf(
			"Insufficient link balance. Withdrawal Amount: %v "+
				"Link Balance: %v",
			wr.Amount.String(), linkBalance.String()))
		return
	}

	hash, err := txm.WithdrawLINK(wr)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponse(c, presenters.NewTx(&models.Tx{Hash: hash}), "transaction")
	}
}
