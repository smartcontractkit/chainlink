package ccipdeployment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployCCIPContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      4,
	})
	// Deploy all the CCIP contracts.
	homeChain := e.AllChainSelectors()[0]
	capRegAddresses, _, err := DeployCapReg(lggr, e.Chains, homeChain)
	require.NoError(t, err)
	s, err := LoadOnchainState(e, capRegAddresses)
	require.NoError(t, err)
	ab, err := DeployCCIPContracts(e, DeployCCIPContractConfig{
		HomeChainSel:     homeChain,
		ChainsToDeploy:   e.AllChainSelectors(),
		CCIPOnChainState: s,
	})
	require.NoError(t, err)
	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.Snapshot(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func TestJobSpecGeneration(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	js, err := NewCCIPJobSpecs(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	for node, jb := range js {
		fmt.Println(node, jb)
	}
	// TODO: Add job assertions
}
