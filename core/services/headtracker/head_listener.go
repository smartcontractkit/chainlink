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
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promNumHeadsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_listener_heads_received",
		Help: "The total number of heads seen",
	})
	promHeadsInQueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "head_listener_heads_in_queue",
		Help: "The number of heads currently waiting to be executed. You can think of this as the 'load' on the head tracker. Should rarely or never be more than 0",
	})
	promNumHeadsDropped = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_listener_num_heads_dropped",
		Help: "The total number of heads dropped",
	})
	promEthConnectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_listener_eth_connection_errors",
		Help: "The total number of eth node connection errors",
	})

	promHeadTimeoutsCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_listener_head_timeouts_count",
		Help: "The total number of receipts fetched",
	})
)

type Config interface {
	EthHeadTrackerMaxBufferSize() uint
	EthereumURL() string
	HeadTimeBudget() time.Duration
	ChainID() *big.Int
}

type HeadListener struct {
	logger           *logger.Logger
	config           Config
	ethClient        eth.Client
	inHeaders        chan *models.Head
	outHeaders       chan models.Head
	headSubscription ethereum.Subscription
	connectedMutex   sync.RWMutex
	connected        bool
	sleeper          utils.Sleeper

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
		logger:    l,
		chStop:    chStop,
		wgDone:    wgDone,
	}
}

func (hl *HeadListener) ListenForNewHeads(handleNewHead func(ctx context.Context, header models.Head) error, connected func()) {
	defer hl.wgDone.Done()
	defer func() {
		if err := hl.unsubscribeFromHead(); err != nil {
			hl.logger.Warn(errors.Wrap(err, "HeadListener failed when unsubscribe from head"))
		}
	}()

	ctx, cancel := utils.ContextFromChan(hl.chStop)
	defer cancel()

	for {
		if !hl.subscribe(connected) {
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

// This should be safe to run concurrently across multiple nodes connected to the same database
// Note: returning nil from receiveHeaders will cause listenForNewHeads to exit completely
func (hl *HeadListener) receiveHeaders(ctx context.Context, handleNewHead func(ctx context.Context, header models.Head) error) error {
	for {
		select {
		case <-hl.chStop:
			return nil
		case blockHeader, open := <-hl.outHeaders:
			if !open {
				return errors.New("HeadListener: outHeaders prematurely closed")
			}
			timeBudget := hl.config.HeadTimeBudget()
			{
				deadlineCtx, cancel := context.WithTimeout(ctx, timeBudget)
				defer cancel()

				err := handleNewHead(ctx, blockHeader)
				if ctx.Err() != nil {
					// the 'ctx' context is closed only on hl.done - on shutdown, so it's safe to return nil
					return nil
				} else if deadlineCtx.Err() != nil {
					promHeadTimeoutsCount.Inc()
					logger.Warnw("HeadListener: handling of new head timed out", "error", ctx.Err(), "timeBudget", timeBudget.String())
					continue
				} else if err != nil {
					return err
				}
			}
		case err, open := <-hl.headSubscription.Err():
			if open && err != nil {
				return err
			}
		}
	}
}

// subscribe periodically attempts to connect to the ethereum node via websocket.
// It returns true on success, and false if cut short by a done request and did not connect.
func (hl *HeadListener) subscribe(connected func()) bool {
	hl.sleeper.Reset()
	for {
		if err := hl.unsubscribeFromHead(); err != nil {
			hl.logger.Error("failed when unsubscribe from head", err)
			return false
		}

		hl.logger.Info("HeadListener: Connecting to ethereum node ", hl.config.EthereumURL(), " in ", hl.sleeper.Duration())
		select {
		case <-hl.chStop:
			return false
		case <-time.After(hl.sleeper.After()):
			err := hl.subscribeToHead(connected)
			if err != nil {
				promEthConnectionErrors.Inc()
				hl.logger.Warnw(fmt.Sprintf("HeadListener: Failed to connect to ethereum node %v", hl.config.EthereumURL()), "err", err)
			} else {
				hl.logger.Info("HeadListener: Connected to ethereum node ", hl.config.EthereumURL())
				return true
			}
		}
	}
}

func (hl *HeadListener) subscribeToHead(connected func()) error {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	hl.inHeaders = make(chan *models.Head)
	var rb *headRingBuffer
	rb, hl.outHeaders = newHeadRingBuffer(hl.inHeaders, int(hl.config.EthHeadTrackerMaxBufferSize()), func() *logger.Logger { return hl.logger })
	// It will autostop when we close inHeaders channel
	rb.Start()

	sub, err := hl.ethClient.SubscribeNewHead(context.Background(), hl.inHeaders)
	if err != nil {
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
	}

	if err := verifyEthereumChainID(hl); err != nil {
		return errors.Wrap(err, "verifyEthereumChainID failed")
	}

	hl.headSubscription = sub
	hl.connected = true

	connected()
	return nil
}

func (hl *HeadListener) unsubscribeFromHead() error {
	hl.connectedMutex.Lock()
	defer hl.connectedMutex.Unlock()

	if !hl.connected {
		return nil
	}

	timedUnsubscribe(hl.headSubscription)

	hl.connected = false
	close(hl.inHeaders)
	// Drain channel and wait for ringbuffer to close it
	for range hl.outHeaders {
	}
	return nil
}

// timedUnsubscribe attempts to unsubscribe but aborts abruptly after a time delay
// unblocking the application. This is an effort to mitigate the occasional
// indefinite block described here from go-ethereum:
// https://chainlink/pull/600#issuecomment-426320971
func timedUnsubscribe(unsubscriber ethereum.Subscription) {
	unsubscribed := make(chan struct{})
	go func() {
		unsubscriber.Unsubscribe()
		close(unsubscribed)
	}()
	select {
	case <-unsubscribed:
	case <-time.After(100 * time.Millisecond):
		logger.Warnf("Subscription %T Unsubscribe timed out.", unsubscriber)
	}
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

// headRingBuffer is a small goroutine that sits between the eth client and the
// head tracker and drops the oldest head if necessary in order to keep to a fixed
// queue size (defined by the buffer size of out channel)
type headRingBuffer struct {
	in     <-chan *models.Head
	out    chan models.Head
	start  sync.Once
	logger func() *logger.Logger
}

func newHeadRingBuffer(in <-chan *models.Head, size int, logger func() *logger.Logger) (r *headRingBuffer, out chan models.Head) {
	out = make(chan models.Head, size)
	return &headRingBuffer{
		in:     in,
		out:    out,
		start:  sync.Once{},
		logger: logger,
	}, out
}

// Start the headRingBuffer goroutine
// It will be stopped implicitly by closing the in channel
func (r *headRingBuffer) Start() {
	r.start.Do(func() {
		go r.run()
	})
}

func (r *headRingBuffer) run() {
	for h := range r.in {
		if h == nil {
			r.logger().Error("HeadListener: got nil block header")
			continue
		}
		promNumHeadsReceived.Inc()
		hInQueue := len(r.out)
		promHeadsInQueue.Set(float64(hInQueue))
		if hInQueue > 0 {
			r.logger().Infof("HeadListener: Head %v is lagging behind, there are %v more heads in the queue. Your node is operating close to its maximum capacity and may start to miss jobs.", h.Number, hInQueue)
		}
		select {
		case r.out <- *h:
		default:
			// Need to select/default here because it's conceivable (although
			// improbable) that between the previous select and now, all heads were drained
			// from r.out by another goroutine
			//
			// NOTE: In this unlikely event, we may drop an extra head unnecessarily.
			// The probability of this seems vanishingly small, and only hits
			// if the queue was already full anyway, so we can live with this
			select {
			case dropped := <-r.out:
				promNumHeadsDropped.Inc()
				r.logger().Errorf("HeadListener: dropping head %v with hash 0x%x because queue is full. WARNING: Your node is overloaded and may start missing jobs.", dropped.Number, h.Hash)
				r.out <- *h
			default:
				r.out <- *h
			}
		}
	}
	close(r.out)
}
