//go:build wasip1

// Confidential counterpart of vaultsecretcron: the cron trigger's handler runs
// in a TEE, so the secret it reads is fetched by the enclave through the
// confidential relay. An empty secret_key skips the secret leg.
package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/confidentialvaultsecretcron/config"
)

func main() {
	wasm.NewRunner(cre.ParseJSON[config.Config]).Run(initWorkflow)
}

func initWorkflow(cfg *config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*config.Config], error) {
	return cre.Workflow[*config.Config]{
		cre.HandlerInTee(
			cron.Trigger(&cron.Config{Schedule: cfg.Schedule}),
			onTrigger,
			cre.AnyTee{},
		),
	}, nil
}

func onTrigger(cfg *config.Config, trt cre.TeeRuntime, _ *cron.Payload) (string, error) {
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
