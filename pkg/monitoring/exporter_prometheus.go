package monitoring

import (
	"context"
	"math/big"
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
		feedConfig.GetID(),
		feedConfig.GetID(),
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
		new(big.Int),
		time.Time{},
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
		contractAddress: feedConfig.GetID(),
		feedID:          feedConfig.GetID(),
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

	prevValue     *big.Int
	prevTimestamp time.Time
	prevMu        sync.Mutex
}

type previousTransmissionDetails struct {
}

func (p *prometheusExporter) Export(ctx context.Context, data interface{}) {
	envelope, isEnvelope := data.(Envelope)
	if !isEnvelope {
		return
	}
	p.updateLabels(prometheusLabels{
		sender: string(envelope.Transmitter),
	})
	p.metrics.SetFeedContractLinkBalance(
		envelope.LinkBalance,
		p.feedConfig.GetID(),
		p.feedConfig.GetID(),
		p.chainConfig.GetChainID(),
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

	isLateAnswer := time.Since(envelope.LatestTimestamp).Seconds() > float64(p.feedConfig.GetHeartbeatSec())
	p.metrics.SetOffchainAggregatorAnswerStalled(
		isLateAnswer,
		p.feedConfig.GetID(),
		p.feedConfig.GetID(),
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)

	if !p.isNewTransmission(envelope.LatestAnswer, envelope.LatestTimestamp) {
		return
	}
	// All the metrics below are only updates if there was a fresh
	// transmission since the last chain read.
	p.metrics.SetOffchainAggregatorAnswers(
		envelope.LatestAnswer,
		p.feedConfig.GetID(),
		p.feedConfig.GetID(),
		p.chainConfig.GetChainID(),
		p.feedConfig.GetContractStatus(),
		p.feedConfig.GetContractType(),
		p.feedConfig.GetName(),
		p.feedConfig.GetPath(),
		p.chainConfig.GetNetworkID(),
		p.chainConfig.GetNetworkName(),
	)
	p.metrics.IncOffchainAggregatorAnswersTotal(
		p.feedConfig.GetID(),
		p.feedConfig.GetID(),
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
		p.feedConfig.GetID(),
		p.feedConfig.GetID(),
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

// isNewTransmission considers four cases:
// - old value == new value && old timestamp == new timestap => return false
// - old value != new value && old timestamp == new timestap => This is probably and error since
//   any new transmission updates the timestamp as well, but, to record the observation, we return true.
// - old value != new value && old timestamp != new timestap => return true
// - old value == new value && old timestamp != new timestap => An unlikely case given the
//   high precision of observations but still a valid update. Return true
func (p prometheusExporter) isNewTransmission(value *big.Int, timestamp time.Time) bool {
	p.prevMu.Lock()
	defer p.prevMu.Unlock()
	if value.Cmp(p.prevValue) == 0 && timestamp.Equal(p.prevTimestamp) {
		return false
	}
	p.prevValue = value
	p.prevTimestamp = timestamp
	return true
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
