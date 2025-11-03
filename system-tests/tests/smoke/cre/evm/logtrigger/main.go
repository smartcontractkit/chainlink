//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"github.com/stretchr/testify/require"

	sdk "github.com/smartcontractkit/cre-sdk-go/cre"

	"gopkg.in/yaml.v3"

	logtrigger "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/logtrigger/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (logtrigger.Config, error) {
		cfg := logtrigger.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return logtrigger.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunSimpleEvmLogTriggerWorkflow)
}

func RunSimpleEvmLogTriggerWorkflow(input logtrigger.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[logtrigger.Config], error) {
	logger.Info("Trigger RunSimpleEvmLogTriggerWorkflow called")

	cfg := &evm.FilterLogTriggerRequest{
		Addresses: toByteSlices(input.Addresses),
		Topics: []*evm.TopicValues{
			{
				Values: toByteSlices(input.Topics[0].Values),
			},
		},
		Confidence: 1, // LATEST
	}

	logger.Info(fmt.Sprintf(
		"FilterLogTriggerRequest content: Addresses=%v Topics=%+v Confidence=%d",
		formatAddresses(cfg.Addresses),
		formatTopics(cfg.Topics),
		cfg.Confidence,
	))
	return cre.Workflow[logtrigger.Config]{
		cre.Handler(
			evm.LogTrigger(input.ChainSelector, cfg),
			onTrigger,
		),
	}, nil
}

func formatAddresses(addresses [][]byte) []string {
	result := make([]string, len(addresses))
	for i, addr := range addresses {
		result[i] = "0x" + hex.EncodeToString(addr)
	}
	return result
}

func formatTopics(topics []*evm.TopicValues) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		vals := make([]string, len(topic.Values))
		for j, v := range topic.Values {
			vals[j] = "0x" + hex.EncodeToString(v)
		}
		result[i] = fmt.Sprintf("values: [%s]", strings.Join(vals, ", "))
	}
	return result
}

func toByteSlices(addresses []string) [][]byte {
	result := make([][]byte, len(addresses))
	for i, addr := range addresses {
		// Assumes addresses are hex strings with or without 0x prefix
		b, err := hex.DecodeString(strings.TrimPrefix(addr, "0x"))
		if err != nil {
			panic(fmt.Sprintf("failed to decode hex string %s: %v", addr, err))
		}
		result[i] = b
	}
	return result
}

func onTrigger(cfg logtrigger.Config, runtime sdk.Runtime, outputs *evm.Log) (string, error) {
	runtime.Logger().Info("Trigger OnTrigger called", "outputs", outputs)
	t := &T{Logger: runtime.Logger()}
	require.NotNil(t, outputs, "Log input should not be nil")

	decodedMessageString, err := printDecodedData(t, runtime, cfg.Abi, cfg.Event, outputs.Data)
	if err != nil {
		runtime.Logger().Info("OnTrigger error decoding log data:", "error", err)
		return "", fmt.Errorf("OnTrigger error decoding log data: %w", err)
	}
	runtime.Logger().Info(fmt.Sprintf("OnTrigger decoded message: %s", decodedMessageString))
	client := evm.Client{ChainSelector: cfg.ChainSelector}
	txHash := sendTx(t, runtime, cfg, client, decodedMessageString)
	runtime.Logger().Info("Successfully sent transaction", "hash", common.Hash(txHash).String())
	return "success", nil
}

func sendTx(t *T, runtime sdk.Runtime, cfg logtrigger.Config, client evm.Client, msg string) []byte {
	// NOTE: This is not a right way to send a transaction. Msg must be properly encoded to trigger a proper receiver contract call.
	// In this case we just need to see transaction on chain, so it's sufficient.
	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: []byte(msg),
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	require.NoError(t, err, "failed to generate report")
	reportReply, err := client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  common.HexToAddress(cfg.Addresses[0]).Bytes(),
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: 500_000},
	}).Await()
	require.NoError(t, err, "failed to write report")
	require.NotNil(t, reportReply)
	return reportReply.TxHash
}

func printDecodedData(t *T, runtime sdk.Runtime, eventABI string, eventName string, data []byte) (string, error) {
	require.NotNil(t, eventABI, "eventABI input should not be nil")
	require.NotNil(t, eventName, "eventName input should not be nil")
	require.NotNil(t, data, "data input should not be nil")

	runtime.Logger().Info("About to read ABI and unpack data")
	parsedABI, err := abi.JSON(strings.NewReader(eventABI))
	if err != nil {
		return "", err
	}
	runtime.Logger().Info("ABI parsed successfully, about to unpack data:", "data ", data)
	event := parsedABI.Events[eventName]
	values := make(map[string]interface{})
	err = event.Inputs.UnpackIntoMap(values, data)
	if err != nil {
		return "", err
	}

	runtime.Logger().Info("Data unpacked successfully, about to format values")
	var sb strings.Builder
	first := true
	for k, v := range values {
		if !first {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("%s:%v", k, v))
		first = false
	}
	decodedData := sb.String()
	runtime.Logger().Info("Values formatted successfully:", "decodedData ", decodedData)
	return decodedData, nil
}

type T struct {
	*slog.Logger
}

func (t *T) Errorf(format string, args ...interface{}) {
	// if the log was produced by require/assert we need to split it, as engine does not allow logs longer than 1k bytes
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok && strings.Contains(msg, "Error:") && strings.Contains(msg, "Error Trace:") {
			for _, line := range strings.Split(msg, "Error:") {
				t.Logger.Error(line)
			}
			return
		}
	}
	t.Logger.Error(fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...)) // panic to stop execution
}

func (t *T) FailNow() {
	panic("Test failed. Panic to stop execution")
}
