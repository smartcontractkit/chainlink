//go:build wasip1

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"main/types"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	workflowpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
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

func RunProofOfReservesWorkflow(config types.WorkflowConfig, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[types.WorkflowConfig], error) {
	return cre.Workflow[types.WorkflowConfig]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "0 * * * * *"}), // every minute
			onTrigger,
		),
	}, nil
}

func onTrigger(config types.WorkflowConfig, runtime cre.Runtime, payload *cron.Payload) (string, error) {
	runtime.Logger().Info("PoR workflow started", "payload", payload)

	// get balance with BalanceAt()
	evmClient := evm.Client{ChainSelector: chain_selectors.GETH_TESTNET.Selector}
	addressesToRead := config.BalanceReaderConfig.AddressesToRead
	if len(addressesToRead) < 2 {
		runtime.Logger().Info("need at least 2 addresses to read balances from:", "addresses", addressesToRead)
		return "", fmt.Errorf("need at least 2 addresses to read balances from: %v", addressesToRead)
	}

	// For testing purposes, there is no handling of index out of range or nil cases.
	// It allows for the configuration of empty addresses, a single address, or zero balances.
	// The happy-path scenario in the system tests guarantees there are at least two addresses present.
	// However, in real-world usage, it is advisable to implement
	// proper validation for the configuration and handle possible errors.
	addressToRead_1 := addressesToRead[0]
	balanceAtOutput, err := evmClient.BalanceAt(runtime, &evm.BalanceAtRequest{
		Account:     addressToRead_1.Bytes(),
		BlockNumber: nil,
	}).Await()
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("[logger] failed to get on-chain balance: %v", err))
		return "", fmt.Errorf("failed to get on-chain balance: %w", err)
	}
	balanceAtResult := pb.NewIntFromBigInt(balanceAtOutput.Balance)
	runtime.Logger().With().Info(fmt.Sprintf("[logger] Got on-chain balance with BalanceAt() for address %s: %s", addressToRead_1, balanceAtResult.String()))

	// get balance with CallContract
	readBalancesParsedABI, err := getReadBalancesContractABI(runtime, balanceReaderABIJson)
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("failed to get ReadBalances ABI: %v", err))
		return "", fmt.Errorf("failed to get ReadBalances ABI: %w", err)
	}

	// To test that reading the contract is operational, it is sufficient to use 1 address.
	// For testing purposes, there is no index out of range or nil handling,
	// see comments above for more details (TL:DR; implement your own proper validation)
	addressToRead_2 := addressesToRead[1]
	readBalancesOutput, err := readBalancesFromContract([]common.Address{addressToRead_2}, readBalancesParsedABI, evmClient, runtime, config)
	if err != nil {
		return "", fmt.Errorf("failed to read balances from contract: %w", err)
	}

	if readBalancesOutput == nil || readBalancesOutput.Data == nil {
		runtime.Logger().Error("CallContract returned nil output, check if BalanceReader contract is deployed")
		return "", fmt.Errorf("CallContract returned nil output")
	}

	var readBalancePrices []*big.Int
	methodName := "getNativeBalances"
	err = readBalancesParsedABI.UnpackIntoInterface(&readBalancePrices, methodName, readBalancesOutput.Data)
	if err != nil {
		return "", fmt.Errorf("failed to read CallContract output: %w", err)
	}
	runtime.Logger().With().Info(fmt.Sprintf("Read on-onchain balances for addresses %v: %v", addressesToRead, &readBalancePrices))

	// get total on-chain balance
	allOnchainBalances := append(readBalancePrices, balanceAtResult)
	var totalOnChainBalance big.Int
	for _, balance := range allOnchainBalances {
		totalOnChainBalance = *totalOnChainBalance.Add(&totalOnChainBalance, balance)
	}
	runtime.Logger().With().Info(fmt.Sprintf("Total on-chain balance for addresses %v", &totalOnChainBalance))

	totalPriceOutput, err := cre.RunInNodeMode(config, runtime,
		func(config types.WorkflowConfig, nodeRuntime cre.NodeRuntime) (priceOutput, error) {
			httpOutput, err := getHTTPPrice(config, nodeRuntime)
			if err != nil {
				return priceOutput{}, fmt.Errorf("failed to get HTTP price: %w", err)
			}
			httpOutput.Price.Add(httpOutput.Price, &totalOnChainBalance)
			return httpOutput, nil
		},
		cre.ConsensusIdenticalAggregation[priceOutput](),
	).Await()
	if err != nil {
		return "", fmt.Errorf("failed to get price: %w", err)
	}
	runtime.Logger().With().Info(fmt.Sprintf("Got price: %s, for feed: %s, at time: %d", totalPriceOutput.Price.String(), common.Bytes2Hex(totalPriceOutput.FeedID[:]), totalPriceOutput.Timestamp))

	encodedPrice, err := encodeReports([]priceOutput{totalPriceOutput})
	if err != nil {
		return "", fmt.Errorf("failed to pack price report: %w", err)
	}

	report, err := runtime.GenerateReport(&workflowpb.ReportRequest{
		EncodedPayload: encodedPrice,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}
	runtime.Logger().With().Info(fmt.Sprintln("final report generated"))

	receiver, err := common.ParseHexOrString(config.DataFeedsCacheAddress)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %w", err)
	}

	wrOutput, err := evmClient.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  receiver,
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: 400000},
	}).Await()
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("[logger] failed to write report on-chain: %v", err))
		return "", fmt.Errorf("failed to write report on-chain: %w", err)
	}
	runtime.Logger().With().Info("Submitted report on-chain")

	var message = "PoR Workflow successfully completed"
	if wrOutput.ErrorMessage != nil {
		message = *wrOutput.ErrorMessage
	}

	return message, nil
}

func getReadBalancesContractABI(runtime cre.Runtime, balanceReaderABI string) (abi.ABI, error) {
	parsedABI, err := abi.JSON(strings.NewReader(balanceReaderABI))
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("failed to parse ABI: %v", err))
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %w", err)
	}
	runtime.Logger().With().Info(fmt.Sprintln("Parsed ABI successfully"))
	return parsedABI, nil
}

func readBalancesFromContract(addresses []common.Address, readBalancesABI abi.ABI, evmClient evm.Client, runtime cre.Runtime, config types.WorkflowConfig) (*evm.CallContractReply, error) {
	methodName := "getNativeBalances"
	packedData, err := readBalancesABI.Pack(methodName, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to pack read balances call: %w", err)
	}
	readBalancesOutput, err := evmClient.CallContract(runtime, &evm.CallContractRequest{
		Call: &evm.CallMsg{
			To:   common.HexToAddress(config.BalanceReaderAddress).Bytes(),
			Data: packedData,
		},
	}).Await()
	if err != nil {
		runtime.Logger().Error(fmt.Sprintf("[logger] failed to get balances %v: %v", addresses, err))
		return nil, fmt.Errorf("failed to get balances for addresses %v: %w", addresses, err)
	}
	runtime.Logger().With().Info(fmt.Sprintf("Got raw CallContract output: %s", hex.EncodeToString(readBalancesOutput.Data)))
	return readBalancesOutput, nil
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (types.WorkflowConfig, error) {
		cfg := types.WorkflowConfig{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return types.WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if cfg.AuthKeySecretName != "" {
			cfg.AuthKey = sdk.SecretValue(cfg.AuthKeySecretName)
		}

		return cfg, nil
	}).Run(RunProofOfReservesWorkflow)
}

type priceOutput struct {
	FeedID    [32]byte
	Timestamp uint32
	Price     *big.Int
}

type trueUSDResponse struct {
	AccountName string    `json:"accountName"`
	TotalTrust  float64   `json:"totalTrust"`
	Ripcord     bool      `json:"ripcord"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func getHTTPPrice(config types.WorkflowConfig, runtime cre.NodeRuntime) (priceOutput, error) {
	httpClient := &http.Client{}

	feedID, err := convertFeedIDtoBytes(config.FeedID)
	if err != nil {
		return priceOutput{}, fmt.Errorf("cannot convert feedID to bytes : %w : %b", err, feedID)
	}

	fetchRequest := http.Request{
		Url:       config.URL + "?feedID=" + config.FeedID,
		Method:    "GET",
		TimeoutMs: 5000,
	}

	if string(config.AuthKey) != "" {
		fetchRequest.Headers = map[string]string{
			"Authorization": string(config.AuthKey),
		}
	}

	r, err := httpClient.SendRequest(runtime, &fetchRequest).Await()
	if err != nil {
		return priceOutput{}, fmt.Errorf("failed to await price response from %s and %v err: %w", fetchRequest.String(), fetchRequest.Headers, err)
	}

	var resp trueUSDResponse
	if err = json.Unmarshal(r.Body, &resp); err != nil {
		return priceOutput{}, fmt.Errorf("failed to unmarshal price response: %w", err)
	}

	runtime.Logger().With().Info(fmt.Sprintf("Response is account name: %s, totalTrust: %.10f, ripcord: %v, updatedAt: %s", resp.AccountName, resp.TotalTrust, resp.Ripcord, resp.UpdatedAt.String()))

	if resp.Ripcord {
		runtime.Logger().With(
			"feedID", config.FeedID,
		).Info(fmt.Sprintf("ripcord flag set for feed ID %s", config.FeedID))
		return priceOutput{}, sdk.BreakErr
	}

	return priceOutput{
		FeedID:    feedID, // TrueUSD
		Timestamp: uint32(resp.UpdatedAt.Unix()),
		Price:     big.NewInt(int64(resp.TotalTrust * 100)), // Convert to integer cents
	}, nil
}

func convertFeedIDtoBytes(feedIDStr string) ([32]byte, error) {
	if feedIDStr == "" {
		return [32]byte{}, fmt.Errorf("feedID string is empty")
	}

	if len(feedIDStr) < 2 {
		return [32]byte{}, fmt.Errorf("feedID string too short: %q", feedIDStr)
	}

	b, err := hex.DecodeString(feedIDStr[2:])
	if err != nil {
		return [32]byte{}, err
	}

	if len(b) < 32 {
		nb := [32]byte{}
		copy(nb[:], b[:])
		return nb, err
	}

	return [32]byte(b), nil
}

func encodeReports(reports []priceOutput) ([]byte, error) {
	typ, err := abi.NewType("tuple[]", "",
		[]abi.ArgumentMarshaling{
			{Name: "FeedID", Type: "bytes32"},
			{Name: "Timestamp", Type: "uint32"},
			{Name: "Price", Type: "uint224"},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	args := abi.Arguments{
		{
			Name: "Reports",
			Type: typ,
		},
	}
	return args.Pack(reports)
}
