package workflow

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	secretsUtils "github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/vault"
	vaulttypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// vaultSecretsConfig defines the structure of the vault secrets YAML file.
type vaultSecretsConfig struct {
	Secrets []vaultSecretEntry `yaml:"secrets"`
}

// vaultSecretEntry represents a single secret to be stored in the vault.
type vaultSecretEntry struct {
	Key       string `yaml:"key"`       // Vault key identifier
	EnvVar    string `yaml:"envVar"`    // Name of the env var containing the secret value
	Namespace string `yaml:"namespace"` // Vault namespace (defaults to "main" if empty)
}

// secretsNamesConfig is the secrets YAML format used by CRE:
//
//	secretsNames:
//	  SECRET_KEY:
//	    - ENV_VAR_NAME
type secretsNamesConfig struct {
	SecretsNames map[string][]string `yaml:"secretsNames"`
}

func newSecretsConfig(configPath string) (*secretsUtils.SecretsConfig, error) {
	secretsConfigFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening secrets config file: %w", err)
	}
	defer secretsConfigFile.Close()

	var config secretsUtils.SecretsConfig
	err = yaml.NewDecoder(secretsConfigFile).Decode(&config)
	if err != nil && errors.Is(err, io.EOF) {
		return &secretsUtils.SecretsConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error decoding secrets config file: %w", err)
	}

	return &config, nil
}

// PrepareSecrets reads the vault secrets YAML file, encrypts each secret using the vault
// public key, and writes the encrypted secrets list to a JSON file. The JSON file path is returned.
//
// Two YAML formats are accepted:
//
// Format 1 (explicit):
//
//	secrets:
//	  - key: "my-secret"
//	    envVar: "MY_SECRET_ENV_VAR"
//	    namespace: "main"   # optional, defaults to "main"
//
// Format 2 (secretsNames, shared with other CRE tools):
//
//	secretsNames:
//	  SECRET_KEY:
//	    - ENV_VAR_NAME
func PrepareSecrets(secretsFilePath, vaultPublicKey string, ownerAddress common.Address, outputFilePath string) (string, error) {
	data, err := os.ReadFile(secretsFilePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read secrets file")
	}

	var cfg vaultSecretsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", errors.Wrap(err, "failed to parse secrets YAML file")
	}

	if len(cfg.Secrets) == 0 {
		// Try the alternative secretsNames format.
		var altCfg secretsNamesConfig
		if altErr := yaml.Unmarshal(data, &altCfg); altErr == nil {
			for key, envVars := range altCfg.SecretsNames {
				if len(envVars) == 0 {
					continue
				}
				cfg.Secrets = append(cfg.Secrets, vaultSecretEntry{
					Key:    key,
					EnvVar: envVars[0],
				})
			}
		}
	}

	if len(cfg.Secrets) == 0 {
		return "", errors.New("no secrets found in secrets file")
	}

	encryptedSecrets := make([]*vault_helpers.EncryptedSecret, 0, len(cfg.Secrets))
	for _, entry := range cfg.Secrets {
		value := os.Getenv(entry.EnvVar)
		if value == "" {
			return "", fmt.Errorf("environment variable %q is not set for secret key %q", entry.EnvVar, entry.Key)
		}

		namespace := entry.Namespace
		if namespace == "" {
			namespace = "main"
		}

		encryptedValue, encErr := crevault.EncryptSecret(value, vaultPublicKey, ownerAddress)
		if encErr != nil {
			return "", errors.Wrapf(encErr, "failed to encrypt secret %q", entry.Key)
		}

		encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
			Id: &vault_helpers.SecretIdentifier{
				Key:       entry.Key,
				Owner:     ownerAddress.Hex(),
				Namespace: namespace,
			},
			EncryptedValue: encryptedValue,
		})
	}

	if outputFilePath == "" {
		outputFilePath = "./vault_secrets.json"
	}

	absPath, absErr := filepath.Abs(outputFilePath)
	if absErr != nil {
		return "", errors.Wrap(absErr, "failed to resolve absolute path for secrets output file")
	}

	jsonData, marshalErr := json.Marshal(encryptedSecrets)
	if marshalErr != nil {
		return "", errors.Wrap(marshalErr, "failed to marshal encrypted secrets to JSON")
	}

	if writeErr := os.WriteFile(absPath, jsonData, 0600); writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write encrypted secrets file")
	}

	return absPath, nil
}

// ExecuteSecrets reads the encrypted secrets JSON file produced by PrepareSecrets,
// allowlists the vault request in the workflow registry, and sends the secrets to the vault gateway.
func ExecuteSecrets(ctx context.Context, encryptedSecretsJSONPath, gatewayURL string, sethClient *seth.Client, workflowRegistryAddress common.Address) error {
	data, err := os.ReadFile(encryptedSecretsJSONPath)
	if err != nil {
		return errors.Wrap(err, "failed to read encrypted secrets file")
	}

	var encryptedSecrets []*vault_helpers.EncryptedSecret
	if err = json.Unmarshal(data, &encryptedSecrets); err != nil {
		return errors.Wrap(err, "failed to unmarshal encrypted secrets")
	}

	uniqueRequestID := uuid.New().String()
	createSecretsRequest := vault_helpers.CreateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: encryptedSecrets,
	}

	requestBody, err := json.Marshal(&createSecretsRequest)
	if err != nil {
		return errors.Wrap(err, "failed to marshal create secrets request")
	}
	requestBodyJSON := json.RawMessage(requestBody)

	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsCreate,
		Params:  &requestBodyJSON,
	}

	requestDigest, err := jsonRequest.Digest()
	if err != nil {
		return errors.Wrap(err, "failed to compute request digest")
	}

	requestDigestBytes, err := hex.DecodeString(requestDigest)
	if err != nil {
		return errors.Wrap(err, "failed to decode request digest hex")
	}
	if len(requestDigestBytes) != 32 {
		return errors.Errorf("invalid request digest length: got %d bytes, want 32", len(requestDigestBytes))
	}

	var reqDigestBytes [32]byte
	copy(reqDigestBytes[:], requestDigestBytes)

	wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(workflowRegistryAddress, sethClient.Client)
	if err != nil {
		return errors.Wrap(err, "failed to instantiate workflow registry v2 wrapper")
	}

	expiry := uint32(time.Now().Add(time.Hour).Unix()) //nolint:gosec // G115: timestamp fits uint32 until year 2106
	_, decErr := sethClient.Decode(wfReg.AllowlistRequest(sethClient.NewTXOpts(), reqDigestBytes, expiry))
	if decErr != nil {
		return errors.Wrap(decErr, "failed to allowlist vault request in workflow registry")
	}

	fmt.Printf("\n✅ Vault request allowlisted in workflow registry\n")

	reqBody, err := json.Marshal(jsonRequest)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON-RPC request")
	}

	statusCode, respBody, sendErr := cre.SendToVaultGateway(ctx, gatewayURL, reqBody)
	if sendErr != nil {
		return errors.Wrap(sendErr, "failed to send request to vault gateway")
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("vault gateway responded with status %d: %s", statusCode, string(respBody))
	}

	var jsonResponse jsonrpc.Response[json.RawMessage]
	if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
		return errors.Wrap(err, "failed to unmarshal vault gateway response")
	}

	if jsonResponse.Error != nil && jsonResponse.Error.Error() != "" {
		return fmt.Errorf("vault gateway returned error: %s", jsonResponse.Error.Error())
	}

	return nil
}

// FetchVaultPublicKey polls the vault gateway until it returns a public key.
func FetchVaultPublicKey(ctx context.Context, gatewayURL string) (string, error) {
	getPublicKeyRequest := jsonrpc.Request[vault_helpers.GetPublicKeyRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uuid.New().String(),
		Method:  vaulttypes.MethodPublicKeyGet,
		Params:  &vault_helpers.GetPublicKeyRequest{},
	}

	reqBody, err := json.Marshal(getPublicKeyRequest)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal public key request")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeout := time.After(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timeout:
			return "", errors.New("timed out waiting for vault public key from gateway")
		case <-ticker.C:
			statusCode, respBody, sendErr := cre.SendToVaultGateway(ctx, gatewayURL, reqBody)
			if sendErr != nil || statusCode != http.StatusOK || respBody == nil {
				continue
			}

			var jsonResponse jsonrpc.Response[vault_helpers.GetPublicKeyResponse]
			if jsonErr := json.Unmarshal(respBody, &jsonResponse); jsonErr != nil {
				continue
			}

			if jsonResponse.Result.PublicKey != "" {
				return jsonResponse.Result.PublicKey, nil
			}
		}
	}
}
