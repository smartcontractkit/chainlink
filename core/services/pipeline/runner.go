package pipeline

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	// "github.com/smartcontractkit/chainlink/core/services"
	// "github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	Runner interface {
		Start()
		Stop()
		CreateRun(specID int64) error
	}

	runner struct {
		// processTasks services.SleeperTask
		orm    ORM
		config Config
		chStop chan struct{}
		chDone chan struct{}
	}
)

func NewRunner(orm ORM, config Config) *runner {
	r := &runner{
		orm:    orm,
		config: config,
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

func (r *runner) CreateRun(specID int64) (int64, error) {
	err = r.orm.CreateRun(specID)
	if err != nil {
		return err
	}
	// r.processTasks.WakeUp()
	return nil
}

// NOTE: This could potentially run on another machine in the cluster
func (r *runner) processIncompleteTaskRuns() error {
	for {
		var runID int64
		err := r.orm.WithNextUnclaimedTaskRun(func(ptRun TaskRun, predecessors []TaskRun) Result {
			runID = ptRun.RunID

			inputs := make([]Result, len(predecessors))
			for i, predecessor := range predecessors {
				inputs[i] = predecessor.Result()
			}

			task, err := UnmarshalTask(taskSpec.TaskType, taskSpec.TaskJson.Value, orm, config)
			if err != nil {
				return nil, err
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

		err = r.orm.NotifyCompletion(runID)
		if err != nil {
			return err
		}
	}
}
