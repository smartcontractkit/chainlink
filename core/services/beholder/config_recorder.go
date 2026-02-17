package beholder

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ConfigRecorder periodically records Beholder config info metric.
type ConfigRecorder struct {
	services.Service
	eng *services.Engine

	interval time.Duration
}

func NewConfigRecorder(logger logger.Logger, interval time.Duration) *ConfigRecorder {
	cr := &ConfigRecorder{
		interval: interval,
	}
	cr.Service, cr.eng = services.Config{
		Name:  "BeholderConfigRecorder",
		Start: cr.start,
	}.NewServiceEngine(logger)
	return cr
}

func (s *ConfigRecorder) start(ctx context.Context) error {
	ticker := services.TickerConfig{}.NewTicker(s.interval)
	s.eng.GoTick(ticker, func(ctx context.Context) {
		if err := beholder.GetClient().RecordConfigMetric(ctx); err != nil {
			s.eng.Errorf("failed to record beholder config metric: %v", err)
		}
	})
	return nil
}
