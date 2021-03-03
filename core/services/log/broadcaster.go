package log

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Broadcaster --output ./mocks/ --case=underscore --structname Broadcaster --filename broadcaster.go
//go:generate mockery --name Listener --output ./mocks/ --case=underscore --structname Listener --filename listener.go
//go:generate mockery --name AbigenContract --output ./mocks/ --case=underscore --structname AbigenContract --filename abigen_contract.go

type (
	// The Broadcaster manages log subscription requests for the Chainlink node.  Instead
	// of creating a new subscription for each request, it multiplexes all subscriptions
	// to all of the relevant contracts over a single connection and forwards the logs to the
	// relevant subscribers.
	Broadcaster interface {
		utils.DependentAwaiter
		Start() error
		Stop() error
		Register(contract AbigenContract, listener Listener) (connected bool)
		Unregister(contract AbigenContract, listener Listener)
	}

	broadcaster struct {
		subscriber *subscriber
		relayer    *relayer
		utils.StartStopOnce
		utils.DependentAwaiter
	}

	// The Listener responds to log events through HandleLog, and contains setup/tear-down
	// callbacks in the On* functions.
	Listener interface {
		OnConnect()
		OnDisconnect()
		HandleLog(lb Broadcast, err error)
		JobID() models.JobID
		JobIDV2() int32
		IsV2Job() bool
	}

	AbigenContract interface {
		Address() common.Address
		ParseLog(log types.Log) (interface{}, error)
	}

	Config interface {
		BlockBackfillDepth() uint64
		TriggerFallbackDBPollInterval() time.Duration
	}
)

var _ Broadcaster = (*broadcaster)(nil)

// NewBroadcaster creates a new instance of the broadcaster
func NewBroadcaster(orm ORM, ethClient eth.Client, config Config) *broadcaster {
	var (
		dependentAwaiter = utils.NewDependentAwaiter()
		relayer          = newRelayer(orm, config, dependentAwaiter)
		subscriber       = newSubscriber(orm, ethClient, config, relayer, dependentAwaiter)
	)
	return &broadcaster{
		subscriber:       subscriber,
		relayer:          relayer,
		DependentAwaiter: dependentAwaiter,
	}
}

func (b *broadcaster) Start() error {
	return b.StartOnce("Log broadcaster", func() error {
		err := b.subscriber.Start()
		if err != nil {
			return err
		}
		return b.relayer.Start()
	})
}

func (b *broadcaster) Stop() error {
	return b.StopOnce("Log broadcaster", func() error {
		err := b.subscriber.Stop()
		if err != nil {
			return err
		}
		return b.relayer.Stop()
	})
}

func (b *broadcaster) Register(contract AbigenContract, listener Listener) (connected bool) {
	b.subscriber.NotifyAddContract(contract.Address())
	b.relayer.NotifyAddListener(contract, listener)
	return b.subscriber.IsConnected()
}

func (b *broadcaster) Unregister(contract AbigenContract, listener Listener) {
	b.subscriber.NotifyRemoveContract(contract.Address())
	b.relayer.NotifyRemoveListener(contract, listener)
}

func (b *broadcaster) OnNewLongestChain(ctx context.Context, head models.Head) {
	b.relayer.OnNewLongestChain(ctx, head)
}

func (b *broadcaster) Connect(head *models.Head) error { return nil }
func (b *broadcaster) Disconnect()                     {}

// ListenerJobID returns the appropriate job ID for a listener
func ListenerJobID(listener Listener) interface{} {
	if listener.IsV2Job() {
		return listener.JobIDV2()
	}
	return listener.JobID()
}
