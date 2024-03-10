package mocks

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type MockRunner struct {
	Trrs pipeline.TaskRunResults
	Err  error
}

func (m *MockRunner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error) {
	return &pipeline.Run{ID: 42}, m.Trrs, m.Err
}

var _ pipeline.Task = &MockTask{}

type MockTask struct {
	result pipeline.Result
}

func (m *MockTask) Type() pipeline.TaskType { return "MockTask" }
func (m *MockTask) ID() int                 { return 0 }
func (m *MockTask) DotID() string           { return "" }
func (m *MockTask) Run(ctx context.Context, lggr logger.Logger, vars pipeline.Vars, inputs []pipeline.Result) (pipeline.Result, pipeline.RunInfo) {
	return m.result, pipeline.RunInfo{}
}
func (m *MockTask) Base() *pipeline.BaseTask           { return nil }
func (m *MockTask) Outputs() []pipeline.Task           { return nil }
func (m *MockTask) Inputs() []pipeline.TaskDependency  { return nil }
func (m *MockTask) OutputIndex() int32                 { return 0 }
func (m *MockTask) TaskTimeout() (time.Duration, bool) { return 0, false }
func (m *MockTask) TaskRetries() uint32                { return 0 }
func (m *MockTask) TaskMinBackoff() time.Duration      { return 0 }
func (m *MockTask) TaskMaxBackoff() time.Duration      { return 0 }
