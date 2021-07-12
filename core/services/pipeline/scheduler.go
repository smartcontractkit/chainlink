package pipeline

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type scheduler struct {
	ctx          context.Context
	pipeline     *Pipeline
	dependencies map[int]uint
	waiting      uint
	results      map[int]TaskRunResult
	vars         Vars

	taskCh   chan *memoryTaskRun
	resultCh chan TaskRunResult
}

func newScheduler(ctx context.Context, p *Pipeline, vars Vars) *scheduler {
	dependencies := make(map[int]uint, len(p.Tasks))
	var roots []Task

	for id, task := range p.Tasks {
		len := len(task.Inputs())
		dependencies[id] = uint(len)

		// no inputs: this is a root
		if len == 0 {
			roots = append(roots, task)
		}
	}
	s := &scheduler{
		ctx:          ctx,
		pipeline:     p,
		dependencies: dependencies,
		results:      make(map[int]TaskRunResult, len(p.Tasks)),
		vars:         vars,

		// taskCh should never block
		taskCh:   make(chan *memoryTaskRun, len(dependencies)),
		resultCh: make(chan TaskRunResult),
	}

	for _, task := range roots {
		run := &memoryTaskRun{task: task, vars: s.vars.Copy()}

		s.taskCh <- run
		s.waiting++
	}

	return s
}

func (s *scheduler) Run() {
Loop:
	for s.waiting > 0 {
		// we don't "for result in resultCh" because it would stall if the
		// pipeline is completely empty

		var result TaskRunResult
		select {
		case result = <-s.resultCh:
		case <-s.ctx.Done():
			now := time.Now()
			// mark remaining jobs as timeout
			for _, task := range s.pipeline.Tasks {
				if _, ok := s.results[task.ID()]; !ok {
					s.results[task.ID()] = TaskRunResult{
						Task:       task,
						Result:     Result{Error: ErrTimeout},
						CreatedAt:  now, // TODO: more accurate start time
						FinishedAt: now,
					}
				}
			}

			break Loop
		}

		s.waiting--

		// mark job as complete
		s.results[result.Task.ID()] = result

		// store the result in vars
		if result.Result.Error != nil {
			s.vars.Set(result.Task.DotID(), result.Result.Error)
		} else {
			s.vars.Set(result.Task.DotID(), result.Result.Value)
		}

		for _, output := range result.Task.Outputs() {
			id := output.ID()
			s.dependencies[id]--

			// if all dependencies are done, schedule task run
			if s.dependencies[id] == 0 {
				task := s.pipeline.Tasks[id]
				run := &memoryTaskRun{task: task, vars: s.vars.Copy()}

				// fill in the inputs
				for _, i := range task.Inputs() {
					run.inputs = append(run.inputs, input{index: int32(i.OutputIndex()), result: s.results[i.ID()].Result})
				}

				s.taskCh <- run
				s.waiting++
			}
		}

	}
	close(s.taskCh)
}

func (s *scheduler) report(ctx context.Context, result TaskRunResult) {
	select {
	case s.resultCh <- result:
	case <-ctx.Done():
		logger.Errorw("pipeline.scheduler: timed out reporting result", "result", result)
	}
}
