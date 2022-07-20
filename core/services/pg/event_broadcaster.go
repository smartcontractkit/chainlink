package pg

import (
	"context"
	"database/sql"
	"net/url"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name EventBroadcaster --output ./mocks/ --case=underscore
//go:generate mockery --name Subscription --output ./mocks/ --case=underscore

// EventBroadcaster opaquely manages a collection of Postgres event listeners
// and broadcasts events to subscribers (with an optional payload filter).
type EventBroadcaster interface {
	services.ServiceCtx
	Subscribe(channel, payloadFilter string) (Subscription, error)
	Notify(channel string, payload string) error
}

type eventBroadcaster struct {
	uri                  string
	minReconnectInterval time.Duration
	maxReconnectDuration time.Duration
	db                   *sql.DB
	listener             *pq.Listener
	subscriptions        map[string]map[Subscription]struct{}
	subscriptionsMu      sync.RWMutex
	chStop               chan struct{}
	chDone               chan struct{}
	lggr                 logger.Logger
	utils.StartStopOnce
}

var _ EventBroadcaster = (*eventBroadcaster)(nil)

type Event struct {
	Channel string
	Payload string
}

func NewEventBroadcaster(uri url.URL, minReconnectInterval time.Duration, maxReconnectDuration time.Duration, lggr logger.Logger, appID uuid.UUID) *eventBroadcaster {
	if minReconnectInterval == time.Duration(0) {
		minReconnectInterval = 1 * time.Second
	}
	if maxReconnectDuration == time.Duration(0) {
		maxReconnectDuration = 1 * time.Minute
	}
	static.SetConsumerName(&uri, "EventBroadcaster", &appID)
	return &eventBroadcaster{
		uri:                  uri.String(),
		minReconnectInterval: minReconnectInterval,
		maxReconnectDuration: maxReconnectDuration,
		subscriptions:        make(map[string]map[Subscription]struct{}),
		chStop:               make(chan struct{}),
		chDone:               make(chan struct{}),
		lggr:                 lggr.Named("EventBroadcaster"),
	}
}

// Start starts EventBroadcaster.
func (b *eventBroadcaster) Start(context.Context) error {
	return b.StartOnce("Postgres event broadcaster", func() (err error) {
		// Explicitly using the lib/pq for notifications so we use the postgres driverName
		// and NOT pgx.
		db, err := sql.Open("postgres", b.uri)
		if err != nil {
			return err
		}
		b.db = db
		b.listener = pq.NewListener(b.uri, b.minReconnectInterval, b.maxReconnectDuration, func(ev pq.ListenerEventType, err error) {
			// These are always connection-related events, and the pq library
			// automatically handles reconnecting to the DB. Therefore, we do not
			// need to terminate, but rather simply log these events for node
			// operators' sanity.
			switch ev {
			case pq.ListenerEventConnected:
				b.lggr.Debug("Postgres event broadcaster: connected")
			case pq.ListenerEventDisconnected:
				b.lggr.Warnw("Postgres event broadcaster: disconnected, trying to reconnect...", "error", err)
			case pq.ListenerEventReconnected:
				b.lggr.Debug("Postgres event broadcaster: reconnected")
			case pq.ListenerEventConnectionAttemptFailed:
				b.lggr.Warnw("Postgres event broadcaster: reconnect attempt failed, trying again...", "error", err)
			}
		})

		go b.runLoop()
		return nil
	})
}

// Stop permanently destroys the EventBroadcaster.  Calling this does not clean
// up any outstanding subscriptions.  Subscribers must explicitly call `.Close()`
// or they will leak goroutines.
func (b *eventBroadcaster) Close() error {
	return b.StopOnce("Postgres event broadcaster", func() (err error) {
		b.subscriptionsMu.RLock()
		defer b.subscriptionsMu.RUnlock()
		b.subscriptions = nil

		err = multierr.Append(err, b.db.Close())
		err = multierr.Append(err, b.listener.Close())
		close(b.chStop)
		<-b.chDone
		return err
	})
}

func (b *eventBroadcaster) runLoop() {
	defer close(b.chDone)
	for {
		select {
		case <-b.chStop:
			return

		case notification, open := <-b.listener.NotificationChannel():
			if !open {
				return
			} else if notification == nil {
				continue
			}
			b.lggr.Debugw("Postgres event broadcaster: received notification",
				"channel", notification.Channel,
				"payload", notification.Extra,
			)
			b.broadcast(notification)
		}
	}
}

func (b *eventBroadcaster) Notify(channel string, payload string) error {
	_, err := b.db.Exec(`SELECT pg_notify($1, $2)`, channel, payload)
	return errors.Wrap(err, "Postgres event broadcaster could not notify")
}

func (b *eventBroadcaster) Subscribe(channel, payloadFilter string) (Subscription, error) {
	b.subscriptionsMu.Lock()
	defer b.subscriptionsMu.Unlock()

	if _, exists := b.subscriptions[channel]; !exists {
		err := b.listener.Listen(channel)
		if err != nil {
			return nil, errors.Wrap(err, "Postgres event broadcaster could not subscribe")
		}
		b.subscriptions[channel] = make(map[Subscription]struct{})
	}

	sub := &subscription{
		channel:          channel,
		payloadFilter:    payloadFilter,
		eventBroadcaster: b,
		queue:            utils.NewBoundedQueue[Event](1000),
		chEvents:         make(chan Event),
		chDone:           make(chan struct{}),
		lggr:             logger.Sugared(b.lggr),
	}
	sub.processQueueWorker = utils.NewSleeperTask(
		utils.SleeperFuncTask(sub.processQueue, "SubscriptionQueueProcessor"),
	)
	b.subscriptions[channel][sub] = struct{}{}
	return sub, nil
}

func (b *eventBroadcaster) removeSubscription(sub Subscription) {
	b.subscriptionsMu.Lock()
	defer b.subscriptionsMu.Unlock()

	// The following conditions can occur on shutdown when .Stop() is called
	// before one or more subscriptions' .Close() methods are called
	if b.subscriptions == nil {
		return
	}
	subs, exists := b.subscriptions[sub.ChannelName()]
	if !exists || subs == nil {
		return
	}

	delete(b.subscriptions[sub.ChannelName()], sub)
	if len(b.subscriptions[sub.ChannelName()]) == 0 {
		err := b.listener.Unlisten(sub.ChannelName())
		if err != nil {
			b.lggr.Errorw("Postgres event broadcaster: failed to unsubscribe", "error", err)
		}
		delete(b.subscriptions, sub.ChannelName())
	}
}

func (b *eventBroadcaster) broadcast(notification *pq.Notification) {
	b.subscriptionsMu.RLock()
	defer b.subscriptionsMu.RUnlock()

	event := Event{
		Channel: notification.Channel,
		Payload: notification.Extra,
	}

	var wg sync.WaitGroup
	for sub := range b.subscriptions[event.Channel] {
		if sub.InterestedIn(event) {
			wg.Add(1)
			go func(sub Subscription) {
				defer wg.Done()
				sub.Send(event)
			}(sub)
		}
	}
	wg.Wait()
}

// Subscription represents a subscription to a Postgres event channel
type Subscription interface {
	Events() <-chan Event
	Close()

	ChannelName() string
	InterestedIn(event Event) bool
	Send(event Event)
}

type subscription struct {
	channel            string
	payloadFilter      string
	eventBroadcaster   *eventBroadcaster
	queue              *utils.BoundedQueue[Event]
	processQueueWorker utils.SleeperTask
	chEvents           chan Event
	chDone             chan struct{}
	lggr               logger.SugaredLogger
}

var _ Subscription = (*subscription)(nil)

func (sub *subscription) InterestedIn(event Event) bool {
	return sub.payloadFilter == event.Payload || sub.payloadFilter == ""
}

func (sub *subscription) Send(event Event) {
	sub.queue.Add(event)
	sub.processQueueWorker.WakeUpIfStarted()
}

const broadcastTimeout = 10 * time.Second

func (sub *subscription) processQueue() {
	ctx, cancel := context.WithTimeout(context.Background(), broadcastTimeout)
	defer cancel()

	for !sub.queue.Empty() {
		event := sub.queue.Take()
		select {
		case sub.chEvents <- event:
		case <-ctx.Done():
		case <-sub.chDone:
		}
	}
}

func (sub *subscription) Events() <-chan Event {
	return sub.chEvents
}

func (sub *subscription) ChannelName() string {
	return sub.channel
}

func (sub *subscription) Close() {
	sub.eventBroadcaster.removeSubscription(sub)
	// Close chDone before stopping the SleeperTask to avoid deadlocks
	close(sub.chDone)
	err := sub.processQueueWorker.Stop()
	if err != nil {
		sub.lggr.Errorw("THIS NEVER RETURNS AN ERROR", "error", err)
	}
}

// NullEventBroadcaster implements null pattern for event broadcaster
type NullEventBroadcaster struct {
	Sub *NullSubscription
}

func NewNullEventBroadcaster() *NullEventBroadcaster {
	sub := &NullSubscription{make(chan (Event))}
	return &NullEventBroadcaster{sub}
}

var _ EventBroadcaster = &NullEventBroadcaster{}

// Start does no-op.
func (*NullEventBroadcaster) Start(context.Context) error { return nil }

// Close does no-op.
func (*NullEventBroadcaster) Close() error { return nil }

// Ready does no-op.
func (*NullEventBroadcaster) Ready() error { return nil }

// Healthy does no-op.
func (*NullEventBroadcaster) Healthy() error { return nil }

func (ne *NullEventBroadcaster) Subscribe(channel, payloadFilter string) (Subscription, error) {
	return ne.Sub, nil
}
func (*NullEventBroadcaster) Notify(channel string, payload string) error { return nil }

var _ Subscription = &NullSubscription{}

type NullSubscription struct {
	Ch chan (Event)
}

func (ns *NullSubscription) Events() <-chan Event          { return ns.Ch }
func (ns *NullSubscription) Close()                        {}
func (ns *NullSubscription) ChannelName() string           { return "" }
func (ns *NullSubscription) InterestedIn(event Event) bool { return false }
func (ns *NullSubscription) Send(event Event)              {}
