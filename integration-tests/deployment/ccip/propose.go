package ccipdeployment

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

var (
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

func NewTestMCMSConfig(t *testing.T, e deployment.Environment) MCMSConfig {
	c := SingleGroupMCMS(t)
	// All deployer keys can execute.
	var executors []common.Address
	for _, chain := range e.Chains {
		executors = append(executors, chain.DeployerKey.From)
	}
	return MCMSConfig{
		Admin:     c,
		Bypasser:  c,
		Canceller: c,
		Executors: executors,
		Proposer:  c,
	}
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
	state CCIPOnChainState, sel uint64) {
	t.Log("Executing proposal on chain", sel)
	// Set the root.
	tx, err2 := executor.SetRootOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, mcms.ChainIdentifier(sel))
	require.NoError(t, err2)
	_, err2 = env.Chains[sel].Confirm(tx)
	require.NoError(t, err2)

	// TODO: This sort of helper probably should move to the MCMS lib.
	// Execute all the transactions in the proposal which are for this chain.
	for _, chainOp := range executor.Operations[mcms.ChainIdentifier(sel)] {
		for idx, op := range executor.ChainAgnosticOps {
			if bytes.Equal(op.Data, chainOp.Data) && op.To == chainOp.To {
				opTx, err3 := executor.ExecuteOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, idx)
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
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	timelockAddresses := make(map[mcms.ChainIdentifier]common.Address)
	for _, sel := range chains {
		chain, _ := chainsel.ChainBySelector(sel)
		acceptOnRamp, err := state.Chains[sel].OnRamp.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return nil, err
		}
		acceptFeeQuoter, err := state.Chains[sel].FeeQuoter.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return nil, err
		}
		chainSel := mcms.ChainIdentifier(chain.Selector)
		opCount, err := state.Chains[sel].ProposerMcm.GetOpCount(nil)
		if err != nil {
			return nil, err
		}
		metaDataPerChain[chainSel] = mcms.ChainMetadata{
			MCMAddress:      state.Chains[sel].ProposerMcm.Address(),
			StartingOpCount: opCount.Uint64(),
		}
		timelockAddresses[chainSel] = state.Chains[sel].Timelock.Address()
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

	acceptCR, err := state.Chains[homeChain].CapabilityRegistry.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return nil, err
	}
	acceptCCIPConfig, err := state.Chains[homeChain].CCIPHome.AcceptOwnership(deployment.SimTransactOpts())
	if err != nil {
		return nil, err
	}
	homeChainID := mcms.ChainIdentifier(homeChain)
	opCount, err := state.Chains[homeChain].ProposerMcm.GetOpCount(nil)
	if err != nil {
		return nil, err
	}
	metaDataPerChain[homeChainID] = mcms.ChainMetadata{
		StartingOpCount: opCount.Uint64(),
		MCMAddress:      state.Chains[homeChain].ProposerMcm.Address(),
	}
	timelockAddresses[homeChainID] = state.Chains[homeChain].Timelock.Address()
	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: homeChainID,
		Batch: []mcms.Operation{
			{
				To:    state.Chains[homeChain].CapabilityRegistry.Address(),
				Data:  acceptCR.Data(),
				Value: big.NewInt(0),
			},
			{
				To:    state.Chains[homeChain].CCIPHome.Address(),
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
		timelockAddresses,
		"blah", // TODO
		batches,
		timelock.Schedule, "0s")
}

func BuildProposalMetadata(state CCIPOnChainState, chains []uint64) (map[mcms.ChainIdentifier]common.Address, map[mcms.ChainIdentifier]mcms.ChainMetadata, error) {
	tlAddressMap := make(map[mcms.ChainIdentifier]common.Address)
	metaDataPerChain := make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	for _, sel := range chains {
		chainId := mcms.ChainIdentifier(sel)
		tlAddressMap[chainId] = state.Chains[sel].Timelock.Address()
		mcm := state.Chains[sel].ProposerMcm
		opCount, err := mcm.GetOpCount(nil)
		if err != nil {
			return nil, nil, err
		}
		metaDataPerChain[chainId] = mcms.ChainMetadata{
			StartingOpCount: opCount.Uint64(),
			MCMAddress:      mcm.Address(),
		}
	}
	return tlAddressMap, metaDataPerChain, nil
}
