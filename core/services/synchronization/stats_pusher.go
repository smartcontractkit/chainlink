package synchronization

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// StatsPusher polls for events and pushes them via a WebSocketClient
type StatsPusher struct {
	ORM      *orm.ORM
	WSClient WebSocketClient
	Period   time.Duration
	cancel   context.CancelFunc
	clock    utils.Afterer
}

const (
	createCallbackName = "sync:run_after_create"
	updateCallbackName = "sync:run_after_update"
)

// NewStatsPusher returns a new event queuer
func NewStatsPusher(orm *orm.ORM, url *url.URL, afters ...utils.Afterer) *StatsPusher {
	var clock utils.Afterer
	if len(afters) == 0 {
		clock = utils.Clock{}
	} else {
		clock = afters[0]
	}

	var wsClient WebSocketClient
	wsClient = noopWebSocketClient{}
	if url != nil {
		wsClient = NewWebSocketClient(url)
		orm.DB.Callback().Create().Register(createCallbackName, createSyncEvent)
		orm.DB.Callback().Update().Register(updateCallbackName, createSyncEvent)
	}
	return &StatsPusher{
		ORM:      orm,
		WSClient: wsClient,
		Period:   60 * time.Second,
		clock:    clock,
	}
}

// Start starts the stats pusher
func (eq *StatsPusher) Start() error {
	err := eq.WSClient.Start()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	eq.cancel = cancel
	go eq.pollEvents(ctx)
	return nil
}

// Close shuts down the stats pusher
func (eq *StatsPusher) Close() error {
	if eq.cancel != nil {
		eq.cancel()
	}
	eq.ORM.DB.Callback().Create().Remove(createCallbackName)
	eq.ORM.DB.Callback().Update().Remove(updateCallbackName)
	return eq.WSClient.Close()
}

func (eq *StatsPusher) pollEvents(parentCtx context.Context) {
	for {
		select {
		case <-parentCtx.Done():
			return
		case <-eq.clock.After(eq.Period):
			err := eq.ORM.AllSyncEvents(func(event *models.SyncEvent) {
				logger.Debugw("StatsPusher got event", "event", event.ID)

				eq.WSClient.Send([]byte(event.Body))

				// TODO: This is fire and forget, we may want to get confirmation
				// before deleting...

				// TODO: This should also likely have backoff logic to avoid the
				// stampeding herd problem on the link stats server

				err := eq.ORM.DB.Delete(event).Error
				if err != nil {
					logger.Errorw("Error deleting event", "event_id", event.ID, "error", err)
				}
			})

			if err != nil {
				logger.Warnf("Error querying for sync events: %v", err)
			}
		}
	}
}

func createSyncEvent(scope *gorm.Scope) {
	if scope.HasError() {
		return
	}

	if scope.TableName() == "job_runs" {
		run, ok := scope.Value.(*models.JobRun)
		if !ok {
			return
		}

		presenter := SyncJobRunPresenter{run}
		bodyBytes, err := json.Marshal(presenter)
		if err != nil {
			scope.Err(err)
			return
		}

		event := models.SyncEvent{
			Body: string(bodyBytes),
		}
		err = scope.DB().Save(&event).Error
		scope.Err(err)
	}
}
