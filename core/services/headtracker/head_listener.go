package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promNumHeadsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
	})
	promEthConnectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_eth_connection_errors",
		Help: "The total number of eth node connection errors",
	})
)

type Config interface {
	ChainID() *big.Int
	EvmHeadTrackerHistoryDepth() uint
	EvmHeadTrackerMaxBufferSize() uint
	EvmHeadTrackerSamplingInterval() time.Duration
	BlockEmissionIdleWarningThreshold() time.Duration
	EthereumURL() string
	EvmFinalityDepth() uint
}

type HeadListener struct {
	config           Config
	ethClient        eth.Client
	headers          chan *models.Head
	headSubscription ethereum.Subscription
	connectedMutex   sync.RWMutex
	connected        bool
	receivesHeads    int32
	sleeper          utils.Sleeper

	log      *logger.Logger
	muLogger sync.RWMutex

	chStop chan struct{}
	wgDone *sync.WaitGroup
}

func NewHeadListener(l *logger.Logger,
	ethClient eth.Client,
	config Config,
	chStop chan struct{},
	wgDone *sync.WaitGroup,
	sleepers ...utils.Sleeper,
) *HeadListener {
	var sleeper utils.Sleeper
	if len(sleepers) > 0 {
		sleeper = sleepers[0]
	} else {
		sleeper = utils.NewBackoffSleeper()
	}
	return &HeadListener{
		config:    config,
		ethClient: ethClient,
		sleeper:   sleeper,
		log:       l,
		chStop:    chStop,
		wgDone:    wgDone,
	}
}

// SetLogger sets and reconfigures the log for the head tracker service
func (hl *HeadListener) SetLogger(logger *logger.Logger) {
	hl.muLogger.Lock()
	defer hl.muLogger.Unlock()
	hl.log = logger
}

func (hl *HeadListener) logger() *logger.Logger {
	hl.muLogger.RLock()
	defer hl.muLogger.RUnlock()
	return hl.log
}

func (hl *HeadListener) ListenForNewHeads(handleNewHead func(ctx context.Context, header models.Head) error) {
	defer hl.wgDone.Done()
	defer func() {
		if err := hl.unsubscribeFromHead(); err != nil {
			hl.logger().Warn(errors.Wrap(err, "HeadListener failed when unsubscribe from head"))
		}
	}()

	ctx, cancel := utils.ContextFromChan(hl.chStop)
	defer cancel()

	for {
		if !hl.subscribe() {
			break
		}
		err := hl.receiveHeaders(ctx, handleNewHead)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			hl.logger().Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
			hl.headers = nil
			continue
		} else {
			break
		}
	}
}

// This should be safe to run concurrently across multiple nodes connected to the same database
// Note: returning nil from receiveHeaders will cause listenForNewHeads to exit completely
func (hl *HeadListener) receiveHeaders(ctx context.Context, handleNewHead func(ctx context.Context, header models.Head) error) error {
	noHeadsAlarmDuration := hl.config.BlockEmissionIdleWarningThreshold()
	t := time.NewTicker(noHeadsAlarmDuration)

	for {
		select {
		case <-hl.chStop:
			return nil
		case blockHeader, open := <-hl.headers:
			// We've received a head, reset the no heads alarm
			t.Stop()
			t = time.NewTicker(noHeadsAlarmDuration)
			atomic.StoreInt32(&hl.receivesHeads, 1)
			if !open {
				return errors.New("HeadTracker: headers prematurely closed")
			}
			if blockHeader == nil {
				hl.logger().Error("HeadTracker: got nil block header")
				continue
			}
			promNumHeadsReceived.Inc()

			err := handleNewHead(ctx, *blockHeader)
			if ctx.Err() != nil {
				// the 'ctx' context is closed only on ht.done - on shutdown, so it's safe to return nil
				return nil
			} else if err != nil {
				return err
			}

		case err, open := <-hl.headSubscription.Err():
			if open && err != nil {
				return err
			}

		case <-t.C:
			// We haven't received a head on the channel for a long time, log a warning
			logger.Warn(fmt.Sprintf("HeadTracker: have not received a head for %v", noHeadsAlarmDuration))
			atomic.StoreInt32(&hl.receivesHeads, 0)
		}
	}
}

// subscribe periodically attempts to connect to the ethereum node via websocket.
// It returns true on success, and false if cut short by a done request and did not connect.
func (hl *HeadListener) subscribe() bool {
	hl.sleeper.Reset()
	for {
		if err := hl.unsubscribeFromHead(); err != nil {
			hl.logger().Error("failed when unsubscribe from head", err)
			return false
		}

		hl.logger().Info("HeadListener: Connecting to ethereum node ", hl.config.EthereumURL(), " in ", hl.sleeper.Duration())
		select {
		case <-hl.chStop:
			return false
		case <-time.After(hl.sleeper.After()):
			err := hl.subscribeToHead()
			if err != nil {
				promEthConnectionErrors.Inc()
				hl.logger().Warnw(fmt.Sprintf("HeadListener: Failed to connect to ethereum node %v", hl.config.EthereumURL()), "err", err)
			} else {
				hl.logger().Info("HeadListener: Connected to ethereum node ", hl.config.EthereumURL())
				return true
			}
		}
	}
}

func (hl *HeadListener) subscribeToHead() error {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	hl.headers = make(chan *models.Head)

	sub, err := hl.ethClient.SubscribeNewHead(context.Background(), hl.headers)
	if err != nil {
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
	}
	if err := verifyEthereumChainID(hl); err != nil {
		return errors.Wrap(err, "verifyEthereumChainID failed")
	}

	hl.headSubscription = sub
	hl.connected = true

	return nil
}

func (hl *HeadListener) unsubscribeFromHead() error {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	if !hl.connected {
		return nil
	}

	hl.headSubscription.Unsubscribe()
	hl.connected = false

	// ht.headers will be nil if subscription failed, channel closed, and
	// receiveHeaders returned from the loop. listenForNewHeads will set it to
	// nil in that case to avoid a double close panic.
	if hl.headers != nil {
		close(hl.headers)
	}
	return nil
}

// Connected returns whether or not this HeadTracker is connected.
func (hl *HeadListener) Connected() bool {
	hl.connectedMutex.RLock()
	defer hl.connectedMutex.RUnlock()

	return hl.connected
}

// chainIDVerify checks whether or not the ChainID from the Chainlink config
// matches the ChainID reported by the ETH node connected to this Chainlink node.
func verifyEthereumChainID(ht *HeadListener) error {
	ethereumChainID, err := ht.ethClient.ChainID(context.Background())
	if err != nil {
		return err
	}

	if ethereumChainID.Cmp(ht.config.ChainID()) != 0 {
		return fmt.Errorf(
			"ethereum ChainID doesn't match chainlink config.ChainID: config ID=%d, eth RPC ID=%d",
			ht.config.ChainID(),
			ethereumChainID,
		)
	}
	return nil
}
