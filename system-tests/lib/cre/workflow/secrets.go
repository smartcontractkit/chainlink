package workflow

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"

	secretsUtils "github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

func PrepareSecrets(sethClient *seth.Client, donID uint32, capabilitiesRegistryAddress, workflowOwnerAddress common.Address, secretsFilePath string) (string, error) {
	secretsConfig, secretsConfigErr := newSecretsConfig(secretsFilePath)
	if secretsConfigErr != nil {
		return "", errors.Wrap(secretsConfigErr, "failed to parse secrets config")
	}

	envSecrets, envSecretsErr := loadSecretsFromEnvironment(secretsConfig)
	if envSecretsErr != nil {
		return "", errors.Wrap(envSecretsErr, "failed to load secrets from environment")
	}

	encryptSecrets, encryptSecretsErr := encryptSecrets(sethClient, donID, capabilitiesRegistryAddress, workflowOwnerAddress, envSecrets, secretsConfig)
	if encryptSecretsErr != nil {
		return "", errors.Wrap(encryptSecretsErr, "failed to encrypt secrets")
	}

	encryptedSecretsFilePath := "./encrypted.secrets.json"
	encryptedSecretsFile, encryptedSecretsFileErr := os.Create(encryptedSecretsFilePath)
	if encryptedSecretsFileErr != nil {
		return "", errors.Wrap(encryptedSecretsFileErr, "failed to create secrets file")
	}

	defer encryptedSecretsFile.Close()

	encoder := json.NewEncoder(encryptedSecretsFile)
	if encoderErr := encoder.Encode(encryptSecrets); encoderErr != nil {
		return "", errors.Wrap(encoderErr, "failed to write to secrets file")
	}

	secretPathAbs, secretPathAbsErr := filepath.Abs(encryptedSecretsFilePath)
	if secretPathAbsErr != nil {
		return "", errors.Wrap(secretPathAbsErr, "failed to get absolute path of the encrypted secrets file")
	}

	return secretPathAbs, nil
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

func loadSecretsFromEnvironment(config *secretsUtils.SecretsConfig) (map[string][]string, error) {
	secrets := make(map[string][]string)
	for secretName, envVars := range config.SecretsNames {
		for _, envVar := range envVars {
			secretValue := os.Getenv(envVar)
			if secretValue == "" {
				return nil, fmt.Errorf("missing environment variable: %s", envVar)
			}
			secrets[secretName] = append(secrets[secretName], secretValue)
		}
	}
	return secrets, nil
}

func encryptSecrets(c *seth.Client, donID uint32, capabilitiesRegistry, workflowOwner common.Address, secrets map[string][]string, config *secretsUtils.SecretsConfig) (secretsUtils.EncryptedSecretsResult, error) {
	cr, err := capabilities_registry.NewCapabilitiesRegistry(capabilitiesRegistry, c.Client)
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("failed to attach to the Capabilities Registry contract: %w", err)
	}

	nodeInfos, err := cr.GetNodes(c.NewCallOpts())
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("failed to get node information from the Capabilities Registry contract: %w", err)
	}

	donInfo, err := cr.GetDON(c.NewCallOpts(), donID)
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("failed to get DON information from the Capabilities Registry contract: %w", err)
	}

	encryptionPublicKeys := make(map[string][32]byte)
	for _, nodeInfo := range nodeInfos {
		// Filter only the nodes that are part of the DON
		if secretsUtils.ContainsP2pId(nodeInfo.P2pId, donInfo.NodeP2PIds) {
			encryptionPublicKeys[hex.EncodeToString(nodeInfo.P2pId[:])] = nodeInfo.EncryptionPublicKey
		}
	}

	if len(encryptionPublicKeys) == 0 {
		return secretsUtils.EncryptedSecretsResult{}, errors.New("no nodes found for the don")
	}

	// Encrypt secrets for each node
	encryptedSecrets, secretsEnvVarsByNode, err := secretsUtils.EncryptSecretsForNodes(
		workflowOwner.String(),
		secrets,
		encryptionPublicKeys,
		secretsUtils.SecretsConfig{SecretsNames: config.SecretsNames},
	)
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("node public keys not found: %w", err)
	}

	// Convert encryptionPublicKey to hex strings for including in the metadata
	nodePublicEncryptionKeys := make(map[string]string)
	for p2pID, encryptionPublicKey := range encryptionPublicKeys {
		nodePublicEncryptionKeys[p2pID] = hex.EncodeToString(encryptionPublicKey[:])
	}

	result := secretsUtils.EncryptedSecretsResult{
		EncryptedSecrets: encryptedSecrets,
		Metadata: secretsUtils.Metadata{
			WorkflowOwner:            workflowOwner.String(),
			CapabilitiesRegistry:     capabilitiesRegistry.String(),
			DonId:                    strconv.FormatUint(uint64(donID), 10),
			DateEncrypted:            time.Now().Format(time.RFC3339),
			NodePublicEncryptionKeys: nodePublicEncryptionKeys,
			EnvVarsAssignedToNodes:   secretsEnvVarsByNode,
		},
	}
	return result, nil
}
