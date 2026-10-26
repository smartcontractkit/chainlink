//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecretcron/config"
)

func main() {
	wasm.NewRunner(func(configBytes []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunVaultSecretCronWorkflow)
}

func RunVaultSecretCronWorkflow(cfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: cfg.Schedule}),
			onTrigger,
		),
	}, nil
}

func onTrigger(cfg config.Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	secret, err := runtime.GetSecret(&cre.SecretRequest{
		Namespace: cfg.SecretNamespace,
		Id:        cfg.SecretKey,
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", cfg.SecretNamespace, cfg.SecretKey, err)
	}

	runtime.Logger().Info(fmt.Sprintf("Vault secret fetched: %s", secret.Value))
	return secret.Value, nil
}
