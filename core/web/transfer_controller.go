package web

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
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
		jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	if !tr.AllowHigherAmounts {
		var balance *big.Int

		balanceMonitor := chain.BalanceMonitor()

		if balanceMonitor != nil {
			balance = balanceMonitor.GetEthBalance(tr.FromAddress).ToInt()
		} else {
			balance, err = chain.Client().BalanceAt(c, tr.FromAddress, nil)
			if err != nil {
				jsonAPIError(c, http.StatusInternalServerError, err)
				return
			}
		}

		zero := big.NewInt(0)

		if balance == nil || balance.Cmp(zero) == 0 {
			jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("balance is too low for this transaction to be executed: %v", balance))
			return
		}

		var gasPrice *big.Int

		gasLimit := chain.Config().EvmGasLimitTransfer()
		estimator := chain.TxManager().GetGasEstimator()

		if chain.Config().EvmEIP1559DynamicFees() {
			gasPrice = chain.Config().EvmMaxGasPriceWei()
			_, gasLimit, err = estimator.GetDynamicFee(gasLimit)
			if err != nil {
				jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "failed to get dynamic gas fee"))
				return
			}
		} else {
			gasPrice, gasLimit, err = estimator.GetLegacyGas(nil, gasLimit)
			if err != nil {
				jsonAPIError(c, http.StatusUnprocessableEntity, errors.Wrap(err, "failed to estimate gas"))
				return
			}
		}

		intAmount := utils.NewBig(tr.Amount.ToInt())
		fee := gasPrice.Mul(gasPrice, utils.NewBigI(int64(gasLimit)).ToInt())

		amountWithFees := intAmount.ToInt().Add(intAmount.ToInt(), fee)

		// ETH balance is less than the sent amount + fees
		if balance.Cmp(amountWithFees) < 0 {
			jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("balance is too low for this transaction to be executed: %v", balance))
			return
		}
	}

	etx, err := chain.TxManager().SendEther(chain.ID(), tr.FromAddress, tr.DestinationAddress, tr.Amount, chain.Config().EvmGasLimitTransfer())
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %v", err))
		return
	}

	jsonAPIResponse(c, presenters.NewEthTxResource(etx), "eth_tx")
}
