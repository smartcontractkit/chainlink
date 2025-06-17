package web

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	commontxmgr "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/gin-gonic/gin"
)

// EVMTransfersController can send LINK tokens to another address
type EVMTransfersController struct {
	App chainlink.Application
}

// Create sends ETH from the Chainlink's account to a specified address.
//
// Example: "<application>/withdrawals"
func (tc *EVMTransfersController) Create(c *gin.Context) {
	var tr models.SendEtherRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	chain, err := getChain(tc.App.GetRelayers().LegacyEVMChains(), tr.EVMChainID.String())
	if err != nil {
		if errors.Is(err, ErrInvalidChainID) || errors.Is(err, ErrMultipleChains) || errors.Is(err, ErrMissingChainID) {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if tr.FromAddress == utils.ZeroAddress {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	if !tr.AllowHigherAmounts {
		err = ValidateEthBalanceForTransfer(c, chain, tr.FromAddress, tr.Amount, tr.DestinationAddress)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("transaction failed: %v", err))
			return
		}
	}

	etx, err := chain.TxManager().SendNativeToken(c, chain.ID(), tr.FromAddress, tr.DestinationAddress, *tr.Amount.ToInt(), chain.Config().EVM().GasEstimator().LimitTransfer())
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, errors.Errorf("transaction failed: %v", err))
		return
	}

	tc.App.GetAuditLogger().Audit(audit.EthTransactionCreated, map[string]interface{}{
		"ethTX": etx,
	})

	// skip waiting for txmgr to create TxAttempt
	if tr.SkipWaitTxAttempt {
		jsonAPIResponse(c, presenters.NewEthTxResource(etx), "eth_tx")
		return
	}

	timeout := 10 * time.Second // default
	if tr.WaitAttemptTimeout != nil {
		timeout = *tr.WaitAttemptTimeout
	}

	// wait and retrieve tx attempt matching tx ID
	attempt, err := FindTxAttempt(c, timeout, etx, tc.App.TxmStorageService().FindTxWithAttempts)
	if err != nil {
		jsonAPIError(c, http.StatusGatewayTimeout, fmt.Errorf("failed to find transaction within timeout: %w", err))
		return
	}
	jsonAPIResponse(c, presenters.NewEthTxResourceFromAttempt(attempt), "eth_tx")
}

// ValidateEthBalanceForTransfer validates that the current balance can cover the transaction amount
func ValidateEthBalanceForTransfer(c *gin.Context, chain legacyevm.Chain, fromAddr common.Address, amount assets.Eth, toAddr common.Address) error {
	var err error
	var balance *big.Int

	balanceMonitor := chain.BalanceMonitor()

	if balanceMonitor != nil {
		balance = balanceMonitor.GetEthBalance(fromAddr).ToInt()
	} else {
		balance, err = chain.Client().BalanceAt(c, fromAddr, nil)
		if err != nil {
			return err
		}
	}

	zero := big.NewInt(0)

	if balance == nil || balance.Cmp(zero) == 0 {
		return errors.Errorf("balance is too low for this transaction to be executed: %v", balance)
	}

	gasLimit := chain.Config().EVM().GasEstimator().LimitTransfer()
	estimator := chain.GasEstimator()

	amountWithFees, err := estimator.GetMaxCost(c, amount, nil, gasLimit, chain.Config().EVM().GasEstimator().PriceMaxKey(fromAddr), &fromAddr, &toAddr)
	if err != nil {
		return err
	}
	if balance.Cmp(amountWithFees) < 0 {
		// ETH balance is less than the sent amount + fees
		return errors.Errorf("balance is too low for this transaction to be executed: %v", balance)
	}

	return nil
}

func FindTxAttempt(ctx context.Context, timeout time.Duration, etx txmgr.Tx, FindTxWithAttempts func(context.Context, int64) (txmgr.Tx, error)) (attempt txmgr.TxAttempt, err error) {
	recheckTime := time.Second
	tick := time.After(0)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return attempt, fmt.Errorf("%w - tx may still have been broadcast", ctx.Err())
		case <-tick:
			etx, err = FindTxWithAttempts(ctx, etx.ID)
			if err != nil {
				return attempt, fmt.Errorf("failed to find transaction: %w", err)
			}
		}

		// exit if tx attempts are found
		// also validate etx.State != unstarted (ensure proper tx state for tx with attempts)
		if len(etx.TxAttempts) > 0 && etx.State != commontxmgr.TxUnstarted {
			break
		}
		tick = time.After(recheckTime)
	}

	// attach original tx to attempt
	attempt = etx.TxAttempts[0]
	attempt.Tx = etx
	return attempt, nil
}
