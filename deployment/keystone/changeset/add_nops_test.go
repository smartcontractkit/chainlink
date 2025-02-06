package changeset_test

import (
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
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
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
		})

		req := &changeset.AddNopsRequest{
			RegistryChainSel: te.RegistrySelector,
			Nops:             nops,
		}
		csOut, err := changeset.AddNops(te.Env, req)
		require.NoError(t, err)
		require.Empty(t, csOut.Proposals)
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
		te := test.SetupTestEnv(t, test.TestConfig{
			WFDonConfig:     test.DonConfig{N: 4},
			AssetDonConfig:  test.DonConfig{N: 4},
			WriterDonConfig: test.DonConfig{N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		req := &changeset.AddNopsRequest{
			RegistryChainSel: te.RegistrySelector,
			Nops:             nops,
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
		}
		csOut, err := changeset.AddNops(te.Env, req)
		require.NoError(t, err)
		require.Len(t, csOut.Proposals, 1)
		require.Nil(t, csOut.AddressBook)

		// now apply the changeset such that the proposal is signed and execed
		contracts := te.ContractSets()[te.RegistrySelector]
		timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
			te.RegistrySelector: {
				Timelock:  contracts.Timelock,
				CallProxy: contracts.CallProxy,
			},
		}
		_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(changeset.AddNops),
				Config:    req,
			},
		})
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
