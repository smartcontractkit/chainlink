package v1_6

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_RMNRemote_Curse_View(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.Chains[e.AllChainSelectors()[0]]
	_, tx, remote, err := rmn_remote.DeployRMNRemote(chain.DeployerKey, chain.Client, e.AllChainSelectors()[0], common.Address{})
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = remote.Curse(chain.DeployerKey, globals.GlobalCurseSubject())
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = remote.Curse(chain.DeployerKey, globals.SelectorToSubject(e.AllChainSelectors()[0]))
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	view, err := GenerateRMNRemoteView(remote)
	require.NoError(t, err)

	require.True(t, view.IsCursed)
	require.Len(t, view.CursedSubjectEntries, 2)
	require.Equal(t, "01000000000000000000000000000001", view.CursedSubjectEntries[0].Subject)
	require.Equal(t, uint64(0), view.CursedSubjectEntries[0].Selector)
	require.Equal(t, e.AllChainSelectors()[0], view.CursedSubjectEntries[1].Selector)
}
