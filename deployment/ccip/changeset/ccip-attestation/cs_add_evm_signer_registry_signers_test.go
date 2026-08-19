package ccip_attestation

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/operations/contract"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	cciputils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
)

func TestEVMSignerRegistryAddSigners_DirectExpansionRetryAndReAdd(t *testing.T) {
	initial := makeContractSigners(1, 5)
	e := newAddSignersEnvironment(t, []uint64{BaseSepoliaSelector}, initial)
	registry := getTestSignerRegistry(t, e, BaseSepoliaSelector)
	additions := makeConfigSigners(101, 5)
	cfg := AddSignersConfig{
		SignersByChain: map[uint64][]Signer{BaseSepoliaSelector: additions},
		MCMS:           validAddSignersMCMSInput(),
	}

	require.NoError(t, EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(*e, cfg))
	out, err := EVMSignerRegistryAddSignersChangeset.Apply(*e, cfg)
	require.NoError(t, err)
	require.Len(t, out.Reports, 1)
	require.Empty(t, out.MCMSTimelockProposals)
	require.Nil(t, out.DataStore)
	require.Equal(t, append(initial, toContractSigners(additions)...), readTestSigners(t, e, registry))

	staleRetry := cfg
	staleRetry.MCMS.ValidUntil = 1
	out, err = EVMSignerRegistryAddSignersChangeset.Apply(*e, staleRetry)
	require.NoError(t, err)
	require.Empty(t, out.Reports)
	require.Empty(t, out.MCMSTimelockProposals)

	chain := e.BlockChains.EVMChains()[BaseSepoliaSelector]
	removed := make([]common.Address, len(additions))
	for i, signer := range additions {
		removed[i] = signer.EVMAddress
	}
	tx, err := registry.RemoveSigners(chain.DeployerKey, removed)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	require.Equal(t, initial, readTestSigners(t, e, registry))
	out, err = EVMSignerRegistryAddSignersChangeset.Apply(*e, cfg)
	require.NoError(t, err)
	require.Len(t, out.Reports, 1)
	require.Len(t, readTestSigners(t, e, registry), 10)
}

func TestEVMSignerRegistryAddSigners_MultiChainDeterministicAndPreflighted(t *testing.T) {
	initial := makeContractSigners(1, 5)
	selectors := []uint64{BaseMainnetSelector, BaseSepoliaSelector}
	e := newAddSignersEnvironment(t, selectors, initial)
	cfg := AddSignersConfig{
		SignersByChain: map[uint64][]Signer{
			BaseMainnetSelector: {makeConfigSigner(301)},
			BaseSepoliaSelector: {makeConfigSigner(201)},
		},
	}

	out, err := EVMSignerRegistryAddSignersChangeset.Apply(*e, cfg)
	require.NoError(t, err)
	require.Len(t, out.Reports, 2)
	firstInput, ok := out.Reports[0].Input.(contract.FunctionInput[[]Signer])
	require.True(t, ok)
	secondInput, ok := out.Reports[1].Input.(contract.FunctionInput[[]Signer])
	require.True(t, ok)
	require.Equal(t, uint64(BaseSepoliaSelector), firstInput.ChainSelector)
	require.Equal(t, uint64(BaseMainnetSelector), secondInput.ChainSelector)

	beforeSepolia := readTestSigners(t, e, getTestSignerRegistry(t, e, BaseSepoliaSelector))
	beforeMainnet := readTestSigners(t, e, getTestSignerRegistry(t, e, BaseMainnetSelector))
	overCapacity := AddSignersConfig{
		SignersByChain: map[uint64][]Signer{
			BaseSepoliaSelector: {makeConfigSigner(401)},
			BaseMainnetSelector: makeConfigSigners(500, 15),
		},
	}
	_, err = EVMSignerRegistryAddSignersChangeset.Apply(*e, overCapacity)
	require.ErrorContains(t, err, "exceeds SignerRegistry capacity")
	require.Equal(t, beforeSepolia, readTestSigners(t, e, getTestSignerRegistry(t, e, BaseSepoliaSelector)))
	require.Equal(t, beforeMainnet, readTestSigners(t, e, getTestSignerRegistry(t, e, BaseMainnetSelector)))
}

func TestEVMSignerRegistryAddSigners_RotationAndKeyReconciliation(t *testing.T) {
	active := testAddress(1)
	e := newAddSignersEnvironment(t, []uint64{BaseSepoliaSelector}, []signer_registry.ISignerRegistrySigner{{EvmAddress: active}})
	registry := getTestSignerRegistry(t, e, BaseSepoliaSelector)
	chain := e.BlockChains.EVMChains()[BaseSepoliaSelector]
	pending := testAddress(2)
	tx, err := registry.SetNewSignerAddresses(chain.DeployerKey, []common.Address{active}, []common.Address{pending})
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	for _, signer := range []Signer{
		{EVMAddress: active},
		{EVMAddress: active, NewEVMAddress: pending},
	} {
		out, applyErr := EVMSignerRegistryAddSignersChangeset.Apply(*e, AddSignersConfig{
			SignersByChain: map[uint64][]Signer{BaseSepoliaSelector: {signer}},
		})
		require.NoError(t, applyErr)
		require.Empty(t, out.Reports)
	}

	_, err = EVMSignerRegistryAddSignersChangeset.Apply(*e, AddSignersConfig{
		SignersByChain: map[uint64][]Signer{
			BaseSepoliaSelector: {{EVMAddress: active, NewEVMAddress: testAddress(3)}},
		},
	})
	require.ErrorContains(t, err, "add-signers cannot change")

	_, err = EVMSignerRegistryAddSignersChangeset.Apply(*e, AddSignersConfig{
		SignersByChain: map[uint64][]Signer{
			BaseSepoliaSelector: {{EVMAddress: testAddress(4), NewEVMAddress: pending}},
		},
	})
	require.ErrorContains(t, err, "already in use")
}

func TestEVMSignerRegistryAddSigners_ValidationDatastoreAndYAML(t *testing.T) {
	valid := makeConfigSigner(10)
	tests := []struct {
		name    string
		signers []Signer
		wantErr string
	}{
		{name: "empty", wantErr: "no signers provided"},
		{name: "zero active", signers: []Signer{{}}, wantErr: "zero EVM address"},
		{name: "same active and pending", signers: []Signer{{EVMAddress: valid.EVMAddress, NewEVMAddress: valid.EVMAddress}}, wantErr: "identical"},
		{name: "duplicate across roles", signers: []Signer{valid, {EVMAddress: testAddress(11), NewEVMAddress: valid.EVMAddress}}, wantErr: "reuses"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorContains(t, validateSigners(tt.signers), tt.wantErr)
		})
	}

	e := newAddSignersEnvironment(t, []uint64{BaseSepoliaSelector}, makeContractSigners(1, 5))
	cfg := AddSignersConfig{
		SignersByChain: map[uint64][]Signer{BaseSepoliaSelector: {valid}},
	}
	withoutRegistry := datastore.NewMemoryDataStore()
	refs, err := e.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	for _, ref := range refs {
		if ref.Type == datastore.ContractType(shared.EVMSignerRegistry) && ref.Version.Equal(&deployment.Version1_0_0) {
			continue
		}
		require.NoError(t, withoutRegistry.Addresses().Add(ref))
	}
	missing := e.Clone()
	missing.DataStore = withoutRegistry.Seal()
	_, err = EVMSignerRegistryAddSignersChangeset.Apply(missing, cfg)
	require.ErrorContains(t, err, "expected to find exactly 1 ref")

	yamlInput := fmt.Sprintf(`signersByChain:
  %d:
    - evmAddress: "%s"
      newEVMAddress: "%s"
`, uint64(BaseSepoliaSelector), testAddress(700).Hex(), testAddress(701).Hex())
	var decoded AddSignersConfig
	require.NoError(t, yaml.Unmarshal([]byte(yamlInput), &decoded))
	require.Equal(t, testAddress(700), decoded.SignersByChain[BaseSepoliaSelector][0].EVMAddress)
	require.Equal(t, testAddress(701), decoded.SignersByChain[BaseSepoliaSelector][0].NewEVMAddress)
	require.Equal(t, mcms.Input{}, decoded.MCMS)
	require.NoError(t, EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(*e, decoded))
}

func TestEVMSignerRegistryAddSigners_ConfigValidation(t *testing.T) {
	e := newAddSignersEnvironment(t, []uint64{BaseSepoliaSelector}, nil)
	tests := []struct {
		name    string
		cfg     AddSignersConfig
		wantErr string
	}{
		{name: "empty config", cfg: AddSignersConfig{}, wantErr: "no signer additions provided"},
		{
			name: "empty signer list",
			cfg: AddSignersConfig{
				SignersByChain: map[uint64][]Signer{BaseSepoliaSelector: nil},
			},
			wantErr: "no signers provided",
		},
		{
			name: "unsupported chain",
			cfg: AddSignersConfig{
				SignersByChain: map[uint64][]Signer{1: {{EVMAddress: testAddress(10)}}},
			},
			wantErr: "not a supported Base chain",
		},
		{
			name: "missing Base chain",
			cfg: AddSignersConfig{
				SignersByChain: map[uint64][]Signer{BaseMainnetSelector: {{EVMAddress: testAddress(10)}}},
			},
			wantErr: "not found in environment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorContains(t, EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(*e, tt.cfg), tt.wantErr)
		})
	}

	withoutDataStore := e.Clone()
	withoutDataStore.DataStore = nil
	require.ErrorContains(t, EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(withoutDataStore, AddSignersConfig{
		SignersByChain: map[uint64][]Signer{BaseSepoliaSelector: {{EVMAddress: testAddress(10)}}},
	}), "environment DataStore is required")
}

func TestEVMSignerRegistryAddSigners_MCMSExecutionAndRetry(t *testing.T) {
	selector := uint64(BaseSepoliaSelector)
	initial := []signer_registry.ISignerRegistrySigner{{EvmAddress: testAddress(1)}}
	e := newAddSignersEnvironment(t, []uint64{selector}, initial)
	registry := getTestSignerRegistry(t, e, selector)
	rt := runtime.NewFromEnvironment(*e)

	qualifier := cciputils.CLLQualifier
	mcmsConfig := cldftesthelpers.SingleGroupTimelockConfig(t)
	mcmsConfig.Qualifier = &qualifier
	require.NoError(t, rt.Exec(runtime.ChangesetTask(
		cldf.CreateLegacyChangeSet(mcmschangesets.DeployMCMSWithTimelockV2),
		map[uint64]cldfproposalutils.MCMSWithTimelockConfig{selector: mcmsConfig},
	)))

	cfg := AddSignersConfig{
		SignersByChain: map[uint64][]Signer{
			selector: {{EVMAddress: testAddress(10), NewEVMAddress: testAddress(11)}},
		},
		MCMS: validAddSignersMCMSInput(),
	}
	require.NoError(t, EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(rt.Environment(), cfg))

	require.NoError(t, rt.Exec(
		runtime.ChangesetTask(
			cldf.CreateLegacyChangeSet(mcmschangesets.TransferToMCMSWithTimelockV2),
			mcmschangesets.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{selector: {registry.Address()}},
				MCMSConfig: cldfproposalutils.TimelockConfig{
					MinDelay: 0,
					TimelockQualifierPerChain: map[uint64]string{
						selector: cciputils.CLLQualifier,
					},
				},
			},
		),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	))
	invalidCfg := cfg
	invalidCfg.MCMS.ValidUntil = 1
	require.ErrorContains(t, EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(rt.Environment(), invalidCfg), "invalid MCMS proposal configuration")

	proposalTask := runtime.ChangesetTask(EVMSignerRegistryAddSignersChangeset, cfg)
	require.NoError(t, rt.Exec(proposalTask))
	require.Equal(t, initial, readTestSigners(t, e, registry))
	require.Len(t, rt.State().GetPendingProposals(), 1)

	proposalOutput := rt.State().Outputs[proposalTask.ID()]
	require.Len(t, proposalOutput.Reports, 1)
	require.Len(t, proposalOutput.MCMSTimelockProposals, 1)
	require.Len(t, proposalOutput.MCMSTimelockProposals[0].Operations, 1)
	require.Len(t, proposalOutput.MCMSTimelockProposals[0].Operations[0].Transactions, 1)

	require.NoError(t, rt.Exec(
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	))
	require.Equal(t, []signer_registry.ISignerRegistrySigner{
		{EvmAddress: testAddress(1)},
		{EvmAddress: testAddress(10), NewEVMAddress: testAddress(11)},
	}, readTestSigners(t, e, registry))

	proposalCount := len(rt.State().Proposals)
	retryCfg := cfg
	retryCfg.MCMS.ValidUntil = 1
	retryTask := runtime.ChangesetTask(EVMSignerRegistryAddSignersChangeset, retryCfg)
	require.NoError(t, rt.Exec(retryTask))
	require.Len(t, rt.State().Proposals, proposalCount)
	require.Empty(t, rt.State().Outputs[retryTask.ID()].Reports)
	require.Empty(t, rt.State().Outputs[retryTask.ID()].MCMSTimelockProposals)
}

func TestEVMSignerRegistryAddSigners_RejectsArbitraryOwner(t *testing.T) {
	selector := uint64(BaseSepoliaSelector)
	e := newAddSignersEnvironment(t, []uint64{selector}, []signer_registry.ISignerRegistrySigner{
		{EvmAddress: testAddress(10)},
	})
	registry := getTestSignerRegistry(t, e, selector)
	chain := e.BlockChains.EVMChains()[selector]
	require.NotEmpty(t, chain.Users)
	rt := runtime.NewFromEnvironment(*e)

	qualifier := cciputils.CLLQualifier
	mcmsConfig := cldftesthelpers.SingleGroupTimelockConfig(t)
	mcmsConfig.Qualifier = &qualifier
	require.NoError(t, rt.Exec(runtime.ChangesetTask(
		cldf.CreateLegacyChangeSet(mcmschangesets.DeployMCMSWithTimelockV2),
		map[uint64]cldfproposalutils.MCMSWithTimelockConfig{selector: mcmsConfig},
	)))

	tx, err := registry.TransferOwnership(chain.DeployerKey, chain.Users[0].From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)
	tx, err = registry.AcceptOwnership(chain.Users[0])
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	err = EVMSignerRegistryAddSignersChangeset.VerifyPreconditions(rt.Environment(), AddSignersConfig{
		SignersByChain: map[uint64][]Signer{
			selector: {{EVMAddress: testAddress(10)}},
		},
		MCMS: validAddSignersMCMSInput(),
	})
	require.ErrorContains(t, err, "unsupported owner")
}

func newAddSignersEnvironment(
	t *testing.T,
	selectors []uint64,
	initial []signer_registry.ISignerRegistrySigner,
) *cldf.Environment {
	t.Helper()
	e, err := environment.New(t.Context(),
		environment.WithEVMSimulatedWithConfig(t, selectors, onchain.EVMSimLoaderConfig{
			NumAdditionalAccounts: 1,
		}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	out, err := EVMSignerRegistryDeploymentChangeset.Apply(*e, SignerRegistryChangesetConfig{
		MaxSigners: MaxSigners,
		Signers:    initial,
	})
	require.NoError(t, err)
	require.NotNil(t, out.DataStore)
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Merge(e.DataStore))
	require.NoError(t, ds.Merge(out.DataStore.Seal()))
	e.DataStore = ds.Seal()
	return e
}

func getTestSignerRegistry(t *testing.T, e *cldf.Environment, selector uint64) *signer_registry.SignerRegistry {
	t.Helper()
	address, err := resolveSignerRegistry(*e, selector)
	require.NoError(t, err)
	registry, err := signer_registry.NewSignerRegistry(address, e.BlockChains.EVMChains()[selector].Client)
	require.NoError(t, err)
	return registry
}

func readTestSigners(t *testing.T, e *cldf.Environment, registry *signer_registry.SignerRegistry) []signer_registry.ISignerRegistrySigner {
	t.Helper()
	signers, err := registry.GetSigners(nil)
	require.NoError(t, err)
	return signers
}

func makeContractSigners(start, count int64) []signer_registry.ISignerRegistrySigner {
	signers := make([]signer_registry.ISignerRegistrySigner, count)
	for i := range signers {
		signers[i] = signer_registry.ISignerRegistrySigner{EvmAddress: testAddress(start + int64(i))}
	}
	return signers
}

func makeConfigSigners(start, count int64) []Signer {
	signers := make([]Signer, count)
	for i := range signers {
		signers[i] = makeConfigSigner(start + int64(i))
	}
	return signers
}

func makeConfigSigner(value int64) Signer {
	return Signer{EVMAddress: testAddress(value)}
}

func testAddress(value int64) common.Address {
	return common.BigToAddress(big.NewInt(value))
}

func validAddSignersMCMSInput() mcms.Input {
	return mcms.Input{
		ValidUntil:     4_000_000_000,
		TimelockDelay:  mcms_types.NewDuration(0),
		TimelockAction: mcms_types.TimelockActionSchedule,
		Qualifier:      cciputils.CLLQualifier,
		Description:    "Add CCIP attestation signer entries",
	}
}
