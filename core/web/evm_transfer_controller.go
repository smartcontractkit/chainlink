package web

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	commontxmgr "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
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

	// If LegacyEVMChains are available, use them; otherwise use the relayer.
	// Note that once we fully deprecate LegacyEVMChains we will switch to the relayer only.
	chain, errLegacy := getChain(tc.App.GetRelayers().LegacyEVMChains(), tr.EVMChainID.String()) //nolint:staticcheck //SA1019 keep the deprecated path for now
	if errLegacy == nil {
		tc.CreateEVMLegacy(c, chain, &tr)
	} else {
		relayer, errRelayer := tc.App.GetRelayers().Get(types.RelayID{Network: relay.NetworkEVM, ChainID: tr.EVMChainID.String()})
		if errRelayer != nil {
			jsonAPIError(c, http.StatusInternalServerError, errors.Join(errLegacy, errRelayer))
			return
		}
		tc.CreateWithRelayer(c, relayer, &tr)
	}
}

// CreateWithRelayer processes an ETH transfer request using the specified relayer and binds the result to the given context.
// As opposed to CreateEVMLegacy, it just sends the transaction without verifying its success.
func (tc *EVMTransfersController) CreateWithRelayer(c *gin.Context, relayer loop.Relayer, tr *models.SendEtherRequest) {
	info, err := relayer.GetChainInfo(c)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	cid, ok := new(big.Int).SetString(info.ChainID, 10)
	if !ok {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("could not parse chain ID: %s", info.ChainID))
		return
	}

	err = relayer.Transact(c.Request.Context(), tr.FromAddress.String(), tr.DestinationAddress.String(), tr.Amount.ToInt(), !tr.AllowHigherAmounts)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	resource := presenters.EthTxResource{
		From:       &tr.FromAddress,
		To:         &tr.DestinationAddress,
		Value:      tr.Amount.String(),
		EVMChainID: *sqlutil.New(cid),
	}

	tc.App.GetAuditLogger().Audit(audit.EthTransactionCreated, map[string]any{
		"ethTX": resource,
	})
	jsonAPIResponse(c, resource, "eth_tx")
}

// CreateEVMLegacy processes the EVM transfer using deprecated legacyevm.Chain.
// Functionally, it is similar to CreateWithRelayer with the difference that it waits for the transaction to be confirmed.
// Note that CreateEVMLegacy is deprecated and is scheduled for deletion.
func (tc *EVMTransfersController) CreateEVMLegacy(c *gin.Context, chain legacyevm.Chain, tr *models.SendEtherRequest) {
	if tr.FromAddress == utils.ZeroAddress {
		jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("withdrawal source address is missing: %v", tr.FromAddress))
		return
	}

	if !tr.AllowHigherAmounts {
		err := ValidateEthBalanceForTransfer(c, chain, tr.FromAddress, tr.Amount, tr.DestinationAddress)
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, fmt.Errorf("transaction failed: %w", err))
			return
		}
	}

	etx, err := chain.TxManager().SendNativeToken(c, chain.ID(), tr.FromAddress, tr.DestinationAddress, *tr.Amount.ToInt(), chain.Config().EVM().GasEstimator().LimitTransfer())
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("transaction failed: %w", err))
		return
	}

	tc.App.GetAuditLogger().Audit(audit.EthTransactionCreated, map[string]any{
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
		return fmt.Errorf("balance is too low for this transaction to be executed: %v", balance)
	}

	gasLimit := chain.Config().EVM().GasEstimator().LimitTransfer()
	estimator := chain.GasEstimator()

	amountWithFees, err := estimator.GetMaxCost(c, amount, nil, gasLimit, chain.Config().EVM().GasEstimator().PriceMaxKey(fromAddr), &fromAddr, &toAddr)
	if err != nil {
		return err
	}
	if balance.Cmp(amountWithFees) < 0 {
		// ETH balance is less than the sent amount + fees
		return fmt.Errorf("balance is too low for this transaction to be executed: %v", balance)
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
