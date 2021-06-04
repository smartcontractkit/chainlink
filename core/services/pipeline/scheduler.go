package pipeline

type scheduler struct {
	pipeline     *Pipeline
	dependencies map[int]uint
	input        interface{}
	waiting      uint
	results      map[int]TaskRunResult
	vars         Vars

	taskCh   chan *memoryTaskRun
	resultCh chan TaskRunResult
}

func newScheduler(p *Pipeline, pipelineInput interface{}) *scheduler {
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
		pipeline:     p,
		dependencies: dependencies,
		input:        pipelineInput,
		results:      make(map[int]TaskRunResult, len(p.Tasks)),
		vars:         NewVarsFrom(map[string]interface{}{"input": pipelineInput}),

		// taskCh should never block
		taskCh:   make(chan *memoryTaskRun, len(dependencies)),
		resultCh: make(chan TaskRunResult),
	}

	for _, task := range roots {
		run := &memoryTaskRun{task: task, vars: s.vars.Copy()}
		// fill in the inputs
		run.inputs = append(run.inputs, input{index: 0, result: Result{Value: s.input}})

		s.taskCh <- run
		s.waiting++
	}

	return s
}

func (s *scheduler) Run() {
	for s.waiting > 0 {
		// we don't "for result in resultCh" because it would stall if the
		// pipeline is completely empty
		result := <-s.resultCh
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
