package blockfetcher

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// HeadSubscription - Managing fetching of data for a single head
type HeadSubscription interface {
	Head() models.Head

	// Block - will deliver the block data
	Block() <-chan types.Block

	// Receipts - will deliver receipt data
	Receipts() <-chan []types.Receipt

	// Unsubscribe - stops inflow of data for this head
	Unsubscribe()
}

// Backend - Managing fetching of heads and block data
type Backend interface {

	// Subscribe - Start receiving newest heads and data for them
	Subscribe() <-chan HeadSubscription

	HeadByNumber(n *big.Int) HeadSubscription
}

type Config interface {
	EthHeadTrackerMaxBufferSize() int

	ChainID() *Int
}

type DefaultBackend struct {
	logger                    *logger.Logger
	config                    Config
	ethClient                 eth.Client
	headTrackerAddressChannel chan common.Address
	inHeaders                 chan *models.Head
	outHeaders                chan models.Head
	headSubscription          ethereum.Subscription
	highestSeenHead           *models.Head
	headMutex                 sync.RWMutex
	connected                 bool
	sleeper                   utils.Sleeper
	done                      chan struct{}
	started                   bool
	listenForNewHeadsWg       sync.WaitGroup
	backfillMB                utils.Mailbox
	subscriptionSucceeded     chan struct{}
	muLogger                  sync.RWMutex
}

func newDefaultBackend(l *logger.Logger,
	ethClient eth.Client,
	config Config,
	headTrackerAddressChannel chan common.Address,
) *DefaultBackend {
	var sleeper utils.Sleeper
	return &DefaultBackend{
		config:                    config,
		headTrackerAddressChannel: headTrackerAddressChannel,
		ethClient:                 ethClient,
		sleeper:                   sleeper,
		logger:                    l,
		backfillMB:                *utils.NewMailbox(1),
		done:                      make(chan struct{}),
	}
}

func (ht *DefaultBackend) subscribeToHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ht.inHeaders = make(chan *models.Head)
	var rb *headRingBuffer
	rb, ht.outHeaders = newHeadRingBuffer(ht.inHeaders, int(ht.config.EthHeadTrackerMaxBufferSize()), func() *logger.Logger { return ht.logger })
	// It will autostop when we close inHeaders channel
	rb.Start()

	sub, err := ht.ethClient.SubscribeNewHead(context.Background(), ht.inHeaders)
	if err != nil {
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
	}

	if err := verifyEthereumChainID(ht); err != nil {
		return errors.Wrap(err, "verifyEthereumChainID failed")
	}

	ht.headSubscription = sub
	ht.connected = true

	ht.connect(ht.highestSeenHead)
	return nil
}

func (ht *DefaultBackend) unsubscribeFromHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	if !ht.connected {
		return nil
	}

	timedUnsubscribe(ht.headSubscription)

	ht.connected = false
	ht.disconnect()
	close(ht.inHeaders)
	// Drain channel and wait for ringbuffer to close it
	for range ht.outHeaders {
	}
	return nil
}

// chainIDVerify checks whether or not the ChainID from the Chainlink config
// matches the ChainID reported by the ETH node connected to this Chainlink node.
func verifyEthereumChainID(ht *DefaultBackend) error {
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
			r.logger().Error("HeadTracker: got nil block header")
			continue
		}
		//promNumHeadsReceived.Inc()
		hInQueue := len(r.out)
		//promHeadsInQueue.Set(float64(hInQueue))
		if hInQueue > 0 {
			r.logger().Infof("HeadTracker: Head %v is lagging behind, there are %v more heads in the queue. Your node is operating close to its maximum capacity and may start to miss jobs.", h.Number, hInQueue)
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
				//promNumHeadsDropped.Inc()
				r.logger().Errorf("HeadTracker: dropping head %v with hash 0x%x because queue is full. WARNING: Your node is overloaded and may start missing jobs.", dropped.Number, h.Hash)
				r.out <- *h
			default:
				r.out <- *h
			}
		}
	}
	close(r.out)
}
