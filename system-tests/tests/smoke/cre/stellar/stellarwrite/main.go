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

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/stellar/stellarwrite/config"
)

// defaultReportPayload is used when the config does not supply a hex payload.
// The Soroban receiver's on_report only needs a non-empty payload for the smoke
// test; content is opaque to the forwarder.
var defaultReportPayload = []byte{0x01, 0x02, 0x03, 0x04}

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunStellarWriteWorkflow)
}

func RunStellarWriteWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onStellarWriteTrigger,
		),
	}, nil
}

func onStellarWriteTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onStellarWriteTrigger called", "workflow", cfg.WorkflowName)

	if strings.TrimSpace(cfg.ReceiverContractID) == "" {
		err := fmt.Errorf("receiverContractID is required")
		runtime.Logger().Info("Stellar write failed: missing receiver", "workflow", cfg.WorkflowName)
		return nil, err
	}

	reportPayload, err := resolveReportPayload(cfg)
	if err != nil {
		runtime.Logger().Info(fmt.Sprintf("Stellar write failed: invalid report payload: %v", err), "workflow", cfg.WorkflowName)
		return nil, err
	}

	// The Soroban CRE forwarder verifies ed25519 signatures over
	// keccak256(keccak256(raw_report) ++ report_context); "stellar" selects the
	// Stellar key bundle in consensus report generation.
	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: reportPayload,
		EncoderName:    "stellar",
		SigningAlgo:    "ed25519",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		runtime.Logger().Info(fmt.Sprintf("Stellar write failed: generate report error: %v", err), "workflow", cfg.WorkflowName)
		return nil, err
	}

	reportResp := report.X_GeneratedCodeOnly_Unwrap()
	if len(reportResp.ReportContext) == 0 {
		err := fmt.Errorf("missing report context from generated report")
		runtime.Logger().Info("Stellar write failed: missing report context", "workflow", cfg.WorkflowName)
		return nil, err
	}
	if len(reportResp.RawReport) == 0 {
		err := fmt.Errorf("missing raw report from generated report")
		runtime.Logger().Info("Stellar write failed: missing raw report", "workflow", cfg.WorkflowName)
		return nil, err
	}
	runtime.Logger().Info(
		"Stellar write: generated report",
		"workflow", cfg.WorkflowName,
		"reportContextLen", len(reportResp.ReportContext),
		"rawReportLen", len(reportResp.RawReport),
		"sigCount", len(reportResp.Sigs),
	)

	// The Soroban forwarder expects exactly f+1 signatures; trim any extras.
	requiredSignatures := cfg.RequiredSignatures
	if requiredSignatures <= 0 {
		requiredSignatures = len(reportResp.Sigs)
	}
	if len(reportResp.Sigs) > requiredSignatures {
		reportResp.Sigs = reportResp.Sigs[:requiredSignatures]
	}
	if len(reportResp.Sigs) < requiredSignatures {
		err := fmt.Errorf("insufficient report signatures: have=%d need=%d", len(reportResp.Sigs), requiredSignatures)
		runtime.Logger().Info("Stellar write failed: too few signatures", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}

	client := stellar.Client{ChainSelector: cfg.ChainSelector}
	reply, err := client.WriteReport(runtime, &stellar.WriteCreReportRequest{
		ContractId: cfg.ReceiverContractID,
		Report:     report,
	}).Await()
	if err != nil {
		if cfg.ExpectFailure {
			return nil, fmt.Errorf("expected failed tx status in WriteReport reply, got call error instead: %w", err)
		}
		runtime.Logger().Info(fmt.Sprintf("Stellar write failed: WriteReport error: %v", err), "workflow", cfg.WorkflowName, "chainSelector", cfg.ChainSelector)
		return nil, err
	}
	if reply == nil {
		runtime.Logger().Info("Stellar write failed: WriteReport reply is nil", "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("nil WriteReport reply")
	}

	errorMsg := ""
	if reply.ErrorMessage != nil {
		errorMsg = *reply.ErrorMessage
	}

	if cfg.ExpectFailure {
		if reply.TxStatus == stellar.TxStatus_TX_STATUS_SUCCESS {
			return nil, fmt.Errorf("expected non-success tx status, got %s", reply.TxStatus.String())
		}
		runtime.Logger().Info(
			fmt.Sprintf("Stellar write failure observed as expected txStatus=%s", reply.TxStatus.String()),
			"workflow", cfg.WorkflowName, "txHash", reply.GetTxHash(), "error", errorMsg,
		)
		return nil, nil
	}

	if reply.TxStatus != stellar.TxStatus_TX_STATUS_SUCCESS {
		failMsg := fmt.Sprintf("Stellar write failed: tx status=%s error=%s", reply.TxStatus.String(), errorMsg)
		runtime.Logger().Info(failMsg, "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("unexpected tx status: %s", reply.TxStatus.String())
	}
	if reply.GetTxHash() == "" {
		runtime.Logger().Info("Stellar write failed: empty tx hash on success", "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("expected non-empty tx hash in successful WriteReport reply")
	}

	runtime.Logger().Info(
		"Stellar write consensus succeeded",
		"workflow", cfg.WorkflowName,
		"txHash", reply.GetTxHash(),
		"ledgerSequence", reply.GetLedgerSequence(),
	)
	return nil, nil
}

func resolveReportPayload(cfg config.Config) ([]byte, error) {
	h := strings.TrimSpace(cfg.ReportPayloadHex)
	if h == "" {
		return defaultReportPayload, nil
	}
	h = strings.TrimPrefix(h, "0x")
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("decode reportPayloadHex: %w", err)
	}
	if len(b) == 0 {
		return defaultReportPayload, nil
	}
	return b, nil
}
