package crib

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

const (
	solFunds = 1000
)

func distributeTransmitterFunds(lggr logger.Logger, nodeInfo []devenv.Node, env cldf.Environment, evmFundingEth uint64) error {
	evmFundingAmount := new(big.Int).Mul(deployment.UBigInt(evmFundingEth), deployment.UBigInt(1e18))

	g := new(errgroup.Group)

	// Handle EVM funding
	evmChains := env.BlockChains.EVMChains()
	if len(evmChains) > 0 {
		for sel, chain := range evmChains {
			g.Go(func() error {
				var evmAccounts []common.Address
				for _, n := range nodeInfo {
					chainID, err := chainsel.GetChainIDFromSelector(sel)
					if err != nil {
						lggr.Errorw("could not get chain id from selector", "selector", sel, "err", err)
						return err
					}
					addr := common.HexToAddress(n.AccountAddr[chainID])
					evmAccounts = append(evmAccounts, addr)
				}

				err := SendFundsToAccounts(env.GetContext(), lggr, chain, evmAccounts, evmFundingAmount, sel)
				if err != nil {
					lggr.Errorw("error funding evm accounts", "selector", sel, "err", err)
					return err
				}
				return nil
			})
		}
	}

	// Handle Solana funding
	solChains := env.BlockChains.SolanaChains()
	if len(solChains) > 0 {
		lggr.Info("Funding solana transmitters")
		for sel, chain := range solChains {
			g.Go(func() error {
				var solanaAddrs []solana.PublicKey
				for _, n := range nodeInfo {
					chainID, err := chainsel.GetChainIDFromSelector(sel)
					if err != nil {
						lggr.Errorw("could not get chain id from selector", "selector", sel, "err", err)
						return err
					}
					base58Addr := n.AccountAddr[chainID]
					lggr.Infof("Found %v solana transmitter address", base58Addr)

					pk, err := solana.PublicKeyFromBase58(base58Addr)
					if err != nil {
						lggr.Errorw("error converting base58 to solana PublicKey", "err", err, "address", base58Addr)
						return err
					}
					solanaAddrs = append(solanaAddrs, pk)
				}

				err := fundSolanaAccountsWithLogging(env.GetContext(), solanaAddrs, solFunds, chain.Client, lggr)
				if err != nil {
					lggr.Errorw("error funding solana accounts", "err", err, "selector", sel)
					return err
				}
				for _, addr := range solanaAddrs {
					res, err := chain.Client.GetBalance(env.GetContext(), addr, rpc.CommitmentFinalized)
					if err != nil {
						lggr.Errorw("failed to fetch transmitter balance", "transmitter", addr, "err", err)
						return err
					} else if res != nil {
						lggr.Infow("got balance for transmitter", "transmitter", addr, "balance", res.Value)
					}
				}
				return nil
			})
		}
	}

	return g.Wait()
}

func SendFundsToAccounts(ctx context.Context, lggr logger.Logger, chain cldf_evm.Chain, accounts []common.Address, fundingAmount *big.Int, sel uint64) error {
	nonce, err := chain.Client.NonceAt(ctx, chain.DeployerKey.From, nil)
	if err != nil {
		return fmt.Errorf("could not get latest nonce for deployer key: %w", err)
	}
	lggr.Infow("Starting funding process", "chain", sel, "fromAddress", chain.DeployerKey.From, "startNonce", nonce)

	tipCap, err := chain.Client.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("could not suggest gas tip cap: %w", err)
	}

	latestBlock, err := chain.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not get latest block: %w", err)
	}
	baseFee := latestBlock.BaseFee

	feeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		tipCap,
	)

	gasLimit, err := chain.Client.EstimateGas(ctx, ethereum.CallMsg{
		From:  chain.DeployerKey.From,
		To:    &accounts[0],
		Value: fundingAmount,
	})
	if err != nil {
		return fmt.Errorf("could not estimate gas for chain %d: %w", sel, err)
	}
	lggr.Infow("Using EIP-1559 fees", "chain", sel, "baseFee", baseFee, "tipCap", tipCap, "feeCap", feeCap, "gasLimit", gasLimit)

	var signedTxs []*gethtypes.Transaction

	chainID, err := chainsel.GetChainIDFromSelector(chain.Selector)
	if err != nil {
		return fmt.Errorf("could not get chainID from selector: %w", err)
	}
	chainIDBig := new(big.Int)
	if _, ok := chainIDBig.SetString(chainID, 10); !ok {
		return fmt.Errorf("could not parse chainID: %s", chainID)
	}

	for i, address := range accounts {
		currentNonce := nonce + uint64(i) //nolint:gosec // G115: i is always positive and within reasonable bounds
		baseTx := &gethtypes.DynamicFeeTx{
			ChainID:   chainIDBig,
			Nonce:     currentNonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &address,
			Value:     fundingAmount,
			Data:      nil,
		}
		tx := gethtypes.NewTx(baseTx)

		signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, tx)
		if err != nil {
			return fmt.Errorf("could not sign transaction for account %s: %w", address.Hex(), err)
		}
		signedTxs = append(signedTxs, signedTx)
	}

	for _, signedTx := range signedTxs {
		lggr.Infow("Sending funding tx", "chain", sel, "hash", signedTx.Hash().Hex(), "nonce", signedTx.Nonce())
		err = chain.Client.SendTransaction(ctx, signedTx)
		if err != nil {
			return fmt.Errorf("could not send transaction %s: %w", signedTx.Hash().Hex(), err)
		}
	}

	g, waitCtx := errgroup.WithContext(ctx)
	for _, tx := range signedTxs {
		g.Go(func() error {
			receipt, err := bind.WaitMined(waitCtx, chain.Client, tx)
			if err != nil {
				return fmt.Errorf("failed to mine transaction %s: %w", tx.Hash().Hex(), err)
			}
			if receipt.Status == gethtypes.ReceiptStatusFailed {
				return fmt.Errorf("transaction %s reverted", tx.Hash().Hex())
			}
			lggr.Infow("Transaction successfully mined", "chain", sel, "hash", tx.Hash().Hex())
			return nil
		})
	}

	return g.Wait()
}

// getTierChainSelectors organizes the provided chain selectors into deterministic tiers based on the supplied number of high and low tier chains.
func getTierChainSelectors(allSelectors []uint64, highTierCount int, lowTierCount int) (highTierSelectors []uint64, lowTierSelectors []uint64) {
	// we prioritize home selector, simulated solana, and evm feed selectors
	prioritySelectors := []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802}
	orderedSelectors := make([]uint64, 0, len(allSelectors))
	for _, sel := range prioritySelectors {
		if slices.Contains(allSelectors, sel) {
			orderedSelectors = append(orderedSelectors, sel)
		}
	}

	// the remaining chains are evm and count up
	evmChainID := 90000001
	for len(orderedSelectors) < len(allSelectors) {
		details, _ := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(evmChainID), chainsel.FamilyEVM)
		orderedSelectors = append(orderedSelectors, details.ChainSelector)
		evmChainID++
	}

	return orderedSelectors[0:highTierCount], orderedSelectors[highTierCount : highTierCount+lowTierCount]
}

// FundSolanaAccountsWithLogging requests airdrops for the provided accounts and waits for confirmation.
// It waits until all transactions reach at least "Confirmed" commitment level with enhanced logging and timeouts.
// Solana commitment levels: Processed < Confirmed < Finalized
// - Processed: Transaction processed by a validator but may be rolled back
// - Confirmed: Transaction confirmed by supermajority of cluster stake
// - Finalized: Transaction finalized and cannot be rolled back
func fundSolanaAccountsWithLogging(
	ctx context.Context, accounts []solana.PublicKey, solAmount uint64, solanaGoClient *rpc.Client,
	lggr logger.Logger,
) error {
	if len(accounts) == 0 {
		return nil
	}

	var sigs = make([]solana.Signature, 0, len(accounts))
	var successfulAccounts = make([]solana.PublicKey, 0, len(accounts))

	lggr.Infow("Starting Solana airdrop requests", "accountCount", len(accounts), "amountSOL", solAmount)

	// Request airdrops with better error tracking
	// Note: Using CommitmentConfirmed here means the RequestAirdrop call itself waits for confirmed status
	for i, account := range accounts {
		sig, err := solanaGoClient.RequestAirdrop(ctx, account, solAmount*solana.LAMPORTS_PER_SOL, rpc.CommitmentFinalized)
		if err != nil {
			// Return partial success information
			if len(sigs) > 0 {
				return fmt.Errorf("airdrop request failed for account %d (%s): %w (note: %d previous requests may have succeeded)",
					i, account.String(), err, len(sigs))
			}
			return fmt.Errorf("airdrop request failed for account %d (%s): %w", i, account.String(), err)
		}
		sigs = append(sigs, sig)
		successfulAccounts = append(successfulAccounts, account)

		lggr.Debugw("Airdrop request completed",
			"progress", fmt.Sprintf("%d/%d", i+1, len(accounts)),
			"account", account.String(),
			"signature", sig.String())

		// small delay to avoid rate limiting issues
		time.Sleep(100 * time.Millisecond)
	}

	// Adaptive timeout based on batch size - each airdrop can take several seconds
	// Base timeout of 30s + 5s per account for larger batches
	baseTimeout := 60 * time.Second
	if len(accounts) > 5 {
		baseTimeout += time.Duration(len(accounts)) * 5 * time.Second
	}
	timeout := baseTimeout
	const pollInterval = 500 * time.Millisecond

	lggr.Infow("Starting confirmation polling", "timeout", timeout, "accounts", len(accounts))

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	remaining := len(sigs)
	pollCount := 0
	for remaining > 0 {
		select {
		case <-timeoutCtx.Done():
			// Log which transactions are still unconfirmed for debugging
			unfinalizedSigs := []string{}
			statusRes, _ := solanaGoClient.GetSignatureStatuses(ctx, true, sigs...)
			if statusRes != nil && statusRes.Value != nil {
				for i, res := range statusRes.Value {
					if res == nil || res.ConfirmationStatus != rpc.ConfirmationStatusFinalized {
						unfinalizedSigs = append(unfinalizedSigs, fmt.Sprintf("%s (account: %s)",
							sigs[i].String(), successfulAccounts[i].String()))
					}
				}
			}
			lggr.Errorw("Timeout waiting for transaction confirmations",
				"remaining", remaining,
				"total", len(sigs),
				"timeout", timeout,
				"unfinalizedSigs", unfinalizedSigs)

			return fmt.Errorf("timeout waiting for transaction confirmations,"+
				"remaining: %d, total: %d, timeout: %s"+
				"unfinalizedSigs: %v",
				remaining, len(sigs), timeout, unfinalizedSigs)
		case <-ticker.C:
			pollCount++
			statusRes, sigErr := solanaGoClient.GetSignatureStatuses(timeoutCtx, true, sigs...)
			if sigErr != nil {
				return fmt.Errorf("failed to get signature statuses: %w", sigErr)
			}
			if statusRes == nil {
				return errors.New("signature status response is nil")
			}
			if statusRes.Value == nil {
				return errors.New("signature status response value is nil")
			}

			unfinalizedTxCount := 0
			for i, res := range statusRes.Value {
				if res == nil {
					// Transaction status not yet available
					unfinalizedTxCount++
					continue
				}

				if res.Err != nil {
					// Transaction failed
					lggr.Errorw("Transaction failed",
						"account", successfulAccounts[i].String(),
						"signature", sigs[i].String(),
						"error", res.Err)
					return fmt.Errorf("transaction failed for account %s (sig: %s): %v",
						successfulAccounts[i].String(), sigs[i].String(), res.Err)
				}

				// Check confirmation status - we want at least "Confirmed" level
				// Solana confirmation levels: Processed < Confirmed < Finalized
				switch res.ConfirmationStatus {
				case rpc.ConfirmationStatusProcessed, rpc.ConfirmationStatusConfirmed:
					// Still only processed, not yet confirmed
					unfinalizedTxCount++
				case rpc.ConfirmationStatusFinalized:
					// Transaction is finalized - we're good
					// Don't increment unfinalizedTxCount
				default:
					// Unknown status, treat as unconfirmed
					unfinalizedTxCount++
				}
			}
			remaining = unfinalizedTxCount

			// Log progress every 10 polls (5 seconds) for large batches
			if pollCount%10 == 0 {
				finalized := len(sigs) - remaining
				lggr.Infow("Confirmation progress",
					"finalized", finalized,
					"total", len(sigs),
					"pollCount", pollCount)
			}
		}
	}

	// Log successful completion
	lggr.Infow("Successfully funded all accounts",
		"accountCount", len(accounts),
		"amountSOL", solAmount)
	return nil
}
