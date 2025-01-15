package crib

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func distributeFunds(lggr logger.Logger, nodeInfo []devenv.Node, env deployment.Environment) {
	transmittersStr := make([]common.Address, 0)
	fundingAmount := big.NewInt(500_000_000_000_000_000) // 0.5 ETH
	minThreshold := big.NewInt(50_000_000_000_000_000)   // 0.05 ETH

	// todo: parallelize this using a waitgroup
	for sel, chain := range env.Chains {
		for _, n := range nodeInfo {

			chainId, err := chainsel.ChainIdFromSelector(sel)
			if err != nil {
				lggr.Errorw("could not get chain id from selector", "selector", sel, "err", err)
				continue
			}
			addr := common.HexToAddress(n.AccountAddr[chainId])
			balance, err := chain.Client.BalanceAt(context.Background(), addr, nil)
			if err != nil {
				lggr.Errorw("error fetching balance for %s: %v\n", n.Name, err)
				continue
			}
			if balance.Cmp(minThreshold) < 0 {
				lggr.Infow(
					"sending funds to",
					"node", n.Name,
					"address", addr.String(),
					"amount", conversions.WeiToEther(fundingAmount).String(),
				)
				transmittersStr = append(transmittersStr, addr)
			}
		}
		latesthdr, err := chain.Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			lggr.Errorw("could not get header, skipping chain", "chain", sel, "err", err)
			continue
		}
		block := latesthdr.Number.Uint64()

		nonce, err := chain.Client.NonceAt(context.Background(), chain.DeployerKey.From, big.NewInt(int64(block)))
		if err != nil {
			lggr.Warnw("could not get latest nonce for deployer key", "err", err)
			continue
		}
		for _, transmitter := range transmittersStr {
			tx := gethtypes.NewTransaction(nonce, transmitter, fundingAmount, uint64(1000000), big.NewInt(1000000), nil)

			signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, tx)
			if err != nil {
				lggr.Errorw("could not sign transaction to transmitter on ", "chain", sel, "transmitter", transmitter, "err", err)
				continue
			}

			lggr.Infow("sending transaction for ", "transmitter", transmitter.String(), "chain", sel)
			err = chain.Client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				lggr.Errorw("could not send transaction to transmitter on ", "chain", sel, "transmitter", transmitter, "err", err)
				continue
			}

			_, err = bind.WaitMined(context.Background(), chain.Client, signedTx)
			if err != nil {
				lggr.Errorw("could not mine transaction to transmitter on ", "chain", sel)
				continue
			}
			nonce++
		}
	}
}
