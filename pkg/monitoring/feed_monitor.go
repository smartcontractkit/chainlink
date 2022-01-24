package monitoring

import (
	"context"
	"sync"
)

type FeedMonitor interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
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

// Run should be executed as a goroutine
func (f *feedMonitor) Run(ctx context.Context, wg *sync.WaitGroup) {
	f.log.Infow("starting feed monitor")
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
		wg.Add(len(f.exporters))
		for _, exp := range f.exporters {
			go func(exp Exporter) {
				defer wg.Done()
				exp.Export(ctx, update)
			}(exp)
		}
	}
}
