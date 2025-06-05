package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

var CCIPHomeABI *abi.ABI

func init() {
	var err error
	CCIPHomeABI, err = ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
}

type DONSequenceDeps struct {
	HomeChain cldf_evm.Chain
}

type DONAddition struct {
	ExpectedID       uint32
	PluginConfig     ccip_home.CCIPHomeOCR3Config
	PeerIDs          [][32]byte
	F                uint8
	IsPublic         bool
	AcceptsWorkflows bool
}

type AddDONAndSetCandidateSequenceInput struct {
	CapabilitiesRegistry common.Address
	NoSend               bool
	DONs                 []DONAddition
}

var AddDONAndSetCandidateSequence = operations.NewSequence(
	"AddDONAndSetCandidateSequence",
	semver.MustParse("1.0.0"),
	"Adds commit / exec DONs for chains and sets their candidates on CCIPHome",
	func(b operations.Bundle, deps DONSequenceDeps, input AddDONAndSetCandidateSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
		opOutputs := make(map[uint64][]opsutil.EVMCallOutput, 1) // Only calls against the home chain will be made
		opOutputs[deps.HomeChain.Selector] = make([]opsutil.EVMCallOutput, len(input.DONs))

		for i, don := range input.DONs {
			encodedSetCandidateCall, err := CCIPHomeABI.Pack(
				"setCandidate",
				don.ExpectedID,
				don.PluginConfig.PluginType,
				don.PluginConfig,
				[32]byte{},
			)
			if err != nil {
				return map[uint64][]opsutil.EVMCallOutput{}, fmt.Errorf("failed to pack set candidate call: %w", err)
			}
			report, err := operations.ExecuteOperation(
				b,
				ccipops.AddDONOp,
				deps.HomeChain,
				opsutil.EVMCallInput[ccipops.AddDONOpInput]{
					Address:       input.CapabilitiesRegistry,
					ChainSelector: deps.HomeChain.Selector,
					CallInput: ccipops.AddDONOpInput{
						Nodes: don.PeerIDs,
						CapabilityConfigurations: []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: shared.CCIPCapabilityID,
								Config:       encodedSetCandidateCall,
							},
						},
						IsPublic:         don.IsPublic,
						AcceptsWorkflows: don.AcceptsWorkflows,
						F:                don.F,
					},
					NoSend: input.NoSend,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to execute AddDONOp for chain with selector %d and plugin type %s: %w", don.PluginConfig.ChainSelector, types.PluginType(don.PluginConfig.PluginType), err)
			}
			opOutputs[deps.HomeChain.Selector][i] = report.Output
		}

		return opOutputs, nil
	})

type DONUpdate struct {
	ID             uint32
	PluginConfig   ccip_home.CCIPHomeOCR3Config
	PeerIDs        [][32]byte
	F              uint8
	IsPublic       bool
	ExistingDigest [32]byte
}

type SetCandidateSequenceInput struct {
	CapabilitiesRegistry common.Address
	NoSend               bool
	DONs                 []DONUpdate
}

var SetCandidateSequence = operations.NewSequence(
	"SetCandidateSequence",
	semver.MustParse("1.0.0"),
	"Updates candidates for existing commit / exec DONs across multiple chains",
	func(b operations.Bundle, deps DONSequenceDeps, input SetCandidateSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
		opOutputs := make(map[uint64][]opsutil.EVMCallOutput, 1) // Only calls against the home chain will be made
		opOutputs[deps.HomeChain.Selector] = make([]opsutil.EVMCallOutput, len(input.DONs))

		for i, don := range input.DONs {
			encodedSetCandidateCall, err := CCIPHomeABI.Pack(
				"setCandidate",
				don.ID,
				don.PluginConfig.PluginType,
				don.PluginConfig,
				don.ExistingDigest,
			)
			if err != nil {
				return map[uint64][]opsutil.EVMCallOutput{}, fmt.Errorf("failed to pack set candidate call: %w", err)
			}
			report, err := operations.ExecuteOperation(
				b,
				ccipops.UpdateDONOp,
				deps.HomeChain,
				opsutil.EVMCallInput[ccipops.UpdateDONOpInput]{
					Address:       input.CapabilitiesRegistry,
					ChainSelector: deps.HomeChain.Selector,
					CallInput: ccipops.UpdateDONOpInput{
						ID:    don.ID,
						Nodes: don.PeerIDs,
						CapabilityConfigurations: []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: shared.CCIPCapabilityID,
								Config:       encodedSetCandidateCall,
							},
						},
						IsPublic: don.IsPublic,
						F:        don.F,
					},
					NoSend: input.NoSend,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to execute UpdateDONOp for chain with selector %d and plugin type %s: %w", don.PluginConfig.ChainSelector, types.PluginType(don.PluginConfig.PluginType), err)
			}
			opOutputs[deps.HomeChain.Selector][i] = report.Output
		}

		return opOutputs, nil
	})

type DONUpdatePromotion struct {
	ID              uint32
	PluginType      uint8
	ChainSelector   uint64
	PeerIDs         [][32]byte
	F               uint8
	IsPublic        bool
	CandidateDigest [32]byte
	ActiveDigest    [32]byte
}

type PromoteCandidateSequenceInput struct {
	CapabilitiesRegistry common.Address
	NoSend               bool
	DONs                 []DONUpdatePromotion
}

var PromoteCandidateSequence = operations.NewSequence(
	"PromoteCandidateSequence",
	semver.MustParse("1.0.0"),
	"Promote candidates for existing commit / exec DONs across multiple chains",
	func(b operations.Bundle, deps DONSequenceDeps, input PromoteCandidateSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
		opOutputs := make(map[uint64][]opsutil.EVMCallOutput, 1) // Only calls against the home chain will be made
		opOutputs[deps.HomeChain.Selector] = make([]opsutil.EVMCallOutput, len(input.DONs))

		for i, don := range input.DONs {
			encodedPromoteCandidateCall, err := CCIPHomeABI.Pack(
				"promoteCandidateAndRevokeActive",
				don.ID,
				don.PluginType,
				don.CandidateDigest,
				don.ActiveDigest,
			)
			if err != nil {
				return map[uint64][]opsutil.EVMCallOutput{}, fmt.Errorf("failed to pack promote candidate call: %w", err)
			}
			report, err := operations.ExecuteOperation(
				b,
				ccipops.UpdateDONOp,
				deps.HomeChain,
				opsutil.EVMCallInput[ccipops.UpdateDONOpInput]{
					Address:       input.CapabilitiesRegistry,
					ChainSelector: deps.HomeChain.Selector,
					CallInput: ccipops.UpdateDONOpInput{
						ID:    don.ID,
						Nodes: don.PeerIDs,
						CapabilityConfigurations: []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
							{
								CapabilityId: shared.CCIPCapabilityID,
								Config:       encodedPromoteCandidateCall,
							},
						},
						IsPublic: don.IsPublic,
						F:        don.F,
					},
					NoSend: input.NoSend,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to execute UpdateDONOp for chain with selector %d and plugin type %s: %w", don.ChainSelector, types.PluginType(don.PluginType), err)
			}
			opOutputs[deps.HomeChain.Selector][i] = report.Output
		}

		return opOutputs, nil
	})
