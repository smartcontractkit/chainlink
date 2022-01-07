package evm

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
)

// configMailboxSanityLimit is the maximum number of configs that can be held
// in the mailbox. Under normal operation there should never be more than 0 or
// 1 configs in the mailbox, this limit is here merely to prevent unbounded usage
// in some kind of unforeseen insane situation.
const configMailboxSanityLimit = 100

var (
	_ ocrtypes.ContractConfigTracker = &ContractTracker{}
	_ httypes.HeadTrackable          = &ContractTracker{}

	OCRContractConfigSet = getEventTopic("ConfigSet")
)

type OCRContractTrackerDB interface {
	SaveLatestRoundRequested(tx pg.Queryer, rr ocr2aggregator.OCR2AggregatorRoundRequested) error
	LoadLatestRoundRequested() (rr ocr2aggregator.OCR2AggregatorRoundRequested, err error)
}

// ContractTracker complies with ContractConfigTracker interface and
// handles log events related to the contract more generally
//go:generate mockery --name OCRContractTrackerDB --output ./mocks/ --case=underscore
type ContractTracker struct {
	utils.StartStopOnce

	ethClient        evmclient.Client
	contract         *offchain_aggregator_wrapper.OffchainAggregator
	contractFilterer *ocr2aggregator.OCR2AggregatorFilterer
	contractCaller   *ocr2aggregator.OCR2AggregatorCaller
	logBroadcaster   log.Broadcaster
	jobID            int32
	logger           logger.Logger
	odb              OCRContractTrackerDB
	q                pg.Q
	blockTranslator  ocrcommon.BlockTranslator
	chain            ocrcommon.Config

	// HeadBroadcaster
	headBroadcaster  httypes.HeadBroadcaster
	unsubscribeHeads func()

	// Start/Stop lifecycle
	ctx             context.Context
	ctxCancel       context.CancelFunc
	wg              sync.WaitGroup
	unsubscribeLogs func()

	// LatestRoundRequested
	latestRoundRequested ocr2aggregator.OCR2AggregatorRoundRequested
	lrrMu                sync.RWMutex

	// ContractConfig
	configsMB utils.Mailbox
	chConfigs chan ocrtypes.ContractConfig

	// LatestBlockHeight
	latestBlockHeight   int64
	latestBlockHeightMu sync.RWMutex
}

// NewOCRContractTracker makes a new ContractTracker
func NewOCRContractTracker(
	contract *offchain_aggregator_wrapper.OffchainAggregator,
	contractFilterer *ocr2aggregator.OCR2AggregatorFilterer,
	contractCaller *ocr2aggregator.OCR2AggregatorCaller,
	ethClient evmclient.Client,
	logBroadcaster log.Broadcaster,
	jobID int32,
	logger logger.Logger,
	db *sqlx.DB,
	odb OCRContractTrackerDB,
	chain ocrcommon.Config,
	headBroadcaster httypes.HeadBroadcaster,
) (o *ContractTracker) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ContractTracker{
		utils.StartStopOnce{},
		ethClient,
		contract,
		contractFilterer,
		contractCaller,
		logBroadcaster,
		jobID,
		logger,
		odb,
		pg.NewQ(db, logger, chain),
		ocrcommon.NewBlockTranslator(chain, ethClient, logger),
		chain,
		headBroadcaster,
		nil,
		ctx,
		cancel,
		sync.WaitGroup{},
		nil,
		ocr2aggregator.OCR2AggregatorRoundRequested{},
		sync.RWMutex{},
		*utils.NewMailbox(configMailboxSanityLimit),
		make(chan ocrtypes.ContractConfig),
		-1,
		sync.RWMutex{},
	}
}

// Start must be called before logs can be delivered
// It ought to be called before starting OCR
func (t *ContractTracker) Start() error {
	return t.StartOnce("ContractTracker", func() (err error) {
		t.latestRoundRequested, err = t.odb.LoadLatestRoundRequested()
		if err != nil {
			return errors.Wrap(err, "ContractTracker#Start: failed to load latest round requested")
		}

		t.unsubscribeLogs = t.logBroadcaster.Register(t, log.ListenerOpts{
			Contract: t.contract.Address(),
			ParseLog: t.contract.ParseLog,
			LogsWithTopics: map[gethCommon.Hash][][]log.Topic{
				offchain_aggregator_wrapper.OffchainAggregatorRoundRequested{}.Topic(): nil,
				offchain_aggregator_wrapper.OffchainAggregatorConfigSet{}.Topic():      nil,
			},
			MinIncomingConfirmations: 1,
		})

		var latestHead *evmtypes.Head
		latestHead, t.unsubscribeHeads = t.headBroadcaster.Subscribe(t)
		if latestHead != nil {
			t.setLatestBlockHeight(*latestHead)
		}

		t.wg.Add(1)
		go t.processLogs()
		return nil
	})
}

// Close should be called after teardown of the OCR job relying on this tracker
func (t *ContractTracker) Close() error {
	return t.StopOnce("ContractTracker", func() error {
		t.ctxCancel()
		t.wg.Wait()
		t.unsubscribeHeads()
		t.unsubscribeLogs()
		close(t.chConfigs)
		return nil
	})
}

// Connect conforms to HeadTrackable
func (t *ContractTracker) Connect(*evmtypes.Head) error { return nil }

// OnNewLongestChain conformed to HeadTrackable and updates latestBlockHeight
func (t *ContractTracker) OnNewLongestChain(_ context.Context, h *evmtypes.Head) {
	t.setLatestBlockHeight(*h)
}

func (t *ContractTracker) setLatestBlockHeight(h evmtypes.Head) {
	var num int64
	if h.L1BlockNumber.Valid {
		num = h.L1BlockNumber.Int64
	} else {
		num = h.Number
	}
	t.latestBlockHeightMu.Lock()
	defer t.latestBlockHeightMu.Unlock()
	if num > t.latestBlockHeight {
		t.latestBlockHeight = num
	}
}

func (t *ContractTracker) getLatestBlockHeight() int64 {
	t.latestBlockHeightMu.RLock()
	defer t.latestBlockHeightMu.RUnlock()
	return t.latestBlockHeight
}

func (t *ContractTracker) processLogs() {
	defer t.wg.Done()
	for {
		select {
		case <-t.configsMB.Notify():
			// NOTE: libocr could take an arbitrary amount of time to process a
			// new config. To avoid blocking the log broadcaster, we use this
			// background thread to deliver them and a mailbox as the buffer.
			for {
				x, exists := t.configsMB.Retrieve()
				if !exists {
					break
				}
				cc, ok := x.(ocrtypes.ContractConfig)
				if !ok {
					panic(fmt.Sprintf("expected ocrtypes.ContractConfig but got %T", x))
				}
				select {
				case t.chConfigs <- cc:
				case <-t.ctx.Done():
					return
				}
			}
		case <-t.ctx.Done():
			return
		}
	}
}

// HandleLog complies with LogListener interface
// It is not thread safe
func (t *ContractTracker) HandleLog(lb log.Broadcast) {
	was, err := t.logBroadcaster.WasAlreadyConsumed(lb)
	if err != nil {
		t.logger.Errorw("OCRContract: could not determine if log was already consumed", "error", err)
		return
	} else if was {
		return
	}

	raw := lb.RawLog()
	if raw.Address != t.contract.Address() {
		t.logger.Errorf("log address of 0x%x does not match configured contract address of 0x%x", raw.Address, t.contract.Address())
		t.logger.ErrorIf(t.logBroadcaster.MarkConsumed(lb), "unable to mark consumed")
		return
	}
	topics := raw.Topics
	if len(topics) == 0 {
		t.logger.ErrorIf(t.logBroadcaster.MarkConsumed(lb), "unable to mark consumed")
		return
	}

	var consumed bool
	switch topics[0] {
	case offchain_aggregator_wrapper.OffchainAggregatorConfigSet{}.Topic():
		var configSet *ocr2aggregator.OCR2AggregatorConfigSet
		configSet, err = t.contractFilterer.ParseConfigSet(raw)
		if err != nil {
			t.logger.Errorw("could not parse config set", "err", err)
			t.logger.ErrorIf(t.logBroadcaster.MarkConsumed(lb), "unable to mark consumed")
			return
		}
		configSet.Raw = lb.RawLog()
		cc := evmutil.ContractConfigFromConfigSetEvent(*configSet)

		wasOverCapacity := t.configsMB.Deliver(cc)
		if wasOverCapacity {
			t.logger.Error("config mailbox is over capacity - dropped the oldest unprocessed item")
		}
	case offchain_aggregator_wrapper.OffchainAggregatorRoundRequested{}.Topic():
		var rr *ocr2aggregator.OCR2AggregatorRoundRequested
		rr, err = t.contractFilterer.ParseRoundRequested(raw)
		if err != nil {
			t.logger.Errorw("could not parse round requested", "err", err)
			t.logger.ErrorIf(t.logBroadcaster.MarkConsumed(lb), "unable to mark consumed")
			return
		}
		if IsLaterThan(raw, t.latestRoundRequested.Raw) {
			err = t.q.Transaction(func(q pg.Queryer) error {
				if err = t.odb.SaveLatestRoundRequested(q, *rr); err != nil {
					return err
				}
				return t.logBroadcaster.MarkConsumed(lb, pg.WithQueryer(q))
			})
			if err != nil {
				t.logger.Error(err)
				return
			}
			consumed = true
			t.lrrMu.Lock()
			t.latestRoundRequested = *rr
			t.lrrMu.Unlock()
			t.logger.Infow("ContractTracker: received new latest RoundRequested event", "latestRoundRequested", *rr)
		} else {
			t.logger.Warnw("ContractTracker: ignoring out of date RoundRequested event", "latestRoundRequested", t.latestRoundRequested, "roundRequested", rr)
		}
	default:
		t.logger.Debugw("ContractTracker: got unrecognised log topic", "topic", topics[0])
	}
	if !consumed {
		t.logger.ErrorIf(t.logBroadcaster.MarkConsumed(lb), "unable to mark consumed")
	}
}

// IsLaterThan returns true if the first log was emitted "after" the second log
// from the blockchain's point of view
func IsLaterThan(incoming gethTypes.Log, existing gethTypes.Log) bool {
	return incoming.BlockNumber > existing.BlockNumber ||
		(incoming.BlockNumber == existing.BlockNumber && incoming.TxIndex > existing.TxIndex) ||
		(incoming.BlockNumber == existing.BlockNumber && incoming.TxIndex == existing.TxIndex && incoming.Index > existing.Index)
}

// IsV2Job complies with LogListener interface
func (t *ContractTracker) IsV2Job() bool {
	return true
}

// JobID complies with LogListener interface
func (t *ContractTracker) JobID() int32 {
	return t.jobID
}

// Notify returns a channel that can wake up the contract tracker to let it
// know when a new config is available
func (t *ContractTracker) Notify() <-chan struct{} {
	return nil
}

// LatestConfigDetails queries the eth node
func (t *ContractTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	var cancel context.CancelFunc
	ctx, cancel = utils.CombinedContext(t.ctx, ctx)
	defer cancel()

	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := t.contract.LatestConfigDetails(&opts)
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting LatestConfigDetails")
	}
	configDigest, err = ocrtypes.BytesToConfigDigest(result.ConfigDigest[:])
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting config digest")
	}
	return uint64(result.BlockNumber), configDigest, err
}

// Return the latest configuration
func (t *ContractTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	fromBlock, toBlock := t.blockTranslator.NumberToQueryRange(ctx, changedInBlock)
	q := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []gethCommon.Address{t.contract.Address()},
		Topics: [][]gethCommon.Hash{
			{OCRContractConfigSet},
		},
	}

	var cancel context.CancelFunc
	ctx, cancel = utils.CombinedContext(t.ctx, ctx)
	defer cancel()

	logs, err := t.ethClient.FilterLogs(ctx, q)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(logs) == 0 {
		return ocrtypes.ContractConfig{}, errors.Errorf("ConfigFromLogs: OCRContract with address 0x%x has no logs", t.contract.Address())
	}

	latest, err := t.contractFilterer.ParseConfigSet(logs[len(logs)-1])
	if err != nil {
		return ocrtypes.ContractConfig{}, errors.Wrap(err, "ConfigFromLogs failed to ParseConfigSet")
	}
	latest.Raw = logs[len(logs)-1]
	if latest.Raw.Address != t.contract.Address() {
		return ocrtypes.ContractConfig{}, errors.Errorf("log address of 0x%x does not match configured contract address of 0x%x", latest.Raw.Address, t.contract.Address())
	}
	return evmutil.ContractConfigFromConfigSetEvent(*latest), err
}

// LatestBlockHeight queries the eth node for the most recent header
func (t *ContractTracker) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	// We skip confirmation checking anyway on Optimism so there's no need to
	// care about the block height; we have no way of getting the L1 block
	// height anyway
	if t.chain.ChainType() != "" {
		return 0, nil
	}
	latestBlockHeight := t.getLatestBlockHeight()
	if latestBlockHeight >= 0 {
		return uint64(latestBlockHeight), nil
	}

	t.logger.Debugw("ContractTracker: still waiting for first head, falling back to on-chain lookup")

	var cancel context.CancelFunc
	ctx, cancel = utils.CombinedContext(t.ctx, ctx)
	defer cancel()

	h, err := t.ethClient.HeadByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	if h == nil {
		return 0, errors.New("got nil head")
	}

	if h.L1BlockNumber.Valid {
		return uint64(h.L1BlockNumber.Int64), nil
	}

	return uint64(h.Number), nil
}

// LatestRoundRequested returns the configDigest, epoch, and round from the latest
// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
// return a result if the latest such event was emitted in a block b such that
// b.timestamp < tip.timestamp - lookback.
//
// If no event is found, LatestRoundRequested should return zero values, not an error.
// An error should only be returned if an actual error occurred during execution,
// e.g. because there was an error querying the blockchain or the database.
//
// As an optimization, this function may also return zero values, if no
// RoundRequested event has been emitted after the latest NewTransmission event.
func (t *ContractTracker) LatestRoundRequested(_ context.Context, lookback time.Duration) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, err error) {
	t.lrrMu.RLock()
	defer t.lrrMu.RUnlock()
	return t.latestRoundRequested.ConfigDigest, t.latestRoundRequested.Epoch, t.latestRoundRequested.Round, nil
}

func getEventTopic(name string) gethCommon.Hash {
	abi, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	event, exists := abi.Events[name]
	if !exists {
		panic(fmt.Sprintf("abi.Events was missing %s", name))
	}
	return event.ID
}
