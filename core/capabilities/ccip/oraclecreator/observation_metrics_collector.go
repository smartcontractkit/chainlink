package oraclecreator

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	prometheus_dto "github.com/prometheus/client_model/go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// ObservationMetricsPublisher is the interface for publishing observation metrics to external destinations
type ObservationMetricsPublisher interface {
	PublishMetric(ctx context.Context, metricName string, value float64, labels map[string]string)
}

// ObservationMetricsCollector creates and wraps OCR3 observation metrics to intercept updates
type ObservationMetricsCollector struct {
	logger         logger.Logger
	publisher      ObservationMetricsPublisher
	cancel         context.CancelFunc
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
	_, cancel := context.WithCancel(context.Background())

	collector := &ObservationMetricsCollector{
		logger:         logger,
		publisher:      publisher,
		cancel:         cancel,
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

// Close stops the collector
func (c *ObservationMetricsCollector) Close() error {
	c.cancel()
	return nil
}

// wrappedCounter wraps a Prometheus collector (which may be a counter or wrappingCollector)
// to intercept Collect() calls and track value changes
type wrappedCounter struct {
	prometheus.Collector
	metricName    string
	labels        map[string]string // Beholder labels (for metrics publishing)
	publisher     ObservationMetricsPublisher
	logger        logger.Logger
	lastValueBits uint64 // stores float64 as bits for atomic operations
}

// Collect intercepts metric collection to detect counter increments
func (w *wrappedCounter) Collect(ch chan<- prometheus.Metric) {
	// Create a channel to intercept metrics
	interceptCh := make(chan prometheus.Metric, 10)

	// Collect from the underlying collector
	go func() {
		w.Collector.Collect(interceptCh)
		close(interceptCh)
	}()

	// Forward metrics and track counter value
	for m := range interceptCh {
		// Try to extract the counter value from the metric
		var metricValue float64
		if err := extractCounterValue(m, &metricValue); err == nil {
			// Load the last value atomically
			lastBits := atomic.LoadUint64(&w.lastValueBits)
			lastValue := math.Float64frombits(lastBits)

			if metricValue > lastValue {
				delta := metricValue - lastValue
				// Store the new value atomically
				atomic.StoreUint64(&w.lastValueBits, math.Float64bits(metricValue))

				w.logger.Debugw("Observation metric incremented",
					"metric", w.metricName,
					"value", metricValue,
					"delta", delta,
					"labels", w.labels,
				)

				if w.publisher != nil {
					// Publish the delta, not the cumulative value
					w.publisher.PublishMetric(context.Background(), w.metricName, delta, w.labels)
				}
			}
		}

		// Forward the metric to the actual channel
		ch <- m
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

// Describe implements prometheus.Collector
func (c *ObservationMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	if c.sentObservationsCounter != nil {
		c.sentObservationsCounter.Describe(ch)
	}
	if c.includedObservationsCounter != nil {
		c.includedObservationsCounter.Describe(ch)
	}
}

// Collect implements prometheus.Collector
func (c *ObservationMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	if c.sentObservationsCounter != nil {
		c.sentObservationsCounter.Collect(ch)
	}
	if c.includedObservationsCounter != nil {
		c.includedObservationsCounter.Collect(ch)
	}
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
