package ocr

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

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

// configMailboxSanityLimit is the maximum number of configs that can be held
// in the mailbox. Under normal operation there should never be more than 0 or
// 1 configs in the mailbox, this limit is here merely to prevent unbounded usage
// in some kind of unforeseen insane situation.
const configMailboxSanityLimit = 100

var (
	_ ocrtypes.ContractConfigTracker = &OCRContractTracker{}
	_ log.Listener                   = &OCRContractTracker{}
	_ httypes.HeadTrackable          = &OCRContractTracker{}

	OCRContractConfigSet            = getEventTopic("ConfigSet")
	OCRContractLatestRoundRequested = getEventTopic("RoundRequested")
)

//go:generate mockery --quiet --name OCRContractTrackerDB --output ./mocks/ --case=underscore
type (
	// OCRContractTracker complies with ContractConfigTracker interface and
	// handles log events related to the contract more generally
	OCRContractTracker struct {
		services.StateMachine

		ethClient        evmclient.Client
		contract         *offchain_aggregator_wrapper.OffchainAggregator
		contractFilterer *offchainaggregator.OffchainAggregatorFilterer
		contractCaller   *offchainaggregator.OffchainAggregatorCaller
		logBroadcaster   log.Broadcaster
		jobID            int32
		logger           logger.Logger
		ocrDB            OCRContractTrackerDB
		ds               sqlutil.DataSource
		blockTranslator  ocrcommon.BlockTranslator
		cfg              ocrcommon.Config
		mailMon          *mailbox.Monitor

		// HeadBroadcaster
		headBroadcaster  httypes.HeadBroadcaster
		unsubscribeHeads func()

		// Start/Stop lifecycle
		chStop          services.StopChan
		wg              sync.WaitGroup
		unsubscribeLogs func()

		// LatestRoundRequested
		latestRoundRequested offchainaggregator.OffchainAggregatorRoundRequested
		lrrMu                sync.RWMutex

		// ContractConfig
		configsMB *mailbox.Mailbox[ocrtypes.ContractConfig]
		chConfigs chan ocrtypes.ContractConfig

		// LatestBlockHeight
		latestBlockHeight   int64
		latestBlockHeightMu sync.RWMutex
	}

	OCRContractTrackerDB interface {
		SaveLatestRoundRequested(ctx context.Context, rr offchainaggregator.OffchainAggregatorRoundRequested) error
		LoadLatestRoundRequested(ctx context.Context) (rr offchainaggregator.OffchainAggregatorRoundRequested, err error)
		WithDataSource(sqlutil.DataSource) OCRContractTrackerDB
	}
)

func (t *OCRContractTracker) HealthReport() map[string]error {
	return map[string]error{t.Name(): t.Healthy()}
}

func (t *OCRContractTracker) Name() string { return t.logger.Name() }

// NewOCRContractTracker makes a new OCRContractTracker
func NewOCRContractTracker(
	contract *offchain_aggregator_wrapper.OffchainAggregator,
	contractFilterer *offchainaggregator.OffchainAggregatorFilterer,
	contractCaller *offchainaggregator.OffchainAggregatorCaller,
	ethClient evmclient.Client,
	logBroadcaster log.Broadcaster,
	jobID int32,
	logger logger.Logger,
	ds sqlutil.DataSource,
	ocrDB OCRContractTrackerDB,
	cfg ocrcommon.Config,
	headBroadcaster httypes.HeadBroadcaster,
	mailMon *mailbox.Monitor,
) (o *OCRContractTracker) {
	logger = logger.Named("OCRContractTracker")
	return &OCRContractTracker{
		ethClient:            ethClient,
		contract:             contract,
		contractFilterer:     contractFilterer,
		contractCaller:       contractCaller,
		logBroadcaster:       logBroadcaster,
		jobID:                jobID,
		logger:               logger,
		ocrDB:                ocrDB,
		ds:                   ds,
		blockTranslator:      ocrcommon.NewBlockTranslator(cfg, ethClient, logger),
		cfg:                  cfg,
		mailMon:              mailMon,
		headBroadcaster:      headBroadcaster,
		chStop:               make(services.StopChan),
		latestRoundRequested: offchainaggregator.OffchainAggregatorRoundRequested{},
		configsMB:            mailbox.New[ocrtypes.ContractConfig](configMailboxSanityLimit),
		chConfigs:            make(chan ocrtypes.ContractConfig),
		latestBlockHeight:    -1,
	}
}

// Start must be called before logs can be delivered
// It ought to be called before starting OCR
func (t *OCRContractTracker) Start(ctx context.Context) error {
	return t.StartOnce("OCRContractTracker", func() (err error) {
		t.latestRoundRequested, err = t.ocrDB.LoadLatestRoundRequested(ctx)
		if err != nil {
			return errors.Wrap(err, "OCRContractTracker#Start: failed to load latest round requested")
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
			t.setLatestBlockHeight(latestHead)
		}

		t.wg.Add(1)
		go t.processLogs()

		t.mailMon.Monitor(t.configsMB, "OCRContractTracker", "Configs", fmt.Sprint(t.jobID))

		return nil
	})
}

// Close should be called after teardown of the OCR job relying on this tracker
func (t *OCRContractTracker) Close() error {
	return t.StopOnce("OCRContractTracker", func() error {
		close(t.chStop)
		t.wg.Wait()
		t.unsubscribeHeads()
		t.unsubscribeLogs()
		close(t.chConfigs)
		return t.configsMB.Close()
	})
}

// OnNewLongestChain conformed to HeadTrackable and updates latestBlockHeight
func (t *OCRContractTracker) OnNewLongestChain(_ context.Context, h *evmtypes.Head) {
	t.setLatestBlockHeight(h)
}

func (t *OCRContractTracker) setLatestBlockHeight(h *evmtypes.Head) {
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

func (t *OCRContractTracker) getLatestBlockHeight() int64 {
	t.latestBlockHeightMu.RLock()
	defer t.latestBlockHeightMu.RUnlock()
	return t.latestBlockHeight
}

func (t *OCRContractTracker) processLogs() {
	defer t.wg.Done()
	for {
		select {
		case <-t.configsMB.Notify():
			// NOTE: libocr could take an arbitrary amount of time to process a
			// new chaintype. To avoid blocking the log broadcaster, we use this
			// background thread to deliver them and a mailbox as the buffer.
			for {
				cc, exists := t.configsMB.Retrieve()
				if !exists {
					break
				}
				select {
				case t.chConfigs <- cc:
				case <-t.chStop:
					return
				}
			}
		case <-t.chStop:
			return
		}
	}
}

// HandleLog complies with LogListener interface
// It is not thread safe
func (t *OCRContractTracker) HandleLog(ctx context.Context, lb log.Broadcast) {
	was, err := t.logBroadcaster.WasAlreadyConsumed(ctx, lb)
	if err != nil {
		t.logger.Errorw("could not determine if log was already consumed", "err", err)
		return
	} else if was {
		return
	}

	raw := lb.RawLog()
	if raw.Address != t.contract.Address() {
		t.logger.Errorf("log address of 0x%x does not match configured contract address of 0x%x", raw.Address, t.contract.Address())
		if err2 := t.logBroadcaster.MarkConsumed(ctx, nil, lb); err2 != nil {
			t.logger.Errorw("failed to mark log consumed", "err", err2)
		}
		return
	}
	topics := raw.Topics
	if len(topics) == 0 {
		if err2 := t.logBroadcaster.MarkConsumed(ctx, nil, lb); err2 != nil {
			t.logger.Errorw("failed to mark log consumed", "err", err2)
		}
		return
	}

	var consumed bool
	switch topics[0] {
	case OCRContractConfigSet:
		var configSet *offchainaggregator.OffchainAggregatorConfigSet
		configSet, err = t.contractFilterer.ParseConfigSet(raw)
		if err != nil {
			t.logger.Errorw("could not parse config set", "err", err)
			if err2 := t.logBroadcaster.MarkConsumed(ctx, nil, lb); err2 != nil {
				t.logger.Errorw("failed to mark log consumed", "err", err2)
			}
			return
		}
		configSet.Raw = lb.RawLog()
		cc := confighelper.ContractConfigFromConfigSetEvent(*configSet)

		wasOverCapacity := t.configsMB.Deliver(cc)
		if wasOverCapacity {
			t.logger.Error("config mailbox is over capacity - dropped the oldest unprocessed item")
		}
	case OCRContractLatestRoundRequested:
		var rr *offchainaggregator.OffchainAggregatorRoundRequested
		rr, err = t.contractFilterer.ParseRoundRequested(raw)
		if err != nil {
			t.logger.Errorw("could not parse round requested", "err", err)
			if err2 := t.logBroadcaster.MarkConsumed(ctx, nil, lb); err2 != nil {
				t.logger.Errorw("failed to mark log consumed", "err", err2)
			}
			return
		}
		if IsLaterThan(raw, t.latestRoundRequested.Raw) {
			err = sqlutil.TransactDataSource(ctx, t.ds, nil, func(tx sqlutil.DataSource) error {
				if err = t.ocrDB.WithDataSource(tx).SaveLatestRoundRequested(ctx, *rr); err != nil {
					return err
				}
				return t.logBroadcaster.MarkConsumed(ctx, tx, lb)
			})
			if err != nil {
				t.logger.Error(err)
				return
			}
			consumed = true
			t.lrrMu.Lock()
			t.latestRoundRequested = *rr
			t.lrrMu.Unlock()
			t.logger.Infow("received new latest RoundRequested event", "latestRoundRequested", *rr)
		} else {
			t.logger.Warnw("ignoring out of date RoundRequested event", "latestRoundRequested", t.latestRoundRequested, "roundRequested", rr)
		}
	default:
		t.logger.Debugw("got unrecognised log topic", "topic", topics[0])
	}
	if !consumed {
		if err := t.logBroadcaster.MarkConsumed(ctx, nil, lb); err != nil {
			t.logger.Errorw("failed to mark log consumed", "err", err)
		}
	}
}

// IsLaterThan returns true if the first log was emitted "after" the second log
// from the blockchain's point of view
func IsLaterThan(incoming gethTypes.Log, existing gethTypes.Log) bool {
	return incoming.BlockNumber > existing.BlockNumber ||
		(incoming.BlockNumber == existing.BlockNumber && incoming.TxIndex > existing.TxIndex) ||
		(incoming.BlockNumber == existing.BlockNumber && incoming.TxIndex == existing.TxIndex && incoming.Index > existing.Index)
}

// JobID complies with LogListener interface
func (t *OCRContractTracker) JobID() int32 {
	return t.jobID
}

// SubscribeToNewConfigs returns the tracker aliased as a ContractConfigSubscription
func (t *OCRContractTracker) SubscribeToNewConfigs(context.Context) (ocrtypes.ContractConfigSubscription, error) {
	return (*OCRContractConfigSubscription)(t), nil
}

// LatestConfigDetails queries the eth node
func (t *OCRContractTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	var cancel context.CancelFunc
	ctx, cancel = t.chStop.Ctx(ctx)
	defer cancel()

	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := t.contractCaller.LatestConfigDetails(&opts)
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting LatestConfigDetails")
	}
	configDigest, err = ocrtypes.BytesToConfigDigest(result.ConfigDigest[:])
	if err != nil {
		return 0, configDigest, errors.Wrap(err, "error getting config digest")
	}
	return uint64(result.BlockNumber), configDigest, err
}

// ConfigFromLogs queries the eth node for logs for this contract
func (t *OCRContractTracker) ConfigFromLogs(ctx context.Context, changedInBlock uint64) (c ocrtypes.ContractConfig, err error) {
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
	ctx, cancel = t.chStop.Ctx(ctx)
	defer cancel()

	logs, err := t.ethClient.FilterLogs(ctx, q)
	if err != nil {
		return c, err
	}
	if len(logs) == 0 {
		return c, errors.Errorf("ConfigFromLogs: OCRContract with address 0x%x has no logs", t.contract.Address())
	}

	latest, err := t.contractFilterer.ParseConfigSet(logs[len(logs)-1])
	if err != nil {
		return c, errors.Wrap(err, "ConfigFromLogs failed to ParseConfigSet")
	}
	latest.Raw = logs[len(logs)-1]
	if latest.Raw.Address != t.contract.Address() {
		return c, errors.Errorf("log address of 0x%x does not match configured contract address of 0x%x", latest.Raw.Address, t.contract.Address())
	}
	return confighelper.ContractConfigFromConfigSetEvent(*latest), err
}

// LatestBlockHeight queries the eth node for the most recent header
func (t *OCRContractTracker) LatestBlockHeight(ctx context.Context) (blockheight uint64, err error) {
	switch t.cfg.ChainType() {
	case chaintype.ChainMetis:
		// We skip confirmation checking anyway on these L2s so there's no need to
		// care about the block height; we have no way of getting the L1 block
		// height anyway
		return 0, nil
	case "", chaintype.ChainArbitrum, chaintype.ChainCelo, chaintype.ChainGnosis, chaintype.ChainKroma, chaintype.ChainOptimismBedrock, chaintype.ChainScroll, chaintype.ChainWeMix, chaintype.ChainXLayer, chaintype.ChainZkEvm, chaintype.ChainZkSync:
		// continue
	}
	latestBlockHeight := t.getLatestBlockHeight()
	if latestBlockHeight >= 0 {
		return uint64(latestBlockHeight), nil
	}

	t.logger.Debugw("still waiting for first head, falling back to on-chain lookup")

	var cancel context.CancelFunc
	ctx, cancel = t.chStop.Ctx(ctx)
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
func (t *OCRContractTracker) LatestRoundRequested(_ context.Context, lookback time.Duration) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, err error) {
	// NOTE: This should be "good enough" 99% of the time.
	// It guarantees validity up to `EVM.BlockBackfillDepth` blocks ago
	// Some further improvements could be made:
	t.lrrMu.RLock()
	defer t.lrrMu.RUnlock()
	return t.latestRoundRequested.ConfigDigest, t.latestRoundRequested.Epoch, t.latestRoundRequested.Round, nil
}

func getEventTopic(name string) gethCommon.Hash {
	abi, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	event, exists := abi.Events[name]
	if !exists {
		panic(fmt.Sprintf("abi.Events was missing %s", name))
	}
	return event.ID
}
