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
	case "EstimateGas - invalid 'to' address":
		// it does not make sense to test with invalid CallMsg.Data because any bytes will be correctly processed
		return runEstimateGasForInvalidToAddress(client, runtime, wfCfg)
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

// readInvalidBalancesFromContract tries to read balances for an invalid address
// eventually it should return an empty array of balances
func runCallContractForInvalidAddressesToRead(evmClient evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.CallContractReply, error) {
	readBalancesABI, _ := getReadBalanceAbi(runtime)
	invalidAddressToRead := wfCfg.InvalidInput
	methodName := "getNativeBalances"
	readBalancesCallWithInvalidAddressToRead, _ := getPackedReadBalancesCall(methodName, invalidAddressToRead, readBalancesABI)

	runtime.Logger().Info("Attempting to read balances using invalid address to read", "invalid_address", invalidAddressToRead)
	validReadBalancesAddress := wfCfg.BalanceReader.BalanceReaderAddress
	readBalancesOutput, err := evmClient.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   validReadBalancesAddress.Bytes(),
			Data: readBalancesCallWithInvalidAddressToRead,
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

	// this line produces the expected 0 balances result: balances=&[+0]
	runtime.Logger().Info("got expected 0 balances for invalid addresses to read with CallContract", "invalid_address", invalidAddressToRead, "balances", &readBalancePrices)
	return readBalancesOutput, nil
}

// runCallContractForInvalidContractAddress is referring to invalid contract address
// evm capability should return an error
func runCallContractForInvalidContractAddress(evmClient evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.CallContractReply, error) {
	// it is a valid 0-address to read, it may be hardcoded
	// it should not make CallContract to error.
	// Instead, it returns either 0 or some balance depending on a chain used.
	addressToRead := "0x0000000000000000000000000000000000000000"
	methodName := "getNativeBalances"
	readBalancesABI, _ := getReadBalanceAbi(runtime)
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
	runtime.Logger().Info("CallContract for invalid balance reader contract address completed", "output_data", readBalancesOutput.Data)
	if err != nil || len(readBalancesOutput.Data) == 0 {
		runtime.Logger().Error("got expected error for invalid balance reader contract address", "invalid_rb_address", invalidReadBalancesContractAddr.String(), "error", err, "output_data", readBalancesOutput.Data)
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

func runEstimateGasForInvalidToAddress(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (any, error) {
	runtime.Logger().Info("Attempting to EstimateGas using invalid 'to' address", "invalid_to_address", wfCfg.InvalidInput)
	marshalledTx := common.FromHex("02f8f18205392084481f228084481f228782608294c3e53f4d16ae77db1c982e75a937b9f60fe6369080b8842ac0df2600000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000026496e697469616c206d65737361676520746f206265207265616420627920776f726b666c6f770000000000000000000000000000000000000000000000000000c080a008a98a170eeeca4d94df4bae10e61b5fc7d0084313cf42761dfc361f23e86d74a02144720570a62b17bbb774a3f083ced13251d1eb9d7f85101ee9d4410479ead9")

	invalidToAddress := common.Address(common.HexToAddress(wfCfg.InvalidInput))
	estimatedGasReply, err := client.EstimateGas(runtime, &evm.EstimateGasRequest{
		Msg: &evm.CallMsg{
			To:   invalidToAddress.Bytes(),
			Data: marshalledTx,
		},
	}).Await()
	runtime.Logger().Info("EstimateGas completed", "output_data", estimatedGasReply)
	if err != nil || estimatedGasReply == nil {
		runtime.Logger().Error("got expected error for GasEstimate invalid 'to' address", "invalid_to_address", invalidToAddress.String(), "error", err, "output_data", estimatedGasReply)
		return nil, fmt.Errorf("expected error for GasEstimate invalid 'to' address '%s': %w", invalidToAddress.String(), err)
	}

	runtime.Logger().Info("this is not expected: GasEstimate for invalid 'to' address should return an error or empty response", "invalid_to_address", invalidToAddress.String(), "output_data", estimatedGasReply)
	return estimatedGasReply, nil
}
