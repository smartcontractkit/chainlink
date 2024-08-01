package plugin

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"strconv"

	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/runner"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

type pluginFactory struct {
	logProvider        commontypes.LogEventProvider
	events             types.TransmitEventProvider
	blocks             commontypes.BlockSubscriber
	rp                 commontypes.RecoverableProvider
	builder            commontypes.PayloadBuilder
	getter             commontypes.ConditionalUpkeepProvider
	runnable           types.Runnable
	runnerConf         runner.RunnerConfig
	encoder            commontypes.Encoder
	upkeepTypeGetter   types.UpkeepTypeGetter
	workIDGenerator    types.WorkIDGenerator
	upkeepStateUpdater commontypes.UpkeepStateUpdater
	logger             *log.Logger
}

func NewReportingPluginFactory(
	logProvider commontypes.LogEventProvider,
	events types.TransmitEventProvider,
	blocks commontypes.BlockSubscriber,
	rp commontypes.RecoverableProvider,
	builder commontypes.PayloadBuilder,
	getter commontypes.ConditionalUpkeepProvider,
	runnable types.Runnable,
	runnerConf runner.RunnerConfig,
	encoder commontypes.Encoder,
	upkeepTypeGetter types.UpkeepTypeGetter,
	workIDGenerator types.WorkIDGenerator,
	upkeepStateUpdater commontypes.UpkeepStateUpdater,
	logger *log.Logger,
) ocr3types.ReportingPluginFactory[AutomationReportInfo] {
	return &pluginFactory{
		logProvider:        logProvider,
		events:             events,
		blocks:             blocks,
		rp:                 rp,
		builder:            builder,
		getter:             getter,
		runnable:           runnable,
		runnerConf:         runnerConf,
		encoder:            encoder,
		upkeepTypeGetter:   upkeepTypeGetter,
		workIDGenerator:    workIDGenerator,
		upkeepStateUpdater: upkeepStateUpdater,
		logger:             logger,
	}
}

func (factory *pluginFactory) NewReportingPlugin(c ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[AutomationReportInfo], ocr3types.ReportingPluginInfo, error) {
	info := ocr3types.ReportingPluginInfo{
		Name: fmt.Sprintf("Oracle: %d: Automation Plugin Instance w/ Digest '%s'", c.OracleID, c.ConfigDigest),
		Limits: ocr3types.ReportingPluginLimits{
			MaxQueryLength:       0,
			MaxObservationLength: ocr2keepers.MaxObservationLength,
			MaxOutcomeLength:     ocr2keepers.MaxOutcomeLength,
			MaxReportLength:      ocr2keepers.MaxReportLength,
			MaxReportCount:       ocr2keepers.MaxReportCount,
		},
	}

	// decode the off-chain config
	conf, err := config.DecodeOffchainConfig(c.OffchainConfig)
	if err != nil {
		return nil, info, err
	}

	parsed, err := strconv.ParseFloat(conf.TargetProbability, 32)
	if err != nil {
		return nil, info, fmt.Errorf("%w: failed to parse configured probability", err)
	}

	sample, err := sampleFromProbability(conf.TargetInRounds, c.N-c.F, float32(parsed))
	if err != nil {
		return nil, info, fmt.Errorf("%w: failed to create plugin", err)
	}

	factory.logProvider.SetConfig(commontypes.LogEventProviderConfig{
		BlockRate: conf.LogProviderConfig.BlockRate,
		LogLimit:  conf.LogProviderConfig.LogLimit,
	})

	// create the plugin; all services start automatically
	p, err := newPlugin(
		c.ConfigDigest,
		factory.logProvider,
		factory.events,
		factory.blocks,
		factory.rp,
		factory.builder,
		sample,
		factory.getter,
		factory.encoder,
		factory.upkeepTypeGetter,
		factory.workIDGenerator,
		factory.upkeepStateUpdater,
		factory.runnable,
		factory.runnerConf,
		conf,
		c.F,
		factory.logger,
	)
	if err != nil {
		return nil, info, err
	}

	return p, info, nil
}

func sampleFromProbability(rounds, nodes int, probability float32) (sampleRatio, error) {
	var ratio sampleRatio

	if rounds <= 0 {
		return ratio, fmt.Errorf("number of rounds must be greater than 0")
	}

	if nodes <= 0 {
		return ratio, fmt.Errorf("number of nodes must be greater than 0")
	}

	if probability > 1 || probability <= 0 {
		return ratio, fmt.Errorf("probability must be less than 1 and greater than 0")
	}

	r := complex(float64(rounds), 0)
	n := complex(float64(nodes), 0)
	p := complex(float64(probability), 0)

	// calculate the probability that x of total selection collectively will
	// cover all of a selection by all nodes over number of rounds
	g := -1.0 * (p - 1.0)
	x := cmplx.Pow(cmplx.Pow(g, 1.0/r), 1.0/n)
	rat := cmplx.Abs(-1.0 * (x - 1.0))
	rat = math.Round(rat/0.01) * 0.01
	ratio = sampleRatio(float32(rat))

	return ratio, nil
}

type sampleRatio float32

func (r sampleRatio) OfInt(count int) int {
	if count == 0 {
		return 0
	}

	// rounds the result using basic rounding op
	value := math.Round(float64(r) * float64(count))
	if value < 1.0 {
		return 1
	}

	return int(value)
}

func (r sampleRatio) String() string {
	return fmt.Sprintf("%.8f", float32(r))
}
