//go:build wasip1

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/aptos"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswriteroundtrip/config"
)

type feedEntry struct {
	FeedID string `json:"feed_id"`
	Feed   struct {
		Benchmark string `json:"benchmark"`
	} `json:"feed"`
}

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunAptosWriteReadRoundtripWorkflow)
}

func RunAptosWriteReadRoundtripWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onAptosWriteReadRoundtripTrigger,
		),
	}, nil
}

func onAptosWriteReadRoundtripTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onAptosWriteReadRoundtripTrigger called", "workflow", cfg.WorkflowName)

	receiverBytes, err := decodeAptosAddressHex(cfg.ReceiverHex)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver address: %w", err)
	}

	reportPayload, err := resolveReportPayload(cfg.ReportPayloadHex)
	if err != nil {
		return nil, fmt.Errorf("invalid report payload: %w", err)
	}

	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: reportPayload,
		EncoderName:    "aptos",
		SigningAlgo:    "ed25519",
		HashingAlgo:    "blake2b_256",
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("generate report error: %w", err)
	}
	reportResp := report.X_GeneratedCodeOnly_Unwrap()
	if len(reportResp.ReportContext) != 96 {
		return nil, fmt.Errorf("invalid report context length: got=%d want=96", len(reportResp.ReportContext))
	}
	if len(reportResp.RawReport) == 0 {
		return nil, fmt.Errorf("missing raw report")
	}

	requiredSignatures := cfg.RequiredSignatures
	if requiredSignatures <= 0 {
		requiredSignatures = len(reportResp.Sigs)
	}
	if len(reportResp.Sigs) < requiredSignatures {
		return nil, fmt.Errorf("insufficient report signatures: have=%d need=%d", len(reportResp.Sigs), requiredSignatures)
	}
	if len(reportResp.Sigs) > requiredSignatures {
		reportResp.Sigs = reportResp.Sigs[:requiredSignatures]
	}

	client := aptos.Client{ChainSelector: cfg.ChainSelector}
	reply, err := client.WriteReport(runtime, &aptos.WriteCreReportRequest{
		Receiver: receiverBytes,
		Report:   report,
		GasConfig: &aptos.GasConfig{
			MaxGasAmount: cfg.MaxGasAmount,
			GasUnitPrice: cfg.GasUnitPrice,
		},
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("WriteReport error: %w", err)
	}
	if reply == nil {
		return nil, fmt.Errorf("nil WriteReport reply")
	}
	if reply.TxStatus != aptos.TxStatus_TX_STATUS_SUCCESS {
		return nil, fmt.Errorf("unexpected tx status: %s", reply.TxStatus.String())
	}

	viewReply, err := client.View(runtime, &aptos.ViewRequest{
		Payload: &aptos.ViewPayload{
			Module: &aptos.ModuleID{
				Address: receiverBytes,
				Name:    "registry",
			},
			Function: "get_feeds",
		},
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("Aptos View(%s::registry::get_feeds): %w", normalizeHex(cfg.ReceiverHex), err)
	}
	if viewReply == nil || len(viewReply.Data) == 0 {
		return nil, fmt.Errorf("empty view reply for %s::registry::get_feeds", normalizeHex(cfg.ReceiverHex))
	}

	benchmark, found, parseErr := parseBenchmark(viewReply.Data, cfg.FeedIDHex)
	if parseErr != nil {
		return nil, fmt.Errorf("parse benchmark view reply: %w", parseErr)
	}
	if !found {
		return nil, fmt.Errorf("feed %s not found in get_feeds reply", cfg.FeedIDHex)
	}
	if benchmark != cfg.ExpectedBenchmark {
		return nil, fmt.Errorf("benchmark mismatch: got=%d want=%d", benchmark, cfg.ExpectedBenchmark)
	}

	runtime.Logger().Info(
		"Aptos write/read consensus succeeded",
		"workflow", cfg.WorkflowName,
		"benchmark", benchmark,
		"feedID", normalizeHex(cfg.FeedIDHex),
	)
	return nil, nil
}

func resolveReportPayload(reportPayloadHex string) ([]byte, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(reportPayloadHex), "0x")
	if trimmed == "" {
		return nil, fmt.Errorf("empty hex payload")
	}
	raw, err := hex.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("decode hex payload: %w", err)
	}
	return raw, nil
}

func decodeAptosAddressHex(in string) ([]byte, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(in), "0x")
	if trimmed == "" {
		return nil, fmt.Errorf("empty address")
	}
	if len(trimmed)%2 != 0 {
		trimmed = "0" + trimmed
	}
	raw, err := hex.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("decode hex address: %w", err)
	}
	if len(raw) > 32 {
		return nil, fmt.Errorf("address too long: %d bytes", len(raw))
	}
	out := make([]byte, 32)
	copy(out[32-len(raw):], raw)
	return out, nil
}

func parseBenchmark(data []byte, feedIDHex string) (uint64, bool, error) {
	normalizedFeedID := normalizeHex(feedIDHex)
	if normalizedFeedID == "" {
		return 0, false, fmt.Errorf("empty feed id")
	}

	var wrapped [][]feedEntry
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped) > 0 {
		for _, entry := range wrapped[0] {
			if normalizeHex(entry.FeedID) != normalizedFeedID {
				continue
			}
			v, convErr := strconv.ParseUint(strings.TrimSpace(entry.Feed.Benchmark), 10, 64)
			if convErr != nil {
				return 0, false, fmt.Errorf("parse benchmark %q: %w", entry.Feed.Benchmark, convErr)
			}
			return v, true, nil
		}
		return 0, false, nil
	}

	var direct []feedEntry
	if err := json.Unmarshal(data, &direct); err != nil {
		return 0, false, fmt.Errorf("decode get_feeds payload: %w", err)
	}
	for _, entry := range direct {
		if normalizeHex(entry.FeedID) != normalizedFeedID {
			continue
		}
		v, convErr := strconv.ParseUint(strings.TrimSpace(entry.Feed.Benchmark), 10, 64)
		if convErr != nil {
			return 0, false, fmt.Errorf("parse benchmark %q: %w", entry.Feed.Benchmark, convErr)
		}
		return v, true, nil
	}
	return 0, false, nil
}

func normalizeHex(in string) string {
	s := strings.TrimSpace(strings.ToLower(in))
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "0x0"
	}
	return "0x" + s
}
