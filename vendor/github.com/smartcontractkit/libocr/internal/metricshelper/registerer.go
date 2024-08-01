package metricshelper

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
)

type PrometheusRegistererWrapper struct {
	registerer prometheus.Registerer
	logger     commontypes.Logger
}

func NewPrometheusRegistererWrapper(registerer prometheus.Registerer, logger commontypes.Logger) *PrometheusRegistererWrapper {
	prw := &PrometheusRegistererWrapper{
		registerer: registerer,
		logger:     logger,
	}
	return prw
}

var _ prometheus.Registerer = (*PrometheusRegistererWrapper)(nil)

func (prw *PrometheusRegistererWrapper) Register(collector prometheus.Collector) error {
	prw.logger.Trace("Registering collector", nil)
	if collector == nil {
		return fmt.Errorf("tried to register nil collector")
	}
	if prw.registerer == nil {
		return fmt.Errorf("nil registerer implementation")
	}
	return prw.registerer.Register(collector)
}

func (prw *PrometheusRegistererWrapper) MustRegister(collectors ...prometheus.Collector) {
	prw.logger.Critical("Should use the Register method instead! Registering collectors with MustRegister will panic if not successful", nil)
	prw.registerer.MustRegister(collectors...)
}

func (prw *PrometheusRegistererWrapper) Unregister(collector prometheus.Collector) bool {
	prw.logger.Trace("Unregistering collector", nil)
	if collector == nil {
		prw.logger.Warn("Unregistering nil collector", nil)
		return false
	}
	if prw.registerer == nil {
		prw.logger.Warn("Nil registerer implementation", nil)
		return false
	}
	return prw.registerer.Unregister(collector)
}

func RegisterOrLogError(logger commontypes.Logger,
	registerer prometheus.Registerer,
	collector prometheus.Collector,
	name string,
) {
	if err := registerer.Register(collector); err != nil {
		logger.Error("PrometheusMetrics: Could not register collector",
			commontypes.LogFields{"name": name, "error": err})
	}
}
