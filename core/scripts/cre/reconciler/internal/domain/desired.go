// Package domain defines the shared value types for griddle: DesiredState, ChartValues, StateFile, NodeRole, etc.
package domain

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// NodeRole classifies a Chainlink node by its function in the DON.
type NodeRole string

const (
	RoleStandard  NodeRole = "standard"
	RoleBootstrap NodeRole = "boot"
	RoleGateway   NodeRole = "gateway"
)

// DesiredState is the parsed desired-state TOML that the user declares.
// It describes the DONs, their node membership, and capabilities.
type DesiredState struct {
	Infra             Infra                       `toml:"infra"`
	JD                JDConfig                    `toml:"jd"`
	Chains            []Chain                     `toml:"chains"`
	DONs              []DON                       `toml:"dons"`
	GatewayNodes      []GatewayNodeAssignment     `toml:"gateway_nodes"`
	CapabilityConfigs map[string]CapabilityConfig `toml:"capability_configs"`
}

// Chain is a user-declared EVM chain: chain ID plus the RPC URLs used to talk
// to it. Exactly one declared chain must be the registry chain — the chain
// where CapabilitiesRegistry/WorkflowRegistry live and where nodes register.
// This is the only source of chain data the reconciler uses; nothing is
// derived from the chart's anvil.instances.
type Chain struct {
	ChainID  uint64 `toml:"chain_id"`
	WSURL    string `toml:"ws_url"`
	HTTPURL  string `toml:"http_url"`
	Registry bool   `toml:"registry"`
}

// Infra describes the Griddle deployment target.
type Infra struct {
	Type        string `toml:"type"`         // must be "griddle"
	ChartValues string `toml:"chart_values"` // path to deploy/config/<service> dir
	Namespace   string `toml:"namespace"`    // K8s namespace
	Kubeconfig  string `toml:"kubeconfig"`   // optional path to kubeconfig
}

// JDConfig describes how to connect to the Job Distributor.
// The access token is deliberately NOT a field here — it is a live bearer
// credential and must come only from the GRIDDLE_JD_ACCESS_TOKEN env var
// (see internal/infra.JDAccessToken), never from desired.toml or the UI.
type JDConfig struct {
	GRPC        string `toml:"grpc"`        // gRPC endpoint
	WSRPC       string `toml:"wsrpc"`       // optional wsRPC endpoint
	Domain      string `toml:"domain"`      // e.g. "cre"
	Environment string `toml:"environment"` // e.g. "dev"
	UseTLS      bool   `toml:"use_tls"`     // use TLS for gRPC (auto-detected from :443 if false)
}

// DON describes a logical DON in the desired state.
// Node membership is always derived from the chart's don-name label.
type DON struct {
	Name                   string                      `toml:"name"`
	DONTypes               []string                    `toml:"don_types"`
	Capabilities           []string                    `toml:"capabilities"`
	BootstrapNode          string                      `toml:"bootstrap_node"` // optional; defaults to nodeType: boot
	ExposesRemoteCaps      bool                        `toml:"exposes_remote_capabilities"`
	RegistryBasedAllowlist []string                    `toml:"registry_based_launch_allowlist"`
	CapabilityConfigs      map[string]CapabilityConfig `toml:"capability_configs"`
}

// ResolveBootstrap returns the bootstrap node name for this DON.
// Priority: explicit BootstrapNode field → chart node with nodeType: boot
// in this DON's node-set → node matching *-bt-* naming convention.
// Returns empty string if no bootstrap is found.
func (d *DON) ResolveBootstrap(cv *ChartValues) string {
	if d.BootstrapNode != "" {
		return d.BootstrapNode
	}
	nodeNames := cv.NodeNamesForDONName(d.Name)
	for _, nodeName := range nodeNames {
		if n := cv.GetNode(nodeName); n != nil && n.NodeType == RoleBootstrap {
			return nodeName
		}
	}
	for _, nodeName := range nodeNames {
		lower := strings.ToLower(nodeName)
		if strings.Contains(lower, "-bt-") || strings.Contains(lower, "bootstrap") {
			return nodeName
		}
	}
	return ""
}

// WorkerNodes returns the nodes in this DON that are NOT the bootstrap.
// Bootstrap nodes are excluded from on-chain signatory registration.
func (d *DON) WorkerNodes(cv *ChartValues) []string {
	bootstrap := d.ResolveBootstrap(cv)
	nodeNames := cv.NodeNamesForDONName(d.Name)
	var workers []string
	for _, nodeName := range nodeNames {
		if nodeName == bootstrap {
			continue
		}
		if n := cv.GetNode(nodeName); n != nil {
			if n.NodeType == RoleBootstrap || n.NodeType == RoleGateway {
				continue
			}
		}
		workers = append(workers, nodeName)
	}
	return workers
}

// IsBootstrapOnly returns true if this DON is a bootstrap-type DON with no
// worker nodes. Gateway-only DONs are never bootstrap-only even when they have
// zero workers (gateway nodes are excluded from WorkerNodes).
func (d *DON) IsBootstrapOnly(cv *ChartValues) bool {
	return d.IsBootstrapDon() && len(d.WorkerNodes(cv)) == 0
}

// HasDONType reports whether the DON has the given don_types entry.
func (d *DON) HasDONType(donType string) bool {
	return slices.Contains(d.DONTypes, donType)
}

// IsWorkflowDon reports whether this DON is a workflow DON.
func (d *DON) IsWorkflowDon() bool {
	return d.HasDONType(cre.WorkflowDON)
}

// IsGatewayDon reports whether this DON is a gateway DON.
func (d *DON) IsGatewayDon() bool {
	return d.HasDONType(cre.GatewayDON)
}

// IsBootstrapDon reports whether this DON is a bootstrap DON.
func (d *DON) IsBootstrapDon() bool {
	return d.HasDONType(cre.BootstrapDON)
}

// NeedsGatewayAccess reports whether this DON has gateway-routed capabilities.
func (d *DON) NeedsGatewayAccess() bool {
	gatewayCaps := map[string]bool{
		cre.VaultCapability:       true,
		cre.HTTPActionCapability:  true,
		cre.HTTPTriggerCapability: true,
	}
	for _, cap := range d.Capabilities {
		if gatewayCaps[cap] {
			return true
		}
	}
	return false
}

// GatewayNodeAssignment optionally maps a gateway node to a specific DON.
// If not provided, the gateway node serves the first workflow DON.
type GatewayNodeAssignment struct {
	Node string `toml:"node"`
	DON  string `toml:"don"`
}

// CapabilityConfig mirrors cre.CapabilityConfig for TOML parsing.
type CapabilityConfig struct {
	BinaryName string         `toml:"binary_name"`
	Values     map[string]any `toml:"values"`
}

// LoadDesiredState reads and parses a desired-state TOML file.
func LoadDesiredState(path string) (*DesiredState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read desired state file %s", path)
	}

	var ds DesiredState
	if err := toml.Unmarshal(data, &ds); err != nil {
		return nil, errors.Wrapf(err, "failed to parse desired state TOML %s", path)
	}

	if err := ds.Validate(); err != nil {
		return nil, errors.Wrapf(err, "invalid desired state in %s", path)
	}

	return &ds, nil
}

// Validate checks the desired state for obvious errors.
// Cross-referencing node names against chart values happens later in discovery.
func (ds *DesiredState) Validate() error {
	if ds.Infra.Type != "griddle" {
		return fmt.Errorf("infra.type must be \"griddle\", got %q", ds.Infra.Type)
	}
	if ds.Infra.ChartValues == "" {
		return errors.New("infra.chart_values is required (path to deploy/config/<service> dir)")
	}
	if ds.Infra.Namespace == "" {
		return errors.New("infra.namespace is required")
	}
	if ds.JD.GRPC == "" {
		return errors.New("jd.grpc is required (gRPC endpoint)")
	}
	if ds.JD.Domain == "" {
		return errors.New("jd.domain is required")
	}
	if ds.JD.Environment == "" {
		return errors.New("jd.environment is required")
	}
	if len(ds.DONs) == 0 {
		return errors.New("at least one [[dons]] entry is required")
	}

	seenNames := make(map[string]bool)
	for i, don := range ds.DONs {
		if don.Name == "" {
			return fmt.Errorf("dons[%d].name is required", i)
		}
		if seenNames[don.Name] {
			return fmt.Errorf("duplicate DON name %q", don.Name)
		}
		seenNames[don.Name] = true

		// Check that every capability has a capability_configs entry
		// (either global or per-DON). Some capabilities (e.g. don-time,
		// consensus) may not have any config — that's valid.
		for _, cap := range don.Capabilities {
			if !ds.hasCapabilityConfig(cap, don) {
				// Not all capabilities require config. These are built-in
				// or auto-configured and don't need explicit config entries.
				if isConfiglessCapability(cap) {
					continue
				}
				return fmt.Errorf("dons[%d] (%s): capability %q has no capability_configs entry (neither global nor per-DON)", i, don.Name, cap)
			}
		}
	}

	return ds.validateChains()
}

// validateChains enforces the explicit chain model: at least one declared
// chain, exactly one registry chain, unique chain IDs, valid RPC URLs, and
// every chain-scoped capability referencing a declared chain.
func (ds *DesiredState) validateChains() error {
	if len(ds.Chains) == 0 {
		return errors.New("at least one [[chains]] entry is required")
	}

	seen := make(map[uint64]bool, len(ds.Chains))
	registryCount := 0
	for i, ch := range ds.Chains {
		if ch.ChainID == 0 {
			return fmt.Errorf("chains[%d]: chain_id is required", i)
		}
		if seen[ch.ChainID] {
			return fmt.Errorf("duplicate chain_id %d in [[chains]]", ch.ChainID)
		}
		seen[ch.ChainID] = true
		if ch.Registry {
			registryCount++
		}
		if ch.WSURL == "" {
			return fmt.Errorf("chains[%d] (chain_id %d): ws_url is required", i, ch.ChainID)
		}
		if _, err := commonconfig.ParseURL(ch.WSURL); err != nil {
			return fmt.Errorf("chains[%d] (chain_id %d): ws_url %q is not a valid URL: %w", i, ch.ChainID, ch.WSURL, err)
		}
		if ch.HTTPURL == "" {
			return fmt.Errorf("chains[%d] (chain_id %d): http_url is required", i, ch.ChainID)
		}
		if _, err := commonconfig.ParseURL(ch.HTTPURL); err != nil {
			return fmt.Errorf("chains[%d] (chain_id %d): http_url %q is not a valid URL: %w", i, ch.ChainID, ch.HTTPURL, err)
		}
	}
	if registryCount != 1 {
		return fmt.Errorf("exactly one [[chains]] entry must set registry = true, found %d", registryCount)
	}

	for _, don := range ds.DONs {
		for _, cap := range don.Capabilities {
			chainID, ok := ParseEVMChainIDFromCapability(cap)
			if !ok {
				continue
			}
			if !seen[chainID] {
				return fmt.Errorf("dons[%s]: capability %q references chain %d which is not declared in [[chains]]", don.Name, cap, chainID)
			}
		}
	}

	return nil
}

// RegistryChain returns the chain declared with registry = true.
func (ds *DesiredState) RegistryChain() (Chain, bool) {
	for _, ch := range ds.Chains {
		if ch.Registry {
			return ch, true
		}
	}
	return Chain{}, false
}

// ChainIDs returns all declared chain IDs.
func (ds *DesiredState) ChainIDs() []uint64 {
	ids := make([]uint64, 0, len(ds.Chains))
	for _, ch := range ds.Chains {
		ids = append(ids, ch.ChainID)
	}
	return ids
}

// ChainByID returns the declared chain with the given ID, or false if absent.
func (ds *DesiredState) ChainByID(id uint64) (Chain, bool) {
	for _, ch := range ds.Chains {
		if ch.ChainID == id {
			return ch, true
		}
	}
	return Chain{}, false
}

// ParseEVMChainIDFromCapability extracts the chain ID suffix from an
// EVM-scoped capability name (e.g. "evm-1337" -> 1337, true).
func ParseEVMChainIDFromCapability(capability string) (uint64, bool) {
	if !strings.HasPrefix(capability, cre.EVMCapability+"-") {
		return 0, false
	}
	chainPart := strings.TrimPrefix(capability, cre.EVMCapability+"-")
	if chainPart == "" {
		return 0, false
	}
	chainID, err := strconv.ParseUint(chainPart, 10, 64)
	if err != nil {
		return 0, false
	}
	return chainID, true
}

func (ds *DesiredState) hasCapabilityConfig(capName string, don DON) bool {
	if _, ok := ds.CapabilityConfigs[capName]; ok {
		return true
	}
	if don.CapabilityConfigs != nil {
		if _, ok := don.CapabilityConfigs[capName]; ok {
			return true
		}
	}
	// Chain-scoped capabilities like "evm-1337" use the base name "evm"
	base := stripChainSuffix(capName)
	if base != capName {
		if _, ok := ds.CapabilityConfigs[base]; ok {
			return true
		}
		if don.CapabilityConfigs != nil {
			if _, ok := don.CapabilityConfigs[base]; ok {
				return true
			}
		}
	}
	return false
}

// stripChainSuffix removes a chain ID suffix from a capability name.
// e.g. "evm-1337" -> "evm", "solana-123" -> "solana".
// Returns the original string if there is no suffix.
func stripChainSuffix(capName string) string {
	for _, base := range []string{cre.EVMCapability, cre.SolanaCapability} {
		if strings.HasPrefix(capName, base+"-") {
			return base
		}
	}
	return capName
}

// configlessCapabilities are capabilities that don't require a
// capability_configs entry. They are built-in or auto-configured.
var configlessCapabilities = map[string]bool{
	cre.DONTimeCapability:   true,
	cre.ConsensusCapability: true,
}

func isConfiglessCapability(capName string) bool {
	if configlessCapabilities[capName] {
		return true
	}
	// Also check the base name (strip chain suffix)
	base := stripChainSuffix(capName)
	if base != capName && configlessCapabilities[base] {
		return true
	}
	return false
}

// NeedsGateway returns true if any DON has gateway-routed capabilities.
func (ds *DesiredState) NeedsGateway() bool {
	gatewayCaps := map[string]bool{
		cre.VaultCapability:       true,
		cre.HTTPActionCapability:  true,
		cre.HTTPTriggerCapability: true,
	}
	for _, don := range ds.DONs {
		for _, cap := range don.Capabilities {
			if gatewayCaps[cap] {
				return true
			}
		}
	}
	return false
}

// DONByName returns the DON with the given name, or nil if not found.
func (ds *DesiredState) DONByName(name string) *DON {
	for i := range ds.DONs {
		if ds.DONs[i].Name == name {
			return &ds.DONs[i]
		}
	}
	return nil
}

// GatewayDONFor returns the DON name that a gateway node should serve.
// If there's an explicit [[gateway_nodes]] assignment, use it.
// Otherwise default to the first workflow DON.
func (ds *DesiredState) GatewayDONFor(gatewayNodeName string) string {
	for _, gwn := range ds.GatewayNodes {
		if gwn.Node == gatewayNodeName {
			return gwn.DON
		}
	}
	for _, don := range ds.DONs {
		if slices.Contains(don.DONTypes, cre.WorkflowDON) {
			return don.Name
		}
	}
	return ds.DONs[0].Name
}

// ChainScopedCapabilities is the set of capability base names that require
// a chain ID suffix when added to a DON (e.g. "evm" → "evm-1337").
var ChainScopedCapabilities = map[string]bool{
	"evm":    true,
	"solana": true,
	"aptos":  true,
}

// LoadCapabilityDefaults parses the capability defaults TOML string.
func LoadCapabilityDefaults(raw string) map[string]CapabilityConfig {
	if raw == "" {
		return nil
	}
	return parseCapabilityDefaultsString(raw)
}

func parseCapabilityDefaultsString(data string) map[string]CapabilityConfig {
	// The TOML structure is [capability_configs.<name>] with binary_name and
	// [capability_configs.<name>.values] sub-tables.
	// pelletier/go-toml/v2 decodes this as nested maps:
	// { "capability_configs": { "evm": { "binary_name": "evm", "values": {...} } } }
	var raw map[string]map[string]any
	if err := toml.Unmarshal([]byte(data), &raw); err != nil {
		return nil
	}

	ccSection, ok := raw["capability_configs"]
	if !ok {
		return nil
	}

	result := make(map[string]CapabilityConfig)
	for capName, val := range ccSection {
		cc := CapabilityConfig{}
		if vals, ok := val.(map[string]any); ok {
			if bn, ok := vals["binary_name"].(string); ok {
				cc.BinaryName = bn
			}
			if v, ok := vals["values"].(map[string]any); ok {
				cc.Values = v
			}
		}
		result[capName] = cc
	}

	return result
}
