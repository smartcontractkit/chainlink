package crib

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/sync/errgroup"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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

				err := memory.FundSolanaAccountsWithLogging(env.GetContext(), solanaAddrs, solFunds, chain.Client, lggr)
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
