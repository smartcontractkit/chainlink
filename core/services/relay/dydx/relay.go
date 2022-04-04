package dydx

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
)

type Logger interface {
	Tracef(format string, values ...interface{})
	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Criticalf(format string, values ...interface{})
	Panicf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})
}

// CL Core OCR2 job spec RelayConfig for dydx
type RelayConfig struct {
	EndpointType string `json:"endpointType"` // required
}

type OCR2Spec struct {
	RelayConfig
	ID          int32
	IsBootstrap bool
}

// Relayer for dydx.
// Note that our dydx integration doesn't have any associated Chain.
// We are just uploading to an API endpoint. This relayer is an interface to
// doing that via OCR2. The implementation just has basic functionality needed
// to make OCR2 work, without any associated chain.
type Relayer struct {
	lggr Logger
}

func NewRelayer(lggr Logger) *Relayer {
	return &Relayer{
		lggr: lggr,
	}
}

// Start starts the relayer respecting the given context.
func (r *Relayer) Start(context.Context) error {
	// No subservices started on relay start, but when the first job is started
	return nil
}

// Close will close all open subservices
func (r *Relayer) Close() error {
	return nil
}

func (r *Relayer) Ready() error {
	// always ready
	return nil
}

// Healthy only if all subservices are healthy
func (r *Relayer) Healthy() error {
	return nil
}

type ocr2Provider struct {
	offchainConfigDigester OffchainConfigDigester
	reportCodec            evmreportcodec.ReportCodec
	tracker                *ContractTracker
}

// NewOCR2Provider creates a new OCR2ProviderCtx instance.
func (r *Relayer) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (relaytypes.OCR2ProviderCtx, error) {
	var provider ocr2Provider
	spec, ok := s.(OCR2Spec)
	if !ok {
		return &provider, errors.New("unsuccessful cast to 'dydx.OCR2Spec'")
	}

	offchainConfigDigester := OffchainConfigDigester{
		endpointType: spec.EndpointType,
	}

	contractTracker := NewTracker(spec, offchainConfigDigester, r.lggr)

	if spec.IsBootstrap {
		// Return early if bootstrap node (doesn't require the full OCR2 provider)
		return &ocr2Provider{
			offchainConfigDigester: offchainConfigDigester,
			tracker:                &contractTracker,
		}, nil
	}

	return &ocr2Provider{
		offchainConfigDigester: offchainConfigDigester,
		reportCodec:            evmreportcodec.ReportCodec{},
		tracker:                &contractTracker,
	}, nil
}

func (p *ocr2Provider) Start(context.Context) error {
	return p.tracker.Start()
}

func (p *ocr2Provider) Close() error {
	return p.tracker.Close()
}

func (p ocr2Provider) Ready() error {
	return p.tracker.Ready()
}

func (p ocr2Provider) Healthy() error {
	return p.tracker.Healthy()
}

func (p ocr2Provider) ContractTransmitter() types.ContractTransmitter {
	return p.tracker
}

func (p ocr2Provider) ContractConfigTracker() types.ContractConfigTracker {
	return p.tracker
}

func (p ocr2Provider) OffchainConfigDigester() types.OffchainConfigDigester {
	return p.offchainConfigDigester
}

func (p ocr2Provider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p ocr2Provider) MedianContract() median.MedianContract {
	return p.tracker
}
