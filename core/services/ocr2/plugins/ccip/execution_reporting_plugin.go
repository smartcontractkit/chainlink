package ccip

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/custom_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/contractutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	// exec Report should make sure to cap returned payload to this limit
	MaxExecutionReportLength = 250_000

	// MaxDataLenPerBatch limits the total length of msg data that can be in a batch.
	MaxDataLenPerBatch = 60_000
)

var (
	_ types.ReportingPluginFactory = &ExecutionReportingPluginFactory{}
	_ types.ReportingPlugin        = &ExecutionReportingPlugin{}
)

type ExecutionPluginConfig struct {
	lggr                     logger.Logger
	sourceLP, destLP         logpoller.LogPoller
	sourceReader             ccipdata.Reader
	destReader               ccipdata.Reader
	onRamp                   evm_2_evm_onramp.EVM2EVMOnRampInterface
	offRamp                  evm_2_evm_offramp.EVM2EVMOffRampInterface
	commitStore              commit_store.CommitStoreInterface
	sourcePriceRegistry      price_registry.PriceRegistryInterface
	sourceWrappedNativeToken common.Address
	destClient               evmclient.Client
	sourceClient             evmclient.Client
	destGasEstimator         gas.EvmFeeEstimator
	leafHasher               hashlib.LeafHasherInterface[[32]byte]
	tokenDataProviders       map[common.Address]tokendata.Reader
}

type ExecutionReportingPlugin struct {
	config                 ExecutionPluginConfig
	F                      int
	lggr                   logger.Logger
	inflightReports        *inflightExecReportsContainer
	snoozedRoots           cache.SnoozedRoots
	destPriceRegistry      price_registry.PriceRegistryInterface
	destWrappedNative      common.Address
	onchainConfig          ccipconfig.ExecOnchainConfig
	offchainConfig         ccipconfig.ExecOffchainConfig
	cachedSourceFeeTokens  cache.AutoSync[[]common.Address]
	cachedDestTokens       cache.AutoSync[cache.CachedTokens]
	customTokenPoolFactory func(ctx context.Context, poolAddress common.Address, bind bind.ContractBackend) (custom_token_pool.CustomTokenPoolInterface, error)
}

type ExecutionReportingPluginFactory struct {
	config ExecutionPluginConfig

	// We keep track of the registered filters
	sourceChainFilters []logpoller.Filter
	destChainFilters   []logpoller.Filter
	filtersMu          *sync.Mutex
}

func NewExecutionReportingPluginFactory(config ExecutionPluginConfig) *ExecutionReportingPluginFactory {
	return &ExecutionReportingPluginFactory{
		config:    config,
		filtersMu: &sync.Mutex{},
	}
}

func (rf *ExecutionReportingPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	onchainConfig, err := abihelpers.DecodeAbiStruct[ccipconfig.ExecOnchainConfig](config.OnchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	offchainConfig, err := ccipconfig.DecodeOffchainConfig[ccipconfig.ExecOffchainConfig](config.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	priceRegistry, err := observability.NewObservedPriceRegistry(onchainConfig.PriceRegistry, ExecPluginLabel, rf.config.destClient)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	destRouter, err := router.NewRouter(onchainConfig.Router, rf.config.destClient)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	destWrappedNative, err := destRouter.GetWrappedNative(nil)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	if err = rf.UpdateLogPollerFilters(onchainConfig.PriceRegistry); err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	cachedSourceFeeTokens := cache.NewCachedFeeTokens(rf.config.sourceLP, rf.config.sourcePriceRegistry, int64(offchainConfig.SourceFinalityDepth))
	cachedDestTokens := cache.NewCachedSupportedTokens(rf.config.destLP, rf.config.offRamp, priceRegistry, int64(offchainConfig.DestOptimisticConfirmations))
	rf.config.lggr.Infow("Starting exec plugin",
		"offchainConfig", offchainConfig,
		"onchainConfig", onchainConfig)

	return &ExecutionReportingPlugin{
			config:                rf.config,
			F:                     config.F,
			lggr:                  rf.config.lggr.Named("ExecutionReportingPlugin"),
			snoozedRoots:          cache.NewSnoozedRoots(onchainConfig.PermissionLessExecutionThresholdDuration(), offchainConfig.RootSnoozeTime.Duration()),
			inflightReports:       newInflightExecReportsContainer(offchainConfig.InflightCacheExpiry.Duration()),
			destPriceRegistry:     priceRegistry,
			destWrappedNative:     destWrappedNative,
			onchainConfig:         onchainConfig,
			offchainConfig:        offchainConfig,
			cachedDestTokens:      cachedDestTokens,
			cachedSourceFeeTokens: cachedSourceFeeTokens,
			customTokenPoolFactory: func(ctx context.Context, poolAddress common.Address, contractBackend bind.ContractBackend) (custom_token_pool.CustomTokenPoolInterface, error) {
				return custom_token_pool.NewCustomTokenPool(poolAddress, contractBackend)
			},
		}, types.ReportingPluginInfo{
			Name: "CCIPExecution",
			// Setting this to false saves on calldata since OffRamp doesn't require agreement between NOPs
			// (OffRamp is only able to execute committed messages).
			UniqueReports: false,
			Limits: types.ReportingPluginLimits{
				MaxObservationLength: MaxObservationLength,
				MaxReportLength:      MaxExecutionReportLength,
			},
		}, nil
}

func (r *ExecutionReportingPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ExecutionReportingPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	lggr := r.lggr.Named("ExecutionObservation")
	if contractutil.IsCommitStoreDownNow(ctx, lggr, r.config.commitStore) {
		return nil, ErrCommitStoreIsDown
	}
	// Expire any inflight reports.
	r.inflightReports.expire(lggr)
	inFlight := r.inflightReports.getAll()

	observationBuildStart := time.Now()

	executableObservations, err := r.getExecutableObservations(ctx, lggr, timestamp, inFlight)
	measureObservationBuildDuration(timestamp, time.Since(observationBuildStart))
	if err != nil {
		return nil, err
	}
	// cap observations which fits MaxObservationLength (after serialized)
	capped := sort.Search(len(executableObservations), func(i int) bool {
		var encoded []byte
		encoded, err = NewExecutionObservation(executableObservations[:i+1]).Marshal()
		if err != nil {
			// false makes Search keep looking to the right, always including any "erroring" ObservedMessage and allowing us to detect in the bottom
			return false
		}
		return len(encoded) > MaxObservationLength
	})
	if err != nil {
		return nil, err
	}
	executableObservations = executableObservations[:capped]
	lggr.Infow("Observation", "executableMessages", executableObservations)
	// Note can be empty
	return NewExecutionObservation(executableObservations).Marshal()
}

// UpdateLogPollerFilters updates the log poller filters for the source and destination chains.
// pass zeroAddress if dstPriceRegistry is unknown, filters with zero address are omitted.
func (rf *ExecutionReportingPluginFactory) UpdateLogPollerFilters(destPriceRegistry common.Address, qopts ...pg.QOpt) error {
	rf.filtersMu.Lock()
	defer rf.filtersMu.Unlock()

	// source chain filters
	sourceFiltersBefore, sourceFiltersNow := rf.sourceChainFilters, getExecutionPluginSourceLpChainFilters(
		rf.config.onRamp.Address(),
		rf.config.sourcePriceRegistry.Address(),
		rf.config.tokenDataProviders,
	)
	created, deleted := logpollerutil.FiltersDiff(sourceFiltersBefore, sourceFiltersNow)
	if err := logpollerutil.UnregisterLpFilters(rf.config.sourceLP, deleted, qopts...); err != nil {
		return err
	}
	if err := logpollerutil.RegisterLpFilters(rf.config.sourceLP, created, qopts...); err != nil {
		return err
	}
	rf.sourceChainFilters = sourceFiltersNow

	// destination chain filters
	destFiltersBefore, destFiltersNow := rf.destChainFilters, getExecutionPluginDestLpChainFilters(rf.config.commitStore.Address(), rf.config.offRamp.Address(), destPriceRegistry)
	created, deleted = logpollerutil.FiltersDiff(destFiltersBefore, destFiltersNow)
	if err := logpollerutil.UnregisterLpFilters(rf.config.destLP, deleted, qopts...); err != nil {
		return err
	}
	if err := logpollerutil.RegisterLpFilters(rf.config.destLP, created, qopts...); err != nil {
		return err
	}
	rf.destChainFilters = destFiltersNow

	return nil
}

func (r *ExecutionReportingPlugin) getExecutableObservations(ctx context.Context, lggr logger.Logger, timestamp types.ReportTimestamp, inflight []InflightInternalExecutionReport) ([]ObservedMessage, error) {
	unexpiredReports, err := getUnexpiredCommitReports(
		ctx,
		r.config.destReader,
		r.config.commitStore,
		r.onchainConfig.PermissionLessExecutionThresholdDuration(),
	)
	if err != nil {
		return nil, err
	}
	lggr.Infow("Unexpired roots", "n", len(unexpiredReports))
	if len(unexpiredReports) == 0 {
		return []ObservedMessage{}, nil
	}

	// This could result in slightly different values on each call as
	// the function returns the allowed amount at the time of the last block.
	// Since this will only increase over time, the highest observed value will
	// always be the lower bound of what would be available on chain
	// since we already account for inflight txs.
	getAllowedTokenAmount := cache.LazyFetch(func() (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
		return r.config.offRamp.CurrentRateLimiterState(&bind.CallOpts{Context: ctx})
	})
	sourceToDestTokens, supportedDestTokens, err := r.sourceDestinationTokens(ctx)
	if err != nil {
		return nil, err
	}
	getSourceTokensPrices := cache.LazyFetch(func() (map[common.Address]*big.Int, error) {
		sourceFeeTokens, err1 := r.cachedSourceFeeTokens.Get(ctx)
		if err1 != nil {
			return nil, err1
		}
		return getTokensPrices(ctx, sourceFeeTokens, r.config.sourcePriceRegistry, []common.Address{r.config.sourceWrappedNativeToken})
	})
	getDestTokensPrices := cache.LazyFetch(func() (map[common.Address]*big.Int, error) {
		dstTokens, err1 := r.cachedDestTokens.Get(ctx)
		if err1 != nil {
			return nil, err1
		}
		return getTokensPrices(ctx, dstTokens.FeeTokens, r.destPriceRegistry, append(supportedDestTokens, r.destWrappedNative))
	})
	getDestGasPrice := cache.LazyFetch(func() (*big.Int, error) {
		return r.estimateDestinationGasPrice(ctx)
	})

	lggr.Infow("Processing unexpired reports", "n", len(unexpiredReports))
	measureNumberOfReportsProcessed(timestamp, len(unexpiredReports))
	reportIterationStart := time.Now()
	defer func() {
		measureReportsIterationDuration(timestamp, time.Since(reportIterationStart))
	}()

	unexpiredReportsWithSendReqs, err := r.getReportsWithSendRequests(ctx, unexpiredReports)
	if err != nil {
		return nil, err
	}

	getDestPoolRateLimits := cache.LazyFetch(func() (map[common.Address]*big.Int, error) {
		return r.destPoolRateLimits(ctx, unexpiredReportsWithSendReqs, sourceToDestTokens)
	})

	for _, rep := range unexpiredReportsWithSendReqs {
		if ctx.Err() != nil {
			lggr.Warn("Processing of roots killed by context")
			break
		}

		merkleRoot := rep.commitReport.MerkleRoot

		rootLggr := lggr.With("root", hexutil.Encode(merkleRoot[:]),
			"minSeqNr", rep.commitReport.Interval.Min,
			"maxSeqNr", rep.commitReport.Interval.Max,
		)

		if r.snoozedRoots.IsSnoozed(merkleRoot) {
			rootLggr.Debug("Skipping snoozed root")
			continue
		}

		if err := rep.validate(); err != nil {
			rootLggr.Errorw("Skipping invalid report", "err", err)
			continue
		}

		// If all messages are already executed and finalized, snooze the root for
		// config.PermissionLessExecutionThresholdSeconds so it will never be considered again.
		if allMsgsExecutedAndFinalized := rep.allRequestsAreExecutedAndFinalized(); allMsgsExecutedAndFinalized {
			rootLggr.Infof("Snoozing root %s forever since there are no executable txs anymore", hex.EncodeToString(merkleRoot[:]))
			r.snoozedRoots.MarkAsExecuted(merkleRoot)
			incSkippedRequests(reasonAllExecuted)
			continue
		}

		blessed, err := r.config.commitStore.IsBlessed(&bind.CallOpts{Context: ctx}, merkleRoot)
		if err != nil {
			return nil, err
		}
		if !blessed {
			rootLggr.Infow("Report is accepted but not blessed")
			incSkippedRequests(reasonNotBlessed)
			continue
		}

		allowedTokenAmountValue, err := getAllowedTokenAmount()
		if err != nil {
			return nil, err
		}
		sourceTokensPricesValue, err := getSourceTokensPrices()
		if err != nil {
			return nil, fmt.Errorf("get source token prices: %w", err)
		}

		destTokensPricesValue, err := getDestTokensPrices()
		if err != nil {
			return nil, fmt.Errorf("get dest token prices: %w", err)
		}

		destPoolRateLimits, err := getDestPoolRateLimits()
		if err != nil {
			return nil, fmt.Errorf("get dest pool rate limits: %w", err)
		}

		buildBatchDuration := time.Now()
		batch := r.buildBatch(
			ctx,
			rootLggr,
			rep,
			inflight,
			allowedTokenAmountValue.Tokens,
			sourceTokensPricesValue,
			destTokensPricesValue,
			getDestGasPrice,
			sourceToDestTokens,
			destPoolRateLimits)
		measureBatchBuildDuration(timestamp, time.Since(buildBatchDuration))
		if len(batch) != 0 {
			return batch, nil
		}
		r.snoozedRoots.Snooze(merkleRoot)
	}
	return []ObservedMessage{}, nil
}

// destPoolRateLimits returns a map that consists of the rate limits of each destination tokens of the provided reports.
// If a token is missing from the returned map it either means that token was not found or token pool is disabled for this token.
func (r *ExecutionReportingPlugin) destPoolRateLimits(ctx context.Context, commitReports []commitReportWithSendRequests, sourceToDestToken map[common.Address]common.Address) (map[common.Address]*big.Int, error) {
	dstTokens := make(map[common.Address]struct{}) // todo: replace with a set or uniqueSlice data structure
	for _, msg := range commitReports {
		for _, req := range msg.sendRequestsWithMeta {
			for _, tk := range req.TokenAmounts {
				if dstToken, exists := sourceToDestToken[tk.Token]; exists {
					dstTokens[dstToken] = struct{}{}
					continue
				}
				r.lggr.Warnw("token not found on destination chain", "sourceToken", tk)
			}
		}
	}

	res := make(map[common.Address]*big.Int, len(dstTokens))

	for dstToken := range dstTokens {
		poolAddress, err := r.config.offRamp.GetPoolByDestToken(&bind.CallOpts{Context: ctx}, dstToken)
		if err != nil {
			return nil, fmt.Errorf("get pool by dest token (%s): %w", dstToken, err)
		}

		tokenPool, err := r.customTokenPoolFactory(ctx, poolAddress, r.config.destClient)
		if err != nil {
			return nil, fmt.Errorf("new custom dest token pool %s: %w", poolAddress, err)
		}

		rateLimiterState, err := tokenPool.CurrentOffRampRateLimiterState(&bind.CallOpts{Context: ctx}, r.config.offRamp.Address())
		if err != nil {
			return nil, fmt.Errorf("get rate off ramp rate limiter state: %w", err)
		}

		if rateLimiterState.IsEnabled {
			res[dstToken] = rateLimiterState.Tokens
		}
	}

	return res, nil
}

func (r *ExecutionReportingPlugin) estimateDestinationGasPrice(ctx context.Context) (*big.Int, error) {
	destGasPriceWei, _, err := r.config.destGasEstimator.GetFee(ctx, nil, 0, assets.NewWei(big.NewInt(int64(r.offchainConfig.MaxGasPrice))))
	if err != nil {
		return nil, errors.Wrap(err, "could not estimate destination gas price")
	}
	destGasPrice := destGasPriceWei.Legacy.ToInt()
	if destGasPriceWei.DynamicFeeCap != nil {
		destGasPrice = destGasPriceWei.DynamicFeeCap.ToInt()
	}
	return destGasPrice, nil
}

func (r *ExecutionReportingPlugin) sourceDestinationTokens(ctx context.Context) (map[common.Address]common.Address, []common.Address, error) {
	destTokens, err := r.cachedDestTokens.Get(ctx)
	if err != nil {
		return nil, nil, err
	}

	sourceToDestTokens := destTokens.SupportedTokens
	supportedDestTokens := make([]common.Address, 0, len(sourceToDestTokens))
	for _, destToken := range sourceToDestTokens {
		supportedDestTokens = append(supportedDestTokens, destToken)
	}
	return sourceToDestTokens, supportedDestTokens, nil
}

// Calculates a map that indicated whether a sequence number has already been executed
// before. It doesn't matter if the executed succeeded, since we don't retry previous
// attempts even if they failed. Value in the map indicates whether the log is finalized or not.
func (r *ExecutionReportingPlugin) getExecutedSeqNrsInRange(ctx context.Context, min, max uint64, latestBlock int64) (map[uint64]bool, error) {
	stateChanges, err := r.config.destReader.GetExecutionStateChangesBetweenSeqNums(
		ctx,
		r.config.offRamp.Address(),
		min,
		max,
		int(r.offchainConfig.DestOptimisticConfirmations),
	)
	if err != nil {
		return nil, err
	}
	executedMp := make(map[uint64]bool, len(stateChanges))
	for _, stateChange := range stateChanges {
		finalized := (latestBlock - stateChange.BlockNumber) >= int64(r.offchainConfig.DestFinalityDepth)
		executedMp[stateChange.Data.SequenceNumber] = finalized
	}
	return executedMp, nil
}

// Builds a batch of transactions that can be executed, takes into account
// the available gas, rate limiting, execution state, nonce state, and
// profitability of execution.
func (r *ExecutionReportingPlugin) buildBatch(
	ctx context.Context,
	lggr logger.Logger,
	report commitReportWithSendRequests,
	inflight []InflightInternalExecutionReport,
	aggregateTokenLimit *big.Int,
	sourceTokenPricesUSD map[common.Address]*big.Int,
	destTokenPricesUSD map[common.Address]*big.Int,
	execGasPriceEstimate cache.LazyFunction[*big.Int],
	sourceToDestToken map[common.Address]common.Address,
	destTokenPoolRateLimits map[common.Address]*big.Int,
) (executableMessages []ObservedMessage) {
	inflightSeqNrs, inflightAggregateValue, maxInflightSenderNonces, inflightTokenAmounts, err := inflightAggregates(inflight, destTokenPricesUSD, sourceToDestToken)
	if err != nil {
		lggr.Errorw("Unexpected error computing inflight values", "err", err)
		return []ObservedMessage{}
	}
	availableGas := uint64(r.offchainConfig.BatchGasLimit)
	expectedNonces := make(map[common.Address]uint64)
	availableDataLen := MaxDataLenPerBatch

	for _, msg := range report.sendRequestsWithMeta {
		msgLggr := lggr.With("messageID", hexutil.Encode(msg.MessageId[:]))
		if msg.Executed {
			msgLggr.Infow("Skipping message already executed", "seqNr", msg.SequenceNumber)
			continue
		}
		if _, isInflight := inflightSeqNrs[msg.SequenceNumber]; isInflight {
			msgLggr.Infow("Skipping message already inflight", "seqNr", msg.SequenceNumber)
			continue
		}
		if _, ok := expectedNonces[msg.Sender]; !ok {
			// First message in batch, need to populate expected nonce
			if maxInflight, ok := maxInflightSenderNonces[msg.Sender]; ok {
				// Sender already has inflight nonce, populate from there
				expectedNonces[msg.Sender] = maxInflight + 1
			} else {
				// Nothing inflight take from chain.
				// Chain holds existing nonce.
				nonce, err := r.config.offRamp.GetSenderNonce(nil, msg.Sender)
				if err != nil {
					lggr.Errorw("unable to get sender nonce", "err", err, "seqNr", msg.SequenceNumber)
					continue
				}
				expectedNonces[msg.Sender] = nonce + 1
			}
		}
		// Check expected nonce is valid
		if msg.Nonce != expectedNonces[msg.Sender] {
			msgLggr.Warnw("Skipping message invalid nonce", "have", msg.Nonce, "want", expectedNonces[msg.Sender])
			continue
		}

		if !r.isRateLimitEnoughForTokenPool(destTokenPoolRateLimits, msg.TokenAmounts, inflightTokenAmounts, sourceToDestToken) {
			msgLggr.Warnw("Skipping message token pool rate limit hit")
			continue
		}

		msgValue, err := aggregateTokenValue(destTokenPricesUSD, sourceToDestToken, msg.TokenAmounts)
		if err != nil {
			msgLggr.Errorw("Skipping message unable to compute aggregate value", "err", err)
			continue
		}

		// if token limit is smaller than message value skip message
		if tokensLeft, hasCapacity := hasEnoughTokens(aggregateTokenLimit, msgValue, inflightAggregateValue); !hasCapacity {
			msgLggr.Warnw("token limit is smaller than message value", "aggregateTokenLimit", tokensLeft.String(), "msgValue", msgValue.String())
			continue
		}

		tokenData, ready, err2 := getTokenData(ctx, msgLggr, msg, r.config.tokenDataProviders)
		if err2 != nil {
			msgLggr.Errorw("Skipping message unable to check token data", "err", err2)
			continue
		}
		if !ready {
			msgLggr.Warnw("Skipping message attestation not ready")
			continue
		}

		// Fee boosting
		execGasPriceEstimateValue, err := execGasPriceEstimate()
		if err != nil {
			msgLggr.Errorw("Unexpected error fetching gas price estimate", "err", err)
			return []ObservedMessage{}
		}

		dstWrappedNativePrice, exists := destTokenPricesUSD[r.destWrappedNative]
		if !exists {
			msgLggr.Errorw("token not in dst token prices", "token", r.destWrappedNative)
			continue
		}

		execCostUsd := computeExecCost(msg.GasLimit, execGasPriceEstimateValue, dstWrappedNativePrice)
		// calculating the source chain fee, dividing by 1e18 for denomination.
		// For example:
		// FeeToken=link; FeeTokenAmount=1e17 i.e. 0.1 link, price is 6e18 USD/link (1 USD = 1e18),
		// availableFee is 1e17*6e18/1e18 = 6e17 = 0.6 USD

		sourceFeeTokenPrice, exists := sourceTokenPricesUSD[msg.FeeToken]
		if !exists {
			msgLggr.Errorw("token not in source token prices", "token", msg.FeeToken)
			continue
		}

		if len(msg.Data) > availableDataLen {
			msgLggr.Infow("Skipping message, insufficient remaining batch data len",
				"msgDataLen", len(msg.Data), "availableBatchDataLen", availableDataLen)
			continue
		}

		availableFee := big.NewInt(0).Mul(msg.FeeTokenAmount, sourceFeeTokenPrice)
		availableFee = availableFee.Div(availableFee, big.NewInt(1e18))
		availableFeeUsd := waitBoostedFee(time.Since(msg.BlockTimestamp), availableFee, r.offchainConfig.RelativeBoostPerWaitHour)
		if availableFeeUsd.Cmp(execCostUsd) < 0 {
			msgLggr.Infow("Insufficient remaining fee", "availableFeeUsd", availableFeeUsd, "execCostUsd", execCostUsd,
				"sourceBlockTimestamp", msg.BlockTimestamp, "waitTime", time.Since(msg.BlockTimestamp), "boost", r.offchainConfig.RelativeBoostPerWaitHour)
			continue
		}

		messageMaxGas, err := calculateMessageMaxGas(
			msg.GasLimit,
			len(report.sendRequestsWithMeta),
			len(msg.Data),
			len(msg.TokenAmounts),
		)
		if err != nil {
			msgLggr.Errorw("calculate message max gas", "err", err)
			continue
		}

		// Check sufficient gas in batch
		if availableGas < messageMaxGas {
			msgLggr.Infow("Insufficient remaining gas in batch limit", "availableGas", availableGas, "messageMaxGas", messageMaxGas)
			continue
		}
		availableGas -= messageMaxGas
		aggregateTokenLimit.Sub(aggregateTokenLimit, msgValue)
		for _, tk := range msg.TokenAmounts {
			dstToken, exists := sourceToDestToken[tk.Token]
			if !exists {
				msgLggr.Warnw("destination token does not exist", "token", tk.Token)
				continue
			}
			if rl, exists := destTokenPoolRateLimits[dstToken]; exists {
				destTokenPoolRateLimits[dstToken] = rl.Sub(rl, tk.Amount)
			}
		}

		msgLggr.Infow("Adding msg to batch", "seqNum", msg.SequenceNumber, "nonce", msg.Nonce,
			"value", msgValue, "aggregateTokenLimit", aggregateTokenLimit)
		executableMessages = append(executableMessages, NewObservedMessage(msg.SequenceNumber, tokenData))

		// after message is added to the batch, decrease the available data length
		availableDataLen -= len(msg.Data)

		expectedNonces[msg.Sender] = msg.Nonce + 1
	}
	return executableMessages
}

func getTokenData(ctx context.Context, lggr logger.Logger, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenDataProviders map[common.Address]tokendata.Reader) (tokenData [][]byte, allReady bool, err error) {
	for _, token := range msg.TokenAmounts {
		offchainTokenDataProvider, ok := tokenDataProviders[token.Token]
		if !ok {
			// No token data required
			tokenData = append(tokenData, []byte{})
			continue
		}
		lggr.Infow("Fetching token data", "token", token.Token.Hex())
		tknData, err2 := offchainTokenDataProvider.ReadTokenData(ctx, msg)
		if err2 != nil {
			if errors.Is(err2, tokendata.ErrNotReady) {
				lggr.Infow("Token data not ready yet", "token", token.Token.Hex())
				return [][]byte{}, false, nil
			}
			return [][]byte{}, false, err2
		}

		lggr.Infow("Token data retrieved", "token", token.Token.Hex())
		tokenData = append(tokenData, tknData)
	}
	return tokenData, true, nil
}

func (r *ExecutionReportingPlugin) isRateLimitEnoughForTokenPool(
	destTokenPoolRateLimits map[common.Address]*big.Int,
	sourceTokenAmounts []evm_2_evm_offramp.ClientEVMTokenAmount,
	inflightTokenAmounts map[common.Address]*big.Int,
	sourceToDestToken map[common.Address]common.Address,
) bool {
	rateLimitsCopy := make(map[common.Address]*big.Int)
	for destToken, rl := range destTokenPoolRateLimits {
		rateLimitsCopy[destToken] = new(big.Int).Set(rl)
	}

	for sourceToken, amount := range inflightTokenAmounts {
		if destToken, exists := sourceToDestToken[sourceToken]; exists {
			if rl, exists := rateLimitsCopy[destToken]; exists {
				rateLimitsCopy[destToken] = rl.Sub(rl, amount)
			}
		}
	}

	for _, sourceToken := range sourceTokenAmounts {
		destToken, exists := sourceToDestToken[sourceToken.Token]
		if !exists {
			r.lggr.Warnw("dest token not found", "sourceToken", sourceToken.Token)
			continue
		}

		rl, exists := rateLimitsCopy[destToken]
		if !exists {
			r.lggr.Debugw("rate limit not applied to token", "token", destToken)
			continue
		}

		if rl.Cmp(sourceToken.Amount) < 0 {
			r.lggr.Warnw("token pool rate limit reached",
				"token", sourceToken.Token, "destToken", destToken, "amount", sourceToken.Amount, "rateLimit", rl)
			return false
		}
		rateLimitsCopy[destToken] = rl.Sub(rl, sourceToken.Amount)
	}

	return true
}

func hasEnoughTokens(tokenLimit *big.Int, msgValue *big.Int, inflightValue *big.Int) (*big.Int, bool) {
	tokensLeft := big.NewInt(0).Sub(tokenLimit, inflightValue)
	return tokensLeft, tokensLeft.Cmp(msgValue) >= 0
}

func calculateMessageMaxGas(gasLimit *big.Int, numRequests, dataLen, numTokens int) (uint64, error) {
	if !gasLimit.IsUint64() {
		return 0, fmt.Errorf("gas limit %s cannot be casted to uint64", gasLimit)
	}

	gasLimitU64 := gasLimit.Uint64()
	gasOverHeadGas := maxGasOverHeadGas(numRequests, dataLen, numTokens)
	messageMaxGas := gasLimitU64 + gasOverHeadGas

	if messageMaxGas < gasLimitU64 || messageMaxGas < gasOverHeadGas {
		return 0, fmt.Errorf("message max gas overflow, gasLimit=%d gasOverHeadGas=%d", gasLimitU64, gasOverHeadGas)
	}

	return messageMaxGas, nil
}

// helper struct to hold the commitReport and the related send requests
type commitReportWithSendRequests struct {
	commitReport         commit_store.CommitStoreCommitReport
	sendRequestsWithMeta []internal.EVM2EVMOnRampCCIPSendRequestedWithMeta
}

func (r *commitReportWithSendRequests) validate() error {
	// make sure that number of messages is the expected
	if exp := int(r.commitReport.Interval.Max - r.commitReport.Interval.Min + 1); len(r.sendRequestsWithMeta) != exp {
		return errors.Errorf(
			"unexpected missing sendRequestsWithMeta in committed root %x have %d want %d", r.commitReport.MerkleRoot, len(r.sendRequestsWithMeta), exp)
	}

	return nil
}

func (r *commitReportWithSendRequests) allRequestsAreExecutedAndFinalized() bool {
	for _, req := range r.sendRequestsWithMeta {
		if !req.Executed || !req.Finalized {
			return false
		}
	}
	return true
}

// checks if the send request fits the commit report interval
func (r *commitReportWithSendRequests) sendReqFits(sendReq internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) bool {
	return sendReq.SequenceNumber >= r.commitReport.Interval.Min &&
		sendReq.SequenceNumber <= r.commitReport.Interval.Max
}

// getReportsWithSendRequests returns the target reports with populated send requests.
func (r *ExecutionReportingPlugin) getReportsWithSendRequests(
	ctx context.Context,
	reports []commit_store.CommitStoreCommitReport,
) ([]commitReportWithSendRequests, error) {
	if len(reports) == 0 {
		return nil, nil
	}

	// find interval from all the reports
	intervalMin := reports[0].Interval.Min
	intervalMax := reports[0].Interval.Max
	for _, report := range reports[1:] {
		if report.Interval.Max > intervalMax {
			intervalMax = report.Interval.Max
		}
		if report.Interval.Min < intervalMin {
			intervalMin = report.Interval.Min
		}
	}

	// use errgroup to fetch send request logs and executed sequence numbers in parallel
	eg := &errgroup.Group{}

	var sendRequests []ccipdata.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested]
	eg.Go(func() error {
		sendReqs, err := r.config.sourceReader.GetSendRequestsBetweenSeqNums(
			ctx,
			r.config.onRamp.Address(),
			intervalMin,
			intervalMax,
			int(r.offchainConfig.SourceFinalityDepth),
		)
		if err != nil {
			return err
		}
		sendRequests = sendReqs
		return nil
	})

	var executedSeqNums map[uint64]bool
	eg.Go(func() error {
		latestBlock, err := r.config.destReader.LatestBlock(ctx)
		if err != nil {
			return err
		}

		// get executable sequence numbers
		executedMp, err := r.getExecutedSeqNrsInRange(ctx, intervalMin, intervalMax, latestBlock)
		if err != nil {
			return err
		}
		executedSeqNums = executedMp
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	reportsWithSendReqs := make([]commitReportWithSendRequests, len(reports))
	for i, report := range reports {
		reportsWithSendReqs[i] = commitReportWithSendRequests{
			commitReport:         report,
			sendRequestsWithMeta: make([]internal.EVM2EVMOnRampCCIPSendRequestedWithMeta, 0, report.Interval.Max-report.Interval.Min+1),
		}
	}

	for _, sendReq := range sendRequests {
		msg := abihelpers.OnRampMessageToOffRampMessage(sendReq.Data.Message)

		// if value exists in the map then it's executed
		// if value exists, and it's true then it's considered finalized
		finalized, executed := executedSeqNums[msg.SequenceNumber]

		reqWithMeta := internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{
			InternalEVM2EVMMessage: msg,
			BlockTimestamp:         sendReq.BlockTimestamp,
			Executed:               executed,
			Finalized:              finalized,
			LogIndex:               sendReq.Data.Raw.Index,
			TxHash:                 sendReq.Data.Raw.TxHash,
		}

		// attach the msg to the appropriate reports
		for i := range reportsWithSendReqs {
			if reportsWithSendReqs[i].sendReqFits(reqWithMeta) {
				reportsWithSendReqs[i].sendRequestsWithMeta = append(reportsWithSendReqs[i].sendRequestsWithMeta, reqWithMeta)
			}
		}
	}

	return reportsWithSendReqs, nil
}

func aggregateTokenValue(destTokenPricesUSD map[common.Address]*big.Int, sourceToDest map[common.Address]common.Address, tokensAndAmount []evm_2_evm_offramp.ClientEVMTokenAmount) (*big.Int, error) {
	sum := big.NewInt(0)
	for i := 0; i < len(tokensAndAmount); i++ {
		price, ok := destTokenPricesUSD[sourceToDest[tokensAndAmount[i].Token]]
		if !ok {
			return nil, errors.Errorf("do not have price for source token %v", tokensAndAmount[i].Token)
		}
		sum.Add(sum, new(big.Int).Quo(new(big.Int).Mul(price, tokensAndAmount[i].Amount), big.NewInt(1e18)))
	}
	return sum, nil
}

// Assumes non-empty report. Messages to execute can span more than one report, but are assumed to be in order of increasing
// sequence number.
func (r *ExecutionReportingPlugin) buildReport(ctx context.Context, lggr logger.Logger, observedMessages []ObservedMessage) ([]byte, error) {
	if err := validateSeqNumbers(ctx, r.config.commitStore, observedMessages); err != nil {
		return nil, err
	}
	commitReport, err := getCommitReportForSeqNum(ctx, r.config.destReader, r.config.commitStore, observedMessages[0].SeqNr)
	if err != nil {
		return nil, err
	}
	lggr.Infow("Building execution report", "observations", observedMessages, "merkleRoot", hexutil.Encode(commitReport.MerkleRoot[:]), "report", commitReport)

	sendReqsInRoot, leaves, tree, err := getProofData(ctx, lggr, r.config.leafHasher, r.config.onRamp.Address(), r.config.sourceReader, commitReport.Interval)
	if err != nil {
		return nil, err
	}

	messages := make([]*evm_2_evm_offramp.InternalEVM2EVMMessage, len(sendReqsInRoot))
	for i, msg := range sendReqsInRoot {
		offRampMsg := abihelpers.OnRampMessageToOffRampMessage(msg.Data.Message)
		messages[i] = &offRampMsg
	}

	// cap messages which fits MaxExecutionReportLength (after serialized)
	capped := sort.Search(len(observedMessages), func(i int) bool {
		report, _, err2 := buildExecutionReportForMessages(messages, leaves, tree, commitReport.Interval, observedMessages[:i+1])
		if err2 != nil {
			r.lggr.Errorw("build execution report", "err", err2)
			return false
		}

		var encoded []byte
		encoded, err = abihelpers.EncodeExecutionReport(report)
		if err != nil {
			// false makes Search keep looking to the right, always including any "erroring" ObservedMessage and allowing us to detect in the bottom
			return false
		}
		return len(encoded) > MaxExecutionReportLength
	})
	if err != nil {
		return nil, err
	}

	execReport, hashes, err := buildExecutionReportForMessages(messages, leaves, tree, commitReport.Interval, observedMessages[:capped])
	if err != nil {
		return nil, err
	}

	encodedReport, err := abihelpers.EncodeExecutionReport(execReport)
	if err != nil {
		return nil, err
	}

	if capped < len(observedMessages) {
		lggr.Warnf(
			"Capping report to fit MaxExecutionReportLength: msgsCount %d -> %d, bytes %d, bytesLimit %d",
			len(observedMessages), capped, len(encodedReport), MaxExecutionReportLength,
		)
	}

	// Double check this verifies before sending.
	res, err := r.config.commitStore.Verify(&bind.CallOpts{Context: ctx}, hashes, execReport.Proofs, execReport.ProofFlagBits)
	if err != nil {
		lggr.Errorw("Unable to call verify", "observations", observedMessages[:capped], "root", commitReport.MerkleRoot[:], "seqRange", commitReport.Interval, "err", err)
		return nil, err
	}
	// No timestamp, means failed to verify root.
	if res.Cmp(big.NewInt(0)) == 0 {
		root := tree.Root()
		lggr.Errorf("Root does not verify for messages: %v, our inner root 0x%x", observedMessages[:capped], root)
		return nil, errors.New("root does not verify")
	}
	return encodedReport, nil
}

func (r *ExecutionReportingPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	lggr := r.lggr.Named("ExecutionReport")
	parsableObservations := getParsableObservations[ExecutionObservation](lggr, observations)
	// Need at least F+1 observations
	if len(parsableObservations) <= r.F {
		lggr.Warn("Non-empty observations <= F, need at least F+1 to continue")
		return false, nil, nil
	}

	observedMessages, err := calculateObservedMessagesConsensus(parsableObservations, r.F)
	if err != nil {
		return false, nil, err
	}
	if len(observedMessages) == 0 {
		return false, nil, nil
	}

	report, err := r.buildReport(ctx, lggr, observedMessages)
	if err != nil {
		return false, nil, err
	}
	lggr.Infow("Report", "executableObservations", observedMessages)
	return true, report, nil
}

type tallyKey struct {
	seqNr         uint64
	tokenDataHash [32]byte
}

type tallyVal struct {
	tally     int
	tokenData [][]byte
}

func calculateObservedMessagesConsensus(observations []ExecutionObservation, f int) ([]ObservedMessage, error) {
	tally := make(map[tallyKey]tallyVal)
	for _, obs := range observations {
		for seqNr, msgData := range obs.Messages {
			tokenDataHash, err := hashlib.BytesOfBytesKeccak(msgData.TokenData)
			if err != nil {
				return nil, fmt.Errorf("bytes of bytes keccak: %w", err)
			}

			key := tallyKey{seqNr: seqNr, tokenDataHash: tokenDataHash}
			if val, ok := tally[key]; ok {
				tally[key] = tallyVal{tally: val.tally + 1, tokenData: msgData.TokenData}
			} else {
				tally[key] = tallyVal{tally: 1, tokenData: msgData.TokenData}
			}
		}
	}

	// We might have different token data for the same sequence number.
	// For that purpose we want to keep the token data with the most occurrences.
	seqNumTally := make(map[uint64]tallyVal)

	// order tally keys to make looping over the entries deterministic
	tallyKeys := make([]tallyKey, 0, len(tally))
	for key := range tally {
		tallyKeys = append(tallyKeys, key)
	}
	sort.Slice(tallyKeys, func(i, j int) bool {
		return hex.EncodeToString(tallyKeys[i].tokenDataHash[:]) < hex.EncodeToString(tallyKeys[j].tokenDataHash[:])
	})

	for _, key := range tallyKeys {
		tallyInfo := tally[key]
		existingTally, exists := seqNumTally[key.seqNr]
		if tallyInfo.tally > f && (!exists || tallyInfo.tally > existingTally.tally) {
			seqNumTally[key.seqNr] = tallyInfo
		}
	}

	finalSequenceNumbers := make([]ObservedMessage, 0, len(seqNumTally))
	for seqNr, tallyInfo := range seqNumTally {
		finalSequenceNumbers = append(finalSequenceNumbers, NewObservedMessage(seqNr, tallyInfo.tokenData))
	}
	// buildReport expects sorted sequence numbers (tally map is non-deterministic).
	sort.Slice(finalSequenceNumbers, func(i, j int) bool {
		return finalSequenceNumbers[i].SeqNr < finalSequenceNumbers[j].SeqNr
	})
	return finalSequenceNumbers, nil
}

func (r *ExecutionReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("ShouldAcceptFinalizedReport")
	messages, err := abihelpers.MessagesFromExecutionReport(report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, err
	}
	lggr = lggr.With("messageIDs", contractutil.GetMessageIDsAsHexString(messages))

	// If the first message is executed already, this execution report is stale, and we do not accept it.
	stale, err := r.isStaleReport(messages)
	if err != nil {
		return false, err
	}
	if stale {
		lggr.Info("Execution report is stale")
		return false, nil
	}
	// Else just assume in flight
	if err = r.inflightReports.add(lggr, messages); err != nil {
		return false, err
	}
	lggr.Info("Accepting finalized report")
	return true, nil
}

func (r *ExecutionReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("ShouldTransmitAcceptedReport")
	messages, err := abihelpers.MessagesFromExecutionReport(report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, nil
	}
	lggr = lggr.With("messageIDs", contractutil.GetMessageIDsAsHexString(messages))

	// If report is not stale we transmit.
	// When the executeTransmitter enqueues the tx for tx manager,
	// we mark it as execution_sent, removing it from the set of inflight messages.
	stale, err := r.isStaleReport(messages)
	if err != nil {
		return false, err
	}
	if stale {
		lggr.Info("Execution report is stale")
		return false, nil
	}

	lggr.Info("Transmitting finalized report")
	return true, err
}

func (r *ExecutionReportingPlugin) isStaleReport(messages []evm_2_evm_offramp.InternalEVM2EVMMessage) (bool, error) {
	if len(messages) == 0 {
		return true, fmt.Errorf("messages are empty")
	}

	// If the first message is executed already, this execution report is stale.
	// Note the default execution state, including for arbitrary seq number not yet committed
	// is ExecutionStateUntouched.
	msgState, err := r.config.offRamp.GetExecutionState(nil, messages[0].SequenceNumber)
	if err != nil {
		return true, err
	}
	if state := abihelpers.MessageExecutionState(msgState); state == abihelpers.ExecutionStateFailure || state == abihelpers.ExecutionStateSuccess {
		return true, nil
	}

	return false, nil
}

func (r *ExecutionReportingPlugin) Close() error {
	return nil
}

func inflightAggregates(
	inflight []InflightInternalExecutionReport,
	destTokenPrices map[common.Address]*big.Int,
	sourceToDest map[common.Address]common.Address,
) (map[uint64]struct{}, *big.Int, map[common.Address]uint64, map[common.Address]*big.Int, error) {
	inflightSeqNrs := make(map[uint64]struct{})
	inflightAggregateValue := big.NewInt(0)
	maxInflightSenderNonces := make(map[common.Address]uint64)
	inflightTokenAmounts := make(map[common.Address]*big.Int)

	for _, rep := range inflight {
		for _, message := range rep.messages {
			inflightSeqNrs[message.SequenceNumber] = struct{}{}
			msgValue, err := aggregateTokenValue(destTokenPrices, sourceToDest, message.TokenAmounts)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			inflightAggregateValue.Add(inflightAggregateValue, msgValue)
			maxInflightSenderNonce, ok := maxInflightSenderNonces[message.Sender]
			if !ok || message.Nonce > maxInflightSenderNonce {
				maxInflightSenderNonces[message.Sender] = message.Nonce
			}

			for _, tk := range message.TokenAmounts {
				if rl, exists := inflightTokenAmounts[tk.Token]; exists {
					inflightTokenAmounts[tk.Token] = rl.Add(rl, tk.Amount)
				} else {
					inflightTokenAmounts[tk.Token] = new(big.Int).Set(tk.Amount)
				}
			}
		}
	}
	return inflightSeqNrs, inflightAggregateValue, maxInflightSenderNonces, inflightTokenAmounts, nil
}

// getTokensPrices returns token prices of the given price registry,
// results include feeTokens and passed-in tokens
// price values are USD per 1e18 of smallest token denomination, in base units 1e18 (e.g. 5$ = 5e18 USD per 1e18 units).
// this function is used for price registry of both source and destination chains.
func getTokensPrices(ctx context.Context, feeTokens []common.Address, priceRegistry price_registry.PriceRegistryInterface, tokens []common.Address) (map[common.Address]*big.Int, error) {
	priceRegistryAddress := priceRegistry.Address()
	prices := make(map[common.Address]*big.Int)

	wantedTokens := append(feeTokens, tokens...)
	fetchedPrices, err := priceRegistry.GetTokenPrices(&bind.CallOpts{Context: ctx}, wantedTokens)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get token prices of %v", wantedTokens)
	}

	// price registry should always return a price per token ordered by input tokens
	if len(fetchedPrices) != len(wantedTokens) {
		return nil, fmt.Errorf("token prices length exp=%d actual=%d", len(wantedTokens), len(fetchedPrices))
	}

	for i, token := range wantedTokens {
		// price of a token can never be zero
		if fetchedPrices[i].Value.BitLen() == 0 {
			return nil, fmt.Errorf("price of token %s is zero (price registry=%s)", token, priceRegistryAddress)
		}

		// price registry should not report different price for the same token
		price, exists := prices[token]
		if exists && fetchedPrices[i].Value.Cmp(price) != 0 {
			return nil, fmt.Errorf("price registry reported different prices (%s and %s) for the same token %s",
				fetchedPrices[i].Value, price, token)
		}

		prices[token] = fetchedPrices[i].Value
	}

	return prices, nil
}

func getUnexpiredCommitReports(
	ctx context.Context,
	destReader ccipdata.Reader,
	commitStore commit_store.CommitStoreInterface,
	permissionExecutionThreshold time.Duration,
) ([]commit_store.CommitStoreCommitReport, error) {
	acceptedReports, err := destReader.GetAcceptedCommitReportsGteTimestamp(
		ctx,
		commitStore.Address(),
		time.Now().Add(-permissionExecutionThreshold),
		0,
	)
	if err != nil {
		return nil, err
	}

	var reports []commit_store.CommitStoreCommitReport
	for _, acceptedReport := range acceptedReports {
		reports = append(reports, acceptedReport.Data.Report)
	}
	return reports, nil
}
