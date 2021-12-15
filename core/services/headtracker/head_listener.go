package headtracker

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
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

//go:generate mockery --name HeadListener --output ./mocks/ --case=underscore

// NewHeadHandler is a callback that handles incoming heads
type NewHeadHandler func(ctx context.Context, header *eth.Head) error

// HeadListener manages eth.Client connection that receives heads from the eth node
type HeadListener interface {
	// ListenForNewHeads kicks off the listen loop (not thread safe)
	// done() must be executed upon leaving ListenForNewHeads()
	ListenForNewHeads(handleNewHead NewHeadHandler, done func())
	// ReceivingHeads returns true if the listener is receiving heads (thread safe)
	ReceivingHeads() bool
	// Connected returns true if the listener is connected (thread safe)
	Connected() bool
}

type headListener struct {
	config           Config
	ethClient        eth.Client
	logger           logger.Logger
	stopCh           chan struct{}
	headersCh        chan *eth.Head
	headSubscription ethereum.Subscription
	connected        atomic.Bool
	receivingHeads   atomic.Bool
}

// NewHeadListener creates a new HeadListener
func NewHeadListener(logger logger.Logger, ethClient eth.Client, config Config, stopCh chan struct{}) HeadListener {
	return &headListener{
		config:    config,
		ethClient: ethClient,
		logger:    logger.Named("listener"),
		stopCh:    stopCh,
	}
}

func (hl *headListener) ListenForNewHeads(handleNewHead NewHeadHandler, done func()) {
	defer done()
	defer hl.unsubscribe()

	ctx, cancel := utils.ContextFromChan(hl.stopCh)
	defer cancel()

	for {
		if !hl.subscribe(ctx) {
			break
		}
		err := hl.receiveHeaders(ctx, handleNewHead)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			hl.logger.Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
			continue
		} else {
			break
		}
	}
}

func (hl *headListener) ReceivingHeads() bool {
	return hl.receivingHeads.Load()
}

func (hl *headListener) Connected() bool {
	return hl.connected.Load()
}

// This should be safe to run concurrently across multiple nodes connected to the same database
// Note: returning nil from receiveHeaders will cause ListenForNewHeads to exit completely
func (hl *headListener) receiveHeaders(ctx context.Context, handleNewHead NewHeadHandler) error {
	noHeadsAlarmDuration := hl.config.BlockEmissionIdleWarningThreshold()
	t := time.NewTicker(noHeadsAlarmDuration)

	for {
		select {
		case <-hl.stopCh:
			return nil

		case blockHeader, open := <-hl.headersCh:
			// We've received a head, reset the no heads alarm
			t.Stop()
			t = time.NewTicker(noHeadsAlarmDuration)
			hl.receivingHeads.Store(true)
			if !open {
				return errors.New("head listener: headersCh prematurely closed")
			}
			if blockHeader == nil {
				hl.logger.Error("got nil block header")
				continue
			}
			if blockHeader.EVMChainID == nil || !utils.NewBig(hl.ethClient.ChainID()).Equal(blockHeader.EVMChainID) {
				hl.logger.Panicf("head listener for %s received block header for %s", hl.ethClient.ChainID(), blockHeader.EVMChainID)
			}
			promNumHeadsReceived.WithLabelValues(hl.ethClient.ChainID().String()).Inc()

			err := handleNewHead(ctx, blockHeader)
			if ctx.Err() != nil {
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
			hl.logger.Warn(fmt.Sprintf("have not received a head for %v", noHeadsAlarmDuration))
			hl.receivingHeads.Store(false)
		}
	}
}

func (hl *headListener) subscribe(ctx context.Context) bool {
	subscribeRetryBackoff := backoff.Backoff{
		Min: 1 * time.Second,
		Max: 10 * time.Second,
	}

	chainID := hl.ethClient.ChainID().String()

	for {
		hl.unsubscribe()

		hl.logger.Debugf("Subscribing to new heads on chain %s", chainID)

		select {
		case <-hl.stopCh:
			return false

		case <-time.After(subscribeRetryBackoff.Duration()):
			err := hl.subscribeToHead(ctx)
			if err != nil {
				promEthConnectionErrors.WithLabelValues(hl.ethClient.ChainID().String()).Inc()
				hl.logger.Warnw(fmt.Sprintf("Failed to subscribe to heads on chain %s", chainID), "err", err)
			} else {
				hl.logger.Debugf("Subscribed to heads on chain %s", chainID)
				return true
			}
		}
	}
}

func (hl *headListener) subscribeToHead(ctx context.Context) error {
	hl.headersCh = make(chan *eth.Head)

	var err error
	hl.headSubscription, err = hl.ethClient.SubscribeNewHead(ctx, hl.headersCh)
	if err != nil {
		close(hl.headersCh)
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
	}

	hl.connected.Store(true)

	return nil
}

func (hl *headListener) unsubscribe() {
	if hl.headSubscription != nil {
		hl.connected.Store(false)
		hl.headSubscription.Unsubscribe()
		hl.headSubscription = nil
	}
}
