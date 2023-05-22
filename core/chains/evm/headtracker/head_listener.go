package headtracker

import (
	"context"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

type headListener[
	H commontypes.Head[BLOCK_HASH],
	S commontypes.Subscription,
	ID commontypes.ID,
	BLOCK_HASH commontypes.Hashable,
	CLIENT htrktypes.Client[H, S, ID, BLOCK_HASH],
] struct {
	config           htrktypes.Config
	client           CLIENT
	logger           logger.Logger
	chStop           utils.StopChan
	chHeaders        chan H
	headSubscription commontypes.Subscription
	connected        atomic.Bool
	receivingHeads   atomic.Bool
	getNilHead       func() H
}

type evmHeadListener = headListener[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash, evmclient.Client]

// NewHeadListener creates a new HeadListener for EVM chains
func NewHeadListener(lggr logger.Logger, ethClient evmclient.Client, config Config, chStop chan struct{}) *evmHeadListener {
	return &evmHeadListener{
		config: NewWrappedConfig(config),
		client: ethClient,
		logger: lggr.Named("HeadListener"),
		chStop: chStop,
		getNilHead: func() *evmtypes.Head {
			return nil
		},
	}
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) Name() string {
	return hl.logger.Name()
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) ListenForNewHeads(handleNewHead commontypes.NewHeadHandler[H, BLOCK_HASH], done func()) {
	defer done()
	defer hl.unsubscribe()

	ctx, cancel := hl.chStop.NewCtx()
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

func (hl *headListener[H, S, ID, BLOCK_HAS, CLIENT]) ReceivingHeads() bool {
	return hl.receivingHeads.Load()
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) Connected() bool {
	return hl.connected.Load()
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) HealthReport() map[string]error {
	var err error
	if !hl.ReceivingHeads() {
		err = errors.New("Listener is not receiving heads")
	}
	if !hl.Connected() {
		err = errors.New("Listener is not connected")
	}
	return map[string]error{hl.Name(): err}
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) receiveHeaders(ctx context.Context, handleNewHead commontypes.NewHeadHandler[H, BLOCK_HASH]) error {
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
			if blockHeader.Equals(hl.getNilHead()) {
				hl.logger.Error("got nil block header")
				continue
			}
			if blockHeader.EVMChainID == nil || !utils.NewBig(chainId).Equal(blockHeader.EVMChainID) {
				hl.logger.Panicf("head listener for %s received block header for %s", chainId, blockHeader.EVMChainID)
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
			hl.logger.Warnf("have not received a head for %v", noHeadsAlarmDuration)
			hl.receivingHeads.Store(false)
		}
	}
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) subscribe(ctx context.Context) bool {
	subscribeRetryBackoff := utils.NewRedialBackoff()

	chainId := hl.client.ConfiguredChainID()

	for {
		hl.unsubscribe()

		hl.logger.Debugf("Subscribing to new heads on chain %s", chainId.String())

		select {
		case <-hl.chStop:
			return false

		case <-time.After(subscribeRetryBackoff.Duration()):
			err := hl.subscribeToHead(ctx)
			if err != nil {
				promEthConnectionErrors.WithLabelValues(chainId.String()).Inc()
				hl.logger.Warnw("Failed to subscribe to heads on chain", "chainID", chainId.String(), "err", err)
			} else {
				hl.logger.Debugf("Subscribed to heads on chain %s", chainId.String())
				return true
			}
		}
	}
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) subscribeToHead(ctx context.Context) error {
	hl.chHeaders = make(chan H)

	var err error
	hl.headSubscription, err = hl.client.SubscribeNewHead(ctx, hl.chHeaders)
	if err != nil {
		close(hl.chHeaders)
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
	}

	hl.connected.Store(true)

	return nil
}

func (hl *headListener[H, S, ID, BLOCK_HASH, CLIENT]) unsubscribe() {
	if hl.headSubscription != nil {
		hl.connected.Store(false)
		hl.headSubscription.Unsubscribe()
		hl.headSubscription = nil
	}
}
