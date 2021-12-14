package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/jpillora/backoff"
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

// Config configures headtracker and related structs
type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	EvmFinalityDepth() uint32
	EvmHeadTrackerHistoryDepth() uint32
	EvmHeadTrackerMaxBufferSize() uint32
	EvmHeadTrackerSamplingInterval() time.Duration
}

// HeadListener manages the websocket connection that receives heads from the
// eth node
type HeadListener struct {
	config                Config
	ethClient             eth.Client
	chainID               big.Int
	headers               chan *eth.Head
	headSubscription      ethereum.Subscription
	connectedMutex        sync.RWMutex
	connected             bool
	receivesHeads         atomic.Bool
	subscribeRetryBackoff backoff.Backoff

	log logger.Logger

	chStop chan struct{}
}

// NewHeadListener creates a new HeadListener
func NewHeadListener(l logger.Logger,
	ethClient eth.Client,
	config Config,
	chStop chan struct{},
) *HeadListener {
	if ethClient == nil {
		panic("head listener requires non-nil ethclient")
	}
	return &HeadListener{
		config:    config,
		ethClient: ethClient,
		chainID:   *ethClient.ChainID(),
		subscribeRetryBackoff: backoff.Backoff{
			Min: 1 * time.Second,
			Max: 10 * time.Second,
		},
		log:    l.Named("listener"),
		chStop: chStop,
	}
}

// ListenForNewHeads kicks off the listen loop
// NOT THREAD SAFE
func (hl *HeadListener) ListenForNewHeads(handleNewHead func(ctx context.Context, header *eth.Head) error, done func()) {
	defer done()
	defer hl.unsubscribeFromHead()

	ctx, cancel := utils.ContextFromChan(hl.chStop)
	defer cancel()

	for {
		if !hl.subscribe(ctx) {
			break
		}
		err := hl.receiveHeaders(ctx, handleNewHead)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			hl.log.Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
			continue
		} else {
			break
		}
	}
}

// This should be safe to run concurrently across multiple nodes connected to the same database
// Note: returning nil from receiveHeaders will cause listenForNewHeads to exit completely
func (hl *HeadListener) receiveHeaders(ctx context.Context, handleNewHead func(ctx context.Context, header *eth.Head) error) error {
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
			if blockHeader.EVMChainID == nil || !utils.NewBig(&hl.chainID).Equal(blockHeader.EVMChainID) {
				panic(fmt.Sprintf("head listener for %s received block header for %s", &hl.chainID, blockHeader.EVMChainID))
			}
			promNumHeadsReceived.WithLabelValues(hl.chainID.String()).Inc()

			err := handleNewHead(ctx, blockHeader)
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
// NOT THREAD SAFE
func (hl *HeadListener) subscribe(ctx context.Context) bool {
	hl.subscribeRetryBackoff.Reset()

	for {
		hl.unsubscribeFromHead()

		hl.log.Debugf("Subscribing to new heads on chain %s", hl.chainID.String())
		select {
		case <-hl.chStop:
			return false
		case <-time.After(hl.subscribeRetryBackoff.Duration()):
			err := hl.subscribeToHead(ctx)
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

func (hl *HeadListener) subscribeToHead(ctx context.Context) error {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	hl.headers = make(chan *eth.Head)

	sub, err := hl.ethClient.SubscribeNewHead(ctx, hl.headers)
	if err != nil {
		close(hl.headers)
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
	}

	hl.headSubscription = sub
	hl.connected = true

	return nil
}

func (hl *HeadListener) unsubscribeFromHead() {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	if !hl.connected {
		return
	}

	hl.headSubscription.Unsubscribe()
	hl.connected = false
}

// Connected returns whether or not this HeadTracker is connected.
func (hl *HeadListener) Connected() bool {
	hl.connectedMutex.RLock()
	defer hl.connectedMutex.RUnlock()

	return hl.connected
}
