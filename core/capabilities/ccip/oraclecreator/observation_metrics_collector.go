package oraclecreator

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	prometheus_dto "github.com/prometheus/client_model/go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// defaultPollingInterval is the interval at which the collector polls counters and publishes
// deltas to Beholder, independent of Prometheus scrapes.
const defaultPollingInterval = 10 * time.Second

// ObservationMetricsPublisher is the interface for publishing observation metrics to external destinations
type ObservationMetricsPublisher interface {
	PublishMetric(ctx context.Context, metricName string, value float64, labels map[string]string)
}

// ObservationMetricsCollector creates and wraps OCR3 observation metrics to intercept updates
type ObservationMetricsCollector struct {
	logger         logger.Logger
	publisher      ObservationMetricsPublisher
	stop           chan struct{}
	startOnce      sync.Once
	stopOnce       sync.Once
	constantLabels map[string]string // Prometheus labels (for WrapRegistererWith)
	beholderLabels map[string]string // Beholder labels (for metrics publishing)

	// Wrapped counters
	sentObservationsCounter     *wrappedCounter
	includedObservationsCounter *wrappedCounter
}

// NewObservationMetricsCollector creates a new collector that wraps OCR3 observation metrics
func NewObservationMetricsCollector(
	logger logger.Logger,
	publisher ObservationMetricsPublisher,
	constantLabels map[string]string,
	beholderLabels map[string]string,
) *ObservationMetricsCollector {
	collector := &ObservationMetricsCollector{
		logger:         logger,
		publisher:      publisher,
		stop:           make(chan struct{}),
		constantLabels: constantLabels,
		beholderLabels: beholderLabels,
	}

	return collector
}

// CreateWrappedRegisterer returns a registerer that intercepts and wraps observation metrics
func (c *ObservationMetricsCollector) CreateWrappedRegisterer(baseRegisterer prometheus.Registerer) prometheus.Registerer {
	return &interceptingRegisterer{
		base:      baseRegisterer,
		collector: c,
	}
}

// Start launches a background goroutine that polls the wrapped counters on the given interval
// and publishes deltas to Beholder, independent of Prometheus scrapes.
// Call Start after the wrapped registerer has been passed to libocr (i.e. after NewOracle),
// so that the counters are already registered before the first poll fires.
// Safe to call multiple times; only the first call starts the goroutine.
func (c *ObservationMetricsCollector) Start(interval time.Duration) {
	if interval <= 0 {
		interval = defaultPollingInterval
	}
	c.startOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					c.poll()
				case <-c.stop:
					return
				}
			}
		}()
	})
}

// poll reads the current value of each wrapped counter and publishes any delta to Beholder.
func (c *ObservationMetricsCollector) poll() {
	if c.sentObservationsCounter != nil {
		c.sentObservationsCounter.readAndPublish()
	}
	if c.includedObservationsCounter != nil {
		c.includedObservationsCounter.readAndPublish()
	}
}

// Close stops the background polling goroutine. Safe to call multiple times.
func (c *ObservationMetricsCollector) Close() error {
	c.stopOnce.Do(func() { close(c.stop) })
	return nil
}

// wrappedCounter wraps a Prometheus collector (which may be a counter or wrappingCollector)
// to intercept Collect() calls and track value changes
type wrappedCounter struct {
	prometheus.Collector
	metricName string
	labels     map[string]string // Beholder labels (for metrics publishing)
	publisher  ObservationMetricsPublisher
	logger     logger.Logger
	lastValue  float64
}

// readAndPublish reads the current counter value and publishes any delta to Beholder.
// Only called sequentially by the background poller goroutine, so lastValue requires
// no synchronisation.
func (w *wrappedCounter) readAndPublish() {
	// Buffer sized to the typical max Prometheus series per counter (1 for a plain
	// Counter, more for a CounterVec). Sized generously to avoid blocking Collect.
	ch := make(chan prometheus.Metric, 16)
	w.Collect(ch)
	close(ch)

	for m := range ch {
		var metricValue float64
		if err := extractCounterValue(m, &metricValue); err != nil {
			continue
		}
		delta := metricValue - w.lastValue
		if delta > 0 {
			w.lastValue = metricValue
			w.logger.Debugw("Observation metric incremented",
				"metric", w.metricName,
				"value", metricValue,
				"delta", delta,
				"labels", w.labels,
			)
			if w.publisher != nil {
				w.publisher.PublishMetric(context.Background(), w.metricName, delta, w.labels)
			}
		}
	}
}

// extractCounterValue extracts the value from a prometheus.Metric
// This uses the prometheus dto.Metric structure
func extractCounterValue(m prometheus.Metric, value *float64) error {
	// Create a DTO metric to write into
	dto := &prometheus_dto.Metric{}
	if err := m.Write(dto); err != nil {
		return err
	}

	// Check if it's a counter
	if dto.Counter != nil {
		*value = dto.Counter.GetValue()
		return nil
	}

	return errors.New("metric is not a counter")
}

// interceptingRegisterer wraps a Prometheus registerer to intercept specific metric registrations
type interceptingRegisterer struct {
	base      prometheus.Registerer
	collector *ObservationMetricsCollector
}

func (r *interceptingRegisterer) Register(c prometheus.Collector) error {
	// Try to intercept counter registration
	// This returns either our wrappedCounter (for observation metrics)
	// or the original collector (for other metrics)
	wrapped := r.maybeWrapCollector(c)

	// If we wrapped it with our wrappedCounter, we still need to add Prometheus labels
	// If we didn't wrap it, we need to add Prometheus labels to maintain existing behavior
	wrappedWithLabels := prometheus.WrapCollectorWith(r.collector.constantLabels, wrapped)

	return r.base.Register(wrappedWithLabels)
}

func (r *interceptingRegisterer) MustRegister(cs ...prometheus.Collector) {
	wrapped := make([]prometheus.Collector, len(cs))
	for i, c := range cs {
		// Try to intercept and wrap with our custom wrapper
		maybeWrapped := r.maybeWrapCollector(c)
		// Add Prometheus labels to maintain existing behavior
		wrapped[i] = prometheus.WrapCollectorWith(r.collector.constantLabels, maybeWrapped)
	}
	r.base.MustRegister(wrapped...)
}

func (r *interceptingRegisterer) Unregister(c prometheus.Collector) bool {
	return r.base.Unregister(c)
}

// maybeWrapCollector checks if this is one of the observation counters and wraps it
func (r *interceptingRegisterer) maybeWrapCollector(c prometheus.Collector) prometheus.Collector {
	// Check if this is a Counter by trying to extract its descriptor
	descChan := make(chan *prometheus.Desc, 10)
	go func() {
		c.Describe(descChan)
		close(descChan)
	}()

	for desc := range descChan {
		descString := desc.String()

		// We need to extract the fqName from the descriptor string
		// Format: Desc{fqName: "metric_name", help: "...", ...}
		// We'll check if the fqName matches exactly, not just contains

		// Check if this is one of our target metrics by matching the fqName field
		if strings.Contains(descString, `fqName: "ocr3_sent_observations_total"`) {
			r.collector.logger.Info("Wrapping ocr3_sent_observations_total counter")

			// Wrap the collector (whether it's a raw Counter or wrappingCollector)
			wrapped := &wrappedCounter{
				Collector:  c,
				metricName: "ocr3_sent_observations_total",
				labels:     r.collector.beholderLabels,
				publisher:  r.collector.publisher,
				logger:     r.collector.logger,
			}
			r.collector.sentObservationsCounter = wrapped
			return wrapped
		}

		if strings.Contains(descString, `fqName: "ocr3_included_observations_total"`) {
			r.collector.logger.Info("Wrapping ocr3_included_observations_total counter")

			// Wrap the collector (whether it's a raw Counter or wrappingCollector)
			wrapped := &wrappedCounter{
				Collector:  c,
				metricName: "ocr3_included_observations_total",
				labels:     r.collector.beholderLabels,
				publisher:  r.collector.publisher,
				logger:     r.collector.logger,
			}
			r.collector.includedObservationsCounter = wrapped
			return wrapped
		}
	}

	// Not a metric we care about, return as-is
	return c
}
