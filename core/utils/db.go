package utils

import (
	"time"

	"github.com/lib/pq"
	"github.com/tevino/abool"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// PostgresEventListener listens to one type of postgres event as emitted by NOTIFY
// TODO(sam): Currently each new listener opens a new connection. This is
// suboptimal. We ought to move to a design similar to LogBroadcaster where we
// hold only one open connection and listen to everything on that.
type PostgresEventListener struct {
	URI                  string
	Event                string
	PayloadFilter        string
	MinReconnectInterval time.Duration
	MaxReconnectDuration time.Duration

	listener *pq.Listener
	started  abool.AtomicBool
	stopped  abool.AtomicBool
	chEvents chan string
	chStop   chan struct{}
	chDone   chan struct{}
}

func (p *PostgresEventListener) Start() error {
	if !p.started.SetToIf(false, true) {
		panic("PostgresEventListener can only be started once")
	}

	if p.MinReconnectInterval == time.Duration(0) {
		p.MinReconnectInterval = 1 * time.Second
	}
	if p.MaxReconnectDuration == time.Duration(0) {
		p.MaxReconnectDuration = 1 * time.Minute
	}

	p.chEvents = make(chan string)
	p.chStop = make(chan struct{})
	p.chDone = make(chan struct{})

	p.listener = pq.NewListener(p.URI, p.MinReconnectInterval, p.MaxReconnectDuration, func(ev pq.ListenerEventType, err error) {
		// These are always connection-related events, and the pq library
		// automatically handles reconnecting to the DB. Therefore, we do not
		// need to terminate, but rather simply log these events for node
		// operators' sanity.
		switch ev {
		case pq.ListenerEventConnected:
			logger.Debug("Postgres listener: connected")
		case pq.ListenerEventDisconnected:
			logger.Warnw("Postgres listener: disconnected, trying to reconnect...", "error", err)
		case pq.ListenerEventReconnected:
			logger.Debug("Postgres listener: reconnected")
		case pq.ListenerEventConnectionAttemptFailed:
			logger.Warnw("Postgres listener: reconnect attempt failed, trying again...", "error", err)
		}
	})
	err := p.listener.Listen(p.Event)
	if err != nil {
		return err
	}

	go func() {
		defer close(p.chDone)

		for {
			select {
			case <-p.chStop:
				return
			case notification, open := <-p.listener.NotificationChannel():
				if !open {
					return
				}
				logger.Debugw("Postgres listener: received notification",
					"channel", notification.Channel,
					"payload", notification.Extra,
				)
				if p.PayloadFilter == "" || p.PayloadFilter == notification.Extra {
					p.chEvents <- notification.Extra
				}
			}
		}
	}()
	return nil
}

func (p *PostgresEventListener) Stop() error {
	if !p.started.IsSet() {
		panic("PostgresEventListener cannot stop before starting")
	}
	if !p.stopped.SetToIf(false, true) {
		panic("PostgresEventListener can only be stopped once")
	}
	err := p.listener.Close()
	close(p.chStop)
	<-p.chDone
	return err
}

func (p *PostgresEventListener) Events() <-chan string {
	return p.chEvents
}
