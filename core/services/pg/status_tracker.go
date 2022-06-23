package pg

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type StatusCallback func(connected bool)

type StatusTracker interface {
	services.ServiceCtx
	Subscribe(callback StatusCallback) (unsubscribe func())
}

type callbackSet map[int]StatusCallback

func (set callbackSet) values() []StatusCallback {
	var values []StatusCallback
	for _, callback := range set {
		values = append(values, callback)
	}
	return values
}

type statusTracker struct {
	db             *sqlx.DB
	callbacks      callbackSet
	lastCallbackID int
	mutex          sync.Mutex
	interval       time.Duration
	lggr           logger.Logger
	utils.StartStopOnce
	chStop chan struct{}
	wgDone sync.WaitGroup
}

func NewStatusTracker(db *sqlx.DB, interval time.Duration, lggr logger.Logger) StatusTracker {
	return &statusTracker{
		db:             db,
		callbacks:      make(callbackSet),
		lastCallbackID: 0,
		mutex:          sync.Mutex{},
		interval:       interval,
		lggr:           lggr.Named("PG.StatusTracker"),
		chStop:         make(chan struct{}),
		wgDone:         sync.WaitGroup{},
	}
}

func (st *statusTracker) Start(context.Context) error {
	return st.StartOnce("PG.StatusTracker", func() error {
		st.wgDone.Add(1)
		go st.loop()
		return nil
	})
}

func (st *statusTracker) Close() error {
	return st.StopOnce("PG.StatusTracker", func() error {
		close(st.chStop)
		st.wgDone.Wait()
		return nil
	})
}

func (st *statusTracker) Subscribe(handler StatusCallback) (unsubscribe func()) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.lastCallbackID++
	callbackID := st.lastCallbackID
	st.callbacks[callbackID] = handler
	unsubscribe = func() {
		st.mutex.Lock()
		defer st.mutex.Unlock()
		delete(st.callbacks, callbackID)
	}

	return
}

func (st *statusTracker) loop() {
	defer st.wgDone.Done()

	connected := st.db.Ping() == nil
	st.runHandlers(connected)

	for {
		select {
		case <-st.chStop:
			return
		case <-time.After(st.interval):
			ok := st.db.Ping() == nil
			if connected != ok {
				if ok {
					st.lggr.Info("Database connection is restored")
				} else {
					st.lggr.Error("Database connection is interrupted")
				}
				connected = ok
				st.runHandlers(connected)
			}
		}
	}
}

func (st *statusTracker) runHandlers(connected bool) {
	st.mutex.Lock()
	callbacks := st.callbacks.values()
	st.mutex.Unlock()

	for _, callback := range callbacks {
		callback(connected)
	}
}
