//go:build wasip1

// Package main is the confidential counterpart of vaultsecretcron: the cron
// trigger's handler runs inside a TEE (cre.HandlerInTee) instead of on the DON,
// so the secret it reads is fetched by the enclave through the confidential
// relay rather than the workflow-DON -> vault-DON remote hop.
//
// An empty secret_key skips the secret leg, so one binary serves both phases of
// the sharded confidential tests: the phase whose executing shard hosts the
// confidential relay proves the enclave -> relay -> shared vault fetch returned
// the exact expected value, while the phase that moved to the other shard
// proves execution alone (see shard_confidential_workflows_test.go for why the
// secret leg cannot follow the swap).
package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type config struct {
	Schedule            string `json:"schedule"`
	SecretNamespace     string `json:"secret_namespace"`
	SecretKey           string `json:"secret_key"`
	ExpectedSecretValue string `json:"expected_secret_value"`
}

func main() {
	wasm.NewRunner(cre.ParseJSON[config]).Run(initWorkflow)
}

func initWorkflow(cfg *config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*config], error) {
	return cre.Workflow[*config]{
		cre.HandlerInTee(
			cron.Trigger(&cron.Config{Schedule: cfg.Schedule}),
			onTrigger,
			cre.AnyTee{},
		),
	}, nil
}

func onTrigger(cfg *config, trt cre.TeeRuntime, _ *cron.Payload) (string, error) {
	if cfg.SecretKey == "" {
		return "Confidential workflow executed in enclave", nil
	}

	secret, err := trt.GetSecret(&cre.SecretRequest{
		Id:        cfg.SecretKey,
		Namespace: cfg.SecretNamespace,
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s in the enclave: %w", cfg.SecretNamespace, cfg.SecretKey, err)
	}
	if secret.Value != cfg.ExpectedSecretValue {
		return "", fmt.Errorf("enclave-fetched secret %s/%s mismatch: got %q, want %q", cfg.SecretNamespace, cfg.SecretKey, secret.Value, cfg.ExpectedSecretValue)
	}

	return "Confidential vault secret fetched: " + secret.Value, nil
}
