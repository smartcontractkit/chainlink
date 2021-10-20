package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type RunResultSaver struct {
	utils.StartStopOnce

	db             *gorm.DB
	runResults     <-chan pipeline.RunWithResults
	pipelineRunner pipeline.Runner
	done           chan struct{}
	logger         logger.Logger
}

func NewResultRunSaver(db *gorm.DB, runResults <-chan pipeline.RunWithResults, pipelineRunner pipeline.Runner, done chan struct{},
	logger logger.Logger,
) *RunResultSaver {
	return &RunResultSaver{
		db:             db,
		runResults:     runResults,
		pipelineRunner: pipelineRunner,
		done:           done,
		logger:         logger,
	}
}

func (r *RunResultSaver) Start() error {
	return r.StartOnce("RunResultSaver", func() error {
		go gracefulpanic.WrapRecover(func() {
			for {
				select {
				case rr := <-r.runResults:
					r.logger.Infow("RunSaver: saving job run", "run", rr.Run, "task results", rr.TaskRunResults)
					// We do not want save successful TaskRuns as OCR runs very frequently so a lot of records
					// are produced and the successful TaskRuns do not provide value.
					ctx, cancel := postgres.DefaultQueryCtx()
					_, err := r.pipelineRunner.InsertFinishedRun(r.db.WithContext(ctx), rr.Run, rr.TaskRunResults, false)
					cancel()
					if err != nil {
						r.logger.Errorw("error inserting finished results", "err", err)
					}
				case <-r.done:
					return
				}
			}
		})
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
			case rr := <-r.runResults:
				r.logger.Infow("RunSaver: saving job run before exiting", "run", rr.Run, "task results", rr.TaskRunResults)
				ctx, cancel := postgres.DefaultQueryCtx()
				_, err := r.pipelineRunner.InsertFinishedRun(r.db.WithContext(ctx), rr.Run, rr.TaskRunResults, false)
				cancel()
				if err != nil {
					r.logger.Errorw("error inserting finished results", "err", err)
				}
			default:
				return nil
			}
		}
	})
}
