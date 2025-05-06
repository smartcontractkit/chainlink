package changeset_test

import (
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func TestAddNops(t *testing.T) {
	t.Parallel()

	nops := []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: gethcommon.HexToAddress("0x01"),
			Name:  "new test nop1",
		},
		{
			Admin: gethcommon.HexToAddress("0x02"),
			Name:  "another test nop2",
		},
	}
	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		req := &changeset.AddNopsRequest{
			RegistryChainSel: te.RegistrySelector,
			Nops:             nops,
			RegistryRef:      te.CapabilityRegistryAddressRef(),
		}
		csOut, err := changeset.AddNops(te.Env, req)
		require.NoError(t, err)
		require.Empty(t, csOut.MCMSTimelockProposals)
		require.Nil(t, csOut.AddressBook)
		assertNopsExist(t, te.CapabilitiesRegistry(), nops...)
		t.Run("idempotent", func(t *testing.T) {
			_, err = changeset.AddNops(te.Env, req)
			require.NoError(t, err)
			assertNopsExist(t, te.CapabilitiesRegistry(), nops...)
		})
		t.Run("deduplication", func(t *testing.T) {
			req.Nops = append(req.Nops, req.Nops...)
			_, err = changeset.AddNops(te.Env, req)
			require.NoError(t, err)
			assertNopsExist(t, te.CapabilitiesRegistry(), nops...)
		})
	})

	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		req := &changeset.AddNopsRequest{
			RegistryChainSel: te.RegistrySelector,
			Nops:             nops,
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
			RegistryRef:      te.CapabilityRegistryAddressRef(),
		}
		csOut, err := changeset.AddNops(te.Env, req)
		require.NoError(t, err)
		require.Len(t, csOut.MCMSTimelockProposals, 1)
		require.Nil(t, csOut.AddressBook)

		err = applyProposal(t, te, commonchangeset.Configure(cldf.CreateLegacyChangeSet(changeset.AddNops), req))
		require.NoError(t, err)

		assertNopsExist(t, te.CapabilitiesRegistry(), nops...)
	})
}

func assertNopsExist(t *testing.T, cr *kcr.CapabilitiesRegistry, expected ...kcr.CapabilitiesRegistryNodeOperator) {
	t.Helper()
	ops, err := cr.GetNodeOperators(nil)
	require.NoError(t, err)
	for _, want := range expected {
		assert.Contains(t, ops, want)
	}
}
