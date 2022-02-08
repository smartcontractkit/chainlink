package monitoring

import (
	"context"
	"errors"
	"time"
)

type Poller interface {
	Run(context.Context)
	// You should never close the channel returned by Updates()!
	// You should always read from the channel returned by Updates() in a
	// select statement with the same context you passed to Run()
	Updates() <-chan interface{}
}

func NewSourcePoller(
	source Source,
	log Logger,
	pollInterval time.Duration,
	fetchTimeout time.Duration,
	bufferCapacity uint32,
) Poller {
	return &sourcePoller{
		log,
		source,
		make(chan interface{}, bufferCapacity),
		pollInterval,
		fetchTimeout,
	}
}

type sourcePoller struct {
	log     Logger
	source  Source
	updates chan interface{}

	pollInterval time.Duration
	fetchTimeout time.Duration
}

// Run should be executed as a goroutine
func (s *sourcePoller) Run(ctx context.Context) {
	s.log.Debugw("poller started")
	defer s.log.Debugw("poller closed")
	// Initial fetch.
	data, err := s.source.Fetch(ctx)
	if err != nil {
		if errors.Is(err, ErrNoUpdate) {
			s.log.Debugw("no update found on initial fetch")
		} else if errors.Is(err, context.Canceled) {
			return
		} else {
			s.log.Errorw("failed initial fetch", "error", err)
		}
	} else {
		select {
		case s.updates <- data:
		case <-ctx.Done():
			return
		}
	}

	reusedTimer := time.NewTimer(s.pollInterval)
	for {
		select {
		case <-reusedTimer.C:
			var data interface{}
			var err error
			func() {
				ctx, cancel := context.WithTimeout(ctx, s.fetchTimeout)
				defer cancel()
				data, err = s.source.Fetch(ctx)
			}()
			if err != nil {
				if errors.Is(err, ErrNoUpdate) {
					s.log.Debugw("no update found")
					reusedTimer.Reset(s.pollInterval)
					continue
				} else if errors.Is(err, context.Canceled) {
					return
				} else {
					s.log.Errorw("failed to fetch from source", "error", err)
					reusedTimer.Reset(s.pollInterval)
					continue
				}
			}
			select {
			case s.updates <- data:
			case <-ctx.Done():
				return
			}
			reusedTimer.Reset(s.pollInterval)
		case <-ctx.Done():
			if !reusedTimer.Stop() {
				<-reusedTimer.C
			}
			return
		}
	}
}

func (s *sourcePoller) Updates() <-chan interface{} {
	return s.updates
}
