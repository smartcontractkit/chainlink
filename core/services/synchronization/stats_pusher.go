package synchronization

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"chainlink/core/logger"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/jinzhu/gorm"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
)

// StatsPusher polls for events and pushes them via a WebSocketClient
type StatsPusher struct {
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
func NewStatsPusher(orm *orm.ORM, url *url.URL, accessKey, secret string, afters ...utils.Afterer) *StatsPusher {
	var clock utils.Afterer
	if len(afters) == 0 {
		clock = utils.Clock{}
	} else {
		clock = afters[0]
	}

	sp := &StatsPusher{
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

// Start starts the stats pusher
func (sp *StatsPusher) Start() error {
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
func (sp *StatsPusher) Close() error {
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
func (sp *StatsPusher) PushNow() {
	select {
	case sp.waker <- struct{}{}:
	default:
	}
}

type response struct {
	Status int `json:"status"`
}

func (sp *StatsPusher) eventLoop(parentCtx context.Context) {
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

func (sp *StatsPusher) pusherLoop(parentCtx context.Context) error {
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

func (sp *StatsPusher) pushEvents() error {
	err := sp.ORM.AllSyncEvents(func(event *models.SyncEvent) error {
		return sp.syncEvent(event)
	})

	if err != nil {
		return errors.Wrap(err, "pushEvents#AllSyncEvents failed")
	}

	sp.backoffSleeper.Reset()
	return nil
}

func (sp *StatsPusher) syncEvent(event *models.SyncEvent) error {
	sp.WSClient.Send([]byte(event.Body))

	message, err := sp.WSClient.Receive()
	if err != nil {
		return errors.Wrap(err, "syncEvent#WSClient.Receive failed")
	}

	var response response
	err = json.Unmarshal(message, &response)
	if err != nil {
		return errors.Wrap(err, "syncEvent#json.Unmarshal failed")
	}

	if response.Status != 201 {
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

func createSyncEventWithStatsPusher(sp *StatsPusher, orm *orm.ORM) func(*gorm.Scope) {
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
			scope.Err(errors.Wrap(err, "createSyncEvent#json.Marshal failed"))
			return
		}

		event := models.SyncEvent{
			Body: string(bodyBytes),
		}
		err = scope.DB().Save(&event).Error
		if err != nil {
			scope.Err(errors.Wrap(err, "createSyncEvent#Save failed"))
			return
		}

		sp.PushNow()
	}
}

var (
	gormCallbacksMutex *sync.Mutex
)

func init() {
	gormCallbacksMutex = new(sync.Mutex)
}
