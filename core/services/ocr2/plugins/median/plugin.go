package median

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type MedianConfig interface {
	JobPipelineMaxSuccessfulRuns() uint64
	relay.EnvConfig
}

func NewMedianServices(ctx context.Context,
	jb job.Job,
	isNewlyCreatedJob bool,
	relayer loop.Relayer,
	pipelineRunner pipeline.Runner,
	runResults chan pipeline.Run,
	lggr logger.Logger,
	argsNoPlugin libocr2.OracleArgs,
	cfg MedianConfig,
	endpoint commontypes.MonitoringEndpoint,
	errorLog loop.ErrorLog,
) (srvs []job.ServiceCtx, err error) {
	var pluginConfig config.PluginConfig
	err = json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return
	}
	spec := jb.OCR2OracleSpec
	//TODO retry https://smartcontract-it.atlassian.net/browse/BCF-2112
	provider, err := relayer.NewMedianProvider(ctx, types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         spec.ID,
		ContractID:    spec.ContractID,
		New:           isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
	}, types.PluginArgs{
		TransmitterID: spec.TransmitterID.String,
		PluginConfig:  spec.PluginConfig.Bytes(),
	})
	if err != nil {
		return
	}
	srvs = append(srvs, provider)
	argsNoPlugin.ContractTransmitter = provider.ContractTransmitter()
	argsNoPlugin.ContractConfigTracker = provider.ContractConfigTracker()
	argsNoPlugin.OffchainConfigDigester = provider.OffchainConfigDigester()

	abort := func() {
		var mc services.MultiClose
		for i := range srvs {
			mc = append(mc, srvs[i])
		}
		if cerr := mc.Close(); err != nil {
			lggr.Errorw("Error closing unused services", "err", cerr)
		}
	}

	var median loop.PluginMedian
	if cmdName := v2.EnvMedianPlugin.Get(); cmdName != "" {
		ms := NewPluginMedianService(cmdName, lggr, cfg)
		if err = ms.Launch(); err != nil {
			abort()
			return
		}
		median = ms
		srvs = append(srvs, ms)
	} else {
		median = &Plugin{lggr}
	}
	argsNoPlugin.ReportingPluginFactory, err = median.NewMedianPluginFactory(ctx, provider, ocrcommon.NewDataSourceV2(pipelineRunner,
		jb,
		*jb.PipelineSpec,
		lggr,
		runResults,
		endpoint,
	), ocrcommon.NewInMemoryDataSource(pipelineRunner, jb, pipeline.Spec{
		ID:           jb.ID,
		DotDagSource: pluginConfig.JuelsPerFeeCoinPipeline,
		CreatedAt:    time.Now(),
	}, lggr), errorLog)
	if err != nil {
		abort()
		return
	}

	var oracle *libocr2.Oracle
	oracle, err = libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		abort()
		return
	}
	runSaver := ocrcommon.NewResultRunSaver(
		runResults,
		pipelineRunner,
		make(chan struct{}),
		lggr,
		cfg.JobPipelineMaxSuccessfulRuns(),
	)
	srvs = append(srvs, runSaver, job.NewServiceAdapter(oracle))
	if !jb.OCR2OracleSpec.CaptureEATelemetry {
		lggr.Infof("Enhanced EA telemetry is disabled for job %s", jb.Name.ValueOrZero())
	}
	return
}

type Plugin struct {
	Lggr logger.Logger
}

func (m *Plugin) NewMedianPluginFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin median.DataSource, errorLog loop.ErrorLog) (ocrtypes.ReportingPluginFactory, error) {
	var ctxVals loop.ContextValues
	ctxVals.SetValues(ctx)
	lggr := m.Lggr.With(ctxVals.Args()...)
	factory := median.NumericalMedianFactory{
		ContractTransmitter:       provider.MedianContract(),
		DataSource:                dataSource,
		JuelsPerFeeCoinDataSource: juelsPerFeeCoin,
		Logger: logger.NewOCRWrapper(lggr, true, func(msg string) {
			if err := errorLog.SaveError(context.TODO(), msg); err != nil {
				lggr.Errorw("Unable to save error", "err", msg)
			}
		}),
		OnchainConfigCodec: provider.OnchainConfigCodec(),
		ReportCodec:        provider.ReportCodec(),
	}
	return factory, nil
}

var _ services.ServiceCtx = (*medianService)(nil)

type medianService struct {
	utils.StartStopOnce

	lggr    logger.Logger
	cfg     relay.EnvConfig
	cmdName string

	client *plugin.Client
	cp     plugin.ClientProtocol
	loop.PluginMedian
}

func NewPluginMedianService(cmdName string, lggr logger.Logger, cfg relay.EnvConfig) *medianService {
	return &medianService{cmdName: cmdName, lggr: lggr.Named("PluginMedianService"), cfg: cfg}
}

func (m *medianService) Start(ctx context.Context) error {
	return m.StartOnce("PluginMedianService", func() error {
		if m.PluginMedian != nil {
			return nil
		}
		return m.Launch()
	})
}

// Launch launces the plugin, and sets the backing [loop.PluginMedian]. If this is called directly, then Start() will nop.
func (m *medianService) Launch() error {
	cc := loop.PluginMedianClientConfig(m.lggr)
	cc.Cmd = exec.Command(m.cmdName) //nolint:gosec
	relay.SetEnv(cc.Cmd, m.cfg)
	client := plugin.NewClient(cc)
	cp, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to create plugin Client: %w", err)
	}
	abort := func() {
		if cerr := cp.Close(); cerr != nil {
			m.lggr.Errorw("Error closing ClientProtocol", "err", cerr)
		}
		client.Kill()
	}
	i, err := cp.Dispense(loop.PluginMedianName)
	if err != nil {
		abort()
		return fmt.Errorf("failed to Dispense %q plugin: %w", loop.PluginMedianName, err)
	}
	plug, ok := i.(loop.PluginMedian)
	if !ok {
		abort()
		return fmt.Errorf("expected PluginMedian but got %T", i)
	}
	m.client = client
	m.cp = cp
	m.PluginMedian = plug
	return nil
}

func (m *medianService) Close() error {
	return m.StopOnce("PluginMedianService", func() error {
		err := m.cp.Close()
		m.client.Kill()
		return err
	})
}

func (m *medianService) Name() string { return m.lggr.Name() }

func (m *medianService) HealthReport() map[string]error {
	return map[string]error{m.lggr.Name(): m.Healthy()}
}

func (m *medianService) ping() error { return m.cp.Ping() }

func (m *medianService) Ready() error { return m.ping() }

func (m *medianService) Healthy() error { return m.ping() }
