package cre

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

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
func (w wfRegTestStub) MaxActivationRetries() int               { return 0 }
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

	registry := newLocalTestMetadataRegistry(nil)
	require.Equal(t, uint8(0), registry.WorkflowDONF)
}
