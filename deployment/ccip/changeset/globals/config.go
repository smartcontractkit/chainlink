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
	PermissionLessExecutionThreshold  = 8 * time.Hour
	RemoteGasPriceBatchWriteFrequency = 30 * time.Minute
	TokenPriceBatchWriteFrequency     = 30 * time.Minute
	BatchGasLimit                     = 6_500_000
	InflightCacheExpiry               = 1 * time.Minute
	RootSnoozeTime                    = 5 * time.Minute
	BatchingStrategyID                = 0
	GasPriceDeviationPPB              = 1000
	DAGasPriceDeviationPPB            = 0
	OptimisticConfirmations           = 1
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
		RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(30 * time.Minute),
		TokenPriceBatchWriteFrequency:      *config.MustNewDuration(30 * time.Minute),
		NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
		MaxReportTransmissionCheckAttempts: 10,
		RMNSignaturesTimeout:               6900 * time.Millisecond,
		RMNEnabled:                         true,
		MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
		SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
		// TransmissionDelayMultiplier for non-ETH (i.e, typically fast) chains should be pretty aggressive.
		// e.g assuming a 2s blocktime, 15 seconds is ~8 blocks.
		TransmissionDelayMultiplier:        15 * time.Second,
		InflightPriceCheckRetries:          10,
		MerkleRootAsyncObserverDisabled:    false,
		MerkleRootAsyncObserverSyncFreq:    4 * time.Second,
		MerkleRootAsyncObserverSyncTimeout: 12 * time.Second,
		ChainFeeAsyncObserverDisabled:      false,
		ChainFeeAsyncObserverSyncFreq:      10 * time.Second,
		ChainFeeAsyncObserverSyncTimeout:   12 * time.Second,
		TokenPriceAsyncObserverDisabled:    false,
		TokenPriceAsyncObserverSyncFreq:    *config.MustNewDuration(10 * time.Second),
		TokenPriceAsyncObserverSyncTimeout: *config.MustNewDuration(12 * time.Second),

		// Remaining fields cannot be statically set:
		// PriceFeedChainSelector: , // Must be configured in CLD
		// TokenInfo: , // Must be configured in CLD
	}

	// DefaultExecuteOffChainCfg represents the default offchain configuration for the Execute plugin
	// on _most_ chains. This should be used as a base for all chains, with overrides only where necessary.
	// Notable overrides are for Ethereum, which has a slower block time.
	DefaultExecuteOffChainCfg = pluginconfig.ExecuteOffchainConfig{
		BatchGasLimit:             6_500_000, // Building batches with 6.5m and transmit with 8m to account for overhead. Clarify with offchain
		InflightCacheExpiry:       *config.MustNewDuration(5 * time.Minute),
		RootSnoozeTime:            *config.MustNewDuration(5 * time.Minute), // does not work now
		MessageVisibilityInterval: *config.MustNewDuration(8 * time.Hour),
		BatchingStrategyID:        0,
		// TransmissionDelayMultiplier for non-ETH (i.e, typically fast) chains should be pretty aggressive.
		TransmissionDelayMultiplier: 25 * time.Second,
	}

	// CommitOffChainCfgForETH represents the offchain configuration for the Commit plugin on Ethereum.
	// This is necessary because Ethereum has a slower block time than other EVM chains.
	// Hence the transmission delay multiplier is increased to account for this.
	CommitOffChainCfgForETH = withCommitOverrides(
		DefaultCommitOffChainCfg,
		pluginconfig.CommitOffchainConfig{
			// 45 seconds is ~4 blocks on Ethereum
			TransmissionDelayMultiplier: 45 * time.Second,
		},
	)

	// ExecuteOffChainCfgForETH represents the offchain configuration for the Execute plugin on Ethereum.
	// This is necessary because Ethereum has a slower block time than other EVM chains.
	// Hence the transmission delay multiplier is increased to account for this.
	ExecuteOffChainCfgForETH = withExecOverrides(
		DefaultExecuteOffChainCfg,
		pluginconfig.ExecuteOffchainConfig{
			// 45 seconds is ~4 blocks on Ethereum
			TransmissionDelayMultiplier: 45 * time.Second,
		},
	)
)

func withCommitOverrides(base pluginconfig.CommitOffchainConfig, overrides pluginconfig.CommitOffchainConfig) pluginconfig.CommitOffchainConfig {
	outcome := base
	if err := mergo.Merge(&outcome, overrides, mergo.WithOverride); err != nil {
		panic(fmt.Sprintf("error while building a commit plugin offchain config %v", err))
	}
	return outcome
}

func withExecOverrides(base pluginconfig.ExecuteOffchainConfig, overrides pluginconfig.ExecuteOffchainConfig) pluginconfig.ExecuteOffchainConfig {
	outcome := base
	if err := mergo.Merge(&outcome, overrides, mergo.WithOverride); err != nil {
		panic(fmt.Sprintf("error while building an execute plugin offchain config %v", err))
	}
	return outcome
}
