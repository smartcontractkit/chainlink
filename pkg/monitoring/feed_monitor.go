package monitoring

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

type FeedMonitor interface {
	Run(ctx context.Context)
}

func NewFeedMonitor(
	log Logger,
	pollers []Poller,
	exporters []Exporter,
) FeedMonitor {
	return &feedMonitor{
		log,
		pollers,
		exporters,
	}
}

type feedMonitor struct {
	log       Logger
	pollers   []Poller
	exporters []Exporter
}

// Run should be executed as a goroutine.
// Signal termination by cancelling ctx; then wait for Run() to exit.
func (f *feedMonitor) Run(ctx context.Context) {
	f.log.Infow("starting feed monitor")
	var subs utils.Subprocesses

	// Listen for updates
	updatesFanIn := make(chan interface{})
	for _, poller := range f.pollers {
		poller := poller
		subs.Go(func() {
			for {
				select {
				case update := <-poller.Updates():
					select {
					case updatesFanIn <- update:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		})
	}

	// Consume updates.
CONSUME_LOOP:
	for {
		var update interface{}
		select {
		case update = <-updatesFanIn:
		case <-ctx.Done():
			break CONSUME_LOOP
		}
		// TODO (dru) do we need a worker pool here?
		for index, exp := range f.exporters {
			index, exp := index, exp
			subs.Go(func() {
				defer func() {
					if err := recover(); err != nil {
						f.log.Errorw("failed Export", "error", err, "index", index)
					}
				}()
				exp.Export(ctx, update)
			})
		}
	}

	// Cleanup happens after all the exporters have finished.
	subs.Wait()
	subs = utils.Subprocesses{}
	defer subs.Wait()
	cleanupContext, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for index, exp := range f.exporters {
		index, exp := index, exp
		subs.Go(func() {
			defer func() {
				if err := recover(); err != nil {
					f.log.Errorw("failed Cleanup", "error", err, "index", index)
				}
			}()
			exp.Cleanup(cleanupContext)
		})
	}
}
