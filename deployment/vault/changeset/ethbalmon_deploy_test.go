package changeset

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/eth_balance_monitor_wrapper"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestDeployEthBalMonValidation(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	selectorOther := chainselectors.TEST_90000002.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		config      types.DeployEthBalMonInput
		wantError   bool
		errorMsg    string
		setupMCMSIn bool
	}{
		{
			name: "empty chains",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{},
			},
			wantError: true,
			errorMsg:  "chains must not be empty",
		},
		{
			name: "unknown chain selector",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					math.MaxUint64: {
						SetKeeperRegistryAddress: "0x1234567890123456789012345678901234567890",
					},
				},
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("unknown chain selector %d", uint64(math.MaxUint64)),
		},
		{
			name: "chain not in environment",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selectorOther: {
						SetKeeperRegistryAddress: "0x1234567890123456789012345678901234567890",
					},
				},
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "empty setKeeperRegistryAddress",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "",
					},
				},
			},
			wantError: true,
			errorMsg:  "setKeeperRegistryAddress must not be empty",
		},
		{
			name: "invalid setKeeperRegistryAddress",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "not-a-valid-address",
					},
				},
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("chain %d: setKeeperRegistryAddress is not a valid hex address: not-a-valid-address", selector),
		},
		{
			name: "zero setKeeperRegistryAddress",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: zeroAddr,
					},
				},
			},
			wantError: true,
			errorMsg:  fmt.Sprintf("chain %d: setKeeperRegistryAddress cannot be zero address", selector),
		},
		{
			name: "missing MCMS and timelock in datastore",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "0x1234567890123456789012345678901234567890",
					},
				},
			},
			wantError: true,
			errorMsg:  "failed to get addresses from datastore",
		},
		{
			name: "valid config",
			config: types.DeployEthBalMonInput{
				Chains: map[uint64]types.DeployEthBalMonChainConfig{
					selector: {
						SetKeeperRegistryAddress: "0x1234567890123456789012345678901234567890",
					},
				},
			},
			wantError:   false,
			setupMCMSIn: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var testEnv cldf.Environment
			if test.setupMCMSIn {
				rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
					environment.WithEVMSimulated(t, []uint64{selector}),
				))
				require.NoError(t, err)
				setupMCMSInfrastructure(t, rt, []uint64{selector})
				testEnv = rt.Environment()
			} else {
				testEnv = *env
			}

			err := ValidateDeployEthBalMonConfig(testEnv.GetContext(), testEnv, test.config)

			if test.wantError {
				require.Error(t, err)
				if test.errorMsg != "" {
					require.Contains(t, err.Error(), test.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBuildAcceptOwnershipTimelockProposal(t *testing.T) {
	t.Parallel()

	t.Run("rejects empty contract set", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{selector}),
		)
		require.NoError(t, err)

		_, err = BuildAcceptOwnershipTimelockProposal(*env, AcceptOwnershipProposalInput{
			ContractsByChain: map[uint64]string{},
			Description:      "test",
			MCMSConfig:       proposalutils.TimelockConfig{MinDelay: 0, MCMSAction: mcmstypes.TimelockActionBypass},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "no contracts provided")
	})

	t.Run("chain not in environment", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{selector}),
		)
		require.NoError(t, err)

		otherSel := chainselectors.TEST_90000002.Selector
		_, err = BuildAcceptOwnershipTimelockProposal(*env, AcceptOwnershipProposalInput{
			ContractsByChain: map[uint64]string{
				otherSel: testAddr1,
			},
			Description: "test",
			MCMSConfig:  proposalutils.TimelockConfig{MinDelay: 0, MCMSAction: mcmstypes.TimelockActionBypass},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in environment")
	})

	t.Run("fails when MCMS or timelock missing from datastore", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{selector}),
		)
		require.NoError(t, err)

		_, err = BuildAcceptOwnershipTimelockProposal(*env, AcceptOwnershipProposalInput{
			ContractsByChain: map[uint64]string{
				selector: testAddr1,
			},
			Description: "test",
			MCMSConfig:  proposalutils.TimelockConfig{MinDelay: 0, MCMSAction: mcmstypes.TimelockActionBypass},
		})
		require.Error(t, err)
	})

	t.Run("builds proposal after deploy with custom description", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {SetKeeperRegistryAddress: testAddr1},
			},
		}
		out, err := DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.NoError(t, err)

		contractsByChain := make(map[uint64]string, len(cfg.Chains))
		for sel := range cfg.Chains {
			addr, err := GetContractAddress(out.DataStore, sel, cldf.ContractType(types.EthBalMonContractType))
			require.NoError(t, err)
			contractsByChain[sel] = addr
		}

		customDesc := "custom EthBalanceMonitor accept ownership proposal"
		prop, err := BuildAcceptOwnershipTimelockProposal(rt.Environment(), AcceptOwnershipProposalInput{
			ContractsByChain: contractsByChain,
			Description:      customDesc,
			MCMSConfig:       proposalutils.TimelockConfig{MinDelay: 0, MCMSAction: mcmstypes.TimelockActionBypass},
		})
		require.NoError(t, err)
		require.NotNil(t, prop)
		require.Equal(t, customDesc, prop.Description)
		require.Len(t, prop.Operations, len(cfg.Chains))
	})

	t.Run("uses default description when empty", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {SetKeeperRegistryAddress: testAddr1},
			},
		}
		out, err := DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.NoError(t, err)

		addr, err := GetContractAddress(out.DataStore, selector, cldf.ContractType(types.EthBalMonContractType))
		require.NoError(t, err)

		prop, err := BuildAcceptOwnershipTimelockProposal(rt.Environment(), AcceptOwnershipProposalInput{
			ContractsByChain: map[uint64]string{selector: addr},
			Description:      "",
			MCMSConfig:       proposalutils.TimelockConfig{MinDelay: 0, MCMSAction: mcmstypes.TimelockActionBypass},
		})
		require.NoError(t, err)
		require.Equal(t, "Accept ownership of EthBalanceMonitor across chains", prop.Description)
	})
}

func TestDeployEthBalMonChangeset(t *testing.T) {
	t.Parallel()

	t.Run("single chain", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		customWait := uint64(120)
		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {
					SetKeeperRegistryAddress: testAddr1,
					SetMinWaitPeriodSeconds:  &customWait,
				},
			},
		}

		require.NoError(t, DeployEthBalMonChangeSet.VerifyPreconditions(rt.Environment(), cfg))

		out, err := DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.NoError(t, err)
		assertEthBalMonDeployOutput(t, rt.Environment(), out, cfg)
	})

	t.Run("default min wait when unset uses 60s", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {
					SetKeeperRegistryAddress: testAddr1,
				},
			},
		}

		out, err := DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.NoError(t, err)

		mds, err := out.DataStore.ContractMetadata().Fetch()
		require.NoError(t, err)
		require.Len(t, mds, 1)
		mdMap := contractMetadataMap(t, mds[0].Metadata)
		require.Equal(t, uint64(60), uint64FromAny(t, mdMap["minWaitPeriodSeconds"]))
	})

	t.Run("explicit zero min wait uses default 60s", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		zero := uint64(0)
		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {
					SetKeeperRegistryAddress: testAddr1,
					SetMinWaitPeriodSeconds:  &zero,
				},
			},
		}

		out, err := DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.NoError(t, err)

		mds, err := out.DataStore.ContractMetadata().Fetch()
		require.NoError(t, err)
		require.Len(t, mds, 1)
		mdMap := contractMetadataMap(t, mds[0].Metadata)
		require.Equal(t, uint64(60), uint64FromAny(t, mdMap["minWaitPeriodSeconds"]))
	})

	t.Run("multiple chains", func(t *testing.T) {
		t.Parallel()

		selector1 := chainselectors.TEST_90000001.Selector
		selector2 := chainselectors.TEST_90000002.Selector
		selectors := []uint64{selector1, selector2}

		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, selectors),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, selectors)
		fundDeployerAccounts(t, rt.Environment(), selectors)

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector1: {SetKeeperRegistryAddress: testAddr1},
				selector2: {SetKeeperRegistryAddress: testAddr2},
			},
		}

		out, err := DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.NoError(t, err)
		assertEthBalMonDeployOutput(t, rt.Environment(), out, cfg)
	})

	t.Run("verify preconditions rejects invalid config", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		err = DeployEthBalMonChangeSet.VerifyPreconditions(rt.Environment(), types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chains must not be empty")
	})

	t.Run("verify preconditions rejects chain not in environment", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		otherSel := chainselectors.TEST_90000002.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		err = DeployEthBalMonChangeSet.VerifyPreconditions(rt.Environment(), types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				otherSel: {SetKeeperRegistryAddress: testAddr1},
			},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in environment")
	})

	t.Run("apply without MCMS infrastructure fails", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {SetKeeperRegistryAddress: testAddr1},
			},
		}

		_, err = DeployEthBalMonChangeSet.Apply(rt.Environment(), cfg)
		require.Error(t, err)
		require.ErrorContains(t, err, "timelock")
	})
}

func TestDeployEthBalMon_RuntimeChangesetTask(t *testing.T) {
	t.Parallel()

	t.Run("exec succeeds and merges EthBalMon into runtime datastore", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector: {SetKeeperRegistryAddress: testAddr1},
			},
		}

		task := runtime.ChangesetTask(DeployEthBalMonChangeSet, cfg)
		err = rt.Exec(task)
		require.NoError(t, err)

		_, hasOutput := rt.State().Outputs[task.ID()]
		require.True(t, hasOutput, "CLDF runtime should store ChangesetOutput under the task id")

		out := rt.State().Outputs[task.ID()]
		require.NotNil(t, out.DataStore)
		require.NotEmpty(t, out.MCMSTimelockProposals)

		records := rt.State().DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(selector),
			datastore.AddressRefByType(datastore.ContractType(types.EthBalMonContractType)),
		)
		require.Len(t, records, 1)
		labelSet := records[0].Labels.List()
		require.Contains(t, labelSet, types.EthBalMonContractType)
		require.Contains(t, labelSet, "EthBalMonV1_0_0")

		assertEthBalMonDeployOutput(t, rt.Environment(), out, cfg)
	})

	t.Run("exec fails on invalid precondition", func(t *testing.T) {
		t.Parallel()

		selector := chainselectors.TEST_90000001.Selector
		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, []uint64{selector}),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, []uint64{selector})
		fundDeployerAccounts(t, rt.Environment(), []uint64{selector})

		task := runtime.ChangesetTask(DeployEthBalMonChangeSet, types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{},
		})
		err = rt.Exec(task)
		require.Error(t, err)
		require.Contains(t, err.Error(), "chains must not be empty")
	})

	t.Run("multiple chains via runtime task", func(t *testing.T) {
		t.Parallel()

		selector1 := chainselectors.TEST_90000001.Selector
		selector2 := chainselectors.TEST_90000002.Selector
		selectors := []uint64{selector1, selector2}

		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, selectors),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, selectors)
		fundDeployerAccounts(t, rt.Environment(), selectors)

		cfg := types.DeployEthBalMonInput{
			Chains: map[uint64]types.DeployEthBalMonChainConfig{
				selector1: {SetKeeperRegistryAddress: testAddr1},
				selector2: {SetKeeperRegistryAddress: testAddr2},
			},
		}

		task := runtime.ChangesetTask(DeployEthBalMonChangeSet, cfg)
		require.NoError(t, rt.Exec(task))

		out := rt.State().Outputs[task.ID()]
		assertEthBalMonDeployOutput(t, rt.Environment(), out, cfg)
	})
}

// assertEthBalMonDeployOutput checks datastore, on-chain owner (post transferOwnership, pre-accept),
// and the accept-ownership timelock proposal.
func assertEthBalMonDeployOutput(
	t *testing.T,
	env cldf.Environment,
	out cldf.ChangesetOutput,
	cfg types.DeployEthBalMonInput,
) {
	t.Helper()

	n := len(cfg.Chains)
	require.NotNil(t, out.DataStore)

	addrs, err := out.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, addrs, n)

	mds, err := out.DataStore.ContractMetadata().Fetch()
	require.NoError(t, err)
	require.Len(t, mds, n)

	bySel := make(map[uint64]datastore.ContractMetadata)
	for _, m := range mds {
		bySel[m.ChainSelector] = m
	}

	for sel, chainCfg := range cfg.Chains {
		timelockAddr, err := GetContractAddress(env.DataStore, sel, commontypes.RBACTimelock)
		require.NoError(t, err)
		mcmsType := ethBalMonMCMSContractTypeForAction(deployEthBalMonAcceptOwnershipMCMSAction(cfg.MCMSConfig))
		mcmsAddr, err := GetContractAddress(env.DataStore, sel, mcmsType)
		require.NoError(t, err)

		meta, ok := bySel[sel]
		require.True(t, ok, "missing contract metadata for chain %d", sel)
		require.NotEmpty(t, meta.Address)

		md := contractMetadataMap(t, meta.Metadata)
		require.Equal(t, chainCfg.SetKeeperRegistryAddress, md["keeperRegistryAddress"])

		minWait := chainCfgMinWaitForEffective(chainCfg)
		require.Equal(t, effectiveMinWaitPeriodSeconds(minWait), uint64FromAny(t, md["minWaitPeriodSeconds"]))

		require.NotEmpty(t, md["deployTxHash"])
		require.NotZero(t, uint64FromAny(t, md["deployBlockNumber"]))
		require.Equal(t, timelockAddr, md["timelockAddress"])
		require.Equal(t, mcmsAddr, md["mcmsAddress"])
		require.NotEmpty(t, md["transferOwnershipTxHash"])

		ebmAddr, err := GetContractAddress(out.DataStore, sel, cldf.ContractType(types.EthBalMonContractType))
		require.NoError(t, err)

		chain := env.BlockChains.EVMChains()[sel]
		c, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(common.HexToAddress(ebmAddr), chain.Client)
		require.NoError(t, err)
		owner, err := c.Owner(nil)
		require.NoError(t, err)
		require.Equal(t, chain.DeployerKey.From, owner,
			"ConfirmedOwner keeps owner until timelock calls acceptOwnership")
	}

	require.Len(t, out.MCMSTimelockProposals, 1)
	prop := out.MCMSTimelockProposals[0]
	require.Contains(t, prop.Description, "EthBalanceMonitor")
	require.Len(t, prop.Operations, n)

	seen := make(map[uint64]bool)
	for _, op := range prop.Operations {
		sel := uint64(op.ChainSelector)
		seen[sel] = true
		require.Len(t, op.Transactions, 1)
		tx := op.Transactions[0]
		require.Equal(t, types.EthBalMonContractType, tx.ContractType)
		require.Contains(t, tx.Tags, "acceptOwnership")

		wantContract, err := GetContractAddress(out.DataStore, sel, cldf.ContractType(types.EthBalMonContractType))
		require.NoError(t, err)
		require.Equal(t, common.HexToAddress(wantContract), common.HexToAddress(tx.To))
	}

	for sel := range cfg.Chains {
		require.True(t, seen[sel], "proposal missing operation for chain %d", sel)
	}
}

func chainCfgMinWaitForEffective(c types.DeployEthBalMonChainConfig) uint64 {
	if c.SetMinWaitPeriodSeconds == nil {
		return 0
	}
	return *c.SetMinWaitPeriodSeconds
}

func contractMetadataMap(t *testing.T, raw any) map[string]any {
	t.Helper()
	m, ok := raw.(map[string]any)
	require.True(t, ok, "expected metadata map[string]any, got %T", raw)
	return m
}

func uint64FromAny(t *testing.T, v any) uint64 {
	t.Helper()
	require.NotNil(t, v)
	switch x := v.(type) {
	case uint64:
		return x
	case uint:
		return uint64(x)
	case uint32:
		return uint64(x)
	case int:
		require.GreaterOrEqual(t, x, 0)
		u, err := strconv.ParseUint(strconv.Itoa(x), 10, 64)
		require.NoError(t, err)
		return u
	case int64:
		require.GreaterOrEqual(t, x, int64(0))
		u, err := strconv.ParseUint(strconv.FormatInt(x, 10), 10, 64)
		require.NoError(t, err)
		return u
	case float64:
		return uint64(x)
	case json.Number:
		i, err := x.Int64()
		require.NoError(t, err)
		require.GreaterOrEqual(t, i, int64(0))
		u, err := strconv.ParseUint(strconv.FormatInt(i, 10), 10, 64)
		require.NoError(t, err)
		return u
	case string:
		u, err := strconv.ParseUint(x, 10, 64)
		require.NoError(t, err)
		return u
	default:
		require.Failf(t, "unexpected type for uint64 metadata field", "%T %#v", v, v)
		return 0
	}
}
