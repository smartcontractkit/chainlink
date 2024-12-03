package changeset

import (
	"fmt"
	"os"
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

// ConfigureNewChains enables new chains as destination(s) for CCIP
// It performs the following steps per chain:
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

type CCIPOCRParams struct {
	OCRParameters types.OCRParameters
	// Note contains pointers to Arb feeds for prices
	CommitOffChainConfig pluginconfig.CommitOffchainConfig
	// Note ontains USDC config
	ExecuteOffChainConfig pluginconfig.ExecuteOffchainConfig
}

func (c CCIPOCRParams) Validate() error {
	if err := c.OCRParameters.Validate(); err != nil {
		return fmt.Errorf("invalid OCR parameters: %w", err)
	}
	if err := c.CommitOffChainConfig.Validate(); err != nil {
		return fmt.Errorf("invalid commit off-chain config: %w", err)
	}
	if err := c.ExecuteOffChainConfig.Validate(); err != nil {
		return fmt.Errorf("invalid execute off-chain config: %w", err)
	}
	return nil
}

type NewChainsConfig struct {
	// Common to all chains
	HomeChainSel uint64
	FeedChainSel uint64
	OCRSecrets   deployment.OCRSecrets
	// Per chain config
	ChainConfigByChain map[uint64]CCIPOCRParams
}

func (c NewChainsConfig) Chains() []uint64 {
	chains := make([]uint64, 0, len(c.ChainConfigByChain))
	for chain := range c.ChainConfigByChain {
		chains = append(chains, chain)
	}
	return chains
}

func (c NewChainsConfig) Validate() error {
	if err := deployment.IsValidChainSelector(c.HomeChainSel); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSel, err)
	}
	if err := deployment.IsValidChainSelector(c.FeedChainSel); err != nil {
		return fmt.Errorf("invalid feed chain selector: %d - %w", c.FeedChainSel, err)
	}
	if c.OCRSecrets.IsEmpty() {
		return fmt.Errorf("no OCR secrets provided")
	}
	// Validate chain config
	for chain, cfg := range c.ChainConfigByChain {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid OCR params for chain %d: %w", chain, err)
		}
		if cfg.CommitOffChainConfig.PriceFeedChainSelector != ccipocr3.ChainSelector(c.FeedChainSel) {
			return fmt.Errorf("chain %d has invalid feed chain selector", chain)
		}
	}
	return nil
}

// DefaultOCRParams returns the default OCR parameters for a chain,
// except for a few values which must be parameterized (passed as arguments).
func DefaultOCRParams(
	feedChainSel uint64,
	tokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	tokenDataObservers []pluginconfig.TokenDataObserverConfig,
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
			TokenDataObservers:        tokenDataObservers,
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
