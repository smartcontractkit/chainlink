package offchainreporting

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type RunResultSaver struct {
	utils.StartStopOnce

	runResults     <-chan pipeline.RunWithResults
	pipelineRunner pipeline.Runner
	done           chan struct{}
	jobID          int32
}

func (r *RunResultSaver) Start() error {
	if !r.OkayToStart() {
		return errors.New("cannot start already started run result saver")
	}
	go func() {
		for {
			select {
			case rr := <-r.runResults:
				logger.Debugw("saving job run", "run", rr.Run.ID, "task results", rr.TaskRunResults)
				if _, err := r.pipelineRunner.InsertFinishedRunWithResults(context.Background(), rr.Run, rr.TaskRunResults); err != nil {
					logger.Errorw(fmt.Sprintf("error inserting finished results for job ID %v", r.jobID), "err", err)
				}
			case <-r.done:
				return
			}
		}
	}()
	return nil
}

func (r *RunResultSaver) Close() error {
	if !r.OkayToStop() {
		return errors.New("cannot close unstarted run result saver")
	}
	r.done <- struct{}{}

	// In the unlikely event that there are remaining runResults to write,
	// drain the channel and save them.
Drain:
	for {
		select {
		case rr := <-r.runResults:
			logger.Debugw("saving job run before exiting", "run", rr.Run.ID, "task results", rr.TaskRunResults)
			if _, err := r.pipelineRunner.InsertFinishedRunWithResults(context.Background(), rr.Run, rr.TaskRunResults); err != nil {
				logger.Errorw(fmt.Sprintf("error inserting finished results for job %v", r.jobID), "err", err)
			}
		default:
			break Drain
		}
	}
	return nil
}
