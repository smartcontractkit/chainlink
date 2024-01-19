package ocrcommon

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type Runner interface {
	InsertFinishedRun(run *pipeline.Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error
}

type RunResultSaver struct {
	services.StateMachine

	maxSuccessfulRuns uint64
	runResults        chan *pipeline.Run
	pipelineRunner    Runner
	done              chan struct{}
	logger            logger.Logger
}

func (r *RunResultSaver) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

func (r *RunResultSaver) Name() string { return r.logger.Name() }

func NewResultRunSaver(pipelineRunner Runner,
	logger logger.Logger, maxSuccessfulRuns uint64, resultsWriteDepth uint64,
) *RunResultSaver {
	return &RunResultSaver{
		maxSuccessfulRuns: maxSuccessfulRuns,
		runResults:        make(chan *pipeline.Run, resultsWriteDepth),
		pipelineRunner:    pipelineRunner,
		done:              make(chan struct{}),
		logger:            logger.Named("RunResultSaver"),
	}
}

// Save sends the run on the internal `runResults` channel for saving.
// IMPORTANT: if the `runResults` pipeline is full, the run will be dropped.
func (r *RunResultSaver) Save(run *pipeline.Run) {
	select {
	case r.runResults <- run:
	default:
		r.logger.Warnw("RunSaver: the write queue was full, dropping run")
	}
}

// Start starts RunResultSaver.
func (r *RunResultSaver) Start(context.Context) error {
	return r.StartOnce("RunResultSaver", func() error {
		go func() {
			for {
				select {
				case run := <-r.runResults:
					if !run.HasErrors() && r.maxSuccessfulRuns == 0 {
						// optimisation: don't bother persisting it if we don't need to save successful runs
						r.logger.Tracew("Skipping save of successful run due to MaxSuccessfulRuns=0", "run", run)
						continue
					}
					r.logger.Tracew("RunSaver: saving job run", "run", run)
					// We do not want save successful TaskRuns as OCR runs very frequently so a lot of records
					// are produced and the successful TaskRuns do not provide value.
					if err := r.pipelineRunner.InsertFinishedRun(run, false); err != nil {
						r.logger.Errorw("error inserting finished results", "err", err)
					}
				case <-r.done:
					return
				}
			}
		}()
		return nil
	})
}

func (r *RunResultSaver) Close() error {
	return r.StopOnce("RunResultSaver", func() error {
		r.done <- struct{}{}

		// In the unlikely event that there are remaining runResults to write,
		// drain the channel and save them.
		for {
			select {
			case run := <-r.runResults:
				r.logger.Infow("RunSaver: saving job run before exiting", "run", run)
				if err := r.pipelineRunner.InsertFinishedRun(run, false); err != nil {
					r.logger.Errorw("error inserting finished results", "err", err)
				}
			default:
				return nil
			}
		}
	})
}
