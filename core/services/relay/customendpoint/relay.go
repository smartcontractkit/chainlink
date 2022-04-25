package customendpoint

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// CL Core OCR2 job spec RelayConfig for customendpoint
type RelayConfig struct {
	// The name of custom endpoint. For example, dydx.
	EndpointName string `json:"endpointName"` // required

	// Endpoint specific transmission target. For example, staging/prod bridge names.
	EndpointTarget string `json:"endpointTarget"` // required

	// The identifier of what payload this job sends.
	// For example, ETHUSD represents the ETH-USD price feed.
	PayloadType string `json:"payloadType"` // required

	// Fields specific to Bridge type targets
	BridgeRequestData string `json:"bridgeRequestData"`
	BridgeInputAtKey  string `json:"bridgeInputAtKey"`

	// The multiplier used in the job spec for storing price feed.
	// Must be a positive integer, which is a power of 10.
	// This should be same as the value used in the multiply task
	// of the Observation phase while reporting final result.
	MultiplierUsed int32 `json:"multiplierUsed"`
}

type OCR2Spec struct {
	RelayConfig
	ID          int32
	IsBootstrap bool
}

// Relayer for customendpoint.
// Note that a customendpoint integration doesn't have any associated Chain.
// We are just uploading to some custom endpoint. This relayer is an interface to
// doing that via OCR2. The implementation just has basic functionality needed
// to make OCR2 work, without any associated chain.
type Relayer struct {
	lggr        logger.Logger
	config      config.GeneralConfig
	pipelineORM pipeline.ORM
	clock       utils.Nower
}

func NewRelayer(lggr logger.Logger,
	config config.GeneralConfig,
	pipelineORM pipeline.ORM,
	clock utils.Nower) *Relayer {
	return &Relayer{
		lggr:        lggr,
		config:      config,
		pipelineORM: pipelineORM,
		clock:       clock,
	}
}

func (r *Relayer) Start(context.Context) error {
	return nil
}

func (r *Relayer) Close() error {
	return nil
}

func (r *Relayer) Ready() error {
	return nil
}

func (r *Relayer) Healthy() error {
	return nil
}

type ocr2Provider struct {
	configDigester offchainConfigDigester
	reportCodec    evmreportcodec.ReportCodec
	tracker        *contractTracker
}

// NewOCR2Provider creates a new OCR2ProviderCtx instance.
func (r *Relayer) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (relaytypes.OCR2ProviderCtx, error) {
	var provider ocr2Provider
	spec, ok := s.(OCR2Spec)
	if !ok {
		return &provider, errors.New("unsuccessful cast to 'customendpoint.OCR2Spec'")
	}
	if spec.MultiplierUsed < 1 {
		return &provider, errors.New("invalid multiplierUsed in 'customendpoint.OCR2Spec'")
	}

	digester := offchainConfigDigester{
		EndpointName:   spec.EndpointName,
		EndpointTarget: spec.EndpointTarget,
		PayloadType:    spec.PayloadType,
	}
	codec := evmreportcodec.ReportCodec{}
	tracker := NewTracker(spec, digester, r.lggr, r.pipelineORM, r.config, codec, r.clock)

	if spec.IsBootstrap {
		// Return early if bootstrap node (doesn't require the full OCR2 provider)
		return &ocr2Provider{
			configDigester: digester,
			tracker:        &tracker,
		}, nil
	}

	return &ocr2Provider{
		configDigester: digester,
		reportCodec:    codec,
		tracker:        &tracker,
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
	return p.configDigester
}

func (p ocr2Provider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p ocr2Provider) MedianContract() median.MedianContract {
	return p.tracker
}
