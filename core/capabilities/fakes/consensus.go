package fakes

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/requests"
	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// This capability simulates consensus by running the OCR plugin on a single node, without libOCR.
type FakeConsensusConfig struct {
	N                       int
	F                       int
	DeltaRound              int
	BatchSize               int
	OutcomePruningThreshold uint64
	RequestTimeout          time.Duration
}

func DefaultFakeConsensusConfig() FakeConsensusConfig {
	return FakeConsensusConfig{
		N:                       10,
		F:                       3,
		DeltaRound:              5000,
		BatchSize:               100,
		OutcomePruningThreshold: 1000,
		RequestTimeout:          time.Second * 20,
	}
}

type fakeConsensus struct {
	services.Service
	eng *services.Engine

	config      FakeConsensusConfig
	plugin      ocr3types.ReportingPlugin[[]byte]
	transmitter *ocr3.ContractTransmitter
	store       *requests.Store
	cap         capIface
	stats       SimpleStats

	previousOutcome []byte
}

type capIface interface {
	commonCap.ConsensusCapability
	services.Service
}

var _ services.Service = (*fakeConsensus)(nil)
var _ commonCap.ExecutableCapability = (*fakeConsensus)(nil)

const consensusCapID = "offchain_reporting@1.0.0"

func NewFakeConsensus(lggr logger.Logger, config FakeConsensusConfig) (*fakeConsensus, error) {
	rpConfig := ocr3types.ReportingPluginConfig{}
	store := requests.NewStore()

	capability := ocr3.NewCapability(store, clockwork.NewRealClock(), config.RequestTimeout, capabilities.NewAggregator, capabilities.NewEncoder, lggr, 100)

	plugin, err := ocr3.NewReportingPlugin(store, capability, config.BatchSize, rpConfig,
		config.OutcomePruningThreshold, lggr)
	if err != nil {
		return nil, err
	}

	transmitter := ocr3.NewContractTransmitter(lggr, nil, "")
	transmitter.SetCapability(capability)

	fc := &fakeConsensus{
		config:      config,
		plugin:      plugin,
		transmitter: transmitter,
		store:       store,
		cap:         capability,
		stats:       *NewSimpleStats(),
	}
	fc.Service, fc.eng = services.Config{
		Name:  "fakeConsensus",
		Start: fc.start,
		Close: fc.close,
	}.NewServiceEngine(lggr)
	return fc, nil
}

func (fc *fakeConsensus) start(ctx context.Context) error {
	ticker := services.TickerConfig{
		Initial:   500 * time.Millisecond,
		JitterPct: 0.0,
	}.NewTicker(time.Duration(fc.config.DeltaRound) * time.Millisecond)
	fc.eng.GoTick(ticker, fc.simulateOCRRound)
	return fc.cap.Start(ctx)
}

func (fc *fakeConsensus) close() error {
	err := fc.cap.Close()
	fc.stats.PrintToStdout("Consensus Capability Stats")
	return err
}

func (fc *fakeConsensus) simulateOCRRound(ctx context.Context) {
	runCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	ocrCtx := ocr3types.OutcomeContext{
		SeqNr:           1,
		PreviousOutcome: fc.previousOutcome,
		Epoch:           1,
		Round:           1,
	}

	query, err := fc.plugin.Query(runCtx, ocrCtx)
	if err != nil {
		fc.eng.Errorw("Error running Query", "error", err)
		return
	}
	fc.eng.Debugw("Query execution complete", "size", len(query))

	now := time.Now()
	observation, err := fc.plugin.Observation(runCtx, ocrCtx, query)
	elapsed := time.Since(now)
	if err != nil {
		fc.eng.Errorw("Error running Observation", "error", err)
		return
	}
	fc.eng.Debugw("Observation execution complete", "size", len(observation), "durationMS", elapsed.Milliseconds())
	fc.stats.UpdateMaxStat("Max observation size (bytes)", int64(len(observation)))
	fc.stats.UpdateMaxStat("Max observation duration (ms)", elapsed.Milliseconds())

	aos := []types.AttributedObservation{}
	oracleID := uint8(0)
	for range 2*fc.config.F + 1 {
		aos = append(aos, types.AttributedObservation{
			Observation: observation,
			Observer:    commontypes.OracleID(oracleID),
		})
		oracleID++
	}

	now = time.Now()
	outcome, err := fc.plugin.Outcome(runCtx, ocrCtx, query, aos)
	elapsed = time.Since(now)
	if err != nil {
		fc.eng.Errorw("Error running Outcome", "error", err)
		return
	}
	fc.eng.Debugw("Outcome execution complete", "size", len(outcome), "durationMS", elapsed.Milliseconds())
	fc.stats.UpdateMaxStat("Max outcome size (bytes)", int64(len(outcome)))
	fc.stats.UpdateMaxStat("Max outcome duration (ms)", elapsed.Milliseconds())

	fc.previousOutcome = outcome

	// calculate the size of the previous outcome minus current reports,
	// which is the data that will be preserved as long as a workflow exists
	// in the system
	unmarshaled := &pbtypes.Outcome{}
	_ = proto.Unmarshal(outcome, unmarshaled)
	unmarshaled.Outcomes = nil
	rawCleaned, _ := proto.Marshal(unmarshaled)
	fc.stats.UpdateMaxStat("Max preserved outcome size (bytes)", int64(len(rawCleaned)))

	now = time.Now()
	reports, err := fc.plugin.Reports(runCtx, 1, outcome)
	elapsed = time.Since(now)
	if err != nil {
		fc.eng.Errorw("Error running Reports", "error", err)
		return
	}

	for _, report := range reports {
		reportSize := len(report.ReportWithInfo.Report)
		fc.eng.Infow("Sending report", "report", report, "size", reportSize)
		fc.stats.UpdateMaxStat("Max report size (bytes)", int64(reportSize))
		fc.stats.UpdateMaxStat("Max report duration (ms)", elapsed.Milliseconds())
		emptyDigest := [32]byte{}
		// TODO: add non-empty signatures
		err = fc.transmitter.Transmit(runCtx, emptyDigest, 1, report.ReportWithInfo, []types.AttributedOnchainSignature{})
		if err != nil {
			fc.eng.Errorw("Error transmitting report", "error", err)
		}
	}
}

func (fc *fakeConsensus) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	return fc.cap.Execute(ctx, request)
}

func (fc *fakeConsensus) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fc.eng.Infow("Registering to Fake Consensus", "workflowID", request.Metadata.WorkflowID)
	return fc.cap.RegisterToWorkflow(ctx, request)
}

func (fc *fakeConsensus) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fc.eng.Infow("Unegistering from Fake Consensus", "workflowID", request.Metadata.WorkflowID)
	return fc.cap.UnregisterFromWorkflow(ctx, request)
}

func (fc *fakeConsensus) Info(ctx context.Context) (commonCap.CapabilityInfo, error) {
	return commonCap.CapabilityInfo{
		ID:             consensusCapID,
		CapabilityType: commonCap.CapabilityTypeConsensus,
		Description:    "Fake OCR Consensus",
		DON:            &commonCap.DON{},
		IsLocal:        true,
	}, nil
}
