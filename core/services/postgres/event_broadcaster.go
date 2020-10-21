package postgres

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
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
	NotifyUsing(db *sql.DB, channel string, payload string) error
}

type eventBroadcaster struct {
	uri                  string
	minReconnectInterval time.Duration
	maxReconnectDuration time.Duration
	db                   *sql.DB

	listeners   map[string]*channelListener
	listenersMu sync.Mutex

	utils.StartStopOnce
}

var _ EventBroadcaster = (*eventBroadcaster)(nil)

type Event struct {
	Channel string
	Payload string
}

func NewEventBroadcaster(uri string, minReconnectInterval time.Duration, maxReconnectDuration time.Duration) *eventBroadcaster {
	if minReconnectInterval == time.Duration(0) {
		minReconnectInterval = 1 * time.Second
	}
	if maxReconnectDuration == time.Duration(0) {
		maxReconnectDuration = 1 * time.Minute
	}
	return &eventBroadcaster{
		uri:                  uri,
		minReconnectInterval: minReconnectInterval,
		maxReconnectDuration: maxReconnectDuration,
		listeners:            make(map[string]*channelListener),
	}
}

func (b *eventBroadcaster) Start() error {
	if !b.OkayToStart() {
		return errors.Errorf("Postgres event broadcaster has already been started")
	}
	db, err := sql.Open("postgres", b.uri)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

func (b *eventBroadcaster) Stop() error {
	if !b.OkayToStop() {
		return errors.Errorf("Postgres event broadcaster has already been stopped")
	}

	err := b.db.Close()

	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()

	for _, listener := range b.listeners {
		err = multierr.Append(err, listener.stop())
	}
	b.listeners = nil // avoid "close of closed channel" panic on shutdown
	return err
}

func (b *eventBroadcaster) Notify(channel string, payload string) error {
	return b.NotifyUsing(b.db, channel, payload)
}

func (b *eventBroadcaster) NotifyUsing(db *sql.DB, channel string, payload string) error {
	_, err := db.Exec(`SELECT pg_notify($1, $2::text)`, channel, payload)
	return errors.Wrap(err, "Postgres event broadcaster could not notify")
}

func (b *eventBroadcaster) Subscribe(channel, payloadFilter string) (Subscription, error) {
	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()

	if b.listeners[channel] == nil {
		listener := newChannelListener(b.uri, b.minReconnectInterval, b.maxReconnectDuration, channel, b)
		err := listener.start()
		if err != nil {
			return nil, err
		}
		b.listeners[channel] = listener
	}
	return b.listeners[channel].subscribe(payloadFilter), nil
}

func (b *eventBroadcaster) unsubscribe(sub Subscription) {
	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()

	listener, exists := b.listeners[sub.channelName()]
	if !exists {
		// This occurs on shutdown when .Stop() is called before one
		// or more subscriptions' .Close() methods are called
		return
	}

	listener.unsubscribe(sub)
	if len(listener.subscriptions) == 0 {
		err := listener.stop()
		if err != nil {
			logger.Errorw("Postgres event broadcaster could not close listener", "error", err)
		}
		delete(b.listeners, sub.channelName())
	}
}

type channelListener struct {
	*pq.Listener
	utils.StartStopOnce
	channel          string
	eventBroadcaster *eventBroadcaster
	subscriptions    map[Subscription]struct{}
	subscriptionsMu  sync.RWMutex
	chStop           chan struct{}
	chDone           chan struct{}
}

func newChannelListener(
	uri string,
	minReconnectInterval time.Duration,
	maxReconnectDuration time.Duration,
	channel string,
	eventBroadcaster *eventBroadcaster,
) *channelListener {
	pqListener := pq.NewListener(uri, minReconnectInterval, maxReconnectDuration, func(ev pq.ListenerEventType, err error) {
		// These are always connection-related events, and the pq library
		// automatically handles reconnecting to the DB. Therefore, we do not
		// need to terminate, but rather simply log these events for node
		// operators' sanity.
		switch ev {
		case pq.ListenerEventConnected:
			logger.Debugw("Postgres listener: connected", "channel", channel)
		case pq.ListenerEventDisconnected:
			logger.Warnw("Postgres listener: disconnected, trying to reconnect...", "channel", channel, "error", err)
		case pq.ListenerEventReconnected:
			logger.Debugw("Postgres listener: reconnected", "channel", channel)
		case pq.ListenerEventConnectionAttemptFailed:
			logger.Warnw("Postgres listener: reconnect attempt failed, trying again...", "channel", channel, "error", err)
		}
	})

	listener := &channelListener{
		Listener:         pqListener,
		channel:          channel,
		eventBroadcaster: eventBroadcaster,
		subscriptions:    make(map[Subscription]struct{}),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
	return listener
}

func (listener *channelListener) start() error {
	if !listener.OkayToStart() {
		return errors.Errorf("Postgres event listener has already been started (channel: %v)", listener.channel)
	}

	err := listener.Listener.Listen(listener.channel)
	if err != nil {
		return err
	}

	go func() {
		defer close(listener.chDone)

		for {
			select {
			case <-listener.chStop:
				return

			case notification, open := <-listener.NotificationChannel():
				if !open {
					return
				} else if notification == nil {
					continue
				}
				logger.Debugw("Postgres listener: received notification",
					"channel", notification.Channel,
					"payload", notification.Extra,
				)
				listener.broadcast(notification)
			}
		}
	}()
	return nil
}

const maxBroadcastDuration = 10 * time.Second

func (listener *channelListener) broadcast(notification *pq.Notification) {
	listener.subscriptionsMu.RLock()
	defer listener.subscriptionsMu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), maxBroadcastDuration)
	defer cancel()

	event := Event{Channel: notification.Channel, Payload: notification.Extra}

	var wg sync.WaitGroup
	for sub := range listener.subscriptions {
		if sub.interestedIn(event) {
			wg.Add(1)
			go func(sub Subscription) {
				defer wg.Done()
				sub.send(ctx, event)
			}(sub)
		}
	}
	wg.Wait()
}

func (listener *channelListener) subscribe(payloadFilter string) *subscription {
	listener.subscriptionsMu.Lock()
	defer listener.subscriptionsMu.Unlock()

	sub := &subscription{
		channel:          listener.channel,
		payloadFilter:    payloadFilter,
		eventBroadcaster: listener.eventBroadcaster,
		chEvents:         make(chan Event),
		chDone:           make(chan struct{}),
	}
	listener.subscriptions[sub] = struct{}{}
	return sub
}

func (listener *channelListener) unsubscribe(sub Subscription) {
	listener.subscriptionsMu.Lock()
	defer listener.subscriptionsMu.Unlock()

	sub.close()
	delete(listener.subscriptions, sub)
}

func (listener *channelListener) stop() error {
	if !listener.OkayToStop() {
		return errors.Errorf("Postgres event listener has already been stopped (channel: %v)", listener.channel)
	}
	err := listener.Listener.Close()

	listener.subscriptionsMu.Lock()
	defer listener.subscriptionsMu.Unlock()

	for sub := range listener.subscriptions {
		sub.close()
	}

	close(listener.chStop)
	<-listener.chDone
	return err
}

// Subscription represents a subscription to a Postgres event channel
type Subscription interface {
	Events() <-chan Event
	Close()

	channelName() string
	interestedIn(event Event) bool
	send(ctx context.Context, event Event)
	close()
}

type subscription struct {
	channel          string
	payloadFilter    string
	eventBroadcaster *eventBroadcaster
	chEvents         chan Event
	chDone           chan struct{}
}

var _ Subscription = (*subscription)(nil)

func (sub *subscription) interestedIn(event Event) bool {
	return sub.payloadFilter == event.Payload || sub.payloadFilter == ""
}

func (sub *subscription) send(ctx context.Context, event Event) {
	select {
	case sub.chEvents <- event:
	case <-ctx.Done():
	case <-sub.chDone:
	}
}

func (sub *subscription) Events() <-chan Event {
	return sub.chEvents
}

func (sub *subscription) channelName() string {
	return sub.channel
}

func (sub *subscription) close() {
	close(sub.chDone)
}

func (sub *subscription) Close() {
	sub.eventBroadcaster.unsubscribe(sub)
}
