// Package v1_5 contains CCIP 1.5 deployment test helpers: lane contracts and OCR2 config only (no Chainlink job specs).
package v1_5

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	ccipcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	deploycciptesthelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/utils/abihelpers"
)

// --- inlined helpers (OCR2 encoding for v1.5 lane tests only)

type jsonCommitOffchainConfig struct {
	SourceFinalityDepth      uint32
	DestFinalityDepth        uint32
	GasPriceHeartBeat        config.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      config.Duration
	TokenPriceDeviationPPB   uint32
	InflightCacheExpiry      config.Duration
	PriceReportingDisabled   bool
}

func (c jsonCommitOffchainConfig) Validate() error {
	if c.GasPriceHeartBeat.Duration() == 0 {
		return errors.New("must set GasPriceHeartBeat")
	}
	if c.ExecGasPriceDeviationPPB == 0 {
		return errors.New("must set ExecGasPriceDeviationPPB")
	}
	if c.TokenPriceHeartBeat.Duration() == 0 {
		return errors.New("must set TokenPriceHeartBeat")
	}
	if c.TokenPriceDeviationPPB == 0 {
		return errors.New("must set TokenPriceDeviationPPB")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	return nil
}

type commitOnchainCfg commit_store.CommitStoreDynamicConfig

func (c commitOnchainCfg) AbiString() string {
	return `
	[
		{
			"components": [
				{"name": "priceRegistry", "type": "address"}
			],
			"type": "tuple"
		}
	]`
}

type jsonExecOffchainConfig struct {
	SourceFinalityDepth         uint32
	DestOptimisticConfirmations uint32
	DestFinalityDepth           uint32
	BatchGasLimit               uint32
	RelativeBoostPerWaitHour    float64
	InflightCacheExpiry         config.Duration
	RootSnoozeTime              config.Duration
	BatchingStrategyID          uint32
	MessageVisibilityInterval   config.Duration
}

func (c jsonExecOffchainConfig) Validate() error {
	if c.DestOptimisticConfirmations == 0 {
		return errors.New("must set DestOptimisticConfirmations")
	}
	if c.BatchGasLimit == 0 {
		return errors.New("must set BatchGasLimit")
	}
	if c.RelativeBoostPerWaitHour == 0 {
		return errors.New("must set RelativeBoostPerWaitHour")
	}
	if c.InflightCacheExpiry.Duration() == 0 {
		return errors.New("must set InflightCacheExpiry")
	}
	if c.RootSnoozeTime.Duration() == 0 {
		return errors.New("must set RootSnoozeTime")
	}
	return nil
}

type execOnchainCfg evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig

func (d execOnchainCfg) AbiString() string {
	return `
	[
		{
			"components": [
				{"name": "permissionLessExecutionThresholdSeconds", "type": "uint32"},
				{"name": "maxDataBytes", "type": "uint32"},
				{"name": "maxNumberOfTokensPerMsg", "type": "uint16"},
				{"name": "router", "type": "address"},
				{"name": "priceRegistry", "type": "address"}
			],
			"type": "tuple"
		}
	]`
}

func linkUSDWei(amount int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(amount))
}

func encodeCommitOffchainBytes(
	gasPriceHeartBeat config.Duration,
	daGasPriceDeviationPPB uint32,
	execGasPriceDeviationPPB uint32,
	tokenPriceHeartBeat config.Duration,
	tokenPriceDeviationPPB uint32,
	inflightCacheExpiry config.Duration,
	priceReportingDisabled bool,
) ([]byte, error) {
	j := jsonCommitOffchainConfig{
		GasPriceHeartBeat:        gasPriceHeartBeat,
		DAGasPriceDeviationPPB:   daGasPriceDeviationPPB,
		ExecGasPriceDeviationPPB: execGasPriceDeviationPPB,
		TokenPriceHeartBeat:      tokenPriceHeartBeat,
		TokenPriceDeviationPPB:   tokenPriceDeviationPPB,
		InflightCacheExpiry:      inflightCacheExpiry,
		PriceReportingDisabled:   priceReportingDisabled,
	}
	if err := j.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(j)
}

func encodeExecOffchainBytes(
	destOptimisticConfirmations uint32,
	batchGasLimit uint32,
	relativeBoostPerWaitHour float64,
	inflightCacheExpiry config.Duration,
	rootSnoozeTime config.Duration,
	batchingStrategyID uint32,
	messageVisibilityInterval config.Duration,
) ([]byte, error) {
	j := jsonExecOffchainConfig{
		DestOptimisticConfirmations: destOptimisticConfirmations,
		BatchGasLimit:               batchGasLimit,
		RelativeBoostPerWaitHour:    relativeBoostPerWaitHour,
		InflightCacheExpiry:         inflightCacheExpiry,
		RootSnoozeTime:              rootSnoozeTime,
		BatchingStrategyID:          batchingStrategyID,
		MessageVisibilityInterval:   messageVisibilityInterval,
	}
	if err := j.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(j)
}


// --- lane contracts (from cs_lane_contracts.go) ---

var _ cldf.ChangeSet[DeployLanesConfig] = DeployLanesChangeset

type DeployLanesConfig struct {
	Configs []DeployLaneConfig
}

func (c *DeployLanesConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	for _, cfg := range c.Configs {
		if err := cfg.Validate(e, state); err != nil {
			return err
		}
	}
	return nil
}

type DeployLaneConfig struct {
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
}

func (c *DeployLaneConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(c.SourceChainSelector); err != nil {
		return err
	}
	if err := cldf.IsValidChainSelector(c.DestinationChainSelector); err != nil {
		return err
	}
	sourceChain, exists := e.BlockChains.EVMChains()[c.SourceChainSelector]
	if !exists {
		return fmt.Errorf("source chain %d not found in environment", c.SourceChainSelector)
	}
	destChain, exists := e.BlockChains.EVMChains()[c.DestinationChainSelector]
	if !exists {
		return fmt.Errorf("destination chain %d not found in environment", c.DestinationChainSelector)
	}
	sourceChainState, exists := state.EVMChainState(c.SourceChainSelector)
	if !exists {
		return fmt.Errorf("source chain %d not found in state", c.SourceChainSelector)
	}
	destChainState, exists := state.EVMChainState(c.DestinationChainSelector)
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

func (c *DeployLaneConfig) populateAddresses(state stateview.CCIPOnChainState) error {
	sourceChainState := state.MustGetEVMChainState(c.SourceChainSelector)
	srcLink, err := sourceChainState.LinkTokenAddress()
	if err != nil {
		return fmt.Errorf("failed to get LINK token address for source chain %d: %w", c.SourceChainSelector, err)
	}
	c.OnRampStaticCfg.LinkToken = srcLink
	c.OnRampStaticCfg.RmnProxy = sourceChainState.RMNProxy.Address()
	c.OnRampStaticCfg.TokenAdminRegistry = sourceChainState.TokenAdminRegistry.Address()

	c.OnRampDynamicCfg.Router = sourceChainState.Router.Address()
	c.OnRampDynamicCfg.PriceRegistry = sourceChainState.PriceRegistry.Address()
	return nil
}

func DeployLanesChangeset(env cldf.Environment, c DeployLanesConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load CCIP onchain state: %w", err)
	}
	if err := c.Validate(env, state); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	// populate addresses from the state
	for i := range c.Configs {
		if err := c.Configs[i].populateAddresses(state); err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	newAddresses := cldf.NewMemoryAddressBook()
	for _, cfg := range c.Configs {
		if err := deployLane(env, state, newAddresses, cfg); err != nil {
			return cldf.ChangesetOutput{
				AddressBook: newAddresses,
			}, err
		}
	}
	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

func deployLane(e cldf.Environment, state stateview.CCIPOnChainState, ab cldf.AddressBook, cfg DeployLaneConfig) error {
	// update prices on the source price registry
	sourceChainState := state.MustGetEVMChainState(cfg.SourceChainSelector)
	destChainState := state.MustGetEVMChainState(cfg.DestinationChainSelector)
	sourceChain := e.BlockChains.EVMChains()[cfg.SourceChainSelector]
	destChain := e.BlockChains.EVMChains()[cfg.DestinationChainSelector]
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
		return fmt.Errorf("failed to confirm price update tx for chain %s: %w", sourceChain.String(), cldf.MaybeDataErr(err))
	}
	// ================================================================
	// │                        Deploy Lane                           │
	// ================================================================
	// Deploy onRamp on source chain
	onRamp, onRampExists := sourceChainState.EVM2EVMOnRamp[cfg.DestinationChainSelector]
	if !onRampExists {
		onRampC, err := cldf.DeployContract(e.Logger, sourceChain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*evm_2_evm_onramp.EVM2EVMOnRamp] {
				onRampAddress, tx2, onRampC, err2 := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
					sourceChain.DeployerKey,
					sourceChain.Client,
					cfg.OnRampStaticCfg,
					cfg.OnRampDynamicCfg,
					cfg.OnRampRateLimiterCfg,
					cfg.OnRampFeeTokenArgs,
					cfg.OnRampTransferTokenCfgs,
					cfg.OnRampNopsAndWeight,
				)
				return cldf.ContractDeploy[*evm_2_evm_onramp.EVM2EVMOnRamp]{
					Address: onRampAddress, Contract: onRampC, Tx: tx2,
					Tv: cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_5_0), Err: err2,
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
		commitStoreC, err := cldf.DeployContract(e.Logger, destChain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*commit_store.CommitStore] {
				commitStoreAddress, tx2, commitStoreC, err2 := commit_store.DeployCommitStore(
					destChain.DeployerKey,
					destChain.Client,
					commit_store.CommitStoreStaticConfig{
						ChainSelector:       destChain.Selector,
						SourceChainSelector: sourceChain.Selector,
						OnRamp:              onRamp.Address(),
						RmnProxy:            destChainState.RMNProxy.Address(),
					},
				)
				return cldf.ContractDeploy[*commit_store.CommitStore]{
					Address: commitStoreAddress, Contract: commitStoreC, Tx: tx2,
					Tv: cldf.NewTypeAndVersion(shared.CommitStore, deployment.Version1_5_0), Err: err2,
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
		offRampC, err := cldf.DeployContract(e.Logger, destChain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*evm_2_evm_offramp.EVM2EVMOffRamp] {
				offRampAddress, tx2, offRampC, err2 := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
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
				return cldf.ContractDeploy[*evm_2_evm_offramp.EVM2EVMOffRamp]{
					Address: offRampAddress, Contract: offRampC, Tx: tx2,
					Tv: cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_5_0), Err: err2,
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
		return fmt.Errorf("failed to apply router updates for source chain %s: %w", sourceChain.String(), cldf.MaybeDataErr(err))
	}
	_, err = sourceChain.Confirm(tx)
	if err != nil {
		return fmt.Errorf("failed to confirm router updates tx %s for source chain %s: %w", tx.Hash().String(), sourceChain.String(), cldf.MaybeDataErr(err))
	}

	tx, err = destChainState.Router.ApplyRampUpdates(destChain.DeployerKey,
		nil,
		nil,
		[]router.RouterOffRamp{{SourceChainSelector: sourceChain.Selector, OffRamp: offRamp.Address()}},
	)
	if err != nil {
		return fmt.Errorf("failed to apply router updates for destination chain %s: %w", destChain.String(), cldf.MaybeDataErr(err))
	}
	_, err = destChain.Confirm(tx)
	if err != nil {
		return fmt.Errorf("failed to confirm router updates tx %s for destination chain %s: %w", tx.Hash().String(), destChain.String(), cldf.MaybeDataErr(err))
	}

	// price registry updates
	prTx, err := destChainState.PriceRegistry.ApplyPriceUpdatersUpdates(
		destChain.DeployerKey,
		[]common.Address{commitStore.Address()},
		[]common.Address{},
	)
	if err != nil {
		return fmt.Errorf("failed to apply price registry updates for destination chain %s: %w", destChain.String(), cldf.MaybeDataErr(err))
	}
	_, err = destChain.Confirm(prTx)
	if err != nil {
		return fmt.Errorf("failed to confirm price registry updates tx %s for destination chain %s: %w", prTx.Hash().String(), destChain.String(), cldf.MaybeDataErr(err))
	}
	return nil
}

func arePrerequisitesMet(chainState evm.CCIPChainState, chain cldf_evm.Chain) error {
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
	if _, err := chainState.LinkTokenAddress(); err != nil {
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

// --- OCR2 config (from cs_ocr2_config.go) ---

var _ cldf.ChangeSet[OCR2Config] = SetOCR2ConfigForTestChangeset

type FinalOCR2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type CommitOCR2ConfigParams struct {
	DestinationChainSelector uint64
	SourceChainSelector      uint64
	OCR2ConfigParams         confighelper.PublicConfig
	GasPriceHeartBeat        config.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      config.Duration
	TokenPriceDeviationPPB   uint32
	InflightCacheExpiry      config.Duration
	PriceReportingDisabled   bool
}

func (c *CommitOCR2ConfigParams) PopulateOffChainAndOnChainCfg(priceReg common.Address) error {
	var err error
	c.OCR2ConfigParams.ReportingPluginConfig, err = encodeCommitOffchainBytes(
		c.GasPriceHeartBeat,
		c.DAGasPriceDeviationPPB,
		c.ExecGasPriceDeviationPPB,
		c.TokenPriceHeartBeat,
		c.TokenPriceDeviationPPB,
		c.InflightCacheExpiry,
		c.PriceReportingDisabled,
	)
	if err != nil {
		return fmt.Errorf("failed to encode offchain config for source chain %d and destination chain %d: %w",
			c.SourceChainSelector, c.DestinationChainSelector, err)
	}
	c.OCR2ConfigParams.OnchainConfig, err = abihelpers.EncodeAbiStruct(commitOnchainCfg{PriceRegistry: priceReg})
	if err != nil {
		return fmt.Errorf("failed to encode onchain config for source chain %d and destination chain %d: %w",
			c.SourceChainSelector, c.DestinationChainSelector, err)
	}
	return nil
}

func (c *CommitOCR2ConfigParams) Validate(state stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(c.DestinationChainSelector); err != nil {
		return fmt.Errorf("invalid DestinationChainSelector: %w", err)
	}
	if err := cldf.IsValidChainSelector(c.SourceChainSelector); err != nil {
		return fmt.Errorf("invalid SourceChainSelector: %w", err)
	}

	chain, exists := state.EVMChainState(c.DestinationChainSelector)
	if !exists {
		return fmt.Errorf("chain %d does not exist in state", c.DestinationChainSelector)
	}
	if chain.CommitStore == nil {
		return fmt.Errorf("chain %d does not have a commit store", c.DestinationChainSelector)
	}
	_, exists = chain.CommitStore[c.SourceChainSelector]
	if !exists {
		return fmt.Errorf("chain %d does not have a commit store for source chain %d", c.DestinationChainSelector, c.SourceChainSelector)
	}
	if chain.PriceRegistry == nil {
		return fmt.Errorf("chain %d does not have a price registry", c.DestinationChainSelector)
	}
	return nil
}

type ExecuteOCR2ConfigParams struct {
	DestinationChainSelector    uint64
	SourceChainSelector         uint64
	DestOptimisticConfirmations uint32
	BatchGasLimit               uint32
	RelativeBoostPerWaitHour    float64
	InflightCacheExpiry         config.Duration
	RootSnoozeTime              config.Duration
	BatchingStrategyID          uint32
	MessageVisibilityInterval   config.Duration
	ExecOnchainConfig           evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig
	OCR2ConfigParams            confighelper.PublicConfig
}

func (e *ExecuteOCR2ConfigParams) PopulateOffChainAndOnChainCfg(router, priceReg common.Address) error {
	var err error
	e.OCR2ConfigParams.ReportingPluginConfig, err = encodeExecOffchainBytes(
		e.DestOptimisticConfirmations,
		e.BatchGasLimit,
		e.RelativeBoostPerWaitHour,
		e.InflightCacheExpiry,
		e.RootSnoozeTime,
		e.BatchingStrategyID,
		e.MessageVisibilityInterval,
	)
	if err != nil {
		return fmt.Errorf("failed to encode offchain config for exec plugin, source chain %d dest chain %d :%w",
			e.SourceChainSelector, e.DestinationChainSelector, err)
	}
	e.OCR2ConfigParams.OnchainConfig, err = abihelpers.EncodeAbiStruct(execOnchainCfg{
		PermissionLessExecutionThresholdSeconds: e.ExecOnchainConfig.PermissionLessExecutionThresholdSeconds,
		Router:                  router,
		PriceRegistry:           priceReg,
		MaxNumberOfTokensPerMsg: e.ExecOnchainConfig.MaxNumberOfTokensPerMsg,
		MaxDataBytes:            e.ExecOnchainConfig.MaxDataBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to encode onchain config for exec plugin, source chain %d dest chain %d :%w",
			e.SourceChainSelector, e.DestinationChainSelector, err)
	}
	return nil
}

func (e *ExecuteOCR2ConfigParams) Validate(state stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(e.SourceChainSelector); err != nil {
		return fmt.Errorf("invalid SourceChainSelector: %w", err)
	}
	if err := cldf.IsValidChainSelector(e.DestinationChainSelector); err != nil {
		return fmt.Errorf("invalid DestinationChainSelector: %w", err)
	}
	chain, exists := state.EVMChainState(e.DestinationChainSelector)
	if !exists {
		return fmt.Errorf("chain %d does not exist in state", e.DestinationChainSelector)
	}
	if chain.EVM2EVMOffRamp == nil {
		return fmt.Errorf("chain %d does not have an EVM2EVMOffRamp", e.DestinationChainSelector)
	}
	_, exists = chain.EVM2EVMOffRamp[e.SourceChainSelector]
	if !exists {
		return fmt.Errorf("chain %d does not have an EVM2EVMOffRamp for source chain %d", e.DestinationChainSelector, e.SourceChainSelector)
	}
	if chain.PriceRegistry == nil {
		return fmt.Errorf("chain %d does not have a price registry", e.DestinationChainSelector)
	}
	if chain.Router == nil {
		return fmt.Errorf("chain %d does not have a router", e.DestinationChainSelector)
	}
	return nil
}

type OCR2Config struct {
	CommitConfigs []CommitOCR2ConfigParams
	ExecConfigs   []ExecuteOCR2ConfigParams
}

func (o OCR2Config) Validate(state stateview.CCIPOnChainState) error {
	for _, c := range o.CommitConfigs {
		if err := c.Validate(state); err != nil {
			return err
		}
	}
	for _, e := range o.ExecConfigs {
		if err := e.Validate(state); err != nil {
			return err
		}
	}
	return nil
}

// SetOCR2ConfigForTestChangeset sets the OCR2 config on the chain for commit and offramp
// This is currently not suitable for prod environments it's only for testing
func SetOCR2ConfigForTestChangeset(env cldf.Environment, c OCR2Config) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load CCIP onchain state: %w", err)
	}
	if err := c.Validate(state); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid OCR2 config: %w", err)
	}
	for _, commit := range c.CommitConfigs {
		if err := commit.PopulateOffChainAndOnChainCfg(state.MustGetEVMChainState(commit.DestinationChainSelector).PriceRegistry.Address()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate offchain and onchain config for commit: %w", err)
		}
		finalCfg, err := deriveOCR2Config(env, commit.DestinationChainSelector, commit.OCR2ConfigParams)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to derive OCR2 config for commit: %w", err)
		}
		commitStore := state.MustGetEVMChainState(commit.DestinationChainSelector).CommitStore[commit.SourceChainSelector]
		chain := env.BlockChains.EVMChains()[commit.DestinationChainSelector]
		tx, err := commitStore.SetOCR2Config(
			chain.DeployerKey,
			finalCfg.Signers,
			finalCfg.Transmitters,
			finalCfg.F,
			finalCfg.OnchainConfig,
			finalCfg.OffchainConfigVersion,
			finalCfg.OffchainConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set OCR2 config for commit store %s on chain %s: %w",
				commitStore.Address().String(), chain.String(), cldf.MaybeDataErr(err))
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm OCR2 for commit store %s config on chain %s: %w",
				commitStore.Address().String(), chain.String(), err)
		}
	}
	for _, exec := range c.ExecConfigs {
		if err := exec.PopulateOffChainAndOnChainCfg(
			state.MustGetEVMChainState(exec.DestinationChainSelector).Router.Address(),
			state.MustGetEVMChainState(exec.DestinationChainSelector).PriceRegistry.Address()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate offchain and onchain config for offramp: %w", err)
		}
		finalCfg, err := deriveOCR2Config(env, exec.DestinationChainSelector, exec.OCR2ConfigParams)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to derive OCR2 config for offramp: %w", err)
		}
		offRamp := state.MustGetEVMChainState(exec.DestinationChainSelector).EVM2EVMOffRamp[exec.SourceChainSelector]
		chain := env.BlockChains.EVMChains()[exec.DestinationChainSelector]
		tx, err := offRamp.SetOCR2Config(
			chain.DeployerKey,
			finalCfg.Signers,
			finalCfg.Transmitters,
			finalCfg.F,
			finalCfg.OnchainConfig,
			finalCfg.OffchainConfigVersion,
			finalCfg.OffchainConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set OCR2 config for offramp %s on chain %s: %w",
				offRamp.Address().String(), chain.String(), err)
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm OCR2 for offramp %s config on chain %s: %w",
				offRamp.Address().String(), chain.String(), err)
		}
	}
	return cldf.ChangesetOutput{}, nil
}

func deriveOCR2Config(
	env cldf.Environment,
	chainSel uint64,
	ocrParams confighelper.PublicConfig,
) (FinalOCR2Config, error) {
	nodeInfo, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil {
		return FinalOCR2Config{}, fmt.Errorf("failed to get node info: %w", err)
	}
	nodes := nodeInfo.NonBootstraps()
	// Get OCR3 Config from helper
	var schedule []int
	var oracles []confighelper.OracleIdentityExtra
	for _, node := range nodes {
		schedule = append(schedule, 1)
		cfg, exists := node.OCRConfigForChainSelector(chainSel)
		if !exists {
			return FinalOCR2Config{}, fmt.Errorf("no OCR config for chain %d", chainSel)
		}
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  cfg.OnchainPublicKey,
				TransmitAccount:   cfg.TransmitAccount,
				OffchainPublicKey: cfg.OffchainPublicKey,
				PeerID:            cfg.PeerID.Raw(),
			},
			ConfigEncryptionPublicKey: cfg.ConfigEncryptionPublicKey,
		})
	}

	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		ocrParams.DeltaProgress,
		ocrParams.DeltaResend,
		ocrParams.DeltaRound,
		ocrParams.DeltaGrace,
		ocrParams.DeltaStage,
		ocrParams.RMax,
		schedule,
		oracles,
		ocrParams.ReportingPluginConfig,
		nil,
		ocrParams.MaxDurationQuery,
		ocrParams.MaxDurationObservation,
		ocrParams.MaxDurationReport,
		ocrParams.MaxDurationShouldAcceptFinalizedReport,
		ocrParams.MaxDurationShouldTransmitAcceptedReport,
		int(nodes.DefaultF()),
		ocrParams.OnchainConfig,
	)
	if err != nil {
		return FinalOCR2Config{}, fmt.Errorf("failed to derive OCR2 config: %w", err)
	}
	var signersAddresses []common.Address
	for _, signer := range signers {
		if len(signer) != 20 {
			return FinalOCR2Config{}, fmt.Errorf("address is not 20 bytes %s", signer)
		}
		signersAddresses = append(signersAddresses, common.BytesToAddress(signer))
	}
	var transmittersAddresses []common.Address
	for _, transmitter := range transmitters {
		bytes, err := hexutil.Decode(string(transmitter))
		if err != nil {
			return FinalOCR2Config{}, fmt.Errorf("given address is not valid %s: %w", transmitter, err)
		}
		if len(bytes) != 20 {
			return FinalOCR2Config{}, fmt.Errorf("address is not 20 bytes %s", transmitter)
		}
		transmittersAddresses = append(transmittersAddresses, common.BytesToAddress(bytes))
	}
	return FinalOCR2Config{
		Signers:               signersAddresses,
		Transmitters:          transmittersAddresses,
		F:                     threshold,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}


// --- test helpers (from test_helpers.go) ---

func AddLanes(t *testing.T, e cldf.Environment, state stateview.CCIPOnChainState, pairs []deploycciptesthelpers.SourceDestPair) cldf.Environment {
	addLanesCfg, commitOCR2Configs, execOCR2Configs := LaneConfigsForChains(t, e, state, pairs)
	var err error
	e, err = commoncs.Apply(t, e, commoncs.Configure(
		cldf.CreateLegacyChangeSet(DeployLanesChangeset),
		DeployLanesConfig{Configs: addLanesCfg},
	), commoncs.Configure(
		cldf.CreateLegacyChangeSet(SetOCR2ConfigForTestChangeset),
		OCR2Config{CommitConfigs: commitOCR2Configs, ExecConfigs: execOCR2Configs},
	))
	require.NoError(t, err)
	return e
}

func LaneConfigsForChains(t *testing.T, env cldf.Environment, state stateview.CCIPOnChainState, pairs []deploycciptesthelpers.SourceDestPair) (
	[]DeployLaneConfig,
	[]CommitOCR2ConfigParams,
	[]ExecuteOCR2ConfigParams,
) {
	addLanesCfg := make([]DeployLaneConfig, 0)
	commitOCR2Configs := make([]CommitOCR2ConfigParams, 0)
	execOCR2Configs := make([]ExecuteOCR2ConfigParams, 0)
	for _, pair := range pairs {
		dest := pair.DestChainSelector
		src := pair.SourceChainSelector
		sourceChainState := state.MustGetEVMChainState(src)
		destChainState := state.MustGetEVMChainState(dest)
		_, err := sourceChainState.LinkTokenAddress()
		require.NoError(t, err)
		require.NotNil(t, sourceChainState.RMNProxy)
		require.NotNil(t, sourceChainState.TokenAdminRegistry)
		require.NotNil(t, sourceChainState.Router)
		require.NotNil(t, sourceChainState.PriceRegistry)
		require.NotNil(t, sourceChainState.Weth9)
		_, err = destChainState.LinkTokenAddress()
		require.NoError(t, err)
		require.NotNil(t, destChainState.RMNProxy)
		require.NotNil(t, destChainState.TokenAdminRegistry)
		srcLinkTokenAddr, err := sourceChainState.LinkTokenAddress()
		require.NoError(t, err)
		addLanesCfg = append(addLanesCfg, DeployLaneConfig{
			SourceChainSelector:      src,
			DestinationChainSelector: dest,
			OnRampStaticCfg: evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
				LinkToken:          srcLinkTokenAddr,
				ChainSelector:      src,
				DestChainSelector:  dest,
				DefaultTxGasLimit:  200_000,
				MaxNopFeesJuels:    big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
				PrevOnRamp:         common.Address{},
				RmnProxy:           sourceChainState.RMNProxy.Address(),
				TokenAdminRegistry: sourceChainState.TokenAdminRegistry.Address(),
			},
			OnRampDynamicCfg: evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
				Router:                            sourceChainState.Router.Address(),
				MaxNumberOfTokensPerMsg:           5,
				DestGasOverhead:                   350_000,
				DestGasPerPayloadByte:             16,
				DestDataAvailabilityOverheadGas:   33_596,
				DestGasPerDataAvailabilityByte:    16,
				DestDataAvailabilityMultiplierBps: 6840,
				PriceRegistry:                     sourceChainState.PriceRegistry.Address(),
				MaxDataBytes:                      1e5,
				MaxPerMsgGasLimit:                 4_000_000,
				DefaultTokenFeeUSDCents:           50,
				DefaultTokenDestGasOverhead:       125_000,
			},
			OnRampFeeTokenArgs: []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
				{
					Token:                      srcLinkTokenAddr,
					NetworkFeeUSDCents:         1_00,
					GasMultiplierWeiPerEth:     1e18,
					PremiumMultiplierWeiPerEth: 9e17,
					Enabled:                    true,
				},
				{
					Token:                      sourceChainState.Weth9.Address(),
					NetworkFeeUSDCents:         1_00,
					GasMultiplierWeiPerEth:     1e18,
					PremiumMultiplierWeiPerEth: 1e18,
					Enabled:                    true,
				},
			},
			OnRampTransferTokenCfgs: []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				{
					Token:                     srcLinkTokenAddr,
					MinFeeUSDCents:            50,           // $0.5
					MaxFeeUSDCents:            1_000_000_00, // $ 1 million
					DeciBps:                   5_0,          // 5 bps
					DestGasOverhead:           350_000,
					DestBytesOverhead:         32,
					AggregateRateLimitEnabled: true,
				},
			},
			OnRampNopsAndWeight: []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
			OnRampRateLimiterCfg: evm_2_evm_onramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  linkUSDWei(100),
				Rate:      linkUSDWei(1),
			},
			OffRampRateLimiterCfg: evm_2_evm_offramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  linkUSDWei(100),
				Rate:      linkUSDWei(1),
			},
			InitialTokenPrices: []price_registry_1_2_0.InternalTokenPriceUpdate{
				{
					SourceToken: srcLinkTokenAddr,
					UsdPerToken: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(20)),
				},
				{
					SourceToken: sourceChainState.Weth9.Address(),
					UsdPerToken: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2000)),
				},
			},
			GasPriceUpdates: []price_registry_1_2_0.InternalGasPriceUpdate{
				{
					DestChainSelector: dest,
					UsdPerUnitGas:     big.NewInt(20000e9),
				},
			},
		})
		commitOCR2Configs = append(commitOCR2Configs, CommitOCR2ConfigParams{
			SourceChainSelector:      src,
			DestinationChainSelector: dest,
			OCR2ConfigParams:         DefaultOCRParams(),
			GasPriceHeartBeat:        *config.MustNewDuration(10 * time.Second),
			DAGasPriceDeviationPPB:   1,
			ExecGasPriceDeviationPPB: 1,
			TokenPriceHeartBeat:      *config.MustNewDuration(10 * time.Second),
			TokenPriceDeviationPPB:   1,
			InflightCacheExpiry:      *config.MustNewDuration(5 * time.Second),
			PriceReportingDisabled:   false,
		})
		execOCR2Configs = append(execOCR2Configs, ExecuteOCR2ConfigParams{
			DestinationChainSelector:    dest,
			SourceChainSelector:         src,
			DestOptimisticConfirmations: 1,
			BatchGasLimit:               5_000_000,
			RelativeBoostPerWaitHour:    0.07,
			InflightCacheExpiry:         *config.MustNewDuration(1 * time.Minute),
			RootSnoozeTime:              *config.MustNewDuration(1 * time.Minute),
			BatchingStrategyID:          0,
			MessageVisibilityInterval:   config.Duration{},
			ExecOnchainConfig: evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig{
				PermissionLessExecutionThresholdSeconds: uint32(24 * time.Hour.Seconds()),
				MaxDataBytes:                            1e5,
				MaxNumberOfTokensPerMsg:                 5,
			},
			OCR2ConfigParams: DefaultOCRParams(),
		})
	}
	return addLanesCfg, commitOCR2Configs, execOCR2Configs
}

func DefaultOCRParams() confighelper.PublicConfig {
	return confighelper.PublicConfig{
		DeltaProgress:                           2 * time.Second,
		DeltaResend:                             1 * time.Second,
		DeltaRound:                              1 * time.Second,
		DeltaGrace:                              500 * time.Millisecond,
		DeltaStage:                              2 * time.Second,
		RMax:                                    3,
		MaxDurationInitialization:               nil,
		MaxDurationQuery:                        50 * time.Millisecond,
		MaxDurationObservation:                  1 * time.Second,
		MaxDurationReport:                       100 * time.Millisecond,
		MaxDurationShouldAcceptFinalizedReport:  100 * time.Millisecond,
		MaxDurationShouldTransmitAcceptedReport: 100 * time.Millisecond,
	}
}

func SendRequest(
	t *testing.T,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	opts ...ccipclient.SendReqOpts,
) (*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested, error) {
	cfg := &ccipclient.CCIPSendReqConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// Set default sender if not provided
	if cfg.Sender == nil {
		cfg.Sender = e.BlockChains.EVMChains()[cfg.SourceChain].DeployerKey
	}
	t.Logf("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, cfg.Sender.From.String())
	tx, blockNum, err := deploycciptesthelpers.CCIPSendRequest(e, state, cfg)
	if err != nil {
		return nil, err
	}

	onRamp := state.MustGetEVMChainState(cfg.SourceChain).EVM2EVMOnRamp[cfg.DestChain]

	it, err := onRamp.FilterCCIPSendRequested(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}

	require.True(t, it.Next())
	t.Logf("CCIP message (id %x) sent from chain selector %d to chain selector %d tx %s seqNum %d sender %s",
		it.Event.Message.MessageId[:],
		cfg.SourceChain,
		cfg.DestChain,
		tx.Hash().String(),
		it.Event.Message.SequenceNumber,
		it.Event.Message.Sender.String(),
	)
	return it.Event, nil
}

func WaitForCommit(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	commitStore *commit_store.CommitStore,
	seqNr uint64,
) {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			minSeqNr, err := commitStore.GetExpectedNextSequenceNumber(nil)
			require.NoError(t, err)
			t.Logf("Waiting for commit for sequence number %d, current min sequence number %d", seqNr, minSeqNr)
			if minSeqNr > seqNr {
				t.Logf("Commit for sequence number %d found", seqNr)
				return
			}
		case <-timer.C:
			t.Fatalf("timed out waiting for commit for sequence number %d for commit store %s ", seqNr, commitStore.Address().String())
			return
		}
	}
}

func WaitForNoCommit(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	commitStore *commit_store.CommitStore,
	seqNr uint64,
) {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			minSeqNr, err := commitStore.GetExpectedNextSequenceNumber(nil)
			require.NoError(t, err)
			t.Logf("Waiting for commit for sequence number %d, current min sequence number %d", seqNr, minSeqNr)
			if minSeqNr > seqNr {
				t.Fatalf("Commit for sequence number %d found while it was not expected", seqNr)
				return
			}
		case <-timer.C:
			t.Logf("Successfully observed no commit for sequence number %d for commit store %s during 30s period", seqNr, commitStore.Address().String())
			return
		}
	}
}

func WaitForNoExec(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	offRamp *evm_2_evm_offramp.EVM2EVMOffRamp,
	seqNr uint64,
) {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			state, err := offRamp.GetExecutionState(nil, seqNr)
			require.NoError(t, err)
			t.Logf("Waiting for no execution for sequence number %d, current state %d", seqNr, state)
			// We expect the state to be untouched or in a failure state
			if ccipcommontypes.MessageExecutionState(state) == ccipcommontypes.ExecutionStateSuccess {
				t.Fatalf("Execution for sequence number %d found while it was not expected, current state %d", seqNr, state)
				return
			}
		case <-timer.C:
			t.Logf("Successfully observed no execution for sequence number %d for offramp %s during 30s period", seqNr, offRamp.Address().String())
			return
		}
	}
}

func WaitForExecute(
	t *testing.T,
	src cldf_evm.Chain,
	dest cldf_evm.Chain,
	offRamp *evm_2_evm_offramp.EVM2EVMOffRamp,
	seqNrs []uint64,
	blockNum uint64,
) {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t.Logf("Waiting for execute for sequence numbers %v", seqNrs)
			it, err := offRamp.FilterExecutionStateChanged(
				&bind.FilterOpts{
					Start: blockNum,
				}, seqNrs, [][32]byte{})
			require.NoError(t, err)
			for it.Next() {
				t.Logf("Execution state changed for sequence number=%d current state=%d", it.Event.SequenceNumber, it.Event.State)
				if ccipcommontypes.MessageExecutionState(it.Event.State) == ccipcommontypes.ExecutionStateSuccess {
					t.Logf("Execution for sequence number %d found", it.Event.SequenceNumber)
					return
				}
				t.Logf("Execution for sequence number %d resulted in status %d", it.Event.SequenceNumber, it.Event.State)
				t.Fail()
			}
		case <-timer.C:
			t.Fatalf("timed out waiting for execute for sequence numbers %v for offramp %s ", seqNrs, offRamp.Address().String())
			return
		}
	}
}
