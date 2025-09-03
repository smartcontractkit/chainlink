package ccip

import (
	"context"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	unexpiredCommitRoots = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_unexpired_commit_roots",
		Help: "Number of unexpired commit roots processed by the plugin",
	}, []string{"plugin", "source", "dest"})
	messagesProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_messages_processed",
		Help: "Number of messages processed by the plugin during different OCR phases",
	}, []string{"plugin", "source", "dest", "ocrPhase"})
	maxSequenceNumber = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_max_sequence_number",
		Help: "Sequence number of the last message processed by the plugin",
	}, []string{"plugin", "source_network_name", "dest_network_name", "ocr_phase", "contract_address"})
	newReportingPluginErrorCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_new_reporting_plugin_error_counter",
		Help: "The count of the number of errors when calling NewReportingPlugin",
	}, []string{"plugin"})
	commitLatestRoundId = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_commit_round_id",
		Help: "The latest round ID observed by the commit plugin",
	}, []string{"source_network_name", "dest_network_name", "contract_address", "plugin"})
	execLatestRoundId = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_exec_round_id",
		Help: "The latest round ID observed by the exec plugin",
	}, []string{"source_network_name", "dest_network_name", "contract_address", "plugin"})
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
	SequenceNumber(phase ocrPhase, seqNr uint64, contractAddress string)
	NewReportingPluginError()
	CommitLatestRoundId(contractAddress string, roundId uint64)
	ExecLatestRoundId(contractAddress string, roundId uint64)
}

type pluginMetricsCollector struct {
	pluginName          string
	source, dest        string
	commitLatestRoundId metric.Int64Gauge
	maxSequenceNumber   metric.Int64Gauge
	execLatestRoundId   metric.Int64Gauge
	messagesProcessed   metric.Int64Gauge
}

func NewPluginMetricsCollector(pluginLabel string, sourceChainId, destChainId int64) (*pluginMetricsCollector, error) {
	commitLatestRoundId, err := beholder.GetMeter().Int64Gauge("ccip_commit_round_id")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_commit_round_id gauge: %w", err)
	}
	execLatestRoundId, err := beholder.GetMeter().Int64Gauge("ccip_exec_round_id")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_exec_round_id gauge: %w", err)
	}
	maxSequenceNumber, err := beholder.GetMeter().Int64Gauge("ccip_max_sequence_number")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_max_sequence_number gauge: %w", err)
	}
	messagesProcessed, err := beholder.GetMeter().Int64Gauge("ccip_messages_processed")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_messages_processed gauge: %w", err)
	}

	return &pluginMetricsCollector{
		pluginName:          pluginLabel,
		source:              strconv.FormatInt(sourceChainId, 10),
		dest:                strconv.FormatInt(destChainId, 10),
		commitLatestRoundId: commitLatestRoundId,
		maxSequenceNumber:   maxSequenceNumber,
		execLatestRoundId:   execLatestRoundId,
		messagesProcessed:   messagesProcessed,
	}, nil
}

func (p *pluginMetricsCollector) NumberOfMessagesProcessed(phase ocrPhase, count int) {
	messagesProcessed.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase)).
		Set(float64(count))
	p.messagesProcessed.Record(context.Background(), int64(count), metric.WithAttributes(
		attribute.String("plugin", p.pluginName),
		attribute.String("source_network_name", p.source),
		attribute.String("dest_network_name", p.dest),
		attribute.String("ocr_phase", string(phase)),
	))
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

func (p *pluginMetricsCollector) SequenceNumber(phase ocrPhase, seqNr uint64, contractAddress string) {
	// Don't publish price reports
	if seqNr == 0 {
		return
	}

	maxSequenceNumber.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase), contractAddress).
		Set(float64(seqNr))
	p.maxSequenceNumber.Record(context.Background(), int64(seqNr), metric.WithAttributes(
		attribute.String("plugin", p.pluginName),
		attribute.String("source_network_name", p.source),
		attribute.String("dest_network_name", p.dest),
		attribute.String("ocr_phase", string(phase)),
		attribute.String("contract_address", contractAddress),
	))
}

func (p *pluginMetricsCollector) NewReportingPluginError() {
	newReportingPluginErrorCounter.
		WithLabelValues(p.pluginName).
		Inc()
}

func (p *pluginMetricsCollector) CommitLatestRoundId(contractAddress string, roundId uint64) {
	commitLatestRoundId.
		WithLabelValues(p.source, p.dest, contractAddress, p.pluginName).
		Set(float64(roundId))
	p.commitLatestRoundId.Record(context.Background(), int64(roundId), metric.WithAttributes(
		attribute.String("source_network_name", p.source),
		attribute.String("dest_network_name", p.dest),
		attribute.String("contract_address", contractAddress),
		attribute.String("plugin", p.pluginName),
	))
}

func (p *pluginMetricsCollector) ExecLatestRoundId(contractAddress string, roundId uint64) {
	execLatestRoundId.
		WithLabelValues(p.source, p.dest, contractAddress, p.pluginName).
		Set(float64(roundId))
	p.execLatestRoundId.Record(context.Background(), int64(roundId), metric.WithAttributes(
		attribute.String("source_network_name", p.source),
		attribute.String("dest_network_name", p.dest),
		attribute.String("contract_address", contractAddress),
		attribute.String("plugin", p.pluginName),
	))
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

func (d noop) SequenceNumber(ocrPhase, uint64, string) {
}

func (d noop) NewReportingPluginError() {
}

func (d noop) CommitLatestRoundId(string, uint64) {
}

func (d noop) ExecLatestRoundId(string, uint64) {
}
