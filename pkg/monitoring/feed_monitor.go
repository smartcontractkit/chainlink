package monitoring

import (
	"context"
	"sync"
)

type FeedMonitor interface {
	Run(ctx context.Context)
}

func NewFeedMonitor(
	log Logger,
	poller Poller,
	exporters []Exporter,
) FeedMonitor {
	return &feedMonitor{
		log,
		poller,
		exporters,
	}
}

type feedMonitor struct {
	log       Logger
	poller    Poller
	exporters []Exporter
}

// Run should be executed as a goroutine.
// Signal termination by cancelling ctx; then wait for Run() to exit.
func (f *feedMonitor) Run(ctx context.Context) {
	f.log.Infow("starting feed monitor")
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	for {
		// Wait for an update.
		var update interface{}
		select {
		case update = <-f.poller.Updates():
		case <-ctx.Done():
			for _, exp := range f.exporters {
				exp.Cleanup()
			}
			return
		}
		// TODO (dru) do we need a worker pool here?
		wg.Add(len(f.exporters))
		for _, exp := range f.exporters {
			go func(exp Exporter) {
				defer wg.Done()
				exp.Export(ctx, update)
			}(exp)
		}
	}
}
