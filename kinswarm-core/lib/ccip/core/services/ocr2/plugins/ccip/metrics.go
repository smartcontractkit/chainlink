package ccip

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	unexpiredCommitRoots = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_unexpired_commit_roots",
		Help: "Number of unexpired commit roots processed by the plugin",
	}, []string{"plugin", "source", "dest"})
	messagesProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_number_of_messages_processed",
		Help: "Number of messages processed by the plugin during different OCR phases",
	}, []string{"plugin", "source", "dest", "ocrPhase"})
	sequenceNumberCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_sequence_number_counter",
		Help: "Sequence number of the last message processed by the plugin",
	}, []string{"plugin", "source", "dest", "ocrPhase"})
	newReportingPluginErrorCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_new_reporting_plugin_error_counter",
		Help: "The count of the number of errors when calling NewReportingPlugin",
	}, []string{"plugin"})
)

type ocrPhase string

const (
	Observation  ocrPhase = "observation"
	Report       ocrPhase = "report"
	ShouldAccept ocrPhase = "shouldAccept"
)

type PluginMetricsCollector interface {
	NumberOfMessagesProcessed(phase ocrPhase, count int)
	NumberOfMessagesBasedOnInterval(phase ocrPhase, seqNrMin, seqNrMax uint64)
	UnexpiredCommitRoots(count int)
	SequenceNumber(phase ocrPhase, seqNr uint64)
	NewReportingPluginError()
}

type pluginMetricsCollector struct {
	pluginName   string
	source, dest string
}

func NewPluginMetricsCollector(pluginLabel string, sourceChainId, destChainId int64) *pluginMetricsCollector {
	return &pluginMetricsCollector{
		pluginName: pluginLabel,
		source:     strconv.FormatInt(sourceChainId, 10),
		dest:       strconv.FormatInt(destChainId, 10),
	}
}

func (p *pluginMetricsCollector) NumberOfMessagesProcessed(phase ocrPhase, count int) {
	messagesProcessed.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase)).
		Set(float64(count))
}

func (p *pluginMetricsCollector) NumberOfMessagesBasedOnInterval(phase ocrPhase, seqNrMin, seqNrMax uint64) {
	messagesProcessed.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase)).
		Set(float64(seqNrMax - seqNrMin + 1))
}

func (p *pluginMetricsCollector) UnexpiredCommitRoots(count int) {
	unexpiredCommitRoots.
		WithLabelValues(p.pluginName, p.source, p.dest).
		Set(float64(count))
}

func (p *pluginMetricsCollector) SequenceNumber(phase ocrPhase, seqNr uint64) {
	// Don't publish price reports
	if seqNr == 0 {
		return
	}

	sequenceNumberCounter.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase)).
		Set(float64(seqNr))
}

func (p *pluginMetricsCollector) NewReportingPluginError() {
	newReportingPluginErrorCounter.
		WithLabelValues(p.pluginName).
		Inc()
}

var (
	// NoopMetricsCollector is a no-op implementation of PluginMetricsCollector
	NoopMetricsCollector PluginMetricsCollector = noop{}
)

type noop struct{}

func (d noop) NumberOfMessagesProcessed(ocrPhase, int) {
}

func (d noop) NumberOfMessagesBasedOnInterval(ocrPhase, uint64, uint64) {
}

func (d noop) UnexpiredCommitRoots(int) {
}

func (d noop) SequenceNumber(ocrPhase, uint64) {
}

func (d noop) NewReportingPluginError() {
}
