package internal_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/test-go/testify/require"
)

func Test_RegisterNOPS(t *testing.T) {
	t.Skip()
}

func Test_AddCapabilities(t *testing.T) {
	var (
		useMCMS      bool
		lggr         = logger.Test(t)
		setupResp    = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry     = setupResp.Registry
		chain        = setupResp.Chain
		capabilities = make([]kcr.CapabilitiesRegistryCapability, 0)
	)

	t.Run("successfully create mcms proposal", func(t *testing.T) {
		useMCMS = true
		capabilities = append(capabilities, kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cap1",
			Version:        "1.0.0",
			CapabilityType: 0,
		})
		ops, err := internal.AddCapabilities(lggr, registry, chain, capabilities, useMCMS)
		require.NoError(t, err)
		require.NotNil(t, ops)
		require.Len(t, ops.Batch, 1)
	})

	t.Run("does nothing if no capabilities", func(t *testing.T) {
		ops, err := internal.AddCapabilities(lggr, registry, chain, nil, useMCMS)
		require.NoError(t, err)
		require.Nil(t, ops)
	})
}
