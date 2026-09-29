//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	httptrigger "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
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

// RunVaultSecretWorkflow exposes an HTTP-triggered vault secret verification workflow.
// Verification parameters arrive per invocation in the trigger input (see
// config.TriggerInput), so tests request verification exactly when a state
// change is expected instead of waiting for the next cron tick (the previous
// cron-trigger variant was gated on its 30s minimum schedule).
func RunVaultSecretWorkflow(cfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			httptrigger.Trigger(&httptrigger.Config{
				AuthorizedKeys: []*httptrigger.AuthorizedKey{
					{
						Type:      httptrigger.KeyType_KEY_TYPE_ECDSA_EVM,
						PublicKey: cfg.AuthorizedKey,
					},
				},
			}),
			onTrigger,
		),
	}, nil
}

func onTrigger(cfg config.Config, runtime cre.Runtime, trigger *httptrigger.Payload) (string, error) {
	var input config.TriggerInput
	if len(trigger.Input) > 0 {
		if err := json.Unmarshal(trigger.Input, &input); err != nil {
			return "", fmt.Errorf("failed to unmarshal trigger input: %w", err)
		}
	}
	if input.ExpectInvalidIdentifier {
		return evaluateInvalidIdentifiers(input, runtime)
	}
	if input.ExpectBatchTooBig {
		return evaluateBatchTooBig(input, runtime)
	}

	phases := input.Phases
	if len(phases) == 0 {
		return "", fmt.Errorf("no vault workflow phases configured")
	}

	// Phases run in declaration order; the FIRST phase whose checks all succeed emits the completion log
	// and returns. Later phases whose checks describe the SAME success state never run—you need a strictly
	// later world state than earlier phases so that earlier phases eventually fail once the vault changes.
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

func evaluateInvalidIdentifiers(input config.TriggerInput, runtime cre.Runtime) (string, error) {
	_, err := runtime.GetSecret(&cre.SecretRequest{
		Namespace: input.SecretNamespace,
		Id:        input.SecretKey,
	}).Await()
	if err == nil {
		runtime.Logger().Error("Expected identifier validation to fail but GetSecret succeeded", "secretKey", input.SecretKey)
		return "", fmt.Errorf("expected identifier validation failure for key=%s, but secret was retrieved", input.SecretKey)
	}
	runtime.Logger().Info("Vault get correctly rejected invalid identifier", "secretKey", input.SecretKey, "error", err)

	if input.SecretKey2 != "" || input.SecretNamespace2 != "" {
		key2 := input.SecretKey2
		if key2 == "" {
			key2 = input.SecretKey
		}
		ns2 := input.SecretNamespace2
		if ns2 == "" {
			ns2 = input.SecretNamespace
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

	return fmt.Sprintf("Invalid identifier correctly rejected: key=%s", input.SecretKey), nil
}

// evaluateBatchTooBig submits a GetSecrets batch larger than the vault request
// batch size limit and verifies the rejection carries the real user-facing
// cause across the workflow-DON -> vault-DON remote capability hop.
func evaluateBatchTooBig(input config.TriggerInput, runtime cre.Runtime) (string, error) {
	if len(input.BatchSecretKeys) == 0 {
		return "", fmt.Errorf("expectBatchTooBig requires batchSecretKeys to be set")
	}

	reqs := make([]*cre.SecretRequest, len(input.BatchSecretKeys))
	for i, key := range input.BatchSecretKeys {
		reqs[i] = &cre.SecretRequest{
			Namespace: input.SecretNamespace,
			Id:        key,
		}
	}

	_, err := runtime.GetSecrets(reqs).Await()
	if err == nil {
		runtime.Logger().Error("Expected batch size validation to fail but GetSecrets succeeded", "count", len(reqs))
		return "", fmt.Errorf("expected batch size validation failure for %d secrets, but batch was retrieved", len(reqs))
	}
	if !strings.Contains(err.Error(), "request batch size exceeds maximum of") {
		runtime.Logger().Error("GetSecrets failed with an unexpected error", "count", len(reqs), "error", err)
		return "", fmt.Errorf("expected batch size rejection for %d secrets, got: %w", len(reqs), err)
	}

	runtime.Logger().Info("Vault get correctly rejected oversized batch", "count", len(reqs), "error", err)
	return fmt.Sprintf("Oversized batch correctly rejected: %d secrets", len(reqs)), nil
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
