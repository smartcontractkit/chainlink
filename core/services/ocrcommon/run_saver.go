package ocrcommon

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type RunResultSaver struct {
	utils.StartStopOnce

	maxSuccessfulRuns uint64
	runResults        <-chan pipeline.Run
	pipelineRunner    pipeline.Runner
	done              chan struct{}
	logger            logger.Logger
}

func NewResultRunSaver(runResults <-chan pipeline.Run, pipelineRunner pipeline.Runner, done chan struct{},
	logger logger.Logger, maxSuccessfulRuns uint64,
) *RunResultSaver {
	return &RunResultSaver{
		maxSuccessfulRuns: maxSuccessfulRuns,
		runResults:        runResults,
		pipelineRunner:    pipelineRunner,
		done:              done,
		logger:            logger,
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
						r.logger.Debugw("Skipping save of successful run due to MaxSuccessfulRuns=0", "run", run)
						continue
					}
					r.logger.Debugw("RunSaver: saving job run", "run", run)
					// We do not want save successful TaskRuns as OCR runs very frequently so a lot of records
					// are produced and the successful TaskRuns do not provide value.
					if err := r.pipelineRunner.InsertFinishedRun(&run, false); err != nil {
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
				if err := r.pipelineRunner.InsertFinishedRun(&run, false); err != nil {
					r.logger.Errorw("error inserting finished results", "err", err)
				}
			default:
				return nil
			}
		}
	})
}
