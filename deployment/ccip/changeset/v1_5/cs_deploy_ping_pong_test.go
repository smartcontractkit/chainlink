package v1_5_test

import (
	"testing"

	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/stretchr/testify/require"
)

func TestDeployPingPongDemoContractsChangeset(t *testing.T) {
	for _, tc := range []struct {
		name           string
		chainsToDeploy []struct {
			ChainSelector uint64
			IsTestRouter  bool
		}
		deployRouter     bool
		deployTestRouter bool

		deployLinkToken       bool
		deployStaticLinkToken bool

		expectError bool
	}{
		{
			name: "Valid configuration",
			chainsToDeploy: []struct {
				ChainSelector uint64
				IsTestRouter  bool
			}{
				{ChainSelector: 1, IsTestRouter: true},
				{ChainSelector: 2, IsTestRouter: false},
			},
			deployRouter:          true,
			deployTestRouter:      true,
			deployLinkToken:       true,
			deployStaticLinkToken: true,
			expectError:           false,
		},
		// {
		// 	name: "Invalid configuration - missing router",
		// 	chainsToDeploy: []struct {
		// 		ChainSelector uint64
		// 		IsTestRouter  bool
		// 	}{
		// 		{ChainSelector: 3, IsTestRouter: false},
		// 	},
		// 	deployRouter:          false,
		// 	deployTestRouter:      false,
		// 	deployLinkToken:       true,
		// 	deployStaticLinkToken: true,
		// },
		// {
		// 	name: "Invalid configuration - missing link token",
		// 	chainsToDeploy: []struct {
		// 		ChainSelector uint64
		// 		IsTestRouter  bool
		// 	}{
		// 		{ChainSelector: 3, IsTestRouter: false},
		// 	},
		// 	deployRouter:          true,
		// 	deployLinkToken:       true,
		// 	deployStaticLinkToken: false,
		// },
		// {
		// 	name: "Invalid configuration - missing link token",
		// 	chainsToDeploy: []struct {
		// 		ChainSelector uint64
		// 		IsTestRouter  bool
		// 	}{
		// 		{ChainSelector: 3, IsTestRouter: false},
		// 	},
		// 	deployRouter:          true,
		// 	deployLinkToken:       false,
		// 	deployStaticLinkToken: true,
		// },
	} {
		t.Run(tc.name, func(t *testing.T) {
			lggr := logger.Test(t)

			tenv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Chains: 2,
			})
			allChains := maps.Keys(tenv.Chains)
			selectorA := allChains[0]
			selectorB := allChains[1]

			// changesetsToRun := make([]deployment.ChangeSetV2[any], 0)
			tEnv, err := commonchangeset.Apply(t, tenv, nil,
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
					allChains,
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
					changeset.DeployPrerequisiteConfig{
						Configs: []changeset.DeployPrerequisiteConfigPerChain{
							{
								ChainSelector: selectorA,
								Opts:          []changeset.PrerequisiteOpt{},
							},
							{
								ChainSelector: selectorB,
								Opts:          []changeset.PrerequisiteOpt{},
							},
						},
					},
				),
				commonchangeset.Configure(
					v1_5.DeployPingPongDemoContractChangeset,
					v1_5.DeployPingPongDemoContractsConfig{
						ChainsToDeploy: []struct {
							ChainSelector uint64
							IsTestRouter  bool
						}{
							{
								ChainSelector: selectorA,
								IsTestRouter:  true,
							},
							{
								ChainSelector: selectorB,
								IsTestRouter:  false,
							},
						},
					},
				),
			)

			_, err = changeset.LoadOnchainState(tEnv)
			require.NoError(t, err)

			// if tc.expectError {
			// 	require.Nil(t, state.Chains[selectorA].PingPongDemo)
			// 	require.Error(t, err)
			// } else {
			// 	require.NotNil(t, state.Chains[selectorB].PingPongDemo.Address())
			// 	require.NoError(t, err)
			// }

		})
	}
}
