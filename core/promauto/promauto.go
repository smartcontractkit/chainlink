package promauto

// package promauto wraps the upstream promauto so we can control our own registration of metrics.
// This is needed due to this issue: https://github.com/prometheus/client_golang/issues/1017
// In the default implementation, vector metrics will not be gather-able until
// at least one WithLabelValues has been called

// IMPORTANT NOTE
// You should use this package everywhere you need promauto and the only import
// of upstream should be in this file

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
)

var _ prometheus.Registerer = &registry{}
var _ prometheus.Gatherer = &registry{}

var (
	customRegistry                         = NewRegistry()
	CustomRegisterer prometheus.Registerer = customRegistry
	CustomGatherer   prometheus.Gatherer   = customRegistry
)

type registry struct {
	promRegisterer prometheus.Registerer
	promGatherer   prometheus.Gatherer

	collectors   []prometheus.Collector
	collectorsMu sync.RWMutex
}

// NewRegistry returns a new custom registry that wraps the prometheus default
// registry
func NewRegistry() *registry {
	promRegistry := prometheus.NewRegistry()

	collectors := make([]prometheus.Collector, 0)
	return &registry{promRegistry, promRegistry, collectors, sync.RWMutex{}}
}

func (r *registry) add(c ...prometheus.Collector) {
	r.collectorsMu.Lock()
	defer r.collectorsMu.Unlock()
	r.collectors = append(r.collectors, c...)
}

func (r *registry) remove(c prometheus.Collector) {
	r.collectorsMu.Lock()
	defer r.collectorsMu.Unlock()
	for i, c2 := range r.collectors {
		if c == c2 {
			UnstableDeleteAt(slice, i)
			return
		}
	}
}

// Register registers a new Collector to be included in metrics
// collection. It returns an error if the descriptors provided by the
// Collector are invalid or if they — in combination with descriptors of
// already registered Collectors — do not fulfill the consistency and
// uniqueness criteria described in the documentation of metric.Desc.
//
// If the provided Collector is equal to a Collector already registered
// (which includes the case of re-registering the same Collector), the
// returned error is an instance of AlreadyRegisteredError, which
// contains the previously registered Collector.
//
// A Collector whose Describe method does not yield any Desc is treated
// as unchecked. Registration will always succeed. No check for
// re-registering (see previous paragraph) is performed. Thus, the
// caller is responsible for not double-registering the same unchecked
// Collector, and for providing a Collector that will not cause
// inconsistent metrics on collection. (This would lead to scrape
// errors.)
func (r *registry) Register(c prometheus.Collector) error {
	r.add(c)
	return errors.WithStack(r.promRegisterer.Register(c))
}

// MustRegister works like Register but registers any number of
// Collectors and panics upon the first registration that causes an
// error.
func (r *registry) MustRegister(c ...prometheus.Collector) {
	r.add(c...)
	r.promRegisterer.MustRegister(c)
}

// Unregister unregisters the Collector that equals the Collector passed
// in as an argument.  (Two Collectors are considered equal if their
// Describe method yields the same set of descriptors.) The function
// returns whether a Collector was unregistered. Note that an unchecked
// Collector cannot be unregistered (as its Describe method does not
// yield any descriptor).
//
// Note that even after unregistering, it will not be possible to
// register a new Collector that is inconsistent with the unregistered
// Collector, e.g. a Collector collecting metrics with the same name but
// a different help string. The rationale here is that the same registry
// instance must only collect consistent metrics throughout its
// lifetime.
func (r *registry) Unregister(c prometheus.Collector) bool {
	r.remove(c...)
	return errors.WithStack(r.promRegisterer.Unregister(c))
}

// Gather calls the Collect method of the registered Collectors and then
// gathers the collected metrics into a lexicographically sorted slice
// of uniquely named MetricFamily protobufs. Gather ensures that the
// returned slice is valid and self-consistent so that it can be used
// for valid exposition. As an exception to the strict consistency
// requirements described for metric.Desc, Gather will tolerate
// different sets of label names for metrics of the same metric family.
//
// Even if an error occurs, Gather attempts to gather as many metrics as
// possible. Hence, if a non-nil error is returned, the returned
// MetricFamily slice could be nil (in case of a fatal error that
// prevented any meaningful metric collection) or contain a number of
// MetricFamily protobufs, some of which might be incomplete, and some
// might be missing altogether. The returned error (which might be a
// MultiError) explains the details. Note that this is mostly useful for
// debugging purposes. If the gathered protobufs are to be used for
// exposition in actual monitoring, it is almost always better to not
// expose an incomplete result and instead disregard the returned
// MetricFamily protobufs in case the returned error is non-nil.
func (r *registry) Gather() ([]*dto.MetricFamily, error) {
	panic("not implemented") // TODO: Implement
}

// promauto package wraps the default promauto with additional registration
// hooks to allow easy display of all registered metrics
//
// Please use this package everywhere instead of github.com/prometheus/client_golang/prometheus/promauto

// NewCounter works like the function of the same name in the prometheus package
// but it automatically registers the Counter with the
// CustomRegisterer. If the registration fails, NewCounter panics.
func NewCounter(opts prometheus.CounterOpts) prometheus.Counter {
	return promauto.With(CustomRegisterer).NewCounter(opts)
}

// NewCounterVec works like the function of the same name in the prometheus
// package but it automatically registers the CounterVec with the
// CustomRegisterer. If the registration fails, NewCounterVec
// panics.
func NewCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	return promauto.With(CustomRegisterer).NewCounterVec(opts, labelNames)
}

// NewCounterFunc works like the function of the same name in the prometheus
// package but it automatically registers the CounterFunc with the
// CustomRegisterer. If the registration fails, NewCounterFunc
// panics.
func NewCounterFunc(opts prometheus.CounterOpts, function func() float64) prometheus.CounterFunc {
	return promauto.With(CustomRegisterer).NewCounterFunc(opts, function)
}

// NewGauge works like the function of the same name in the prometheus package
// but it automatically registers the Gauge with the
// CustomRegisterer. If the registration fails, NewGauge panics.
func NewGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	return promauto.With(CustomRegisterer).NewGauge(opts)
}

// NewGaugeVec works like the function of the same name in the prometheus
// package but it automatically registers the GaugeVec with the
// CustomRegisterer. If the registration fails, NewGaugeVec panics.
func NewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	return promauto.With(CustomRegisterer).NewGaugeVec(opts, labelNames)
}

// NewHistogram works like the function of the same name in the prometheus
// package but it automatically registers the Histogram with the
// CustomRegisterer. If the registration fails, NewHistogram panics.
func NewHistogram(opts prometheus.HistogramOpts) prometheus.Histogram {
	return promauto.With(CustomRegisterer).NewHistogram(opts)
}

// NewHistogramVec works like the function of the same name in the prometheus
// package but it automatically registers the HistogramVec with the
// CustomRegisterer. If the registration fails, NewHistogramVec
// panics.
func NewHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	return promauto.With(CustomRegisterer).NewHistogramVec(opts, labelNames)
}
