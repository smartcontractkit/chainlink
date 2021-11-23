package solana

import (
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type relaySolana struct{}

func NewRelay() *relaySolana {
	return &relaySolana{}
}

func (r relaySolana) Start() error {
	// No subservices started on relay start, but when the first job is started
	return nil
}

func (r relaySolana) Close() error {
	// TODO: close all subservices
	return nil
}

func (r relaySolana) Ready() error {
	// always ready
	return nil
}

func (r relaySolana) Healthy() error {
	// TODO: only if all subservices are healthy
	return nil
}

type OCR2ServiceConfig struct {
	NodeURL string
	Address string
	JobID   int32
}

func (r relaySolana) NewOCR2Service(c interface{}) relay.OCR2Service {
	// TODO: connect with smartcontractkit/solana-integration impl
	config := c.(OCR2ServiceConfig)
	return &ocr2Service{
		client: rpc.New(config.NodeURL),
		config: config,
	}
}

type ocr2Service struct {
	client *rpc.Client
	config OCR2ServiceConfig
}

func (r ocr2Service) Start() error {
	// TODO: start all needed subservices
	return nil
}

func (r ocr2Service) Close() error {
	// TODO: close all subservices
	return nil
}

func (r ocr2Service) Ready() error {
	// always ready
	return nil
}

func (r ocr2Service) Healthy() error {
	// TODO: only if all subservices are healthy
	return nil
}

func (r ocr2Service) ContractTransmitter() types.ContractTransmitter {
	return nil
}

func (r ocr2Service) ContractConfigTracker() types.ContractConfigTracker {
	return nil
}

func (r ocr2Service) OffchainConfigDigester() types.OffchainConfigDigester {
	return nil
}

func (r ocr2Service) ReportCodec() median.ReportCodec {
	return nil
}

func (r ocr2Service) MedianContract() median.MedianContract {
	return nil
}
