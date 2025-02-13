package proposalutils

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
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

// Deprecated: Use SignMCMSTimelockProposal instead.
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

// Deprecated: Use ExecuteMCMSTimelockProposalV2 instead.
func ExecuteProposal(t *testing.T, env deployment.Environment, executor *mcms.Executor,
	timelockContracts *TimelockExecutionContracts, sel uint64) error {
	t.Log("Executing proposal on chain", sel)
	// Set the root.
	tx, err2 := executor.SetRootOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, mcms.ChainIdentifier(sel))
	if err2 != nil {
		require.NoError(t, deployment.MaybeDataErr(err2), "failed to set root")
	}

	_, err2 = env.Chains[sel].Confirm(tx)
	require.NoError(t, err2)
	cfg := RunTimelockExecutorConfig{
		Executor:          executor,
		TimelockContracts: timelockContracts,
		ChainSelector:     sel,
	}
	// return the error so devs can ensure expected reversions
	return RunTimelockExecutor(env, cfg)
}

// SignMCMSTimelockProposal - Signs an MCMS timelock proposal.
func SignMCMSTimelockProposal(t *testing.T, env deployment.Environment, proposal *mcmslib.TimelockProposal) *mcmslib.Proposal {
	converters := make(map[types.ChainSelector]sdk.TimelockConverter)
	inspectorsMap := make(map[types.ChainSelector]sdk.Inspector)
	for _, chain := range env.Chains {
		_, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := types.ChainSelector(chain.Selector)
		converters[chainSel] = &evm.TimelockConverter{}
		inspectorsMap[chainSel] = evm.NewInspector(chain.Client)
	}

	for _, chain := range env.SolChains {
		_, exists := chainsel.SolanaChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := types.ChainSelector(chain.Selector)
		converters[chainSel] = &solana.TimelockConverter{}
		inspectorsMap[chainSel] = solana.NewInspector(chain.Client)
	}

	p, _, err := proposal.Convert(env.GetContext(), converters)
	require.NoError(t, err)

	p.UseSimulatedBackend(true)

	signable, err := mcmslib.NewSignable(&p, inspectorsMap)
	require.NoError(t, err)

	err = signable.ValidateConfigs(env.GetContext())
	require.NoError(t, err)

	signer := mcmslib.NewPrivateKeySigner(TestXXXMCMSSigner)
	_, err = signable.SignAndAppend(signer)
	require.NoError(t, err)

	quorumMet, err := signable.ValidateSignatures(env.GetContext())
	require.NoError(t, err)
	require.True(t, quorumMet)

	return &p
}

// SignMCMSProposal - Signs an MCMS proposal. For timelock proposal, use SignMCMSTimelockProposal instead.
func SignMCMSProposal(t *testing.T, env deployment.Environment, proposal *mcmslib.Proposal) *mcmslib.Proposal {
	converters := make(map[types.ChainSelector]sdk.TimelockConverter)
	inspectorsMap := make(map[types.ChainSelector]sdk.Inspector)
	for _, chain := range env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := types.ChainSelector(chainselc.Selector)
		converters[chainSel] = &evm.TimelockConverter{}
		inspectorsMap[chainSel] = evm.NewInspector(chain.Client)
	}

	for _, chain := range env.SolChains {
		_, exists := chainsel.SolanaChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := types.ChainSelector(chain.Selector)
		converters[chainSel] = &solana.TimelockConverter{}
		inspectorsMap[chainSel] = solana.NewInspector(chain.Client)
	}

	proposal.UseSimulatedBackend(true)

	signable, err := mcmslib.NewSignable(proposal, inspectorsMap)
	require.NoError(t, err)

	err = signable.ValidateConfigs(env.GetContext())
	require.NoError(t, err)

	signer := mcmslib.NewPrivateKeySigner(TestXXXMCMSSigner)
	_, err = signable.SignAndAppend(signer)
	require.NoError(t, err)

	quorumMet, err := signable.ValidateSignatures(env.GetContext())
	require.NoError(t, err)
	require.True(t, quorumMet)

	return proposal
}

// ExecuteMCMSProposalV2 - Executes an MCMS proposal on a chain. For timelock proposal, use ExecuteMCMSTimelockProposalV2 instead.
func ExecuteMCMSProposalV2(t *testing.T, env deployment.Environment, proposal *mcmslib.Proposal) {
	t.Log("Executing proposal")

	encoders, err := proposal.GetEncoders()
	require.NoError(t, err)

	// build a map with chainSelector => executor
	executorsMap := map[types.ChainSelector]sdk.Executor{}
	for _, op := range proposal.Operations {
		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		switch family {
		case chainsel.FamilyEVM:
			encoder := encoders[op.ChainSelector].(*evm.Encoder)
			executorsMap[op.ChainSelector] = evm.NewExecutor(
				encoder,
				env.Chains[uint64(op.ChainSelector)].Client,
				env.Chains[uint64(op.ChainSelector)].DeployerKey)
		case chainsel.FamilySolana:
			encoder := encoders[op.ChainSelector].(*solana.Encoder)
			executorsMap[op.ChainSelector] = solana.NewExecutor(
				encoder,
				env.SolChains[uint64(op.ChainSelector)].Client,
				*env.SolChains[uint64(op.ChainSelector)].DeployerKey)
		default:
			require.FailNow(t, "unsupported chain family")
		}
	}

	executable, err := mcmslib.NewExecutable(proposal, executorsMap)
	require.NoError(t, err)

	// call SetRoot for each chain
	for chainSelector := range executorsMap {
		root, err := executable.SetRoot(env.GetContext(), chainSelector)
		require.NoError(t, deployment.MaybeDataErr(err))

		family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
		require.NoError(t, err)

		// no need to confirm transaction on solana as the MCMS sdk confirms it internally
		if family == chainsel.FamilyEVM {
			chain := env.Chains[uint64(chainSelector)]
			evmTransaction := root.RawTransaction.(*gethtypes.Transaction)
			_, err = chain.Confirm(evmTransaction)
			require.NoError(t, err)
		}
	}

	// execute each operation sequentially
	for i, op := range proposal.Operations {
		result, err := executable.Execute(env.GetContext(), i)
		require.NoError(t, err)

		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		if family == chainsel.FamilyEVM {
			chain := env.Chains[uint64(op.ChainSelector)]
			evmTransaction := result.RawTransaction.(*gethtypes.Transaction)
			_, err = chain.Confirm(evmTransaction)
			require.NoError(t, err)
		}
	}
}

// ExecuteMCMSTimelockProposalV2 - Includes an option to set callProxy to execute the calls through a proxy.
// If the callProxy is not set, the calls will be executed directly to the timelock.
func ExecuteMCMSTimelockProposalV2(t *testing.T, env deployment.Environment, timelockProposal *mcmslib.TimelockProposal, opts ...mcmslib.Option) {
	t.Log("Executing timelock proposal")

	// build a "chainSelector => executor" map
	executorsMap := map[types.ChainSelector]sdk.TimelockExecutor{}
	for _, op := range timelockProposal.Operations {
		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		switch family {
		case chainsel.FamilyEVM:
			executorsMap[op.ChainSelector] = evm.NewTimelockExecutor(
				env.Chains[uint64(op.ChainSelector)].Client,
				env.Chains[uint64(op.ChainSelector)].DeployerKey)
		case chainsel.FamilySolana:
			executorsMap[op.ChainSelector] = solana.NewTimelockExecutor(
				env.SolChains[uint64(op.ChainSelector)].Client,
				*env.SolChains[uint64(op.ChainSelector)].DeployerKey)
		default:
			require.FailNow(t, "unsupported chain family")
		}
	}

	timelockExecutable, err := mcmslib.NewTimelockExecutable(timelockProposal, executorsMap)
	require.NoError(t, err)

	err = timelockExecutable.IsReady(env.GetContext())
	require.NoError(t, err)

	// execute each operation sequentially
	var tx = types.TransactionResult{}
	for i, op := range timelockProposal.Operations {
		tx, err = timelockExecutable.Execute(env.GetContext(), i, opts...)
		require.NoError(t, err)

		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		// no need to confirm transaction on solana as the MCMS sdk confirms it internally
		if family == chainsel.FamilyEVM {
			chain := env.Chains[uint64(op.ChainSelector)]
			evmTransaction := tx.RawTransaction.(*gethtypes.Transaction)
			_, err = chain.Confirm(evmTransaction)
			require.NoError(t, err)
		}
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
