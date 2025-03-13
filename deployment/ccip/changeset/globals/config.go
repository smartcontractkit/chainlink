package globals

import (
	"time"

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
	DefaultCommitOffChainCfg = pluginconfig.CommitOffchainConfig{
		RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(30 * time.Minute),
		TokenPriceBatchWriteFrequency:      *config.MustNewDuration(30 * time.Minute),
		NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
		MaxReportTransmissionCheckAttempts: 5,
		RMNSignaturesTimeout:               6900 * time.Millisecond,
		RMNEnabled:                         true,
		MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
		SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
		TransmissionDelayMultiplier:        1 * time.Minute,
		InflightPriceCheckRetries:          10,
		MerkleRootAsyncObserverDisabled:    false,
		MerkleRootAsyncObserverSyncFreq:    4 * time.Second,
		MerkleRootAsyncObserverSyncTimeout: 12 * time.Second,
		ChainFeeAsyncObserverSyncFreq:      10 * time.Second,
		ChainFeeAsyncObserverSyncTimeout:   12 * time.Second,
	}
	DefaultExecuteOffChainCfg = pluginconfig.ExecuteOffchainConfig{
		BatchGasLimit:               6_500_000, // Building batches with 6.5m and transmit with 8m to account for overhead. Clarify with offchain
		InflightCacheExpiry:         *config.MustNewDuration(5 * time.Minute),
		RootSnoozeTime:              *config.MustNewDuration(5 * time.Minute), // does not work now
		MessageVisibilityInterval:   *config.MustNewDuration(8 * time.Hour),
		BatchingStrategyID:          0,
		TransmissionDelayMultiplier: 1 * time.Minute, // Clarify with offchain
	}
)
