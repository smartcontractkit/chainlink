package v1_6

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

func Test_RMNRemote_Curse_View(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[selector]

	_, tx, remote, err := rmn_remote.DeployRMNRemote(chain.DeployerKey, chain.Client, selector, common.Address{})
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = remote.Curse(chain.DeployerKey, globals.GlobalCurseSubject())
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = remote.Curse(chain.DeployerKey, globals.FamilyAwareSelectorToSubject(selector, chain.Family()))
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	view, err := GenerateRMNRemoteView(remote)
	require.NoError(t, err)

	require.True(t, view.IsCursed)
	require.Len(t, view.CursedSubjectEntries, 2)
	require.Equal(t, "01000000000000000000000000000001", view.CursedSubjectEntries[0].Subject)
	require.Equal(t, uint64(0), view.CursedSubjectEntries[0].Selector)
	require.Equal(t, selector, view.CursedSubjectEntries[1].Selector)
}

func Test_RMN_Selector_To_Solana_Subject(t *testing.T) {
	subject := globals.FamilyAwareSelectorToSubject(chainselectors.BINANCE_SMART_CHAIN_TESTNET.Selector, chainselectors.FamilySolana)
	require.Equal(t, []byte{251, 150, 143, 3, 112, 145, 21, 184, 0, 0, 0, 0, 0, 0, 0, 0}, subject[:])
}

func Test_RMN_Subject_To_Solana_Selector(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector([16]byte{251, 150, 143, 3, 112, 145, 21, 184, 0, 0, 0, 0, 0, 0, 0, 0}, chainselectors.FamilySolana)
	require.Equal(t, chainselectors.BINANCE_SMART_CHAIN_TESTNET.Selector, selector)
}

func Test_RMN_Selector_To_Subject(t *testing.T) {
	subject := globals.FamilyAwareSelectorToSubject(chainselectors.BINANCE_SMART_CHAIN_TESTNET.Selector, chainselectors.FamilyEVM)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 184, 21, 145, 112, 3, 143, 150, 251}, subject[:])
}

func Test_RMN_Subject_To_Selector(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 184, 21, 145, 112, 3, 143, 150, 251}, chainselectors.FamilyEVM)
	require.Equal(t, chainselectors.BINANCE_SMART_CHAIN_TESTNET.Selector, selector)
}

func Test_GlobalSubject_To_Selector(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector(globals.GlobalCurseSubject(), chainselectors.FamilyEVM)
	require.Equal(t, uint64(0), selector)
}

func Test_GlobalSubject_To_Selector_Solana(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector(globals.GlobalCurseSubject(), chainselectors.FamilySolana)
	require.Equal(t, uint64(0), selector)
}
