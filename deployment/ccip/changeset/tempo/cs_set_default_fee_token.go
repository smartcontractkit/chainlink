package hyperliquid

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var EnableBigBlockChangeset = cldf.CreateChangeSet(setFeeTokenLogic, setFeeTokenPreCondition)

type SetFeeTokenConfig struct {
	ChainSel          uint64
	FeeTokenAddress   string
	FeeManagerAddress string
}

func setFeeTokenPreCondition(env cldf.Environment, cfg SetFeeTokenConfig) error {
	_, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	return nil
}

func setFeeTokenLogic(env cldf.Environment, cfg SetFeeTokenConfig) (cldf.ChangesetOutput, error) {
	out := cldf.ChangesetOutput{}
	ctx := context.Background()

	tokenAddress := common.HexToAddress(cfg.FeeTokenAddress)
	feeManagerAddress := common.HexToAddress(cfg.FeeManagerAddress)

	chain, err := findChainBySelector(env, cfg.ChainSel)
	if err != nil {
		return out, fmt.Errorf("error: %w finding chain by selector: %d", err, cfg.ChainSel)
	}

	evmChain, ok := chain.(evm.Chain)
	if !ok {
		return out, fmt.Errorf("not an EVM chain")
	}

	methodSig := []byte("setUserToken(address)")
	selector := crypto.Keccak256(methodSig)[:4]

	addressType, _ := abi.NewType("address", "", nil)
	arguments := abi.Arguments{{Type: addressType}}
	encodedArgs, err := arguments.Pack(tokenAddress)
	if err != nil {
		return out, fmt.Errorf("abi pack: %v", err)
	}

	data := append(selector, encodedArgs...)

	gas, err := evmChain.Client.SuggestGasPrice(ctx)
	if err != nil {
		return out, fmt.Errorf("could not estimate gas: %w", err)
	}

	nonce, err := evmChain.Client.PendingNonceAt(ctx, evmChain.DeployerKey.From)
	if err != nil {
		return out, fmt.Errorf("could not get pending nonce: %v", err)
	}

	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    nonce,
		To:       &feeManagerAddress,
		Value:    big.NewInt(0),
		Gas:      200000,
		GasPrice: gas,
		Data:     data,
	})

	err = evmChain.Client.SendTransaction(context.Background(), tx)

	if err != nil {
		log.Fatalf("failed to send tx: %v", err)
	}

	return out, nil
}

func findChainBySelector(e cldf.Environment, selector uint64) (chain.BlockChain, error) {
	evmChains := e.BlockChains.EVMChains()

	for _, chain := range evmChains {
		if chain.ChainSelector() == selector {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("error finding chain with selector: %d", selector)
}
