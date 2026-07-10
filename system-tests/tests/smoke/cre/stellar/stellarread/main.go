//go:build wasip1

package main

import (
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/stellar"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/stellar/stellarread/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunReadWorkflow)
}

func RunReadWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onStellarReadTrigger,
		),
	}, nil
}

func onStellarReadTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onStellarReadTrigger called", "workflow", cfg.WorkflowName)
	defer func() {
		if r := recover(); r != nil {
			runtime.Logger().Info("Stellar read failed: panic in onStellarReadTrigger", "workflow", cfg.WorkflowName, "panic", fmt.Sprintf("%v", r))
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	client := stellar.Client{ChainSelector: cfg.ChainSelector}
	reply, err := client.GetLatestLedger(runtime, &stellar.GetLatestLedgerRequest{}).Await()
	if err != nil {
		msg := fmt.Sprintf("Stellar read failed: GetLatestLedger error: %v", err)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName, "chainSelector", cfg.ChainSelector)
		return nil, fmt.Errorf("Stellar GetLatestLedger: %w", err)
	}
	if reply == nil {
		runtime.Logger().Info("Stellar read failed: GetLatestLedger reply is nil", "workflow", cfg.WorkflowName)
		return nil, errors.New("GetLatestLedger reply is nil")
	}
	if uint64(reply.Sequence) < cfg.MinLedgerSequence {
		msg := fmt.Sprintf("Stellar read failed: ledger sequence %d below expected minimum %d", reply.Sequence, cfg.MinLedgerSequence)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("ledger sequence %d below expected minimum %d", reply.Sequence, cfg.MinLedgerSequence)
	}

	runtime.Logger().Info("Stellar read consensus succeeded", "sequence", reply.Sequence, "workflow", cfg.WorkflowName)
	return nil, nil
}
