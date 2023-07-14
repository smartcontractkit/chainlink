package headmanager

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	hmtypes "github.com/smartcontractkit/chainlink/v2/common/headmanager/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	promNumHeadsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
	}, []string{"ChainID"})
	promEthConnectionErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_connection_errors",
		Help: "The total number of node connection errors",
	}, []string{"ChainID"})
)

type Listener[
	HMH hmtypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	config           hmtypes.Config
	client           hmtypes.Client[HMH, S, ID, BLOCK_HASH]
	logger           logger.Logger
	chStop           utils.StopChan
	chHeaders        chan HMH
	headSubscription types.Subscription
	connected        atomic.Bool
	receivingHeads   atomic.Bool
}

func NewListener[
	HMH hmtypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
	CLIENT hmtypes.Client[HMH, S, ID, BLOCK_HASH],
](
	lggr logger.Logger,
	client CLIENT,
	config hmtypes.Config,
	chStop chan struct{},
) *Listener[HMH, S, ID, BLOCK_HASH] {
	return &Listener[HMH, S, ID, BLOCK_HASH]{
		config: config,
		client: client,
		logger: lggr.Named("Listener"),
		chStop: chStop,
	}
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) Name() string {
	return l.logger.Name()
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) ListenForNewHeads(handleNewHead types.NewHeadHandler[HMH, BLOCK_HASH], done func()) {
	defer done()
	defer l.unsubscribe()

	ctx, cancel := l.chStop.NewCtx()
	defer cancel()

	for {
		if !l.subscribe(ctx) {
			break
		}
		err := l.receiveHeaders(ctx, handleNewHead)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			l.logger.Errorw("Error in new head subscription, unsubscribed", "err", err)
			continue
		} else {
			break
		}
	}
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) ReceivingHeads() bool {
	return l.receivingHeads.Load()
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) Connected() bool {
	return l.connected.Load()
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) HealthReport() map[string]error {
	var err error
	if !l.ReceivingHeads() {
		err = errors.New("Listener is not receiving heads")
	}
	if !l.Connected() {
		err = errors.New("Listener is not connected")
	}
	return map[string]error{l.Name(): err}
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) receiveHeaders(ctx context.Context, handleNewHead types.NewHeadHandler[HMH, BLOCK_HASH]) error {
	var noHeadsAlarmC <-chan time.Time
	var noHeadsAlarmT *time.Ticker
	noHeadsAlarmDuration := l.config.BlockEmissionIdleWarningThreshold()
	if noHeadsAlarmDuration > 0 {
		noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
		noHeadsAlarmC = noHeadsAlarmT.C
	}

	for {
		select {
		case <-l.chStop:
			return nil

		case blockHeader, open := <-l.chHeaders:
			chainId := l.client.ConfiguredChainID()
			if noHeadsAlarmT != nil {
				// We've received a head, reset the no heads alarm
				noHeadsAlarmT.Stop()
				noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
				noHeadsAlarmC = noHeadsAlarmT.C
			}
			l.receivingHeads.Store(true)
			if !open {
				return errors.New("head listener: chHeaders prematurely closed")
			}
			if !blockHeader.IsValid() {
				l.logger.Error("got nil block header")
				continue
			}

			// Compare the chain ID of the block header to the chain ID of the client
			if !blockHeader.HasChainID() || blockHeader.ChainID().String() != chainId.String() {
				l.logger.Panicf("head listener for %s received block header for %s", chainId, blockHeader.ChainID())
			}
			promNumHeadsReceived.WithLabelValues(chainId.String()).Inc()

			err := handleNewHead(ctx, blockHeader)
			if ctx.Err() != nil {
				return nil
			} else if err != nil {
				return err
			}

		case err, open := <-l.headSubscription.Err():
			// err can be nil, because of using chainIDSubForwarder
			if !open || err == nil {
				return errors.New("head listener: subscription Err channel prematurely closed")
			}
			return err

		case <-noHeadsAlarmC:
			// We haven't received a head on the channel for a long time, log a warning
			l.logger.Warnf("have not received a head for %v", noHeadsAlarmDuration)
			l.receivingHeads.Store(false)
		}
	}
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) subscribe(ctx context.Context) bool {
	subscribeRetryBackoff := utils.NewRedialBackoff()

	chainId := l.client.ConfiguredChainID()

	for {
		l.unsubscribe()

		l.logger.Debugf("Subscribing to new heads on chain %s", chainId.String())

		select {
		case <-l.chStop:
			return false

		case <-time.After(subscribeRetryBackoff.Duration()):
			err := l.subscribeToHead(ctx)
			if err != nil {
				promEthConnectionErrors.WithLabelValues(chainId.String()).Inc()
				l.logger.Warnw("Failed to subscribe to heads on chain", "chainID", chainId.String(), "err", err)
			} else {
				l.logger.Debugf("Subscribed to heads on chain %s", chainId.String())
				return true
			}
		}
	}
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) subscribeToHead(ctx context.Context) error {
	l.chHeaders = make(chan HMH)

	var err error
	l.headSubscription, err = l.client.SubscribeNewHead(ctx, l.chHeaders)
	if err != nil {
		close(l.chHeaders)
		return errors.Wrap(err, "Client#SubscribeNewHead")
	}

	l.connected.Store(true)

	return nil
}

func (l *Listener[HMH, S, ID, BLOCK_HASH]) unsubscribe() {
	if l.headSubscription != nil {
		l.connected.Store(false)
		l.headSubscription.Unsubscribe()
		l.headSubscription = nil
	}
}
