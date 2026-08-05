package creregistry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type stubProxyConfig struct {
	enabled bool
	port    uint16
}

func (s stubProxyConfig) Enabled() bool { return s.enabled }
func (s stubProxyConfig) Port() uint16  { return s.port }

func TestSelect_DisabledUsesInProcessRegistry(t *testing.T) {
	ctx := context.Background()

	reg, closeFn, err := Select(logger.Test(t), stubProxyConfig{enabled: false})
	require.NoError(t, err)
	require.NotNil(t, reg)
	assert.Nil(t, closeFn, "nothing to release when no remote connection was made")

	// The in-process base registry accepts capability values directly, which is
	// the behaviour every existing caller depends on.
	require.NoError(t, reg.Add(ctx, &fakeExecutable{
		info: capabilities.MustNewCapabilityInfo("act@1.0.0", capabilities.CapabilityTypeAction, "act"),
	}))

	got, err := reg.GetExecutable(ctx, "act@1.0.0")
	require.NoError(t, err)
	info, err := got.Info(ctx)
	require.NoError(t, err)
	assert.Equal(t, "act@1.0.0", info.ID)
}

func TestSelect_NilConfigUsesInProcessRegistry(t *testing.T) {
	// A nil config must not panic: it reads as "no proxy configured".
	reg, closeFn, err := Select(logger.Test(t), nil)
	require.NoError(t, err)
	require.NotNil(t, reg)
	assert.Nil(t, closeFn)
}

func TestSelect_EnabledUsesRemoteRegistry(t *testing.T) {
	ctx := context.Background()

	reg, closeFn, err := Select(logger.Test(t), stubProxyConfig{enabled: true, port: 50051})
	require.NoError(t, err)
	require.NotNil(t, reg)
	require.NotNil(t, closeFn, "the caller must be able to release the remote connection")
	t.Cleanup(func() { _ = closeFn() })

	// Metadata now comes from the remote process. Nothing is listening in this
	// test, so the call must fail rather than fall back to an empty local view —
	// a silent fallback would look like "this node is in no DON".
	_, err = reg.LocalNode(ctx)
	require.Error(t, err)

	// Same for capability lookups: no local map is consulted.
	_, err = reg.GetExecutable(ctx, "act@1.0.0")
	require.Error(t, err)
}

func TestSelect_EnabledStillReturnsConcreteRegistryType(t *testing.T) {
	// cre.Opts takes *capabilities.Registry, not the interface, and cre.go calls
	// SetLocalRegistry on it. Select must therefore keep returning the concrete
	// type in both modes or the wiring will not compile.
	reg, closeFn, err := Select(logger.Test(t), stubProxyConfig{enabled: true, port: 50051})
	require.NoError(t, err)
	if closeFn != nil {
		t.Cleanup(func() { _ = closeFn() })
	}
	require.NotNil(t, reg)
}
