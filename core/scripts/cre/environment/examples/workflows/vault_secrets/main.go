//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/vault_secrets/types"
)

func main() {
	wasm.NewRunner(func(configBytes []byte) (types.WorkflowConfig, error) {
		cfg := types.WorkflowConfig{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return types.WorkflowConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunVaultSecretsWorkflow)
}

func RunVaultSecretsWorkflow(config types.WorkflowConfig, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[types.WorkflowConfig], error) {
	workflows := cre.Workflow[types.WorkflowConfig]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(config types.WorkflowConfig, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	secret, err := runtime.GetSecret(&cre.SecretRequest{
		Namespace: config.SecretNamespace,
		Id:        config.SecretKey,
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", config.SecretNamespace, config.SecretKey, err)
	}

	// The plaintext value is part of the user log on purpose: it is the only
	// observable proof that the decrypted secret actually reached the guest.
	runtime.Logger().Info(fmt.Sprintf("Vault secret fetched: %s", secret.Value))
	return secret.Value, nil
}
