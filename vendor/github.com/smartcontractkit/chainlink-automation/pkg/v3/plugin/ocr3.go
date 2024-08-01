package plugin

import (
	"context"
	"errors"
	"fmt"
	"log"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2plustypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/plugin/hooks"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/prommetrics"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/random"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/service"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

type AutomationReportInfo struct{}

type ocr3Plugin struct {
	ConfigDigest                ocr2plustypes.ConfigDigest
	ReportEncoder               ocr2keepers.Encoder
	Coordinator                 types.Coordinator
	UpkeepTypeGetter            types.UpkeepTypeGetter
	WorkIDGenerator             types.WorkIDGenerator
	RemoveFromStagingHook       hooks.RemoveFromStagingHook
	RemoveFromMetadataHook      hooks.RemoveFromMetadataHook
	AddToProposalQHook          hooks.AddToProposalQHook
	AddBlockHistoryHook         hooks.AddBlockHistoryHook
	AddFromStagingHook          hooks.AddFromStagingHook
	AddConditionalProposalsHook hooks.AddConditionalProposalsHook
	AddLogProposalsHook         hooks.AddLogProposalsHook
	Services                    []service.Recoverable
	Config                      config.OffchainConfig
	F                           int
	Logger                      *log.Logger
}

func (plugin *ocr3Plugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (ocr2plustypes.Query, error) {
	return nil, nil
}

func (plugin *ocr3Plugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query ocr2plustypes.Query) (ocr2plustypes.Observation, error) {
	plugin.Logger.Printf("inside Observation for seqNr %d", outctx.SeqNr)
	// first round outcome will be nil or empty so no processing should be done
	if outctx.PreviousOutcome != nil || len(outctx.PreviousOutcome) != 0 {
		// Decode the outcome to AutomationOutcome
		automationOutcome, err := ocr2keepersv3.DecodeAutomationOutcome(outctx.PreviousOutcome, plugin.UpkeepTypeGetter, plugin.WorkIDGenerator)
		if err != nil {
			prommetrics.AutomationPluginError.WithLabelValues(prommetrics.PluginStepObservation, prommetrics.PluginErrorTypeDecodeOutcome).Inc()
			return nil, err
		}

		// Execute pre-build hooks
		plugin.RemoveFromStagingHook.RunHook(automationOutcome)
		plugin.RemoveFromMetadataHook.RunHook(automationOutcome)
		plugin.AddToProposalQHook.RunHook(automationOutcome)
	}
	// Create new AutomationObservation
	observation := ocr2keepersv3.AutomationObservation{}

	plugin.AddBlockHistoryHook.RunHook(&observation, ocr2keepersv3.ObservationBlockHistoryLimit)

	if err := plugin.AddLogProposalsHook.RunHook(&observation, ocr2keepersv3.ObservationLogRecoveryProposalsLimit, getRandomKeySource(plugin.ConfigDigest, outctx.SeqNr)); err != nil {
		return nil, err
	}
	if err := plugin.AddConditionalProposalsHook.RunHook(&observation, ocr2keepersv3.ObservationConditionalsProposalsLimit, getRandomKeySource(plugin.ConfigDigest, outctx.SeqNr)); err != nil {
		return nil, err
	}

	// using the OCR seq number for randomness of the performables ordering.
	// high randomness results in expesive ordering, therefore we reduce
	// the range of the randomness by dividing the seq number by 10
	randSrcSeq := outctx.SeqNr / 10

	// The AddFromStagingHook should always be the last hook that is called as it ensures the size constraints of the observation are met
	if err := plugin.AddFromStagingHook.RunHook(&observation, ocr2keepersv3.ObservationPerformablesLimit, getRandomKeySource(plugin.ConfigDigest, randSrcSeq)); err != nil {
		return nil, err
	}
	prommetrics.AutomationPluginPerformables.WithLabelValues(prommetrics.PluginStepResultStore).Set(float64(len(observation.Performable)))

	plugin.Logger.Printf("built an observation in sequence nr %d with %d performables, %d upkeep proposals and %d block history", outctx.SeqNr, len(observation.Performable), len(observation.UpkeepProposals), len(observation.BlockHistory))
	prommetrics.AutomationPluginPerformables.WithLabelValues(prommetrics.PluginStepObservation).Set(float64(len(observation.Performable)))

	// Encode the observation to bytes
	return observation.Encode()
}

func (plugin *ocr3Plugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query ocr2plustypes.Query) (ocr3types.Quorum, error) {
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (plugin *ocr3Plugin) ValidateObservation(outctx ocr3types.OutcomeContext, query ocr2plustypes.Query, ao ocr2plustypes.AttributedObservation) error {
	plugin.Logger.Printf("inside ValidateObservation for seqNr %d", outctx.SeqNr)
	_, err := ocr2keepersv3.DecodeAutomationObservation(ao.Observation, plugin.UpkeepTypeGetter, plugin.WorkIDGenerator)
	return err
}

func (plugin *ocr3Plugin) Outcome(outctx ocr3types.OutcomeContext, query ocr2plustypes.Query, attributedObservations []ocr2plustypes.AttributedObservation) (ocr3types.Outcome, error) {
	plugin.Logger.Printf("inside Outcome for seqNr %d", outctx.SeqNr)
	p := newPerformables(plugin.F+1, ocr2keepersv3.OutcomeAgreedPerformablesLimit, getRandomKeySource(plugin.ConfigDigest, outctx.SeqNr), plugin.Logger)
	c := newCoordinatedBlockProposals(plugin.F+1, ocr2keepersv3.OutcomeSurfacedProposalsRoundHistoryLimit, ocr2keepersv3.OutcomeSurfacedProposalsLimit, getRandomKeySource(plugin.ConfigDigest, outctx.SeqNr), plugin.Logger)

	for _, attributedObservation := range attributedObservations {
		observation, err := ocr2keepersv3.DecodeAutomationObservation(attributedObservation.Observation, plugin.UpkeepTypeGetter, plugin.WorkIDGenerator)
		if err != nil {
			plugin.Logger.Printf("invalid observation from oracle %d in seqNr %d err %v", attributedObservation.Observer, outctx.SeqNr, err)
			prommetrics.AutomationPluginError.WithLabelValues(prommetrics.PluginStepOutcome, prommetrics.PluginErrorTypeInvalidOracleObservation).Inc()
			// Ignore this observation and continue with further observations. It is expected we will get
			// at least f+1 valid observations
			continue
		}

		plugin.Logger.Printf("adding observation from oracle %d in sequence %d with %d performables, %d upkeep proposals and %d block history in seqNr %d",
			attributedObservation.Observer, outctx.SeqNr, len(observation.Performable), len(observation.UpkeepProposals), len(observation.BlockHistory), outctx.SeqNr)

		p.add(observation)
		c.add(observation)
	}

	outcome := ocr2keepersv3.AutomationOutcome{}
	prevOutcome := ocr2keepersv3.AutomationOutcome{}
	if outctx.PreviousOutcome != nil || len(outctx.PreviousOutcome) != 0 {
		// Decode the outcome to AutomationOutcome
		ao, err := ocr2keepersv3.DecodeAutomationOutcome(outctx.PreviousOutcome, plugin.UpkeepTypeGetter, plugin.WorkIDGenerator)
		if err != nil {
			prommetrics.AutomationPluginError.WithLabelValues(prommetrics.PluginStepOutcome, prommetrics.PluginErrorTypeDecodeOutcome).Inc()
			return nil, err
		}
		prevOutcome = ao
	}

	p.set(&outcome)
	// Important to maintain the order here. Performables should be set before creating new proposals
	c.set(&outcome, prevOutcome)

	newProposals := 0
	if len(outcome.SurfacedProposals) > 0 {
		newProposals = len(outcome.SurfacedProposals[0])
	}
	plugin.Logger.Printf("returning outcome with %d performables and %d new proposals in seqNr %d", len(outcome.AgreedPerformables), newProposals, outctx.SeqNr)
	prommetrics.AutomationPluginPerformables.WithLabelValues(prommetrics.PluginStepOutcome).Set(float64(len(outcome.AgreedPerformables)))

	return outcome.Encode()
}

func (plugin *ocr3Plugin) Reports(seqNr uint64, raw ocr3types.Outcome) ([]ocr3types.ReportWithInfo[AutomationReportInfo], error) {
	plugin.Logger.Printf("inside Reports for seqNr %d", seqNr)
	var (
		reports []ocr3types.ReportWithInfo[AutomationReportInfo]
		outcome ocr2keepersv3.AutomationOutcome
		err     error
	)

	if outcome, err = ocr2keepersv3.DecodeAutomationOutcome(raw, plugin.UpkeepTypeGetter, plugin.WorkIDGenerator); err != nil {
		prommetrics.AutomationPluginError.WithLabelValues(prommetrics.PluginStepReports, prommetrics.PluginErrorTypeDecodeOutcome).Inc()
		return nil, err
	}
	plugin.Logger.Printf("creating report from outcome with %d agreed performables; max batch size: %d; report gas limit %d", len(outcome.AgreedPerformables), plugin.Config.MaxUpkeepBatchSize, plugin.Config.GasLimitPerReport)

	toPerform := []ocr2keepers.CheckResult{}
	var gasUsed uint64
	seenUpkeepIDs := make(map[string]bool)

	performablesAdded := 0
	for i, result := range outcome.AgreedPerformables {
		if len(toPerform) >= plugin.Config.MaxUpkeepBatchSize ||
			gasUsed+result.GasAllocated+uint64(plugin.Config.GasOverheadPerUpkeep) > uint64(plugin.Config.GasLimitPerReport) ||
			seenUpkeepIDs[result.UpkeepID.String()] {

			// If report has reached capacity or has existing upkeepID, encode and append this report
			report, err := plugin.getReportFromPerformables(toPerform)
			if err != nil {
				prommetrics.AutomationPluginError.WithLabelValues(prommetrics.PluginStepReports, prommetrics.PluginErrorTypeEncodeReport).Inc()
				prommetrics.AutomationPluginPerformables.WithLabelValues(prommetrics.PluginStepReports).Set(0)
				return reports, fmt.Errorf("error encountered while encoding: %w", err)
			}
			// append to reports and reset collection
			reports = append(reports, report)
			performablesAdded += len(toPerform)
			toPerform = []ocr2keepers.CheckResult{}
			gasUsed = 0
			seenUpkeepIDs = make(map[string]bool)
		}

		// Add the result to current report
		gasUsed += result.GasAllocated + uint64(plugin.Config.GasOverheadPerUpkeep)
		toPerform = append(toPerform, outcome.AgreedPerformables[i])
		seenUpkeepIDs[result.UpkeepID.String()] = true
	}

	// if there are still values to add
	if len(toPerform) > 0 {
		report, err := plugin.getReportFromPerformables(toPerform)
		if err != nil {
			prommetrics.AutomationPluginError.WithLabelValues(prommetrics.PluginStepReports, prommetrics.PluginErrorTypeEncodeReport).Inc()
			prommetrics.AutomationPluginPerformables.WithLabelValues(prommetrics.PluginStepReports).Set(0)
			return reports, fmt.Errorf("error encountered while encoding: %w", err)
		}
		reports = append(reports, report)
		performablesAdded += len(toPerform)
	}

	plugin.Logger.Printf("%d reports created for sequence number %d", len(reports), seqNr)
	prommetrics.AutomationPluginPerformables.WithLabelValues(prommetrics.PluginStepReports).Set(float64(performablesAdded))
	return reports, nil
}

func (plugin *ocr3Plugin) ShouldAcceptAttestedReport(_ context.Context, seqNr uint64, report ocr3types.ReportWithInfo[AutomationReportInfo]) (bool, error) {
	plugin.Logger.Printf("inside ShouldAcceptAttestedReport for seqNr %d", seqNr)
	upkeeps, err := plugin.ReportEncoder.Extract(report.Report)
	if err != nil {
		return false, err
	}

	plugin.Logger.Printf("%d upkeeps found in report for should accept attested for sequence number %d", len(upkeeps), seqNr)

	accept := false
	// If any upkeep can be accepted, then accept
	for _, upkeep := range upkeeps {
		shouldAccept := plugin.Coordinator.Accept(upkeep)
		plugin.Logger.Printf("checking shouldAccept of upkeep '%s', trigger %s in sequence number %d returned %t", upkeep.UpkeepID, upkeep.Trigger, seqNr, shouldAccept)

		if shouldAccept {
			accept = true
		}
	}

	return accept, nil
}

func (plugin *ocr3Plugin) ShouldTransmitAcceptedReport(_ context.Context, seqNr uint64, report ocr3types.ReportWithInfo[AutomationReportInfo]) (bool, error) {
	plugin.Logger.Printf("inside ShouldTransmitAcceptedReport for seqNr %d", seqNr)
	upkeeps, err := plugin.ReportEncoder.Extract(report.Report)
	if err != nil {
		return false, err
	}

	plugin.Logger.Printf("%d upkeeps found in report for should transmit for sequence number %d", len(upkeeps), seqNr)

	transmit := false
	// If any upkeep should be transmitted, then transmit
	for _, upkeep := range upkeeps {
		shouldTransmit := plugin.Coordinator.ShouldTransmit(upkeep)
		plugin.Logger.Printf("checking transmit of upkeep '%s', trigger %s in sequence number %d returned %t", upkeep.UpkeepID, upkeep.Trigger, seqNr, shouldTransmit)

		if shouldTransmit {
			transmit = true
		}
	}

	return transmit, nil
}

func (plugin *ocr3Plugin) Close() error {
	var err error

	for i := range plugin.Services {
		err = errors.Join(err, plugin.Services[i].Close())
	}

	return err
}

// this start function should not block
func (plugin *ocr3Plugin) startServices() {
	for i := range plugin.Services {
		go func(svc service.Recoverable) {
			if err := svc.Start(context.Background()); err != nil {
				plugin.Logger.Printf("error starting plugin services: %s", err)
			}
		}(plugin.Services[i])
	}
}

func (plugin *ocr3Plugin) getReportFromPerformables(toPerform []ocr2keepers.CheckResult) (ocr3types.ReportWithInfo[AutomationReportInfo], error) {
	encoded, err := plugin.ReportEncoder.Encode(toPerform...)
	return ocr3types.ReportWithInfo[AutomationReportInfo]{
		Report: ocr2plustypes.Report(encoded),
	}, err
}

// Generates a randomness source derived from the config and seq # so
// that it's the same across the network for the same round.
// similar key building as libocr transmit selector.
func getRandomKeySource(cd ocr2plustypes.ConfigDigest, seqNr uint64) [16]byte {
	return random.GetRandomKeySource(cd[:], seqNr)
}
