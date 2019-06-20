package synchronization

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ExplorerPusher polls for events and pushes them via a WebSocketClient
type ExplorerPusher struct {
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

// NewExplorerPusher returns a new event queuer
func NewExplorerPusher(orm *orm.ORM, url *url.URL, accessKey, secret string, afters ...utils.Afterer) *ExplorerPusher {
	var clock utils.Afterer
	if len(afters) == 0 {
		clock = utils.Clock{}
	} else {
		clock = afters[0]
	}

	ep := &ExplorerPusher{
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
		ep.WSClient = NewWebSocketClient(url, accessKey, secret)
		gormCallbacksMutex.Lock()
		orm.DB.Callback().Create().Register(createCallbackName, createSyncEventWithExplorerPusher(ep))
		orm.DB.Callback().Update().Register(updateCallbackName, createSyncEventWithExplorerPusher(ep))
		gormCallbacksMutex.Unlock()
	}
	return ep
}

// Start starts the stats pusher
func (ep *ExplorerPusher) Start() error {
	err := ep.WSClient.Start()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	ep.cancel = cancel
	go ep.eventLoop(ctx)
	logger.Infow("Started ExplorerPusher")
	return nil
}

// Close shuts down the stats pusher
func (ep *ExplorerPusher) Close() error {
	if ep.cancel != nil {
		ep.cancel()
	}
	gormCallbacksMutex.Lock()
	callbacks := ep.ORM.DB.Callback()
	callbacks.Create().Remove(createCallbackName)
	callbacks.Update().Remove(updateCallbackName)
	gormCallbacksMutex.Unlock()
	return ep.WSClient.Close()
}

// PushNow wakes up the stats pusher, asking it to push all queued events immediately.
func (ep *ExplorerPusher) PushNow() {
	select {
	case ep.waker <- struct{}{}:
	default:
	}
}

type response struct {
	Status int `json:"status"`
}

func (ep *ExplorerPusher) eventLoop(parentCtx context.Context) {
	logger.Debugw("Entered ExplorerPusher event loop")
	for {
		err := ep.pusherLoop(parentCtx)
		if err == nil {
			return
		}

		duration := ep.backoffSleeper.Duration()
		logger.Warnw("Failure during event synchronization", "error", err.Error(), "sleep_duration", duration)

		select {
		case <-parentCtx.Done():
			return
		case <-ep.clock.After(duration):
			continue
		}
	}
}

func (ep *ExplorerPusher) pusherLoop(parentCtx context.Context) error {
	logger.Debugw("Entered ExplorerPusher push loop")

	for {
		select {
		case <-ep.waker:
			err := ep.pushEvents()
			if err != nil {
				return err
			}
		case <-ep.clock.After(ep.Period):
			err := ep.pushEvents()
			if err != nil {
				return err
			}
		case <-parentCtx.Done():
			logger.Debugw("ExplorerPusher got done signal, shutting down")
			return nil
		}
	}
}

func (ep *ExplorerPusher) pushEvents() error {
	err := ep.ORM.AllSyncEvents(func(event *models.SyncEvent) error {
		logger.Debugw("ExplorerPusher got event", "event", event.ID)
		return ep.syncEvent(event)
	})

	if err != nil {
		return errors.Wrap(err, "pushEvents#AllSyncEvents failed")
	}

	ep.backoffSleeper.Reset()
	return nil
}

func (ep *ExplorerPusher) syncEvent(event *models.SyncEvent) error {
	ep.WSClient.Send([]byte(event.Body))

	message, err := ep.WSClient.Receive()
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

	err = ep.ORM.DB.Delete(event).Error
	if err != nil {
		return errors.Wrap(err, "syncEvent#DB.Delete failed")
	}

	return nil
}

func createSyncEventWithExplorerPusher(ep *ExplorerPusher) func(*gorm.Scope) {
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

		ep.PushNow()
	}
}

var (
	gormCallbacksMutex *sync.Mutex
)

func init() {
	gormCallbacksMutex = new(sync.Mutex)
}
