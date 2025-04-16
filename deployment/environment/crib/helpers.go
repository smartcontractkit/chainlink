package crib

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"

	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

func distributeTransmitterFunds(lggr logger.Logger, nodeInfo []devenv.Node, env deployment.Environment) error {
	transmittersStr := make([]common.Address, 0)
	fundingAmount := new(big.Int).Mul(deployment.UBigInt(100), deployment.UBigInt(1e18)) // 100 ETH

	g := new(errgroup.Group)
	for sel, chain := range env.Chains {
		sel, chain := sel, chain
		g.Go(func() error {
			for _, n := range nodeInfo {
				chainID, err := chainsel.GetChainIDFromSelector(sel)
				if err != nil {
					lggr.Errorw("could not get chain id from selector", "selector", sel, "err", err)
					return err
				}
				addr := common.HexToAddress(n.AccountAddr[chainID])
				transmittersStr = append(transmittersStr, addr)
			}
			return SendFundsToAccounts(env.GetContext(), lggr, chain, transmittersStr, fundingAmount, sel)
		})
	}

	return g.Wait()
}

func SendFundsToAccounts(ctx context.Context, lggr logger.Logger, chain deployment.Chain, accounts []common.Address, fundingAmount *big.Int, sel uint64) error {
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
		tx := gethtypes.NewTransaction(nonce, address, fundingAmount, uint64(1000000), big.NewInt(100000000), nil)

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
