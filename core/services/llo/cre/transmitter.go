package cre

import (
	"context"
	"errors"
	"strconv"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

const (
	defaultCapabilityName        = "streams-trigger"
	defaultCapabilityVersion     = "2.0.0" // v2 = LLO
	defaultTickerResolutionMs    = 1000
	defaultSendChannelBufferSize = 1000
)

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

type TransmitterConfig struct {
	Logger               logger.Logger                  `json:"-"`
	CapabilitiesRegistry coretypes.CapabilitiesRegistry `json:"-"`
	DonID                uint32                         `json:"-"`

	TriggerCapabilityName        string `json:"triggerCapabilityName"`
	TriggerCapabilityVersion     string `json:"triggerCapabilityVersion"`
	TriggerTickerMinResolutionMs int    `json:"triggerTickerMinResolutionMs"`
	TriggerSendChannelBufferSize int    `json:"triggerSendChannelBufferSize"`
}

var _ Transmitter = &transmitter{}

type transmitter struct {
	services.Service
	eng *services.Engine

	config      TransmitterConfig
	fromAccount ocr2types.Account
}

func (c TransmitterConfig) NewTransmitter() Transmitter {
	return c.newTransmitter(c.Logger)
}

func (c TransmitterConfig) newTransmitter(lggr logger.Logger) *transmitter {
	t := &transmitter{
		config:      c,
		fromAccount: ocr2types.Account(lggr.Name() + strconv.FormatUint(uint64(c.DonID), 10)),
	}
	if t.config.TriggerCapabilityName == "" {
		t.config.TriggerCapabilityName = defaultCapabilityName
	}
	if t.config.TriggerCapabilityVersion == "" {
		t.config.TriggerCapabilityVersion = defaultCapabilityVersion
	}
	if t.config.TriggerTickerMinResolutionMs == 0 {
		t.config.TriggerTickerMinResolutionMs = defaultTickerResolutionMs
	}
	if t.config.TriggerSendChannelBufferSize == 0 {
		t.config.TriggerSendChannelBufferSize = defaultSendChannelBufferSize
	}

	t.Service, t.eng = services.Config{
		Name: "CRETransmitter",
		// TODO(CAPPL-595): Implement Trigger processing based on common:mercury_trigger.go
	}.NewServiceEngine(lggr)

	return t
}

func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	return t.fromAccount, nil
}

func (t *transmitter) Transmit(
	ctx context.Context,
	cd ocr2types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) error {
	switch report.Info.ReportFormat {
	case llotypes.ReportFormatCapabilityTrigger:
	default:
		// NOTE: Silently ignore non-capability format reports here. All
		// channels are broadcast to all transmitters but this transmitter only
		// cares about channels of type ReportFormatCapabilityTrigger
		return nil
	}
	switch report.Info.LifeCycleStage {
	case datastreamsllo.LifeCycleStageProduction:
	default:
		// NOTE: Ignore retirement and staging reports; for now we assume that
		// we only care about sending production reports.
		//
		// Support could be added in future e.g. for verifying blue-green
		// deploys etc.
		return nil
	}
	pbSigs := make([]*capabilitiespb.OCRAttributedOnchainSignature, len(sigs))
	for i, sig := range sigs {
		pbSigs[i] = &capabilitiespb.OCRAttributedOnchainSignature{
			Signer:    uint32(sig.Signer),
			Signature: sig.Signature,
		}
	}
	ev := &capabilitiespb.OCRTriggerEvent{
		ConfigDigest: cd[:],
		SeqNr:        seqNr,
		Report:       report.Report,
		Sigs:         pbSigs,
	}
	return t.processNewEvent(ctx, ev)
}

func (t *transmitter) processNewEvent(ctx context.Context, event *capabilitiespb.OCRTriggerEvent) error {
	// TODO(CAPPL-595): Implement Trigger processing based on common:mercury_trigger.go
	return errors.New("not implemented")
}
