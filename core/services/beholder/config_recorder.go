package beholder

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ConfigRecorder periodically records Beholder config info metric.
type ConfigRecorder struct {
	services.StateMachine
	chStop   services.StopChan
	wgDone   sync.WaitGroup
	interval time.Duration
	logger   logger.Logger
}

func NewConfigRecorder(logger logger.Logger, interval time.Duration) *ConfigRecorder {
	return &ConfigRecorder{
		chStop:   make(services.StopChan),
		interval: interval,
		logger:   logger,
	}
}

func (s *ConfigRecorder) Start(ctx context.Context) error {
	return s.StartOnce("BeholderConfigRecorder", func() error {
		s.wgDone.Add(1)
		go s.run(ctx)
		return nil
	})
}

func (s *ConfigRecorder) Close() error {
	return s.StopOnce("BeholderConfigRecorder", func() error {
		close(s.chStop)
		s.wgDone.Wait()
		return nil
	})
}

func (s *ConfigRecorder) run(ctx context.Context) {
	defer s.wgDone.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		if err := beholder.GetClient().RecordConfigMetric(ctx); err != nil {
			s.logger.Errorf("failed to record beholder config metric: %v", err)
		}
		select {
		case <-ticker.C:
		case <-s.chStop:
			return
		case <-ctx.Done():
			return
		}
	}
}
