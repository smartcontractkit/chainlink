//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecret/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunVaultSecretWorkflow)
}

func RunVaultSecretWorkflow(cfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onTrigger,
		),
	}, nil
}

func onTrigger(cfg config.Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	runtime.Logger().Info("Vault secret workflow triggered",
		"secretKey", cfg.SecretKey,
		"secretNamespace", cfg.SecretNamespace,
		"expectNotFound", cfg.ExpectNotFound,
	)

	secret, err := runtime.GetSecret(&cre.SecretRequest{
		Namespace: cfg.SecretNamespace,
		Id:        cfg.SecretKey,
	}).Await()

	if cfg.ExpectNotFound {
		if err != nil && strings.Contains(err.Error(), "key does not exist") {
			runtime.Logger().Info("Vault secret correctly not found after deletion", "secretKey", cfg.SecretKey)
			return fmt.Sprintf("Secret correctly not found: key=%s", cfg.SecretKey), nil
		}
		if err != nil {
			runtime.Logger().Error("Expected 'key does not exist' but got a different error",
				"error", err, "secretKey", cfg.SecretKey)
			return "", fmt.Errorf("expected 'key does not exist' for key=%s, but got: %w", cfg.SecretKey, err)
		}
		runtime.Logger().Error("Expected secret to be gone but retrieval succeeded", "secretKey", cfg.SecretKey)
		return "", fmt.Errorf("expected secret key=%s to be deleted, but it was still found", cfg.SecretKey)
	}

	if err != nil {
		runtime.Logger().Error("Failed to get secret via workflow", "error", err)
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	if secret.Value == "" {
		runtime.Logger().Error("Secret value is empty")
		return "", fmt.Errorf("secret value is empty for key=%s namespace=%s", cfg.SecretKey, cfg.SecretNamespace)
	}

	runtime.Logger().Info("Vault secret retrieved successfully via workflow", "secretKey", cfg.SecretKey)
	return fmt.Sprintf("Secret retrieved: key=%s", cfg.SecretKey), nil
}
