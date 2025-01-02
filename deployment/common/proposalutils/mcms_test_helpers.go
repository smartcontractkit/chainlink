package proposalutils

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	// "github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	// TestXXXMCMSSigner is a throwaway private key used for signing MCMS proposals.
	// in tests.
	TestXXXMCMSSigner *ecdsa.PrivateKey
)

func init() {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	TestXXXMCMSSigner = key
}

func SingleGroupMCMS(t *testing.T) config.Config {
	publicKey := TestXXXMCMSSigner.Public().(*ecdsa.PublicKey)
	// Convert the public key to an Ethereum address
	address := crypto.PubkeyToAddress(*publicKey)
	c, err := config.NewConfig(1, []common.Address{address}, []config.Config{})
	require.NoError(t, err)
	return *c
}

func SignProposal(t *testing.T, env deployment.Environment, proposal *timelock.MCMSWithTimelockProposal) *mcms.Executor {
	executorClients := make(map[mcms.ChainIdentifier]mcms.ContractDeployBackend)
	for _, chain := range env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := mcms.ChainIdentifier(chainselc.Selector)
		executorClients[chainSel] = chain.Client
	}
	executor, err := proposal.ToExecutor(true)
	require.NoError(t, err)
	payload, err := executor.SigningHash()
	require.NoError(t, err)
	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), TestXXXMCMSSigner)
	require.NoError(t, err)
	mcmSig, err := mcms.NewSignatureFromBytes(sig)
	require.NoError(t, err)
	executor.Proposal.AddSignature(mcmSig)
	require.NoError(t, executor.Proposal.Validate())
	return executor
}

func ExecuteProposal(t *testing.T, env deployment.Environment, executor *mcms.Executor,
	timelockContracts *TimelockExecutionContracts, sel uint64) {
	t.Log("Executing proposal on chain", sel)
	// Set the root.
	tx, err2 := executor.SetRootOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, mcms.ChainIdentifier(sel))
	if err2 != nil {
		require.NoError(t, deployment.MaybeDataErr(err2))
	}

	_, err2 = env.Chains[sel].Confirm(tx)
	require.NoError(t, err2)
	cfg := RunTimelockExecutorConfig{
		Executor:          executor,
		TimelockContracts: timelockContracts,
		ChainSelector:     sel,
	}
	require.NoError(t, RunTimelockExecutor(env, cfg))
}

func SingleGroupTimelockConfig(t *testing.T) commontypes.MCMSWithTimelockConfig {
	return commontypes.MCMSWithTimelockConfig{
		Canceller:        SingleGroupMCMS(t),
		Bypasser:         SingleGroupMCMS(t),
		Proposer:         SingleGroupMCMS(t),
		TimelockMinDelay: big.NewInt(0),
	}
}
