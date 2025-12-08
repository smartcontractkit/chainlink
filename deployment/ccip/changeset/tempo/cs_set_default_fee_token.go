package tempo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var SetFeeTokenChangeset = cldf.CreateLegacyChangeSet(setFeeTokenLogic)

type SetFeeTokenConfig struct {
	ChainSel          uint64
	FeeTokenAddress   common.Address
	FeeManagerAddress common.Address
}

func setFeeTokenLogic(env cldf.Environment, cfg SetFeeTokenConfig) (cldf.ChangesetOutput, error) {
	out := cldf.ChangesetOutput{}
	ctx := context.Background()

	evmChain := env.BlockChains.EVMChains()[cfg.ChainSel]

	methodSig := []byte("setUserToken(address)")
	data := crypto.Keccak256(methodSig)[:4]

	addressType, _ := abi.NewType("address", "", nil)
	arguments := abi.Arguments{{Type: addressType}}
	encodedArgs, err := arguments.Pack(cfg.FeeTokenAddress)
	if err != nil {
		return out, fmt.Errorf("abi pack: %w", err)
	}

	data = append(data, encodedArgs...)

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
		return out, fmt.Errorf("could not get pending nonce: %w", err)
	}

	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       200000,
		To:        &cfg.FeeManagerAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})

	signedTx, err := evmChain.DeployerKey.Signer(evmChain.DeployerKey.From, tx)
	if err != nil {
		return out, errors.New("could not sign transaction")
	}

	err = evmChain.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("failed to send tx: %v", err)
	}

	return out, nil
}
