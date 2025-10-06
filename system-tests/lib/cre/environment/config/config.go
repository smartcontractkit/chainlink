package config

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	billingplatformservice "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/billing_platform_service"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Config struct {
	Blockchains       []blockchain.Input              `toml:"blockchains" validate:"required"`
	NodeSets          []*cre.CapabilitiesAwareNodeSet `toml:"nodesets" validate:"required"`
	JD                *jd.Input                       `toml:"jd" validate:"required"`
	Infra             *infra.Provider                 `toml:"infra" validate:"required"`
	Fake              *fake.Input                     `toml:"fake" validate:"required"`
	S3ProviderInput   *s3provider.Input               `toml:"s3provider"`
	CapabilityConfigs map[string]cre.CapabilityConfig `toml:"capability_configs"` // capability flag -> capability config

	mu     sync.Mutex
	loaded bool
}

// Validate performs validation checks on the configuration, ensuring all required fields
// are present and all referenced capabilities are known to the system.
func (c *Config) Validate(envDependencies cre.CLIEnvironmentDependencies) error {
	if c.JD.CSAEncryptionKey == "" {
		return errors.New("jd.csa_encryption_key must be provided")
	}

	if len(c.Blockchains) == 0 {
		return errors.New("at least one blockchain must be configured")
	}

	if len(c.NodeSets) == 0 {
		return errors.New("at least one nodeset must be configured")
	}

	if c.Infra == nil {
		return errors.New("infra configuration must be provided")
	}

	for _, nodeSet := range c.NodeSets {
		for _, capability := range nodeSet.Capabilities {
			if !slices.Contains(envDependencies.GlobalCapabilityFlags(), capability) {
				return errors.New("unknown global capability: " + capability + ". Valid ones are: " + strings.Join(envDependencies.GlobalCapabilityFlags(), ", ") + ". If it is a new capability make sure you have added it to the capabilityFlagsProvider. If it's chain-specific add it under [nodesets.chain_capabilities] TOML table.")
			}
		}

		for capability := range nodeSet.ChainCapabilities {
			if !slices.Contains(envDependencies.ChainSpecificCapabilityFlags(), capability) {
				return errors.New("unknown chain-specific capability: " + capability + ". Valid ones are: " + strings.Join(envDependencies.ChainSpecificCapabilityFlags(), ", ") + ". If it is a new capability make sure you have added it to the capabilityFlagsProvider. If it's a global capability add it under 'capabilities' TOML key.")
			}
		}
	}

	if err := validateContractVersions(envDependencies); err != nil {
		return fmt.Errorf("failed to validate initial contract set: %w", err)
	}

	return nil
}

func validateContractVersions(envDependencies cre.CLIEnvironmentDependencies) error {
	supportedSet := DefaultContractSet(envDependencies.WithV2Registries())
	cv := envDependencies.ContractVersions()
	for k, v := range supportedSet {
		version, ok := cv[k]
		if !ok {
			return fmt.Errorf("required contract %s not configured for deployment", k)
		}

		if version != v {
			return fmt.Errorf("requested version %s for contract %s yet expected %s", version, k, v)
		}
	}
	return nil
}

const (
	WorkflowRegistryV2Semver   = "2.0.0"
	CapabilityRegistryV2Semver = "2.0.0"
	DefaultDONFamily           = "test-don-family"
)

func DefaultContractSet(withV2Registries bool) map[string]string {
	supportedSet := map[string]string{
		keystone_changeset.OCR3Capability.String():       "1.0.0",
		keystone_changeset.WorkflowRegistry.String():     "1.0.0",
		keystone_changeset.CapabilitiesRegistry.String(): "1.1.0",
		keystone_changeset.KeystoneForwarder.String():    "1.0.0",
	}

	if withV2Registries {
		supportedSet[keystone_changeset.WorkflowRegistry.String()] = WorkflowRegistryV2Semver
		supportedSet[keystone_changeset.CapabilitiesRegistry.String()] = CapabilityRegistryV2Semver
	}

	return supportedSet
}

func (c *Config) Load(absPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.loaded {
		return nil
	}

	previousCTFconfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		_ = os.Setenv("CTF_CONFIGS", previousCTFconfigs)
	}()

	_ = os.Setenv("CTF_CONFIGS", absPath)

	in, loadErr := framework.Load[Config](nil)
	if loadErr != nil {
		return errors.Wrap(loadErr, "failed to load environment configuration")
	}

	for _, nodeSet := range in.NodeSets {
		if err := nodeSet.ParseChainCapabilities(); err != nil {
			return errors.Wrap(err, "failed to parse chain capabilities")
		}

		if err := nodeSet.ValidateChainCapabilities(in.Blockchains); err != nil {
			return errors.Wrap(err, "failed to validate chain capabilities")
		}
	}

	copyExportedFields(c, in)
	c.loaded = true

	return nil
}

const (
	StateDirname          = "core/scripts/cre/environment/state"
	LocalCREStateFilename = "local_cre.toml"
)

func (c *Config) Store(absPath string) error {
	// change override mode to "each" for all node sets, because config contains unique secrets for each node
	// if we later load it with "all" mode, all nodes in the nodeset will have the same configuration as the first node
	// and they will fail to start (because they will all have the same P2P keys)
	for idx, nodeSet := range c.NodeSets {
		if nodeSet.OverrideMode == "all" {
			c.NodeSets[idx].OverrideMode = "each"
		}
	}

	framework.L.Info().Msgf("Storing local CRE state file: %s", absPath)
	return storeLocalArtifact(c, absPath)
}

func MustLocalCREStateFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, StateDirname, LocalCREStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for local CRE state file: %w", err))
	}

	return absPath
}

func LocalCREStateFileExists(relativePathToRepoRoot string) bool {
	_, statErr := os.Stat(MustLocalCREStateFileAbsPath(relativePathToRepoRoot))
	return statErr == nil
}

// ResolveCapabilityForChain merges defaults with chain override for a capability on a given chain.
// Returns (enabled, mergedConfig).
func ResolveCapabilityForChain(
	capName string,
	caps map[string]*cre.ChainCapabilityConfig,
	defaults map[string]any,
	chainID uint64,
) (bool, map[string]any, error) {
	if caps == nil {
		return false, nil, nil
	}
	cfg, ok := caps[capName]
	if !ok {
		return false, nil, nil
	}
	enabled := slices.Contains(cfg.EnabledChains, chainID)
	if !enabled {
		return false, nil, nil
	}
	merged := map[string]any{}
	if defaults != nil {
		// copy defaults
		maps.Copy(merged, defaults)
	}
	if co, ok := cfg.ChainOverrides[chainID]; ok {
		// override with chain-specific values
		maps.Copy(merged, co)
	}
	return true, merged, nil
}

// ResolveCapabilityConfigForDON merges global defaults with DON-specific overrides for capabilities
// that don't have chain-specific configuration (like cron, web-api-target, web-api-trigger).
// Returns the merged configuration.
func ResolveCapabilityConfigForDON(
	capabilityName string,
	globalDefaults map[string]any,
	donOverrides map[string]map[string]any,
) map[string]any {
	merged := map[string]any{}

	// Start with global defaults
	if globalDefaults != nil {
		maps.Copy(merged, globalDefaults)
	}

	// Apply DON-specific overrides
	if donOverrides != nil {
		if overrides, ok := donOverrides[capabilityName]; ok {
			maps.Copy(merged, overrides)
		}
	}

	return merged
}

type ChipIngressConfig struct {
	ChipIngress *chipingressset.Input `toml:"chip_ingress"`
	Kafka       *KafkaConfig          `toml:"kafka"`

	mu     sync.Mutex
	loaded bool
}

type KafkaConfig struct {
	Topics []string `toml:"topics"`
}

func (c *ChipIngressConfig) Load(absPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.loaded {
		return nil
	}

	previousCTFconfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", previousCTFconfigs)
		if setErr != nil {
			framework.L.Warn().Err(setErr).Msg("failed to restore previous CTF_CONFIGS env var")
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", absPath)
	if setErr != nil {
		return errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
	}

	in, err := framework.Load[ChipIngressConfig](nil)
	if err != nil {
		return errors.Wrap(err, "failed to load chip ingress config")
	}

	copyExportedFields(c, in)
	c.loaded = true

	return nil
}

const (
	ChipIngressStateFilename      = "chip_ingress.toml"
	BillingStateFilename          = "billing-platform-service.toml"
	WorkflowRegistryStateFilename = "workflow_registry.toml"
)

func (c *ChipIngressConfig) Store(absPath string) error {
	framework.L.Info().Msgf("Storing Chip Ingress state file: %s", absPath)
	return storeLocalArtifact(c, absPath)
}

func MustChipIngressStateFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, StateDirname, ChipIngressStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for local CRE state file: %w", err))
	}

	return absPath
}

func ChipIngressStateFileExists(relativePathToRepoRoot string) bool {
	_, statErr := os.Stat(MustChipIngressStateFileAbsPath(relativePathToRepoRoot))
	return statErr == nil
}

func storeLocalArtifact(artifact any, absPath string) error {
	dErr := os.MkdirAll(filepath.Dir(absPath), 0o755)
	if dErr != nil {
		return errors.Wrap(dErr, "failed to create directory for the environment artifact")
	}

	d, mErr := toml.Marshal(artifact)
	if mErr != nil {
		return errors.Wrap(mErr, "failed to marshal environment artifact to TOML")
	}

	return os.WriteFile(absPath, d, 0o600)
}

func RemoveAllEnvironmentStateDir(relativePathToRepoRoot string) error {
	framework.L.Info().Msgf("Removing environment state directory: %s", StateDirname)
	return os.RemoveAll(filepath.Join(relativePathToRepoRoot, StateDirname))
}

// copyExportedFields copies all exported fields from src to dst (same concrete type).
// Unexported fields (like once/mu/loaded) are skipped automatically.
func copyExportedFields(dst, src any) {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src).Elem()
	dt := dv.Type()

	for i := 0; i < dt.NumField(); i++ {
		f := dt.Field(i)
		if f.PkgPath != "" { // unexported
			continue
		}
		dv.Field(i).Set(sv.Field(i))
	}
}

type BillingConfig struct {
	BillingService *billingplatformservice.Input `toml:"billing_platform_service"`

	mu     sync.Mutex
	loaded bool
}

func (c *BillingConfig) Store(absPath string) error {
	framework.L.Info().Msgf("Storing Billing state file: %s", absPath)
	return storeLocalArtifact(c, absPath)
}

func (c *BillingConfig) Load(absPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.loaded {
		return nil
	}

	previousCTFconfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", previousCTFconfigs)
		if setErr != nil {
			framework.L.Warn().Err(setErr).Msg("failed to restore previous CTF_CONFIGS env var")
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", absPath)
	if setErr != nil {
		return errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
	}

	in, err := framework.Load[BillingConfig](nil)
	if err != nil {
		return errors.Wrap(err, "failed to load billing config")
	}

	copyExportedFields(c, in)
	c.loaded = true

	return nil
}

func MustBillingStateFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, StateDirname, BillingStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for local CRE state file: %w", err))
	}

	return absPath
}

func BillingStateFileExists(relativePathToRepoRoot string) bool {
	_, statErr := os.Stat(MustBillingStateFileAbsPath(relativePathToRepoRoot))
	return statErr == nil
}

func MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, StateDirname, WorkflowRegistryStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for local CRE state file: %w", err))
	}

	return absPath
}
