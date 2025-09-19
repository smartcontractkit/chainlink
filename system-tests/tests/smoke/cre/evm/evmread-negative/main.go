//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread-negative/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		wfCfg := config.Config{}
		if err := yaml.Unmarshal(b, &wfCfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return wfCfg, nil
	}).Run(RunReadWorkflow)
}

func RunReadWorkflow(wfCfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onEVMReadTrigger,
		),
	}, nil
}

func onEVMReadTrigger(wfCfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, _ error) {
	runtime.Logger().Info("onEVMReadFailTrigger called", "payload", payload)

	client := evm.Client{ChainSelector: wfCfg.ChainSelector}

	switch wfCfg.FunctionToTest {
	case "BalanceAt":
		_, err := client.BalanceAt(runtime, &evm.BalanceAtRequest{
			Account:     []byte(wfCfg.InvalidInput),
			BlockNumber: nil,
		}).Await()
		if err != nil {
			runtime.Logger().Error("balanceAt errored", "error", err)
			return nil, fmt.Errorf("balanceAt errored: %w", err)
		}
		return
	// case "LatestBlockNumber":
	// 	latestHeadNumber := requireLatestBlockNumber(t, runtime, client)
	// 	runtime.Logger().Info("Successfully got latestHeadNumber", "blockNumber", latestHeadNumber)
	// 	return
	// case "requireEvent":
	// 	latestHeadNumber := requireLatestBlockNumber(t, runtime, client)
	// 	requireEvent(t, wfCfg, runtime, latestHeadNumber, client)
	// 	runtime.Logger().Info("Successfully got event")
	// 	return
	// case "requireContractCall":
	// 	requireContractCall(t, wfCfg, runtime, client)
	// 	runtime.Logger().Info("Successfully called contract")
	// 	return
	// case "requireReceipt":
	// 	requireReceipt(t, runtime, wfCfg, client)
	// 	runtime.Logger().Info("Successfully got receipt")
	// 	return
	// case "requireTx":
	// 	var expectedTx types.Transaction
	// 	err := expectedTx.UnmarshalBinary(wfCfg.ExpectedBinaryTx)
	// 	require.NoError(t, err)
	// 	requireTx(t, runtime, &expectedTx, client)
	// 	runtime.Logger().Info("Successfully got transaction")
	// 	return
	// case "requireEstimatedGas":
	// 	var expectedTx types.Transaction
	// 	err := expectedTx.UnmarshalBinary(wfCfg.ExpectedBinaryTx)
	// 	require.NoError(t, err)
	// 	requireEstimatedGas(t, runtime, wfCfg, expectedTx.Data(), client)
	// 	runtime.Logger().Info("Successfully estimated gas")
	// 	return
	// case "requireError":
	// 	requireError(t, runtime, wfCfg, client)
	// 	runtime.Logger().Info("Successfully got error for non-existing transaction")
	// 	return
	// case "sendTx":
	// 	txHash := sendTx(t, runtime, wfCfg, client, "EVM read workflow executed successfully")
	// 	runtime.Logger().Info("Successfully sent transaction", "hash", common.Hash(txHash).String())
	// 	return
	default:
		runtime.Logger().Warn("The provided name for function to execute did not match any known functions", "functionToTest", wfCfg.FunctionToTest)
	}
	return
}
