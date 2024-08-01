package mercuryshim

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type MercuryReportInfo struct {
	Epoch uint32
	Round uint8
}

type MercuryOCR3OnchainKeyring struct {
	ocr2OnchainKeyring types.OnchainKeyring
}

var _ ocr3types.OnchainKeyring[MercuryReportInfo] = &MercuryOCR3OnchainKeyring{}

func NewMercuryOCR3OnchainKeyring(ocr2OnchainKeyring types.OnchainKeyring) *MercuryOCR3OnchainKeyring {
	return &MercuryOCR3OnchainKeyring{ocr2OnchainKeyring}
}

func (ok *MercuryOCR3OnchainKeyring) MaxSignatureLength() int {
	return ok.ocr2OnchainKeyring.MaxSignatureLength()
}

func (ok *MercuryOCR3OnchainKeyring) Sign(configDigest types.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[MercuryReportInfo]) (signature []byte, err error) {
	return ok.ocr2OnchainKeyring.Sign(
		types.ReportContext{
			types.ReportTimestamp{
				configDigest,
				reportWithInfo.Info.Epoch,
				reportWithInfo.Info.Round,
			},
			[32]byte{},
		},
		reportWithInfo.Report,
	)
}

func (ok *MercuryOCR3OnchainKeyring) PublicKey() types.OnchainPublicKey {
	return ok.ocr2OnchainKeyring.PublicKey()
}

func (ok *MercuryOCR3OnchainKeyring) Verify(pubkey types.OnchainPublicKey, configDigest types.ConfigDigest, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[MercuryReportInfo], sig []byte) bool {
	return ok.ocr2OnchainKeyring.Verify(
		pubkey,
		types.ReportContext{
			types.ReportTimestamp{
				configDigest,
				reportWithInfo.Info.Epoch,
				reportWithInfo.Info.Round,
			},
			[32]byte{},
		},
		reportWithInfo.Report,
		sig,
	)
}

type MercuryOCR3ContractTransmitter struct {
	ocr2ContractTransmitter types.ContractTransmitter
}

var _ ocr3types.ContractTransmitter[MercuryReportInfo] = &MercuryOCR3ContractTransmitter{}

func NewMercuryOCR3ContractTransmitter(ocr2ContractTransmitter types.ContractTransmitter) *MercuryOCR3ContractTransmitter {
	return &MercuryOCR3ContractTransmitter{ocr2ContractTransmitter}
}

func (t *MercuryOCR3ContractTransmitter) Transmit(
	ctx context.Context,
	configDigest types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[MercuryReportInfo],
	aoss []types.AttributedOnchainSignature,
) error {
	return t.ocr2ContractTransmitter.Transmit(
		ctx,
		types.ReportContext{
			types.ReportTimestamp{
				configDigest,
				reportWithInfo.Info.Epoch,
				reportWithInfo.Info.Round,
			},
			[32]byte{},
		},
		reportWithInfo.Report,
		aoss,
	)
}

func (t *MercuryOCR3ContractTransmitter) FromAccount() (types.Account, error) {
	return t.ocr2ContractTransmitter.FromAccount()
}

func ocr3MaxOutcomeLength(maxReportLength int) int {
	return 100 + maxReportLength + maxReportLength/2
}

func ReportingPluginLimits(mercuryPluginLimits ocr3types.MercuryPluginLimits) ocr3types.ReportingPluginLimits {
	return ocr3types.ReportingPluginLimits{
		0,
		mercuryPluginLimits.MaxObservationLength,
		ocr3MaxOutcomeLength(mercuryPluginLimits.MaxReportLength),
		mercuryPluginLimits.MaxReportLength,
		1,
	}
}

type MercuryReportingPlugin struct {
	Config       ocr3types.ReportingPluginConfig
	Logger       loghelper.LoggerWithContext
	Plugin       ocr3types.MercuryPlugin
	PluginLimits ocr3types.MercuryPluginLimits
}

var _ ocr3types.ReportingPlugin[MercuryReportInfo] = &MercuryReportingPlugin{}

type mercuryReportingPluginOutcome struct {
	Epoch        uint32
	Round        uint8
	ShouldReport bool
	Report       []byte
}

func deserializeMercuryReportingPluginOutcome(outcome ocr3types.Outcome) (mercuryReportingPluginOutcome, error) {
	var result mercuryReportingPluginOutcome
	if len(outcome) == 0 {
		return result, nil
	}
	err := json.Unmarshal(outcome, &result)
	if err != nil {
		return mercuryReportingPluginOutcome{}, err
	}

	return result, nil
}

func serializeMercuryReportingPluginOutcome(outcome mercuryReportingPluginOutcome) ocr3types.Outcome {
	serialized, err := json.Marshal(outcome)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return serialized
}

func (p *MercuryReportingPlugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	return nil, nil
}

func (p *MercuryReportingPlugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	p.Logger.Debug("MercuryReportingPlugin: Observation", commontypes.LogFields{
		"seqNr": outctx.SeqNr,
		"epoch": outctx.Epoch, // nolint: staticcheck
		"round": outctx.Round, // nolint: staticcheck
	})

	previousOutcomeDeserialized, err := deserializeMercuryReportingPluginOutcome(outctx.PreviousOutcome)
	if err != nil {
		return nil, err
	}

	//nolint:staticcheck
	observation, err := p.Plugin.Observation(ctx, types.ReportTimestamp{p.Config.ConfigDigest, uint32(outctx.Epoch), uint8(outctx.Round)}, previousOutcomeDeserialized.Report)
	if err != nil {
		return nil, err
	}

	if !(len(observation) <= p.PluginLimits.MaxObservationLength) {
		return nil, fmt.Errorf("MercuryReportingPlugin: underlying plugin returned oversize observation (%v vs %v)", len(observation), p.PluginLimits.MaxObservationLength)
	}

	return observation, nil
}

func (p *MercuryReportingPlugin) ValidateObservation(outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	return nil
}

func (p *MercuryReportingPlugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query types.Query) (ocr3types.Quorum, error) {
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (p *MercuryReportingPlugin) Outcome(outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	p.Logger.Debug("MercuryReportingPlugin: Outcome", commontypes.LogFields{
		"seqNr": outctx.SeqNr,
		"epoch": outctx.Epoch, // nolint: staticcheck
		"round": outctx.Round, // nolint: staticcheck
	})

	previousOutcomeDeserialized, err := deserializeMercuryReportingPluginOutcome(outctx.PreviousOutcome)
	if err != nil {
		return nil, err
	}

	//nolint:staticcheck
	shouldReport, report, err := p.Plugin.Report(types.ReportTimestamp{p.Config.ConfigDigest, uint32(outctx.Epoch), uint8(outctx.Round)}, previousOutcomeDeserialized.Report, aos)
	if err != nil {
		return nil, err
	}

	if !(len(report) <= p.PluginLimits.MaxReportLength) {
		return nil, fmt.Errorf("MercuryReportingPlugin: underlying plugin returned oversize report (%v vs %v)", len(report), p.PluginLimits.MaxReportLength)
	}

	if !shouldReport {
		report = previousOutcomeDeserialized.Report
	}

	//nolint:staticcheck
	outcomeDeserialized := mercuryReportingPluginOutcome{uint32(outctx.Epoch), uint8(outctx.Round), shouldReport, report}
	return serializeMercuryReportingPluginOutcome(outcomeDeserialized), nil
}

func (p *MercuryReportingPlugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[MercuryReportInfo], error) {
	outcomeDeserialized, err := deserializeMercuryReportingPluginOutcome(outcome)
	if err != nil {
		return nil, err
	}

	if outcomeDeserialized.ShouldReport {
		//nolint:staticcheck
		return []ocr3types.ReportWithInfo[MercuryReportInfo]{{
			outcomeDeserialized.Report,
			MercuryReportInfo{
				outcomeDeserialized.Epoch,
				outcomeDeserialized.Round,
			},
		}}, nil
	} else {
		return nil, nil
	}
}

func (p *MercuryReportingPlugin) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[MercuryReportInfo]) (bool, error) {
	return true, nil
}

func (p *MercuryReportingPlugin) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[MercuryReportInfo]) (bool, error) {
	return true, nil
}

func (p *MercuryReportingPlugin) Close() error {
	return p.Plugin.Close()
}
