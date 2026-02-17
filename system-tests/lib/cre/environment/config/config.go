package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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

// GetAddresses returns the addresses as datastore.AddressRef slice
func (c *Config) GetAddresses() ([]datastore.AddressRef, error) {
	addresses := make([]datastore.AddressRef, len(c.Addresses))
	for i, addr := range c.Addresses {
		in := []byte(addr)
		var addrRef datastore.AddressRef
		if err := json.Unmarshal(in, &addrRef); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address at index %d: %w", i, err)
		}
		addresses[i] = addrRef
	}
	return addresses, nil
}

// SetAddresses sets the addresses from datastore.AddressRef slice
func (c *Config) SetAddresses(refs []datastore.AddressRef) error {
	c.Addresses = make([]string, len(refs))
	for i, ref := range refs {
		asBytes, err := json.Marshal(ref)
		if err != nil {
			return fmt.Errorf("failed to marshal address at index %d: %w", i, err)
		}
		c.Addresses[i] = string(asBytes)
	}
	return nil
}

type Config struct {
	Blockchains       []*blockchain.Input             `toml:"blockchains" validate:"required"`
	NodeSets          []*cre.NodeSet                  `toml:"nodesets" validate:"required"`
	JD                *jd.Input                       `toml:"jd" validate:"required"`
	Infra             *infra.Provider                 `toml:"infra" validate:"required"`
	Fake              *fake.Input                     `toml:"fake" validate:"required"`
	FakeHTTP          *fake.Input                     `toml:"fake_http" validate:"required"`
	S3ProviderInput   *s3provider.Input               `toml:"s3provider"`
	CapabilityConfigs map[string]cre.CapabilityConfig `toml:"capability_configs"` // capability flag -> capability config
	Addresses         []string                        `toml:"addresses"`

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
			capability = removeChainIDFromFlag(capability)
			if !slices.Contains(envDependencies.SupportedCapabilityFlags(), capability) {
				return errors.New("unknown capability: " + capability + ". Valid ones are: " + strings.Join(envDependencies.SupportedCapabilityFlags(), ", ") + ". If it is a new capability make sure you have added it to the capabilityFlagsProvider")
			}
		}
	}

	if err := validateContractVersions(envDependencies); err != nil {
		return fmt.Errorf("failed to validate initial contract set: %w", err)
	}

	return nil
}

func removeChainIDFromFlag(flag string) string {
	lastIdx := strings.LastIndex(flag, "-")
	if lastIdx == -1 {
		return flag
	}

	maybeChainIDStr := flag[lastIdx+1:]
	_, err := strconv.Atoi(maybeChainIDStr)
	if err != nil {
		return flag
	}

	return flag[:lastIdx]
}

func validateContractVersions(envDependencies cre.CLIEnvironmentDependencies) error {
	supportedSet := DefaultContractSet(envDependencies.WithV2Registries())
	cv := envDependencies.ContractVersions()
	for k, v := range supportedSet {
		version, ok := cv[k]
		if !ok {
			return fmt.Errorf("required contract %s not configured for deployment", k)
		}

		if !version.Equal(v) {
			return fmt.Errorf("requested version %s for contract %s yet expected %s", version, k, v)
		}
	}
	return nil
}

var (
	WorkflowRegistryV2Semver   = semver.MustParse("2.0.0")
	CapabilityRegistryV2Semver = semver.MustParse("2.0.0")
)

const (
	DefaultDONFamily = "test-don-family"
)

func DefaultContractSet(withV2Registries bool) map[cre.ContractType]*semver.Version {
	supportedSet := map[cre.ContractType]*semver.Version{
		keystone_changeset.OCR3Capability.String():       semver.MustParse("1.0.0"),
		keystone_changeset.WorkflowRegistry.String():     semver.MustParse("1.0.0"),
		keystone_changeset.CapabilitiesRegistry.String(): semver.MustParse("1.1.0"),
		keystone_changeset.KeystoneForwarder.String():    semver.MustParse("1.0.0"),
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

		// Clear the embedded Input.NodeSpecs to avoid storing duplicate node specs without roles.
		// We only want to persist NodeSpecs (NodeSpecWithRole[]) which contains role information.
		// The Input.NodeSpecs field is populated at runtime in dons.go for passing to the CTF library.
		c.NodeSets[idx].Input.NodeSpecs = nil
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

	// WORKAROUND: Remove the empty "node_specs = []" line that gets marshaled from the embedded ns.Input
	// This conflicts with our [[nodesets.node_specs]] array tables that include role information.
	// We use regex to remove the problematic line while preserving the actual node_specs tables.
	// TOML library we use doesn't support omitting empty slices nor custom (un)marshalling.
	tomlStr := string(d)
	lines := strings.Split(tomlStr, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		// Skip the "node_specs = []" line but keep [[nodesets.node_specs]] sections
		if strings.TrimSpace(line) == "node_specs = []" {
			continue
		}
		filtered = append(filtered, line)
	}
	d = []byte(strings.Join(filtered, "\n"))

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
