package changeset

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

var _ cldf.ChangeSetV2[TransferNativeInput] = TransferNative{}

type TransferNativeInput struct {
	ChainSel uint64   `json:"chainSel" yaml:"chainSel"`
	Address  string   `json:"address" yaml:"address"`
	Amount   *big.Int `json:"amount" yaml:"amount"`
}

type TransferNative struct{}

func (t TransferNative) VerifyPreconditions(e cldf.Environment, config TransferNativeInput) error {
	if config.Address == "" {
		return errors.New("address cannot be empty")
	}
	valid := common.IsHexAddress(config.Address)
	if !valid {
		return fmt.Errorf("address string %s cannot be converted to ETH hex address", config.Address)
	}

	if config.Amount.Cmp(big.NewInt(0)) < 1 {
		return fmt.Errorf("amount must be positive value: %d", config.Amount)
	}

	chain, ok := e.BlockChains.EVMChains()[config.ChainSel]
	if !ok {
		return fmt.Errorf("chain not found for selector %d", config.ChainSel)
	}

	account := common.HexToAddress(config.Address)
	// Validate the transfer will succeed with simulation
	_, err := chain.Client.EstimateGas(e.GetContext(), ethereum.CallMsg{
		From:  chain.DeployerKey.From,
		To:    &account,
		Value: config.Amount,
	})
	if err != nil {
		return fmt.Errorf("transaction simulation failed: %w", err)
	}

	return nil
}

func (t TransferNative) Apply(e cldf.Environment, config TransferNativeInput) (cldf.ChangesetOutput, error) {
	chain, ok := e.BlockChains.EVMChains()[config.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", config.ChainSel)
	}

	e.Logger.Infow("Starting transfer of native funds", "chainSelector", config.ChainSel, "fromAddress", chain.DeployerKey.From, "toAddress", config.Address, "amount", config.Amount)

	_, err := operations.ExecuteSequence(
		e.OperationsBundle,
		TransferNativeSeq,
		TransferNativeDeps{
			Env: &e,
		},
		TransferNativeOpsInput{
			ChainSel: config.ChainSel,
			Address:  common.HexToAddress(config.Address),
			Amount:   config.Amount,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	e.Logger.Infow("Completed transfer of native funds", "chainSelector", config.ChainSel, "fromAddress", chain.DeployerKey.From, "toAddress", config.Address, "amount", config.Amount)

	return cldf.ChangesetOutput{}, nil
}

var TransferNativeSeq = operations.NewSequence(
	"transfer-native-seq",
	semver.MustParse("1.0.0"),
	"Sequence to transfer native funds from the deployer key to another address",
	func(b operations.Bundle, deps TransferNativeDeps, input TransferNativeOpsInput) (TransferNativeOutput, error) {
		_, err := operations.ExecuteOperation(
			b,
			TransferNativeOp,
			TransferNativeDeps{
				Env: deps.Env,
			},
			input,
		)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("failed to transfer funds: %w", err)
		}

		return TransferNativeOutput{}, nil
	},
)

type TransferNativeOpsInput struct {
	ChainSel uint64
	Address  common.Address
	Amount   *big.Int
}

type TransferNativeOutput struct{}

type TransferNativeDeps struct {
	Env *cldf.Environment
}

var TransferNativeOp = operations.NewOperation(
	"transfer-native-op",
	semver.MustParse("1.0.0"),
	"Operation to transfer funds from the deployer key to another address",
	func(b operations.Bundle, deps TransferNativeDeps, input TransferNativeOpsInput) (TransferNativeOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSel]
		if !ok {
			return TransferNativeOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSel)
		}

		nonce, err := chain.Client.NonceAt(b.GetContext(), chain.DeployerKey.From, nil)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("could not get latest nonce for deployer key: %w", err)
		}

		tipCap, err := chain.Client.SuggestGasTipCap(b.GetContext())
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("could not suggest gas tip cap: %w", err)
		}

		latestBlock, err := chain.Client.HeaderByNumber(b.GetContext(), nil)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("could not get latest block: %w", err)
		}
		baseFee := latestBlock.BaseFee

		feeCap := new(big.Int).Add(
			new(big.Int).Mul(baseFee, big.NewInt(2)),
			tipCap,
		)

		account := input.Address

		gasLimit, err := chain.Client.EstimateGas(b.GetContext(), ethereum.CallMsg{
			From:  chain.DeployerKey.From,
			To:    &account,
			Value: input.Amount,
		})
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("could not estimate gas for chain %d: %w", chain.Selector, err)
		}

		gasCost := new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), feeCap)
		gasPlusValue := new(big.Int).Add(gasCost, input.Amount)

		bal, err := chain.Client.BalanceAt(b.GetContext(), chain.DeployerKey.From, nil)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("could not get balance for deployer key: %w", err)
		}

		if bal.Cmp(gasPlusValue) < 0 {
			return TransferNativeOutput{}, fmt.Errorf("deployer key balance %d is insufficient to cover transfer amount %d plus max gas cost %d", bal, input.Amount, gasCost)
		}

		baseTx := &gethtypes.DynamicFeeTx{
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &account,
			Value:     input.Amount,
		}
		tx := gethtypes.NewTx(baseTx)

		signedTx, err := chain.DeployerKey.Signer(chain.DeployerKey.From, tx)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("could not sign transaction for account %s: %w", account.Hex(), err)
		}

		err = chain.Client.SendTransaction(b.GetContext(), signedTx)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("failed to send transfer to %s on chain %d: %w", account.Hex(), chain.Selector, err)
		}

		_, err = chain.Confirm(signedTx)
		if err != nil {
			return TransferNativeOutput{}, fmt.Errorf("failed to confirm transfer to %s on chain %d (tx %s): %w", account.Hex(), chain.Selector, signedTx.Hash().Hex(), err)
		}

		return TransferNativeOutput{}, nil
	},
)
