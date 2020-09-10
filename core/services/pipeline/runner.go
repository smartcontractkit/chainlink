package pipeline

import (
	"database/sql"
	"github.com/golangci/golangci-lint/pkg/result"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	// "github.com/smartcontractkit/chainlink/core/services"
)

type Runner interface {
	Start()
	Stop()
	CreatePipelineRun(pipelineSpecID int64) error
}

type runner struct {
	// processTasks services.SleeperTask
	orm    ORM
	chStop chan struct{}
	chDone chan struct{}
}

func NewRunner(orm ORM) *runner {
	r := &runner{
		orm:    orm,
		chStop: make(chan struct{}),
		chDone: make(chan struct{}),
	}
	// r.processTasks = services.NewSleeperTask(services.SleeperTaskFuncWorker(r.processIncompleteTaskRuns))
	return r
}

func (r *runner) Start() {
	go func() {
		defer close(r.chDone)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.chStop:
				return
			case <-ticker.C:
				r.processIncompleteTaskRuns()
			}
		}
	}()
}

func (r *runner) Stop() {
	close(r.chStop)
	<-r.chDone
}

func (r *runner) CreatePipelineRun(specID int64) (int64, error) {
	err = r.orm.CreatePipelineRun(specID)
	if err != nil {
		return err
	}
	// r.processTasks.WakeUp()
	return nil
}

// NOTE: This could potentially run on another machine in the cluster
func (r *runner) processIncompleteTaskRuns() error {
	for {
		var pipelineRunID int64
		err := r.orm.WithNextUnclaimedTaskRun(func(ptRun PipelineTaskRun, predecessors []PipelineTaskRun) Result {
			pipelineRunID = ptRun.PipelineRunID

			inputs := make([]Result, len(predecessors))
			for i, predecessor := range predecessors {
				inputs[i] = predecessor.Result()
			}

			task, err := UnmarshalTaskJSON(ptRun)
			if err != nil {
				return err
			}

			result := task.Run(inputs)
			if result.Error != nil {
				logger.Errorw("Pipeline task run errored", "error", result.Error)
			}
			return result
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// All task runs complete
			break
		} else if err != nil {
			return err
		}

		err = r.orm.NotifyCompletion(pipelineRunID)
		if err != nil {
			return err
		}
	}
}
