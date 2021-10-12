package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/atomic"
)

var (
	promNumHeadsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
	}, []string{"evmChainID"})
	promEthConnectionErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_eth_connection_errors",
		Help: "The total number of eth node connection errors",
	}, []string{"evmChainID"})
)

//go:generate mockery --name Config --output ./mocks/ --case=underscore
type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	EvmFinalityDepth() uint32
	EvmHeadTrackerHistoryDepth() uint32
	EvmHeadTrackerMaxBufferSize() uint32
	EvmHeadTrackerSamplingInterval() time.Duration
}

type HeadListener struct {
	config           Config
	ethClient        eth.Client
	chainID          big.Int
	headers          chan *eth.Head
	headSubscription ethereum.Subscription
	connectedMutex   sync.RWMutex
	connected        bool
	receivesHeads    atomic.Bool
	sleeper          utils.Sleeper

	log logger.Logger

	chStop chan struct{}
}

func NewHeadListener(l logger.Logger,
	ethClient eth.Client,
	config Config,
	chStop chan struct{},
	sleepers ...utils.Sleeper,
) *HeadListener {
	if ethClient == nil {
		panic("head listener requires non-nil ethclient")
	}
	var sleeper utils.Sleeper
	if len(sleepers) > 0 {
		sleeper = sleepers[0]
	} else {
		sleeper = utils.NewBackoffSleeper()
	}
	return &HeadListener{
		config:    config,
		ethClient: ethClient,
		chainID:   *ethClient.ChainID(),
		sleeper:   sleeper,
		log:       l.Named("listener"),
		chStop:    chStop,
	}
}

func (hl *HeadListener) ListenForNewHeads(handleNewHead func(ctx context.Context, header eth.Head) error, done func()) {
	defer done()
	defer func() {
		if err := hl.unsubscribeFromHead(); err != nil {
			hl.log.Warn(errors.Wrap(err, "HeadListener failed when unsubscribe from head"))
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
			hl.log.Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
			hl.headers = nil
			continue
		} else {
			break
		}
	}
}

// This should be safe to run concurrently across multiple nodes connected to the same database
// Note: returning nil from receiveHeaders will cause listenForNewHeads to exit completely
func (hl *HeadListener) receiveHeaders(ctx context.Context, handleNewHead func(ctx context.Context, header eth.Head) error) error {
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
			hl.receivesHeads.Store(true)
			if !open {
				return errors.New("HeadTracker: headers prematurely closed")
			}
			if blockHeader == nil {
				hl.log.Error("got nil block header")
				continue
			}
			blockHeader.EVMChainID = utils.NewBig(&hl.chainID)
			promNumHeadsReceived.WithLabelValues(hl.chainID.String()).Inc()

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
			hl.log.Warn(fmt.Sprintf("have not received a head for %v", noHeadsAlarmDuration))
			hl.receivesHeads.Store(false)
		}
	}
}

// subscribe periodically attempts to connect to the ethereum node via websocket.
// It returns true on success, and false if cut short by a done request and did not connect.
func (hl *HeadListener) subscribe() bool {
	hl.sleeper.Reset()
	for {
		if err := hl.unsubscribeFromHead(); err != nil {
			hl.log.Error("failed when unsubscribe from head", err)
			return false
		}

		hl.log.Debugf("Subscribing to new heads on chain %s (in %s)", hl.chainID.String(), hl.sleeper.Duration())
		select {
		case <-hl.chStop:
			return false
		case <-time.After(hl.sleeper.After()):
			err := hl.subscribeToHead()
			if err != nil {
				promEthConnectionErrors.WithLabelValues(hl.chainID.String()).Inc()
				hl.log.Warnw(fmt.Sprintf("Failed to subscribe to heads on chain %s", hl.chainID.String()), "err", err)
			} else {
				hl.log.Debugf("Subscribed to heads on chain %s", hl.chainID.String())
				return true
			}
		}
	}
}

func (hl *HeadListener) subscribeToHead() error {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	hl.headers = make(chan *eth.Head)

	sub, err := hl.ethClient.SubscribeNewHead(context.Background(), hl.headers)
	if err != nil {
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
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
