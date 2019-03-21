package store

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

// EventQueuer polls for events and pushes them via a StatsPusher
type EventQueuer struct {
	ORM         *orm.ORM
	StatsPusher StatsPusher
	cancel      context.CancelFunc
	Period      time.Duration
}

// NewEventQueuer returns a new event queuer
func NewEventQueuer(orm *orm.ORM, statsPusher StatsPusher) *EventQueuer {
	return &EventQueuer{
		ORM:         orm,
		StatsPusher: statsPusher,
		Period:      60 * time.Second,
	}
}

// Start starts the event queuer
func (eq EventQueuer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	eq.cancel = cancel
	go eq.pollEvents(ctx)
	return nil
}

// Shutdown stops the event queuer
func (eq EventQueuer) Shutdown() {
	eq.cancel()
}

func (eq EventQueuer) pollEvents(parentCtx context.Context) {
	pollTicker := time.NewTicker(eq.Period)

	for {
		select {
		case <-parentCtx.Done():
			return
		case <-pollTicker.C:
			err := eq.ORM.AllSyncEvents(func(event *models.SyncEvent) {
				fmt.Println("EventQueuer got event", event)

				eq.StatsPusher.Send([]byte(event.Body))

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
