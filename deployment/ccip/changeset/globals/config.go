package globals

import (
	"fmt"
	"time"

	"dario.cat/mergo"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
)

type ConfigType string

const (
	ConfigTypeActive    ConfigType = "active"
	ConfigTypeCandidate ConfigType = "candidate"
	// ========= Changeset Defaults =========
	PermissionLessExecutionThreshold  = 1 * time.Hour
	RemoteGasPriceBatchWriteFrequency = 20 * time.Minute
	TokenPriceBatchWriteFrequency     = 2 * time.Hour
	// Building batches with 6.5m and transmit with 8m to account for overhead.
	BatchGasLimit               = 6_500_000
	InflightCacheExpiry         = 1 * time.Minute
	RootSnoozeTime              = 5 * time.Minute
	BatchingStrategyID          = 0
	OptimisticConfirmations     = 1
	TransmissionDelayMultiplier = 15 * time.Second
	MaxCommitReportsToFetch     = 250
	// ======================================

	// ========= Onchain consts =========
	// CCIPLockOrBurnV1RetBytes Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES
	// Reference: https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/ccip/libraries/Pool.sol#L17
	CCIPLockOrBurnV1RetBytes = 32
	// ======================================
)

var (
	// DefaultCommitOffChainCfg represents the default offchain configuration for the Commit plugin
	// on _most_ chains. This should be used as a base for all chains, with overrides only where necessary.
	// Notable overrides are for Ethereum, which has a slower block time.
	DefaultCommitOffChainCfg = pluginconfig.CommitOffchainConfig{
		RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(RemoteGasPriceBatchWriteFrequency),
		TokenPriceBatchWriteFrequency:      *config.MustNewDuration(TokenPriceBatchWriteFrequency),
		NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
		MaxReportTransmissionCheckAttempts: 10,
		RMNSignaturesTimeout:               6900 * time.Millisecond,
		RMNEnabled:                         true,
		MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
		SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
		// TransmissionDelayMultiplier for non-ETH (i.e, typically fast) chains should be pretty aggressive.
		// e.g assuming a 2s blocktime, 15 seconds is ~8 blocks.
		TransmissionDelayMultiplier:        TransmissionDelayMultiplier,
		InflightPriceCheckRetries:          10,
		MerkleRootAsyncObserverDisabled:    false,
		MerkleRootAsyncObserverSyncFreq:    4 * time.Second,
		MerkleRootAsyncObserverSyncTimeout: 12 * time.Second,

		// Disabling the chainfee + tokenprice async observers because the low cache TTL + low timeout
		// is currently not a viable combo.
		// Super aggressive frequency and timeout causes rpc timeouts more frequently.
		ChainFeeAsyncObserverDisabled: true,
		// TODO: revisit
		// ChainFeeAsyncObserverSyncFreq:      1*time.Second + 500*time.Millisecond,
		// ChainFeeAsyncObserverSyncTimeout:   1 * time.Second,
		TokenPriceAsyncObserverDisabled: true,
		// TODO: revisit
		// TokenPriceAsyncObserverSyncFreq:    *config.MustNewDuration(1*time.Second + 500*time.Millisecond),
		// TokenPriceAsyncObserverSyncTimeout: *config.MustNewDuration(1 * time.Second),

		// Remaining fields cannot be statically set:
		// PriceFeedChainSelector: , // Must be configured in CLD
		// TokenInfo: , // Must be configured in CLD
	}

	// DefaultExecuteOffChainCfg represents the default offchain configuration for the Execute plugin
	// on _most_ chains. This should be used as a base for all chains, with overrides only where necessary.
	// Notable overrides are for Ethereum, which has a slower block time.
	DefaultExecuteOffChainCfg = pluginconfig.ExecuteOffchainConfig{
		BatchGasLimit:               BatchGasLimit,
		InflightCacheExpiry:         *config.MustNewDuration(InflightCacheExpiry),
		RootSnoozeTime:              *config.MustNewDuration(RootSnoozeTime),
		MessageVisibilityInterval:   *config.MustNewDuration(8 * time.Hour),
		BatchingStrategyID:          BatchingStrategyID,
		TransmissionDelayMultiplier: TransmissionDelayMultiplier,
		MaxReportMessages:           0,
		MaxSingleChainReports:       0,
		MaxCommitReportsToFetch:     MaxCommitReportsToFetch,
		// Remaining fields cannot be statically set:
		// TokenDataObservers: , // Must be configured in CLD
	}

	DefaultCommitOffChainCfgForEth = withCommitOffchainConfigOverrides(
		DefaultCommitOffChainCfg,
		pluginconfig.CommitOffchainConfig{
			RemoteGasPriceBatchWriteFrequency: *config.MustNewDuration(2 * time.Hour),
			TokenPriceBatchWriteFrequency:     *config.MustNewDuration(12 * time.Hour),
		},
	)
)

func withCommitOffchainConfigOverrides(base pluginconfig.CommitOffchainConfig, overrides pluginconfig.CommitOffchainConfig) pluginconfig.CommitOffchainConfig {
	outcome := base
	if err := mergo.Merge(&outcome, overrides, mergo.WithOverride); err != nil {
		panic(fmt.Sprintf("error while building an OCR config %v", err))
	}
	return outcome
}
