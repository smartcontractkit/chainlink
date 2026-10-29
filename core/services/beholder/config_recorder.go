package beholder

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// ConfigRecorder periodically records Beholder config info metric.
type ConfigRecorder struct {
	services.Service
	eng *services.Engine

	client   *beholder.Client
	interval time.Duration
}

// NewConfigRecorder records client's config metric every interval. A nil
// client means a no-op client.
func NewConfigRecorder(lggr logger.Logger, client *beholder.Client, interval time.Duration) *ConfigRecorder {
	if client == nil {
		client = beholder.NewNoopClient()
	}
	cr := &ConfigRecorder{
		client:   client,
		interval: interval,
	}
	cr.Service, cr.eng = services.Config{
		Name:  "BeholderConfigRecorder",
		Start: cr.start,
	}.NewServiceEngine(lggr)
	return cr
}

func (s *ConfigRecorder) start(ctx context.Context) error {
	ticker := services.TickerConfig{}.NewTicker(s.interval)
	s.eng.GoTick(ticker, func(ctx context.Context) {
		if err := s.client.RecordConfigMetric(ctx); err != nil {
			s.eng.Errorf("failed to record beholder config metric: %v", err)
		}
	})
	return nil
}
