package solana

import (
	"errors"

	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ service.Service = (*relayer)(nil)
var _ relay.Relayer = (*relayer)(nil)

type relayer struct {
	keystoreOCR2 keystore.OCR2
}

func NewRelayer(keystoreOCR2 keystore.OCR2) *relayer {
	return &relayer{
		keystoreOCR2: keystoreOCR2,
	}
}

func (r relayer) Start() error {
	// No subservices started on relay start, but when the first job is started
	return nil
}

func (r relayer) Close() error {
	// TODO: close all subservices
	return nil
}

func (r relayer) Ready() error {
	// always ready
	return nil
}

func (r relayer) Healthy() error {
	// TODO: only if all subservices are healthy
	return nil
}

type OCR2ProviderConfig struct {
	NodeURL     string
	Address     string
	JobID       int32
	KeyBundleID string
}

func (r relayer) NewOCR2Provider(c interface{}) (relay.OCR2Provider, error) {
	// TODO: connect with smartcontractkit/solana-integration impl
	config, ok := c.(OCR2ProviderConfig)
	if !ok {
		return nil, errors.New("unsuccessful cast to 'solana.OCR2ProviderConfig'")
	}

	kb, err := r.keystoreOCR2.Get(config.KeyBundleID)
	if err != nil {
		return nil, err
	}

	return &ocr2Provider{
		client:    rpc.New(config.NodeURL),
		keyBundle: kb,
		config:    config,
	}, nil
}

type ocr2Provider struct {
	client    *rpc.Client
	keyBundle ocr2key.KeyBundle
	config    OCR2ProviderConfig
}

func (p ocr2Provider) Start() error {
	// TODO: start all needed subservices
	return nil
}

func (p ocr2Provider) Close() error {
	// TODO: close all subservices
	return nil
}

func (p ocr2Provider) Ready() error {
	// always ready
	return nil
}

func (p ocr2Provider) Healthy() error {
	// TODO: only if all subservices are healthy
	return nil
}

func (p ocr2Provider) OffchainKeyring() types.OffchainKeyring {
	return &p.keyBundle.OffchainKeyring
}

func (p ocr2Provider) OnchainKeyring() types.OnchainKeyring {
	return &p.keyBundle.OnchainKeyring
}

func (p ocr2Provider) ContractTransmitter() types.ContractTransmitter {
	return nil
}

func (p ocr2Provider) ContractConfigTracker() types.ContractConfigTracker {
	return nil
}

func (p ocr2Provider) OffchainConfigDigester() types.OffchainConfigDigester {
	return nil
}

func (p ocr2Provider) ReportCodec() median.ReportCodec {
	return nil
}

func (p ocr2Provider) MedianContract() median.MedianContract {
	return nil
}
