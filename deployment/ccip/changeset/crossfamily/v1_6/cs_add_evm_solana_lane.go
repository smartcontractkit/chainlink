package v1_6

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	solstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"

	solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var (
	AddEVMAndSolanaLaneChangeset = cldf.CreateChangeSet(addEVMAndSolanaLaneLogic, addEVMSolanaPreconditions)

	postOps = operations.NewOperation(
		"postOpsToAggregateProposals",
		semver.MustParse("1.0.0"),
		"Post ops to aggregate proposals",
		func(b operations.Bundle, deps Dependencies, input postOpsInput) ([]mcmslib.TimelockProposal, error) {
			allProposals := input.Proposals
			if input.MCMSConfig == nil || len(allProposals) == 0 {
				return input.Proposals, nil
			}
			proposal, err := proposeutils.AggregateProposals( //nolint:staticcheck //SA1019 ignoring deprecated
				deps.Env,
				deps.EVMMCMSState,
				deps.SolanaMCMSState,
				allProposals,
				"Adding EVM and Solana lane",
				input.MCMSConfig,
			)
			if err != nil {
				return nil, err
			}
			if proposal != nil {
				input.Proposals = []mcmslib.TimelockProposal{*proposal}
			}
			return input.Proposals, nil
		},
	)

	addEVMAndSolanaLaneSequence = operations.NewSequence(
		"addEVMAndSolanaLane",
		semver.MustParse("1.0.0"),
		"Adds bi-directional lane between EVM chain and Solana",
		func(b operations.Bundle, deps Dependencies, input AddMultiEVMSolanaLaneConfig) (OpsOutput, error) {
			allEVMChains := make([]uint64, 0, len(input.Configs))
			for _, cfg := range input.Configs {
				allEVMChains = append(allEVMChains, cfg.EVMChainSelector)
			}
			deps.Env.Logger.Infow("Adding EVM and Solana lane",
				"EVMChainSelector", allEVMChains, "SolanaChainSelector", input.SolanaChainSelector)
			finalOutput := &OpsOutput{}
			updateEVMOnRampReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMOnRamp",
				semver.MustParse("1.0.0"),
				"Updates EVM OnRamps with Destination Chain Configs for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateOnRampDestsConfig) ([]mcmslib.TimelockProposal, error) {
					output, err := v1_6.UpdateOnRampsDestsChangeset(deps.Env, input)
					if err != nil {
						return nil, err
					}
					return output.MCMSTimelockProposals, nil
				},
			), deps, deps.changesetInput.evmOnRampInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the changeset outputs
			if len(updateEVMOnRampReport.Output) > 0 {
				finalOutput.Proposals = append(finalOutput.Proposals, updateEVMOnRampReport.Output...)
			}
			// update EVM fee quoter dest chain
			updateEVMFeeQuoterDestChainReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMFeeQuoterDestChain",
				semver.MustParse("1.0.0"),
				"Updates EVM Fee Quoter with Destination Chain Configs for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateFeeQuoterDestsConfig) ([]mcmslib.TimelockProposal, error) {
					output, err := v1_6.UpdateFeeQuoterDestsChangeset(deps.Env, input)
					if err != nil {
						return nil, err
					}
					return output.MCMSTimelockProposals, nil
				},
			), deps, deps.changesetInput.evmFeeQuoterDestChainInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the changeset outputs
			if len(updateEVMFeeQuoterDestChainReport.Output) > 0 {
				finalOutput.Proposals = append(finalOutput.Proposals, updateEVMFeeQuoterDestChainReport.Output...)
			}
			// update EVM fee quoter prices
			updateEVMFeeQuoterPricesReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMFeeQuoterPrices",
				semver.MustParse("1.0.0"),
				"Updates EVM Fee Quoter with Prices for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateFeeQuoterPricesConfig) ([]mcmslib.TimelockProposal, error) {
					output, err := v1_6.UpdateFeeQuoterPricesChangeset(deps.Env, input)
					if err != nil {
						return nil, err
					}
					return output.MCMSTimelockProposals, nil
				},
			), deps, deps.changesetInput.evmFeeQuoterPriceInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the changeset outputs
			if len(updateEVMFeeQuoterPricesReport.Output) > 0 {
				finalOutput.Proposals = append(finalOutput.Proposals, updateEVMFeeQuoterPricesReport.Output...)
			}
			// update EVM off ramp
			updateEVMOffRampReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMOffRamp",
				semver.MustParse("1.0.0"),
				"Updates EVM OffRamps with Source Chain Configs for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateOffRampSourcesConfig) ([]mcmslib.TimelockProposal, error) {
					output, err := v1_6.UpdateOffRampSourcesChangeset(deps.Env, input)
					if err != nil {
						return nil, err
					}
					return output.MCMSTimelockProposals, nil
				},
			), deps, deps.changesetInput.evmOffRampInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the changeset outputs
			if len(updateEVMOffRampReport.Output) > 0 {
				finalOutput.Proposals = append(finalOutput.Proposals, updateEVMOffRampReport.Output...)
			}

			// update EVM router
			updateEVMRouterReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMRouter",
				semver.MustParse("1.0.0"),
				"Updates EVM Router with onRamp and OffRamp for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateRouterRampsConfig) ([]mcmslib.TimelockProposal, error) {
					output, err := v1_6.UpdateRouterRampsChangeset(deps.Env, input)
					if err != nil {
						return nil, err
					}
					return output.MCMSTimelockProposals, nil
				},
			), deps, deps.changesetInput.evmRouterInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the changeset outputs
			if len(updateEVMRouterReport.Output) > 0 {
				finalOutput.Proposals = append(finalOutput.Proposals, updateEVMRouterReport.Output...)
			}
			// update Solana router
			updateSolanaRouterReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateSolanaRouter",
				semver.MustParse("1.0.0"),
				"Updates Solana Router with EVM Chain Configs",
				func(b operations.Bundle, deps Dependencies, input solana.AddRemoteChainToRouterConfig) (OpsOutput, error) {
					output, err := solana.AddRemoteChainToRouter(deps.Env, input)
					if err != nil {
						return OpsOutput{}, err
					}
					return OpsOutput{
						Proposals:   output.MCMSTimelockProposals,
						AddressBook: output.AddressBook, //nolint:staticcheck //SA1019 ignoring deprecated
						DataStore:   output.DataStore,
					}, nil
				},
			), deps, deps.changesetInput.solanaRouterInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the output
			err = finalOutput.Merge(updateSolanaRouterReport.Output, deps.Env)
			if err != nil {
				return OpsOutput{}, err
			}
			// update Solana off ramp
			updateSolanaOffRampReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateSolanaOffRamp",
				semver.MustParse("1.0.0"),
				"Updates Solana OffRamps with EVM Chain Configs",
				func(b operations.Bundle, deps Dependencies, input solana.AddRemoteChainToOffRampConfig) (OpsOutput, error) {
					output, err := solana.AddRemoteChainToOffRamp(deps.Env, input)
					if err != nil {
						return OpsOutput{}, err
					}
					return OpsOutput{
						Proposals:   output.MCMSTimelockProposals,
						AddressBook: output.AddressBook, //nolint:staticcheck //SA1019 ignoring deprecated
						DataStore:   output.DataStore,
					}, nil
				},
			), deps, deps.changesetInput.solanaOffRampInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the output
			err = finalOutput.Merge(updateSolanaOffRampReport.Output, deps.Env)
			if err != nil {
				return OpsOutput{}, err
			}
			// update Solana fee quoter
			updateSolanaFeeQuoterReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateSolanaFeeQuoter",
				semver.MustParse("1.0.0"),
				"Updates Solana Fee Quoter with EVM Chain Configs",
				func(b operations.Bundle, deps Dependencies, input solana.AddRemoteChainToFeeQuoterConfig) (OpsOutput, error) {
					output, err := solana.AddRemoteChainToFeeQuoter(deps.Env, input)
					if err != nil {
						return OpsOutput{}, err
					}
					return OpsOutput{
						Proposals:   output.MCMSTimelockProposals,
						AddressBook: output.AddressBook, //nolint:staticcheck //SA1019 ignoring deprecated
						DataStore:   output.DataStore,
					}, nil
				},
			), deps, deps.changesetInput.solanaFeeQuoterInput)
			if err != nil {
				return OpsOutput{}, err
			}
			// merge the output
			err = finalOutput.Merge(updateSolanaFeeQuoterReport.Output, deps.Env)
			if err != nil {
				return OpsOutput{}, err
			}
			var mcmsCfg *cldfproposalutils.TimelockConfig
			if input.MCMSConfig != nil {
				mcmsCfg = input.MCMSConfig
			}
			// post ops where we merge all the proposals into one
			postOpsReport, err := operations.ExecuteOperation(b, postOps, deps, postOpsInput{
				MCMSConfig: mcmsCfg,
				Proposals:  finalOutput.Proposals,
			})
			return OpsOutput{
				Proposals:   postOpsReport.Output,
				AddressBook: finalOutput.AddressBook,
				DataStore:   finalOutput.DataStore,
			}, err
		},
	)
)

type Dependencies struct {
	Env             cldf.Environment
	EVMMCMSState    map[uint64]evmstate.MCMSWithTimelockState
	SolanaMCMSState map[uint64]solstate.MCMSWithTimelockState

	changesetInput csInputs
}

type postOpsInput struct {
	MCMSConfig *cldfproposalutils.TimelockConfig
	Proposals  []mcmslib.TimelockProposal
}

type OpsOutput struct {
	Proposals   []mcmslib.TimelockProposal
	AddressBook cldf.AddressBook
	// DataStore carries the lane refs the Solana sub-changesets record at the write site
	// (RemoteDest/RemoteSource, qualified by the remote EVM chain selector). It is merged
	// alongside the address book so the qualified refs reach the final changeset output
	// instead of being re-derived from the address book — which could not carry qualifiers.
	DataStore datastore.MutableDataStore
}

func (o *OpsOutput) Merge(other OpsOutput, env cldf.Environment) error {
	if o == nil {
		o = &OpsOutput{}
	}
	if o.AddressBook == nil {
		o.AddressBook = other.AddressBook
		if other.AddressBook != nil {
			if err := env.ExistingAddresses.Merge(other.AddressBook); err != nil {
				return fmt.Errorf("failed to merge existing addresses to environment: %w", err)
			}
		}
	} else if other.AddressBook != nil {
		if err := o.AddressBook.Merge(other.AddressBook); err != nil {
			return fmt.Errorf("failed to merge address book: %w", err)
		}
		if err := env.ExistingAddresses.Merge(other.AddressBook); err != nil {
			return fmt.Errorf("failed to merge existing addresses to environment: %w", err)
		}
	}
	if other.DataStore != nil {
		if o.DataStore == nil {
			o.DataStore = datastore.NewMemoryDataStore()
		}
		if err := o.DataStore.Merge(other.DataStore.Seal()); err != nil {
			return fmt.Errorf("failed to merge data store: %w", err)
		}
		// Publish to the environment as with the address book above, so an op later in the
		// sequence resolves the refs an earlier one recorded. A sealed environment store is
		// skipped: the refs still ride to the final output in o.DataStore.
		if env.DataStore != nil {
			mutable, ok := env.DataStore.Addresses().(datastore.MutableAddressRefStore)
			if !ok {
				if env.Logger != nil {
					env.Logger.Warnw("Environment data store is not writable; lane refs will not be visible through env.DataStore until the changeset output is applied")
				}
			} else {
				refs, err := other.DataStore.Addresses().Fetch()
				if err != nil {
					return fmt.Errorf("failed to read merged data store: %w", err)
				}
				for _, ref := range refs {
					if err := mutable.Upsert(ref); err != nil {
						return fmt.Errorf("failed to publish address ref to environment: %w", err)
					}
				}
			}
		}
	}
	if o.Proposals == nil {
		o.Proposals = other.Proposals
		return nil
	}
	if len(other.Proposals) > 0 {
		o.Proposals = append(o.Proposals, other.Proposals...)
	}
	return nil
}

type csInputs struct {
	evmOnRampInput             v1_6.UpdateOnRampDestsConfig
	evmFeeQuoterDestChainInput v1_6.UpdateFeeQuoterDestsConfig
	evmFeeQuoterPriceInput     v1_6.UpdateFeeQuoterPricesConfig
	evmOffRampInput            v1_6.UpdateOffRampSourcesConfig
	evmRouterInput             v1_6.UpdateRouterRampsConfig
	solanaRouterInput          solana.AddRemoteChainToRouterConfig
	solanaOffRampInput         solana.AddRemoteChainToOffRampConfig
	solanaFeeQuoterInput       solana.AddRemoteChainToFeeQuoterConfig
}

type AddMultiEVMSolanaLaneConfig struct {
	Configs             []AddRemoteChainE2EConfig
	SolanaChainSelector uint64
	MCMSConfig          *cldfproposalutils.TimelockConfig
}

// PlannedRefs returns the datastore refs this lane configuration writes: one RemoteDest and
// one RemoteSource per remote EVM chain, keyed on the Solana chain and qualified by the remote
// chain selector the same keys the write sites record in solana_v0_1_0. Declaring them up
// front means a duplicated lane config fails validation, before any instruction is sent.
func (multiCfg *AddMultiEVMSolanaLaneConfig) PlannedRefs() []datastore.AddressRef {
	refs := make([]datastore.AddressRef, 0, 2*len(multiCfg.Configs))
	version := semver.MustParse("1.0.0")
	for _, cfg := range multiCfg.Configs {
		qualifier := strconv.FormatUint(cfg.EVMChainSelector, 10)
		for _, laneType := range []datastore.ContractType{
			datastore.ContractType(shared.RemoteDest),
			datastore.ContractType(shared.RemoteSource),
		} {
			refs = append(refs, datastore.AddressRef{
				ChainSelector: multiCfg.SolanaChainSelector,
				Type:          laneType,
				Version:       version,
				Qualifier:     qualifier,
			})
		}
	}
	return refs
}

func (multiCfg *AddMultiEVMSolanaLaneConfig) populateAndValidateIndividualCSConfig(env cldf.Environment, evmState stateview.CCIPOnChainState) (csInputs, error) {
	var timelockConfig *cldfproposalutils.TimelockConfig
	var input csInputs
	if multiCfg.MCMSConfig != nil {
		timelockConfig = multiCfg.MCMSConfig
	}
	solChainSelector := multiCfg.SolanaChainSelector
	input.evmOnRampInput = v1_6.UpdateOnRampDestsConfig{
		MCMS: timelockConfig,
	}
	input.evmFeeQuoterDestChainInput = v1_6.UpdateFeeQuoterDestsConfig{
		MCMS: timelockConfig,
	}
	input.evmFeeQuoterDestChainInput = v1_6.UpdateFeeQuoterDestsConfig{
		MCMS: timelockConfig,
	}
	input.evmFeeQuoterPriceInput = v1_6.UpdateFeeQuoterPricesConfig{
		MCMS: timelockConfig,
	}
	input.evmOffRampInput = v1_6.UpdateOffRampSourcesConfig{
		MCMS: timelockConfig,
	}
	input.evmRouterInput = v1_6.UpdateRouterRampsConfig{
		MCMS: timelockConfig,
	}
	input.solanaRouterInput = solana.AddRemoteChainToRouterConfig{
		MCMS:          timelockConfig,
		ChainSelector: solChainSelector,
	}
	input.solanaOffRampInput = solana.AddRemoteChainToOffRampConfig{
		ChainSelector: solChainSelector,
		MCMS:          timelockConfig,
	}
	input.solanaFeeQuoterInput = solana.AddRemoteChainToFeeQuoterConfig{
		ChainSelector: solChainSelector,
		MCMS:          timelockConfig,
	}
	for _, cfg := range multiCfg.Configs {
		if input.evmOnRampInput.UpdatesByChain == nil {
			input.evmOnRampInput.UpdatesByChain = make(map[uint64]map[uint64]v1_6.OnRampDestinationUpdate)
		}
		input.evmOnRampInput.UpdatesByChain[cfg.EVMChainSelector] = map[uint64]v1_6.OnRampDestinationUpdate{
			solChainSelector: {
				IsEnabled:        true,
				TestRouter:       cfg.IsTestRouter,
				AllowListEnabled: cfg.EVMOnRampAllowListEnabled,
			},
		}
		if input.evmFeeQuoterDestChainInput.UpdatesByChain == nil {
			input.evmFeeQuoterDestChainInput.UpdatesByChain = make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
		}
		input.evmFeeQuoterDestChainInput.UpdatesByChain[cfg.EVMChainSelector] = map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			solChainSelector: cfg.EVMFeeQuoterDestChainInput,
		}
		if input.evmFeeQuoterDestChainInput.UpdatesByChain == nil {
			input.evmFeeQuoterDestChainInput.UpdatesByChain = make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
		}
		input.evmFeeQuoterDestChainInput.UpdatesByChain[cfg.EVMChainSelector] = map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			solChainSelector: cfg.EVMFeeQuoterDestChainInput,
		}
		if input.evmFeeQuoterPriceInput.PricesByChain == nil {
			input.evmFeeQuoterPriceInput.PricesByChain = make(map[uint64]v1_6.FeeQuoterPriceUpdatePerSource)
		}
		input.evmFeeQuoterPriceInput.PricesByChain[cfg.EVMChainSelector] = v1_6.FeeQuoterPriceUpdatePerSource{
			GasPrices: map[uint64]*big.Int{
				solChainSelector: cfg.InitialSolanaGasPriceForEVMFeeQuoter,
			},
			TokenPrices: cfg.InitialEVMTokenPricesForEVMFeeQuoter,
		}
		if input.evmOffRampInput.UpdatesByChain == nil {
			input.evmOffRampInput.UpdatesByChain = make(map[uint64]map[uint64]v1_6.OffRampSourceUpdate)
		}
		input.evmOffRampInput.UpdatesByChain[cfg.EVMChainSelector] = map[uint64]v1_6.OffRampSourceUpdate{
			solChainSelector: {
				IsEnabled:                 true,
				TestRouter:                cfg.IsTestRouter,
				IsRMNVerificationDisabled: cfg.IsRMNVerificationDisabledOnEVMOffRamp,
			},
		}
		if input.evmRouterInput.UpdatesByChain == nil {
			input.evmRouterInput.UpdatesByChain = make(map[uint64]v1_6.RouterUpdates)
		}
		updates := v1_6.RouterUpdates{
			OffRampUpdates: map[uint64]bool{
				solChainSelector: true,
			},
			OnRampUpdates: map[uint64]bool{
				solChainSelector: true,
			},
		}
		input.evmRouterInput.TestRouter = cfg.IsTestRouter
		input.evmRouterInput.UpdatesByChain[cfg.EVMChainSelector] = updates

		// router is always owned by deployer key so no need to pass MCMS
		if cfg.IsTestRouter {
			input.evmRouterInput.MCMS = nil
		}

		if input.solanaRouterInput.UpdatesByChain == nil {
			input.solanaRouterInput.UpdatesByChain = make(map[uint64]*solana.RouterConfig)
		}
		input.solanaRouterInput.UpdatesByChain[cfg.EVMChainSelector] = &cfg.SolanaRouterConfig
		if input.solanaOffRampInput.UpdatesByChain == nil {
			input.solanaOffRampInput.UpdatesByChain = make(map[uint64]*solana.OffRampConfig)
		}
		input.solanaOffRampInput.UpdatesByChain[cfg.EVMChainSelector] = &cfg.SolanaOffRampConfig
		if input.solanaFeeQuoterInput.UpdatesByChain == nil {
			input.solanaFeeQuoterInput.UpdatesByChain = make(map[uint64]*solana.FeeQuoterConfig)
		}
		input.solanaFeeQuoterInput.UpdatesByChain[cfg.EVMChainSelector] = &cfg.SolanaFeeQuoterConfig
	}

	if err := input.evmOnRampInput.Validate(env); err != nil {
		return input, fmt.Errorf("failed to validate evm on ramp input: %w", err)
	}
	if err := input.evmFeeQuoterDestChainInput.Validate(env); err != nil {
		return input, fmt.Errorf("failed to validate evm fee quoter dest chain input: %w", err)
	}
	if err := input.evmFeeQuoterPriceInput.Validate(env); err != nil {
		return input, fmt.Errorf("failed to validate evm fee quoter price input: %w", err)
	}
	if err := input.evmRouterInput.Validate(env, evmState); err != nil {
		return input, fmt.Errorf("failed to validate evm router input: %w", err)
	}
	if err := input.evmOffRampInput.Validate(env, evmState); err != nil {
		return input, fmt.Errorf("failed to validate evm off ramp input: %w", err)
	}
	if err := input.solanaRouterInput.Validate(env, evmState); err != nil {
		return input, fmt.Errorf("failed to validate solana router input: %w", err)
	}
	if err := input.solanaOffRampInput.Validate(env, evmState); err != nil {
		return input, fmt.Errorf("failed to validate solana off ramp input: %w", err)
	}
	if err := input.solanaFeeQuoterInput.Validate(env, evmState); err != nil {
		return input, fmt.Errorf("failed to validate solana fee quoter input: %w", err)
	}
	return input, nil
}

type AddRemoteChainE2EConfig struct {
	// inputs to be filled by user
	EVMChainSelector                      uint64
	IsTestRouter                          bool
	EVMOnRampAllowListEnabled             bool
	EVMFeeQuoterDestChainInput            fee_quoter.FeeQuoterDestChainConfig
	InitialSolanaGasPriceForEVMFeeQuoter  *big.Int
	InitialEVMTokenPricesForEVMFeeQuoter  map[common.Address]*big.Int
	IsRMNVerificationDisabledOnEVMOffRamp bool
	SolanaRouterConfig                    solana.RouterConfig
	SolanaOffRampConfig                   solana.OffRampConfig
	SolanaFeeQuoterConfig                 solana.FeeQuoterConfig
}

func addEVMSolanaPreconditions(env cldf.Environment, input AddMultiEVMSolanaLaneConfig) error {
	evmState, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain evm state: %w", err)
	}
	var timelockConfig *cldfproposalutils.TimelockConfig
	if input.MCMSConfig != nil {
		timelockConfig = input.MCMSConfig
	}
	for _, cfg := range input.Configs {
		// Verify evm Chain
		if err := stateview.ValidateChain(env, evmState, cfg.EVMChainSelector, timelockConfig); err != nil {
			return fmt.Errorf("failed to validate EVM chain %d: %w", cfg.EVMChainSelector, err)
		}
	}
	if _, ok := env.BlockChains.SolanaChains()[input.SolanaChainSelector]; !ok {
		return fmt.Errorf("failed to find Solana chain in env %d", input.SolanaChainSelector)
	}
	solanaState, err := stateview.LoadOnchainStateSolana(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain solana state: %w", err)
	}
	if _, exists := solanaState.SolChains[input.SolanaChainSelector]; !exists {
		return fmt.Errorf("failed to find Solana chain in state %d", input.SolanaChainSelector)
	}
	if _, err := shared.ReserveRefs(input.PlannedRefs()); err != nil {
		return err
	}
	if err := shared.ValidateAddressRefs(env, input.PlannedRefs()); err != nil {
		return fmt.Errorf("lane datastore refs conflict: %w", err)
	}
	return nil
}

func addEVMAndSolanaLaneLogic(env cldf.Environment, input AddMultiEVMSolanaLaneConfig) (cldf.ChangesetOutput, error) {
	evmState, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load evm onchain state: %w", err)
	}
	addresses, err := env.ExistingAddresses.AddressesForChain(input.SolanaChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get addresses for Solana chain: %w", err)
	}
	mcmState, err := solstate.MaybeLoadMCMSWithTimelockChainState(env.BlockChains.SolanaChains()[input.SolanaChainSelector], addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Solana MCMS state: %w", err)
	}
	if mcmState == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Solana MCMS state: %w", err)
	}
	// now populate individual inputs from the config
	changesetInputs, err := input.populateAndValidateIndividualCSConfig(env, evmState)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	env.Logger.Infow("router input", "input", changesetInputs.solanaRouterInput)
	deps := Dependencies{
		Env:          env,
		EVMMCMSState: evmState.EVMMCMSStateByChain(),
		SolanaMCMSState: map[uint64]solstate.MCMSWithTimelockState{
			input.SolanaChainSelector: *mcmState,
		},
		changesetInput: changesetInputs,
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, addEVMAndSolanaLaneSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute addEVMAndSolanaLane sequence: %w", err)
	}

	// The lane refs were recorded at the write site (RemoteDest/RemoteSource, qualified by the
	// remote EVM chain selector) and carried through the sequence's OpsOutput. Nothing is
	// re-derived from the address book, which could not carry the qualifiers.
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: report.Output.Proposals,
		AddressBook:           report.Output.AddressBook, //nolint:staticcheck // SA1019 AddressBook is deprecated
		DataStore:             report.Output.DataStore,
	}, nil
}
