package headtracker

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
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

type headListener struct {
	config           Config
	ethClient        evmclient.Client
	logger           logger.Logger
	chStop           chan struct{}
	chHeaders        chan *evmtypes.Head
	headSubscription ethereum.Subscription
	connected        atomic.Bool
	receivingHeads   atomic.Bool
}

// NewHeadListener creates a new HeadListener
func NewHeadListener(lggr logger.Logger, ethClient evmclient.Client, config Config, chStop chan struct{}) httypes.HeadListener {
	return &headListener{
		config:    config,
		ethClient: ethClient,
		logger:    lggr.Named("HeadListener"),
		chStop:    chStop,
	}
}

func (hl *headListener) ListenForNewHeads(handleNewHead httypes.NewHeadHandler, done func()) {
	defer done()
	defer hl.unsubscribe()

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
			hl.logger.Errorw("Error in new head subscription, unsubscribed", "err", err)
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

func (hl *headListener) receiveHeaders(ctx context.Context, handleNewHead httypes.NewHeadHandler) error {
	var noHeadsAlarmC <-chan time.Time
	var noHeadsAlarmT *time.Ticker
	noHeadsAlarmDuration := hl.config.BlockEmissionIdleWarningThreshold()
	if noHeadsAlarmDuration > 0 {
		noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
		noHeadsAlarmC = noHeadsAlarmT.C
	}

	for {
		select {
		case <-hl.chStop:
			return nil

		case blockHeader, open := <-hl.chHeaders:
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
			// err can be nil, because of using chainIDSubForwarder
			if !open || err == nil {
				return errors.New("head listener: subscription Err channel prematurely closed")
			}
			return err

		case <-noHeadsAlarmC:
			// We haven't received a head on the channel for a long time, log a warning
			hl.logger.Warnf("have not received a head for %v", noHeadsAlarmDuration)
			hl.receivingHeads.Store(false)
		}
	}
}

func (hl *headListener) subscribe(ctx context.Context) bool {
	subscribeRetryBackoff := utils.NewRedialBackoff()

	chainID := hl.ethClient.ChainID().String()

	for {
		hl.unsubscribe()

		hl.logger.Debugf("Subscribing to new heads on chain %s", chainID)

		select {
		case <-hl.chStop:
			return false

		case <-time.After(subscribeRetryBackoff.Duration()):
			err := hl.subscribeToHead(ctx)
			if err != nil {
				promEthConnectionErrors.WithLabelValues(hl.ethClient.ChainID().String()).Inc()
				hl.logger.Warnw("Failed to subscribe to heads on chain", "chainID", chainID, "err", err)
			} else {
				hl.logger.Debugf("Subscribed to heads on chain %s", chainID)
				return true
			}
		}
	}
}

func (hl *headListener) subscribeToHead(ctx context.Context) error {
	hl.chHeaders = make(chan *evmtypes.Head)

	var err error
	hl.headSubscription, err = hl.ethClient.SubscribeNewHead(ctx, hl.chHeaders)
	if err != nil {
		close(hl.chHeaders)
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
