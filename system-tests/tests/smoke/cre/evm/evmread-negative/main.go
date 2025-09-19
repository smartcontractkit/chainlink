//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread-negative/config"
)

const balanceReaderABIJson = `[
   {
      "inputs":[
         {
            "internalType":"address[]",
            "name":"addresses",
            "type":"address[]"
         }
      ],
      "name":"getNativeBalances",
      "outputs":[
         {
            "internalType":"uint256[]",
            "name":"",
            "type":"uint256[]"
         }
      ],
      "stateMutability":"view",
      "type":"function"
   }
]`

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
	default:
		runtime.Logger().Warn("The provided name for function to execute did not match any known functions", "functionToTest", wfCfg.FunctionToTest)
	}
	return
}

func runCallContractForInvalidAddressesToRead(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (any, error) {
	readBalancesParsedABI, err := getReadBalancesContractABI(runtime, balanceReaderABIJson)
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("failed to get ReadBalances ABI: %v", err))
		return nil, fmt.Errorf("failed to get ReadBalances ABI: %w", err)
	}

	reply, err := readInvalidBalancesFromContract(readBalancesParsedABI, client, runtime, wfCfg)
	if err != nil {
		runtime.Logger().Error("callContract errored - invalid address to read", "address", wfCfg.InvalidInput, "error", err)
		return nil, fmt.Errorf("callContract errored - invalid address to read: %w", err)
	}
	return reply, nil
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

func getReadBalancesContractABI(runtime sdk.Runtime, balanceReaderABI string) (abi.ABI, error) {
	parsedABI, err := abi.JSON(strings.NewReader(balanceReaderABI))
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("failed to parse ABI: %v", err))
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %w", err)
	}
	runtime.Logger().With().Info(fmt.Sprintln("Parsed ABI successfully"))
	return parsedABI, nil
}

// readInvalidBalancesFromContract tries to read balances for an invalid address
// eventually it should return an empty array of balances
func readInvalidBalancesFromContract(readBalancesABI abi.ABI, evmClient evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.CallContractReply, error) {
	invalidAddressToRead := common.HexToAddress(wfCfg.InvalidInput)
	methodName := "getNativeBalances"
	packedData, err := readBalancesABI.Pack(methodName, []common.Address{invalidAddressToRead})
	if err != nil {
		return nil, fmt.Errorf("failed to pack read balances call: %w", err)
	}
	readBalancesOutput, err := evmClient.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   wfCfg.BalanceReader.BalanceReaderAddress.Bytes(),
			Data: packedData,
		},
	}).Await()
	if err != nil {
		runtime.Logger().Error("this is not expected: reading invalid balances should return 0", "address", invalidAddressToRead.String(), "error", err)
		return nil, fmt.Errorf("failed to get balances for address '%s': %w", invalidAddressToRead.String(), err)
	}

	var readBalancePrices []*big.Int
	err = readBalancesABI.UnpackIntoInterface(&readBalancePrices, methodName, readBalancesOutput.Data)
	if err != nil {
		runtime.Logger().Error("this is not expected: reading the CallContract output should return empty array", "address", invalidAddressToRead.String(), "error", err)
		return nil, fmt.Errorf("failed to read CallContract output: %w", err)
	}
	runtime.Logger().Info("Read on-chain balances", "address", invalidAddressToRead.String(), "balances", &readBalancePrices)
	return readBalancesOutput, nil
}
