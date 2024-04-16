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

	trigger := triggers.NewMercuryTriggerService(0, d.logger)
	err = d.registry.Add(context.Background(), trigger)
	if err != nil {
		d.logger.Errorw("could not add mercury trigger to registry", err)
	} else {
		go mercuryEventLoop(trigger, d.logger)
	}

	cfg := Config{
		Lggr:       d.logger,
		Spec:       spec.WorkflowSpec.Workflow,
		WorkflowID: spec.WorkflowSpec.WorkflowID,
		Registry:   d.registry,
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

func emitReports(logger logger.Logger, trigger *triggers.MercuryTriggerService, t int64, prices []int64) ([]mercury.FeedReport, error) {
	reports := []mercury.FeedReport{
		{
			FeedID:               "0x1111111111111111111100000000000000000000000000000000000000000000",
			FullReport:           []byte(fmt.Sprintf(`{ "feed": "ETH", "price": %d }`, prices[0])),
			BenchmarkPrice:       prices[0],
			ObservationTimestamp: t,
		},
		{
			FeedID:               "0x2222222222222222222200000000000000000000000000000000000000000000",
			FullReport:           []byte(fmt.Sprintf(`{ "feed": "LINK", "price": %d }`, prices[1])),
			BenchmarkPrice:       prices[1],
			ObservationTimestamp: t,
		},
		{
			FeedID:               "0x3333333333333333333300000000000000000000000000000000000000000000",
			FullReport:           []byte(fmt.Sprintf(`{ "feed": "BTC", "price": %d }`, prices[2])),
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

	var spec job.WorkflowSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on job: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return jb, err
	}

	jb.WorkflowSpec = &spec
	if jb.Type != job.Workflow {
		return jb, fmt.Errorf("unsupported type %s", jb.Type)
	}

	return jb, nil
}
