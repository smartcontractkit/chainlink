package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const hardcodedWorkflow = `
triggers:
  - type: "mercury-trigger"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - type: "offchain_reporting"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: "write_polygon-testnet-mumbai"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - type: "write_ethereum-testnet-sepolia"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

type Delegate struct {
	registry        types.CapabilitiesRegistry
	logger          logger.Logger
	legacyEVMChains legacyevm.LegacyChainContainer
}

var _ job.Delegate = (*Delegate)(nil)

func (d *Delegate) JobType() job.Type {
	return job.Workflow
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}

func (d *Delegate) AfterJobCreated(jb job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, jb job.Job, q pg.Queryer) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	// NOTE: we temporarily do registration inside ServicesForSpec, this will be moved out of job specs in the future
	err := targets.InitializeWrite(d.registry, d.legacyEVMChains, d.logger)
	if err != nil {
		d.logger.Errorw("could not initialize writes", err)
	}

	trigger := triggers.NewMercuryTriggerService(d.logger)
	err = d.registry.Add(context.Background(), trigger)
	if err != nil {
		d.logger.Errorw("could not add mercury trigger to registry", err)
	} else {
		go mercuryEventLoop(trigger, d.logger)
	}

	cfg := Config{
		Lggr:       d.logger,
		Spec:       hardcodedWorkflow,
		Registry:   d.registry,
		WorkflowID: mockedWorkflowID,
	}
	engine, err := NewEngine(cfg)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{engine}, nil
}

func NewDelegate(logger logger.Logger, registry types.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer) *Delegate {
	return &Delegate{logger: logger, registry: registry, legacyEVMChains: legacyEVMChains}
}

func mercuryEventLoop(trigger *triggers.MercuryTriggerService, logger logger.Logger) {
	sleepSec := 60 * time.Second
	ticker := time.NewTicker(sleepSec)
	defer ticker.Stop()

	prices := []int64{300000, 2000, 5000000}

	for range ticker.C {
		for i := range prices {
			prices[i] = prices[i] + 1
		}

		t := time.Now().Round(sleepSec).Unix()
		reports, err := emitReports(logger, trigger, t, prices)
		if err != nil {
			logger.Errorw("failed to process Mercury reports", "err", err, "timestamp", time.Now().Unix(), "payload", reports)
		}
	}
}

func emitReports(logger logger.Logger, trigger *triggers.MercuryTriggerService, t int64, prices []int64) ([]triggers.FeedReport, error) {
	reports := []triggers.FeedReport{
		{
			FeedID:               mercury.FeedID("0x1111111111111111111100000000000000000000000000000000000000000000").Bytes(),
			FullReport:           []byte{},
			BenchmarkPrice:       prices[0],
			ObservationTimestamp: t,
		},
		{
			FeedID:               mercury.FeedID("0x2222222222222222222200000000000000000000000000000000000000000000").Bytes(),
			FullReport:           []byte{},
			BenchmarkPrice:       prices[1],
			ObservationTimestamp: t,
		},
		{
			FeedID:               mercury.FeedID("0x3333333333333333333300000000000000000000000000000000000000000000").Bytes(),
			FullReport:           []byte{},
			BenchmarkPrice:       prices[2],
			ObservationTimestamp: t,
		},
	}

	logger.Infow("New set of Mercury reports", "timestamp", time.Now().Unix(), "payload", reports)
	return reports, trigger.ProcessReport(reports)
}

func ValidatedWorkflowSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}

	if jb.Type != job.Workflow {
		return jb, fmt.Errorf("unsupported type %s", jb.Type)
	}

	return jb, nil
}
