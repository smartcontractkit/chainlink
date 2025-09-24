//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/balance_reader"
	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmread-negative/config"
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
		return runBalanceAt(client, runtime, wfCfg)
	case "CallContract - invalid address to read":
		// it does not error, but returns empty array of balances
		return runCallContractForInvalidAddressesToRead(client, runtime, wfCfg)
	case "CallContract - invalid balance reader contract address":
		return runCallContractForInvalidContractAddress(client, runtime, wfCfg)
	default:
		runtime.Logger().Warn("The provided name for function to execute did not match any known functions", "functionToTest", wfCfg.FunctionToTest)
	}
	return
}

func runBalanceAt(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (_ any, _ error) {
	_, err := client.BalanceAt(runtime, &evm.BalanceAtRequest{
		Account:     []byte(wfCfg.InvalidInput),
		BlockNumber: nil,
	}).Await()
	if err != nil {
		runtime.Logger().Error("balanceAt errored", "error", err)
		return nil, fmt.Errorf("balanceAt errored: %w", err)
	}
	return
}

func runCallContractForInvalidContractAddress(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (any, error) {
	reply, err := readWithInvalidReaderContractAddress(client, runtime, wfCfg)
	if err != nil {
		runtime.Logger().Error("callContract errored - invalid contract address", "address", wfCfg.InvalidInput, "error", err)
		return nil, fmt.Errorf("callContract errored - invalid contract address: %w", err)
	}
	return reply, nil
}

func runCallContractForInvalidAddressesToRead(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (any, error) {
	reply, err := readInvalidBalancesFromContract(client, runtime, wfCfg)
	if err != nil {
		runtime.Logger().Error("callContract errored - invalid address to read", "address", wfCfg.InvalidInput, "error", err)
		return nil, fmt.Errorf("callContract errored - invalid address to read: %w", err)
	}
	return reply, nil
}

// readInvalidBalancesFromContract tries to read balances for an invalid address
// eventually it should return an empty array of balances
func readInvalidBalancesFromContract(evmClient evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.CallContractReply, error) {
	readBalancesABI, _ := getReadBalanceAbi(runtime)
	invalidAddressToRead := wfCfg.InvalidInput
	methodName := "getNativeBalances"
	readBalancesCall, _ := getPackedReadBalancesCall(methodName, invalidAddressToRead, readBalancesABI)

	runtime.Logger().Info("Attempting to read balances using invalid address to read", "invalid_address", invalidAddressToRead)
	readBalancesAddress := wfCfg.BalanceReader.BalanceReaderAddress
	readBalancesOutput, err := evmClient.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   readBalancesAddress.Bytes(),
			Data: readBalancesCall,
		},
	}).Await()
	if err != nil {
		runtime.Logger().Error("this is not expected: reading invalid balances should return 0", "invalid_address", invalidAddressToRead, "error", err)
		return nil, fmt.Errorf("failed to get balances for address '%s': %w", invalidAddressToRead, err)
	}

	var readBalancePrices []*big.Int
	err = readBalancesABI.UnpackIntoInterface(&readBalancePrices, methodName, readBalancesOutput.Data)
	if err != nil {
		runtime.Logger().Error("this is not expected: reading the CallContract output should return empty array", "invalid_address", invalidAddressToRead, "error", err)
		return nil, fmt.Errorf("failed to read CallContract output: %w", err)
	}
	runtime.Logger().Info("Read on-chain balances", "invalid_address", invalidAddressToRead, "balances", &readBalancePrices)
	return readBalancesOutput, nil
}

// readWithInvalidReaderContractAddress is referring to invalid contract address
// evm capability should return an error
func readWithInvalidReaderContractAddress(evmClient evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.CallContractReply, error) {
	readBalancesABI, _ := getReadBalanceAbi(runtime)
	// it is a valid 0-address to read,
	// it should not make CallContract to error.
	// Instead, it returns either 0 or some balance depending on a chain used.
	addressToRead := "0x0000000000000000000000000000000000000000"
	methodName := "getNativeBalances"
	readBalancesCall, _ := getPackedReadBalancesCall(methodName, addressToRead, readBalancesABI)

	runtime.Logger().Info("Attempting to read balances using invalid balance reader contract address", "invalid_rb_address", wfCfg.InvalidInput)
	invalidReadBalancesContractAddr := common.Address(common.HexToAddress(wfCfg.InvalidInput))
	runtime.Logger().Info("Starting CallContract request with parsed address", "invalid_rb_address", invalidReadBalancesContractAddr.String())
	readBalancesOutput, err := evmClient.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   invalidReadBalancesContractAddr.Bytes(),
			Data: readBalancesCall,
		},
	}).Await()
	runtime.Logger().Info("CallContract completed", "output_data", readBalancesOutput.Data)
	if err != nil || len(readBalancesOutput.Data) == 0 {
		runtime.Logger().Error("expected error for invalid balance reader contract address", "invalid_rb_address", invalidReadBalancesContractAddr.String(), "error", err, "output_data", readBalancesOutput.Data)
		return nil, fmt.Errorf("failed to get balances for address '%s': %w", invalidReadBalancesContractAddr.String(), err)
	}

	runtime.Logger().Info("this is not expected: reading from invalid balance reader contract address should return an error or empty response", "invalid_rb_address", invalidReadBalancesContractAddr.String(), "output", readBalancesOutput.Data)
	return readBalancesOutput, nil
}

func getPackedReadBalancesCall(methodName, addressToRead string, readBalancesABI *abi.ABI) ([]byte, error) {
	packedData, err := readBalancesABI.Pack(methodName, []common.Address{common.HexToAddress(addressToRead)})
	if err != nil {
		return nil, fmt.Errorf("failed to pack Read Balances call: %w", err)
	}
	return packedData, nil
}

func getReadBalanceAbi(runtime sdk.Runtime) (*abi.ABI, error) {
	readBalancesABI, abiErr := balance_reader.BalanceReaderMetaData.GetAbi()
	if abiErr != nil {
		runtime.Logger().Error("failed to get Balance Reader contract ABI", "error", abiErr)
		return nil, fmt.Errorf("failed to get Balance Reader contract ABI: %w", abiErr)
	}
	return readBalancesABI, nil
}
