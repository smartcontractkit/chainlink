package v1_5

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
)

var _ deployment.ChangeSet[AddLanesConfig] = AddLanes

type AddLanesConfig struct {
	Configs []AddLaneConfig
}

func (c *AddLanesConfig) Validate(e deployment.Environment, state changeset.CCIPOnChainState) error {
	for _, cfg := range c.Configs {
		if err := cfg.Validate(e, state); err != nil {
			return err
		}
	}
	return nil
}

type AddLaneConfig struct {
	SourceChainSelector      uint64
	DestinationChainSelector uint64

	// onRamp specific configuration
	OnRampStaticCfg         evm_2_evm_onramp.EVM2EVMOnRampStaticConfig
	OnRampDynamicCfg        evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig
	OnRampFeeTokenArgs      []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs
	OnRampTransferTokenCfgs []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs
	OnRampNopsAndWeight     []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight
	OnRampRateLimiterCfg    evm_2_evm_onramp.RateLimiterConfig

	// offRamp specific configuration
	OffRampRateLimiterCfg evm_2_evm_offramp.RateLimiterConfig

	// Price Registry specific configuration
	InitialTokenPrices []price_registry_1_2_0.InternalTokenPriceUpdate
	GasPriceUpdates    []price_registry_1_2_0.InternalGasPriceUpdate

	// Job specific configuration
	TokenPricesUSDPipeline string
	PriceGetterConfigJson  string
	USDCAttestationAPI     string
	USDCCfg                *config.USDCConfig
}

func (c *AddLaneConfig) Validate(e deployment.Environment, state changeset.CCIPOnChainState) error {
	if err := deployment.IsValidChainSelector(c.SourceChainSelector); err != nil {
		return err
	}
	if err := deployment.IsValidChainSelector(c.DestinationChainSelector); err != nil {
		return err
	}
	sourceChain, exists := e.Chains[c.SourceChainSelector]
	if !exists {
		return fmt.Errorf("source chain %d not found in environment", c.SourceChainSelector)
	}
	destChain, exists := e.Chains[c.DestinationChainSelector]
	if !exists {
		return fmt.Errorf("destination chain %d not found in environment", c.DestinationChainSelector)
	}
	sourceChainState, exists := state.Chains[c.SourceChainSelector]
	if !exists {
		return fmt.Errorf("source chain %d not found in state", c.SourceChainSelector)
	}
	destChainState, exists := state.Chains[c.DestinationChainSelector]
	if !exists {
		return fmt.Errorf("destination chain %d not found in state", c.DestinationChainSelector)
	}
	// check for existing chain contracts on both source and destination chains
	if err := arePrerequisitesMet(sourceChainState, sourceChain); err != nil {
		return err
	}
	if err := arePrerequisitesMet(destChainState, destChain); err != nil {
		return err
	}
	// TODO: Add rest of the config validation
	return nil
}

func AddLanes(env deployment.Environment, c AddLanesConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load CCIP onchain state: %w", err)
	}
	if err := c.Validate(env, state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	blocksByDest := make(map[uint64]uint64)
	for _, cfg := range c.Configs {
		number, err := env.Chains[cfg.DestinationChainSelector].Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return deployment.ChangesetOutput{
				AddressBook: newAddresses,
			}, err
		}
		blocksByDest[cfg.DestinationChainSelector] = number.Number.Uint64()
		if err := addLane(env, state, newAddresses, cfg); err != nil {
			return deployment.ChangesetOutput{
				AddressBook: newAddresses,
			}, err
		}
	}
	jobSpecs, err := jobSpecsForLane(env, state, c, blocksByDest)
	if err != nil {
		return deployment.ChangesetOutput{
			AddressBook: newAddresses,
		}, err
	}
	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
		JobSpecs:    jobSpecs,
	}, nil
}

func jobSpecsForLane(
	env deployment.Environment,
	state changeset.CCIPOnChainState,
	addLanesCfg AddLanesConfig,
	blocksByDest map[uint64]uint64,
) (map[string][]string, error) {
	nodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil {
		return nil, err
	}
	nodesToJobSpecs := make(map[string][]string)
	for _, node := range nodes {
		var specs []string
		for _, cfg := range addLanesCfg.Configs {
			var err error
			destChainState := state.Chains[cfg.DestinationChainSelector]
			sourceChain := env.Chains[cfg.SourceChainSelector]
			destChain := env.Chains[cfg.DestinationChainSelector]
			destEVMChainIdStr, err := chain_selectors.GetChainIDFromSelector(cfg.DestinationChainSelector)
			if err != nil {
				return nil, fmt.Errorf("failed to get chain ID from selector for chain %s: %w", destChain.String(), err)
			}
			destEVMChainId, err := strconv.ParseUint(destEVMChainIdStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse chain ID %s for chain %s: %w", destEVMChainIdStr, destChain.String(), err)
			}
			ccipJobParam := integrationtesthelpers.CCIPJobSpecParams{
				OffRamp:                destChainState.EVM2EVMOffRamp[cfg.SourceChainSelector].Address(),
				CommitStore:            destChainState.CommitStore[cfg.SourceChainSelector].Address(),
				SourceChainName:        sourceChain.Name(),
				DestChainName:          destChain.Name(),
				DestEvmChainId:         destEVMChainId,
				TokenPricesUSDPipeline: cfg.TokenPricesUSDPipeline,
				PriceGetterConfig:      cfg.PriceGetterConfigJson,
				DestStartBlock:         blocksByDest[cfg.DestinationChainSelector],
				USDCAttestationAPI:     cfg.USDCAttestationAPI,
				USDCConfig:             cfg.USDCCfg,
				P2PV2Bootstrappers:     nodes.BootstrapLocators(),
			}
			if !node.IsBootstrap {
				commitSpec, err := ccipJobParam.CommitJobSpec()
				if err != nil {
					return nil, fmt.Errorf("failed to generate commit job spec for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				commitSpecStr, err := commitSpec.String()
				if err != nil {
					return nil, fmt.Errorf("failed to convert commit job spec to string for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				execSpec, err := ccipJobParam.ExecutionJobSpec()
				if err != nil {
					return nil, fmt.Errorf("failed to generate execution job spec for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				execSpecStr, err := execSpec.String()
				if err != nil {
					return nil, fmt.Errorf("failed to convert execution job spec to string for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				specs = append(specs, commitSpecStr, execSpecStr)
			} else {
				bootstrapSpec := ccipJobParam.BootstrapJob(destChainState.CommitStore[cfg.SourceChainSelector].Address().String())
				bootstrapSpecStr, err := bootstrapSpec.String()
				if err != nil {
					return nil, fmt.Errorf("failed to convert bootstrap job spec to string for source %s and destination %s: %w",
						sourceChain.String(), destChain.String(), err)
				}
				specs = append(specs, bootstrapSpecStr)
			}
		}
		nodesToJobSpecs[node.NodeID] = append(nodesToJobSpecs[node.NodeID], specs...)
	}
	return nodesToJobSpecs, nil
}

func addLane(e deployment.Environment, state changeset.CCIPOnChainState, ab deployment.AddressBook, cfg AddLaneConfig) error {
	// update prices on the source price registry
	sourceChainState := state.Chains[cfg.SourceChainSelector]
	destChainState := state.Chains[cfg.DestinationChainSelector]
	sourceChain := e.Chains[cfg.SourceChainSelector]
	destChain := e.Chains[cfg.DestinationChainSelector]
	sourcePriceReg := sourceChainState.PriceRegistry
	tx, err := sourcePriceReg.UpdatePrices(sourceChain.DeployerKey, price_registry_1_2_0.InternalPriceUpdates{
		TokenPriceUpdates: cfg.InitialTokenPrices,
		GasPriceUpdates:   cfg.GasPriceUpdates,
	})
	if err != nil {
		return err
	}
	_, err = sourceChain.Confirm(tx)
	if err != nil {
		return fmt.Errorf("failed to confirm price update tx for chain %s: %w", sourceChain.String(), deployment.MaybeDataErr(err))
	}
	// ================================================================
	// │                        Deploy Lane                           │
	// ================================================================
	// Deploy onRamp on source chain
	onRamp, onRampExists := sourceChainState.EVM2EVMOnRamp[cfg.DestinationChainSelector]
	if !onRampExists {
		onRampC, err := deployment.DeployContract(e.Logger, sourceChain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*evm_2_evm_onramp.EVM2EVMOnRamp] {
				onRampAddress, tx, onRampC, err2 := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
					sourceChain.DeployerKey,
					sourceChain.Client,
					cfg.OnRampStaticCfg,
					cfg.OnRampDynamicCfg,
					cfg.OnRampRateLimiterCfg,
					cfg.OnRampFeeTokenArgs,
					cfg.OnRampTransferTokenCfgs,
					cfg.OnRampNopsAndWeight,
				)
				return deployment.ContractDeploy[*evm_2_evm_onramp.EVM2EVMOnRamp]{
					Address: onRampAddress, Contract: onRampC, Tx: tx,
					Tv: deployment.NewTypeAndVersion(changeset.OnRamp, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy EVM2EVMOnRamp", "chain", sourceChain.String(), "err", err)
			return err
		}
		onRamp = onRampC.Contract
	} else {
		e.Logger.Infow("EVM2EVMOnRamp already exists",
			"source chain", sourceChain.String(), "destination chain", destChain.String(),
			"address", onRamp.Address().String())
	}

	// Deploy commit store on source chain
	commitStore, commitStoreExists := destChainState.CommitStore[cfg.SourceChainSelector]
	if !commitStoreExists {
		commitStoreC, err := deployment.DeployContract(e.Logger, sourceChain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*commit_store.CommitStore] {
				commitStoreAddress, tx, commitStoreC, err2 := commit_store.DeployCommitStore(
					destChain.DeployerKey,
					destChain.Client,
					commit_store.CommitStoreStaticConfig{
						ChainSelector:       destChain.Selector,
						SourceChainSelector: sourceChain.Selector,
						OnRamp:              onRamp.Address(),
						RmnProxy:            destChainState.RMNProxy.Address(),
					},
				)
				return deployment.ContractDeploy[*commit_store.CommitStore]{
					Address: commitStoreAddress, Contract: commitStoreC, Tx: tx,
					Tv: deployment.NewTypeAndVersion(changeset.CommitStore, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy CommitStore", "chain", sourceChain.String(), "err", err)
			return err
		}
		commitStore = commitStoreC.Contract
	} else {
		e.Logger.Infow("CommitStore already exists",
			"source chain", sourceChain.String(), "destination chain", destChain.String(),
			"address", commitStore.Address().String())
	}

	// Deploy offRamp on destination chain
	offRamp, offRampExists := destChainState.EVM2EVMOffRamp[cfg.SourceChainSelector]
	if !offRampExists {
		offRampC, err := deployment.DeployContract(e.Logger, sourceChain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*evm_2_evm_offramp.EVM2EVMOffRamp] {
				offRampAddress, tx, offRampC, err2 := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
					destChain.DeployerKey,
					destChain.Client,
					evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
						CommitStore:         commitStore.Address(),
						ChainSelector:       destChain.Selector,
						SourceChainSelector: sourceChain.Selector,
						OnRamp:              onRamp.Address(),
						PrevOffRamp:         common.HexToAddress(""),
						RmnProxy:            destChainState.RMNProxy.Address(), // RMN, formerly ARM
						TokenAdminRegistry:  destChainState.TokenAdminRegistry.Address(),
					},
					cfg.OffRampRateLimiterCfg,
				)
				return deployment.ContractDeploy[*evm_2_evm_offramp.EVM2EVMOffRamp]{
					Address: offRampAddress, Contract: offRampC, Tx: tx,
					Tv: deployment.NewTypeAndVersion(changeset.OffRamp, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy EVM2EVMOffRamp", "chain", sourceChain.String(), "err", err)
			return err
		}
		offRamp = offRampC.Contract
	} else {
		e.Logger.Infow("EVM2EVMOffRamp already exists",
			"source chain", sourceChain.String(), "destination chain", destChain.String(),
			"address", offRamp.Address().String())
	}

	// Apply Router updates
	tx, err = sourceChainState.Router.ApplyRampUpdates(sourceChain.DeployerKey,
		[]router.RouterOnRamp{{DestChainSelector: destChain.Selector, OnRamp: onRamp.Address()}}, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to apply router updates for source chain %s: %w", sourceChain.String(), deployment.MaybeDataErr(err))
	}
	_, err = sourceChain.Confirm(tx)
	if err != nil {
		return fmt.Errorf("failed to confirm router updates tx %s for source chain %s: %w", tx.Hash().String(), sourceChain.String(), deployment.MaybeDataErr(err))
	}

	tx, err = destChainState.Router.ApplyRampUpdates(destChain.DeployerKey,
		nil,
		nil,
		[]router.RouterOffRamp{{SourceChainSelector: sourceChain.Selector, OffRamp: offRamp.Address()}},
	)
	if err != nil {
		return fmt.Errorf("failed to apply router updates for destination chain %s: %w", destChain.String(), deployment.MaybeDataErr(err))
	}
	_, err = destChain.Confirm(tx)
	if err != nil {
		return fmt.Errorf("failed to confirm router updates tx %s for destination chain %s: %w", tx.Hash().String(), destChain.String(), deployment.MaybeDataErr(err))
	}

	// price registry updates
	_, err = destChainState.PriceRegistry.ApplyPriceUpdatersUpdates(
		destChain.DeployerKey,
		[]common.Address{commitStore.Address()},
		[]common.Address{},
	)
	if err != nil {
		return fmt.Errorf("failed to apply price registry updates for destination chain %s: %w", destChain.String(), deployment.MaybeDataErr(err))
	}
	_, err = destChain.Confirm(tx)
	if err != nil {
		return fmt.Errorf("failed to confirm price registry updates tx %s for destination chain %s: %w", tx.Hash().String(), destChain.String(), deployment.MaybeDataErr(err))
	}
	return nil
}

func arePrerequisitesMet(chainState changeset.CCIPChainState, chain deployment.Chain) error {
	if chainState.Router == nil {
		return fmt.Errorf("router not found for chain %s", chain.String())
	}
	if chainState.PriceRegistry == nil {
		return fmt.Errorf("price registry not found for chain %s", chain.String())
	}
	if chainState.RMN == nil && chainState.MockRMN == nil {
		return fmt.Errorf("neither RMN nor mockRMN found for chain %s", chain.String())
	}
	if chainState.Weth9 == nil {
		return fmt.Errorf("WETH9 not found for chain %s", chain.String())
	}
	if chainState.LinkToken == nil {
		return fmt.Errorf("LINK token not found for chain %s", chain.String())
	}
	if chainState.TokenAdminRegistry == nil {
		return fmt.Errorf("token admin registry not found for chain %s", chain.String())
	}
	if chainState.RMNProxy == nil {
		return fmt.Errorf("RMNProxy not found for chain %s", chain.String())
	}
	return nil
}
