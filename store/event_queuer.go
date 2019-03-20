package store

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/orm"
)

type EventQueuer struct {
	ORM    *orm.ORM
	cancel context.CancelFunc
}

func NewEventQueuer(orm *orm.ORM) *EventQueuer {
	eq := &EventQueuer{
		ORM: orm,
	}

	ctx, cancel := context.WithCancel(context.Background())
	eq.cancel = cancel
	go eq.pollEvents(ctx)

	return eq
}

func (eq EventQueuer) Shutdown() {
	eq.cancel()
}

func (eq EventQueuer) pollEvents(parentCtx context.Context) {
	pollTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-parentCtx.Done():
			return
		case <-pollTicker.C:
			events, err := eq.ORM.AllSyncEvents()
			if err != nil {
				logger.Warn("Error querying for sync events: %v", err)
				continue
			}

			for event := range events {
				fmt.Println("syncing event", event)
			}
		}
	}
}
