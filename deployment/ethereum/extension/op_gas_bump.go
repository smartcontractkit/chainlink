package deployment_ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
)

type gasConfig struct {
	GasLimit uint64
	GasPrice *big.Int
}

type GasBump struct {
	RetryLimit      uint
	RetryIntervalMs uint
	BumpPercentage  uint
}

type GasBumpOpInput struct {
	GasBump
	Tx struct {
		To    *common.Address
		Data  []byte
		Value *big.Int
	}
}

var SendTxWithGasBumpOp = deployment.NewOperation(
	"v1",
	"sends a raw transaction, bumping gas if tx fails for gas conditions",
	func(ctx deployment.OpContext, deps EthereumDeps, input GasBumpOpInput) (EthereumTxOutput, error) {

		bumpGas := func(gas gasConfig, attempt uint) gasConfig {
			bumpFactor := 100 + input.BumpPercentage*attempt
			return gasConfig{
				GasPrice: new(big.Int).Mul(gas.GasPrice, big.NewInt(int64(bumpFactor))),
				GasLimit: gas.GasLimit,
			}
		}

		buildRawTx := func(gas gasConfig) *types.Transaction {
			return types.NewTx(&types.LegacyTx{
				To:       input.Tx.To,
				Data:     input.Tx.Data,
				Value:    input.Tx.Value,
				Gas:      gas.GasLimit,
				GasPrice: gas.GasPrice,
			})
		}

		var receipt *types.Receipt
		var err error
		gas := gasConfig{
			GasLimit: 210000,
			// TODO: This should come estimated by the client
			GasPrice: big.NewInt(1000000000),
		}
		// Make a similar implementation to attemptSend but without recursion
		for attempt := uint(1); attempt <= input.RetryLimit; attempt++ {
			err := error(nil)
			// Sleep for the retry interval
			if attempt > 1 {
				ctx.Log.Info("Sleeping before retry", " attempt ", attempt, " interval ", input.RetryIntervalMs)
				time.Sleep(time.Duration(input.RetryIntervalMs) * time.Millisecond)
			}

			rawTx := buildRawTx(gas)

			tx, err := deps.Auth.Signer(deps.Auth.From, rawTx)
			if err != nil {
				// Error with the signature
				break
			}

			err = deps.Client.SendTransaction(context.Background(), tx)
			if err != nil {
				// TODO: Check if gas error
				gas = bumpGas(gas, attempt)
				continue
			}

			// Wait for the transaction to be mined
			receipt, err = deps.Confirm(deps.Client, tx.Hash())
			if err != nil {
				// TODO: Check if gas error
				gas = bumpGas(gas, attempt)
				continue
			} else {
				break
			}
		}

		if err != nil {
			return EthereumTxOutput{}, err
		}

		return EthereumTxOutput{Hash: receipt.TxHash, RawReceipt: receipt}, nil
	},
)
