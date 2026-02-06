package oraclecreator

import (
	"context"
	"strings"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"

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

// wrappedCounter wraps a Prometheus counter to intercept Inc() calls
type wrappedCounter struct {
	prometheus.Counter
	metricName string
	labels     map[string]string // Beholder labels (for metrics publishing)
	publisher  ObservationMetricsPublisher
	logger     logger.Logger
	value      atomic.Uint64
}

// Inc increments the counter and publishes the delta (1)
func (w *wrappedCounter) Inc() {
	w.Counter.Inc()
	newValue := w.value.Add(1)

	w.logger.Debugw("Observation metric incremented",
		"metric", w.metricName,
		"value", newValue,
		"delta", 1,
		"labels", w.labels,
	)

	if w.publisher != nil {
		// Publish the delta (1), not the cumulative value, since Beholder counters are cumulative
		w.publisher.PublishMetric(context.Background(), w.metricName, 1.0, w.labels)
	}
}

// Add increments the counter by the given value and publishes the delta
func (w *wrappedCounter) Add(val float64) {
	w.Counter.Add(val)
	newValue := w.value.Add(uint64(val))

	w.logger.Debugw("Observation metric incremented",
		"metric", w.metricName,
		"value", newValue,
		"delta", val,
		"labels", w.labels,
	)

	if w.publisher != nil {
		// Publish the delta (val), not the cumulative value, since Beholder counters are cumulative
		w.publisher.PublishMetric(context.Background(), w.metricName, val, w.labels)
	}
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
	wrapped := r.maybeWrapCollector(c)
	return r.base.Register(wrapped)
}

func (r *interceptingRegisterer) MustRegister(cs ...prometheus.Collector) {
	wrapped := make([]prometheus.Collector, len(cs))
	for i, c := range cs {
		wrapped[i] = r.maybeWrapCollector(c)
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
		// Try to get counter type
		if counter, ok := c.(prometheus.Counter); ok {
			// Check the descriptor's fqName to see if it matches our target metrics
			descString := desc.String()

			if strings.Contains(descString, "ocr3_sent_observations_total") {
				r.collector.logger.Info("Wrapping ocr3_sent_observations_total counter")
				wrapped := &wrappedCounter{
					Counter:    counter,
					metricName: "ocr3_sent_observations_total",
					labels:     r.collector.beholderLabels,
					publisher:  r.collector.publisher,
					logger:     r.collector.logger,
				}
				r.collector.sentObservationsCounter = wrapped
				return wrapped
			}

			if strings.Contains(descString, "ocr3_included_observations_total") {
				r.collector.logger.Info("Wrapping ocr3_included_observations_total counter")
				wrapped := &wrappedCounter{
					Counter:    counter,
					metricName: "ocr3_included_observations_total",
					labels:     r.collector.beholderLabels,
					publisher:  r.collector.publisher,
					logger:     r.collector.logger,
				}
				r.collector.includedObservationsCounter = wrapped
				return wrapped
			}
		}
	}

	// Not a metric we care about, return as-is
	return c
}
