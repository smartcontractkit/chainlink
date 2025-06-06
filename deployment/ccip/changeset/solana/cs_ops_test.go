package solana_test

import (
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	tutils "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

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
			if test.Msg == "with mcms" {
				tutils.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-437")
			}
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
			solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
			e := tenv.Env

			var mcmsConfig *proposalutils.TimelockConfig
			if test.Mcms {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 1 * time.Second,
				}
			}

			e, _, err := commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetDefaultCodeVersion),
					ccipChangesetSolana.SetDefaultCodeVersionConfig{
						ChainSelector: solChain,
						VersionEnum:   1,
						MCMS:          mcmsConfig,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(ccipChangesetSolana.UpdateEnableManualExecutionAfter),
					ccipChangesetSolana.UpdateEnableManualExecutionAfterConfig{
						ChainSelector:         solChain,
						EnableManualExecution: 1,
						MCMS:                  mcmsConfig,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(ccipChangesetSolana.UpdateSvmChainSelector),
					ccipChangesetSolana.UpdateSvmChainSelectorConfig{
						OldChainSelector: solChain,
						NewChainSelector: solChain + 1,
						MCMS:             mcmsConfig,
					},
				),
			},
			)
			require.NoError(t, err)
		})
	}
}
