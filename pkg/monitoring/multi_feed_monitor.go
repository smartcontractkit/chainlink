package monitoring

import (
	"context"
	"sync"
)

type MultiFeedMonitor interface {
	Start(ctx context.Context, wg *sync.WaitGroup, feeds []FeedConfig)
}

func NewMultiFeedMonitor(
	chainConfig ChainConfig,

	log Logger,
	sourceFactory SourceFactory,
	producer Producer,
	metrics Metrics,

	transmissionTopic string,
	configSetSimplifiedTopic string,

	transmissionSchema Schema,
	configSetSimplifiedSchema Schema,
) MultiFeedMonitor {
	return &multiFeedMonitor{
		chainConfig,

		log,
		sourceFactory,
		producer,
		metrics,

		transmissionTopic,
		configSetSimplifiedTopic,

		transmissionSchema,
		configSetSimplifiedSchema,
	}
}

type multiFeedMonitor struct {
	chainConfig ChainConfig

	log           Logger
	sourceFactory SourceFactory
	producer      Producer
	metrics       Metrics

	transmissionTopic        string
	configSetSimplifiedTopic string

	transmissionSchema        Schema
	configSetSimplifiedSchema Schema
}

const bufferCapacity = 100

// Start should be executed as a goroutine.
func (m *multiFeedMonitor) Start(ctx context.Context, wg *sync.WaitGroup, feeds []FeedConfig) {
	wg.Add(len(feeds))
	for _, feedConfig := range feeds {
		go func(feedConfig FeedConfig) {
			defer wg.Done()

			feedLogger := m.log.With(
				"feed", feedConfig.GetName(),
				"network", m.chainConfig.GetNetworkName(),
			)
			source, err := m.sourceFactory.NewSource(m.chainConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create new source", "error", err)
				return
			}
			poller := NewSourcePoller(
				source,
				feedLogger.With("component", "chain-poller"),
				m.chainConfig.GetPollInterval(),
				m.chainConfig.GetReadTimeout(),
				bufferCapacity,
			)

			wg.Add(1)
			go func() {
				defer wg.Done()
				poller.Start(ctx)
			}()

			exporters := []Exporter{
				NewPrometheusExporter(
					m.chainConfig,
					feedConfig,
					feedLogger.With("component", "prometheus-exporter"),
					m.metrics,
				),
				NewKafkaExporter(
					m.chainConfig,
					feedConfig,
					feedLogger.With("component", "kafka-exporter"),
					m.producer,

					m.transmissionSchema,
					m.configSetSimplifiedSchema,

					m.transmissionTopic,
					m.configSetSimplifiedTopic,
				),
			}

			feedMonitor := NewFeedMonitor(
				feedLogger.With("component", "feed-monitor"),
				poller,
				exporters,
			)
			feedMonitor.Start(ctx, wg)
		}(feedConfig)
	}
}
