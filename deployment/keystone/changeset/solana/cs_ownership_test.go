package solana

import (
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_TransferOwnershipForwarder(t *testing.T) {
	// deploy forwarder
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1, // nodes unused but required in config
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	registrySel := env.AllChainSelectorsSolana()[0]
	ab := cldf.NewMemoryAddressBook()

	chain := env.SolChains[registrySel]
	mcfg := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		registrySel: {
			Canceller:        proposalutils.SingleGroupMCMSV2(t),
			Proposer:         proposalutils.SingleGroupMCMSV2(t),
			Bypasser:         proposalutils.SingleGroupMCMSV2(t),
			TimelockMinDelay: big.NewInt(0),
		},
	}
	//out, err := commonchangeset.DeployMCMSWithTimelockV2(env, mcfg)
	//env.ExistingAddresses.Merge(out.AddressBook)
	//deploy MCMSWithTimelock
	//env = deployMCMSwithTimelock(t, env, chain)
	env = shouldDeployForwarder(t, env, registrySel, ab)

	_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, env.SolChains[registrySel], env.ExistingAddresses, mcfg[registrySel])
	require.NoError(t, err)

	addresses, _ := env.ExistingAddresses.AddressesForChain(registrySel)

	mcmState, _ := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if mcmState.TimelockProgram.IsZero() {
		log.Fatal("timelock was lost??")
	}
	result, err := TransferOwnershipForwarder(env, &TransferOwnershipForwarderRequest{
		ChainSel: registrySel,
		MCMSCfg:  proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
	})

	require.NoError(t, err)
	require.NotNil(t, result.MCMSTimelockProposals)
}

func deployMCMSwithTimelock(t *testing.T, e cldf.Environment, chain cldf.SolChain) cldf.Environment {
	e = deployAccessController(t, e, chain)
	e = initAccessController(t, e, chain)
	e = deployMCM(t, e, chain)
	e = deployTimelock(t, e, chain)
	return e
}

func initAccessController(t *testing.T, e cldf.Environment, chain cldf.SolChain) cldf.Environment {
	return cldf.Environment{}
}
func deployAccessController(t *testing.T, e cldf.Environment, chain cldf.SolChain) cldf.Environment {
	deployedProgramID, err := chain.DeployProgram(e.Logger, cldf.SolProgramInfo{
		Name:  deployment.AccessControllerProgramName,
		Bytes: deployment.SolanaProgramBytes[deployment.AccessControllerProgramName],
	}, false, true)
	require.NoError(t, err)

	programID, err := solana.PublicKeyFromBase58(deployedProgramID)
	require.NoError(t, err)

	typeAndVersion := cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(chain.Selector, programID.String(), typeAndVersion)
	require.NoError(t, err)
	return e
}

func deployMCM(t *testing.T, e cldf.Environment, chain cldf.SolChain) cldf.Environment {
	deployedProgramID, err := chain.DeployProgram(e.Logger, cldf.SolProgramInfo{
		Name:  deployment.McmProgramName,
		Bytes: deployment.SolanaProgramBytes[deployment.McmProgramName],
	}, false, true)
	require.NoError(t, err)

	programID, err := solana.PublicKeyFromBase58(deployedProgramID)
	require.NoError(t, err)

	typeAndVersion := cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(chain.Selector, programID.String(), typeAndVersion)
	require.NoError(t, err)

	return e
}

func deployTimelock(t *testing.T, e cldf.Environment, chain cldf.SolChain) cldf.Environment {
	typeAndVersion := cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	deployedProgramID, err := chain.DeployProgram(e.Logger, cldf.SolProgramInfo{
		Name:  deployment.TimelockProgramName,
		Bytes: deployment.SolanaProgramBytes[deployment.TimelockProgramName],
	}, false, true)
	require.NoError(t, err)

	programID, err := solana.PublicKeyFromBase58(deployedProgramID)
	require.NoError(t, err)

	err = e.ExistingAddresses.Save(chain.Selector, programID.String(), typeAndVersion)
	require.NoError(t, err)

	return e
}
