package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type Delegate struct {
	registry types.CapabilitiesRegistry
	logger   logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

func (d *Delegate) JobType() job.Type {
	return job.Workflow
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}

func (d *Delegate) AfterJobCreated(jb job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(jb job.Job, q pg.Queryer) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	engine, err := NewEngine(d.logger, d.registry)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{engine}, nil
}

func NewDelegate(logger logger.Logger, registry types.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer) *Delegate {
	// NOTE: we temporarily do registration inside NewDelegate, this will be moved out of job specs in the future
	_ = targets.InitializeWrite(registry, legacyEVMChains, logger)
	//trigger := triggers.NewOnDemand()
	trigger := triggers.NewMercuryTriggerService(logger)
	registry.Add(context.Background(), trigger)
	//go eventLoop(trigger, logger)
	go mercuryEventLoop(trigger, logger)

	return &Delegate{logger: logger, registry: registry}
}

func mercuryEventLoop(trigger *triggers.MercuryTriggerService, logger logger.Logger) {
	sleepSec := 60
	ticker := time.NewTicker(time.Duration(sleepSec) * time.Second)
	defer ticker.Stop()

	prices := []int64{300000, 2000, 5000000}

	for range ticker.C {
		for i := range prices {
			prices[i] = prices[i] + 1
		}

		reports := []mercury.FeedReport{
			{
				FeedID:               837699011992234352,
				Fullreport:           []byte{},
				BenchmarkPrice:       prices[0],
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               199223435283769901,
				Fullreport:           []byte{},
				BenchmarkPrice:       prices[1],
				ObservationTimestamp: time.Now().Unix(),
			},
			{
				FeedID:               352837699011992234,
				Fullreport:           []byte{},
				BenchmarkPrice:       prices[2],
				ObservationTimestamp: time.Now().Unix(),
			},
		}

		logger.Infow("New set of Mercury reports", "timestamp", time.Now().Unix(), "payload", reports)
		trigger.ProcessReport(reports)
	}
}

func eventLoop(trigger *triggers.OnDemand, logger logger.Logger) {
	sleepSec := 60
	ticker := time.NewTicker(time.Duration(sleepSec) * time.Second)
	defer ticker.Stop()

	prices := []float64{3000.0, 20.0, 50000.0}

	for range ticker.C {
		for i := range prices {
			prices[i] = prices[i] + 0.01
		}
		resp, _ := values.NewMap(map[string]any{
			"0x1111111111111111111100000000000000000000000000000000000000000000": decimal.NewFromFloat(prices[0]),
			"0x2222222222222222222200000000000000000000000000000000000000000000": decimal.NewFromFloat(prices[1]),
			"0x3333333333333333333300000000000000000000000000000000000000000000": decimal.NewFromFloat(prices[2]),
		})
		cr := capabilities.CapabilityResponse{
			Value: resp,
		}
		logger.Infow("New set of Mercury reports", "timestamp", time.Now().Unix(), "payload", resp)
		trigger.FanOutEvent(context.Background(), cr)
	}
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
