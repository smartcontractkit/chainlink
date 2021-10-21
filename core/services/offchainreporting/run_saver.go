package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type RunResultSaver struct {
	utils.StartStopOnce

	db             *sqlx.DB
	runResults     <-chan pipeline.Run
	pipelineRunner pipeline.Runner
	done           chan struct{}
	logger         logger.Logger
}

func NewResultRunSaver(db *sqlx.DB, runResults <-chan pipeline.Run, pipelineRunner pipeline.Runner, done chan struct{},
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
				case run := <-r.runResults:
					r.logger.Infow("RunSaver: saving job run", "run", run)
					// We do not want save successful TaskRuns as OCR runs very frequently so a lot of records
					// are produced and the successful TaskRuns do not provide value.
					_, err := r.pipelineRunner.InsertFinishedRun(r.db, run, false)
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
			case run := <-r.runResults:
				r.logger.Infow("RunSaver: saving job run before exiting", "run", run, "task results")
				_, err := r.pipelineRunner.InsertFinishedRun(r.db, run, false)
				if err != nil {
					r.logger.Errorw("error inserting finished results", "err", err)
				}
			default:
				return nil
			}
		}
	})
}
