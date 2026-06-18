//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/aptos"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswrite/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunAptosWriteWorkflow)
}

func RunAptosWriteWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onAptosWriteTrigger,
		),
	}, nil
}

func onAptosWriteTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onAptosWriteTrigger called", "workflow", cfg.WorkflowName)

	receiver, err := decodeAptosAddressHex(cfg.ReceiverHex)
	if err != nil {
		msg := fmt.Sprintf("Aptos write failed: invalid receiver address: %v", err)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName)
		return nil, err
	}

	reportPayload, err := resolveReportPayload(cfg)
	if err != nil {
		failMsg := fmt.Sprintf("Aptos write failed: invalid report payload: %v", err)
		runtime.Logger().Info(failMsg, "workflow", cfg.WorkflowName)
		return nil, err
	}

	report, err := runtime.GenerateReport(&sdkpb.ReportRequest{
		EncodedPayload: reportPayload,
		// Select Aptos key bundle path in consensus report generation.
		EncoderName: "aptos",
		// Aptos forwarder verifies ed25519 signatures over blake2b_256(raw_report).
		SigningAlgo: "ed25519",
		HashingAlgo: "blake2b_256",
	}).Await()
	if err != nil {
		failMsg := fmt.Sprintf("Aptos write failed: generate report error: %v", err)
		runtime.Logger().Info(failMsg, "workflow", cfg.WorkflowName)
		return nil, err
	}
	reportResp := report.X_GeneratedCodeOnly_Unwrap()
	if len(reportResp.ReportContext) == 0 {
		err := fmt.Errorf("missing report context from generated report")
		runtime.Logger().Info("Aptos write failed: missing report context", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}
	if len(reportResp.ReportContext) != 96 {
		err := fmt.Errorf("unexpected report context length: got=%d want=96", len(reportResp.ReportContext))
		runtime.Logger().Info("Aptos write failed: invalid report context length", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}
	if len(reportResp.RawReport) == 0 {
		err := fmt.Errorf("missing raw report from generated report")
		runtime.Logger().Info("Aptos write failed: missing raw report", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}
	// Preserve generated report bytes as-is; Aptos capability handles wire-format packing.
	reportVersion := int(reportResp.RawReport[0])
	runtime.Logger().Info(
		"Aptos write: generated report details",
		"workflow", cfg.WorkflowName,
		"reportContextLen", len(reportResp.ReportContext),
		"rawReportLen", len(reportResp.RawReport),
		"reportVersion", reportVersion,
	)

	runtime.Logger().Info(
		"Aptos write: generated report",
		"workflow", cfg.WorkflowName,
		"sigCount", len(reportResp.Sigs),
	)
	if len(reportResp.Sigs) > 0 {
		runtime.Logger().Info(
			"Aptos write: first signature details",
			"workflow", cfg.WorkflowName,
			"firstSigLen", len(reportResp.Sigs[0].Signature),
			"firstSignerID", reportResp.Sigs[0].SignerId,
		)
	}
	requiredSignatures := cfg.RequiredSignatures
	if requiredSignatures <= 0 {
		requiredSignatures = len(reportResp.Sigs)
	}
	if len(reportResp.Sigs) > requiredSignatures {
		reportResp.Sigs = reportResp.Sigs[:requiredSignatures]
		runtime.Logger().Info(
			"Aptos write: trimmed report signatures for forwarder",
			"workflow", cfg.WorkflowName,
			"requiredSignatures", requiredSignatures,
			"sigCount", len(reportResp.Sigs),
		)
	}
	if len(reportResp.Sigs) < requiredSignatures {
		err := fmt.Errorf("insufficient report signatures: have=%d need=%d", len(reportResp.Sigs), requiredSignatures)
		runtime.Logger().Info("Aptos write failed: report has fewer signatures than required", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, err
	}

	client := aptos.Client{ChainSelector: cfg.ChainSelector}
	runtime.Logger().Info(
		"Aptos write: using gas config",
		"workflow", cfg.WorkflowName,
		"chainSelector", cfg.ChainSelector,
		"maxGasAmount", cfg.MaxGasAmount,
		"gasUnitPrice", cfg.GasUnitPrice,
	)
	reply, err := client.WriteReport(runtime, &aptos.WriteCreReportRequest{
		Receiver: receiver,
		Report:   report,
		GasConfig: &aptos.GasConfig{
			MaxGasAmount: cfg.MaxGasAmount,
			GasUnitPrice: cfg.GasUnitPrice,
		},
	}).Await()
	if err != nil {
		if cfg.ExpectFailure {
			runtime.Logger().Info(
				"Aptos write failed: expected failure path requires non-empty failed tx hash",
				"workflow", cfg.WorkflowName,
				"txStatus", "call_error",
				"txHash", "",
				"error", err.Error(),
			)
			return nil, fmt.Errorf("expected failed tx hash in WriteReport reply, got error instead: %w", err)
		}
		failMsg := fmt.Sprintf("Aptos write failed: WriteReport error: %v", err)
		runtime.Logger().Info(failMsg, "workflow", cfg.WorkflowName, "chainSelector", cfg.ChainSelector)
		return nil, err
	}
	if reply == nil {
		runtime.Logger().Info("Aptos write failed: WriteReport reply is nil", "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("nil WriteReport reply")
	}
	if cfg.ExpectFailure {
		if reply.TxStatus == aptos.TxStatus_TX_STATUS_SUCCESS {
			errorMsg := ""
			if reply.ErrorMessage != nil {
				errorMsg = *reply.ErrorMessage
			}
			runtime.Logger().Info(
				"Aptos write failed: expected non-success tx status",
				"workflow", cfg.WorkflowName,
				"txStatus", reply.TxStatus.String(),
				"error", errorMsg,
			)
			return nil, fmt.Errorf("expected non-success tx status, got %s", reply.TxStatus.String())
		}
		txHashRaw := reply.GetTxHash()
		if txHashRaw == "" {
			runtime.Logger().Info(
				"Aptos write failed: expected failed tx hash but got empty hash",
				"workflow", cfg.WorkflowName,
			)
			return nil, fmt.Errorf("expected failed tx hash in WriteReport reply")
		}

		txHash, err := normalizeTxHash(txHashRaw)
		if err != nil {
			runtime.Logger().Info("Aptos write failed: invalid failed tx hash format", "workflow", cfg.WorkflowName, "error", err.Error())
			return nil, fmt.Errorf("invalid failed tx hash format: %w", err)
		}

		errorMsg := ""
		if reply.ErrorMessage != nil {
			errorMsg = *reply.ErrorMessage
		}
		runtime.Logger().Info(
			fmt.Sprintf("Aptos write failure observed as expected txHash=%s", txHash),
			"workflow", cfg.WorkflowName,
			"txStatus", reply.TxStatus.String(),
			"txHash", txHash,
			"error", errorMsg,
		)
		return nil, nil
	}
	if reply.TxStatus != aptos.TxStatus_TX_STATUS_SUCCESS {
		errorMsg := ""
		if reply.ErrorMessage != nil {
			errorMsg = *reply.ErrorMessage
		}
		failMsg := fmt.Sprintf("Aptos write failed: tx status=%s error=%s", reply.TxStatus.String(), errorMsg)
		runtime.Logger().Info(failMsg, "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("unexpected tx status: %s", reply.TxStatus.String())
	}
	txHashRaw := reply.GetTxHash()
	if txHashRaw == "" {
		runtime.Logger().Info(
			"Aptos write failed: expected successful tx hash but got empty hash",
			"workflow", cfg.WorkflowName,
			"txStatus", reply.TxStatus.String(),
		)
		return nil, fmt.Errorf("expected non-empty tx hash in successful WriteReport reply")
	}

	txHash, err := normalizeTxHash(txHashRaw)
	if err != nil {
		runtime.Logger().Info("Aptos write failed: invalid tx hash format", "workflow", cfg.WorkflowName, "error", err.Error())
		return nil, fmt.Errorf("invalid tx hash format: %w", err)
	}

	runtime.Logger().Info("Aptos write capability succeeded", "workflow", cfg.WorkflowName, "txHash", txHash)
	return nil, nil
}

func resolveReportPayload(cfg config.Config) ([]byte, error) {
	if strings.TrimSpace(cfg.ReportPayloadHex) != "" {
		trimmed := strings.TrimPrefix(strings.TrimSpace(cfg.ReportPayloadHex), "0x")
		if trimmed == "" {
			return nil, fmt.Errorf("empty hex payload")
		}
		raw, err := hex.DecodeString(trimmed)
		if err != nil {
			return nil, fmt.Errorf("decode hex payload: %w", err)
		}
		return raw, nil
	}

	msg := cfg.ReportMessage
	if msg == "" {
		msg = "Aptos write workflow executed successfully"
	}
	return []byte(msg), nil
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

var aptosHashRe = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)

func normalizeTxHash(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	if !aptosHashRe.MatchString(s) {
		return "", fmt.Errorf("expected 32-byte tx hash, got %q", raw)
	}
	return "0x" + s, nil
}
