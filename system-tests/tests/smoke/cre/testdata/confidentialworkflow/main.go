//go:build wasip1

// Package main is the E2E test workflow used by the engine-path test.
// It exercises three capability paths from inside a confidential workflow:
//   - GetSecret → VaultDON (remote dispatch through the relay DON)
//   - http.SendRequest + ConsensusMedianAggregation → intercepted locally
//     by the enclave (http-actions + consensus/Simple both handled in-process)
//   - ReportFromDon + evm.WriteReport → routed *out* of the enclave to the DONs,
//     which is how a TEE handler reaches consensus-bound capabilities. The report
//     lands in a PermissionlessFeedsConsumer the test then reads back on-chain.
//
// Each success is marked in the workflow engine logs for the test to scrape.
package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// writeGasLimit mirrors the upstream proof-of-reserve example's limit.
const writeGasLimit = 400_000

type config struct {
	EchoURL string `json:"echo_url"`

	// Chain-write leg. Empty ConsumerAddress disables it, so the test can run the
	// secret + http legs alone if it ever needs to.
	ConsumerAddress string `json:"consumer_address"`
	ChainSelector   uint64 `json:"chain_selector"`
	FeedID          string `json:"feed_id"`
	Price           uint64 `json:"price"`
}

// feedReport matches PermissionlessFeedsConsumer's ReceivedFeedReport struct, which
// its onReport decodes with abi.decode(rawReport, (ReceivedFeedReport[])).
type feedReport struct {
	FeedID    [32]byte
	Timestamp uint32
	Price     *big.Int
}

func main() {
	wasm.NewRunner(cre.ParseJSON[config]).Run(initWorkflow)
}

func initWorkflow(_ *config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*config], error) {
	return cre.Workflow[*config]{
		cre.HandlerInTee(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			handleTrigger,
			cre.AnyTee{},
		),
	}, nil
}

func handleTrigger(cfg *config, trt cre.TeeRuntime, payload *cron.Payload) (any, error) {
	secret, err := trt.GetSecret(&sdkpb.SecretRequest{Id: "MOCK_SECRET"}).Await()
	if err != nil {
		return nil, err
	}
	// DO NOT log secrets in production workflows. We only do it here so the
	// test can scrape the value out of workflow-DON logs and confirm that the
	// VaultDON-routed secret-fetch path actually delivered the right value
	// into the WASM. Real users: don't follow this pattern.
	trt.Logger().Info("engine-test-secret", "value", secret.Value)

	result := map[string]any{"secret": secret.Value}

	if cfg.EchoURL != "" {
		client := &http.Client{}
		status, httpErr := fetchEchoStatus(cfg, trt, client)
		if httpErr != nil {
			trt.Logger().Error("engine-test-http-failed", "error", httpErr.Error())
			return nil, httpErr
		}
		trt.Logger().Info("engine-test-http", "status", status, "url", cfg.EchoURL)
		result["http_status"] = status
	}

	if cfg.ConsumerAddress != "" {
		txHash, writeErr := writeFeedReport(cfg, trt, payload)
		if writeErr != nil {
			trt.Logger().Error("engine-test-write-failed", "error", writeErr.Error())
			return nil, writeErr
		}
		trt.Logger().Info("engine-test-write", "txHash", txHash, "receiver", cfg.ConsumerAddress, "price", cfg.Price)
		result["tx_hash"] = txHash
	}

	return result, nil
}

func fetchEchoStatus(cfg *config, trt cre.TeeRuntime, client *http.Client) (int32, error) {
	resp, err := client.SendRequestInTee(trt, &http.Request{
		Url:    cfg.EchoURL,
		Method: "POST",
		Body:   []byte("hello from engine-test"),
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}).Await()
	if err != nil {
		return 0, err
	}
	return int32(resp.StatusCode), nil
}

// writeFeedReport generates a DON-signed report over a single feed value and writes
// it to the consumer contract. Both legs leave the enclave: cre.TeeRuntime exposes
// no GenerateReport (report signing needs the DONs) and the evm write capability
// takes a cre.Runtime, so the report comes from ReportFromDon and the write goes
// through UsingTheDons.
func writeFeedReport(cfg *config, trt cre.TeeRuntime, payload *cron.Payload) (string, error) {
	feedID, err := parseFeedID(cfg.FeedID)
	if err != nil {
		return "", err
	}

	// Use the trigger's scheduled time rather than a wall clock inside the enclave,
	// matching the chain_write canary.
	encoded, err := encodeFeedReports([]feedReport{{
		FeedID:    feedID,
		Timestamp: uint32(payload.ScheduledExecutionTime.AsTime().Unix()),
		Price:     new(big.Int).SetUint64(cfg.Price),
	}})
	if err != nil {
		return "", fmt.Errorf("encode feed report: %w", err)
	}

	report, err := trt.ReportFromDon(&cre.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return "", fmt.Errorf("report from don: %w", err)
	}

	evmClient := &evm.Client{ChainSelector: cfg.ChainSelector}
	out, err := evmClient.WriteReport(trt.UsingTheDons(), &evm.WriteCreReportRequest{
		Receiver:  common.HexToAddress(cfg.ConsumerAddress).Bytes(),
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: writeGasLimit},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("write report: %w", err)
	}
	if out.ErrorMessage != nil && *out.ErrorMessage != "" {
		return "", fmt.Errorf("write report rejected: %s", *out.ErrorMessage)
	}

	return "0x" + hex.EncodeToString(out.TxHash), nil
}

func parseFeedID(s string) ([32]byte, error) {
	var id [32]byte
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return id, fmt.Errorf("decode feed id %q: %w", s, err)
	}
	if len(b) != 32 {
		return id, fmt.Errorf("feed id %q decoded to %d bytes, want 32", s, len(b))
	}
	copy(id[:], b)
	return id, nil
}

// encodeFeedReports packs reports the way PermissionlessFeedsConsumer.onReport
// decodes them: a single dynamic array of (bytes32, uint32, uint224) tuples.
func encodeFeedReports(reports []feedReport) ([]byte, error) {
	typ, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "FeedID", Type: "bytes32"},
		{Name: "Timestamp", Type: "uint32"},
		{Name: "Price", Type: "uint224"},
	})
	if err != nil {
		return nil, fmt.Errorf("build abi type: %w", err)
	}
	return abi.Arguments{{Name: "Reports", Type: typ}}.Pack(reports)
}
