package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/commontypes"
)

var _ commontypes.Metrics = (*OracleMetrics)(nil)

type OracleMetrics struct {
}

func (s *OracleMetrics) NewMetricVec(name string, help string, labelNames ...string) (commontypes.MetricVec, error) {
	return &OracleMetricVec{m: promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labelNames)}, nil
}

var _ commontypes.MetricVec = (*OracleMetricVec)(nil)

type OracleMetricVec struct {
	m *prometheus.GaugeVec
}

func (o *OracleMetricVec) GetMetricWith(labels map[string]string) (commontypes.Metric, error) {
	return o.m.GetMetricWith(labels)
}
