package protocol

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
)

type pacemakerMetrics struct {
	registerer prometheus.Registerer
	epoch      prometheus.Gauge
	leader     prometheus.Gauge
}

func newPacemakerMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) *pacemakerMetrics {
	epoch := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_epoch",
		Help: "The total number of initialized epochs",
	})
	metricshelper.RegisterOrLogError(logger, registerer, epoch, "ocr2_epoch")

	leader := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ocr2_experimental_leader_oid",
		Help: "The leader oracle id",
	})
	metricshelper.RegisterOrLogError(logger, registerer, leader, "ocr2_experimental_leader_oid")

	return &pacemakerMetrics{
		registerer,
		epoch,
		leader,
	}
}

func (pm *pacemakerMetrics) Close() {
	pm.registerer.Unregister(pm.epoch)
	pm.registerer.Unregister(pm.leader)
}

type reportGenerationMetrics struct {
	registerer                prometheus.Registerer
	sentObservationsTotal     prometheus.Counter
	includedObservationsTotal prometheus.Counter
	completedRoundsTotal      prometheus.Counter
	ledCompletedRoundsTotal   prometheus.Counter
}

func newReportGenerationMetrics(registerer prometheus.Registerer,
	logger commontypes.Logger) *reportGenerationMetrics {

	sentObservationsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr2_sent_observations_total",
		Help: "The total number of observations by this oracle sent to the leader. Note that a " +
			"sent observation might not arrive at the leader in time, or not be included in a " +
			"report request for other reasons. This metric is useful for checking an oracle's " +
			"ability to make observations.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, sentObservationsTotal, "ocr2_sent_observations_total")

	includedObservationsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr2_included_observations_total",
		Help: "The total number of observations by this oracle included in a report request " +
			"from the leader. Note that there is no guarantee that the report request will " +
			"actually lead to a report; for instance, because the leader crashes or maliciously " +
			"equivocates to make this oracle believe that an observation was included. This " +
			"metric is useful in comparison with ocr2_sent_observations_total to check whether " +
			"an oracle is able to regularly make observations that are included in report requests.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, includedObservationsTotal, "ocr2_included_observations_total")

	completedRoundsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr2_completed_rounds_total",
		Help: "The total number of rounds completed by this oracle. A round can be completed by " +
			"the oracle's ReportingPlugin deciding to skip the round or by the creation of a report.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, completedRoundsTotal, "ocr2_completed_rounds_total")

	ledCompletedRoundsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr2_led_completed_rounds_total",
		Help: "The total number of completed rounds led by this oracle. This metric is useful " +
			"for checking an oracle's ability to act as leader.",
	})
	metricshelper.RegisterOrLogError(logger, registerer, ledCompletedRoundsTotal, "ocr2_led_completed_rounds_total")

	return &reportGenerationMetrics{
		registerer,
		sentObservationsTotal,
		includedObservationsTotal,
		completedRoundsTotal,
		ledCompletedRoundsTotal,
	}
}

func (m *reportGenerationMetrics) Close() {
	m.registerer.Unregister(m.sentObservationsTotal)
	m.registerer.Unregister(m.includedObservationsTotal)
	m.registerer.Unregister(m.completedRoundsTotal)
	m.registerer.Unregister(m.ledCompletedRoundsTotal)
}
