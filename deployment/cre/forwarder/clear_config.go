package forwarder

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	cap_reg_v2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type ClearConfigSeqDeps struct {
	Env *cldf.Environment
}

type ClearConfigSeqInput struct {
	DonID         uint32 // the DON id whose config should be cleared from the forwarder
	ConfigVersion uint32 // the config version to clear. Must match the version the forwarder was configured with
	Qualifier     string // used to differentiate Forwarder contracts deployed to the same chain.

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *contracts.MCMSConfig
	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains map[uint64]struct{}
}

func (i ClearConfigSeqInput) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ClearConfigSeqOutput struct {
	MCMSTimelockProposals []mcmslib.TimelockProposal
}

// ClearConfigSeq is a sequence that clears the config of a given DON from Keystone Forwarder contracts.
var ClearConfigSeq = operations.NewSequence[ClearConfigSeqInput, ClearConfigSeqOutput, ClearConfigSeqDeps](
	"clear-forwarders-config-seq",
	semver.MustParse("1.0.0"),
	"Clear Keystone Forwarders Config",
	func(b operations.Bundle, deps ClearConfigSeqDeps, input ClearConfigSeqInput) (ClearConfigSeqOutput, error) {
		evmChain := deps.Env.BlockChains.EVMChains()

		var out ClearConfigSeqOutput
		for _, chain := range evmChain {
			if _, shouldInclude := input.Chains[chain.Selector]; len(input.Chains) > 0 && !shouldInclude {
				continue
			}

			filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
				datastore.AddressRefByChainSelector(chain.Selector),
				datastore.AddressRefByType(datastore.ContractType(contracts.KeystoneForwarder)),
				datastore.AddressRefByQualifier(input.Qualifier),
			}

			addressesRefs := deps.Env.DataStore.Addresses().Filter(filters...)
			if len(addressesRefs) == 0 {
				return ClearConfigSeqOutput{}, fmt.Errorf("clear-forwarders-config-seq failed: no KeystoneForwarder contract found for chain selector %d and qualifier %q", chain.Selector, input.Qualifier)
			}

			if len(addressesRefs) > 1 {
				return ClearConfigSeqOutput{}, fmt.Errorf("clear-forwarders-config-seq failed: expected one KeystoneForwarder contract for chain selector %d and qualifier %q, found %d", chain.Selector, input.Qualifier, len(addressesRefs))
			}

			var mcmsContracts *evmstate.MCMSWithTimelockState
			if input.MCMSConfig != nil {
				var err error
				mcmsContracts, err = strategies.GetMCMSContracts(*deps.Env, chain.Selector, *input.MCMSConfig)
				if err != nil {
					return ClearConfigSeqOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
				}
			}

			for _, addrRef := range addressesRefs {
				// Create the appropriate strategy
				strategy, err := strategies.CreateStrategy(
					chain,
					*deps.Env,
					input.MCMSConfig,
					mcmsContracts,
					common.HexToAddress(addrRef.Address),
					cap_reg_v2.ClearForwarderConfigDescription,
				)
				if err != nil {
					return ClearConfigSeqOutput{}, fmt.Errorf("failed to create strategy: %w", err)
				}

				contract, err := contracts.GetOwnedContractV2[*forwarder.KeystoneForwarder](deps.Env.DataStore.Addresses(), chain, addrRef.Address, addrRef.Qualifier)
				if err != nil {
					return ClearConfigSeqOutput{}, fmt.Errorf("clear-forwarders-config-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", chain.Selector, err)
				}

				fwrReport, err := operations.ExecuteOperation(b, ClearConfigOp, ClearConfigOpDeps{
					Chain:    &chain,
					Contract: contract.Contract,
					Strategy: strategy,
				}, ClearConfigOpInput{
					DonID:         input.DonID,
					ConfigVersion: input.ConfigVersion,
					UseMCMS:       input.UseMCMS(),
					ChainSelector: chain.Selector, // here to skip the check for the previous report, since unless inputs are different they are treated as the same and skipped
				})
				if err != nil {
					return ClearConfigSeqOutput{}, fmt.Errorf("clear-forwarders-config-seq failed for chain selector %d: %w", chain.Selector, err)
				}

				if input.UseMCMS() {
					if contract.McmsContracts == nil {
						return ClearConfigSeqOutput{}, fmt.Errorf("clear-forwarders-config-seq failed: expected forwarder contract %s to be owned by MCMS for chain selector %d", contract.Contract.Address(), chain.Selector)
					}
					proposal, err := strategy.BuildProposal([]mcmstypes.BatchOperation{*fwrReport.Output.MCMSOperation})
					if err != nil {
						return ClearConfigSeqOutput{}, fmt.Errorf("clear-forwarders-config-seq failed to build proposal for chain selector %d: %w", chain.Selector, err)
					}

					out.MCMSTimelockProposals = append(out.MCMSTimelockProposals, *proposal)
				}
			}
		}

		if input.UseMCMS() && len(out.MCMSTimelockProposals) == 0 {
			return ClearConfigSeqOutput{}, errors.New("clear-forwarders-config-seq failed: no proposals generated for MCMS")
		}

		return out, nil
	},
)

type ClearConfigOpDeps struct {
	Chain    *evm.Chain
	Contract *forwarder.KeystoneForwarder
	Strategy strategies.TransactionStrategy
}

type ClearConfigOpInput struct {
	UseMCMS       bool
	ChainSelector uint64
	DonID         uint32
	ConfigVersion uint32
}

type ClearConfigOpOutput struct {
	MCMSOperation *mcmstypes.BatchOperation // if using MCMS, the batch operation to propose the change

	Forwarder     common.Address
	DonID         uint32
	ConfigVersion uint32
}

// ClearConfigOp is an operation that clears a DON config from a Keystone Forwarder contract.
var ClearConfigOp = operations.NewOperation[ClearConfigOpInput, ClearConfigOpOutput, ClearConfigOpDeps](
	"clear-forwarder-config-op",
	semver.MustParse("1.0.0"),
	"Clear Keystone Forwarder Config",
	func(b operations.Bundle, deps ClearConfigOpDeps, input ClearConfigOpInput) (ClearConfigOpOutput, error) {
		r, err := clearForwarderConfig(b.Logger, *deps.Chain, deps.Contract, input.DonID, input.ConfigVersion, input.UseMCMS, deps.Strategy)
		if err != nil {
			return ClearConfigOpOutput{}, fmt.Errorf("clear-forwarder-config-op failed: failed to clear forwarder config for chain selector %d: %w", deps.Chain.Selector, err)
		}
		return ClearConfigOpOutput{
			MCMSOperation: r.MCMSOperation,
			DonID:         input.DonID,
			ConfigVersion: input.ConfigVersion,
			Forwarder:     deps.Contract.Address(),
		}, nil
	},
)

// ClearForwardersConfig is a changeset that clears a DON config from Keystone Forwarder contracts.
type ClearForwardersConfig struct{}

var _ cldf.ChangeSetV2[ClearConfigSeqInput] = ClearForwardersConfig{}

func (c ClearForwardersConfig) VerifyPreconditions(e cldf.Environment, config ClearConfigSeqInput) error {
	for chainSel := range config.Chains {
		if _, ok := e.BlockChains.EVMChains()[chainSel]; !ok {
			return fmt.Errorf("chain selector %d not found in environment", chainSel)
		}
	}

	if config.DonID == 0 {
		return errors.New("DON ID must be non-zero")
	}
	if config.ConfigVersion == 0 {
		return errors.New("config version must be non-zero")
	}

	return nil
}

func (c ClearForwardersConfig) Apply(e cldf.Environment, config ClearConfigSeqInput) (cldf.ChangesetOutput, error) {
	deps := ClearConfigSeqDeps{
		Env: &e,
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ClearConfigSeq,
		deps,
		config,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               report.ExecutionReports,
		MCMSTimelockProposals: report.Output.MCMSTimelockProposals,
	}, nil
}
