package headtracker

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/common/internal/utils"
	"github.com/smartcontractkit/chainlink/v2/common/types"
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

// HeadHandler is a callback that handles incoming heads
type HeadHandler[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] func(ctx context.Context, header H) error

// HeadListener is a chain agnostic interface that manages connection of Client that receives heads from the blockchain node
type HeadListener[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	services.Service

	// ListenForNewHeads runs the listen loop (not thread safe)
	ListenForNewHeads(ctx context.Context)

	// ReceivingHeads returns true if the listener is receiving heads (thread safe)
	ReceivingHeads() bool

	// Connected returns true if the listener is connected (thread safe)
	Connected() bool

	// HealthReport returns report of errors within HeadListener
	HealthReport() map[string]error
}

type headListener[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	services.Service
	eng *services.Engine

	config           htrktypes.Config
	client           htrktypes.Client[HTH, S, ID, BLOCK_HASH]
	onSubscription   func(context.Context)
	handleNewHead    HeadHandler[HTH, BLOCK_HASH]
	chHeaders        chan HTH
	headSubscription types.Subscription
	connected        atomic.Bool
	receivingHeads   atomic.Bool
}

func NewHeadListener[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
	CLIENT htrktypes.Client[HTH, S, ID, BLOCK_HASH],
](
	lggr logger.Logger,
	client CLIENT,
	config htrktypes.Config,
	onSubscription func(context.Context),
	handleNewHead HeadHandler[HTH, BLOCK_HASH],
) HeadListener[HTH, BLOCK_HASH] {
	hl := &headListener[HTH, S, ID, BLOCK_HASH]{
		config:         config,
		client:         client,
		onSubscription: onSubscription,
		handleNewHead:  handleNewHead,
	}
	hl.Service, hl.eng = services.Config{
		Name:  "HeadListener",
		Start: hl.start,
	}.NewServiceEngine(lggr)
	return hl
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) start(context.Context) error {
	hl.eng.Go(hl.ListenForNewHeads)
	return nil
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) ListenForNewHeads(ctx context.Context) {
	defer hl.unsubscribe()

	for {
		if !hl.subscribe(ctx) {
			break
		}

		if hl.onSubscription != nil {
			hl.onSubscription(ctx)
		}
		err := hl.receiveHeaders(ctx, hl.handleNewHead)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			hl.eng.Errorw("Error in new head subscription, unsubscribed", "err", err)
			continue
		}
		break
	}
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) ReceivingHeads() bool {
	return hl.receivingHeads.Load()
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) Connected() bool {
	return hl.connected.Load()
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) HealthReport() map[string]error {
	var err error
	if !hl.ReceivingHeads() {
		err = errors.New("Listener is not receiving heads")
	}
	if !hl.Connected() {
		err = errors.New("Listener is not connected")
	}
	return map[string]error{hl.Name(): err}
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) receiveHeaders(ctx context.Context, handleNewHead HeadHandler[HTH, BLOCK_HASH]) error {
	var noHeadsAlarmC <-chan time.Time
	var noHeadsAlarmT *time.Ticker
	noHeadsAlarmDuration := hl.config.BlockEmissionIdleWarningThreshold()
	if noHeadsAlarmDuration > 0 {
		noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
		noHeadsAlarmC = noHeadsAlarmT.C
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case blockHeader, open := <-hl.chHeaders:
			chainId := hl.client.ConfiguredChainID()
			if noHeadsAlarmT != nil {
				// We've received a head, reset the no heads alarm
				noHeadsAlarmT.Stop()
				noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
				noHeadsAlarmC = noHeadsAlarmT.C
			}
			hl.receivingHeads.Store(true)
			if !open {
				return errors.New("head listener: chHeaders prematurely closed")
			}
			if !blockHeader.IsValid() {
				hl.eng.Error("got nil block header")
				continue
			}

			// Compare the chain ID of the block header to the chain ID of the client
			if !blockHeader.HasChainID() || blockHeader.ChainID().String() != chainId.String() {
				hl.eng.Panicf("head listener for %s received block header for %s", chainId, blockHeader.ChainID())
			}
			promNumHeadsReceived.WithLabelValues(chainId.String()).Inc()

			err := handleNewHead(ctx, blockHeader)
			if ctx.Err() != nil {
				return nil
			} else if err != nil {
				return err
			}

		case err, open := <-hl.headSubscription.Err():
			// err can be nil, because of using chainIDSubForwarder
			if !open || err == nil {
				return errors.New("head listener: subscription Err channel prematurely closed")
			}
			return err

		case <-noHeadsAlarmC:
			// We haven't received a head on the channel for a long time, log a warning
			hl.eng.Warnf("have not received a head for %v", noHeadsAlarmDuration)
			hl.receivingHeads.Store(false)
		}
	}
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) subscribe(ctx context.Context) bool {
	subscribeRetryBackoff := utils.NewRedialBackoff()

	chainId := hl.client.ConfiguredChainID()

	for {
		hl.unsubscribe()

		hl.eng.Debugf("Subscribing to new heads on chain %s", chainId.String())

		select {
		case <-ctx.Done():
			return false

		case <-time.After(subscribeRetryBackoff.Duration()):
			err := hl.subscribeToHead(ctx)
			if err != nil {
				promEthConnectionErrors.WithLabelValues(chainId.String()).Inc()
				hl.eng.Warnw("Failed to subscribe to heads on chain", "chainID", chainId.String(), "err", err)
			} else {
				hl.eng.Debugf("Subscribed to heads on chain %s", chainId.String())
				return true
			}
		}
	}
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) subscribeToHead(ctx context.Context) error {
	hl.chHeaders = make(chan HTH)

	var err error
	hl.headSubscription, err = hl.client.SubscribeNewHead(ctx, hl.chHeaders)
	if err != nil {
		close(hl.chHeaders)
		return fmt.Errorf("Client#SubscribeNewHead: %w", err)
	}

	hl.connected.Store(true)

	return nil
}

func (hl *headListener[HTH, S, ID, BLOCK_HASH]) unsubscribe() {
	if hl.headSubscription != nil {
		hl.connected.Store(false)
		hl.headSubscription.Unsubscribe()
		hl.headSubscription = nil
	}
}
