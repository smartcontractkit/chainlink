package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestUpdateOnRampsDests(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			// Default env just has 2 chains with all contracts
			// deployed but no lanes.
			tenv := NewMemoryEnvironment(t)
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
					Config: UpdateOnRampsDestsConfig{
						UpdatesByChain: map[uint64][]OnRampDestinationUpdate{
							source: {
								{
									DestinationSelector: dest,
									IsEnabled:           true,
									TestRouter:          true,
									AllowListEnabled:    false,
								},
							},
							dest: {
								{
									DestinationSelector: source,
									IsEnabled:           true,
									TestRouter:          false,
									AllowListEnabled:    true,
								},
							},
						},
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// Assert the onramp configuration is as we expect.
			sourceCfg, err := state.Chains[source].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, dest)
			require.NoError(t, err)
			require.Equal(t, state.Chains[source].TestRouter.Address(), sourceCfg.Router)
			require.Equal(t, false, sourceCfg.AllowlistEnabled)
			destCfg, err := state.Chains[dest].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, source)
			require.NoError(t, err)
			require.Equal(t, state.Chains[dest].Router.Address(), destCfg.Router)
			require.Equal(t, true, destCfg.AllowlistEnabled)
		})
	}
}
