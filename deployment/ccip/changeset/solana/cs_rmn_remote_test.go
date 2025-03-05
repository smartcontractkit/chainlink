package solana_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestRMNRemoteCurse(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	_, err := ccipChangeset.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)

	testCases := []struct {
		curseConfig ccipChangesetSolana.CurseConfig
		shouldError bool
		cs          func(e deployment.Environment, cfg ccipChangesetSolana.CurseConfig) (deployment.ChangesetOutput, error)
	}{
		{
			curseConfig: ccipChangesetSolana.CurseConfig{
				ChainSelector:       solChain,
				GlobalCurse:         true,
				RemoteChainSelector: evmChain,
			},
			shouldError: true, // incorrect config
			cs:          ccipChangesetSolana.ApplyCurse,
		},
		{
			curseConfig: ccipChangesetSolana.CurseConfig{
				ChainSelector: solChain,
				GlobalCurse:   false,
			},
			shouldError: true, // incorrect config
			cs:          ccipChangesetSolana.ApplyCurse,
		},
		{
			curseConfig: ccipChangesetSolana.CurseConfig{
				ChainSelector: solChain,
				GlobalCurse:   true,
			},
			shouldError: false, // apply global curse
			cs:          ccipChangesetSolana.ApplyCurse,
		},
		{
			curseConfig: ccipChangesetSolana.CurseConfig{
				ChainSelector: solChain,
				GlobalCurse:   true,
			},
			shouldError: false, // remove global curse
			cs:          ccipChangesetSolana.RemoveCurse,
		},
		{
			curseConfig: ccipChangesetSolana.CurseConfig{
				ChainSelector:       solChain,
				GlobalCurse:         false,
				RemoteChainSelector: evmChain,
			},
			shouldError: false, // apply chain curse
			cs:          ccipChangesetSolana.ApplyCurse,
		},
		{
			curseConfig: ccipChangesetSolana.CurseConfig{
				ChainSelector:       solChain,
				GlobalCurse:         false,
				RemoteChainSelector: evmChain,
			},
			shouldError: false, // remove chain curse
			cs:          ccipChangesetSolana.RemoveCurse,
		},
	}

	// register evm chain on router
	e, err := commonchangeset.ApplyChangesetsV2(t, tenv.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
			ccipChangesetSolana.AddRemoteChainToRouterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.RouterConfig{
					evmChain: {
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: true,
						},
					},
				},
			},
		),
	})
	require.NoError(t, err)

	for _, testCase := range testCases {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				deployment.CreateLegacyChangeSet(testCase.cs),
				testCase.curseConfig,
			),
		})
		if testCase.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
