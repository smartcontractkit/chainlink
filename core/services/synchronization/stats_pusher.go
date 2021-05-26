package synchronization

import (
	"encoding/json"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

var (
	numberEventsSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stats_pusher_events_sent",
		Help: "The number of events pushed up to explorer",
	})

	gormCallbacksMutex *sync.RWMutex
)

func init() {
	gormCallbacksMutex = new(sync.RWMutex)
}

//go:generate mockery --name StatsPusher --output ../../internal/mocks/ --case=underscore

// StatsPusher polls for events and pushes them via a WebSocketClient. Events
// are consumed by the Explorer. Currently there is only one event type: an
// encoding of a JobRun.
type StatsPusher interface {
	service.Service
	PushNow()
	GetURL() url.URL
	GetStatus() ConnectionStatus
	AllSyncEvents(cb func(models.SyncEvent) error) error
}

type NoopStatsPusher struct{}

func (NoopStatsPusher) Start() error                                        { return nil }
func (NoopStatsPusher) Close() error                                        { return nil }
func (NoopStatsPusher) Ready() error                                        { return nil }
func (NoopStatsPusher) Healthy() error                                      { return nil }
func (NoopStatsPusher) PushNow()                                            {}
func (NoopStatsPusher) GetURL() url.URL                                     { return url.URL{} }
func (NoopStatsPusher) GetStatus() ConnectionStatus                         { return ConnectionStatusDisconnected }
func (NoopStatsPusher) AllSyncEvents(cb func(models.SyncEvent) error) error { return nil }

type statsPusher struct {
	DB             *gorm.DB
	ExplorerClient ExplorerClient
	Period         time.Duration
	clock          utils.Afterer
	backoffSleeper backoff.Backoff
	done           chan struct{}
	waker          chan struct{}

	utils.StartStopOnce
}

const (
	createCallbackName = "sync:run_after_create"
	updateCallbackName = "sync:run_after_update"
)

// NewStatsPusher returns a new StatsPusher service
func NewStatsPusher(db *gorm.DB, explorerClient ExplorerClient, afters ...utils.Afterer) StatsPusher {
	var clock utils.Afterer
	if len(afters) == 0 {
		clock = utils.Clock{}
	} else {
		clock = afters[0]
	}

	return &statsPusher{
		DB:             db,
		ExplorerClient: explorerClient,
		Period:         30 * time.Minute,
		clock:          clock,
		backoffSleeper: backoff.Backoff{
			Min: 1 * time.Second,
			Max: 5 * time.Minute,
		},
		done:  make(chan struct{}),
		waker: make(chan struct{}, 1),
	}
}

// GetURL returns the URL where stats are being pushed
func (sp *statsPusher) GetURL() url.URL {
	return sp.ExplorerClient.Url()
}

// GetStatus returns the ExplorerClient connection status
func (sp *statsPusher) GetStatus() ConnectionStatus {
	return sp.ExplorerClient.Status()
}

// Start starts the stats pusher
func (sp *statsPusher) Start() error {
	return sp.StartOnce("StatsPusher", func() error {
		gormCallbacksMutex.Lock()
		err := sp.DB.Callback().Create().Register(createCallbackName, createSyncEventWithStatsPusher(sp))
		if err != nil {
			return err
		}
		err = sp.DB.Callback().Update().Register(updateCallbackName, createSyncEventWithStatsPusher(sp))
		if err != nil {
			return err
		}
		gormCallbacksMutex.Unlock()

		go sp.eventLoop()
		return nil
	})
}

// Close shuts down the stats pusher
func (sp *statsPusher) Close() error {
	return sp.StopOnce("StatsPusher", func() error {
		close(sp.done)

		gormCallbacksMutex.Lock()
		if err := sp.DB.Callback().Create().Remove(createCallbackName); err != nil {
			return err
		}
		if err := sp.DB.Callback().Update().Remove(updateCallbackName); err != nil {
			return err
		}
		gormCallbacksMutex.Unlock()
		return nil
	})
}

// PushNow wakes up the stats pusher, asking it to push all queued events immediately.
func (sp *statsPusher) PushNow() {
	select {
	case sp.waker <- struct{}{}:
	default:
	}
}

type response struct {
	Status int `json:"status"`
}

func (sp *statsPusher) eventLoop() {
	logger.Debugw("Entered StatsPusher event loop")

	for {
		err := sp.pusherLoop()
		if err == nil {
			return
		}

		duration := sp.backoffSleeper.Duration()
		logger.Warnw("Failure during event synchronization", "error", err.Error(), "sleep_duration", duration)

		select {
		case <-sp.done:
			return
		case <-sp.clock.After(duration):
			continue
		}
	}
}

func (sp *statsPusher) pusherLoop() error {
	for {
		select {
		case <-sp.waker:
			err := sp.pushEvents()
			if err != nil {
				return err
			}
		case <-sp.clock.After(sp.Period):
			err := sp.pushEvents()
			if err != nil {
				return err
			}
		case <-sp.done:
			return nil
		}
	}
}

func (sp *statsPusher) pushEvents() error {
	gormCallbacksMutex.RLock()
	defer gormCallbacksMutex.RUnlock()
	err := sp.AllSyncEvents(func(event models.SyncEvent) error {
		return sp.syncEvent(event)
	})

	if err != nil {
		return errors.Wrap(err, "pushEvents#AllSyncEvents failed")
	}

	sp.backoffSleeper.Reset()
	return nil
}

func (sp *statsPusher) AllSyncEvents(cb func(models.SyncEvent) error) error {
	var events []models.SyncEvent
	err := sp.DB.
		Order("id, created_at asc").
		Find(&events).Error
	if err != nil {
		return err
	}

	for _, event := range events {
		err = cb(event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sp *statsPusher) syncEvent(event models.SyncEvent) error {
	ctx, cancel := utils.ContextFromChan(sp.done)
	defer cancel()

	sp.ExplorerClient.Send(ctx, []byte(event.Body))
	if ctx.Err() != nil {
		return nil
	}
	numberEventsSent.Inc()

	message, err := sp.ExplorerClient.Receive(ctx)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "syncEvent#ExplorerClient.Receive failed")
	}

	var resp response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return errors.Wrap(err, "syncEvent#json.Unmarshal failed")
	}

	if resp.Status != 201 {
		return errors.New("event not created")
	}

	err = sp.DB.WithContext(ctx).Delete(event).Error
	if ctx.Err() != nil {
	} else if err != nil {
		return errors.Wrap(err, "syncEvent#DB.Delete failed")
	}

	return nil
}

func createSyncEventWithStatsPusher(sp StatsPusher) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Error != nil {
			return
		}

		if db.Statement.Table != "job_runs" {
			return
		}

		if db.Statement.ReflectValue.Type() != reflect.TypeOf(models.JobRun{}) {
			logger.Errorf("Invariant violated scope.Value %T is not type models.JobRun, but TableName was job_runs", db.Statement.ReflectValue.Type())
			return
		}

		run, ok := db.Statement.ReflectValue.Interface().(models.JobRun)
		if !ok {
			db.Error = errors.Errorf("expected models.JobRun")
			return
		}

		// Note we have to use a separate db instance here
		// as the argument is part of a chain already targeting the job_runs table.
		db.Error = InsertSyncEventForJobRun(sp.(*statsPusher).DB, &run)
	}
}

// InsertSyncEventForJobRun generates an event for a change to job_runs
func InsertSyncEventForJobRun(db *gorm.DB, run *models.JobRun) error {
	presenter := SyncJobRunPresenter{run}
	bodyBytes, err := json.Marshal(presenter)
	if err != nil {
		return errors.Wrap(err, "createSyncEvent#json.Marshal failed")
	}

	event := models.SyncEvent{
		Body: string(bodyBytes),
	}

	err = db.Create(&event).Error
	return errors.Wrap(err, "createSyncEvent#Create failed")
}
