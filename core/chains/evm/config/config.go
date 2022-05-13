package config

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ChainScopedOnlyConfig interface {
	evmclient.NodeConfig

	BalanceMonitorEnabled() bool
	BlockEmissionIdleWarningThreshold() time.Duration
	BlockHistoryEstimatorBatchSize() (size uint32)
	BlockHistoryEstimatorBlockDelay() uint16
	BlockHistoryEstimatorBlockHistorySize() uint16
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16
	BlockHistoryEstimatorTransactionPercentile() uint16
	ChainID() *big.Int
	EvmEIP1559DynamicFees() bool
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	EvmFinalityDepth() uint32
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint16
	EvmGasBumpWei() *big.Int
	EvmGasFeeCapDefault() *big.Int
	EvmGasLimitDefault() uint64
	EvmGasLimitMultiplier() float32
	EvmGasLimitTransfer() uint64
	EvmGasPriceDefault() *big.Int
	EvmGasTipCapDefault() *big.Int
	EvmGasTipCapMinimum() *big.Int
	EvmHeadTrackerHistoryDepth() uint32
	EvmHeadTrackerMaxBufferSize() uint32
	EvmHeadTrackerSamplingInterval() time.Duration
	EvmLogBackfillBatchSize() uint32
	EvmLogPollInterval() time.Duration
	EvmMaxGasPriceWei() *big.Int
	EvmMaxInFlightTransactions() uint32
	EvmMaxQueuedTransactions() uint64
	EvmMinGasPriceWei() *big.Int
	EvmNonceAutoSync() bool
	EvmUseForwarders() bool
	EvmRPCDefaultBatchSize() uint32
	FlagsContractAddress() string
	GasEstimatorMode() string
	ChainType() config.ChainType
	KeySpecificMaxGasPriceWei(addr gethcommon.Address) *big.Int
	LinkContractAddress() string
	MinIncomingConfirmations() uint32
	MinRequiredOutgoingConfirmations() uint64
	MinimumContractPayment() *assets.Link
	NodeNoNewHeadsThreshold() time.Duration

	// OCR2 chain specific config
	OCR2ContractConfirmations() uint16

	// OCR1 chain specific config
	OCRContractConfirmations() uint16
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRObservationGracePeriod() time.Duration
	OCRDatabaseTimeout() time.Duration

	SetEvmGasPriceDefault(value *big.Int) error
}

//go:generate mockery --name ChainScopedConfig --output ./mocks/ --case=underscore
type ChainScopedConfig interface {
	config.GeneralConfig
	ChainScopedOnlyConfig
	Validate() error
	// Both Configure() and PersistedConfig() should be accessed through ChainSet methods only.
	Configure(config evmtypes.ChainCfg)
	PersistedConfig() evmtypes.ChainCfg
}

var _ ChainScopedConfig = &chainScopedConfig{}

type chainScopedConfig struct {
	config.GeneralConfig
	logger     logger.Logger
	defaultSet chainSpecificConfigDefaultSet

	persistedCfg evmtypes.ChainCfg
	persistMu    sync.RWMutex
	orm          *chainScopedConfigORM // calls should be paired with persistedCfg updates while holding write lock

	id      *big.Int
	knownID bool // part of the default set

	onceMap   map[string]struct{}
	onceMapMu sync.RWMutex
}

func NewChainScopedConfig(chainID *big.Int, cfg evmtypes.ChainCfg, orm evmtypes.ChainConfigORM, lggr logger.Logger, gcfg config.GeneralConfig) ChainScopedConfig {
	csorm := &chainScopedConfigORM{*utils.NewBig(chainID), orm}
	defaultSet, exists := chainSpecificConfigDefaultSets[chainID.Int64()]
	if !exists {
		lggr.Warnf("Unrecognised chain %d, falling back to generic default configuration", chainID)
		defaultSet = fallbackDefaultSet
	}
	css := chainScopedConfig{
		GeneralConfig: gcfg,
		logger:        lggr,
		defaultSet:    defaultSet,
		orm:           csorm,
		id:            chainID,
		knownID:       exists,
		onceMap:       make(map[string]struct{})}
	css.Configure(cfg)
	return &css
}

func (c *chainScopedConfig) Validate() (err error) {
	return multierr.Combine(
		c.GeneralConfig.Validate(),
		c.validate(),
	)
}

func (c *chainScopedConfig) Configure(config evmtypes.ChainCfg) {
	c.persistMu.Lock()
	defer c.persistMu.Unlock()
	c.persistedCfg = config
	return
}

func (c *chainScopedConfig) PersistedConfig() evmtypes.ChainCfg {
	c.persistMu.RLock()
	defer c.persistMu.RUnlock()
	return c.persistedCfg
}

func (c *chainScopedConfig) validate() (err error) {
	ethGasBumpPercent := c.EvmGasBumpPercent()
	if uint64(ethGasBumpPercent) < core.DefaultTxPoolConfig.PriceBump {
		err = multierr.Combine(err, errors.Errorf(
			"ETH_GAS_BUMP_PERCENT of %v may not be less than Geth's default of %v",
			c.EvmGasBumpPercent(),
			core.DefaultTxPoolConfig.PriceBump,
		))
	}

	if uint32(c.EvmGasBumpTxDepth()) > c.EvmMaxInFlightTransactions() {
		err = multierr.Combine(err, errors.New("ETH_GAS_BUMP_TX_DEPTH must be less than or equal to ETH_MAX_IN_FLIGHT_TRANSACTIONS"))
	}
	if c.EvmGasTipCapDefault().Cmp(c.EvmGasTipCapMinimum()) < 0 {
		err = multierr.Combine(err, errors.Errorf("EVM_GAS_TIP_CAP_DEFAULT (%s) must be greater than or equal to EVM_GAS_TIP_CAP_MINIMUM (%s)", c.EvmGasTipCapDefault(), c.EvmGasTipCapMinimum()))
	}
	if c.EvmGasFeeCapDefault().Cmp(c.EvmGasTipCapDefault()) < 0 {
		err = multierr.Combine(err, errors.Errorf("EVM_GAS_FEE_CAP_DEFAULT (%s) must be greater than or equal to EVM_GAS_TIP_CAP_DEFAULT (%s)", c.EvmGasFeeCapDefault(), c.EvmGasTipCapDefault()))
	}
	if c.EvmGasFeeCapDefault().Cmp(c.EvmMaxGasPriceWei()) > 0 {
		err = multierr.Combine(err, errors.Errorf("EVM_GAS_FEE_CAP_DEFAULT (%s) must be less than or equal to ETH_MAX_GAS_PRICE_WEI (%s)", c.EvmGasFeeCapDefault(), c.EvmMaxGasPriceWei()))
	}
	if c.EvmMinGasPriceWei().Cmp(c.EvmGasPriceDefault()) > 0 {
		err = multierr.Combine(err, errors.New("ETH_MIN_GAS_PRICE_WEI must be less than or equal to ETH_GAS_PRICE_DEFAULT"))
	}
	if c.EvmMaxGasPriceWei().Cmp(c.EvmGasPriceDefault()) < 0 {
		err = multierr.Combine(err, errors.New("ETH_MAX_GAS_PRICE_WEI must be greater than or equal to ETH_GAS_PRICE_DEFAULT"))
	}
	if c.EvmHeadTrackerHistoryDepth() < c.EvmFinalityDepth() {
		err = multierr.Combine(err, errors.New("ETH_HEAD_TRACKER_HISTORY_DEPTH must be equal to or greater than ETH_FINALITY_DEPTH"))
	}
	if c.GasEstimatorMode() == "BlockHistory" && c.BlockHistoryEstimatorBlockHistorySize() <= 0 {
		err = multierr.Combine(err, errors.New("BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE must be greater than or equal to 1 if block history estimator is enabled"))
	}
	if c.EvmFinalityDepth() < 1 {
		err = multierr.Combine(err, errors.New("ETH_FINALITY_DEPTH must be greater than or equal to 1"))
	}
	if c.MinIncomingConfirmations() < 1 {
		err = multierr.Combine(err, errors.New("MIN_INCOMING_CONFIRMATIONS must be greater than or equal to 1"))
	}
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      c.OCRBlockchainTimeout(),
		ContractConfigConfirmations:            c.OCRContractConfirmations(),
		ContractConfigTrackerPollInterval:      c.OCRContractPollInterval(),
		ContractConfigTrackerSubscribeInterval: c.OCRContractSubscribeInterval(),
		ContractTransmitterTransmitTimeout:     c.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        c.OCRDatabaseTimeout(),
		DataSourceTimeout:                      c.OCRObservationTimeout(),
		DataSourceGracePeriod:                  c.OCRObservationGracePeriod(),
	}
	if ocrerr := ocr.SanityCheckLocalConfig(lc); ocrerr != nil {
		err = multierr.Combine(err, ocrerr)
	}

	chainType := c.ChainType()
	if !chainType.IsValid() {
		err = multierr.Combine(err, errors.Errorf("CHAIN_TYPE %q unrecognised", chainType))
	} else if c.knownID && c.defaultSet.chainType != chainType {
		err = multierr.Combine(err, errors.Errorf("CHAIN_TYPE %q cannot be used with chain ID %d", chainType, c.ChainID()))
	} else {
		switch chainType {
		case config.ChainArbitrum:
			if gasEst := c.GasEstimatorMode(); gasEst != "FixedPrice" {
				err = multierr.Combine(err, errors.Errorf("GAS_ESTIMATOR_MODE %q is not allowed with chain type %q - "+
					"must be %q", gasEst, config.ChainArbitrum, "FixedPrice"))
			}

		case config.ChainOptimism:
			gasEst := c.GasEstimatorMode()
			switch gasEst {
			case "Optimism", "Optimism2":
			default:
				err = multierr.Combine(err, errors.Errorf("GAS_ESTIMATOR_MODE %q is not allowed with chain type %q - "+
					"must be %q or %q", gasEst, config.ChainOptimism, "Optimism", "Optimism2"))
			}
		case config.ChainXDai:

		}
	}

	return err
}

func (c *chainScopedConfig) ChainID() *big.Int {
	return c.id
}

func (c *chainScopedConfig) logEnvOverrideOnce(name string, envVal interface{}) {
	k := fmt.Sprintf("env-%s", name)
	c.onceMapMu.RLock()
	if _, ok := c.onceMap[k]; ok {
		c.onceMapMu.RUnlock()
		return
	}
	c.onceMapMu.RUnlock()
	c.onceMapMu.Lock()
	defer c.onceMapMu.Unlock()
	if _, ok := c.onceMap[k]; ok {
		return
	}
	c.logger.Warnf("Global ENV var set %s=%v, overriding all other values for %s", envvar.TryName(name), envVal, name)
	c.onceMap[k] = struct{}{}
}

func (c *chainScopedConfig) logPersistedOverrideOnce(name string, pstVal interface{}) {
	k := fmt.Sprintf("pst-%s", name)
	c.onceMapMu.RLock()
	if _, ok := c.onceMap[k]; ok {
		c.onceMapMu.RUnlock()
		return
	}
	c.onceMapMu.RUnlock()
	c.onceMapMu.Lock()
	defer c.onceMapMu.Unlock()
	if _, ok := c.onceMap[k]; ok {
		return
	}
	c.logger.Infof("User-specified var set %s=%v, overriding chain-specific default value for %s", name, pstVal, name)
	c.onceMap[k] = struct{}{}
}

func (c *chainScopedConfig) logKeySpecificOverrideOnce(name string, addr gethcommon.Address, pstVal interface{}) {
	k := fmt.Sprintf("ksp-%s", name)
	c.onceMapMu.RLock()
	if _, ok := c.onceMap[k]; ok {
		c.onceMapMu.RUnlock()
		return
	}
	c.onceMapMu.RUnlock()
	c.onceMapMu.Lock()
	defer c.onceMapMu.Unlock()
	if _, ok := c.onceMap[k]; ok {
		return
	}
	c.logger.Infof("Key-specific var set %s=%v for key %s, overriding chain-specific values for %s", name, pstVal, addr.Hex(), name)
	c.onceMap[k] = struct{}{}
}

// EvmBalanceMonitorBlockDelay is the number of blocks that the balance monitor
// trails behind head. This is required e.g. for Infura because they will often
// announce a new head, then route a request to a different node which does not
// have this head yet.
func (c *chainScopedConfig) EvmBalanceMonitorBlockDelay() uint16 {
	return c.defaultSet.balanceMonitorBlockDelay
}

// EvmGasBumpThreshold is the number of blocks to wait before bumping gas again on unconfirmed transactions
// Set to 0 to disable gas bumping
func (c *chainScopedConfig) EvmGasBumpThreshold() uint64 {
	val, ok := c.GeneralConfig.GlobalEvmGasBumpThreshold()
	if ok {
		c.logEnvOverrideOnce("EvmGasBumpThreshold", val)
		return val
	}
	return c.defaultSet.gasBumpThreshold
}

// EvmGasBumpWei is the minimum fixed amount of wei by which gas is bumped on each transaction attempt
func (c *chainScopedConfig) EvmGasBumpWei() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmGasBumpWei()
	if ok {
		c.logEnvOverrideOnce("EvmGasBumpWei", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasBumpWei
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmGasBumpWei", p)
		return p.ToInt()
	}
	n := c.defaultSet.gasBumpWei
	return &n
}

// EvmMaxInFlightTransactions controls how many transactions are allowed to be
// "in-flight" i.e. broadcast but unconfirmed at any one time
// 0 value disables the limit
func (c *chainScopedConfig) EvmMaxInFlightTransactions() uint32 {
	val, ok := c.GeneralConfig.GlobalEvmMaxInFlightTransactions()
	if ok {
		c.logEnvOverrideOnce("EvmMaxInFlightTransactions", val)
		return val
	}
	return c.defaultSet.maxInFlightTransactions
}

// EvmMaxGasPriceWei is the maximum amount in Wei that a transaction will be
// bumped to before abandoning it and marking it as errored.
func (c *chainScopedConfig) EvmMaxGasPriceWei() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmMaxGasPriceWei()
	if ok {
		c.logEnvOverrideOnce("EvmMaxGasPriceWei", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmMaxGasPriceWei
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmMaxGasPriceWei", p)
		return p.ToInt()
	}
	n := c.defaultSet.maxGasPriceWei
	return &n
}

// EvmMaxQueuedTransactions is the maximum number of unbroadcast
// transactions per key that are allowed to be enqueued before jobs will start
// failing and rejecting send of any further transactions.
// 0 value disables
func (c *chainScopedConfig) EvmMaxQueuedTransactions() uint64 {
	val, ok := c.GeneralConfig.GlobalEvmMaxQueuedTransactions()
	if ok {
		c.logEnvOverrideOnce("EvmMaxGasPriceWei", val)
		return val
	}
	return c.defaultSet.maxQueuedTransactions
}

// EvmMinGasPriceWei is the minimum amount in Wei that a transaction may be priced.
// Chainlink will never send a transaction priced below this amount.
func (c *chainScopedConfig) EvmMinGasPriceWei() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmMinGasPriceWei()
	if ok {
		c.logEnvOverrideOnce("EvmMinGasPriceWei", val)
		return val
	}
	n := c.defaultSet.minGasPriceWei
	return &n
}

// EvmGasLimitDefault sets the default gas limit for outgoing transactions.
func (c *chainScopedConfig) EvmGasLimitDefault() uint64 {
	val, ok := c.GeneralConfig.GlobalEvmGasLimitDefault()
	if ok {
		c.logEnvOverrideOnce("EvmGasLimitDefault", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasLimitDefault
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmGasLimitDefault", p.Int64)
		return uint64(p.Int64)
	}
	return c.defaultSet.gasLimitDefault
}

// EvmGasLimitTransfer is the gas limit for an ordinary eth->eth transfer
func (c *chainScopedConfig) EvmGasLimitTransfer() uint64 {
	val, ok := c.GeneralConfig.GlobalEvmGasLimitTransfer()
	if ok {
		c.logEnvOverrideOnce("EvmGasLimitTransfer", val)
		return val
	}
	return c.defaultSet.gasLimitTransfer
}

// EvmGasPriceDefault is the starting gas price for every transaction
func (c *chainScopedConfig) EvmGasPriceDefault() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmGasPriceDefault()
	if ok {
		c.logEnvOverrideOnce("EvmGasPriceDefault", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasPriceDefault
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmGasPriceDefault", p)
		return p.ToInt()
	}
	n := c.defaultSet.gasPriceDefault
	return &n
}

// SetEvmGasPriceDefault saves a runtime value for the default gas price for transactions
// nil or negative value clears
func (c *chainScopedConfig) SetEvmGasPriceDefault(value *big.Int) error {
	if value == nil || value.Cmp(big.NewInt(0)) < 0 {
		c.persistMu.Lock()
		defer c.persistMu.Unlock()
		c.persistedCfg.EvmGasPriceDefault = nil
		return c.orm.clear("EvmGasPriceDefault")
	}
	min := c.EvmMinGasPriceWei()
	max := c.EvmMaxGasPriceWei()
	if value.Cmp(min) < 0 {
		return errors.Errorf("cannot set default gas price to %s, it is below the minimum allowed value of %s", value.String(), min.String())
	}
	if value.Cmp(max) > 0 {
		return errors.Errorf("cannot set default gas price to %s, it is above the maximum allowed value of %s", value.String(), max.String())
	}
	c.persistMu.Lock()
	defer c.persistMu.Unlock()
	c.persistedCfg.EvmGasPriceDefault = utils.NewBig(value)
	return c.orm.storeString("EvmGasPriceDefault", value.String())
}

// EvmFinalityDepth is the number of blocks after which an ethereum transaction is considered "final"
// BlocksConsideredFinal determines how deeply we look back to ensure that transactions are confirmed onto the longest chain
// There is not a large performance penalty to setting this relatively high (on the order of hundreds)
// It is practically limited by the number of heads we store in the database and should be less than this with a comfortable margin.
// If a transaction is mined in a block more than this many blocks ago, and is reorged out, we will NOT retransmit this transaction and undefined behaviour can occur including gaps in the nonce sequence that require manual intervention to fix.
// Therefore this number represents a number of blocks we consider large enough that no re-org this deep will ever feasibly happen.
//
// Special cases:
// ETH_FINALITY_DEPTH=0 would imply that transactions can be final even before they were mined into a block. This is not supported.
// ETH_FINALITY_DEPTH=1 implies that transactions are final after we see them in one block.
//
// Examples:
//
// Transaction sending:
// A transaction is sent at block height 42
//
// ETH_FINALITY_DEPTH is set to 5
// A re-org occurs at height 44 starting at block 41, transaction is marked for rebroadcast
// A re-org occurs at height 46 starting at block 41, transaction is marked for rebroadcast
// A re-org occurs at height 47 starting at block 41, transaction is NOT marked for rebroadcast
func (c *chainScopedConfig) EvmFinalityDepth() uint32 {
	val, ok := c.GeneralConfig.GlobalEvmFinalityDepth()
	if ok {
		c.logEnvOverrideOnce("EvmFinalityDepth", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmFinalityDepth
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmFinalityDepth", p.Int64)
		return uint32(p.Int64)
	}
	return c.defaultSet.finalityDepth
}

// EvmHeadTrackerHistoryDepth tracks the top N block numbers to keep in the `heads` database table.
// Note that this can easily result in MORE than N records since in the case of re-orgs we keep multiple heads for a particular block height.
// This number should be at least as large as `EvmFinalityDepth`.
// There may be a small performance penalty to setting this to something very large (10,000+)
func (c *chainScopedConfig) EvmHeadTrackerHistoryDepth() uint32 {
	val, ok := c.GeneralConfig.GlobalEvmHeadTrackerHistoryDepth()
	if ok {
		c.logEnvOverrideOnce("EvmHeadTrackerHistoryDepth", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmHeadTrackerHistoryDepth
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmHeadTrackerHistoryDepth", p.Int64)
		return uint32(p.Int64)
	}
	return c.defaultSet.headTrackerHistoryDepth
}

// EvmHeadTrackerSamplingInterval is the interval between sampled head callbacks
// to services that are only interested in the latest head every some time
// Setting it to a zero duration disables sampling (every head will be delivered)
func (c *chainScopedConfig) EvmHeadTrackerSamplingInterval() time.Duration {
	val, ok := c.GeneralConfig.GlobalEvmHeadTrackerSamplingInterval()
	if ok {
		c.logEnvOverrideOnce("EvmHeadTrackerSamplingInterval", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmHeadTrackerSamplingInterval
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmHeadTrackerSamplingInterval", p.Duration())
		return p.Duration()
	}
	return c.defaultSet.headTrackerSamplingInterval
}

// BlockEmissionIdleWarningThreshold is the duration of time since last received head
// to print a warning log message indicating not receiving heads
func (c *chainScopedConfig) BlockEmissionIdleWarningThreshold() time.Duration {
	val, ok := c.GeneralConfig.GlobalBlockEmissionIdleWarningThreshold()
	if ok {
		c.logEnvOverrideOnce("BlockEmissionIdleWarningThreshold", val)
		return val
	}
	return c.defaultSet.blockEmissionIdleWarningThreshold
}

// EthTxResendAfterThreshold controls how long the ethResender will wait before
// re-sending the latest eth_tx_attempt. This is designed a as a fallback to
// protect against the eth nodes dropping txes (it has been anecdotally
// observed to happen), networking issues or txes being ejected from the
// mempool.
// See eth_resender.go for more details
func (c *chainScopedConfig) EthTxResendAfterThreshold() time.Duration {
	val, ok := c.GeneralConfig.GlobalEthTxResendAfterThreshold()
	if ok {
		c.logEnvOverrideOnce("EthTxResendAfterThreshold", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EthTxResendAfterThreshold
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EthTxResendAfterThreshold", p.Duration())
		return p.Duration()
	}
	return c.defaultSet.ethTxResendAfterThreshold
}

// BlockHistoryEstimatorBatchSize sets the maximum number of blocks to fetch in one batch in the block history estimator
// If the env var GAS_UPDATER_BATCH_SIZE is set to 0, it defaults to ETH_RPC_DEFAULT_BATCH_SIZE
func (c *chainScopedConfig) BlockHistoryEstimatorBatchSize() (size uint32) {
	val, ok := c.GeneralConfig.GlobalBlockHistoryEstimatorBatchSize()
	if ok {
		c.logEnvOverrideOnce("BlockHistoryEstimatorBatchSize", val)
		size = val
	} else {
		valLegacy, set := lookupEnv(c, "GAS_UPDATER_BATCH_SIZE", parse.Uint32)
		if set {
			c.logEnvOverrideOnce("GAS_UPDATER_BATCH_SIZE", valLegacy)
			c.logger.Error("GAS_UPDATER_BATCH_SIZE is deprecated, please use BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE instead (or simply remove to use the default)")
			size = valLegacy
		} else {
			size = c.defaultSet.blockHistoryEstimatorBatchSize
		}
	}

	if size > 0 {
		return size
	}
	return c.EvmRPCDefaultBatchSize()
}

// BlockHistoryEstimatorBlockDelay is the number of blocks that the block history estimator trails behind head.
// E.g. if this is set to 3, and we receive block 10, block history estimator will
// fetch block 7.
// CAUTION: You might be tempted to set this to 0 to use the latest possible
// block, but it is possible to receive a head BEFORE that block is actually
// available from the connected node via RPC. In this case you will get false
// "zero" blocks that are missing transactions.
func (c *chainScopedConfig) BlockHistoryEstimatorBlockDelay() uint16 {
	val, ok := c.GeneralConfig.GlobalBlockHistoryEstimatorBlockDelay()
	if ok {
		c.logEnvOverrideOnce("BlockHistoryEstimatorBlockDelay", val)
		return val
	}
	valLegacy, set := lookupEnv(c, "GAS_UPDATER_BLOCK_DELAY", parse.Uint16)

	if set {
		c.logEnvOverrideOnce("GAS_UPDATER_BLOCK_DELAY", valLegacy)
		c.logger.Error("GAS_UPDATER_BLOCK_DELAY is deprecated, please use BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY instead (or simply remove to use the default)")
		return valLegacy
	}
	c.persistMu.RLock()
	p := c.persistedCfg.BlockHistoryEstimatorBlockDelay
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("BlockHistoryEstimatorBlockDelay", p.Int64)
		return uint16(p.Int64)
	}
	return c.defaultSet.blockHistoryEstimatorBlockDelay
}

// BlockHistoryEstimatorBlockHistorySize is the number of past blocks to keep in memory to
// use as a basis for calculating a percentile gas price
func (c *chainScopedConfig) BlockHistoryEstimatorBlockHistorySize() uint16 {
	val, ok := c.GeneralConfig.GlobalBlockHistoryEstimatorBlockHistorySize()
	if ok {
		c.logEnvOverrideOnce("BlockHistoryEstimatorBlockHistorySize", val)
		return val
	}
	valLegacy, set := lookupEnv(c, "GAS_UPDATER_BLOCK_HISTORY_SIZE", parse.Uint16)
	if set {
		c.logEnvOverrideOnce("GAS_UPDATER_BLOCK_HISTORY_SIZE", valLegacy)
		c.logger.Error("GAS_UPDATER_BLOCK_HISTORY_SIZE is deprecated, please use BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE instead (or simply remove to use the default)")
		return valLegacy
	}
	c.persistMu.RLock()
	p := c.persistedCfg.BlockHistoryEstimatorBlockHistorySize
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("BlockHistoryEstimatorBlockHistorySize", p.Int64)
		return uint16(p.Int64)
	}
	return c.defaultSet.blockHistoryEstimatorBlockHistorySize
}

func (c *chainScopedConfig) BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16 {
	val, ok := c.GeneralConfig.GlobalBlockHistoryEstimatorEIP1559FeeCapBufferBlocks()
	if ok {
		c.logEnvOverrideOnce("BlockHistoryEstimatorBlockHistorySize", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("BlockHistoryEstimatorBlockHistorySize", p.Int64)
		return uint16(p.Int64)
	}
	if c.defaultSet.blockHistoryEstimatorEIP1559FeeCapBufferBlocks != nil {
		return *c.defaultSet.blockHistoryEstimatorEIP1559FeeCapBufferBlocks
	}
	// Default is the gas bump threshold + 1 block
	return uint16(c.EvmGasBumpThreshold() + 1)
}

// BlockHistoryEstimatorTransactionPercentile is the percentile gas price to choose. E.g.
// if the past transaction history contains four transactions with gas prices:
// [100, 200, 300, 400], picking 25 for this number will give a value of 200
func (c *chainScopedConfig) BlockHistoryEstimatorTransactionPercentile() uint16 {
	val, ok := c.GeneralConfig.GlobalBlockHistoryEstimatorTransactionPercentile()
	if ok {
		c.logEnvOverrideOnce("BlockHistoryEstimatorTransactionPercentile", val)
		return val
	}
	valLegacy, set := lookupEnv(c, "GAS_UPDATER_TRANSACTION_PERCENTILE", parse.Uint16)
	if set {
		c.logEnvOverrideOnce("GAS_UPDATER_TRANSACTION_PERCENTILE", valLegacy)
		c.logger.Error("GAS_UPDATER_TRANSACTION_PERCENTILE is deprecated, please use BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE instead (or simply remove to use the default)")
		return valLegacy
	}
	return c.defaultSet.blockHistoryEstimatorTransactionPercentile
}

// GasEstimatorMode controls what type of gas estimator is used
func (c *chainScopedConfig) GasEstimatorMode() string {
	val, ok := c.GeneralConfig.GlobalGasEstimatorMode()
	if ok {
		c.logEnvOverrideOnce("GasEstimatorMode", val)
		return val
	}
	enabled, set := lookupEnv(c, "GAS_UPDATER_ENABLED", parse.Bool)
	if set {
		c.logEnvOverrideOnce("GAS_UPDATER_ENABLED", enabled)
		if enabled.(bool) {
			c.logger.Error("GAS_UPDATER_ENABLED has been deprecated, to enable the block history estimator, please use GAS_ESTIMATOR_MODE=BlockHistory instead (or simply remove to use the default)")
			return "BlockHistory"
		}
		c.logger.Error("GAS_UPDATER_ENABLED has been deprecated, to disable the block history estimator, please use GAS_ESTIMATOR_MODE=FixedPrice instead (or simply remove to use the default)")
		return "FixedPrice"
	}
	c.persistMu.RLock()
	p := c.persistedCfg.GasEstimatorMode
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("GasEstimatorMode", p.String)
		return p.String
	}
	return c.defaultSet.gasEstimatorMode
}

func (c *chainScopedConfig) KeySpecificMaxGasPriceWei(addr gethcommon.Address) *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmMaxGasPriceWei()
	if ok {
		c.logEnvOverrideOnce("EvmMaxGasPriceWei", val)
		return val
	}
	c.persistMu.RLock()
	keySpecific := c.persistedCfg.KeySpecific[addr.Hex()].EvmMaxGasPriceWei
	c.persistMu.RUnlock()
	if keySpecific != nil && !keySpecific.Equal(utils.NewBigI(0)) {
		c.logKeySpecificOverrideOnce("EvmMaxGasPriceWei", addr, keySpecific)
		return keySpecific.ToInt()
	}
	return c.EvmMaxGasPriceWei()
}

func (c *chainScopedConfig) ChainType() config.ChainType {
	val, ok := c.GeneralConfig.GlobalChainType()
	if ok {
		c.logEnvOverrideOnce("ChainType", val)
		return config.ChainType(val)
	}
	c.persistMu.RLock()
	p := c.persistedCfg.ChainType
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("ChainType", p.String)
		return config.ChainType(p.String)
	}
	return c.defaultSet.chainType
}

// LinkContractAddress represents the address of the official LINK token
// contract on the current Chain
func (c *chainScopedConfig) LinkContractAddress() string {
	val, ok := c.GeneralConfig.GlobalLinkContractAddress()
	if ok {
		c.logEnvOverrideOnce("LinkContractAddress", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.LinkContractAddress
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("LinkContractAddress", p.String)
		return p.String
	}
	return c.defaultSet.linkContractAddress
}

// MinIncomingConfirmations represents the minimum number of block
// confirmations that need to be recorded since a job run started before a task
// can proceed.
// MIN_INCOMING_CONFIRMATIONS=1 would kick off a job after seeing the transaction in a block
// MIN_INCOMING_CONFIRMATIONS=0 would kick off a job even before the transaction is mined, which is not supported
func (c *chainScopedConfig) MinIncomingConfirmations() uint32 {
	val, ok := c.GeneralConfig.GlobalMinIncomingConfirmations()
	if ok {
		c.logEnvOverrideOnce("MinIncomingConfirmations", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.MinIncomingConfirmations
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("MinIncomingConfirmations", p.Int64)
		return uint32(p.Int64)
	}
	return c.defaultSet.minIncomingConfirmations
}

// MinRequiredOutgoingConfirmations represents the default minimum number of block
// confirmations that need to be recorded on an outgoing ethtx task before the run can move onto the next task.
// This can be overridden on a per-task basis by setting the `MinRequiredOutgoingConfirmations` parameter.
// MIN_OUTGOING_CONFIRMATIONS=1 considers a transaction as "done" once it has been mined into one block
// MIN_OUTGOING_CONFIRMATIONS=0 would consider a transaction as "done" even before it has been mined
func (c *chainScopedConfig) MinRequiredOutgoingConfirmations() uint64 {
	val, ok := c.GeneralConfig.GlobalMinRequiredOutgoingConfirmations()
	if ok {
		c.logEnvOverrideOnce("MinRequiredOutgoingConfirmations", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.MinRequiredOutgoingConfirmations
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("MinRequiredOutgoingConfirmations", p.Int64)
		return uint64(p.Int64)
	}
	return c.defaultSet.minRequiredOutgoingConfirmations
}

// MinimumContractPayment represents the minimum amount of LINK that must be
// supplied for a contract to be considered.
func (c *chainScopedConfig) MinimumContractPayment() *assets.Link {
	val, ok := c.GeneralConfig.GlobalMinimumContractPayment()
	if ok {
		c.logEnvOverrideOnce("MinimumContractPayment", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.MinimumContractPayment
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("MinimumContractPayment", p)
		return p
	}
	return c.defaultSet.minimumContractPayment
}

// EvmGasBumpTxDepth is the number of transactions to gas bump starting from oldest.
// Set to 0 for no limit (i.e. bump all)
func (c *chainScopedConfig) EvmGasBumpTxDepth() uint16 {
	val, ok := c.GeneralConfig.GlobalEvmGasBumpTxDepth()
	if ok {
		c.logEnvOverrideOnce("EvmGasBumpTxDepth", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasBumpTxDepth
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmGasBumpTxDepth", p.Int64)
		return uint16(p.Int64)
	}
	return c.defaultSet.gasBumpTxDepth
}

// EvmGasBumpPercent is the minimum percentage by which gas is bumped on each transaction attempt
// Change with care since values below geth's default will fail with "underpriced replacement transaction"
func (c *chainScopedConfig) EvmGasBumpPercent() uint16 {
	val, ok := c.GeneralConfig.GlobalEvmGasBumpPercent()
	if ok {
		c.logEnvOverrideOnce("EvmGasBumpPercent", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasBumpPercent
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmGasBumpPercent", p.Int64)
		return uint16(p.Int64)
	}
	return c.defaultSet.gasBumpPercent
}

// EvmNonceAutoSync enables/disables running the NonceSyncer on application start
func (c *chainScopedConfig) EvmNonceAutoSync() bool {
	val, ok := c.GeneralConfig.GlobalEvmNonceAutoSync()
	if ok {
		c.logEnvOverrideOnce("EvmNonceAutoSync", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmNonceAutoSync
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmNonceAutoSync", p.Bool)
		return p.Bool
	}
	return c.defaultSet.nonceAutoSync
}

// EvmUseForwarders enables/disables sending transactions through forwarder contracts
func (c *chainScopedConfig) EvmUseForwarders() bool {
	val, ok := c.GeneralConfig.GlobalEvmUseForwarders()
	if ok {
		c.logEnvOverrideOnce("EvmUseForwarders", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmUseForwarders
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmUseForwarders", p.Bool)
		return p.Bool
	}
	return c.defaultSet.useForwarders
}

// EvmGasLimitMultiplier is a factor by which a transaction's GasLimit is
// multiplied before transmission. So if the value is 1.1, and the GasLimit for
// a transaction is 10, 10% will be added before transmission.
//
// This factor is always applied, so includes Optimism L2 transactions which
// uses a default gas limit of 1 and is also applied to EvmGasLimitDefault.
func (c *chainScopedConfig) EvmGasLimitMultiplier() float32 {
	val, ok := c.GeneralConfig.GlobalEvmGasLimitMultiplier()
	if ok {
		c.logEnvOverrideOnce("EvmGasLimitMultiplier", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasLimitMultiplier
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmGasLimitMultiplier", p.Float64)
		return float32(p.Float64)
	}
	return c.defaultSet.gasLimitMultiplier
}

// EvmHeadTrackerMaxBufferSize is the maximum number of heads that may be
// buffered in front of the head tracker before older heads start to be
// dropped. You may think of it as something like the maximum permittable "lag"
// for the head tracker before we start dropping heads to keep up.
func (c *chainScopedConfig) EvmHeadTrackerMaxBufferSize() uint32 {
	val, ok := c.GeneralConfig.GlobalEvmHeadTrackerMaxBufferSize()
	if ok {
		c.logEnvOverrideOnce("EvmHeadTrackerMaxBufferSize", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmHeadTrackerMaxBufferSize
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmHeadTrackerMaxBufferSize", p.Int64)
		return uint32(p.Int64)
	}
	return c.defaultSet.headTrackerMaxBufferSize
}

// EthTxReaperInterval controls how often the eth tx reaper should run
func (c *chainScopedConfig) EthTxReaperInterval() time.Duration {
	val, ok := c.GeneralConfig.GlobalEthTxReaperInterval()
	if ok {
		c.logEnvOverrideOnce("EthTxReaperInterval", val)
		return val
	}
	return c.defaultSet.ethTxReaperInterval
}

// EthTxReaperThreshold represents how long any confirmed/fatally_errored eth_txes will hang around in the database.
// If the eth_tx is confirmed but still below ETH_FINALITY_DEPTH it will not be deleted even if it was created at a time older than this value.
// EXAMPLE
// With:
// EthTxReaperThreshold=1h
// EvmFinalityDepth=50
//
// Current head is 142, any eth_tx confirmed in block 91 or below will be reaped as long as its created_at was more than EthTxReaperThreshold ago
// Set to 0 to disable eth_tx reaping
func (c *chainScopedConfig) EthTxReaperThreshold() time.Duration {
	val, ok := c.GeneralConfig.GlobalEthTxReaperThreshold()
	if ok {
		c.logEnvOverrideOnce("EthTxReaperThreshold", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EthTxReaperThreshold
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EthTxReaperThreshold", p.Duration())
		return p.Duration()
	}
	return c.defaultSet.ethTxReaperThreshold
}

// EvmLogPollInterval how fast we poll for new logs.
func (c *chainScopedConfig) EvmLogPollInterval() time.Duration {
	val, ok := c.GeneralConfig.GlobalEvmLogPollInterval()
	if ok {
		c.logEnvOverrideOnce("EvmLogPollInterval", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmLogPollInterval
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmLogPollInterval", *p)
		return p.Duration()
	}
	return c.defaultSet.logPollInterval
}

// EvmLogBackfillBatchSize sets the batch size for calling FilterLogs when we backfill missing logs
func (c *chainScopedConfig) EvmLogBackfillBatchSize() uint32 {
	val, ok := c.GeneralConfig.GlobalEvmLogBackfillBatchSize()
	if ok {
		c.logEnvOverrideOnce("EvmLogBackfillBatchSize", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmLogBackfillBatchSize
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmLogBackfillBatchSize", p.Int64)
		return uint32(p.Int64)
	}
	return c.defaultSet.logBackfillBatchSize
}

// EvmRPCDefaultBatchSize controls the number of receipts fetched in each
// request in the EthConfirmer
func (c *chainScopedConfig) EvmRPCDefaultBatchSize() uint32 {
	val, ok := c.GeneralConfig.GlobalEvmRPCDefaultBatchSize()
	if ok {
		c.logEnvOverrideOnce("EvmRPCDefaultBatchSize", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmRPCDefaultBatchSize
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmRPCDefaultBatchSize", p.Int64)
		return uint32(p.Int64)
	}
	return c.defaultSet.rpcDefaultBatchSize
}

// FlagsContractAddress represents the Flags contract address
func (c *chainScopedConfig) FlagsContractAddress() string {
	val, ok := c.GeneralConfig.GlobalFlagsContractAddress()
	if ok {
		c.logEnvOverrideOnce("FlagsContractAddress", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.FlagsContractAddress
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("FlagsContractAddress", p.String)
		return p.String
	}
	return c.defaultSet.flagsContractAddress
}

// BalanceMonitorEnabled enables the balance monitor
func (c *chainScopedConfig) BalanceMonitorEnabled() bool {
	val, ok := c.GeneralConfig.GlobalBalanceMonitorEnabled()
	if ok {
		c.logEnvOverrideOnce("BalanceMonitorEnabled", val)
		return val
	}
	return c.defaultSet.balanceMonitorEnabled
}

// EvmEIP1559DynamicFees will send transactions with the 0x2 dynamic fee EIP-2718
// type and gas fields when enabled
func (c *chainScopedConfig) EvmEIP1559DynamicFees() bool {
	val, ok := c.GeneralConfig.GlobalEvmEIP1559DynamicFees()
	if ok {
		c.logEnvOverrideOnce("EvmEIP1559DynamicFees", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmEIP1559DynamicFees
	c.persistMu.RUnlock()
	if p.Valid {
		c.logPersistedOverrideOnce("EvmEIP1559DynamicFees", p.Bool)
		return p.Bool
	}
	return c.defaultSet.eip1559DynamicFees
}

// EvmGasFeeCapDefault is the fixed amount to set the fee cap on DynamicFee transactions
func (c *chainScopedConfig) EvmGasFeeCapDefault() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmGasFeeCapDefault()
	if ok {
		c.logEnvOverrideOnce("EvmGasFeeCapDefault", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasFeeCapDefault
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmGasFeeCapDefault", p)
		return p.ToInt()
	}
	return &c.defaultSet.gasFeeCapDefault
}

// EvmGasTipCapDefault is the default value to use for the gas tip on DynamicFee transactions
// This is analogous to EthGasPriceDefault except the base fee is excluded
func (c *chainScopedConfig) EvmGasTipCapDefault() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmGasTipCapDefault()
	if ok {
		c.logEnvOverrideOnce("EvmGasTipCapDefault", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasTipCapDefault
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmGasTipCapDefault", p)
		return p.ToInt()
	}
	return &c.defaultSet.gasTipCapDefault
}

// EvmGasTipCapMinimum is the minimum allowed value to use for the gas tip on DynamicFee transactions
// This is analogous to EthMinGasPriceWei except the base fee is excluded
func (c *chainScopedConfig) EvmGasTipCapMinimum() *big.Int {
	val, ok := c.GeneralConfig.GlobalEvmGasTipCapMinimum()
	if ok {
		c.logEnvOverrideOnce("EvmGasTipCapMinimum", val)
		return val
	}
	c.persistMu.RLock()
	p := c.persistedCfg.EvmGasTipCapMinimum
	c.persistMu.RUnlock()
	if p != nil {
		c.logPersistedOverrideOnce("EvmGasTipCapMinimum", p)
		return p.ToInt()
	}
	return &c.defaultSet.gasTipCapMinimum
}

// NodeNoNewHeadsThreshold controls how long to wait after receiving no new
// heads before marking the node as out-of-sync
// Set to zero to disable out-of-sync checking
func (c *chainScopedConfig) NodeNoNewHeadsThreshold() time.Duration {
	val, ok := c.GeneralConfig.GlobalNodeNoNewHeadsThreshold()
	if ok {
		c.logEnvOverrideOnce("NodeNoNewHeadsThreshold", val)
		return val
	}
	return c.defaultSet.nodeDeadAfterNoNewHeadersThreshold
}

// NodePollFailureThreshold indicates how many consecutive polls must fail in
// order to mark a node as unreachable.
// Set to zero to disable poll checking.
func (c *chainScopedConfig) NodePollFailureThreshold() uint32 {
	val, ok := c.GeneralConfig.GlobalNodePollFailureThreshold()
	if ok {
		c.logEnvOverrideOnce("NodePollFailureThreshold", val)
		return val
	}
	return c.defaultSet.nodePollFailureThreshold
}

// NodePollInterval controls how often to poll the node to check for liveness.
// Set to zero to disable poll checking.
func (c *chainScopedConfig) NodePollInterval() time.Duration {
	val, ok := c.GeneralConfig.GlobalNodePollInterval()
	if ok {
		c.logEnvOverrideOnce("NodePollInterval", val)
		return val
	}
	return c.defaultSet.nodePollInterval
}

func lookupEnv[T any](c *chainScopedConfig, k string, parse func(string) (T, error)) (t T, ok bool) {
	s, ok := os.LookupEnv(k)
	if !ok {
		return
	}
	val, err := parse(s)
	if err == nil {
		return val, true
	}
	c.logger.Errorw(fmt.Sprintf("Invalid value provided for %s, falling back to default.", s),
		"value", s, "key", k, "error", err)
	return
}
