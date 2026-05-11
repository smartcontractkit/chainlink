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
	if cfg.ExpectInvalidIdentifier {
		return evaluateInvalidIdentifiers(cfg, runtime)
	}

	phases := cfg.EffectivePhases()
	if len(phases) == 0 {
		return "", fmt.Errorf("no vault workflow phases configured")
	}

	var lastErr error
	for _, phase := range phases {
		if err := evaluatePhase(runtime, phase); err != nil {
			lastErr = err
			runtime.Logger().Warn("Vault secret workflow phase not yet satisfied",
				"phaseName", phase.Name,
				"error", err,
			)
			continue
		}

		runtime.Logger().Info(fmt.Sprintf("Vault secret workflow phase completed: %s", phase.Name),
			"phaseName", phase.Name,
			"checkCount", len(phase.Checks),
		)
		return fmt.Sprintf("Validated phase %s", phase.Name), nil
	}

	return "", fmt.Errorf("no vault workflow phase matched current state: %w", lastErr)
}

func evaluateInvalidIdentifiers(cfg config.Config, runtime cre.Runtime) (string, error) {
	_, err := runtime.GetSecret(&cre.SecretRequest{
		Namespace: cfg.SecretNamespace,
		Id:        cfg.SecretKey,
	}).Await()
	if err == nil {
		runtime.Logger().Error("Expected identifier validation to fail but GetSecret succeeded", "secretKey", cfg.SecretKey)
		return "", fmt.Errorf("expected identifier validation failure for key=%s, but secret was retrieved", cfg.SecretKey)
	}
	runtime.Logger().Info("Vault get correctly rejected invalid identifier", "secretKey", cfg.SecretKey, "error", err)

	if cfg.SecretKey2 != "" || cfg.SecretNamespace2 != "" {
		key2 := cfg.SecretKey2
		if key2 == "" {
			key2 = cfg.SecretKey
		}
		ns2 := cfg.SecretNamespace2
		if ns2 == "" {
			ns2 = cfg.SecretNamespace
		}
		_, err2 := runtime.GetSecret(&cre.SecretRequest{
			Namespace: ns2,
			Id:        key2,
		}).Await()
		if err2 == nil {
			runtime.Logger().Error("Expected identifier validation to fail for secondary identifier but GetSecret succeeded",
				"secretKey2", key2, "secretNamespace2", ns2)
			return "", fmt.Errorf("expected identifier validation failure for key=%s namespace=%s, but secret was retrieved", key2, ns2)
		}
		runtime.Logger().Info("Vault get correctly rejected invalid identifier", "secretKey", key2, "error", err2)
	}

	return fmt.Sprintf("Invalid identifier correctly rejected: key=%s", cfg.SecretKey), nil
}

func evaluatePhase(runtime cre.Runtime, phase config.Phase) error {
	if len(phase.Checks) == 0 {
		return fmt.Errorf("phase %s has no checks", phase.Name)
	}

	for _, check := range phase.Checks {
		runtime.Logger().Info("Vault secret workflow triggered",
			"phaseName", phase.Name,
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
					"phaseName", phase.Name,
					"checkName", check.Name,
					"secretKey", check.SecretKey,
				)
				continue
			}
			if err != nil {
				return fmt.Errorf("phase %s check %s expected not found for key=%s but got: %w", phase.Name, check.Name, check.SecretKey, err)
			}
			return fmt.Errorf("phase %s check %s expected deleted secret key=%s, but it was still found", phase.Name, check.Name, check.SecretKey)
		}

		if err != nil {
			return fmt.Errorf("phase %s check %s failed to get secret: %w", phase.Name, check.Name, err)
		}

		if secret.Value == "" {
			return fmt.Errorf("phase %s check %s secret value is empty for key=%s namespace=%s", phase.Name, check.Name, check.SecretKey, check.SecretNamespace)
		}

		if check.ExpectedValue != "" && secret.Value != check.ExpectedValue {
			return fmt.Errorf("phase %s check %s secret value mismatch for key=%s namespace=%s", phase.Name, check.Name, check.SecretKey, check.SecretNamespace)
		}

		runtime.Logger().Info("Vault secret retrieved successfully via workflow",
			"phaseName", phase.Name,
			"checkName", check.Name,
			"secretKey", check.SecretKey,
		)
	}

	return nil
}
