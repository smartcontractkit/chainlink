package v1_6

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_RMNRemote_Curse_View(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.BlockChains.EVMChains()[e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))[0]]
	_, tx, remote, err := rmn_remote.DeployRMNRemote(chain.DeployerKey, chain.Client, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))[0], common.Address{})
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = remote.Curse(chain.DeployerKey, globals.GlobalCurseSubject())
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	family, err := chainsel.GetSelectorFamily(e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))[0])
	require.NoError(t, err)

	tx, err = remote.Curse(chain.DeployerKey, globals.FamilyAwareSelectorToSubject(e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))[0], family))
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	view, err := GenerateRMNRemoteView(remote)
	require.NoError(t, err)

	require.True(t, view.IsCursed)
	require.Len(t, view.CursedSubjectEntries, 2)
	require.Equal(t, "01000000000000000000000000000001", view.CursedSubjectEntries[0].Subject)
	require.Equal(t, uint64(0), view.CursedSubjectEntries[0].Selector)
	require.Equal(t, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))[0], view.CursedSubjectEntries[1].Selector)
}

func Test_RMN_Selector_To_Solana_Subject(t *testing.T) {
	subject := globals.FamilyAwareSelectorToSubject(chainsel.BINANCE_SMART_CHAIN_TESTNET.Selector, chainsel.FamilySolana)
	require.Equal(t, []byte{251, 150, 143, 3, 112, 145, 21, 184, 0, 0, 0, 0, 0, 0, 0, 0}, subject[:])
}

func Test_RMN_Subject_To_Solana_Selector(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector([16]byte{251, 150, 143, 3, 112, 145, 21, 184, 0, 0, 0, 0, 0, 0, 0, 0}, chainsel.FamilySolana)
	require.Equal(t, chainsel.BINANCE_SMART_CHAIN_TESTNET.Selector, selector)
}

func Test_RMN_Selector_To_Subject(t *testing.T) {
	subject := globals.FamilyAwareSelectorToSubject(chainsel.BINANCE_SMART_CHAIN_TESTNET.Selector, chainsel.FamilyEVM)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 184, 21, 145, 112, 3, 143, 150, 251}, subject[:])
}

func Test_RMN_Subject_To_Selector(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 184, 21, 145, 112, 3, 143, 150, 251}, chainsel.FamilyEVM)
	require.Equal(t, chainsel.BINANCE_SMART_CHAIN_TESTNET.Selector, selector)
}

func Test_GlobalSubject_To_Selector(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector(globals.GlobalCurseSubject(), chainsel.FamilyEVM)
	require.Equal(t, uint64(0), selector)
}

func Test_GlobalSubject_To_Selector_Solana(t *testing.T) {
	selector := globals.FamilyAwareSubjectToSelector(globals.GlobalCurseSubject(), chainsel.FamilySolana)
	require.Equal(t, uint64(0), selector)
}
