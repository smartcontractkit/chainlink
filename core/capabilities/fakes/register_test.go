package fakes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	corecaps "github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

func TestRegisterFakeStreamsTrigger(t *testing.T) {
	registry := corecaps.NewRegistry(logger.Test(t))

	trigger, err := RegisterFakeStreamsTrigger(t.Context(), logger.Test(t), registry, 4)
	require.NoError(t, err)
	require.NotNil(t, trigger)

	capability, err := registry.Get(t.Context(), "streams-trigger@1.0.0")
	require.NoError(t, err)

	info, err := capability.Info(t.Context())
	require.NoError(t, err)
	require.Equal(t, "streams-trigger@1.0.0", info.ID)
}

func TestNewFakeStreamsTrigger_UsesDeterministicSigners(t *testing.T) {
	triggerA := NewFakeStreamsTrigger(logger.Test(t), 4)
	triggerB := NewFakeStreamsTrigger(logger.Test(t), 4)

	require.Equal(t, triggerA.meta.Signers, triggerB.meta.Signers)
	require.Equal(t, triggerA.meta.MinRequiredSignatures, triggerB.meta.MinRequiredSignatures)
}
