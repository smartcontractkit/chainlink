package crib

import (
	"context"
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
	latesthdr, err := chain.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		lggr.Errorw("could not get header, skipping chain", "chain", sel, "err", err)
		return err
	}
	block := latesthdr.Number

	nonce, err := chain.Client.NonceAt(context.Background(), chain.DeployerKey.From, block)
	if err != nil {
		lggr.Warnw("could not get latest nonce for deployer key", "err", err)
		return err
	}
	for _, address := range accounts {
		gasPrice, err := chain.Client.SuggestGasPrice(ctx)
		if err != nil {
			lggr.Warnw("could not suggest gas price, using default", "err", err)
			gasPrice = big.NewInt(100000000) // 1 Gwei as default
		}

		gasLimit, err := chain.Client.EstimateGas(ctx, ethereum.CallMsg{
			From:     chain.DeployerKey.From,
			To:       &address,
			Value:    fundingAmount,
			GasPrice: gasPrice,
		})
		if err != nil {
			lggr.Warnw("could not estimate gas, using default", "err", err)
			gasLimit = uint64(1000000)
		}

		tx := gethtypes.NewTransaction(nonce, address, fundingAmount, gasLimit, gasPrice, nil)

		signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, tx)
		if err != nil {
			lggr.Errorw("could not sign transaction for sending funds to ", "chain", sel, "account", address, "err", err)
			return err
		}

		lggr.Infow("sending transaction for ", "account", address.String(), "chain", sel)
		err = chain.Client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			lggr.Errorw("could not send transaction to address on ", "chain", sel, "address", address, "err", err)
			return err
		}

		_, err = bind.WaitMined(context.Background(), chain.Client, signedTx)
		if err != nil {
			lggr.Errorw("could not mine transaction to address on ", "chain", sel)
			return err
		}
		nonce++
	}
	return nil
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
