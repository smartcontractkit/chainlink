package changeset

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[NewChainsConfig] = ConfigureNewChains

// ConfigureNewChains enables new chains as destination for CCIP
// It performs the following steps:
// - AddChainConfig + AddDON (candidate->primary promotion i.e. init) on the home chain
// - SetOCR3Config on the remote chain
// ConfigureNewChains assumes that the home chain is already enabled and all CCIP contracts are already deployed.
func ConfigureNewChains(env deployment.Environment, c NewChainsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid NewChainsConfig: %w", err)
	}
	err := configureChain(env, c)
	if err != nil {
		env.Logger.Errorw("Failed to configure chain", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    nil,
	}, nil
}

type USDCConfig struct {
	EnabledChains []uint64
	USDCAttestationConfig
	CCTPTokenConfig map[ccipocr3.ChainSelector]pluginconfig.USDCCCTPTokenConfig
}

func (cfg USDCConfig) EnabledChainMap() map[uint64]bool {
	m := make(map[uint64]bool)
	for _, chain := range cfg.EnabledChains {
		m[chain] = true
	}
	return m
}

func (cfg USDCConfig) ToTokenDataObserverConfig() []pluginconfig.TokenDataObserverConfig {
	return []pluginconfig.TokenDataObserverConfig{{
		Type:    pluginconfig.USDCCCTPHandlerType,
		Version: "1.0",
		USDCCCTPObserverConfig: &pluginconfig.USDCCCTPObserverConfig{
			Tokens:                 cfg.CCTPTokenConfig,
			AttestationAPI:         cfg.API,
			AttestationAPITimeout:  cfg.APITimeout,
			AttestationAPIInterval: cfg.APIInterval,
		},
	}}
}

type USDCAttestationConfig struct {
	API         string
	APITimeout  *config.Duration
	APIInterval *config.Duration
}

type CCIPOCRParams struct {
	OCRParameters         types.OCRParameters
	CommitOffChainConfig  pluginconfig.CommitOffchainConfig
	ExecuteOffChainConfig pluginconfig.ExecuteOffchainConfig
}

func (p CCIPOCRParams) Validate() error {
	if err := p.OCRParameters.Validate(); err != nil {
		return fmt.Errorf("invalid OCR parameters: %w", err)
	}
	if err := p.CommitOffChainConfig.Validate(); err != nil {
		return fmt.Errorf("invalid commit off-chain config: %w", err)
	}
	if err := p.ExecuteOffChainConfig.Validate(); err != nil {
		return fmt.Errorf("invalid execute off-chain config: %w", err)
	}
	return nil
}

type NewChainsConfig struct {
	HomeChainSel   uint64
	FeedChainSel   uint64
	ChainsToDeploy []uint64
	TokenConfig    TokenConfig
	USDCConfig     USDCConfig
	// For setting OCR configuration
	OCRSecrets deployment.OCRSecrets
	OCRParams  map[uint64]CCIPOCRParams
}

func (c NewChainsConfig) Validate() error {
	if err := deployment.IsValidChainSelector(c.HomeChainSel); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSel, err)
	}
	if err := deployment.IsValidChainSelector(c.FeedChainSel); err != nil {
		return fmt.Errorf("invalid feed chain selector: %d - %w", c.FeedChainSel, err)
	}
	mapChainsToDeploy := make(map[uint64]bool)
	for _, cs := range c.ChainsToDeploy {
		mapChainsToDeploy[cs] = true
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	for token := range c.TokenConfig.TokenSymbolToInfo {
		if err := c.TokenConfig.TokenSymbolToInfo[token].Validate(); err != nil {
			return fmt.Errorf("invalid token config for token %s: %w", token, err)
		}
	}
	if c.OCRSecrets.IsEmpty() {
		return fmt.Errorf("no OCR secrets provided")
	}
	usdcEnabledChainMap := c.USDCConfig.EnabledChainMap()
	for chain := range usdcEnabledChainMap {
		if _, exists := mapChainsToDeploy[chain]; !exists {
			return fmt.Errorf("chain %d is not in chains to deploy", chain)
		}
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	for chain := range c.USDCConfig.CCTPTokenConfig {
		if _, exists := mapChainsToDeploy[uint64(chain)]; !exists {
			return fmt.Errorf("chain %d is not in chains to deploy", chain)
		}
		if _, exists := usdcEnabledChainMap[uint64(chain)]; !exists {
			return fmt.Errorf("chain %d is not enabled in USDC config", chain)
		}
		if err := deployment.IsValidChainSelector(uint64(chain)); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	// Validate OCR params
	var ocrChains []uint64
	for chain, ocrParams := range c.OCRParams {
		ocrChains = append(ocrChains, chain)
		if _, exists := mapChainsToDeploy[chain]; !exists {
			return fmt.Errorf("chain %d is not in chains to deploy", chain)
		}
		if err := ocrParams.Validate(); err != nil {
			return fmt.Errorf("invalid OCR params for chain %d: %w", chain, err)
		}
	}
	sort.Slice(ocrChains, func(i, j int) bool { return ocrChains[i] < ocrChains[j] })
	sort.Slice(c.ChainsToDeploy, func(i, j int) bool { return c.ChainsToDeploy[i] < c.ChainsToDeploy[j] })
	if !slices.Equal(ocrChains, c.ChainsToDeploy) {
		return fmt.Errorf("mismatch in given OCR params and chains to deploy")
	}
	return nil
}

func DefaultOCRParams(
	feedChainSel uint64,
	tokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
) CCIPOCRParams {
	return CCIPOCRParams{
		OCRParameters: types.OCRParameters{
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
		ExecuteOffChainConfig: pluginconfig.ExecuteOffchainConfig{
			BatchGasLimit:             internal.BatchGasLimit,
			RelativeBoostPerWaitHour:  internal.RelativeBoostPerWaitHour,
			InflightCacheExpiry:       *config.MustNewDuration(internal.InflightCacheExpiry),
			RootSnoozeTime:            *config.MustNewDuration(internal.RootSnoozeTime),
			MessageVisibilityInterval: *config.MustNewDuration(internal.FirstBlockAge),
			BatchingStrategyID:        internal.BatchingStrategyID,
		},
		CommitOffChainConfig: pluginconfig.CommitOffchainConfig{
			RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(internal.RemoteGasPriceBatchWriteFrequency),
			TokenPriceBatchWriteFrequency:      *config.MustNewDuration(internal.TokenPriceBatchWriteFrequency),
			TokenInfo:                          tokenInfo,
			PriceFeedChainSelector:             ccipocr3.ChainSelector(feedChainSel),
			NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
			MaxReportTransmissionCheckAttempts: 5,
			RMNEnabled:                         os.Getenv("ENABLE_RMN") == "true", // only enabled in manual test
			RMNSignaturesTimeout:               30 * time.Minute,
			MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
			SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
		},
	}
}
