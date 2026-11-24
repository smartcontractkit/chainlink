package config

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func derefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func TestMeteringNodeConfig(t *testing.T) {
	t.Parallel()

	dm := &cre.DonMetadata{Name: "workflow"}
	got := meteringNodeConfig(dm, 0)

	assert.True(t, derefBool(got.MeterRecordsEnabled))
	assert.True(t, derefBool(got.MeterSnapshotsEnabled))
	assert.Equal(t, "cre", derefStr(got.Product))
	assert.Equal(t, "local-cre", derefStr(got.Tenant))
	assert.Equal(t, "1", derefStr(got.NumericTenantID))
	assert.Equal(t, "local", derefStr(got.Environment))
	assert.Equal(t, "workflow", derefStr(got.Zone))
	assert.Equal(t, "workflow-node-0", derefStr(got.NodeID))
}

func TestMeteringNodeConfig_NodeIDUniqueness(t *testing.T) {
	t.Parallel()

	dm := &cre.DonMetadata{Name: "capabilities"}
	seen := make(map[string]struct{})
	for i := range 4 {
		got := meteringNodeConfig(dm, i)
		nodeID := derefStr(got.NodeID)
		expected := "capabilities-node-" + strconv.Itoa(i)
		assert.Equal(t, expected, nodeID)
		_, dup := seen[nodeID]
		assert.False(t, dup, "NodeID %q duplicated", nodeID)
		seen[nodeID] = struct{}{}
	}
}

func TestValidateUserConfigOverrides(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		overrides string
		wantErr   string
	}{
		{
			name:      "allows non-managed sections",
			overrides: "[Log]\nLevel = 'debug'\n\n[CRE.WorkflowFetcher]\nURL = 'file:///home/chainlink/workflows'\n",
		},
		{
			name:      "allows comments mentioning managed sections",
			overrides: "# do not set [Telemetry] or [Metering] here\n[Log]\nLevel = 'debug'\n",
		},
		{
			name:      "rejects Telemetry",
			overrides: "[Telemetry]\nChipIngressEndpoint = 'chip-ingress:50051'\n",
			wantErr:   "[Telemetry] is framework-managed",
		},
		{
			name:      "rejects Metering",
			overrides: "[Metering]\nMeterRecordsEnabled = true\n",
			wantErr:   "[Metering] is framework-managed",
		},
		{
			name:      "rejects Billing",
			overrides: "[Billing]\nURL = 'host.docker.internal:2223'\n",
			wantErr:   "[Billing] is framework-managed",
		},
		{
			name:      "rejects invalid TOML",
			overrides: "[Log\nLevel = 'debug'\n",
			wantErr:   "not valid TOML",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validateUserConfigOverrides(tc.overrides)
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

// TestAddWorkerNodeConfig_MeteringAndChipIngress is a thin integration test
// through addWorkerNodeConfig on a capabilities DON (not WorkflowDON/ShardDON)
// to verify:
//  1. [Metering] is populated when enableMetering is true
//  2. Telemetry.ChipIngressEndpoint equals the chip-router URL passed in
//
// topology is nil because the capabilities path never accesses it: the
// WorkflowDON and ShardDON branches are both skipped when Flags does not
// contain "workflow" or "shard".
func TestAddWorkerNodeConfig_MeteringAndChipIngress(t *testing.T) {
	t.Parallel()

	const routerURL = "chip-router:50051"

	inputs := capabilitiesCommonInputs(routerURL)
	inputs.enableMetering = true

	config, err := addWorkerNodeConfig(corechainlink.Config{}, nil, testPeeringData(), inputs, capabilitiesDonMetadata(), &cre.NodeMetadata{
		Index: 1,
		Roles: []string{cre.WorkerNode},
	})
	require.NoError(t, err)

	// [Metering] is present with the expected NodeID
	require.NotNil(t, config.Metering.MeterRecordsEnabled)
	assert.True(t, *config.Metering.MeterRecordsEnabled)
	assert.Equal(t, "capabilities-node-1", *config.Metering.NodeID)

	// Telemetry.ChipIngressEndpoint equals the router URL — regression test for
	// the original bug where user_config_overrides hardcoded chip-ingress:50051
	require.NotNil(t, config.Telemetry.ChipIngressEndpoint)
	assert.Equal(t, routerURL, *config.Telemetry.ChipIngressEndpoint)
}

// TestAddWorkerNodeConfig_MeteringDisabled verifies no [Metering] is injected
// when enableMetering is false, while ChipIngressEndpoint still points at the
// router.
func TestAddWorkerNodeConfig_MeteringDisabled(t *testing.T) {
	t.Parallel()

	const routerURL = "chip-router:50051"

	config, err := addWorkerNodeConfig(corechainlink.Config{}, nil, testPeeringData(), capabilitiesCommonInputs(routerURL), capabilitiesDonMetadata(), &cre.NodeMetadata{
		Index: 0,
		Roles: []string{cre.WorkerNode},
	})
	require.NoError(t, err)

	assert.Nil(t, config.Metering.MeterRecordsEnabled, "Metering should not be set when disabled")
	require.NotNil(t, config.Telemetry.ChipIngressEndpoint)
	assert.Equal(t, routerURL, *config.Telemetry.ChipIngressEndpoint)
}

func TestAddWorkerNodeConfig_HTTPRegistryLaunch(t *testing.T) {
	t.Parallel()
	don := capabilitiesDonMetadata()
	don.Flags = append(don.Flags, cre.HTTPActionCapability, cre.HTTPTriggerCapability)
	don.RegistryBasedLaunchAllowlist = []string{"cron-trigger@1.0.0"}
	config, err := addWorkerNodeConfig(corechainlink.Config{}, nil, testPeeringData(), capabilitiesCommonInputs("chip-router:50051"), don, &cre.NodeMetadata{Index: 0, Roles: []string{cre.WorkerNode}})
	require.NoError(t, err)
	assert.Equal(t, []string{"cron-trigger@1.0.0", "http-actions@1.0.0-alpha", "http-trigger@1.0.0-alpha"}, config.Capabilities.Local.RegistryBasedLaunchAllowlist)
	encoded, err := toml.Marshal(config)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), "[Capabilities.HTTPAction")
	assert.NotContains(t, string(encoded), "[Capabilities.HTTPTrigger")
	// Reapplying configuration must not duplicate registry launch entries.
	config, err = addWorkerNodeConfig(config, nil, testPeeringData(), capabilitiesCommonInputs("chip-router:50051"), don, &cre.NodeMetadata{Index: 0, Roles: []string{cre.WorkerNode}})
	require.NoError(t, err)
	assert.Len(t, config.Capabilities.Local.RegistryBasedLaunchAllowlist, 3)
}

func TestAddWorkerNodeConfig_ChainCapabilitiesRegistryLaunch(t *testing.T) {
	t.Parallel()
	don := capabilitiesDonMetadata()
	don.Flags = append(don.Flags, cre.EVMCapability+"-1337", cre.AptosCapability+"-4", cre.StellarCapability)
	don.RegistryBasedLaunchAllowlist = []string{"cron-trigger@1.0.0"}
	inputs := capabilitiesCommonInputs("chip-router:50051")
	inputs.evmChains = []*evmChain{{Name: "evm", ChainID: 1337}, {Name: "evm-other", ChainID: 2337}}
	inputs.aptosChains = []*aptosChain{
		{ChainID: "4", ChainSelector: 4457093679053095497},
		{ChainID: "8", ChainSelector: 999}, // chain without the capability enabled
	}
	inputs.stellarChain = &stellarChain{Name: "stellar", ChainID: "net", ChainSelector: 7777}
	node := &cre.NodeMetadata{Index: 0, Roles: []string{cre.WorkerNode}}

	config, err := addWorkerNodeConfig(corechainlink.Config{}, nil, testPeeringData(), inputs, don, node)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"cron-trigger@1.0.0",
		"^evm:ChainSelector:",
		`^aptos:ChainSelector:4457093679053095497@1\.0\.0$`,
		`^stellar:ChainSelector:7777@1\.0\.0$`,
	}, config.Capabilities.Local.RegistryBasedLaunchAllowlist)

	// Reapplying configuration must not duplicate registry launch entries.
	config, err = addWorkerNodeConfig(config, nil, testPeeringData(), inputs, don, node)
	require.NoError(t, err)
	assert.Len(t, config.Capabilities.Local.RegistryBasedLaunchAllowlist, 4)
}

func TestAddWorkerNodeConfig_ChainCapabilitiesNotAllowlistedWithoutFlags(t *testing.T) {
	t.Parallel()
	inputs := capabilitiesCommonInputs("chip-router:50051")
	inputs.evmChains = []*evmChain{{Name: "evm", ChainID: 1337}}
	inputs.aptosChains = []*aptosChain{{ChainID: "4", ChainSelector: 1}}
	inputs.stellarChain = &stellarChain{Name: "stellar", ChainID: "net", ChainSelector: 2}

	config, err := addWorkerNodeConfig(corechainlink.Config{}, nil, testPeeringData(), inputs, capabilitiesDonMetadata(), &cre.NodeMetadata{Index: 0, Roles: []string{cre.WorkerNode}})
	require.NoError(t, err)
	assert.Empty(t, config.Capabilities.Local.RegistryBasedLaunchAllowlist)
}

func TestChainCapabilityAllowlistMatchesCapabilityIDs(t *testing.T) {
	t.Parallel()
	don := capabilitiesDonMetadata()
	don.Flags = append(don.Flags, cre.EVMCapability+"-1337", cre.AptosCapability+"-4", cre.StellarCapability)
	inputs := capabilitiesCommonInputs("chip-router:50051")
	inputs.evmChains = []*evmChain{{ChainID: 1337}}
	inputs.aptosChains = []*aptosChain{{ChainID: "4", ChainSelector: 11}}
	inputs.stellarChain = &stellarChain{ChainSelector: 22}

	matches := func(id string) bool {
		for _, p := range chainCapabilityAllowlistEntries(don, inputs) {
			if regexp.MustCompile(p).MatchString(id) {
				return true
			}
		}
		return false
	}
	assert.True(t, matches("evm:ChainSelector:12345@1.0.0"))
	assert.True(t, matches("aptos:ChainSelector:11@1.0.0"))
	assert.True(t, matches("stellar:ChainSelector:22@1.0.0"))
	assert.False(t, matches("aptos:ChainSelector:111@1.0.0"), "must not match another selector sharing a prefix")
	assert.False(t, matches("aptos:ChainSelector:11@1.0.1"))
	assert.False(t, matches("solana:ChainSelector:11@1.0.0"), "solana stays on job specs")
}

// TestAddWorkerNodeConfig_ConsensusQueueRegistryLaunch asserts a DON with the consensus-queue flag gets it on the
// registry launch allowlist exactly once, next to the topology's own entries.
func TestAddWorkerNodeConfig_ConsensusQueueRegistryLaunch(t *testing.T) {
	t.Parallel()
	don := capabilitiesDonMetadata()
	don.Flags = append(don.Flags, cre.ConsensusQueueCapability)
	don.RegistryBasedLaunchAllowlist = []string{"cron-trigger@1.0.0", "dontime@1.0.0"}

	config, err := addWorkerNodeConfig(corechainlink.Config{}, nil, testPeeringData(), capabilitiesCommonInputs("chip-router:50051"), don, &cre.NodeMetadata{Index: 0, Roles: []string{cre.WorkerNode}})
	require.NoError(t, err)
	want := []string{"cron-trigger@1.0.0", "dontime@1.0.0", "consensus-queue@1.0.0"}
	assert.Equal(t, want, config.Capabilities.Local.RegistryBasedLaunchAllowlist)

	config, err = addWorkerNodeConfig(config, nil, testPeeringData(), capabilitiesCommonInputs("chip-router:50051"), don, &cre.NodeMetadata{Index: 0, Roles: []string{cre.WorkerNode}})
	require.NoError(t, err)
	assert.Equal(t, want, config.Capabilities.Local.RegistryBasedLaunchAllowlist, "reapplying must not duplicate entries")
}

func capabilitiesDonMetadata() *cre.DonMetadata {
	return &cre.DonMetadata{
		Name:  "capabilities",
		Flags: []string{"capabilities"},
	}
}

func capabilitiesCommonInputs(routerURL string) *commonInputs {
	return &commonInputs{
		provider:                  infra.Provider{Type: "docker"},
		chipRouterInternalGRPCURL: routerURL,
		registryChainID:           1337,
		registryChainSelector:     1,
		capabilityRegistry:        versionedAddress{address: "0xdead", version: mustSemVer("1.0.0")},
		workflowRegistry:          versionedAddress{address: "0xbeef", version: mustSemVer("1.0.0")},
	}
}

func testPeeringData() cre.PeeringData {
	return cre.PeeringData{
		BootstrapperPeerID: "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X",
		BootstrapperHost:   "bootstrap",
		Port:               4222,
	}
}

func mustSemVer(v string) *semver.Version {
	ver, err := semver.NewVersion(v)
	if err != nil {
		panic(err)
	}
	return ver
}
