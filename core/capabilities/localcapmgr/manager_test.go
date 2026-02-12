package localcapmgr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestConfigHash(t *testing.T) {
	t.Run("empty config returns empty string", func(t *testing.T) {
		assert.Equal(t, "", configHash(nil))
		assert.Equal(t, "", configHash([]byte{}))
	})

	t.Run("same config produces same hash", func(t *testing.T) {
		cfg := []byte(`{"key": "value"}`)
		h1 := configHash(cfg)
		h2 := configHash(cfg)
		assert.Equal(t, h1, h2)
		assert.NotEmpty(t, h1)
	})

	t.Run("different config produces different hash", func(t *testing.T) {
		h1 := configHash([]byte(`{"key": "value1"}`))
		h2 := configHash([]byte(`{"key": "value2"}`))
		assert.NotEqual(t, h1, h2)
	})
}

func TestRunningKey(t *testing.T) {
	assert.Equal(t, "cron@1.0.0:1", runningKey("cron@1.0.0", 1))
	assert.Equal(t, "consensus@1.0.0:42", runningKey("consensus@1.0.0", 42))
}

func testLogger(t *testing.T) corelogger.Logger {
	return corelogger.TestLogger(t)
}

func TestBuildDesiredState(t *testing.T) {
	lggr := testLogger(t)
	mgr := &localCapabilityManager{
		lggr:     lggr,
		localCfg: &testLocalCapabilities{allowlisted: map[string]bool{"cron@1.0.0": true, "consensus@1.0.0": true}},
	}

	dons := []registrysyncer.DON{
		{
			DON: capabilities.DON{ID: 1},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"cron@1.0.0":      {Config: []byte(`{"interval": 60}`)},
				"consensus@1.0.0": {Config: []byte(`{"key": "evm"}`)},
				"unknown@1.0.0":   {Config: []byte(`{}`)}, // not allowlisted
			},
		},
		{
			DON: capabilities.DON{ID: 2},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"cron@1.0.0": {Config: []byte(`{"interval": 30}`)},
			},
		},
	}

	desired := mgr.buildDesiredState(dons)

	assert.Len(t, desired, 3)
	assert.Contains(t, desired, "cron@1.0.0:1")
	assert.Contains(t, desired, "consensus@1.0.0:1")
	assert.Contains(t, desired, "cron@1.0.0:2")
	assert.NotContains(t, desired, "unknown@1.0.0:1")
}

func TestBuildDesiredState_NilLocalConfig(t *testing.T) {
	lggr := testLogger(t)
	mgr := &localCapabilityManager{
		lggr:     lggr,
		localCfg: nil,
	}

	dons := []registrysyncer.DON{
		{
			DON: capabilities.DON{ID: 1},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"cron@1.0.0": {Config: []byte(`{}`)},
			},
		},
	}

	desired := mgr.buildDesiredState(dons)
	assert.Empty(t, desired, "nil config should not allow any capabilities")
}

func noopServiceBuilder(_ context.Context, _ string, _ string, _ string) ([]job.ServiceCtx, error) {
	return []job.ServiceCtx{&mockService{}}, nil
}

func failingServiceBuilder(_ context.Context, _ string, _ string, _ string) ([]job.ServiceCtx, error) {
	return nil, assert.AnError
}

func TestReconcile_StartsNewCapabilities(t *testing.T) {
	ctx := context.Background()
	lggr := testLogger(t)
	metrics, err := newMetrics()
	require.NoError(t, err)

	mgr := &localCapabilityManager{
		lggr:                lggr,
		localCfg:            &testLocalCapabilities{allowlisted: map[string]bool{"test-cap@1.0.0": true}},
		buildServices:       noopServiceBuilder,
		runningCapabilities: make(map[string]*runningCapability),
		metrics:             metrics,
	}

	dons := []registrysyncer.DON{
		{
			DON: capabilities.DON{ID: 1},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"test-cap@1.0.0": {Config: []byte(`{"test": true}`)},
			},
		},
	}

	err = mgr.Reconcile(ctx, dons, &registrysyncer.LocalRegistry{})
	require.NoError(t, err)

	assert.Len(t, mgr.runningCapabilities, 1)
	assert.Contains(t, mgr.runningCapabilities, "test-cap@1.0.0:1")
}

func TestReconcile_StopsRemovedCapabilities(t *testing.T) {
	ctx := context.Background()
	lggr := testLogger(t)
	metrics, err := newMetrics()
	require.NoError(t, err)

	mockSvc := &mockService{}
	mgr := &localCapabilityManager{
		lggr:          lggr,
		localCfg:      &testLocalCapabilities{allowlisted: map[string]bool{"test-cap@1.0.0": true}},
		buildServices: noopServiceBuilder,
		runningCapabilities: map[string]*runningCapability{
			"test-cap@1.0.0:1": {
				capID:      "test-cap@1.0.0",
				donID:      1,
				services:   []job.ServiceCtx{mockSvc},
				configHash: configHash([]byte(`{}`)),
			},
			"removed-cap@1.0.0:1": {
				capID:    "removed-cap@1.0.0",
				donID:    1,
				services: []job.ServiceCtx{&mockService{}},
			},
		},
		metrics: metrics,
	}

	// Only test-cap@1.0.0 is in the desired state.
	dons := []registrysyncer.DON{
		{
			DON: capabilities.DON{ID: 1},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"test-cap@1.0.0": {Config: []byte(`{}`)},
			},
		},
	}

	err = mgr.Reconcile(ctx, dons, &registrysyncer.LocalRegistry{})
	require.NoError(t, err)

	// removed-cap should be stopped and removed.
	assert.Len(t, mgr.runningCapabilities, 1)
	assert.Contains(t, mgr.runningCapabilities, "test-cap@1.0.0:1")
	assert.NotContains(t, mgr.runningCapabilities, "removed-cap@1.0.0:1")
}

func TestReconcile_DetectsConfigChange(t *testing.T) {
	ctx := context.Background()
	lggr := testLogger(t)
	metrics, err := newMetrics()
	require.NoError(t, err)

	oldCfg := []byte(`{"version": 1}`)
	newCfg := []byte(`{"version": 2}`)

	mockSvc := &mockService{}
	mgr := &localCapabilityManager{
		lggr:          lggr,
		localCfg:      &testLocalCapabilities{allowlisted: map[string]bool{"test-cap@1.0.0": true}},
		buildServices: noopServiceBuilder,
		runningCapabilities: map[string]*runningCapability{
			"test-cap@1.0.0:1": {
				capID:      "test-cap@1.0.0",
				donID:      1,
				services:   []job.ServiceCtx{mockSvc},
				configHash: configHash(oldCfg),
			},
		},
		metrics: metrics,
	}

	dons := []registrysyncer.DON{
		{
			DON: capabilities.DON{ID: 1},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"test-cap@1.0.0": {Config: newCfg},
			},
		},
	}

	err = mgr.Reconcile(ctx, dons, &registrysyncer.LocalRegistry{})
	require.NoError(t, err)

	// Old capability closed, new one started.
	assert.True(t, mockSvc.closed, "old capability should be closed on config change")
	assert.Len(t, mgr.runningCapabilities, 1)
	assert.Contains(t, mgr.runningCapabilities, "test-cap@1.0.0:1")

	// Verify new hash.
	rc := mgr.runningCapabilities["test-cap@1.0.0:1"]
	assert.Equal(t, configHash(newCfg), rc.configHash)
}

func TestReconcile_ContinuesOnStartFailure(t *testing.T) {
	ctx := context.Background()
	lggr := testLogger(t)
	metrics, err := newMetrics()
	require.NoError(t, err)

	mgr := &localCapabilityManager{
		lggr:                lggr,
		localCfg:            &testLocalCapabilities{allowlisted: map[string]bool{"failing-cap@1.0.0": true, "good-cap@1.0.0": true}},
		buildServices:       failingServiceBuilder,
		runningCapabilities: make(map[string]*runningCapability),
		metrics:             metrics,
	}

	dons := []registrysyncer.DON{
		{
			DON: capabilities.DON{ID: 1},
			CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
				"failing-cap@1.0.0": {Config: []byte(`{}`)},
				"good-cap@1.0.0":    {Config: []byte(`{}`)},
			},
		},
	}

	err = mgr.Reconcile(ctx, dons, &registrysyncer.LocalRegistry{})
	require.NoError(t, err) // Reconcile is best-effort, doesn't return errors for individual caps.

	// Both failed since we use failingLOOPStarter.
	assert.Empty(t, mgr.runningCapabilities)
}

func TestResolveCapabilityBinary(t *testing.T) {
	lggr := testLogger(t)

	t.Run("uses TOML override when set", func(t *testing.T) {
		override := "/opt/chainlink/binaries/cron"
		mgr := &localCapabilityManager{
			lggr: lggr,
			localCfg: &testLocalCapabilities{
				allowlisted: map[string]bool{"cron@1.0.0": true},
				configs: map[string]*testCapabilityNodeConfig{
					"cron@1.0.0": {binaryPath: override},
				},
			},
		}
		path := mgr.resolveCapabilityBinary("cron@1.0.0")
		assert.Equal(t, override, path)
	})

	t.Run("falls back to capID when no override", func(t *testing.T) {
		mgr := &localCapabilityManager{
			lggr: lggr,
			localCfg: &testLocalCapabilities{
				allowlisted: map[string]bool{"cron@1.0.0": true},
			},
		}
		path := mgr.resolveCapabilityBinary("cron@1.0.0")
		assert.Equal(t, "cron@1.0.0", path)
	})

	t.Run("falls back to capID with nil config", func(t *testing.T) {
		mgr := &localCapabilityManager{
			lggr:     lggr,
			localCfg: nil,
		}
		path := mgr.resolveCapabilityBinary("cron@1.0.0")
		assert.Equal(t, "cron@1.0.0", path)
	})
}

func TestClose_StopsAllRunningCapabilities(t *testing.T) {
	lggr := testLogger(t)
	metrics, err := newMetrics()
	require.NoError(t, err)

	svc1 := &mockService{}
	svc2 := &mockService{}
	mgr := &localCapabilityManager{
		lggr: lggr,
		runningCapabilities: map[string]*runningCapability{
			"cap1@1.0.0:1": {capID: "cap1@1.0.0", donID: 1, services: []job.ServiceCtx{svc1}},
			"cap2@1.0.0:2": {capID: "cap2@1.0.0", donID: 2, services: []job.ServiceCtx{svc2}},
		},
		metrics: metrics,
	}

	// Start it first so Close works.
	require.NoError(t, mgr.Start(context.Background()))
	require.NoError(t, mgr.Close())

	assert.True(t, svc1.closed)
	assert.True(t, svc2.closed)
	assert.Empty(t, mgr.runningCapabilities)
}

// --- Test helpers ---

type mockService struct {
	started bool
	closed  bool
}

func (m *mockService) Start(context.Context) error { m.started = true; return nil }
func (m *mockService) Close() error                { m.closed = true; return nil }

// testLocalCapabilities implements config.LocalCapabilities for testing.
type testLocalCapabilities struct {
	allowlisted map[string]bool
	configs     map[string]*testCapabilityNodeConfig
}

func (t *testLocalCapabilities) RegistryBasedLaunchAllowlist() []string {
	result := make([]string, 0, len(t.allowlisted))
	for k := range t.allowlisted {
		result = append(result, k)
	}
	return result
}

func (t *testLocalCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	return nil
}

func (t *testLocalCapabilities) IsAllowlisted(capabilityID string) bool {
	return t.allowlisted[capabilityID]
}

func (t *testLocalCapabilities) GetCapabilityConfig(capabilityID string) config.CapabilityNodeConfig {
	if t.configs == nil {
		return nil
	}
	c, ok := t.configs[capabilityID]
	if !ok {
		return nil
	}
	return c
}

type testCapabilityNodeConfig struct {
	binaryPath string
	cfg        map[string]string
}

func (c *testCapabilityNodeConfig) BinaryPathOverride() string { return c.binaryPath }
func (c *testCapabilityNodeConfig) Config() map[string]string  { return c.cfg }
