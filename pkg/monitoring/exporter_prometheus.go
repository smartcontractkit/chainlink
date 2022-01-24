package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func NewPrometheusExporterFactory(
	log Logger,
	metrics Metrics,
) ExporterFactory {
	return &prometheusExporterFactory{
		log,
		metrics,
	}
}

type prometheusExporterFactory struct {
	log     Logger
	metrics Metrics
}

func (f *prometheusExporterFactory) NewExporter(
	chainConfig ChainConfig,
	feedConfig FeedConfig,
) (Exporter, error) {
	f.metrics.SetFeedContractMetadata(
		chainConfig.GetChainID(),
		feedConfig.GetContractAddress(),
		feedConfig.GetContractAddress(),
		feedConfig.GetContractStatus(),
		feedConfig.GetContractType(),
		feedConfig.GetName(),
		feedConfig.GetPath(),
		chainConfig.GetNetworkID(),
		chainConfig.GetNetworkName(),
		feedConfig.GetSymbol(),
	)
	p := &prometheusExporter{
		chainConfig,
		feedConfig,
		f.log,
		f.metrics,
		prometheusLabels{},
		sync.Mutex{},
	}
	p.updateLabels(prometheusLabels{
		networkName:     chainConfig.GetNetworkName(),
		networkID:       chainConfig.GetNetworkID(),
		chainID:         chainConfig.GetChainID(),
		feedName:        feedConfig.GetName(),
		feedPath:        feedConfig.GetPath(),
		symbol:          feedConfig.GetSymbol(),
		contractType:    feedConfig.GetContractType(),
		contractStatus:  feedConfig.GetContractStatus(),
		contractAddress: feedConfig.GetContractAddress(),
		feedID:          feedConfig.GetContractAddress(),
	})
	return p, nil
}

type prometheusExporter struct {
	chainConfig ChainConfig
	feedConfig  FeedConfig

	log     Logger
	metrics Metrics

	labels   prometheusLabels
	labelsMu sync.Mutex
}

func (p *prometheusExporter) Export(ctx context.Context, data interface{}) {
	envelope, ok := data.(Envelope)
	if !ok {
		p.log.Errorw("unexpected type for export", "type", fmt.Sprintf("%T", data))
		return
	}
	p.updateLabels(prometheusLabels{
		sender: string(envelope.Transmitter),
	})
	p.metrics.SetFeedContractNativeTokenBalance(
		envelope.NativeTokenBalance,
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
	p.metrics.SetFeedContractLinkBalance(
		envelope.LinkBalance,
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
	p.metrics.SetNodeMetadata(
		p.chainConfig.GetChainID(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
		string(envelope.Transmitter), // oracleName
		string(envelope.Transmitter), // sender
	)
	p.metrics.SetHeadTrackerCurrentHead(
		envelope.BlockNumber,
		p.chainConfig.GetNetworkName(),
		p.chainConfig.GetChainID(),
		p.chainConfig.GetNetworkID(),
	)
	p.metrics.SetOffchainAggregatorAnswers(
		envelope.LatestAnswer,
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractAddress(),
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
	p.metrics.IncOffchainAggregatorAnswersTotal(
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractAddress(),
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
	isLateAnswer := time.Since(envelope.LatestTimestamp).Seconds() > float64(p.feedConfig.GetHeartbeatSec())
	p.metrics.SetOffchainAggregatorAnswerStalled(
		isLateAnswer,
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractAddress(),
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
	p.metrics.SetOffchainAggregatorSubmissionReceivedValues(
		envelope.LatestAnswer,
		p.feedConfig.GetContractAddress(),
		p.feedConfig.GetContractAddress(),
		string(envelope.Transmitter),
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
}

func (p *prometheusExporter) Cleanup(_ context.Context) {
	p.labelsMu.Lock()
	defer p.labelsMu.Unlock()
	for sender := range p.labels.senders {
		p.metrics.Cleanup(
			p.labels.networkName,
			p.labels.networkID,
			p.labels.chainID,
			sender,
			sender,
			p.labels.feedName,
			p.labels.feedPath,
			p.labels.symbol,
			p.labels.contractType,
			p.labels.contractStatus,
			p.labels.contractAddress,
			p.labels.feedID,
		)
	}
}

// Labels

// prometheusLabels is a helper which stores the labels used an instance of this exporter.
// They are useful at Cleanup time, when this exporter needs to delete all the labels it created.
type prometheusLabels struct {
	networkName     string
	networkID       string
	chainID         string
	sender          string
	senders         map[string]struct{} // A set of unique senders!
	feedName        string
	feedPath        string
	symbol          string
	contractType    string
	contractStatus  string
	contractAddress string
	feedID          string
}

func (p *prometheusExporter) updateLabels(newLabels prometheusLabels) {
	p.labelsMu.Lock()
	defer p.labelsMu.Unlock()
	if newLabels.networkName != "" {
		p.labels.networkName = newLabels.networkName
	}
	if newLabels.networkID != "" {
		p.labels.networkID = newLabels.networkID
	}
	if newLabels.chainID != "" {
		p.labels.chainID = newLabels.chainID
	}
	if p.labels.senders == nil {
		p.labels.senders = map[string]struct{}{}
	}
	if newLabels.sender != "" {
		p.labels.senders[newLabels.sender] = struct{}{}
	}
	if newLabels.feedName != "" {
		p.labels.feedName = newLabels.feedName
	}
	if newLabels.feedPath != "" {
		p.labels.feedPath = newLabels.feedPath
	}
	if newLabels.symbol != "" {
		p.labels.symbol = newLabels.symbol
	}
	if newLabels.contractType != "" {
		p.labels.contractType = newLabels.contractType
	}
	if newLabels.contractStatus != "" {
		p.labels.contractStatus = newLabels.contractStatus
	}
	if newLabels.contractAddress != "" {
		p.labels.contractAddress = newLabels.contractAddress
	}
	if newLabels.feedID != "" {
		p.labels.feedID = newLabels.feedID
	}
}
