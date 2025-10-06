//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
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
	case "FilterLogs - invalid addresses":
		return runFilterLogsWithInvalidAddresses(client, runtime, wfCfg)
	case "FilterLogs - invalid FromBlock":
		return runFilterLogsWithInvalidFromBlock(client, runtime, wfCfg)
	case "FilterLogs - invalid ToBlock":
		return runFilterLogsWithInvalidToBlock(client, runtime, wfCfg)
	case "GetTransactionByHash - invalid hash":
		return runGetTransactionByHashWithInvalidHash(client, runtime, wfCfg)
	case "GetTransactionReceipt - invalid hash":
		return runGetTransactionReceiptWithInvalidHash(client, runtime, wfCfg)
	case "HeaderByNumber - invalid block number":
		return runHeaderByNumberWithInvalidBlock(client, runtime, wfCfg)
	default:
		runtime.Logger().Warn("The provided name for function to test in regression EVM Read Workflow did not match any known functions", "functionToTest", wfCfg.FunctionToTest)
		return nil, fmt.Errorf("the provided name for function to test in regression EVM Read Workflow did not match any known functions: %s", wfCfg.FunctionToTest)
	}
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
	runtime.Logger().Info("CallContract balance reading completed", "output_data", readBalancesOutput.Data)
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
	runtime.Logger().Info("CallContract for invalid balance reader contract address completed", "balance_reader_output", readBalancesOutput)
	if err != nil || readBalancesOutput == nil {
		runtime.Logger().Error("got expected error for invalid balance reader contract address", "invalid_rb_address", invalidReadBalancesContractAddr.String(), "balance_reader_output", readBalancesOutput, "error", err)
		return nil, fmt.Errorf("failed to get balances for address '%s': %w", invalidReadBalancesContractAddr.String(), err)
	} else if len(readBalancesOutput.Data) == 0 {
		runtime.Logger().Error("got expected empty response for invalid balance reader contract address", "invalid_rb_address", invalidReadBalancesContractAddr.String(), "balance_reader_output", readBalancesOutput, "error", err)
		return nil, fmt.Errorf("failed to get balances for address '%s': %w", invalidReadBalancesContractAddr.String(), err)
	}

	runtime.Logger().Info("this is not expected: reading from invalid balance reader contract address should return an error or empty response", "invalid_rb_address", invalidReadBalancesContractAddr.String(), "balance_reader_output", readBalancesOutput)
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
	runtime.Logger().Info("getting Balance Reader contract ABI")
	readBalancesABI, abiErr := balance_reader.BalanceReaderMetaData.GetAbi()
	if abiErr != nil {
		runtime.Logger().Error("failed to get Balance Reader contract ABI", "error", abiErr)
		return nil, fmt.Errorf("failed to get Balance Reader contract ABI: %w", abiErr)
	}
	runtime.Logger().Info("successfully got Balance Reader contract ABI")
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

// runFilterLogsWithInvalidAddresses tries to filter logs using invalid addresses in the request
// it should return an error or empty logs
func runFilterLogsWithInvalidAddresses(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.FilterLogsReply, error) {
	invalidAddress := common.HexToAddress(wfCfg.InvalidInput)
	runtime.Logger().Info("Attempting to filter logs using invalid addresses", "invalid_address", invalidAddress)

	filterLogsOutput, err := client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{invalidAddress.Bytes()},
			FromBlock: pb.NewBigIntFromInt(big.NewInt(100)), // 100 blocks is a max valid range between blocks
			ToBlock:   pb.NewBigIntFromInt(big.NewInt(200)),
		},
	}).Await()
	runtime.Logger().Info("FilterLogs completed", "filtered_logs_output", filterLogsOutput)
	if err != nil || len(filterLogsOutput.Logs) == 0 {
		runtime.Logger().Error("got expected error or empty logs for FilterLogs with invalid addresses", "invalid_address", invalidAddress, "filter_logs_output", filterLogsOutput.Logs, "error", err)
		return filterLogsOutput, fmt.Errorf("expected error or empty logs for FilterLogs with invalid address '%s': %w", invalidAddress, err)
	}

	runtime.Logger().Info("this is not expected: FilterLogs with invalid addresses in the request should return an error or empty logs", "invalid_address", invalidAddress, "filter_logs_output", filterLogsOutput.Logs)
	return filterLogsOutput, nil
}

// runFilterLogsWithInvalidFromBlock tries to filter logs using invalid fromBlock values
// it should return an error
func runFilterLogsWithInvalidFromBlock(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.FilterLogsReply, error) {
	return runFilterLogsWithInvalidBlock(client, runtime, wfCfg, "fromBlock")
}

// runFilterLogsWithInvalidToBlock tries to filter logs using invalid toBlock values
// it should return an error
func runFilterLogsWithInvalidToBlock(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.FilterLogsReply, error) {
	return runFilterLogsWithInvalidBlock(client, runtime, wfCfg, "toBlock")
}

// runFilterLogsWithInvalidBlock tries to filter logs using invalid block values
// it should return an error for invalid block values
func runFilterLogsWithInvalidBlock(client evm.Client, runtime sdk.Runtime, wfCfg config.Config, blockType string) (*evm.FilterLogsReply, error) {
	invalidBlockStr := wfCfg.InvalidInput
	runtime.Logger().Info("Attempting to filter logs using invalid block", "block_type", blockType, "invalid_block", invalidBlockStr)

	// Parse the invalid block string to big.Int
	newBlock := big.NewInt(0)
	invalidBlock, _ := newBlock.SetString(invalidBlockStr, 10)

	// A valid address for FilterLogs
	validAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	// Set up the filter query based on which block type is being tested
	var filterQuery *evm.FilterQuery
	if blockType == "fromBlock" {
		filterQuery = &evm.FilterQuery{
			Addresses: [][]byte{validAddress.Bytes()},
			FromBlock: pb.NewBigIntFromInt(invalidBlock),
			ToBlock:   pb.NewBigIntFromInt(big.NewInt(150)),
		}
	} else { // toBlock
		filterQuery = &evm.FilterQuery{
			Addresses: [][]byte{validAddress.Bytes()},
			FromBlock: pb.NewBigIntFromInt(big.NewInt(2)),
			ToBlock:   pb.NewBigIntFromInt(invalidBlock),
		}
	}

	filterLogsOutput, err := client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: filterQuery,
	}).Await()
	runtime.Logger().Info("FilterLogs with invalid block completed", "block_type", blockType, "filtered_logs_output", filterLogsOutput)
	if err != nil || filterLogsOutput == nil {
		runtime.Logger().Error("got expected error for FilterLogs with invalid block", "block_type", blockType, "invalid_block", invalidBlockStr, "filter_logs_output", filterLogsOutput, "error", err)
		return filterLogsOutput, fmt.Errorf("expected error for FilterLogs with invalid %s '%s': %w", blockType, invalidBlockStr, err)
	}

	runtime.Logger().Info("this is not expected: FilterLogs with invalid block should return an error or nil", "block_type", blockType, "invalid_block", invalidBlockStr, "filter_logs_output", filterLogsOutput)
	return filterLogsOutput, nil
}

// runGetTransactionByHashWithInvalidHash tries to get a transaction using an invalid hash
func runGetTransactionByHashWithInvalidHash(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.GetTransactionByHashReply, error) {
	runtime.Logger().Info("Attempting to get transaction using invalid hash", "invalid_hash", wfCfg.InvalidInput)

	invalidHash := common.FromHex(wfCfg.InvalidInput)
	runtime.Logger().Info("Starting GetTransactionByHash request with parsed hash", "invalid_hash", invalidHash)
	txByHashOutput, err := client.GetTransactionByHash(runtime, &evm.GetTransactionByHashRequest{
		Hash: invalidHash,
	}).Await()
	runtime.Logger().Info("GetTransactionByHash completed", "tx_by_hash_output", txByHashOutput)
	if err != nil || txByHashOutput == nil {
		runtime.Logger().Error("got expected error for GetTransactionByHash with invalid hash", "invalid_hash", invalidHash, "tx_by_hash_output", txByHashOutput, "error", err)
		return nil, fmt.Errorf("expected error for GetTransactionByHash with invalid hash '%s': %w", invalidHash, err)
	}

	runtime.Logger().Info("this is not expected: GetTransactionByHash with invalid hash should return an error or nil", "invalid_hash", invalidHash, "tx_by_hash_output", txByHashOutput)
	return txByHashOutput, nil
}

// runGetTransactionReceiptWithInvalidHash tries to get transaction receipt using an invalid hash
// it should return an error
func runGetTransactionReceiptWithInvalidHash(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.GetTransactionReceiptReply, error) {
	runtime.Logger().Info("Attempting to GetTransactionReceipt using invalid hash", "invalid_hash", wfCfg.InvalidInput)

	// Convert the invalid input to bytes - this will handle various invalid formats
	invalidHash := common.FromHex(wfCfg.InvalidInput)
	runtime.Logger().Info("Starting GetTransactionReceipt request with parsed hash", "invalid_hash", invalidHash)
	txReceiptOutput, err := client.GetTransactionReceipt(runtime, &evm.GetTransactionReceiptRequest{
		Hash: invalidHash,
	}).Await()
	runtime.Logger().Info("GetTransactionReceipt completed", "tx_receipt_output", txReceiptOutput)
	if err != nil || txReceiptOutput == nil {
		runtime.Logger().Error("got expected error for GetTransactionReceipt with invalid hash", "invalid_hash", invalidHash, "tx_receipt_output", txReceiptOutput, "error", err)
		return nil, fmt.Errorf("expected error for GetTransactionReceipt with invalid hash '%s': %w", invalidHash, err)
	}

	runtime.Logger().Info("this is not expected: GetTransactionReceipt with invalid hash should return an error or nil", "invalid_hash", invalidHash, "tx_receipt_output", txReceiptOutput)
	return txReceiptOutput, nil
}

// runHeaderByNumberWithInvalidBlock tries to get header using an invalid block number
func runHeaderByNumberWithInvalidBlock(client evm.Client, runtime sdk.Runtime, wfCfg config.Config) (*evm.HeaderByNumberReply, error) {
	invalidBlockStr := wfCfg.InvalidInput
	runtime.Logger().Info("Attempting to get header using invalid block number", "invalid_block", invalidBlockStr)

	// convert to big.Int
	newBlock := big.NewInt(0)
	invalidBlock, _ := newBlock.SetString(invalidBlockStr, 10)

	runtime.Logger().Info("Starting HeaderByNumber request with parsed block number", "invalid_block", invalidBlock.String())
	headerOutput, err := client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
		BlockNumber: pb.NewBigIntFromInt(invalidBlock),
	}).Await()
	runtime.Logger().Info("HeaderByNumber with invalid block completed", "header_output", headerOutput)
	if err != nil || headerOutput == nil {
		runtime.Logger().Error("got expected error for HeaderByNumber with invalid block", "invalid_block", invalidBlockStr, "header_output", headerOutput, "error", err)
		return nil, fmt.Errorf("expected error for HeaderByNumber with invalid block '%s': %w", invalidBlockStr, err)
	}

	runtime.Logger().Info("this is not expected: HeaderByNumber with invalid block should return an error or nil", "invalid_block", invalidBlockStr, "header_output", headerOutput)
	return headerOutput, nil
}
