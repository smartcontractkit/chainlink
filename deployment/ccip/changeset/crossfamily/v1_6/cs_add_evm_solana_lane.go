package v1_6

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	AddEVMAndSolanaLaneChangeset = cldf.CreateChangeSet(addEVMAndSolanaLaneLogic, addEVMSolanaPreconditions)

	postOps = operations.NewOperation(
		"postOpsToAggregateProposals",
		semver.MustParse("1.0.0"),
		"Post ops to aggregate proposals",
		func(b operations.Bundle, deps Dependencies, input postOpsInput) ([]mcmslib.TimelockProposal, error) {
			allProposals := input.Proposals
			proposal, err := proposalutils.AggregateProposals(
				deps.Env, deps.EVMMCMSState, deps.SolanaMCMSState, allProposals,
				"Adding EVM and Solana lane", input.MCMSConfig)
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
		func(b operations.Bundle, deps Dependencies, input AddRemoteChainE2EConfig) (OpsOutput, error) {
			deps.Env.Logger.Infow("Adding EVM and Solana lane", "EVMChainSelector", input.EVMChainSelector, "SolanaChainSelector", input.SolanaChainSelector)
			var finalOutput *OpsOutput
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
			var mcmsCfg *proposalutils.TimelockConfig
			if input.MCMSConfig != nil {
				mcmsCfg = input.MCMSConfig
			}
			// post ops where we merge all the proposals into one
			postOpsReport, err := operations.ExecuteOperation(b, postOps, deps, postOpsInput{
				SolanaChainSelector: input.SolanaChainSelector,
				EVMChainSelector:    input.EVMChainSelector,
				MCMSConfig:          mcmsCfg,
				Proposals:           finalOutput.Proposals,
			})
			return OpsOutput{
				Proposals:   postOpsReport.Output,
				AddressBook: finalOutput.AddressBook,
			}, err
		},
	)
)

type Dependencies struct {
	Env             cldf.Environment
	EVMMCMSState    map[uint64]commonstate.MCMSWithTimelockState
	SolanaMCMSState map[uint64]commonstate.MCMSWithTimelockStateSolana

	changesetInput csInputs
}

type postOpsInput struct {
	SolanaChainSelector uint64
	EVMChainSelector    uint64
	MCMSConfig          *proposalutils.TimelockConfig
	Proposals           []mcmslib.TimelockProposal
}

type OpsOutput struct {
	Proposals   []mcmslib.TimelockProposal
	AddressBook cldf.AddressBook
}

func (o *OpsOutput) Merge(other OpsOutput, env cldf.Environment) error {
	if o.AddressBook == nil {
		o.AddressBook = other.AddressBook
	} else if other.AddressBook != nil {
		if err := o.AddressBook.Merge(other.AddressBook); err != nil {
			return fmt.Errorf("failed to merge address book: %w", err)
		}
		if err := env.ExistingAddresses.Merge(other.AddressBook); err != nil {
			return fmt.Errorf("failed to merge existing addresses to environment: %w", err)
		}
	}
	o.Proposals = append(o.Proposals, other.Proposals...)
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

type AddRemoteChainE2EConfig struct {
	// inputs to be filled by user
	SolanaChainSelector                   uint64
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

	MCMSConfig *proposalutils.TimelockConfig
}

func (cfg *AddRemoteChainE2EConfig) populateAndValidateIndividualCSConfig(env cldf.Environment, evmState stateview.CCIPOnChainState) (csInputs, error) {
	var timelockConfig *proposalutils.TimelockConfig
	if cfg.MCMSConfig != nil {
		timelockConfig = cfg.MCMSConfig
	}
	var input csInputs
	input.evmOnRampInput = v1_6.UpdateOnRampDestsConfig{
		MCMS: timelockConfig,
		UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
			cfg.EVMChainSelector: {
				cfg.SolanaChainSelector: {
					IsEnabled:        true,
					TestRouter:       cfg.IsTestRouter,
					AllowListEnabled: cfg.EVMOnRampAllowListEnabled,
				},
			},
		},
	}
	input.evmFeeQuoterDestChainInput = v1_6.UpdateFeeQuoterDestsConfig{
		MCMS: timelockConfig,
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			cfg.EVMChainSelector: {
				cfg.SolanaChainSelector: cfg.EVMFeeQuoterDestChainInput,
			},
		},
	}
	input.evmFeeQuoterPriceInput = v1_6.UpdateFeeQuoterPricesConfig{
		MCMS: timelockConfig,
		PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
			cfg.EVMChainSelector: {
				GasPrices: map[uint64]*big.Int{
					cfg.SolanaChainSelector: cfg.InitialSolanaGasPriceForEVMFeeQuoter,
				},
				TokenPrices: cfg.InitialEVMTokenPricesForEVMFeeQuoter,
			},
		},
	}
	input.evmOffRampInput = v1_6.UpdateOffRampSourcesConfig{
		MCMS: timelockConfig,
		UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
			cfg.EVMChainSelector: {
				cfg.SolanaChainSelector: {
					IsEnabled:                 true,
					TestRouter:                cfg.IsTestRouter,
					IsRMNVerificationDisabled: cfg.IsRMNVerificationDisabledOnEVMOffRamp,
				},
			},
		},
	}
	input.evmRouterInput = v1_6.UpdateRouterRampsConfig{
		MCMS:       timelockConfig,
		TestRouter: cfg.IsTestRouter,
		UpdatesByChain: map[uint64]v1_6.RouterUpdates{
			cfg.EVMChainSelector: {
				OffRampUpdates: map[uint64]bool{
					cfg.SolanaChainSelector: true,
				},
				OnRampUpdates: map[uint64]bool{
					cfg.SolanaChainSelector: true,
				},
			},
		},
	}
	// router is always owned by deployer key so no need to pass MCMS
	if cfg.IsTestRouter {
		input.evmRouterInput.MCMS = nil
	}
	input.solanaRouterInput = solana.AddRemoteChainToRouterConfig{
		ChainSelector: cfg.SolanaChainSelector,
		MCMS:          cfg.MCMSConfig,
		UpdatesByChain: map[uint64]*solana.RouterConfig{
			cfg.EVMChainSelector: &cfg.SolanaRouterConfig,
		},
	}
	input.solanaOffRampInput = solana.AddRemoteChainToOffRampConfig{
		ChainSelector: cfg.SolanaChainSelector,
		MCMS:          cfg.MCMSConfig,
		UpdatesByChain: map[uint64]*solana.OffRampConfig{
			cfg.EVMChainSelector: &cfg.SolanaOffRampConfig,
		},
	}
	input.solanaFeeQuoterInput = solana.AddRemoteChainToFeeQuoterConfig{
		ChainSelector: cfg.SolanaChainSelector,
		MCMS:          cfg.MCMSConfig,
		UpdatesByChain: map[uint64]*solana.FeeQuoterConfig{
			cfg.EVMChainSelector: &cfg.SolanaFeeQuoterConfig,
		},
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
	if err := input.solanaRouterInput.Validate(env); err != nil {
		return input, fmt.Errorf("failed to validate solana router input: %w", err)
	}
	if err := input.solanaOffRampInput.Validate(env); err != nil {
		return input, fmt.Errorf("failed to validate solana off ramp input: %w", err)
	}
	if err := input.solanaFeeQuoterInput.Validate(env); err != nil {
		return input, fmt.Errorf("failed to validate solana fee quoter input: %w", err)
	}
	return input, nil
}

func addEVMSolanaPreconditions(env cldf.Environment, input AddRemoteChainE2EConfig) error {
	evmState, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain evm state: %w", err)
	}
	var timelockConfig *proposalutils.TimelockConfig
	if input.MCMSConfig != nil {
		timelockConfig = input.MCMSConfig
	}
	// Verify evm Chain
	if err := stateview.ValidateChain(env, evmState, input.EVMChainSelector, timelockConfig); err != nil {
		return fmt.Errorf("failed to validate EVM chain %d: %w", input.EVMChainSelector, err)
	}
	if _, ok := env.SolChains[input.SolanaChainSelector]; !ok {
		return fmt.Errorf("failed to find Solana chain in env %d", input.SolanaChainSelector)
	}
	solanaState, err := stateview.LoadOnchainStateSolana(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain solana state: %w", err)
	}
	if _, exists := solanaState.SolChains[input.SolanaChainSelector]; !exists {
		return fmt.Errorf("failed to find Solana chain in state %d", input.SolanaChainSelector)
	}
	return nil
}

func addEVMAndSolanaLaneLogic(env cldf.Environment, input AddRemoteChainE2EConfig) (cldf.ChangesetOutput, error) {
	evmState, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load evm onchain state: %w", err)
	}
	addresses, err := env.ExistingAddresses.AddressesForChain(input.SolanaChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get addresses for Solana chain: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(env.SolChains[input.SolanaChainSelector], addresses)
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
		SolanaMCMSState: map[uint64]commonstate.MCMSWithTimelockStateSolana{
			input.SolanaChainSelector: *mcmState,
		},
		changesetInput: changesetInputs,
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, addEVMAndSolanaLaneSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute addEVMAndSolanaLane sequence: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: report.Output.Proposals,
		AddressBook:           report.Output.AddressBook,
	}, nil
}
