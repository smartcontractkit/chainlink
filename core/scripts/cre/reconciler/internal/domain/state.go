package domain

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

// Phase represents a stage in the reconcile lifecycle.
type Phase string

const (
	PhaseNone          Phase = ""
	PhaseDiscovery     Phase = "discovery"      // node runtime collected from K8s + node API
	PhaseProvisioned   Phase = "provisioned"    // contracts deployed and configured
	PhaseConfigWritten Phase = "config_written" // node config patch written and applied
	PhaseJobsSynced    Phase = "jobs_synced"    // jobs proposed + approved + PostEnvStartup run
	PhaseDone          Phase = "done"           // converged

	// Legacy phase names for backwards compatibility (old states still parse)
	PhaseOnChain Phase = "on-chain"
	PhaseTOML    Phase = "toml"
	PhaseReroll  Phase = "reroll"
	PhaseJobs    Phase = "jobs"
)

// Lifecycle: discovery → provisioned → config_written → jobs_synced → done

// StateFile is persisted between runs to enable idempotent reconciliation.
// It is committed to the consumer repo so multiple machines/CI can share state.
type StateFile struct {
	Phase                 Phase                       `toml:"phase"`
	TOMLPatchApplied      bool                        `toml:"toml_patch_applied"`
	Addresses             []AddressRef                `toml:"addresses"`
	DONIDs                map[string]uint64           `toml:"don_ids"`
	JDNodeIDs             map[string]string           `toml:"jd_node_ids"`
	WorkflowReg           *WorkflowRegState           `toml:"workflow_registry"`
	NodeRuntime           map[string]NodeRuntimeInfo  `toml:"node_runtime"`
	GatewayConnectors     []GatewayConnectorState     `toml:"gateway_connectors"`
	GatewayServiceConfigs []GatewayServiceConfigState `toml:"gateway_service_configs"`
	NodeConfigFiles       map[string]string           `toml:"node_config_files"` // key: namespace/nodeName
}

// GatewayConnectorState stores per-gateway connector info from topology for TOML generation.
type GatewayConnectorState struct {
	NodeUUID      string `toml:"node_uuid"`
	AuthGatewayID string `toml:"auth_gateway_id"`
	WebSocketURL  string `toml:"web_socket_url"`
	GatewayDonID  string `toml:"gateway_don_id"`
	DONName       string `toml:"don_name"`
}

// GatewayServiceConfigState persists a topology GatewayServiceConfig so a resumed
// or two-invocation run can restore the handler set that Features' PreEnvStartup
// added during the on-chain phase, instead of rebuilding a topology that carries
// only the default web-api-capabilities handler.
type GatewayServiceConfigState struct {
	ServiceName string                    `toml:"service_name"`
	Handlers    []string                  `toml:"handlers"`
	DONs        []string                  `toml:"dons"`
	Auth0       *GatewayServiceAuth0State `toml:"auth0,omitempty"`
}

// GatewayServiceAuth0State mirrors cre.GatewayServiceAuth0Config for persistence.
type GatewayServiceAuth0State struct {
	IssuerURL string `toml:"issuer_url"`
	Audience  string `toml:"audience"`
	TenantID  uint64 `toml:"tenant_id"` // tenant IDs are small; no uint64>MaxInt64 TOML issue
}

// ChainSelector is a uint64 chain selector. TOML integers are signed 64-bit
// (github.com/pelletier/go-toml/v2 refuses to encode a uint64 above
// math.MaxInt64), but real chain selectors routinely exceed that — so this
// type round-trips through TOML as a decimal string instead of a native
// integer.
type ChainSelector uint64

func (c ChainSelector) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(c), 10)), nil
}

func (c *ChainSelector) UnmarshalText(text []byte) error {
	v, err := strconv.ParseUint(string(text), 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalid chain selector %q", text)
	}
	*c = ChainSelector(v)
	return nil
}

// AddressRef is a deployed contract reference, compatible with CLDF datastore.
type AddressRef struct {
	ChainSelector ChainSelector `toml:"chain_selector"`
	Address       string        `toml:"address"`
	Type          string        `toml:"type"`
	Version       string        `toml:"version"`
	Qualifier     string        `toml:"qualifier"`
}

// WorkflowRegState stores the Workflow Registry configuration.
type WorkflowRegState struct {
	UseCache       bool          `toml:"use_cache"`
	ChainSelector  ChainSelector `toml:"chain_selector"`
	AllowedDonIDs  []uint32      `toml:"allowed_don_ids"`
	WorkflowOwners []string      `toml:"workflow_owners"`
}

// NodeRuntimeInfo caches per-node runtime data discovered from K8s/JD.
type NodeRuntimeInfo struct {
	PeerID        string            `toml:"peer_id"`
	APIURL        string            `toml:"api_url"`
	CSAKey        string            `toml:"csa_key"`
	EVMAddress    map[string]string `toml:"evm_addresses"`
	NodeType      string            `toml:"node_type"`
	OCR2BundleIDs map[string]string `toml:"ocr2_bundle_ids"` // chain family (lowercase, e.g. "evm") -> bundle ID
}

// LoadState reads a state file from disk. Returns nil, nil if the file
// does not exist (fresh run).
func LoadState(path string) (*StateFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to read state file %s", path)
	}

	var s StateFile
	if err := toml.Unmarshal(data, &s); err != nil {
		return nil, errors.Wrapf(err, "failed to parse state TOML %s", path)
	}

	return &s, nil
}

// Store writes the state file atomically.
func (s *StateFile) Store(path string) error {
	data, err := toml.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "failed to marshal state")
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create state dir %s", dir)
		}
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return errors.Wrapf(err, "failed to write state file %s", tmp)
	}

	if err := os.Rename(tmp, path); err != nil {
		return errors.Wrapf(err, "failed to rename state file %s -> %s", tmp, path)
	}

	return nil
}

// GetAddress returns the address for a contract type, or empty if not found.
func (s *StateFile) GetAddress(contractType string) string {
	for _, ref := range s.Addresses {
		if ref.Type == contractType {
			return ref.Address
		}
	}
	return ""
}

// SetAddress adds or updates an address reference.
func (s *StateFile) SetAddress(ref AddressRef) {
	for i, existing := range s.Addresses {
		if existing.Type == ref.Type && existing.ChainSelector == ref.ChainSelector && existing.Qualifier == ref.Qualifier {
			s.Addresses[i] = ref
			return
		}
	}
	s.Addresses = append(s.Addresses, ref)
}

// HasAddress returns true if a contract of the given type is in the state.
func (s *StateFile) HasAddress(contractType string) bool {
	return s.GetAddress(contractType) != ""
}
