package crib

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

var wg = &sync.WaitGroup{}

func distributeFunds(lggr logger.Logger, nodeInfo []devenv.Node, env deployment.Environment) {
	transmittersStr := make([]common.Address, 0)
	fundingAmount := new(big.Int).Mul(deployment.UBigInt(5), deployment.UBigInt(1e17)) // 0.5 ETH
	minThreshold := new(big.Int).Mul(deployment.UBigInt(5), deployment.UBigInt(1e16))  // 0.05 ETH

	for sel, chain := range env.Chains {
		wg.Add(1)
		go func(sel uint64, chain deployment.Chain) {
			defer wg.Done()
			for _, n := range nodeInfo {
				chainID, err := chainsel.ChainIdFromSelector(sel)
				if err != nil {
					lggr.Errorw("could not get chain id from selector", "selector", sel, "err", err)
					continue
				}
				addr := common.HexToAddress(n.AccountAddr[chainID])
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
			SendFundsToAccounts(lggr, chain, transmittersStr, fundingAmount, sel)
		}(sel, chain)
	}

	wg.Wait()
}

func SendFundsToAccounts(lggr logger.Logger, chain deployment.Chain, accounts []common.Address, fundingAmount *big.Int, sel uint64) {
	latesthdr, err := chain.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		lggr.Errorw("could not get header, skipping chain", "chain", sel, "err", err)
		return
	}
	block := latesthdr.Number

	nonce, err := chain.Client.NonceAt(context.Background(), chain.DeployerKey.From, block)
	if err != nil {
		lggr.Warnw("could not get latest nonce for deployer key", "err", err)
		return
	}
	for _, address := range accounts {
		tx := gethtypes.NewTransaction(nonce, address, fundingAmount, uint64(1000000), big.NewInt(1000000), nil)

		signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, tx)
		if err != nil {
			lggr.Errorw("could not sign transaction for sending funds to ", "chain", sel, "account", address, "err", err)
			continue
		}

		lggr.Infow("sending transaction for ", "account", address.String(), "chain", sel)
		err = chain.Client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			lggr.Errorw("could not send transaction to address on ", "chain", sel, "address", address, "err", err)
			continue
		}

		_, err = bind.WaitMined(context.Background(), chain.Client, signedTx)
		if err != nil {
			lggr.Errorw("could not mine transaction to address on ", "chain", sel)
			continue
		}
		nonce++
	}
}
