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

	tipCap, err := evmChain.Client.SuggestGasTipCap(ctx)
	if err != nil {
		return out, fmt.Errorf("could not suggest gas tip cap: %w", err)
	}

	latestBlock, err := evmChain.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return out, fmt.Errorf("could not get latest block: %w", err)
	}
	baseFee := latestBlock.BaseFee

	feeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		tipCap,
	)

	nonce, err := evmChain.Client.PendingNonceAt(ctx, evmChain.DeployerKey.From)
	if err != nil {
		return out, fmt.Errorf("could not get pending nonce: %v", err)
	}

	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       200000,
		To:        &feeManagerAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})

	signedTx, err := evmChain.DeployerKey.Signer(evmChain.DeployerKey.From, tx)
	if err != nil {
		return out, fmt.Errorf("could not sign transaction")
	}

	err = evmChain.Client.SendTransaction(context.Background(), signedTx)

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
