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
	checks := cfg.EffectiveChecks()
	if len(checks) == 0 {
		return "", fmt.Errorf("no vault workflow checks configured")
	}

	for _, check := range checks {
		runtime.Logger().Info("Vault secret workflow triggered",
			"checkName", check.Name,
			"secretKey", check.SecretKey,
			"secretNamespace", check.SecretNamespace,
			"expectNotFound", check.ExpectNotFound,
		)

		secret, err := runtime.GetSecret(&cre.SecretRequest{
			Namespace: check.SecretNamespace,
			Id:        check.SecretKey,
		}).Await()

		if check.ExpectNotFound {
			if err != nil && strings.Contains(err.Error(), "key does not exist") {
				runtime.Logger().Info("Vault secret correctly not found after deletion",
					"checkName", check.Name,
					"secretKey", check.SecretKey,
				)
				continue
			}
			if err != nil {
				runtime.Logger().Error("Expected 'key does not exist' but got a different error",
					"checkName", check.Name,
					"error", err,
					"secretKey", check.SecretKey,
				)
				return "", fmt.Errorf("expected 'key does not exist' for key=%s, but got: %w", check.SecretKey, err)
			}
			runtime.Logger().Error("Expected secret to be gone but retrieval succeeded",
				"checkName", check.Name,
				"secretKey", check.SecretKey,
			)
			return "", fmt.Errorf("expected secret key=%s to be deleted, but it was still found", check.SecretKey)
		}

		if err != nil {
			runtime.Logger().Error("Failed to get secret via workflow",
				"checkName", check.Name,
				"error", err,
			)
			return "", fmt.Errorf("failed to get secret: %w", err)
		}

		if secret.Value == "" {
			runtime.Logger().Error("Secret value is empty",
				"checkName", check.Name,
				"secretKey", check.SecretKey,
				"secretNamespace", check.SecretNamespace,
			)
			return "", fmt.Errorf("secret value is empty for key=%s namespace=%s", check.SecretKey, check.SecretNamespace)
		}

		runtime.Logger().Info("Vault secret retrieved successfully via workflow",
			"checkName", check.Name,
			"secretKey", check.SecretKey,
		)
	}

	runtime.Logger().Info("Vault secret workflow batch completed", "checkCount", len(checks))
	return fmt.Sprintf("Validated %d secret checks", len(checks)), nil
}
