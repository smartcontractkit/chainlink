package pipeline

import (
	"context"
	"sort"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"gopkg.in/guregu/null.v4"
)

func (s *scheduler) newMemoryTaskRun(task Task) *memoryTaskRun {
	run := &memoryTaskRun{task: task, vars: s.vars.Copy()}

	// fill in the inputs, fast path for no inputs
	if len(task.Inputs()) != 0 {
		// construct a list of inputs, sorted by OutputIndex
		type input struct {
			index  int32
			result Result
		}
		inputs := make([]input, 0, len(task.Inputs()))
		// NOTE: we could just allocate via make, then assign directly to run.inputs[i.OutputIndex()]
		// if we're confident that indices are within range
		for _, i := range task.Inputs() {
			inputs = append(inputs, input{index: int32(i.OutputIndex()), result: s.results[i.ID()].Result})
		}
		sort.Slice(inputs, func(i, j int) bool {
			return inputs[i].index < inputs[j].index
		})
		run.inputs = make([]Result, len(inputs))
		for i, input := range inputs {
			run.inputs[i] = input.result
		}
	}

	return run
}

type scheduler struct {
	ctx          context.Context
	cancel       context.CancelFunc
	pipeline     *Pipeline
	run          *Run
	dependencies map[int]uint
	waiting      uint
	results      map[int]TaskRunResult
	vars         Vars

	pending bool
	exiting bool

	taskCh   chan *memoryTaskRun
	resultCh chan TaskRunResult
}

func newScheduler(ctx context.Context, p *Pipeline, run *Run, vars Vars) *scheduler {
	dependencies := make(map[int]uint, len(p.Tasks))

	for id, task := range p.Tasks {
		len := len(task.Inputs())
		dependencies[id] = uint(len)
	}

	ctx, cancel := context.WithCancel(ctx)

	s := &scheduler{
		ctx:          ctx,
		cancel:       cancel,
		pipeline:     p,
		run:          run,
		dependencies: dependencies,
		results:      make(map[int]TaskRunResult, len(p.Tasks)),
		vars:         vars,

		// taskCh should never block
		taskCh:   make(chan *memoryTaskRun, len(dependencies)),
		resultCh: make(chan TaskRunResult),
	}

	// if there's results already present on Run, then this is a resumption. Loop over them and fill results table
	s.reconstructResults()

	// immediately schedule all doable tasks
	for id, task := range p.Tasks {
		// skip tasks that are not ready
		if s.dependencies[id] != 0 {
			continue
		}

		// skip finished tasks
		if _, exists := s.results[id]; exists {
			continue
		}

		run := s.newMemoryTaskRun(task)

		logger.Debugw("scheduling task run", "dot_id", task.DotID(), "attempts", run.attempts)

		s.taskCh <- run
		s.waiting++
	}

	return s
}

func (s *scheduler) reconstructResults() {
	// if there's results already present on Run, then this is a resumption. Loop over them and fill results table
	for _, r := range s.run.PipelineTaskRuns {
		task := s.pipeline.ByDotID(r.DotID)

		if task == nil {
			panic("can't find task by dot id")
		}

		if r.IsPending() {
			continue
		}

		result := Result{}

		if r.Error.Valid {
			result.Error = errors.New(r.Error.String)
		}

		if !r.Output.Null {
			result.Value = r.Output.Val
		}

		s.results[task.ID()] = TaskRunResult{
			Task:       task,
			Result:     result,
			CreatedAt:  r.CreatedAt,
			FinishedAt: r.FinishedAt,
		}

		// store the result in vars
		if result.Error != nil {
			s.vars.Set(task.DotID(), result.Error)
		} else {
			s.vars.Set(task.DotID(), result.Value)
		}

		// mark all outputs as complete
		for _, output := range task.Outputs() {
			id := output.ID()
			s.dependencies[id]--
		}
	}
}

func (s *scheduler) Run() {
	for s.waiting > 0 {
		// we don't "for result in resultCh" because it would stall if the
		// pipeline is completely empty

		result := <-s.resultCh
		// var result TaskRunResult
		// select {
		// case result = <-s.resultCh:
		// case <-s.ctx.Done():
		// }

		s.waiting--

		// TODO: this is temporary until task_bridge can return a proper pending result
		if result.Result.Error == ErrPending {
			result.Result = Result{}        // no output, no error
			result.FinishedAt = null.Time{} // not finished
		}

		// retrieve previous attempt count
		result.Attempts = s.results[result.Task.ID()].Attempts

		// only count as an attempt if the job actually ran. If we're exiting then it got cancelled
		if !s.exiting {
			result.Attempts++
		}

		// store task run
		s.results[result.Task.ID()] = result

		// catch the pending state, we will keep the pipeline running until no more progress is made
		if result.IsPending() {
			s.pending = true

			// skip output wrangling because this task isn't actually complete yet
			continue
		}

		// store the result in vars
		if result.Result.Error != nil {
			s.vars.Set(result.Task.DotID(), result.Result.Error)
		} else {
			s.vars.Set(result.Task.DotID(), result.Result.Value)
		}

		// if the task was marked as fail early, and the result is a fail
		if result.Result.Error != nil && result.Task.Base().FailEarly {
			// drain remaining jobs (continue the loop until waiting = 0) then exit
			s.exiting = true
			s.cancel() // cleanup: terminate pending retries

			// mark remaining jobs as cancelled
			s.markRemaining(ErrCancelled)
		}

		if s.exiting {
			// skip scheduling dependencies if we're exiting early
			continue
		}

		// if task hasn't reached it's max retry count yet, we schedule it again
		if result.Attempts < uint(result.Task.TaskRetries()) && result.Result.Error != nil {
			// we immediately increase the in-flight counter so the pipeline doesn't terminate
			// while we wait for the next retry
			s.waiting++

			backoff := backoff.Backoff{
				Factor: 2,
				Min:    result.Task.TaskMinBackoff(),
				Max:    result.Task.TaskMaxBackoff(),
			}

			go func() {
				select {
				case <-s.ctx.Done():
					// report back so the waiting counter gets decreased
					now := time.Now()
					s.report(context.Background(), TaskRunResult{
						Task:       result.Task,
						Result:     Result{Error: ErrCancelled},
						CreatedAt:  now, // TODO: more accurate start time
						FinishedAt: null.TimeFrom(now),
					})
				case <-time.After(backoff.ForAttempt(float64(result.Attempts - 1))): // we subtract 1 because backoff 0-indexes
					// schedule a new attempt
					run := s.newMemoryTaskRun(result.Task)
					run.attempts = result.Attempts
					logger.Debugw("scheduling task run", "dot_id", run.task.DotID(), "attempts", run.attempts)
					s.taskCh <- run
				}
			}()

			// skip scheduling dependencies since it's the task is not complete yet
			continue
		}

		for _, output := range result.Task.Outputs() {
			id := output.ID()
			s.dependencies[id]--

			// if all dependencies are done, schedule task run
			if s.dependencies[id] == 0 {
				task := s.pipeline.Tasks[id]
				run := s.newMemoryTaskRun(task)

				logger.Debugw("scheduling task run", "dot_id", run.task.DotID(), "attempts", run.attempts)
				s.taskCh <- run
				s.waiting++
			}
		}

	}

	close(s.taskCh)
}

func (s *scheduler) markRemaining(err error) {
	now := time.Now()
	for _, task := range s.pipeline.Tasks {
		if _, ok := s.results[task.ID()]; !ok {
			s.results[task.ID()] = TaskRunResult{
				Task:       task,
				Result:     Result{Error: err},
				CreatedAt:  now, // TODO: more accurate start time
				FinishedAt: null.TimeFrom(now),
			}
		}
	}
}

func (s *scheduler) report(ctx context.Context, result TaskRunResult) {
	select {
	case s.resultCh <- result:
	case <-ctx.Done():
		logger.Errorw("pipeline.scheduler: timed out reporting result", "result", result)
	}
}
