package v1_6

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	AddEVMAndSolanaLaneChangeset = cldf.CreateChangeSet(addEVMAndSolanaLaneLogic, addEVMSolanaPreconditions)

	postOps = operations.NewOperation(
		"postOpsToAggregateProposals",
		semver.MustParse("1.0.0"),
		"Post ops to aggregate proposals",
		func(b operations.Bundle, deps Dependencies, input postOpsInput) (deployment.ChangesetOutput, error) {
			allProposals := input.ConsolidatedCSOutput.MCMSTimelockProposals
			proposal, err := proposalutils.AggregateProposals(
				deps.Env, deps.EVMMCMSState, deps.SolanaMCMSState, allProposals, nil,
				"Adding EVM and Solana lane", input.MCMSConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			if proposal != nil {
				input.ConsolidatedCSOutput.MCMSTimelockProposals = []mcmslib.TimelockProposal{*proposal}
			}
			return input.ConsolidatedCSOutput, nil
		},
	)

	addEVMAndSolanaLaneSequence = operations.NewSequence(
		"addEVMAndSolanaLane",
		semver.MustParse("1.0.0"),
		"Adds bi-directional lane between EVM chain and Solana",
		func(b operations.Bundle, deps Dependencies, input AddRemoteChainE2EConfig) (deployment.ChangesetOutput, error) {
			var finalOutput deployment.ChangesetOutput
			updateEVMOnRampReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMOnRamp",
				semver.MustParse("1.0.0"),
				"Updates EVM OnRamps with Destination Chain Configs for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateOnRampDestsConfig) (deployment.ChangesetOutput, error) {
					return v1_6.UpdateOnRampsDestsChangeset(deps.Env, input)
				},
			), deps, input.evmOnRampInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateEVMOnRampReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after EVMOnRampUpdate: %w", err)
			}
			// update EVM fee quoter dest chain
			updateEVMFeeQuoterDestChainReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMFeeQuoterDestChain",
				semver.MustParse("1.0.0"),
				"Updates EVM Fee Quoter with Destination Chain Configs for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateFeeQuoterDestsConfig) (deployment.ChangesetOutput, error) {
					return v1_6.UpdateFeeQuoterDestsChangeset(deps.Env, input)
				},
			), deps, input.evmFeeQuoterDestChainInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateEVMFeeQuoterDestChainReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after EVMFeeQuoterDestChainUpdate: %w", err)
			}

			// update EVM fee quoter prices
			updateEVMFeeQuoterPricesReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMFeeQuoterPrices",
				semver.MustParse("1.0.0"),
				"Updates EVM Fee Quoter with Prices for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateFeeQuoterPricesConfig) (deployment.ChangesetOutput, error) {
					return v1_6.UpdateFeeQuoterPricesChangeset(deps.Env, input)
				},
			), deps, input.evmFeeQuoterPriceInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateEVMFeeQuoterPricesReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after EVMFeeQuoterPricesUpdate: %w", err)
			}

			// update EVM off ramp
			updateEVMOffRampReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMOffRamp",
				semver.MustParse("1.0.0"),
				"Updates EVM OffRamps with Source Chain Configs for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateOffRampSourcesConfig) (deployment.ChangesetOutput, error) {
					return v1_6.UpdateOffRampSourcesChangeset(deps.Env, input)
				},
			), deps, input.evmOffRampInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateEVMOffRampReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after EVMOffRampUpdate: %w", err)
			}

			// update EVM router
			updateEVMRouterReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateEVMRouter",
				semver.MustParse("1.0.0"),
				"Updates EVM Router with onRamp and OffRamp for Solana",
				func(b operations.Bundle, deps Dependencies, input v1_6.UpdateRouterRampsConfig) (deployment.ChangesetOutput, error) {
					return v1_6.UpdateRouterRampsChangeset(deps.Env, input)
				},
			), deps, input.evmRouterInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateEVMRouterReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after EVMRouterUpdate: %w", err)
			}

			// update Solana router
			updateSolanaRouterReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateSolanaRouter",
				semver.MustParse("1.0.0"),
				"Updates Solana Router with EVM Chain Configs",
				func(b operations.Bundle, deps Dependencies, input solana.AddRemoteChainToRouterConfig) (deployment.ChangesetOutput, error) {
					return solana.AddRemoteChainToRouter(deps.Env, input)
				},
			), deps, input.solanaRouterInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateSolanaRouterReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after SolanaRouterUpdate: %w", err)
			}
			// update Solana off ramp
			updateSolanaOffRampReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateSolanaOffRamp",
				semver.MustParse("1.0.0"),
				"Updates Solana OffRamps with EVM Chain Configs",
				func(b operations.Bundle, deps Dependencies, input solana.AddRemoteChainToOffRampConfig) (deployment.ChangesetOutput, error) {
					return solana.AddRemoteChainToOffRamp(deps.Env, input)
				},
			), deps, input.solanaOffRampInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateSolanaOffRampReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after SolanaOffRampUpdate: %w", err)
			}

			// update Solana fee quoter
			updateSolanaFeeQuoterReport, err := operations.ExecuteOperation(b, operations.NewOperation(
				"updateSolanaFeeQuoter",
				semver.MustParse("1.0.0"),
				"Updates Solana Fee Quoter with EVM Chain Configs",
				func(b operations.Bundle, deps Dependencies, input solana.AddRemoteChainToFeeQuoterConfig) (deployment.ChangesetOutput, error) {
					return solana.AddRemoteChainToFeeQuoter(deps.Env, input)
				},
			), deps, input.solanaFeeQuoterInput)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			// merge the changeset outputs
			err = deployment.MergeChangesetOutput(deps.Env, &finalOutput, updateSolanaFeeQuoterReport.Output)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge changeset outputs after SolanaFeeQuoterUpdate: %w", err)
			}

			// post ops where we merge all the proposals into one
			postOpsReport, err := operations.ExecuteOperation(b, postOps, deps, postOpsInput{
				SolanaChainSelector:  input.SolanaChainSelector,
				EVMChainSelector:     input.EVMChainSelector,
				MCMSConfig:           input.MCMSConfig.MCMS,
				ConsolidatedCSOutput: finalOutput,
			})
			return postOpsReport.Output, nil
		},
	)
)

type Dependencies struct {
	Env             deployment.Environment
	EVMMCMSState    map[uint64]commonstate.MCMSWithTimelockState
	SolanaMCMSState map[uint64]commonstate.MCMSWithTimelockStateSolana
}

type postOpsInput struct {
	SolanaChainSelector  uint64
	EVMChainSelector     uint64
	MCMSConfig           *proposalutils.TimelockConfig
	ConsolidatedCSOutput deployment.ChangesetOutput
}

type AddRemoteChainE2EConfig struct {
	// inputs to be filled by user
	SolanaChainSelector                  uint64
	EVMChainSelector                     uint64
	IsTestRouter                         bool
	EVMOnRampAllowListEnabled            bool
	EVMFeeQuoterDestChainInput           fee_quoter.FeeQuoterDestChainConfig
	InitialSolanaGasPriceForEVMFeeQuoter *big.Int
	InitialEVMTokenPricesForEVMFeeQuoter map[common.Address]*big.Int
	IsRMNVerificationEnabledOnEVMOffRamp bool
	SolanaRouterConfig                   solana.RouterConfig
	SolanaOffRampConfig                  solana.OffRampConfig
	SolanaFeeQuoterConfig                solana.FeeQuoterConfig

	MCMSConfig *solana.MCMSConfigSolana

	// inputs not to be filled by user
	// but populated by the precondition function
	// we do this to avoid having mismatch of chain selectors in setting these inputs
	evmOnRampInput             v1_6.UpdateOnRampDestsConfig
	evmFeeQuoterDestChainInput v1_6.UpdateFeeQuoterDestsConfig
	evmFeeQuoterPriceInput     v1_6.UpdateFeeQuoterPricesConfig
	evmOffRampInput            v1_6.UpdateOffRampSourcesConfig
	evmRouterInput             v1_6.UpdateRouterRampsConfig
	solanaRouterInput          solana.AddRemoteChainToRouterConfig
	solanaOffRampInput         solana.AddRemoteChainToOffRampConfig
	solanaFeeQuoterInput       solana.AddRemoteChainToFeeQuoterConfig
}

func (input *AddRemoteChainE2EConfig) populateAndValidateIndividualCSConfig(env deployment.Environment, evmState ccipchangeset.CCIPOnChainState) error {
	var timelockConfig *proposalutils.TimelockConfig
	if input.MCMSConfig != nil {
		timelockConfig = input.MCMSConfig.MCMS
	}
	input.evmOnRampInput = v1_6.UpdateOnRampDestsConfig{
		MCMS: timelockConfig,
		UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
			input.EVMChainSelector: {
				input.SolanaChainSelector: {
					IsEnabled:        true,
					TestRouter:       input.IsTestRouter,
					AllowListEnabled: input.EVMOnRampAllowListEnabled,
				},
			},
		},
	}
	input.evmFeeQuoterDestChainInput = v1_6.UpdateFeeQuoterDestsConfig{
		MCMS: timelockConfig,
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			input.EVMChainSelector: {
				input.SolanaChainSelector: input.EVMFeeQuoterDestChainInput,
			},
		},
	}
	input.evmFeeQuoterPriceInput = v1_6.UpdateFeeQuoterPricesConfig{
		MCMS: timelockConfig,
		PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
			input.EVMChainSelector: {
				GasPrices: map[uint64]*big.Int{
					input.SolanaChainSelector: input.InitialSolanaGasPriceForEVMFeeQuoter,
				},
				TokenPrices: input.InitialEVMTokenPricesForEVMFeeQuoter,
			},
		},
	}
	input.evmOffRampInput = v1_6.UpdateOffRampSourcesConfig{
		MCMS: timelockConfig,
		UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
			input.EVMChainSelector: {
				input.SolanaChainSelector: {
					IsEnabled:                 true,
					TestRouter:                input.IsTestRouter,
					IsRMNVerificationDisabled: input.IsRMNVerificationEnabledOnEVMOffRamp,
				},
			},
		},
	}
	input.evmRouterInput = v1_6.UpdateRouterRampsConfig{
		MCMS:       timelockConfig,
		TestRouter: input.IsTestRouter,
		UpdatesByChain: map[uint64]v1_6.RouterUpdates{
			input.EVMChainSelector: {
				OffRampUpdates: map[uint64]bool{
					input.SolanaChainSelector: true,
				},
				OnRampUpdates: map[uint64]bool{
					input.SolanaChainSelector: true,
				},
			},
		},
	}
	input.solanaRouterInput = solana.AddRemoteChainToRouterConfig{
		ChainSelector: input.SolanaChainSelector,
		MCMSSolana:    input.MCMSConfig,
		UpdatesByChain: map[uint64]solana.RouterConfig{
			input.EVMChainSelector: input.SolanaRouterConfig,
		},
	}
	input.solanaOffRampInput = solana.AddRemoteChainToOffRampConfig{
		ChainSelector: input.SolanaChainSelector,
		MCMSSolana:    input.MCMSConfig,
		UpdatesByChain: map[uint64]solana.OffRampConfig{
			input.EVMChainSelector: input.SolanaOffRampConfig,
		},
	}
	input.solanaFeeQuoterInput = solana.AddRemoteChainToFeeQuoterConfig{
		ChainSelector: input.SolanaChainSelector,
		MCMSSolana:    input.MCMSConfig,
		UpdatesByChain: map[uint64]solana.FeeQuoterConfig{
			input.EVMChainSelector: input.SolanaFeeQuoterConfig,
		},
	}
	if err := input.evmOnRampInput.Validate(env); err != nil {
		return fmt.Errorf("failed to validate evm on ramp input: %w", err)
	}
	if err := input.evmFeeQuoterDestChainInput.Validate(env); err != nil {
		return fmt.Errorf("failed to validate evm fee quoter dest chain input: %w", err)
	}
	if err := input.evmFeeQuoterPriceInput.Validate(env); err != nil {
		return fmt.Errorf("failed to validate evm fee quoter price input: %w", err)
	}
	if err := input.evmRouterInput.Validate(env, evmState); err != nil {
		return fmt.Errorf("failed to validate evm router input: %w", err)
	}
	if err := input.evmOffRampInput.Validate(env, evmState); err != nil {
		return fmt.Errorf("failed to validate evm off ramp input: %w", err)
	}
	if err := input.solanaRouterInput.Validate(env); err != nil {
		return fmt.Errorf("failed to validate solana router input: %w", err)
	}
	if err := input.solanaOffRampInput.Validate(env); err != nil {
		return fmt.Errorf("failed to validate solana off ramp input: %w", err)
	}
	if err := input.solanaFeeQuoterInput.Validate(env); err != nil {
		return fmt.Errorf("failed to validate solana fee quoter input: %w", err)
	}
	return nil
}

func addEVMSolanaPreconditions(env deployment.Environment, input *AddRemoteChainE2EConfig) error {
	evmState, err := ccipchangeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain evm state: %w", err)
	}
	var timelockConfig *proposalutils.TimelockConfig
	if input.MCMSConfig != nil {
		timelockConfig = input.MCMSConfig.MCMS
	}
	// Verify evm Chain
	if err := ccipchangeset.ValidateChain(env, evmState, input.EVMChainSelector, timelockConfig); err != nil {
		return fmt.Errorf("failed to validate EVM chain %d: %w", input.EVMChainSelector, err)
	}
	if _, ok := env.SolChains[input.SolanaChainSelector]; !ok {
		return fmt.Errorf("failed to find Solana chain in env %d", input.SolanaChainSelector)
	}
	solanaState, err := ccipchangeset.LoadOnchainStateSolana(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain solana state: %w", err)
	}
	if _, exists := solanaState.SolChains[input.SolanaChainSelector]; !exists {
		return fmt.Errorf("failed to find Solana chain in state %d", input.SolanaChainSelector)
	}
	// now populate individual inputs from the config
	err = input.populateAndValidateIndividualCSConfig(env, evmState)
	if err != nil {
		return err
	}
	return nil
}

func addEVMAndSolanaLaneLogic(env deployment.Environment, input *AddRemoteChainE2EConfig) (deployment.ChangesetOutput, error) {
	evmState, err := ccipchangeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load evm onchain state: %w", err)
	}
	addresses, err := env.ExistingAddresses.AddressesForChain(input.SolanaChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get addresses for Solana chain: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(env.SolChains[input.SolanaChainSelector], addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load Solana MCMS state: %w", err)
	}
	if mcmState == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load Solana MCMS state: %w", err)
	}
	deps := Dependencies{
		Env:          env,
		EVMMCMSState: evmState.EVMMCMSStateByChain(),
		SolanaMCMSState: map[uint64]commonstate.MCMSWithTimelockStateSolana{
			input.SolanaChainSelector: *mcmState,
		},
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, addEVMAndSolanaLaneSequence, deps, *input)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to execute addEVMAndSolanaLane sequence: %w", err)
	}
	return report.Output, nil
}
