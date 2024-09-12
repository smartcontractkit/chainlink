package ccipdeployment

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

// TODO: Pull up to deploy
func SimTransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		return transaction, nil
	}, From: common.HexToAddress("0x0"), NoSend: true, GasLimit: 1_000_000}
}

func SignProposal(t *testing.T, env deployment.Environment, proposal *timelock.MCMSWithTimelockProposal) *mcms.Executor {
	executorClients := make(map[mcms.ChainIdentifier]mcms.ContractDeployBackend)
	for _, chain := range env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := mcms.ChainIdentifier(chainselc.Selector)
		executorClients[chainSel] = chain.Client
	}
	realProposal, err := proposal.ToMCMSOnlyProposal()
	require.NoError(t, err)
	executor, err := realProposal.ToExecutor(executorClients)
	require.NoError(t, err)
	payload, err := executor.SigningHash()
	require.NoError(t, err)
	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), TestXXXMCMSSigner)
	require.NoError(t, err)
	mcmSig, err := mcms.NewSignatureFromBytes(sig)
	require.NoError(t, err)
	executor.Proposal.Signatures = append(executor.Proposal.Signatures, mcmSig)
	require.NoError(t, executor.Proposal.Validate())
	return executor
}

func ExecuteProposal(t *testing.T, env deployment.Environment, executor *mcms.Executor,
	state CCIPOnChainState, sel uint64) {
	// Set the root.
	tx, err2 := executor.SetRootOnChain(env.Chains[sel].DeployerKey, mcms.ChainIdentifier(sel))
	require.NoError(t, err2)
	_, err2 = env.Chains[sel].Confirm(tx)
	require.NoError(t, err2)

	// TODO: This sort of helper probably should move to the MCMS lib.
	// Execute all the transactions in the proposal which are for this chain.
	for _, chainOp := range executor.Operations[mcms.ChainIdentifier(sel)] {
		for idx, op := range executor.ChainAgnosticOps {
			if bytes.Equal(op.Data, chainOp.Data) && op.To == chainOp.To {
				opTx, err3 := executor.ExecuteOnChain(env.Chains[sel].DeployerKey, idx)
				require.NoError(t, err3)
				block, err3 := env.Chains[sel].Confirm(opTx)
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
					env.Chains[sel].DeployerKey, calls, pred, salt)
				require.NoError(t, err)
				_, err = env.Chains[sel].Confirm(tx)
				require.NoError(t, err)
			}
		}
	}
}

func GenerateAcceptOwnershipProposal(
	state CCIPOnChainState,
	homeChain uint64,
	chains []uint64,
) (*timelock.MCMSWithTimelockProposal, error) {
	// TODO: Accept rest of contracts
	var batches []timelock.BatchChainOperation
	metaDataPerChain := make(map[mcms.ChainIdentifier]timelock.MCMSWithTimelockChainMetadata)
	for _, sel := range chains {
		chain, _ := chainsel.ChainBySelector(sel)
		acceptOnRamp, err := state.Chains[sel].OnRamp.AcceptOwnership(SimTransactOpts())
		if err != nil {
			return nil, err
		}
		acceptFeeQuoter, err := state.Chains[sel].FeeQuoter.AcceptOwnership(SimTransactOpts())
		if err != nil {
			return nil, err
		}
		chainSel := mcms.ChainIdentifier(chain.Selector)
		metaDataPerChain[chainSel] = timelock.MCMSWithTimelockChainMetadata{
			ChainMetadata: mcms.ChainMetadata{
				NonceOffset: 0,
				MCMAddress:  state.Chains[sel].Mcm.Address(),
			},
			TimelockAddress: state.Chains[sel].Timelock.Address(),
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: chainSel,
			Batch: []mcms.Operation{
				{
					To:    state.Chains[sel].OnRamp.Address(),
					Data:  acceptOnRamp.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[sel].FeeQuoter.Address(),
					Data:  acceptFeeQuoter.Data(),
					Value: big.NewInt(0),
				},
			},
		})
	}
	acceptCR, err := state.Chains[homeChain].CapabilityRegistry.AcceptOwnership(SimTransactOpts())
	if err != nil {
		return nil, err
	}
	acceptCCIPConfig, err := state.Chains[homeChain].CCIPConfig.AcceptOwnership(SimTransactOpts())
	if err != nil {
		return nil, err
	}
	homeChainID := mcms.ChainIdentifier(homeChain)
	metaDataPerChain[homeChainID] = timelock.MCMSWithTimelockChainMetadata{
		ChainMetadata: mcms.ChainMetadata{
			NonceOffset: 0,
			MCMAddress:  state.Chains[homeChain].Mcm.Address(),
		},
		TimelockAddress: state.Chains[homeChain].Timelock.Address(),
	}
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: homeChainID,
		Batch: []mcms.Operation{
			{
				To:    state.Chains[homeChain].CapabilityRegistry.Address(),
				Data:  acceptCR.Data(),
				Value: big.NewInt(0),
			},
			{
				To:    state.Chains[homeChain].CCIPConfig.Address(),
				Data:  acceptCCIPConfig.Data(),
				Value: big.NewInt(0),
			},
		},
	})
	return timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO
		[]mcms.Signature{},
		false,
		metaDataPerChain,
		"blah", // TODO
		batches,
		timelock.Schedule, "0s")
}
