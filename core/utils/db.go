package utils

import (
	"time"

	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/tevino/abool"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

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
	chEvents chan string

	StartStopOnce
	chStop chan struct{}
	chDone chan struct{}
}

func (p *PostgresEventListener) Start() error {
	p.AssertNeverStarted()

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
				logger.Infow("Postgres listener: received notification",
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
	p.AssertNeverStopped()

	err := p.listener.Close()
	close(p.chStop)
	<-p.chDone
	return err
}

func (p *PostgresEventListener) Events() <-chan string {
	return p.chEvents
}

const (
	AdvisoryLockClassID_EthBroadcaster int32 = 0
	AdvisoryLockClassID_JobSpawner     int32 = 1
)

func GormTransaction(db *gorm.DB, fc func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("%+v", r)
			tx.Rollback()
		}
	}()

	err = fc(tx)

	if err == nil {
		err = errors.WithStack(tx.Commit().Error)
	}

	// Makesure rollback when Block error or Commit error
	if err != nil {
		tx.Rollback()
	}
	return
}

type PostgresAdvisoryLock struct {
	URI  string
	conn *sql.Conn
	db   *sql.DB
	mu   sync.Mutex
}

func (lock *PostgresAdvisoryLock) Close() error {
	var connErr, dbErr error

	if lock.conn != nil {
		connErr = lock.conn.Close()
		if connErr == sql.ErrConnDone {
			connErr = nil
		}
	}
	if lock.db != nil {
		dbErr = lock.db.Close()
		if dbErr == sql.ErrConnDone {
			dbErr = nil
		}
	}

	lock.db = nil
	lock.conn = nil

	return multierr.Combine(connErr, dbErr)
}

func (lock *PostgresAdvisoryLock) TryLock(ctx context.Context, classID int32, objectID int32) (err error) {
	lock.mu.Lock()
	defer lock.mu.Unlock()
	defer WrapIfError(&err, "TryAdvisoryLock failed")

	if lock.conn == nil {
		db, err := sql.Open("postgres", lock.URI)
		if err != nil {
			return err
		}
		lock.db = db

		// `database/sql`.DB does opaque connection pooling, but PG advisory locks are per-connection
		conn, err := db.Conn(ctx)
		if err != nil {
			lock.db.Close()
			lock.db = nil
			return err
		}
		lock.conn = conn
	}

	gotLock := false
	rows, err := lock.conn.QueryContext(ctx, "SELECT pg_try_advisory_lock($1, $2)", classID, objectID)
	if err != nil {
		return err
	}
	defer logger.ErrorIfCalling(rows.Close)
	gotRow := rows.Next()
	if !gotRow {
		return errors.New("query unexpectedly returned 0 rows")
	}
	if err := rows.Scan(&gotLock); err != nil {
		return err
	}
	if gotLock {
		return nil
	}
	return errors.Errorf("could not get advisory lock for classID, objectID %v, %v", classID, objectID)
}

func (lock *PostgresAdvisoryLock) Unlock(ctx context.Context, classID int32, objectID int32) error {
	lock.mu.Lock()
	defer lock.mu.Unlock()

	if lock.conn == nil {
		return nil
	}
	_, err := lock.conn.ExecContext(ctx, "SELECT pg_advisory_unlock($1, $2)", classID, objectID)
	return errors.Wrap(err, "AdvisoryUnlock failed")
}

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
		// automatically handles reconnecting to the DB.  Therefore, we do
		// not need to terminate `AwaitRun`, but rather simply log these
		// events for node operators' sanity.
		switch ev {
		case pq.ListenerEventConnected:
			logger.Infow("Postgres listener: connected", "channel", p.Event)
		case pq.ListenerEventDisconnected:
			logger.Warnw("Postgres listener: disconnected, trying to reconnect...", "channel", p.Event, "error", err)
		case pq.ListenerEventReconnected:
			logger.Info("Postgres listener: reconnected", "channel", p.Event)
		case pq.ListenerEventConnectionAttemptFailed:
			logger.Warnw("Postgres listener: reconnect attempt failed, trying again...", "channel", p.Event, "error", err)
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
