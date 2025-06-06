package proposalutils

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"slices"
	"testing"
	"time"

	aptosapi "github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsaptossdk "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

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

func SingleGroupMCMSV2(t *testing.T) mcmstypes.Config {
	publicKey := TestXXXMCMSSigner.Public().(*ecdsa.PublicKey)
	// Convert the public key to an Ethereum address
	address := crypto.PubkeyToAddress(*publicKey)
	c, err := mcmstypes.NewConfig(1, []common.Address{address}, []mcmstypes.Config{})
	require.NoError(t, err)
	return c
}

// SignMCMSTimelockProposal - Signs an MCMS timelock proposal.
func SignMCMSTimelockProposal(t *testing.T, env cldf.Environment, proposal *mcmslib.TimelockProposal) *mcmslib.Proposal {
	converters := make(map[mcmstypes.ChainSelector]mcmssdk.TimelockConverter)
	inspectorsMap := make(map[mcmstypes.ChainSelector]mcmssdk.Inspector)
	evmChains := env.BlockChains.EVMChains()
	solanaChains := env.BlockChains.SolanaChains()
	for _, chain := range evmChains {
		_, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := mcmstypes.ChainSelector(chain.Selector)
		converters[chainSel] = &mcmsevmsdk.TimelockConverter{}
		inspectorsMap[chainSel] = mcmsevmsdk.NewInspector(chain.Client)
	}
	for chainSelector, chain := range solanaChains {
		_, err := chainsel.SolanaChainIdFromSelector(chainSelector)
		require.NoError(t, err)
		chainSel := mcmstypes.ChainSelector(chainSelector)
		converters[chainSel] = mcmssolanasdk.TimelockConverter{}
		inspectorsMap[chainSel] = mcmssolanasdk.NewInspector(chain.Client)
	}
	for chainSelector, chain := range env.BlockChains.AptosChains() {
		_, err := chainsel.AptosChainIdFromSelector(chainSelector)
		require.NoError(t, err)
		chainSel := mcmstypes.ChainSelector(chainSelector)
		converters[chainSel] = mcmsaptossdk.NewTimelockConverter()
		roleFromAction := map[mcmstypes.TimelockAction]mcmsaptossdk.TimelockRole{
			mcmstypes.TimelockActionSchedule: mcmsaptossdk.TimelockRoleProposer,
			mcmstypes.TimelockActionBypass:   mcmsaptossdk.TimelockRoleBypasser,
			mcmstypes.TimelockActionCancel:   mcmsaptossdk.TimelockRoleCanceller,
		}
		inspectorsMap[chainSel] = mcmsaptossdk.NewInspector(chain.Client, roleFromAction[proposal.Action])
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
func SignMCMSProposal(t *testing.T, env cldf.Environment, proposal *mcmslib.Proposal) *mcmslib.Proposal {
	converters := make(map[mcmstypes.ChainSelector]mcmssdk.TimelockConverter)
	inspectorsMap := make(map[mcmstypes.ChainSelector]mcmssdk.Inspector)
	evmChains := env.BlockChains.EVMChains()
	solanaChains := env.BlockChains.SolanaChains()
	for _, chain := range evmChains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := mcmstypes.ChainSelector(chainselc.Selector)
		converters[chainSel] = &mcmsevmsdk.TimelockConverter{}
		inspectorsMap[chainSel] = mcmsevmsdk.NewInspector(chain.Client)
	}

	for _, chain := range solanaChains {
		_, exists := chainsel.SolanaChainBySelector(chain.Selector)
		require.True(t, exists)
		chainSel := mcmstypes.ChainSelector(chain.Selector)
		converters[chainSel] = &mcmssolanasdk.TimelockConverter{}
		inspectorsMap[chainSel] = mcmssolanasdk.NewInspector(chain.Client)
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
func ExecuteMCMSProposalV2(t *testing.T, env cldf.Environment, proposal *mcmslib.Proposal) error {
	t.Log("Executing proposal")

	encoders, err := proposal.GetEncoders()
	require.NoError(t, err, "[ExecuteMCMSProposalV2] failed to get encoders")

	// build a map with chainSelector => executor
	executorsMap := map[mcmstypes.ChainSelector]mcmssdk.Executor{}
	aptosChains := env.BlockChains.AptosChains()
	evmChains := env.BlockChains.EVMChains()
	solChains := env.BlockChains.SolanaChains()
	for _, op := range proposal.Operations {
		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		switch family {
		case chainsel.FamilyEVM:
			encoder := encoders[op.ChainSelector].(*mcmsevmsdk.Encoder)
			executorsMap[op.ChainSelector] = mcmsevmsdk.NewExecutor(
				encoder,
				evmChains[uint64(op.ChainSelector)].Client,
				evmChains[uint64(op.ChainSelector)].DeployerKey)
			t.Logf("[ExecuteMCMSProposalV2] Using EVM chain with chainID=%d", uint64(op.ChainSelector))
		case chainsel.FamilySolana:
			encoder := encoders[op.ChainSelector].(*mcmssolanasdk.Encoder)
			executorsMap[op.ChainSelector] = mcmssolanasdk.NewExecutor(
				encoder,
				solChains[uint64(op.ChainSelector)].Client,
				*solChains[uint64(op.ChainSelector)].DeployerKey)
			t.Logf("[ExecuteMCMSProposalV2] Using Solana chain with chainID=%d. RPC=%s. Authority=%s",
				uint64(op.ChainSelector),
				solChains[uint64(op.ChainSelector)].URL,
				solChains[uint64(op.ChainSelector)].DeployerKey.PublicKey().String(),
			)
		case chainsel.FamilyAptos:
			encoder := encoders[op.ChainSelector].(*mcmsaptossdk.Encoder)
			executorsMap[op.ChainSelector] = mcmsaptossdk.NewExecutor(
				aptosChains[uint64(op.ChainSelector)].Client,
				aptosChains[uint64(op.ChainSelector)].DeployerSigner,
				encoder,
				mcmsaptossdk.TimelockRoleProposer,
			)
			t.Logf("[ExecuteMCMSProposalV2] Using Aptos chain with chainSelector=%d", uint64(op.ChainSelector))

		default:
			require.FailNow(t, "unsupported chain family")
		}
	}

	executable, err := mcmslib.NewExecutable(proposal, executorsMap)
	require.NoError(t, err, "[ExecuteMCMSProposalV2] failed to build executable")

	// call SetRoot for each chain
	for chainSelector := range executorsMap {
		t.Logf("[ExecuteMCMSProposalV2] Setting root on chain %d...", chainSelector)
		root, err := executable.SetRoot(env.GetContext(), chainSelector)
		if err != nil {
			return fmt.Errorf("[ExecuteMCMSProposalV2] SetRoot failed: %w", err)
		}

		family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
		require.NoError(t, err)

		// no need to confirm transaction on solana as the MCMS sdk confirms it internally
		if family == chainsel.FamilyEVM {
			chain := evmChains[uint64(chainSelector)]
			evmTransaction := root.RawData.(*gethtypes.Transaction)
			t.Logf("[ExecuteMCMSProposalV2] SetRoot EVM tx hash: %s", evmTransaction.Hash().String())
			_, err = chain.Confirm(evmTransaction)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSProposalV2] Confirm failed: %w", err)
			}
		}
		if family == chainsel.FamilyAptos {
			chain := aptosChains[uint64(chainSelector)]
			tx := root.RawData.(*aptosapi.PendingTransaction)
			t.Logf("[ExecuteMCMSProposalV2] SetRoot Aptos tx hash: %s", tx.Hash)
			err = chain.Confirm(tx.Hash)
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
			chain := evmChains[uint64(op.ChainSelector)]
			evmTransaction := result.RawData.(*gethtypes.Transaction)
			t.Logf("[ExecuteMCMSProposalV2] Operation %d EVM tx hash: %s", i, evmTransaction.Hash().String())
			_, err = chain.Confirm(evmTransaction)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSProposalV2] Confirm failed: %w", err)
			}
		}
		if family == chainsel.FamilyAptos {
			chain := aptosChains[uint64(op.ChainSelector)]
			tx := result.RawData.(*aptosapi.PendingTransaction)
			t.Logf("[ExecuteMCMSProposalV2] Operation %d Aptos tx hash: %s", i, tx.Hash)
			err = chain.Confirm(tx.Hash)
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

	// build a "chainSelector => executor" map
	executorsMap := map[mcmstypes.ChainSelector]mcmssdk.TimelockExecutor{}
	callProxies := make([]string, len(timelockProposal.Operations))
	aptosChains := env.BlockChains.AptosChains()
	evmChains := env.BlockChains.EVMChains()
	solChains := env.BlockChains.SolanaChains()
	for i, op := range timelockProposal.Operations {
		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		switch family {
		case chainsel.FamilyEVM:
			executorsMap[op.ChainSelector] = mcmsevmsdk.NewTimelockExecutor(
				evmChains[uint64(op.ChainSelector)].Client,
				evmChains[uint64(op.ChainSelector)].DeployerKey)
			callProxies[i] = findCallProxyAddress(t, env, uint64(op.ChainSelector))

		case chainsel.FamilySolana:
			executorsMap[op.ChainSelector] = mcmssolanasdk.NewTimelockExecutor(
				solChains[uint64(op.ChainSelector)].Client,
				*solChains[uint64(op.ChainSelector)].DeployerKey)

		case chainsel.FamilyAptos:
			executorsMap[op.ChainSelector] = mcmsaptossdk.NewTimelockExecutor(
				aptosChains[uint64(op.ChainSelector)].Client,
				aptosChains[uint64(op.ChainSelector)].DeployerSigner)

		default:
			require.FailNow(t, "unsupported chain family")
		}
	}

	timelockExecutable, err := mcmslib.NewTimelockExecutable(env.GetContext(), timelockProposal, executorsMap)
	require.NoError(t, err)

	isReady := func() error {
		err := timelockExecutable.IsReady(env.GetContext())
		return err
	}
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.NoErrorf(collect, isReady(), "Proposal is not ready")
	}, 100*time.Second, 50*time.Millisecond)

	// execute each operation sequentially
	var tx = mcmstypes.TransactionResult{}
	for i, op := range timelockProposal.Operations {
		opOpts := slices.Clone(opts)
		if callProxies[i] != "" {
			opOpts = append(opOpts, mcmslib.WithCallProxy(callProxies[i]))
		}

		tx, err = timelockExecutable.Execute(env.GetContext(), i, opOpts...)
		if err != nil {
			return fmt.Errorf("[ExecuteMCMSTimelockProposalV2] Execute failed: %w", err)
		}
		t.Logf("[ExecuteMCMSTimelockProposalV2] Executed timelock operation index=%d on chain %d", i, uint64(op.ChainSelector))
		family, err := chainsel.GetSelectorFamily(uint64(op.ChainSelector))
		require.NoError(t, err)

		// no need to confirm transaction on solana as the MCMS sdk confirms it internally
		if family == chainsel.FamilyEVM {
			chain := evmChains[uint64(op.ChainSelector)]
			evmTransaction := tx.RawData.(*gethtypes.Transaction)
			_, err = chain.Confirm(evmTransaction)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSTimelockProposalV2] Confirm failed: %w", err)
			}
		}
		if family == chainsel.FamilyAptos {
			chain := aptosChains[uint64(op.ChainSelector)]
			aptosTx := tx.RawData.(*aptosapi.PendingTransaction)
			err = chain.Confirm(aptosTx.Hash)
			if err != nil {
				return fmt.Errorf("[ExecuteMCMSTimelockProposalV2] Confirm failed: %w", err)
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

func findCallProxyAddress(t *testing.T, env cldf.Environment, chainSelector uint64) string {
	addressesForChain, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)

	for address, tvStr := range addressesForChain {
		if tvStr.Type == commontypes.CallProxy && tvStr.Version == deployment.Version1_0_0 {
			return address
		}
	}

	require.FailNow(t, "unable to find call proxy address")
	return ""
}
