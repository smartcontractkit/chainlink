package internal_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kstest "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/stretchr/testify/require"
)

func Test_RegisterNOPS(t *testing.T) {
	var (
		useMCMS   bool
		lggr      = logger.Test(t)
		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry  = setupResp.Registry
		chain     = setupResp.Chain
		nops      = make([]kcr.CapabilitiesRegistryNodeOperator, 0)
	)
	t.Run("success create add NOPs mcms proposal", func(t *testing.T) {
		nops = append(nops, kcr.CapabilitiesRegistryNodeOperator{
			Name: "test-nop",
		})
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterNOPS(context.TODO(), lggr, internal.RegisterNOPSRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			Nops:                  nops,
			UseMCMS:               useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 1)
	})
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

func Test_RegisterNodes(t *testing.T) {
	var (
		useMCMS   bool
		lggr      = logger.Test(t)
		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry  = setupResp.Registry
		chain     = setupResp.Chain
	)
	t.Run("success create add nodes mcms proposal", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterNodes(lggr, &internal.RegisterNodesRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			UseMCMS:               useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 1)
	})
}

func Test_RegisterDons(t *testing.T) {
	var (
		useMCMS   bool
		lggr      = logger.Test(t)
		setupResp = kstest.SetupTestRegistry(t, lggr, &kstest.SetupTestRegistryRequest{})
		registry  = setupResp.Registry
		chain     = setupResp.Chain
	)
	t.Run("success create add DONs mcms proposal", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterDons(lggr, internal.RegisterDonsRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			DonToCapabilities: map[string][]internal.RegisteredCapability{
				"test-don": {},
			},
			DonsToRegister: []internal.DONToRegister{
				{
					Name: "test-don",
					F:    2,
				},
			},
			UseMCMS: useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 1)
	})

	t.Run("success create add DONs mcms proposal with multiple DONs", func(t *testing.T) {
		useMCMS = true
		env := &deployment.Environment{
			Logger: lggr,
			Chains: map[uint64]deployment.Chain{
				chain.Selector: chain,
			},
			ExistingAddresses: deployment.NewMemoryAddressBookFromMap(map[uint64]map[string]deployment.TypeAndVersion{
				chain.Selector: {
					registry.Address().String(): deployment.TypeAndVersion{
						Type:    internal.CapabilitiesRegistry,
						Version: deployment.Version1_0_0,
					},
				},
			}),
		}
		resp, err := internal.RegisterDons(lggr, internal.RegisterDonsRequest{
			Env:                   env,
			RegistryChainSelector: chain.Selector,
			DonToCapabilities: map[string][]internal.RegisteredCapability{
				"test-don-1": {},
				"test-don-2": {},
			},
			DonsToRegister: []internal.DONToRegister{
				{
					Name: "test-don-1",
					F:    2,
				},
				{
					Name: "test-don-2",
					F:    2,
				},
			},
			UseMCMS: useMCMS,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Ops)
		require.Len(t, resp.Ops.Batch, 2)
	})
}
