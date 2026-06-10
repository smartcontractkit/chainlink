package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	capStreams "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// wfRegTestStub implements config.CapabilitiesWorkflowRegistry for tests.
type wfRegTestStub struct {
	addr           string
	additionalURLs []string
}

type wfRegAddSrcStub struct{ u string }

func (a wfRegAddSrcStub) GetURL() string      { return a.u }
func (a wfRegAddSrcStub) GetTLSEnabled() bool { return false }
func (a wfRegAddSrcStub) GetName() string     { return "" }

type wfRegStorageStub struct{}

func (wfRegStorageStub) ArtifactStorageHost() string { return "" }
func (wfRegStorageStub) URL() string                 { return "" }
func (wfRegStorageStub) TLSEnabled() bool            { return false }

type wfRegModuleCacheStub struct{}

func (wfRegModuleCacheStub) Enabled() bool              { return false }
func (wfRegModuleCacheStub) DiskMonitorEnabled() bool   { return false }
func (wfRegModuleCacheStub) IdleEviction() bool         { return true }
func (wfRegModuleCacheStub) IdleTimeout() time.Duration { return 10 * time.Minute }
func (wfRegModuleCacheStub) MaxLoaded() int             { return 200 }
func (wfRegModuleCacheStub) CacheDir() string           { return "" }

func (w wfRegTestStub) Address() string                         { return w.addr }
func (w wfRegTestStub) NetworkID() string                       { return "" }
func (w wfRegTestStub) ChainID() string                         { return "" }
func (w wfRegTestStub) ContractVersion() string                 { return "" }
func (w wfRegTestStub) MaxEncryptedSecretsSize() utils.FileSize { return 0 }
func (w wfRegTestStub) MaxBinarySize() utils.FileSize           { return 0 }
func (w wfRegTestStub) MaxConfigSize() utils.FileSize           { return 0 }
func (w wfRegTestStub) RelayID() commontypes.RelayID            { return commontypes.RelayID{} }
func (w wfRegTestStub) SyncStrategy() string                    { return "" }
func (w wfRegTestStub) MaxConcurrency() int                     { return 0 }
func (w wfRegTestStub) WorkflowStorage() config.WorkflowStorage { return wfRegStorageStub{} }
func (w wfRegTestStub) ModuleCache() config.ModuleCache         { return wfRegModuleCacheStub{} }
func (w wfRegTestStub) AdditionalSources() []config.AdditionalWorkflowSource {
	out := make([]config.AdditionalWorkflowSource, len(w.additionalURLs))
	for i, u := range w.additionalURLs {
		out[i] = wfRegAddSrcStub{u: u}
	}
	return out
}

func testWorkflowRegistry(addr string, urls ...string) config.CapabilitiesWorkflowRegistry {
	return wfRegTestStub{addr: addr, additionalURLs: urls}
}

type testLocalCapabilities struct {
	cfgs map[string]config.CapabilityNodeConfig
}

func (t testLocalCapabilities) RegistryBasedLaunchAllowlist() []string {
	return nil
}

func (t testLocalCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	return t.cfgs
}

func (t testLocalCapabilities) IsAllowlisted(string) bool {
	return false
}

func (t testLocalCapabilities) GetCapabilityConfig(capabilityID string) config.CapabilityNodeConfig {
	if t.cfgs == nil {
		return nil
	}

	return t.cfgs[capabilityID]
}

type testCapabilityNodeConfig struct{}

func (testCapabilityNodeConfig) BinaryPathOverride() string {
	return ""
}

func (testCapabilityNodeConfig) Config() map[string]string {
	return nil
}

func TestWorkflowRegistrySemverMajor(t *testing.T) {
	t.Parallel()

	major, err := workflowRegistrySemverMajor("")
	require.NoError(t, err)
	require.Equal(t, uint64(2), major)

	major, err = workflowRegistrySemverMajor("   ")
	require.NoError(t, err)
	require.Equal(t, uint64(2), major)

	major, err = workflowRegistrySemverMajor("2.0.0")
	require.NoError(t, err)
	require.Equal(t, uint64(2), major)

	major, err = workflowRegistrySemverMajor("1.0.0")
	require.NoError(t, err)
	require.Equal(t, uint64(1), major)

	_, err = workflowRegistrySemverMajor("not-a-version")
	require.Error(t, err)
}

func TestWorkflowRegistryConfigured(t *testing.T) {
	t.Parallel()

	require.False(t, workflowRegistryConfigured(testWorkflowRegistry(""), 1))
	require.False(t, workflowRegistryConfigured(testWorkflowRegistry("", "", "  "), 1))
	require.True(t, workflowRegistryConfigured(testWorkflowRegistry("0xabc"), 1))

	require.False(t, workflowRegistryConfigured(testWorkflowRegistry(""), 2))
	require.True(t, workflowRegistryConfigured(testWorkflowRegistry("0xdef"), 2))
	require.True(t, workflowRegistryConfigured(testWorkflowRegistry("", "https://example"), 2))
	require.True(t, workflowRegistryConfigured(testWorkflowRegistry("", "", "grpc://x"), 2))
}

func TestNewLocalTestMetadataRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		localCfg  config.LocalCapabilities
		expectedF uint8
	}{
		{
			name:      "default workflow DON fault tolerance",
			localCfg:  nil,
			expectedF: 0,
		},
		{
			name: "mock trigger opt-in uses workflow DON fault tolerance one",
			localCfg: testLocalCapabilities{
				cfgs: map[string]config.CapabilityNodeConfig{
					capStreams.MockTriggerCapabilityID: testCapabilityNodeConfig{},
				},
			},
			expectedF: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := newLocalTestMetadataRegistry(tt.localCfg)
			require.Equal(t, tt.expectedF, registry.WorkflowDONF)
		})
	}
}
