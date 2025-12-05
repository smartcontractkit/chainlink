//go:build !integration

package solana_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

// use this for a quick deploy test
func TestDeployChainContractsChangesetPreload(t *testing.T) {
	quarantine.Flaky(t, "DX-1729")
	t.Parallel()

	homeChainSel := chain_selectors.TEST_90000001.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath := t.TempDir()
	progIDs := soltestutils.LoadCCIPPrograms(t, programsPath)
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{homeChainSel}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, progIDs),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	testhelpers.RegisterNodes(t, env, 4, homeChainSel)

	err = testhelpers.SavePreloadedSolAddresses(*env, solSelector)
	require.NoError(t, err)

	e := *env

	// empty build config means, if artifacts are not present, resolve the artifact from github based on go.mod version
	// for a simple local in memory test, they will always be present, because we need them to spin up the in memory chain
	e, _, err = commonchangeset.ApplyChangesets(t, e, initialDeployCS(t, e, nil))
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, []uint64{solSelector})
	require.NoError(t, err)
}
