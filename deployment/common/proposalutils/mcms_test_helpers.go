package proposalutils

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcms2 "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/types"
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

func SignMCMSTimelockProposal(t *testing.T, env deployment.Environment, proposal *mcms2.TimelockProposal) *mcms2.Proposal {
	converters := make(map[types.ChainSelector]sdk.TimelockConverter)
	inspectorsMap := make(map[types.ChainSelector]sdk.Inspector)
	for _, chain := range env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := types.ChainSelector(chainselc.Selector)
		converters[chainSel] = &evm.TimelockConverter{}
		inspectorsMap[chainSel] = evm.NewInspector(chain.Client)
	}

	p, _, err := proposal.Convert(env.GetContext(), converters)
	require.NoError(t, err)
	p.UseSimulatedBackend(true)

	signable, err := mcms2.NewSignable(&p, inspectorsMap)
	require.NoError(t, err)

	err = signable.ValidateConfigs(env.GetContext())
	require.NoError(t, err)

	signer := mcms2.NewPrivateKeySigner(TestXXXMCMSSigner)
	_, err = signable.SignAndAppend(signer)
	require.NoError(t, err)

	quorumMet, err := signable.ValidateSignatures(env.GetContext())
	require.NoError(t, err)
	require.True(t, quorumMet)

	return &p
}

func SignMCMSProposal(t *testing.T, env deployment.Environment, p *mcms2.Proposal) *mcms2.Proposal {
	converters := make(map[types.ChainSelector]sdk.TimelockConverter)
	inspectorsMap := make(map[types.ChainSelector]sdk.Inspector)
	for _, chain := range env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := types.ChainSelector(chainselc.Selector)
		converters[chainSel] = &evm.TimelockConverter{}
		inspectorsMap[chainSel] = evm.NewInspector(chain.Client)
	}

	signable, err := mcms2.NewSignable(p, inspectorsMap)
	require.NoError(t, err)

	err = signable.ValidateConfigs(env.GetContext())
	require.NoError(t, err)

	signer := mcms2.NewPrivateKeySigner(TestXXXMCMSSigner)
	_, err = signable.SignAndAppend(signer)
	require.NoError(t, err)

	quorumMet, err := signable.ValidateSignatures(env.GetContext())
	require.NoError(t, err)
	require.True(t, quorumMet)

	return p
}

func ExecuteProposalV2(t *testing.T, env deployment.Environment, proposal *mcms2.Proposal, sel uint64) {
	t.Log("Executing proposal on chain", sel)

	encoders, err := proposal.GetEncoders()
	require.NoError(t, err)

	selector := types.ChainSelector(sel)
	encoder := encoders[selector].(*evm.Encoder)
	evmExecutor := evm.NewExecutor(encoder, env.Chains[sel].Client, env.Chains[sel].DeployerKey)
	executorsMap := map[types.ChainSelector]sdk.Executor{
		selector: evmExecutor,
	}
	executable, err := mcms2.NewExecutable(proposal, executorsMap)
	require.NoError(t, err)

	chain := env.Chains[sel]
	root, err := executable.SetRoot(env.GetContext(), selector)
	require.NoError(t, deployment.MaybeDataErr(err))

	evmTransaction := root.RawTransaction.(*types2.Transaction)
	_, err = chain.Confirm(evmTransaction)
	require.NoError(t, err)

	for i := 0; i < len(proposal.Operations); i++ {
		result, err := executable.Execute(env.GetContext(), i)
		require.NoError(t, err)

		evmTransaction = result.RawTransaction.(*types2.Transaction)
		_, err = chain.Confirm(evmTransaction)
		require.NoError(t, err)
	}
}

func SingleGroupTimelockConfig(t *testing.T) commontypes.MCMSWithTimelockConfig {
	return commontypes.MCMSWithTimelockConfig{
		Canceller:        SingleGroupMCMS(t),
		Bypasser:         SingleGroupMCMS(t),
		Proposer:         SingleGroupMCMS(t),
		TimelockMinDelay: big.NewInt(0),
	}
}
