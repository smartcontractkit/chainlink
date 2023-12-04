package dummyplugin

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func NewPlugin(lggr logger.Logger) *Plugin {
	return &Plugin{
		Plugin:               loop.Plugin{Logger: lggr},
		MedianProviderServer: reportingplugins.MedianProviderServer{},
		stop:                 make(services.StopChan),
	}
}

type Plugin struct {
	loop.Plugin
	stop services.StopChan
	reportingplugins.MedianProviderServer
}

func (p *Plugin) NewReportingPluginFactory(
	ctx context.Context,
	config types.ReportingPluginServiceConfig,
	provider types.MedianProvider,
	pipelineRunner types.PipelineRunnerService,
	telemetry types.TelemetryClient,
	errorLog types.ErrorLog,
) (types.ReportingPluginFactory, error) {
	f, err := p.newFactory(ctx, config, provider, pipelineRunner, telemetry, errorLog)
	if err != nil {
		return nil, err
	}
	s := &reportingPluginFactoryService{lggr: p.Logger, ReportingPluginFactory: f}
	p.SubService(s)
	return s, nil
}

func (p *Plugin) newFactory(ctx context.Context, config types.ReportingPluginServiceConfig, provider types.MedianProvider, pipelineRunner types.PipelineRunnerService, telemetry types.TelemetryClient, errorLog types.ErrorLog) (*dummyPluginFactory, error) {

	factory := &dummyPluginFactory{}
	return factory, nil
}

type reportingPluginFactoryService struct {
	services.StateMachine
	lggr logger.Logger
	ocrtypes.ReportingPluginFactory
}

func (r *reportingPluginFactoryService) Name() string { return r.lggr.Name() }

func (r *reportingPluginFactoryService) Start(ctx context.Context) error {
	return r.StartOnce("ReportingPluginFactory", func() error { return nil })
}

func (r *reportingPluginFactoryService) Close() error {
	return r.StopOnce("ReportingPluginFactory", func() error { return nil })
}

func (r *reportingPluginFactoryService) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

var _ ocrtypes.ReportingPluginFactory = (*dummyPluginFactory)(nil)

type dummyPluginFactory struct {
}

var _ ocrtypes.ReportingPlugin = (*dummyPlugin)(nil)

type dummyPlugin struct {
	f                   int
	maxReportLength     int
	offchainConfig      median.OffchainConfig
	onchainConfig       median.OnchainConfig
	contractTransmitter median.MedianContract
}

func (d dummyPlugin) Query(ctx context.Context, timestamp ocrtypes.ReportTimestamp) (ocrtypes.Query, error) {
	return nil, nil
}

func (d dummyPlugin) Observation(ctx context.Context, timestamp ocrtypes.ReportTimestamp, query ocrtypes.Query) (ocrtypes.Observation, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
	return b, nil
}

func (d dummyPlugin) Report(ctx context.Context, timestamp ocrtypes.ReportTimestamp, query ocrtypes.Query, observations []ocrtypes.AttributedObservation) (bool, ocrtypes.Report, error) {
	if !(d.f+1 <= len(observations)) {
		return false, nil, fmt.Errorf("only received %v valid attributed observations, but need at least f+1 (%v)", len(observations), d.f+1)
	}

	var resultTransmissionDetails struct {
		configDigest    ocrtypes.ConfigDigest
		epoch           uint32
		round           uint8
		latestAnswer    *big.Int
		latestTimestamp time.Time
		err             error
	}
	var resultRoundRequested struct {
		configDigest ocrtypes.ConfigDigest
		epoch        uint32
		round        uint8
		err          error
	}

	var subs subprocesses.Subprocesses
	subs.Go(func() {
		resultTransmissionDetails.configDigest,
			resultTransmissionDetails.epoch,
			resultTransmissionDetails.round,
			resultTransmissionDetails.latestAnswer,
			resultTransmissionDetails.latestTimestamp,
			resultTransmissionDetails.err =
			d.contractTransmitter.LatestTransmissionDetails(ctx)
	})
	subs.Go(func() {
		resultRoundRequested.configDigest,
			resultRoundRequested.epoch,
			resultRoundRequested.round,
			resultRoundRequested.err =
			d.contractTransmitter.LatestRoundRequested(ctx, d.offchainConfig.DeltaC)
	})
	subs.Wait()

	if err := multierr.Combine(resultTransmissionDetails.err, resultRoundRequested.err); err != nil {
		return false, nil, fmt.Errorf("error during LatestTransmissionDetails/LatestRoundRequested: %w", err)
	}

	type timestampObservations struct {
		timestamp uint64
		observer  commontypes.OracleID
	}

	timestamps := make([]timestampObservations, 0)
	for _, o := range observations {
		timestamps = append(timestamps, timestampObservations{
			timestamp: binary.LittleEndian.Uint64(o.Observation),
			observer:  o.Observer,
		})
	}

	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].timestamp > timestamps[j].timestamp
	})

	answer := timestamps[len(timestamps)/2]

	if resultTransmissionDetails.latestAnswer.Uint64() < answer.timestamp {
		observers := [32]byte{}
		var reportObservations []*big.Int

		for i, t := range timestamps {
			observers[i] = byte(t.observer)
			reportObservations = append(reportObservations, big.NewInt(int64(t.timestamp)))
		}

		reportBytes, err := reportTypes.Pack(timestamp, observers, reportObservations, big.NewInt(0))
		return true, reportBytes, err
	}

	return false, nil, nil
}

func (d dummyPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp ocrtypes.ReportTimestamp, report ocrtypes.Report) (bool, error) {
	return true, nil
}

func (d dummyPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp ocrtypes.ReportTimestamp, report ocrtypes.Report) (bool, error) {
	return true, nil
}

func (d dummyPlugin) Close() error {
	return nil
}

func (d dummyPluginFactory) NewReportingPlugin(config ocrtypes.ReportingPluginConfig) (ocrtypes.ReportingPlugin, ocrtypes.ReportingPluginInfo, error) {
	return &dummyPlugin{}, ocrtypes.ReportingPluginInfo{}, nil

}

var reportTypes = getReportTypes()

func getReportTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "rawObservers", Type: mustNewType("bytes32")},
		{Name: "observations", Type: mustNewType("int192[]")},
		{Name: "juelsPerFeeCoin", Type: mustNewType("int192")},
	})
}
