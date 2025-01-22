package changeset

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

var (
	_ deployment.ChangeSet[AddDonAndSetCandidateChangesetConfig] = AddDonAndSetCandidateChangeset
	_ deployment.ChangeSet[PromoteCandidateChangesetConfig]      = PromoteCandidateChangeset
	_ deployment.ChangeSet[SetCandidateChangesetConfig]          = SetCandidateChangeset
	_ deployment.ChangeSet[RevokeCandidateChangesetConfig]       = RevokeCandidateChangeset
	_ deployment.ChangeSet[UpdateChainConfigConfig]              = UpdateChainConfigChangeset
)

type tokenInfo interface {
	Address() common.Address
	Symbol(opts *bind.CallOpts) (string, error)
	Decimals(opts *bind.CallOpts) (uint8, error)
}

func findTokenInfo(tokens []tokenInfo, address common.Address) (string, uint8, error) {
	for _, token := range tokens {
		if token.Address() == address {
			tokenSymbol, err := token.Symbol(nil)
			if err != nil {
				return "", 0, fmt.Errorf("fetch token symbol for token %s: %w", address, err)
			}
			tokenDecimals, err := token.Decimals(nil)
			if err != nil {
				return "", 0, fmt.Errorf("fetch token decimals for token %s: %w", address, err)
			}
			return tokenSymbol, tokenDecimals, nil
		}
	}
	return "", 0, fmt.Errorf("token %s not found in available tokens", address)
}

func validateCommitOffchainConfig(c *pluginconfig.CommitOffchainConfig, selector uint64, feedChainSel uint64, state CCIPOnChainState) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid commit off-chain config: %w", err)
	}
	for tokenAddr, tokenConfig := range c.TokenInfo {
		tokenUnknownAddr, err := ccipocr3.NewUnknownAddressFromHex(string(tokenAddr))
		if err != nil {
			return fmt.Errorf("invalid token address %s: %w", tokenAddr, err)
		}

		aggregatorAddr := common.HexToAddress(string(tokenConfig.AggregatorAddress))
		token := common.HexToAddress(tokenUnknownAddr.String())
		tokenInfos := make([]tokenInfo, 0)
		onchainState := state.Chains[selector]
		for _, tk := range onchainState.BurnMintTokens677 {
			tokenInfos = append(tokenInfos, tk)
		}
		for _, tk := range onchainState.ERC20Tokens {
			tokenInfos = append(tokenInfos, tk)
		}
		for _, tk := range onchainState.ERC677Tokens {
			tokenInfos = append(tokenInfos, tk)
		}
		var linkTokenInfo tokenInfo
		linkTokenInfo = onchainState.LinkToken
		if onchainState.LinkToken == nil {
			linkTokenInfo = onchainState.StaticLinkToken
		}
		tokenInfos = append(tokenInfos, linkTokenInfo)
		tokenInfos = append(tokenInfos, onchainState.Weth9)
		symbol, decimal, err := findTokenInfo(tokenInfos, token)
		if err != nil {
			return err
		}
		if decimal != tokenConfig.Decimals {
			return fmt.Errorf("token %s -address %s has %d decimals in provided token config, expected %d",
				symbol, token.String(), tokenConfig.Decimals, decimal)
		}
		feedChainState := state.Chains[feedChainSel]
		aggregatorInState := feedChainState.USDFeeds[TokenSymbol(symbol)]
		if aggregatorAddr == (common.Address{}) {
			return fmt.Errorf("token %s -address %s has no aggregator in provided token config", symbol, token.String())
		}
		if aggregatorInState == nil {
			return fmt.Errorf("token %s -address %s has no aggregator in state,"+
				" but the aggregator %s is provided in token config", symbol, token.String(), aggregatorAddr.String())
		}
		if aggregatorAddr != aggregatorInState.Address() {
			return fmt.Errorf("token %s -address %s has aggregator %s in provided token config, expected %s",
				symbol, token.String(), aggregatorAddr.String(), aggregatorInState.Address().String())
		}
	}
	return nil
}

func validateUSDCConfig(usdcConfig *pluginconfig.USDCCCTPObserverConfig, state CCIPOnChainState) error {
	for sel, token := range usdcConfig.Tokens {
		onchainState, ok := state.Chains[uint64(sel)]
		if !ok {
			return fmt.Errorf("chain %d does not exist in state but provided in USDCCCTPObserverConfig", sel)
		}
		if onchainState.USDCTokenPool == nil {
			return fmt.Errorf("chain %d does not have USDC token pool deployed", sel)
		}
		if common.HexToAddress(token.SourcePoolAddress) != onchainState.USDCTokenPool.Address() {
			return fmt.Errorf("chain %d has USDC token pool deployed at %s, "+
				"but SourcePoolAddress %s is provided in USDCCCTPObserverConfig",
				sel, onchainState.USDCTokenPool.Address().String(), token.SourcePoolAddress)
		}
	}
	return nil
}

type CCIPOCRParams struct {
	OCRParameters commontypes.OCRParameters
	// Note contains pointers to Arb feeds for prices
	CommitOffChainConfig *pluginconfig.CommitOffchainConfig
	// Note contains USDC config
	ExecuteOffChainConfig *pluginconfig.ExecuteOffchainConfig
}

func (c CCIPOCRParams) Validate(selector uint64, feedChainSel uint64, state CCIPOnChainState) error {
	if err := c.OCRParameters.Validate(); err != nil {
		return fmt.Errorf("invalid OCR parameters: %w", err)
	}
	if c.CommitOffChainConfig == nil && c.ExecuteOffChainConfig == nil {
		return errors.New("at least one of CommitOffChainConfig or ExecuteOffChainConfig must be set")
	}
	if c.CommitOffChainConfig != nil {
		if err := validateCommitOffchainConfig(c.CommitOffChainConfig, selector, feedChainSel, state); err != nil {
			return fmt.Errorf("invalid commit off-chain config: %w", err)
		}
	}
	if c.ExecuteOffChainConfig != nil {
		if err := c.ExecuteOffChainConfig.Validate(); err != nil {
			return fmt.Errorf("invalid execute off-chain config: %w", err)
		}
		for _, observerConfig := range c.ExecuteOffChainConfig.TokenDataObservers {
			switch observerConfig.Type {
			case pluginconfig.USDCCCTPHandlerType:
				if err := validateUSDCConfig(observerConfig.USDCCCTPObserverConfig, state); err != nil {
					return fmt.Errorf("invalid USDC config: %w", err)
				}
			default:
				return fmt.Errorf("unknown token observer config type: %s", observerConfig.Type)
			}
		}
	}
	return nil
}

type CCIPOCROpts func(params *CCIPOCRParams)

// WithOCRParamOverride can be used if you want to override the default OCR parameters with your custom function.
func WithOCRParamOverride(override func(params *CCIPOCRParams)) CCIPOCROpts {
	return func(params *CCIPOCRParams) {
		if override != nil {
			override(params)
		}
	}
}

// WithDefaultCommitOffChainConfig can be used to add token info to the existing commit off-chain config. If no commit off-chain config is set, it will be created with default values.
func WithDefaultCommitOffChainConfig(feedChainSel uint64, tokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo) CCIPOCROpts {
	return func(params *CCIPOCRParams) {
		if params.CommitOffChainConfig == nil {
			params.CommitOffChainConfig = &pluginconfig.CommitOffchainConfig{
				RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(internal.RemoteGasPriceBatchWriteFrequency),
				TokenPriceBatchWriteFrequency:      *config.MustNewDuration(internal.TokenPriceBatchWriteFrequency),
				TokenInfo:                          tokenInfo,
				PriceFeedChainSelector:             ccipocr3.ChainSelector(feedChainSel),
				NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
				MaxReportTransmissionCheckAttempts: 5,
				RMNEnabled:                         false,
				RMNSignaturesTimeout:               30 * time.Minute,
				MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
				SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
			}
		} else {
			if params.CommitOffChainConfig.TokenInfo == nil {
				params.CommitOffChainConfig.TokenInfo = make(map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo)
			}
			for k, v := range tokenInfo {
				params.CommitOffChainConfig.TokenInfo[k] = v
			}
		}
	}
}

// WithDefaultExecuteOffChainConfig can be used to add token data observers to the execute off-chain config. If no execute off-chain config is set, it will be created with default values.
func WithDefaultExecuteOffChainConfig(tokenDataObservers []pluginconfig.TokenDataObserverConfig) CCIPOCROpts {
	return func(params *CCIPOCRParams) {
		if params.ExecuteOffChainConfig == nil {
			params.ExecuteOffChainConfig = &pluginconfig.ExecuteOffchainConfig{
				BatchGasLimit:             internal.BatchGasLimit,
				RelativeBoostPerWaitHour:  internal.RelativeBoostPerWaitHour,
				InflightCacheExpiry:       *config.MustNewDuration(internal.InflightCacheExpiry),
				RootSnoozeTime:            *config.MustNewDuration(internal.RootSnoozeTime),
				MessageVisibilityInterval: *config.MustNewDuration(internal.FirstBlockAge),
				BatchingStrategyID:        internal.BatchingStrategyID,
				TokenDataObservers:        tokenDataObservers,
			}
		} else if tokenDataObservers != nil {
			params.ExecuteOffChainConfig.TokenDataObservers = append(params.ExecuteOffChainConfig.TokenDataObservers, tokenDataObservers...)
		}
	}
}

// DeriveCCIPOCRParams derives the default OCR parameters for a chain, with the option to override them.
func DeriveCCIPOCRParams(
	opts ...CCIPOCROpts,
) CCIPOCRParams {
	params := CCIPOCRParams{
		OCRParameters: commontypes.OCRParameters{
			DeltaProgress:                           internal.DeltaProgress,
			DeltaResend:                             internal.DeltaResend,
			DeltaInitial:                            internal.DeltaInitial,
			DeltaRound:                              internal.DeltaRound,
			DeltaGrace:                              internal.DeltaGrace,
			DeltaCertifiedCommitRequest:             internal.DeltaCertifiedCommitRequest,
			DeltaStage:                              internal.DeltaStage,
			Rmax:                                    internal.Rmax,
			MaxDurationQuery:                        internal.MaxDurationQuery,
			MaxDurationObservation:                  internal.MaxDurationObservation,
			MaxDurationShouldAcceptAttestedReport:   internal.MaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: internal.MaxDurationShouldTransmitAcceptedReport,
		},
	}
	for _, opt := range opts {
		opt(&params)
	}
	return params
}

type PromoteCandidatePluginInfo struct {
	// RemoteChainSelectors is the chain selector of the DONs that we want to promote the candidate config of.
	// Note that each (chain, ccip capability version) pair has a unique DON ID.
	RemoteChainSelectors []uint64
	PluginType           types.PluginType
}

type PromoteCandidateChangesetConfig struct {
	HomeChainSelector uint64

	PluginInfo []PromoteCandidatePluginInfo
	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *MCMSConfig
}

func (p PromoteCandidateChangesetConfig) Validate(e deployment.Environment) (map[uint64]uint32, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	if err := deployment.IsValidChainSelector(p.HomeChainSelector); err != nil {
		return nil, fmt.Errorf("home chain selector invalid: %w", err)
	}
	homeChainState, exists := state.Chains[p.HomeChainSelector]
	if !exists {
		return nil, fmt.Errorf("home chain %d does not exist", p.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), p.MCMS != nil, e.Chains[p.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CapabilityRegistry); err != nil {
		return nil, err
	}

	donIDs := make(map[uint64]uint32)
	for _, plugin := range p.PluginInfo {
		if plugin.PluginType != types.PluginTypeCCIPCommit &&
			plugin.PluginType != types.PluginTypeCCIPExec {
			return nil, errors.New("PluginType must be set to either CCIPCommit or CCIPExec")
		}
		for _, chainSelector := range plugin.RemoteChainSelectors {
			if err := deployment.IsValidChainSelector(chainSelector); err != nil {
				return nil, fmt.Errorf("don chain selector invalid: %w", err)
			}
			chainState, exists := state.Chains[chainSelector]
			if !exists {
				return nil, fmt.Errorf("chain %d does not exist", chainSelector)
			}
			if chainState.OffRamp == nil {
				// should not be possible, but a defensive check.
				return nil, errors.New("OffRamp contract does not exist")
			}

			donID, err := internal.DonIDForChain(
				state.Chains[p.HomeChainSelector].CapabilityRegistry,
				state.Chains[p.HomeChainSelector].CCIPHome,
				chainSelector,
			)
			if err != nil {
				return nil, fmt.Errorf("fetch don id for chain: %w", err)
			}
			if donID == 0 {
				return nil, fmt.Errorf("don doesn't exist in CR for chain %d", chainSelector)
			}
			// Check that candidate digest and active digest are not both zero - this is enforced onchain.
			pluginConfigs, err := state.Chains[p.HomeChainSelector].CCIPHome.GetAllConfigs(&bind.CallOpts{
				Context: e.GetContext(),
			}, donID, uint8(plugin.PluginType))
			if err != nil {
				return nil, fmt.Errorf("fetching %s configs from cciphome: %w", plugin.PluginType.String(), err)
			}

			if pluginConfigs.ActiveConfig.ConfigDigest == [32]byte{} &&
				pluginConfigs.CandidateConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("%s active and candidate config digests are both zero", plugin.PluginType.String())
			}
			donIDs[chainSelector] = donID
		}
	}
	if len(e.NodeIDs) == 0 {
		return nil, errors.New("NodeIDs must be set")
	}
	if state.Chains[p.HomeChainSelector].CCIPHome == nil {
		return nil, errors.New("CCIPHome contract does not exist")
	}
	if state.Chains[p.HomeChainSelector].CapabilityRegistry == nil {
		return nil, errors.New("CapabilityRegistry contract does not exist")
	}

	return donIDs, nil
}

// PromoteCandidateChangeset generates a proposal to call promoteCandidate on the CCIPHome through CapReg.
// Note that a DON must exist prior to being able to use this changeset effectively,
// i.e AddDonAndSetCandidateChangeset must be called first.
// This can also be used to promote a 0x0 candidate config to be the active, effectively shutting down the DON.
// At that point you can call the RemoveDON changeset to remove the DON entirely from the capability registry.
// PromoteCandidateChangeset is NOT idempotent, once candidate config is promoted to active, if it's called again,
// It might promote empty candidate config to active, which is not desired.
func PromoteCandidateChangeset(
	e deployment.Environment,
	cfg PromoteCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	donIDs, err := cfg.Validate(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("fetch node info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	homeChain := e.Chains[cfg.HomeChainSelector]

	var ops []mcms.Operation
	for _, plugin := range cfg.PluginInfo {
		for _, donID := range donIDs {
			promoteCandidateOps, err := promoteCandidateForChainOps(
				txOpts,
				homeChain,
				state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
				state.Chains[cfg.HomeChainSelector].CCIPHome,
				nodes.NonBootstraps(),
				donID,
				plugin.PluginType,
				cfg.MCMS != nil,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("generating promote candidate ops: %w", err)
			}
			ops = append(ops, promoteCandidateOps)
		}
	}

	// Disabled MCMS means that we already executed the txes, so just return early w/out the proposals.
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	prop, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch:           ops,
		}},
		"promoteCandidate",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil
}

type SetCandidatePluginInfo struct {
	// OCRConfigPerRemoteChainSelector is the chain selector of the chain where the DON will be added.
	OCRConfigPerRemoteChainSelector map[uint64]CCIPOCRParams
	PluginType                      types.PluginType
}

func (p SetCandidatePluginInfo) String() string {
	allchains := maps.Keys(p.OCRConfigPerRemoteChainSelector)
	return fmt.Sprintf("PluginType: %s, Chains: %v", p.PluginType.String(), allchains)
}

func (p SetCandidatePluginInfo) Validate(state CCIPOnChainState, homeChain uint64, feedChain uint64) error {
	if p.PluginType != types.PluginTypeCCIPCommit &&
		p.PluginType != types.PluginTypeCCIPExec {
		return errors.New("PluginType must be set to either CCIPCommit or CCIPExec")
	}
	for chainSelector, params := range p.OCRConfigPerRemoteChainSelector {
		_, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain %d does not exist in state", chainSelector)
		}
		if err := deployment.IsValidChainSelector(chainSelector); err != nil {
			return fmt.Errorf("don chain selector invalid: %w", err)
		}
		if state.Chains[chainSelector].OffRamp == nil {
			// should not be possible, but a defensive check.
			return fmt.Errorf("OffRamp contract does not exist on don chain selector %d", chainSelector)
		}
		if p.PluginType == types.PluginTypeCCIPCommit && params.CommitOffChainConfig == nil {
			return errors.New("commit off-chain config must be set")
		}
		if p.PluginType == types.PluginTypeCCIPExec && params.ExecuteOffChainConfig == nil {
			return errors.New("execute off-chain config must be set")
		}

		chainConfig, err := state.Chains[homeChain].CCIPHome.GetChainConfig(nil, chainSelector)
		if err != nil {
			return fmt.Errorf("can't get chain config for %d: %w", chainSelector, err)
		}
		// FChain should never be zero if a chain config is set in CCIPHome
		if chainConfig.FChain == 0 {
			return fmt.Errorf("chain config not set up for new chain %d", chainSelector)
		}
		if len(chainConfig.Readers) == 0 {
			return errors.New("readers must be set")
		}
		decodedChainConfig, err := chainconfig.DecodeChainConfig(chainConfig.Config)
		if err != nil {
			return fmt.Errorf("can't decode chain config: %w", err)
		}
		if err := decodedChainConfig.Validate(); err != nil {
			return fmt.Errorf("invalid chain config: %w", err)
		}
		err = params.Validate(chainSelector, feedChain, state)
		if err != nil {
			return fmt.Errorf("invalid ccip ocr params: %w", err)
		}
	}
	return nil
}

// SetCandidateConfigBase is a common base config struct for AddDonAndSetCandidateChangesetConfig and SetCandidateChangesetConfig.
// This is extracted to deduplicate most of the validation logic.
// Remaining validation logic is done in the specific config structs that inherit from this.
type SetCandidateConfigBase struct {
	HomeChainSelector uint64
	FeedChainSelector uint64

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *MCMSConfig
}

func (s SetCandidateConfigBase) Validate(e deployment.Environment, state CCIPOnChainState) error {
	if err := deployment.IsValidChainSelector(s.HomeChainSelector); err != nil {
		return fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(s.FeedChainSelector); err != nil {
		return fmt.Errorf("feed chain selector invalid: %w", err)
	}
	homeChainState, exists := state.Chains[s.HomeChainSelector]
	if !exists {
		return fmt.Errorf("home chain %d does not exist", s.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), s.MCMS != nil, e.Chains[s.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CapabilityRegistry); err != nil {
		return err
	}

	if len(e.NodeIDs) == 0 {
		return errors.New("nodeIDs must be set")
	}
	if state.Chains[s.HomeChainSelector].CCIPHome == nil {
		return errors.New("CCIPHome contract does not exist")
	}
	if state.Chains[s.HomeChainSelector].CapabilityRegistry == nil {
		return errors.New("CapabilityRegistry contract does not exist")
	}

	if e.OCRSecrets.IsEmpty() {
		return errors.New("OCR secrets must be set")
	}

	return nil
}

// AddDonAndSetCandidateChangesetConfig is a separate config struct
// because the validation is slightly different from SetCandidateChangesetConfig.
// In particular, we check to make sure we don't already have a DON for the chain.
type AddDonAndSetCandidateChangesetConfig struct {
	SetCandidateConfigBase

	// Only set one plugin at a time while you are adding the DON for the first time.
	// For subsequent SetCandidate call use SetCandidateChangeset as that fetches the already added DONID and sets the candidate.
	PluginInfo SetCandidatePluginInfo
}

func (a AddDonAndSetCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
	if err := a.SetCandidateConfigBase.Validate(e, state); err != nil {
		return err
	}

	if err := a.PluginInfo.Validate(state, a.HomeChainSelector, a.FeedChainSelector); err != nil {
		return fmt.Errorf("validate plugin info %s: %w", a.PluginInfo.String(), err)
	}
	for chainSelector := range a.PluginInfo.OCRConfigPerRemoteChainSelector {
		// check if a DON already exists for this chain
		donID, err := internal.DonIDForChain(
			state.Chains[a.HomeChainSelector].CapabilityRegistry,
			state.Chains[a.HomeChainSelector].CCIPHome,
			chainSelector,
		)
		if err != nil {
			return fmt.Errorf("fetch don id for chain: %w", err)
		}
		// if don already exists use SetCandidateChangeset instead
		if donID != 0 {
			return fmt.Errorf("don already exists in CR for chain %d, it has id %d", chainSelector, donID)
		}
	}

	return nil
}

// AddDonAndSetCandidateChangeset adds new DON for destination to home chain
// and sets the plugin config as candidateConfig for the don.
//
// This is the first step to creating a CCIP DON and must be executed before any
// other changesets (SetCandidateChangeset, PromoteCandidateChangeset)
// can be executed.
//
// Note that these operations must be done together because the createDON call
// in the capability registry calls the capability config contract, so we must
// provide suitable calldata for CCIPHome.
// AddDonAndSetCandidateChangeset is not idempotent, if AddDON is called more than once for the same chain,
// it will throw an error because the DON would already exist for that chain.
func AddDonAndSetCandidateChangeset(
	e deployment.Environment,
	cfg AddDonAndSetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	err = cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("get node info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}
	var donOps []mcms.Operation

	for chainSelector, params := range cfg.PluginInfo.OCRConfigPerRemoteChainSelector {
		newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
			e.OCRSecrets,
			state.Chains[chainSelector].OffRamp,
			e.Chains[chainSelector],
			nodes.NonBootstraps(),
			state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
			params.OCRParameters,
			params.CommitOffChainConfig,
			params.ExecuteOffChainConfig,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		latestDon, err := internal.LatestCCIPDON(state.Chains[cfg.HomeChainSelector].CapabilityRegistry)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		pluginOCR3Config, ok := newDONArgs[cfg.PluginInfo.PluginType]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("missing plugin %s in ocr3Configs",
				cfg.PluginInfo.PluginType.String())
		}

		expectedDonID := latestDon.Id + 1
		addDonOp, err := newDonWithCandidateOp(
			txOpts,
			e.Chains[cfg.HomeChainSelector],
			expectedDonID,
			pluginOCR3Config,
			state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
			nodes.NonBootstraps(),
			cfg.MCMS != nil,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		donOps = append(donOps, addDonOp)
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	prop, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch:           donOps,
		}},
		"addDON on new Chain && setCandidate for plugin "+cfg.PluginInfo.PluginType.String(),
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

// newDonWithCandidateOp sets the candidate commit config by calling setCandidate on CCIPHome contract through the AddDON call on CapReg contract
// This should be done first before calling any other UpdateDON calls
// This proposes to set up OCR3 config for the commit plugin for the DON
func newDonWithCandidateOp(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes,
	mcmsEnabled bool,
) (mcms.Operation, error) {
	encodedSetCandidateCall, err := internal.CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("pack set candidate call: %w", err)
	}

	addDonTx, err := capReg.AddDON(
		txOpts,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false, // isPublic
		false, // acceptsWorkflows
		nodes.DefaultF(),
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("could not generate add don tx w/ %s config: %w",
			types.PluginType(pluginConfig.PluginType).String(), err)
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, addDonTx, err)
		if err != nil {
			return mcms.Operation{}, fmt.Errorf("error confirming addDon call: %w", err)
		}
	}

	return mcms.Operation{
		To:    capReg.Address(),
		Data:  addDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}

type SetCandidateChangesetConfig struct {
	SetCandidateConfigBase

	PluginInfo []SetCandidatePluginInfo
}

func (s SetCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (map[uint64]uint32, error) {
	err := s.SetCandidateConfigBase.Validate(e, state)
	if err != nil {
		return nil, err
	}

	chainToDonIDs := make(map[uint64]uint32)
	for _, plugin := range s.PluginInfo {
		if err := plugin.Validate(state, s.HomeChainSelector, s.FeedChainSelector); err != nil {
			return nil, fmt.Errorf("validate plugin info %s: %w", plugin.String(), err)
		}
		for chainSelector := range plugin.OCRConfigPerRemoteChainSelector {
			donID, err := internal.DonIDForChain(
				state.Chains[s.HomeChainSelector].CapabilityRegistry,
				state.Chains[s.HomeChainSelector].CCIPHome,
				chainSelector,
			)
			if err != nil {
				return nil, fmt.Errorf("fetch don id for chain: %w", err)
			}
			// if don doesn't exist use AddDonAndSetCandidateChangeset instead
			if donID == 0 {
				return nil, fmt.Errorf("don doesn't exist in CR for chain %d", chainSelector)
			}
			chainToDonIDs[chainSelector] = donID
		}
	}
	return chainToDonIDs, nil
}

// SetCandidateChangeset generates a proposal to call setCandidate on the CCIPHome through the capability registry.
// A DON must exist in order to use this changeset effectively, i.e AddDonAndSetCandidateChangeset must be called first.
func SetCandidateChangeset(
	e deployment.Environment,
	cfg SetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chainToDonIDs, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("get node info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}
	var setCandidateOps []mcms.Operation
	pluginInfos := make([]string, 0)
	for _, plugin := range cfg.PluginInfo {
		pluginInfos = append(pluginInfos, plugin.String())
		for chainSelector, params := range plugin.OCRConfigPerRemoteChainSelector {
			newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
				e.OCRSecrets,
				state.Chains[chainSelector].OffRamp,
				e.Chains[chainSelector],
				nodes.NonBootstraps(),
				state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
				params.OCRParameters,
				params.CommitOffChainConfig,
				params.ExecuteOffChainConfig,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			config, ok := newDONArgs[plugin.PluginType]
			if !ok {
				return deployment.ChangesetOutput{}, fmt.Errorf("missing %s plugin in ocr3Configs", plugin.PluginType.String())
			}

			setCandidateMCMSOps, err := setCandidateOnExistingDon(
				txOpts,
				e.Chains[cfg.HomeChainSelector],
				state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
				nodes.NonBootstraps(),
				chainToDonIDs[chainSelector],
				config,
				cfg.MCMS != nil,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			setCandidateOps = append(setCandidateOps, setCandidateMCMSOps...)
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	prop, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch:           setCandidateOps,
		}},
		fmt.Sprintf("SetCandidate for plugin details %v", pluginInfos),
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil
}

// setCandidateOnExistingDon calls setCandidate on CCIPHome contract through the UpdateDON call on CapReg contract
// This proposes to set up OCR3 config for the provided plugin for the DON
func setCandidateOnExistingDon(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes,
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	mcmsEnabled bool,
) ([]mcms.Operation, error) {
	if donID == 0 {
		return nil, errors.New("donID is zero")
	}

	encodedSetCandidateCall, err := internal.CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return nil, fmt.Errorf("pack set candidate call: %w", err)
	}

	// set candidate call
	updateDonTx, err := capReg.UpdateDON(
		txOpts,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return nil, fmt.Errorf("update don w/ setCandidate call: %w", err)
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming updateDon call: %w", err)
		}
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming updateDon call: %w", err)
		}
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}}, nil
}

// promoteCandidateOp will create the MCMS Operation for `promoteCandidateAndRevokeActive` directed towards the capabilityRegistry
func promoteCandidateOp(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	pluginType uint8,
	mcmsEnabled bool,
) (mcms.Operation, error) {
	allConfigs, err := ccipHome.GetAllConfigs(nil, donID, pluginType)
	if err != nil {
		return mcms.Operation{}, err
	}

	encodedPromotionCall, err := internal.CCIPHomeABI.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		pluginType,
		allConfigs.CandidateConfig.ConfigDigest,
		allConfigs.ActiveConfig.ConfigDigest,
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("pack promotion call: %w", err)
	}

	updateDonTx, err := capReg.UpdateDON(
		txOpts,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("error creating updateDon op for donID(%d) and plugin type (%s): %w",
			donID, types.PluginType(pluginType).String(), err)
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return mcms.Operation{}, fmt.Errorf("error confirming updateDon call for donID(%d) and plugin type (%d): %w", donID, pluginType, err)
		}
	}

	return mcms.Operation{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}

// promoteCandidateForChainOps promotes the candidate commit and exec configs to active by calling promoteCandidateAndRevokeActive on CCIPHome through the UpdateDON call on CapReg contract
func promoteCandidateForChainOps(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	pluginType cctypes.PluginType,
	mcmsEnabled bool,
) (mcms.Operation, error) {
	if donID == 0 {
		return mcms.Operation{}, errors.New("donID is zero")
	}
	digest, err := ccipHome.GetCandidateDigest(nil, donID, uint8(pluginType))
	if err != nil {
		return mcms.Operation{}, err
	}
	fmt.Println("Promoting candidate for plugin", pluginType.String(), "with digest", digest)
	updatePluginOp, err := promoteCandidateOp(
		txOpts,
		homeChain,
		capReg,
		ccipHome,
		nodes,
		donID,
		uint8(pluginType),
		mcmsEnabled,
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("promote candidate op for plugin %s: %w", pluginType.String(), err)
	}
	return updatePluginOp, nil
}

type RevokeCandidateChangesetConfig struct {
	HomeChainSelector uint64

	// RemoteChainSelector is the chain selector whose candidate config we want to revoke.
	RemoteChainSelector uint64
	PluginType          types.PluginType

	// MCMS is optional MCMS configuration, if provided the changeset will generate an MCMS proposal.
	// If nil, the changeset will execute the commands directly using the deployer key
	// of the provided environment.
	MCMS *MCMSConfig
}

func (r RevokeCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (donID uint32, err error) {
	if err := deployment.IsValidChainSelector(r.HomeChainSelector); err != nil {
		return 0, fmt.Errorf("home chain selector invalid: %w", err)
	}
	if err := deployment.IsValidChainSelector(r.RemoteChainSelector); err != nil {
		return 0, fmt.Errorf("don chain selector invalid: %w", err)
	}
	if len(e.NodeIDs) == 0 {
		return 0, errors.New("NodeIDs must be set")
	}
	if state.Chains[r.HomeChainSelector].CCIPHome == nil {
		return 0, errors.New("CCIPHome contract does not exist")
	}
	if state.Chains[r.HomeChainSelector].CapabilityRegistry == nil {
		return 0, errors.New("CapabilityRegistry contract does not exist")
	}
	homeChainState, exists := state.Chains[r.HomeChainSelector]
	if !exists {
		return 0, fmt.Errorf("home chain %d does not exist", r.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), r.MCMS != nil, e.Chains[r.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CapabilityRegistry); err != nil {
		return 0, err
	}

	// check that the don exists for this chain
	donID, err = internal.DonIDForChain(
		state.Chains[r.HomeChainSelector].CapabilityRegistry,
		state.Chains[r.HomeChainSelector].CCIPHome,
		r.RemoteChainSelector,
	)
	if err != nil {
		return 0, fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return 0, fmt.Errorf("don doesn't exist in CR for chain %d", r.RemoteChainSelector)
	}

	// check that candidate digest is not zero - this is enforced onchain.
	candidateDigest, err := state.Chains[r.HomeChainSelector].CCIPHome.GetCandidateDigest(nil, donID, uint8(r.PluginType))
	if err != nil {
		return 0, fmt.Errorf("fetching candidate digest from cciphome: %w", err)
	}
	if candidateDigest == [32]byte{} {
		return 0, errors.New("candidate config digest is zero, can't revoke it")
	}

	return donID, nil
}

func RevokeCandidateChangeset(e deployment.Environment, cfg RevokeCandidateChangesetConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	donID, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("fetch nodes info: %w", err)
	}

	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	homeChain := e.Chains[cfg.HomeChainSelector]
	ops, err := revokeCandidateOps(
		txOpts,
		homeChain,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		state.Chains[cfg.HomeChainSelector].CCIPHome,
		nodes.NonBootstraps(),
		donID,
		uint8(cfg.PluginType),
		cfg.MCMS != nil,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("revoke candidate ops: %w", err)
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	prop, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch:           ops,
		}},
		fmt.Sprintf("revokeCandidate for don %d", cfg.RemoteChainSelector),
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{
			*prop,
		},
	}, nil
}

func revokeCandidateOps(
	txOpts *bind.TransactOpts,
	homeChain deployment.Chain,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	nodes deployment.Nodes,
	donID uint32,
	pluginType uint8,
	mcmsEnabled bool,
) ([]mcms.Operation, error) {
	if donID == 0 {
		return nil, errors.New("donID is zero")
	}

	candidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, pluginType)
	if err != nil {
		return nil, fmt.Errorf("fetching candidate digest from cciphome: %w", err)
	}

	encodedRevokeCandidateCall, err := internal.CCIPHomeABI.Pack(
		"revokeCandidate",
		donID,
		pluginType,
		candidateDigest,
	)
	if err != nil {
		return nil, fmt.Errorf("pack set candidate call: %w", err)
	}

	updateDonTx, err := capReg.UpdateDON(
		txOpts,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: internal.CCIPCapabilityID,
				Config:       encodedRevokeCandidateCall,
			},
		},
		false, // isPublic
		nodes.DefaultF(),
	)
	if err != nil {
		return nil, fmt.Errorf("update don w/ revokeCandidate call: %w", deployment.MaybeDataErr(err))
	}
	if !mcmsEnabled {
		_, err = deployment.ConfirmIfNoError(homeChain, updateDonTx, err)
		if err != nil {
			return nil, fmt.Errorf("error confirming updateDon call: %w", err)
		}
	}

	return []mcms.Operation{{
		To:    capReg.Address(),
		Data:  updateDonTx.Data(),
		Value: big.NewInt(0),
	}}, nil
}

type ChainConfig struct {
	Readers              [][32]byte
	FChain               uint8
	EncodableChainConfig chainconfig.ChainConfig
}

type UpdateChainConfigConfig struct {
	HomeChainSelector  uint64
	RemoteChainRemoves []uint64
	RemoteChainAdds    map[uint64]ChainConfig
	MCMS               *MCMSConfig
}

func (c UpdateChainConfigConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	if err := deployment.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("home chain selector invalid: %w", err)
	}
	if len(c.RemoteChainRemoves) == 0 && len(c.RemoteChainAdds) == 0 {
		return errors.New("no chain adds or removes")
	}
	homeChainState, exists := state.Chains[c.HomeChainSelector]
	if !exists {
		return fmt.Errorf("home chain %d does not exist", c.HomeChainSelector)
	}
	if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, e.Chains[c.HomeChainSelector].DeployerKey.From, homeChainState.Timelock.Address(), homeChainState.CCIPHome); err != nil {
		return err
	}
	for _, remove := range c.RemoteChainRemoves {
		if err := deployment.IsValidChainSelector(remove); err != nil {
			return fmt.Errorf("chain remove selector invalid: %w", err)
		}
		if _, ok := state.SupportedChains()[remove]; !ok {
			return fmt.Errorf("chain to remove %d is not supported", remove)
		}
	}
	for add, ccfg := range c.RemoteChainAdds {
		if err := deployment.IsValidChainSelector(add); err != nil {
			return fmt.Errorf("chain remove selector invalid: %w", err)
		}
		if _, ok := state.SupportedChains()[add]; !ok {
			return fmt.Errorf("chain to add %d is not supported", add)
		}
		if ccfg.FChain == 0 {
			return errors.New("FChain must be set")
		}
		if len(ccfg.Readers) == 0 {
			return errors.New("Readers must be set")
		}
		if err := ccfg.EncodableChainConfig.Validate(); err != nil {
			return fmt.Errorf("invalid chain config for selector %d: %w", add, err)
		}
	}
	return nil
}

func UpdateChainConfigChangeset(e deployment.Environment, cfg UpdateChainConfigConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	txOpts := e.Chains[cfg.HomeChainSelector].DeployerKey
	txOpts.Context = e.GetContext()
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}
	var adds []ccip_home.CCIPHomeChainConfigArgs
	for chain, ccfg := range cfg.RemoteChainAdds {
		encodedChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
			GasPriceDeviationPPB:    ccfg.EncodableChainConfig.GasPriceDeviationPPB,
			DAGasPriceDeviationPPB:  ccfg.EncodableChainConfig.DAGasPriceDeviationPPB,
			OptimisticConfirmations: ccfg.EncodableChainConfig.OptimisticConfirmations,
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("encoding chain config: %w", err)
		}
		chainConfig := ccip_home.CCIPHomeChainConfig{
			Readers: ccfg.Readers,
			FChain:  ccfg.FChain,
			Config:  encodedChainConfig,
		}
		existingCfg, err := state.Chains[cfg.HomeChainSelector].CCIPHome.GetChainConfig(nil, chain)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("get chain config for selector %d: %w", chain, err)
		}
		if isChainConfigEqual(existingCfg, chainConfig) {
			e.Logger.Infow("Chain config already exists, not applying again",
				"addedChain", chain,
				"chainConfig", chainConfig,
			)
			continue
		}
		adds = append(adds, ccip_home.CCIPHomeChainConfigArgs{
			ChainSelector: chain,
			ChainConfig:   chainConfig,
		})
	}

	tx, err := state.Chains[cfg.HomeChainSelector].CCIPHome.ApplyChainConfigUpdates(txOpts, cfg.RemoteChainRemoves, adds)
	if cfg.MCMS == nil {
		_, err = deployment.ConfirmIfNoError(e.Chains[cfg.HomeChainSelector], tx, err)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		e.Logger.Infof("Updated chain config on chain %d removes %v, adds %v", cfg.HomeChainSelector, cfg.RemoteChainRemoves, cfg.RemoteChainAdds)
		return deployment.ChangesetOutput{}, nil
	}

	p, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		},
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch: []mcms.Operation{
				{
					To:    state.Chains[cfg.HomeChainSelector].CCIPHome.Address(),
					Data:  tx.Data(),
					Value: big.NewInt(0),
				},
			},
		}},
		"Update chain config",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infof("Proposed chain config update on chain %d removes %v, adds %v", cfg.HomeChainSelector, cfg.RemoteChainRemoves, cfg.RemoteChainAdds)
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

func isChainConfigEqual(a, b ccip_home.CCIPHomeChainConfig) bool {
	mapReader := make(map[[32]byte]struct{})
	for i := range a.Readers {
		mapReader[a.Readers[i]] = struct{}{}
	}
	for i := range b.Readers {
		if _, ok := mapReader[b.Readers[i]]; !ok {
			return false
		}
	}
	return bytes.Equal(a.Config, b.Config) &&
		a.FChain == b.FChain
}

// ValidateCCIPHomeConfigSetUp checks that the commit and exec active and candidate configs are set up correctly
// TODO: Utilize this
func ValidateCCIPHomeConfigSetUp(
	lggr logger.Logger,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	chainSel uint64,
) error {
	// fetch DONID
	donID, err := internal.DonIDForChain(capReg, ccipHome, chainSel)
	if err != nil {
		return fmt.Errorf("fetch don id for chain: %w", err)
	}
	if donID == 0 {
		return fmt.Errorf("don id for chain (%d) does not exist", chainSel)
	}

	// final sanity checks on configs.
	commitConfigs, err := ccipHome.GetAllConfigs(&bind.CallOpts{
		//Pending: true,
	}, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get all commit configs: %w", err)
	}
	commitActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get active commit digest: %w", err)
	}
	lggr.Debugw("Fetched active commit digest", "commitActiveDigest", hex.EncodeToString(commitActiveDigest[:]))
	commitCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get commit candidate digest: %w", err)
	}
	lggr.Debugw("Fetched candidate commit digest", "commitCandidateDigest", hex.EncodeToString(commitCandidateDigest[:]))
	if commitConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return fmt.Errorf(
			"active config digest is empty for commit, expected nonempty, donID: %d, cfg: %+v, config digest from GetActiveDigest call: %x, config digest from GetCandidateDigest call: %x",
			donID, commitConfigs.ActiveConfig, commitActiveDigest, commitCandidateDigest)
	}
	if commitConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return fmt.Errorf(
			"candidate config digest is nonempty for commit, expected empty, donID: %d, cfg: %+v, config digest from GetCandidateDigest call: %x, config digest from GetActiveDigest call: %x",
			donID, commitConfigs.CandidateConfig, commitCandidateDigest, commitActiveDigest)
	}

	execConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return fmt.Errorf("get all exec configs: %w", err)
	}
	lggr.Debugw("Fetched exec configs",
		"ActiveConfig.ConfigDigest", hex.EncodeToString(execConfigs.ActiveConfig.ConfigDigest[:]),
		"CandidateConfig.ConfigDigest", hex.EncodeToString(execConfigs.CandidateConfig.ConfigDigest[:]),
	)
	if execConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return fmt.Errorf("active config digest is empty for exec, expected nonempty, cfg: %v", execConfigs.ActiveConfig)
	}
	if execConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return fmt.Errorf("candidate config digest is nonempty for exec, expected empty, cfg: %v", execConfigs.CandidateConfig)
	}
	return nil
}
