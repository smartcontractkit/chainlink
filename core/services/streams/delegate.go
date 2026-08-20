package streams

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type DelegateConfig interface {
	MaxSuccessfulRuns() uint64
	ResultWriteQueueDepth() uint64
}

type Delegate struct {
	lggr     common.Logger
	registry Registry
	runner   ocrcommon.Runner
	cfg      DelegateConfig
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(lggr common.Logger, registry Registry, runner ocrcommon.Runner, cfg DelegateConfig) *Delegate {
	return &Delegate{lggr: common.Named(lggr, "StreamsDelegate"), registry: registry, runner: runner, cfg: cfg}
}

func (d *Delegate) JobType() job.Type {
	return job.Stream
}

func (d *Delegate) BeforeJobCreated(jb job.Job)                {}
func (d *Delegate) AfterJobCreated(jb job.Job)                 {}
func (d *Delegate) BeforeJobDeleted(jb job.Job)                {}
func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

func (d *Delegate) ServicesForSpec(ctx context.Context, jb job.Job) (services []job.ServiceCtx, err error) {
	lggr := common.With(common.Named(d.lggr, fmt.Sprintf("Job.%d", jb.ID)), "jobID", jb.ID)

	rrs := ocrcommon.NewResultRunSaver(d.runner, lggr, d.cfg.MaxSuccessfulRuns(), d.cfg.ResultWriteQueueDepth())
	services = append(services, rrs, &StreamService{
		d.registry,
		jb,
		lggr,
		rrs,
	})
	return services, nil
}

type ResultRunSaver interface {
	Save(run *pipeline.Run)
}

type StreamService struct {
	registry Registry
	jb       job.Job
	lggr     common.Logger
	rrs      ResultRunSaver
}

func (s *StreamService) Start(_ context.Context) error {
	if s.jb.PipelineSpec == nil {
		return errors.New("pipeline spec unexpectedly missing for stream")
	}
	return s.registry.Register(s.jb, s.rrs)
}

func (s *StreamService) Close() error {
	s.registry.Unregister(s.jb.ID)
	return nil
}

func ValidatedStreamSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	r := strings.NewReader(tomlString)
	d := toml.NewDecoder(r)
	d.DisallowUnknownFields()
	err := d.Decode(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	if jb.Type != job.Stream {
		return jb, errors.Errorf("unsupported type: %q", jb.Type)
	}

	// The spec stream ID is optional, but if provided represents the final output of the pipeline run.
	// nodes in the DAG may also contain streamID tags.
	// Every spec must have at least one streamID.
	var streamIDs []StreamID

	if jb.StreamID != nil {
		streamIDs = append(streamIDs, *jb.StreamID)
	}

	for _, t := range jb.Pipeline.Tasks {
		if streamID := t.TaskStreamID(); streamID != nil {
			streamIDs = append(streamIDs, *streamID)
		}
	}

	if len(streamIDs) == 0 {
		return jb, errors.New("no streamID found in spec (must be either specified as top-level key 'streamID' or at least one streamID tag must be provided in the pipeline)")
	}

	return jb, nil
}
