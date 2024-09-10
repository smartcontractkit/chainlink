package ccipdeployment

import (
	"bytes"
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAddChain(t *testing.T) {
	// 4 chains where the 4th is added after initial deployment.
	e := NewEnvironmentWithCRAndJobs(t, logger.TestLogger(t), 4)
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	sels := e.Env.AllChainSelectors()
	initialDeploy := sels[0:3]
	newChain := sels[3]

	ab, err := DeployCCIPContracts(e.Env, DeployCCIPContractConfig{
		HomeChainSel:     e.HomeChainSel,
		ChainsToDeploy:   initialDeploy,
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(ab))
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// Contracts deployed and initial DONs set up.
	// Connect all the lanes
	for _, source := range initialDeploy {
		for _, dest := range initialDeploy {
			if source != dest {
				require.NoError(t, AddLane(e.Env, state, uint64(source), uint64(dest)))
			}
		}
	}

	executorClients := make(map[mcms.ChainIdentifier]mcms.ContractDeployBackend)
	for _, chain := range e.Env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := mcms.ChainIdentifier(chainselc.Selector)
		executorClients[chainSel] = chain.Client
	}
	//  Deploy contracts to new chain
	newAddresses, err := DeployChainContracts(e.Env, e.Env.Chains[newChain], deployment.NewMemoryAddressBook())
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(newAddresses))
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// We can directly enable the sources on the new chain with deployer key.
	var offRampEnables []offramp.OffRampSourceChainConfigArgs
	for _, source := range initialDeploy {
		offRampEnables = append(offRampEnables, offramp.OffRampSourceChainConfigArgs{
			Router:              state.Chains[newChain].Router.Address(),
			SourceChainSelector: source,
			IsEnabled:           true,
			OnRamp:              common.LeftPadBytes(state.Chains[source].OnRamp.Address().Bytes(), 32),
		})
	}
	tx, err := state.Chains[newChain].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[newChain].DeployerKey, offRampEnables)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[newChain], tx, err)
	require.NoError(t, err)

	// Generate and sign inbound proposal to new 4th chain.
	proposals, err := NewChainInboundProposal(e.Env, state, e.HomeChainSel, newChain, initialDeploy)
	require.NoError(t, err)
	realProposal, err := proposals[0].ToMCMSOnlyProposal()
	require.NoError(t, err)
	executor, err := realProposal.ToExecutor(executorClients)
	require.NoError(t, err)
	payload, err := executor.SigningHash()
	require.NoError(t, err)
	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), TestXXXMCMSSigner)
	require.NoError(t, err)
	mcmSig, err := mcms.NewSignatureFromBytes(sig)
	executor.Proposal.Signatures = append(executor.Proposal.Signatures, mcmSig)
	require.NoError(t, executor.Proposal.Validate())

	// Apply the proposal to all the chains.
	for _, sel := range sels {
		if sel == newChain {
			continue
		}
		// Set the root.
		tx, err2 := executor.SetRootOnChain(e.Env.Chains[sel].DeployerKey, mcms.ChainIdentifier(sel))
		require.NoError(t, err2)
		_, err2 = e.Env.Chains[sel].Confirm(tx.Hash())
		require.NoError(t, err2)

		// Execute all the transactions in the proposal which are for this chain.
		for _, chainOp := range executor.Operations[mcms.ChainIdentifier(sel)] {
			for idx, op := range executor.ChainAgnosticOps {
				if bytes.Equal(op.Data, chainOp.Data) && op.To == chainOp.To {
					opTx, err3 := executor.ExecuteOnChain(e.Env.Chains[sel].DeployerKey, idx)
					require.NoError(t, err3)
					block, err3 := e.Env.Chains[sel].Confirm(opTx.Hash())
					require.NoError(t, err3)
					t.Log("executed", chainOp)
					it, err3 := state.Chains[sel].Timelock.FilterCallScheduled(&bind.FilterOpts{
						Start:   block,
						End:     &block,
						Context: context.Background(),
					}, nil, nil)
					require.NoError(t, err3)
					var calls []owner_helpers.RBACTimelockCall
					var pred, salt [32]byte
					for it.Next() {
						// Note these are the same for the whole batch, can overwrite
						pred = it.Event.Predecessor
						salt = it.Event.Salt
						t.Log("scheduled", it.Event)
						calls = append(calls, owner_helpers.RBACTimelockCall{
							Target: it.Event.Target,
							Data:   it.Event.Data,
							Value:  it.Event.Value,
						})
					}
					tx, err := state.Chains[sel].Timelock.ExecuteBatch(
						e.Env.Chains[sel].DeployerKey, calls, pred, salt)
					require.NoError(t, err)
					_, err = e.Env.Chains[sel].Confirm(tx.Hash())
					require.NoError(t, err)
				}
			}
		}
	}

	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, chain := range initialDeploy {
		cfg, err2 := state.Chains[chain].OnRamp.GetDestChainConfig(nil, newChain)
		require.NoError(t, err2)
		t.Log("config", cfg)
		s, err2 := state.Chains[newChain].OffRamp.GetSourceChainConfig(nil, chain)
		require.NoError(t, err2)
		t.Log("config", s)
	}

	// Now that the proposal has been executed we expect to be able to send traffic to this new 4th chain.
	//SendRequest(t, e.Env, state, initialDeploy[0], newChain)
	//e.Env.initialDeploy[0]

}
