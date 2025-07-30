package solana_test

import (
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestGenericOpsWithMcms(t *testing.T) {
	t.Parallel()
	skipInCI(t) // takes too long in CI
	doTestGenericOps(t, true)
}

func TestGenericOpsWithoutMcms(t *testing.T) {
	t.Parallel()
	skipInCI(t)
	doTestGenericOps(t, false)
}

func doTestGenericOps(t *testing.T, mcms bool) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	e := tenv.Env

	var mcmsConfig *proposalutils.TimelockConfig
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChain, true,
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
}
