package postgres

import (
	"context"
	"database/sql"
	"net/url"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name EventBroadcaster --output ./mocks/ --case=underscore
//go:generate mockery --name Subscription --output ./mocks/ --case=underscore

// EventBroadcaster opaquely manages a collection of Postgres event listeners
// and broadcasts events to subscribers (with an optional payload filter).
type EventBroadcaster interface {
	Start() error
	Stop() error
	Subscribe(channel, payloadFilter string) (Subscription, error)
	Notify(channel string, payload string) error
	NotifyInsideGormTx(tx *gorm.DB, channel string, payload string) error
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
	utils.StartStopOnce
}

var _ EventBroadcaster = (*eventBroadcaster)(nil)

type Event struct {
	Channel string
	Payload string
}

func NewEventBroadcaster(uri url.URL, minReconnectInterval time.Duration, maxReconnectDuration time.Duration) *eventBroadcaster {
	if minReconnectInterval == time.Duration(0) {
		minReconnectInterval = 1 * time.Second
	}
	if maxReconnectDuration == time.Duration(0) {
		maxReconnectDuration = 1 * time.Minute
	}
	static.SetConsumerName(&uri, "EventBroadcaster")
	return &eventBroadcaster{
		uri:                  uri.String(),
		minReconnectInterval: minReconnectInterval,
		maxReconnectDuration: maxReconnectDuration,
		subscriptions:        make(map[string]map[Subscription]struct{}),
		chStop:               make(chan struct{}),
		chDone:               make(chan struct{}),
	}
}

func (b *eventBroadcaster) Start() error {
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
				logger.Debug("Postgres event broadcaster: connected")
			case pq.ListenerEventDisconnected:
				logger.Warnw("Postgres event broadcaster: disconnected, trying to reconnect...", "error", err)
			case pq.ListenerEventReconnected:
				logger.Debug("Postgres event broadcaster: reconnected")
			case pq.ListenerEventConnectionAttemptFailed:
				logger.Warnw("Postgres event broadcaster: reconnect attempt failed, trying again...", "error", err)
			}
		})

		go b.runLoop()
		return nil
	})
}

// Stop permanently destroys the EventBroadcaster.  Calling this does not clean
// up any outstanding subscriptions.  Subscribers must explicitly call `.Close()`
// or they will leak goroutines.
func (b *eventBroadcaster) Stop() error {
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
			logger.Debugw("Postgres event broadcaster: received notification",
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

func (b *eventBroadcaster) NotifyInsideGormTx(tx *gorm.DB, channel string, payload string) error {
	err := tx.Exec(`SELECT pg_notify(?, ?)`, channel, payload).Error
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
		queue:            utils.NewBoundedQueue(1000),
		chEvents:         make(chan Event),
		chDone:           make(chan struct{}),
	}
	sub.processQueueWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(sub.processQueue),
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
	subs, exists := b.subscriptions[sub.channelName()]
	if !exists || subs == nil {
		return
	}

	delete(b.subscriptions[sub.channelName()], sub)
	if len(b.subscriptions[sub.channelName()]) == 0 {
		err := b.listener.Unlisten(sub.channelName())
		if err != nil {
			logger.Errorw("Postgres event broadcaster: failed to unsubscribe", "error", err)
		}
		delete(b.subscriptions, sub.channelName())
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
		if sub.interestedIn(event) {
			wg.Add(1)
			go func(sub Subscription) {
				defer wg.Done()
				sub.send(event)
			}(sub)
		}
	}
	wg.Wait()
}

// Subscription represents a subscription to a Postgres event channel
type Subscription interface {
	Events() <-chan Event
	Close()

	channelName() string
	interestedIn(event Event) bool
	send(event Event)
}

type subscription struct {
	channel            string
	payloadFilter      string
	eventBroadcaster   *eventBroadcaster
	queue              *utils.BoundedQueue
	processQueueWorker utils.SleeperTask
	chEvents           chan Event
	chDone             chan struct{}
}

var _ Subscription = (*subscription)(nil)

func (sub *subscription) interestedIn(event Event) bool {
	return sub.payloadFilter == event.Payload || sub.payloadFilter == ""
}

func (sub *subscription) send(event Event) {
	sub.queue.Add(event)
	sub.processQueueWorker.WakeUpIfStarted()
}

const broadcastTimeout = 10 * time.Second

func (sub *subscription) processQueue() {
	ctx, cancel := context.WithTimeout(context.Background(), broadcastTimeout)
	defer cancel()

	for !sub.queue.Empty() {
		event, ok := sub.queue.Take().(Event)
		if !ok {
			logger.Errorf("Postgres event broadcaster subscription expected an Event, got %T", event)
			continue
		}
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

func (sub *subscription) channelName() string {
	return sub.channel
}

func (sub *subscription) Close() {
	sub.eventBroadcaster.removeSubscription(sub)
	// Close chDone before stopping the SleeperTask to avoid deadlocks
	close(sub.chDone)
	err := sub.processQueueWorker.Stop()
	if err != nil {
		logger.Errorw("THIS NEVER RETURNS AN ERROR", "error", err)
	}
}
