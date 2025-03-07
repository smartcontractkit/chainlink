package solana_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestGenericOps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Msg  string
		Mcms bool
	}{
		{
			Msg:  "with mcms",
			Mcms: true,
		},
		{
			Msg:  "without mcms",
			Mcms: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
			solChain := tenv.Env.AllChainSelectorsSolana()[0]
			e := tenv.Env

			var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
			if test.Mcms {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
					RouterOwnedByTimelock:    true,
					FeeQuoterOwnedByTimelock: true,
					OffRampOwnedByTimelock:   true,
				}
			}

			e, err := commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetDefaultCodeVersion),
					ccipChangesetSolana.SetDefaultCodeVersionConfig{
						ChainSelector: solChain,
						VersionEnum:   1,
						MCMSSolana:    mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.UpdateEnableManualExecutionAfter),
					ccipChangesetSolana.UpdateEnableManualExecutionAfterConfig{
						ChainSelector:         solChain,
						EnableManualExecution: 1,
						MCMSSolana:            mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.UpdateSvmChainSelector),
					ccipChangesetSolana.UpdateSvmChainSelectorConfig{
						OldChainSelector: solChain,
						NewChainSelector: solChain + 1,
						MCMSSolana:       mcmsConfig,
					},
				),
			},
			)
			require.NoError(t, err)
		})
	}
}
