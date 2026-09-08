//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/stellar"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/stellar/datafeeds/write/config"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/stellar/datafeeds/write/report"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunStellarDataFeedsWriteWorkflow)
}

func RunStellarDataFeedsWriteWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onStellarDataFeedsWriteTrigger,
		),
	}, nil
}

func onStellarDataFeedsWriteTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onStellarDataFeedsWriteTrigger called", "workflow", cfg.WorkflowName)

	if strings.TrimSpace(cfg.CacheContractID) == "" {
		err := fmt.Errorf("cacheContractID is required")
		runtime.Logger().Info("Stellar DF write failed: missing cache contract id", "workflow", cfg.WorkflowName)
		return nil, err
	}
	dataIDBytes, err := hex.DecodeString(strings.TrimPrefix(strings.TrimSpace(cfg.DataIDHex), "0x"))
	if err != nil {
		runtime.Logger().Info("Stellar DF write failed: invalid data id", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, fmt.Errorf("dataIDHex is not valid hex: %w", err)
	}
	if len(dataIDBytes) != 32 {
		err := fmt.Errorf("dataIDHex must decode to 32 bytes, got %d", len(dataIDBytes))
		runtime.Logger().Info("Stellar DF write failed: invalid data id", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}

	ts := uint64(payload.ScheduledExecutionTime.AsTime().Unix())
	if ts == 0 {
		err := fmt.Errorf("missing scheduled execution time")
		runtime.Logger().Info("Stellar DF write failed: missing scheduled execution time", "workflow", cfg.WorkflowName)
		return nil, err
	}
	reportPayload, err := report.EncodeEntries([32]byte(dataIDBytes), cfg.Answer, ts)
	if err != nil {
		runtime.Logger().Info(fmt.Sprintf("Stellar DF write failed: encode report entries: %v", err), "workflow", cfg.WorkflowName)
		return nil, err
	}

	generated, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: reportPayload,
		EncoderName:    "stellar",
		SigningAlgo:    "ed25519",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		runtime.Logger().Info(fmt.Sprintf("Stellar DF write failed: generate report error: %v", err), "workflow", cfg.WorkflowName)
		return nil, err
	}

	reportResp := generated.X_GeneratedCodeOnly_Unwrap()
	if len(reportResp.ReportContext) == 0 {
		err := fmt.Errorf("missing report context from generated report")
		runtime.Logger().Info("Stellar DF write failed: missing report context", "workflow", cfg.WorkflowName)
		return nil, err
	}
	if len(reportResp.RawReport) == 0 {
		err := fmt.Errorf("missing raw report from generated report")
		runtime.Logger().Info("Stellar DF write failed: missing raw report", "workflow", cfg.WorkflowName)
		return nil, err
	}

	requiredSignatures := cfg.RequiredSignatures
	if requiredSignatures <= 0 {
		requiredSignatures = len(reportResp.Sigs)
	}
	if len(reportResp.Sigs) > requiredSignatures {
		reportResp.Sigs = reportResp.Sigs[:requiredSignatures]
	}
	if len(reportResp.Sigs) < requiredSignatures {
		err := fmt.Errorf("insufficient report signatures: have=%d need=%d", len(reportResp.Sigs), requiredSignatures)
		runtime.Logger().Info("Stellar DF write failed: too few signatures", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}

	client := stellar.Client{ChainSelector: cfg.ChainSelector}
	reply, err := client.WriteReport(runtime, &stellar.WriteCreReportRequest{
		ContractId: cfg.CacheContractID,
		Report:     generated,
	}).Await()
	if err != nil {
		runtime.Logger().Info(fmt.Sprintf("Stellar DF write failed: WriteReport error: %v", err), "workflow", cfg.WorkflowName, "chainSelector", cfg.ChainSelector)
		return nil, err
	}
	if reply == nil {
		runtime.Logger().Info("Stellar DF write failed: WriteReport reply is nil", "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("nil WriteReport reply")
	}

	if reply.TxStatus != stellar.TxStatus_TX_STATUS_SUCCESS {
		runtime.Logger().Info(fmt.Sprintf("Stellar DF write failed: tx status=%s error=%s", reply.TxStatus.String(), reply.GetErrorMessage()), "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("unexpected tx status: %s", reply.TxStatus.String())
	}
	if reply.GetTxHash() == "" {
		runtime.Logger().Info("Stellar DF write failed: empty tx hash on success", "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("expected non-empty tx hash in successful WriteReport reply")
	}

	runtime.Logger().Info(
		"Stellar DF write succeeded",
		"workflow", cfg.WorkflowName,
		"txHash", reply.GetTxHash(),
		"ledgerSequence", reply.GetLedgerSequence(),
	)
	return nil, nil
}
