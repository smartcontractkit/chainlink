package proposalutils

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmschainwrappers "github.com/smartcontractkit/mcms/chainwrappers"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldfmcmsadapters "github.com/smartcontractkit/chainlink-deployments-framework/chain/mcms/adapters"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TestXXXMCMSSigner is a throwaway private key used for signing MCMS proposals.
// in tests.
var TestXXXMCMSSigner *ecdsa.PrivateKey

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

func SingleGroupMCMSV2(t *testing.T) mcmstypes.Config {
	publicKey := TestXXXMCMSSigner.Public().(*ecdsa.PublicKey)
	// Convert the public key to an Ethereum address
	address := crypto.PubkeyToAddress(*publicKey)
	c, err := mcmstypes.NewConfig(1, []common.Address{address}, []mcmstypes.Config{})
	require.NoError(t, err)
	return c
}

// SignMCMSTimelockProposal - Signs an MCMS timelock proposal.
func SignMCMSTimelockProposal(t *testing.T, env cldf.Environment, proposal *mcmslib.TimelockProposal, realBackend bool) *mcmslib.Proposal {
	converters, err := mcmschainwrappers.BuildConverters(proposal.ChainMetadata)
	require.NoError(t, err)

	mcmsChains := cldfmcmsadapters.Wrap(env.BlockChains)
	inspectors, err := mcmschainwrappers.BuildInspectors(&mcmsChains, proposal.ChainMetadata, proposal.Action)
	require.NoError(t, err)

	p, _, err := proposal.Convert(env.GetContext(), converters)
	require.NoError(t, err)

	p.UseSimulatedBackend(!realBackend)

	signable, err := mcmslib.NewSignable(&p, inspectors)
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
func SignMCMSProposal(t *testing.T, env cldf.Environment, proposal *mcmslib.Proposal) *mcmslib.Proposal {
	mcmsChains := cldfmcmsadapters.Wrap(env.BlockChains)
	inspectors, err := mcmschainwrappers.BuildInspectors(&mcmsChains, proposal.ChainMetadata, mcmstypes.TimelockActionSchedule)
	require.NoError(t, err)

	proposal.UseSimulatedBackend(true)

	signable, err := mcmslib.NewSignable(proposal, inspectors)
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
func ExecuteMCMSProposalV2(t *testing.T, env cldf.Environment, proposal *mcmslib.Proposal) error {
	t.Log("Executing proposal")

	encoders, err := proposal.GetEncoders()
	require.NoError(t, err, "[ExecuteMCMSProposalV2] failed to get encoders")

	mcmsChains := cldfmcmsadapters.Wrap(env.BlockChains)
	executors, err := mcmschainwrappers.BuildExecutors(&mcmsChains, proposal.ChainMetadata, encoders, mcmstypes.TimelockActionSchedule)
	require.NoError(t, err)

	executable, err := mcmslib.NewExecutable(proposal, executors)
	require.NoError(t, err, "[ExecuteMCMSProposalV2] failed to build executable")

	// call SetRoot for each chain
	for chainSelector := range executors {
		t.Logf("[ExecuteMCMSProposalV2] Setting root on chain %d...", chainSelector)
		root, err := executable.SetRoot(env.GetContext(), chainSelector)
		if err != nil {
			return fmt.Errorf("[ExecuteMCMSProposalV2] SetRoot failed: %w", err)
		}

		family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
		require.NoError(t, err)

		// no need to confirm transaction on solana as the MCMS sdk confirms it internally
		if family == chainsel.FamilyEVM {
			chain := env.BlockChains.EVMChains()[uint64(chainSelector)]
			evmTransaction := root.RawData.(*gethtypes.Transaction)
			t.Logf("[ExecuteMCMSProposalV2] SetRoot EVM tx hash: %s", evmTransaction.Hash().String())
			_, err = chain.Confirm(evmTransaction)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSProposalV2] Confirm failed: %w", err)
			}
		}
		if family == chainsel.FamilyAptos {
			chain := env.BlockChains.AptosChains()[uint64(chainSelector)]
			t.Logf("[ExecuteMCMSProposalV2] SetRoot Aptos tx hash: %s", root.Hash)
			err = chain.Confirm(root.Hash)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSProposalV2] Confirm failed: %w", err)
			}
		}
	}

	// execute each operation sequentially
	for i, op := range proposal.Operations {
		t.Logf("[ExecuteMCMSProposalV2] Executing operation index=%d on chain %d...", i, uint64(op.ChainSelector))
		result, err := executable.Execute(env.GetContext(), i)
		if err != nil {
			return fmt.Errorf("[ExecuteMCMSProposalV2] Execute failed: %w", err)
		}

		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		if family == chainsel.FamilyEVM {
			chain := env.BlockChains.EVMChains()[uint64(op.ChainSelector)]
			evmTransaction := result.RawData.(*gethtypes.Transaction)
			t.Logf("[ExecuteMCMSProposalV2] Operation %d EVM tx hash: %s", i, evmTransaction.Hash().String())
			_, err = chain.Confirm(evmTransaction)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSProposalV2] Confirm failed: %w", err)
			}
		}
		if family == chainsel.FamilyAptos {
			chain := env.BlockChains.AptosChains()[uint64(op.ChainSelector)]
			t.Logf("[ExecuteMCMSProposalV2] Operation %d Aptos tx hash: %s", i, result.Hash)
			err = chain.Confirm(result.Hash)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSProposalV2] Confirm failed: %w", err)
			}
		}
	}

	return nil
}

// ExecuteMCMSTimelockProposalV2 - Includes an option to set callProxy to execute the calls through a proxy.
// If the callProxy is not set, the calls will be executed directly to the timelock.
func ExecuteMCMSTimelockProposalV2(t *testing.T, env cldf.Environment, timelockProposal *mcmslib.TimelockProposal, opts ...mcmslib.Option) error {
	t.Log("Executing timelock proposal")

	mcmsChains := cldfmcmsadapters.Wrap(env.BlockChains)
	executors, err := mcmschainwrappers.BuildTimelockExecutors(&mcmsChains, timelockProposal.ChainMetadata,
		timelockProposal.Action)
	require.NoError(t, err)

	timelockExecutable, err := mcmslib.NewTimelockExecutable(env.GetContext(), timelockProposal, executors)
	require.NoError(t, err)

	isReady := func() error {
		err := timelockExecutable.IsReady(env.GetContext())
		return err
	}
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.NoErrorf(collect, isReady(), "Proposal is not ready")
	}, 100*time.Second, 50*time.Millisecond, "timelock proposal not ready after 100s")

	// execute each operation sequentially
	tx := mcmstypes.TransactionResult{}
	for i, op := range timelockProposal.Operations {
		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		opOpts := slices.Clone(opts)
		if family == chainsel.FamilyEVM {
			callProxy := findCallProxyAddress(t, env, uint64(op.ChainSelector), timelockProposal.TimelockAddresses[op.ChainSelector])
			opOpts = append(opOpts, mcmslib.WithCallProxy(callProxy))
			t.Logf("[ExecuteMCMSTimelockProposalV2] Using EVM chain with chainID=%d, timelock address %s call proxy %s",
				uint64(op.ChainSelector), timelockProposal.TimelockAddresses[op.ChainSelector], callProxy)
		}

		tx, err = timelockExecutable.Execute(env.GetContext(), i, opOpts...)
		if err != nil {
			return fmt.Errorf("[ExecuteMCMSTimelockProposalV2] Execute failed: %w", err)
		}
		t.Logf("[ExecuteMCMSTimelockProposalV2] Executed timelock operation index=%d on chain %d (tx %v)", i, uint64(op.ChainSelector), tx.Hash)

		// no need to confirm transaction on solana as the MCMS sdk confirms it internally
		if family == chainsel.FamilyEVM {
			chain := env.BlockChains.EVMChains()[uint64(op.ChainSelector)]
			evmTransaction := tx.RawData.(*gethtypes.Transaction)
			_, err = chain.Confirm(evmTransaction)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSTimelockProposalV2] Confirm on EVM failed: %w", err)
			}
		}
		if family == chainsel.FamilyAptos {
			chain := env.BlockChains.AptosChains()[uint64(op.ChainSelector)]
			err = chain.Confirm(tx.Hash)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSTimelockProposalV2] Confirm on Aptos failed: %w", err)
			}
		}
	}

	return nil
}

func SingleGroupTimelockConfig(t *testing.T) commontypes.MCMSWithTimelockConfig {
	return commontypes.MCMSWithTimelockConfig{
		Canceller:        SingleGroupMCMS(t),
		Bypasser:         SingleGroupMCMS(t),
		Proposer:         SingleGroupMCMS(t),
		TimelockMinDelay: big.NewInt(0),
	}
}

func SingleGroupTimelockConfigV2(t *testing.T) commontypes.MCMSWithTimelockConfigV2 {
	return commontypes.MCMSWithTimelockConfigV2{
		Canceller:        SingleGroupMCMSV2(t),
		Bypasser:         SingleGroupMCMSV2(t),
		Proposer:         SingleGroupMCMSV2(t),
		TimelockMinDelay: big.NewInt(0),
	}
}

func findCallProxyAddress(t *testing.T, env cldf.Environment, chainSelector uint64, timelockAddr string) string {
	timelock, err := bindings.NewRBACTimelock(common.HexToAddress(timelockAddr), env.BlockChains.EVMChains()[chainSelector].Client)
	require.NoError(t, err)
	role, err := timelock.EXECUTORROLE(&bind.CallOpts{
		Context: env.GetContext(),
	})
	require.NoError(t, err)
	addr, err := timelock.GetRoleMember(&bind.CallOpts{
		Context: env.GetContext(),
	}, role, big.NewInt(0)) // we expect only one member in the executor role
	require.NoError(t, err)
	require.NotEqual(t, common.Address{}, addr, "executor role has no members; is the timelock initialized?")
	return addr.Hex()
}
