package synchronization

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/jinzhu/gorm"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	Start() error
	Close() error
	PushNow()
	GetURL() url.URL
	GetStatus() ConnectionStatus
}

type statsPusher struct {
	ORM            *orm.ORM
	WSClient       WebSocketClient
	Period         time.Duration
	cancel         context.CancelFunc
	clock          utils.Afterer
	backoffSleeper backoff.Backoff
	waker          chan struct{}
}

const (
	createCallbackName = "sync:run_after_create"
	updateCallbackName = "sync:run_after_update"
)

// NewStatsPusher returns a new event queuer
func NewStatsPusher(orm *orm.ORM, url *url.URL, accessKey, secret string, afters ...utils.Afterer) StatsPusher {
	var clock utils.Afterer
	if len(afters) == 0 {
		clock = utils.Clock{}
	} else {
		clock = afters[0]
	}

	sp := &statsPusher{
		ORM:      orm,
		WSClient: noopWebSocketClient{},
		Period:   30 * time.Minute,
		clock:    clock,
		backoffSleeper: backoff.Backoff{
			Min: 1 * time.Second,
			Max: 5 * time.Minute,
		},
		waker: make(chan struct{}, 1),
	}

	if url != nil {
		sp.WSClient = NewWebSocketClient(url, accessKey, secret)
		gormCallbacksMutex.Lock()
		_ = orm.RawDB(func(db *gorm.DB) error {
			db.Callback().Create().Register(createCallbackName, createSyncEventWithStatsPusher(sp, orm))
			db.Callback().Update().Register(updateCallbackName, createSyncEventWithStatsPusher(sp, orm))
			return nil
		})
		gormCallbacksMutex.Unlock()
	}
	return sp
}

func (sp *statsPusher) GetURL() url.URL {
	return sp.WSClient.Url()
}

func (sp *statsPusher) GetStatus() ConnectionStatus {
	return sp.WSClient.Status()
}

// Start starts the stats pusher
func (sp *statsPusher) Start() error {
	err := sp.WSClient.Start()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	sp.cancel = cancel
	go sp.eventLoop(ctx)
	return nil
}

// Close shuts down the stats pusher
func (sp *statsPusher) Close() error {
	if sp.cancel != nil {
		sp.cancel()
	}
	gormCallbacksMutex.Lock()
	_ = sp.ORM.RawDB(func(db *gorm.DB) error {
		db.Callback().Create().Remove(createCallbackName)
		db.Callback().Update().Remove(updateCallbackName)
		return nil
	})
	gormCallbacksMutex.Unlock()
	return sp.WSClient.Close()
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

func (sp *statsPusher) eventLoop(parentCtx context.Context) {
	logger.Debugw("Entered StatsPusher event loop")
	for {
		err := sp.pusherLoop(parentCtx)
		if err == nil {
			return
		}

		duration := sp.backoffSleeper.Duration()
		logger.Warnw("Failure during event synchronization", "error", err.Error(), "sleep_duration", duration)

		select {
		case <-parentCtx.Done():
			return
		case <-sp.clock.After(duration):
			continue
		}
	}
}

func (sp *statsPusher) pusherLoop(parentCtx context.Context) error {
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
		case <-parentCtx.Done():
			return nil
		}
	}
}

func (sp *statsPusher) pushEvents() error {
	gormCallbacksMutex.RLock()
	defer gormCallbacksMutex.RUnlock()
	err := sp.ORM.AllSyncEvents(func(event models.SyncEvent) error {
		return sp.syncEvent(event)
	})

	if err != nil {
		return errors.Wrap(err, "pushEvents#AllSyncEvents failed")
	}

	sp.backoffSleeper.Reset()
	return nil
}

func (sp *statsPusher) syncEvent(event models.SyncEvent) error {
	sp.WSClient.Send([]byte(event.Body))
	numberEventsSent.Inc()

	message, err := sp.WSClient.Receive()
	if err != nil {
		return errors.Wrap(err, "syncEvent#WSClient.Receive failed")
	}

	var resp response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return errors.Wrap(err, "syncEvent#json.Unmarshal failed")
	}

	if resp.Status != 201 {
		return errors.New("event not created")
	}

	err = sp.ORM.RawDB(func(db *gorm.DB) error {
		return db.Delete(event).Error
	})
	if err != nil {
		return errors.Wrap(err, "syncEvent#DB.Delete failed")
	}

	return nil
}

func createSyncEventWithStatsPusher(sp StatsPusher, orm *orm.ORM) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		if scope.HasError() {
			return
		}

		if scope.TableName() != "job_runs" {
			return
		}

		run, ok := scope.Value.(*models.JobRun)
		if !ok {
			logger.Error("Invariant violated scope.Value is not type *models.JobRun, but TableName was job_runes")
			return
		}

		orm.MustEnsureAdvisoryLock()

		presenter := SyncJobRunPresenter{run}
		bodyBytes, err := json.Marshal(presenter)
		if err != nil {
			_ = scope.Err(errors.Wrap(err, "createSyncEvent#json.Marshal failed"))
			return
		}

		event := models.SyncEvent{
			Body: string(bodyBytes),
		}
		err = scope.DB().Create(&event).Error
		if err != nil {
			_ = scope.Err(errors.Wrap(err, "createSyncEvent#Create failed"))
			return
		}
	}
}
