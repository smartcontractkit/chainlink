package ccip

import (
	"context"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

var (
	unexpiredCommitRoots = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_unexpired_commit_roots",
		Help: "Number of unexpired commit roots processed by the plugin",
	}, []string{"plugin", "source", "dest", "source_network_name", "dest_network_name"})
	messagesProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_number_of_messages_processed",
		Help: "Number of messages processed by the plugin during different OCR phases",
	}, []string{"plugin", "source", "dest", "ocrPhase", "source_network_name", "dest_network_name"})
	maxSequenceNumber = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_max_sequence_number",
		Help: "Sequence number of the last message processed by the plugin",
	}, []string{"plugin", "source", "dest", "ocr_phase", "contract_address", "source_network_name", "dest_network_name"})
	newReportingPluginErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ccip_new_reporting_plugin_error_counter",
		Help: "The count of the number of errors when calling NewReportingPlugin",
	}, []string{"plugin"})
	commitLatestRoundID = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_commit_round_id",
		Help: "The latest round ID observed by the commit plugin",
	}, []string{"source_network_name", "dest_network_name", "contract_address", "plugin"})
	execLatestRoundID = promauto.NewGaugeVec(prometheus.GaugeOpts{
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
	CommitLatestRoundID(contractAddress string, roundID uint64)
	ExecLatestRoundID(contractAddress string, roundID uint64)
}

type pluginMetricsCollector struct {
	pluginName                         string
	source, dest, sourceName, destName string
	bhClient                           beholder.Client
	unexpiredCommitRoots               metric.Int64Gauge
	messagesProcessed                  metric.Int64Gauge
	maxSequenceNumber                  metric.Int64Gauge
	newReportingPluginErrorCounter     metric.Int64Counter
	commitLatestRoundID                metric.Int64Gauge
	execLatestRoundID                  metric.Int64Gauge
}

func NewPluginMetricsCollector(pluginLabel string, bhClient beholder.Client, sourceChainID, destChainID int64, srcChainName string, destChainName string) (*pluginMetricsCollector, error) {
	unexpiredCommitRoots, err := bhClient.Meter.Int64Gauge("ccip_unexpired_commit_roots")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_unexpired_commit_roots gauge: %w", err)
	}
	messagesProcessed, err := bhClient.Meter.Int64Gauge("ccip_number_of_messages_processed")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_messages_processed gauge: %w", err)
	}
	maxSequenceNumber, err := bhClient.Meter.Int64Gauge("ccip_max_sequence_number")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_max_sequence_number gauge: %w", err)
	}
	newReportingPluginErrorCounter, err := bhClient.Meter.Int64Counter("ccip_new_reporting_plugin_error_counter")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_new_reporting_plugin_error_counter counter: %w", err)
	}
	commitLatestRoundID, err := bhClient.Meter.Int64Gauge("ccip_commit_round_id")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_commit_round_id gauge: %w", err)
	}
	execLatestRoundID, err := bhClient.Meter.Int64Gauge("ccip_exec_round_id")
	if err != nil {
		return nil, fmt.Errorf("failed to register ccip_exec_round_id gauge: %w", err)
	}

	return &pluginMetricsCollector{
		pluginName:                     pluginLabel,
		source:                         strconv.FormatInt(sourceChainID, 10),
		dest:                           strconv.FormatInt(destChainID, 10),
		sourceName:                     srcChainName,
		destName:                       destChainName,
		bhClient:                       bhClient,
		unexpiredCommitRoots:           unexpiredCommitRoots,
		messagesProcessed:              messagesProcessed,
		maxSequenceNumber:              maxSequenceNumber,
		newReportingPluginErrorCounter: newReportingPluginErrorCounter,
		commitLatestRoundID:            commitLatestRoundID,
		execLatestRoundID:              execLatestRoundID,
	}, nil
}

func (p *pluginMetricsCollector) NumberOfMessagesProcessed(phase ocrPhase, count int) {
	messagesProcessed.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase), p.sourceName, p.destName).
		Set(float64(count))
	p.messagesProcessed.Record(context.Background(), int64(count), metric.WithAttributes(
		attribute.String("plugin", p.pluginName),
		attribute.String("source", p.source),
		attribute.String("dest", p.dest),
		attribute.String("ocr_phase", string(phase)),
		attribute.String("source_network_name", p.sourceName),
		attribute.String("dest_network_name", p.destName),
	))
}

func (p *pluginMetricsCollector) NumberOfMessagesBasedOnInterval(phase ocrPhase, seqNrMin, seqNrMax uint64) {
	messagesProcessed.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase), p.sourceName, p.destName).
		Set(float64(seqNrMax - seqNrMin + 1))
	p.messagesProcessed.Record(context.Background(), int64(seqNrMax-seqNrMin+1), metric.WithAttributes( //nolint:gosec // Number will not be negative
		attribute.String("plugin", p.pluginName),
		attribute.String("source", p.source),
		attribute.String("dest", p.dest),
		attribute.String("ocr_phase", string(phase)),
		attribute.String("source_network_name", p.sourceName),
		attribute.String("dest_network_name", p.destName),
	))
}

func (p *pluginMetricsCollector) UnexpiredCommitRoots(count int) {
	unexpiredCommitRoots.
		WithLabelValues(p.pluginName, p.source, p.dest, p.sourceName, p.destName).
		Set(float64(count))
	p.unexpiredCommitRoots.Record(context.Background(), int64(count), metric.WithAttributes(
		attribute.String("plugin", p.pluginName),
		attribute.String("source", p.source),
		attribute.String("dest", p.dest),
		attribute.String("source_network_name", p.sourceName),
		attribute.String("dest_network_name", p.destName),
	))
}

func (p *pluginMetricsCollector) SequenceNumber(phase ocrPhase, seqNr uint64, contractAddress string) {
	// Don't publish price reports
	if seqNr == 0 {
		return
	}

	maxSequenceNumber.
		WithLabelValues(p.pluginName, p.source, p.dest, string(phase), contractAddress, p.sourceName, p.destName).
		Set(float64(seqNr))
	p.maxSequenceNumber.Record(context.Background(), int64(seqNr), metric.WithAttributes( //nolint:gosec // Number will not be negative
		attribute.String("plugin", p.pluginName),
		attribute.String("source", p.source),
		attribute.String("dest", p.dest),
		attribute.String("ocr_phase", string(phase)),
		attribute.String("contract_address", contractAddress),
		attribute.String("source_network_name", p.sourceName),
		attribute.String("dest_network_name", p.destName),
	))
}

func (p *pluginMetricsCollector) NewReportingPluginError() {
	newReportingPluginErrorCounter.
		WithLabelValues(p.pluginName).
		Inc()
	p.newReportingPluginErrorCounter.Add(context.Background(), 1, metric.WithAttributes(
		attribute.String("plugin", p.pluginName),
	))
}

func (p *pluginMetricsCollector) CommitLatestRoundID(contractAddress string, roundID uint64) {
	commitLatestRoundID.
		WithLabelValues(p.source, p.dest, contractAddress, p.pluginName).
		Set(float64(roundID))
	p.commitLatestRoundID.Record(context.Background(), int64(roundID), metric.WithAttributes( //nolint:gosec // Number will not be negative
		attribute.String("source_network_name", p.sourceName),
		attribute.String("dest_network_name", p.destName),
		attribute.String("contract_address", contractAddress),
		attribute.String("plugin", p.pluginName),
	))
}

func (p *pluginMetricsCollector) ExecLatestRoundID(contractAddress string, roundID uint64) {
	execLatestRoundID.
		WithLabelValues(p.source, p.dest, contractAddress, p.pluginName).
		Set(float64(roundID))
	p.execLatestRoundID.Record(context.Background(), int64(roundID), metric.WithAttributes( //nolint:gosec // Number will not be negative
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

func (d noop) CommitLatestRoundID(string, uint64) {
}

func (d noop) ExecLatestRoundID(string, uint64) {
}
