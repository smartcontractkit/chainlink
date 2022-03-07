package types

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type relayerAdapter struct {
	relayer Relayer
}

// NewRelayerCtx creates a new RelayerCtx instance using adapter.
func NewRelayerCtx(relayer Relayer) RelayerCtx {
	return &relayerAdapter{
		relayer,
	}
}

// NewOCR2Provider forwards the call to the underlying relayer.NewOCR2Provider().
func (a relayerAdapter) NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2ProviderCtx, error) {
	provider, err := a.relayer.NewOCR2Provider(externalJobID, spec)
	if err != nil {
		return nil, err
	}
	return NewOCR2ProviderCtx(provider), nil
}

// Start forwards the call to the underlying relayer.Start().
// Context is not used in this case.
func (a relayerAdapter) Start(context.Context) error {
	return a.relayer.Start()
}

// Close forwards the call to the underlying relayer.Close().
func (a relayerAdapter) Close() error {
	return a.relayer.Close()
}

// Ready forwards the call to the underlying relayer.Ready().
func (a relayerAdapter) Ready() error {
	return a.relayer.Ready()
}

// Healthy forwards the call to the underlying relayer.Healthy().
func (a relayerAdapter) Healthy() error {
	return a.relayer.Healthy()
}

type ocr2ProviderAdapter struct {
	provider OCR2Provider
}

// NewOCR2ProviderCtx creates a new OCR2ProviderCtx instance using adapter.
func NewOCR2ProviderCtx(provider OCR2Provider) OCR2ProviderCtx {
	return &ocr2ProviderAdapter{
		provider,
	}
}

// Start forwards the call to the underlying provider.Start().
func (o ocr2ProviderAdapter) Start(context.Context) error {
	return o.provider.Start()
}

// Close forwards the call to the underlying provider.Close().
func (o ocr2ProviderAdapter) Close() error {
	return o.provider.Close()
}

// Ready forwards the call to the underlying provider.Ready().
func (o ocr2ProviderAdapter) Ready() error {
	return o.provider.Ready()
}

// Healthy forwards the call to the underlying provider.Healthy().
func (o ocr2ProviderAdapter) Healthy() error {
	return o.provider.Healthy()
}

// ContractTransmitter forwards the call to the underlying provider.ContractTransmitter().
func (o ocr2ProviderAdapter) ContractTransmitter() ocrtypes.ContractTransmitter {
	return o.provider.ContractTransmitter()
}

// ContractConfigTracker forwards the call to the underlying provider.ContractConfigTracker().
func (o ocr2ProviderAdapter) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return o.provider.ContractConfigTracker()
}

// OffchainConfigDigester forwards the call to the underlying provider.OffchainConfigDigester().
func (o ocr2ProviderAdapter) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return o.provider.OffchainConfigDigester()
}

// ReportCodec forwards the call to the underlying provider.ReportCodec().
func (o ocr2ProviderAdapter) ReportCodec() median.ReportCodec {
	return o.provider.ReportCodec()
}

// MedianContract forwards the call to the underlying provider.MedianContract().
func (o ocr2ProviderAdapter) MedianContract() median.MedianContract {
	return o.provider.MedianContract()
}
