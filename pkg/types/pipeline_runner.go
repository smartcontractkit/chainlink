package types

import (
	"context"
	"time"
)

type Vars struct {
	Vars map[string]interface{}
}

type Options struct {
	MaxTaskDuration time.Duration
}

type TaskResult struct {
	ID    string
	Type  string
	Value interface{}
	Error error
	Index int
}

type PipelineRunnerService interface {
	ExecuteRun(ctx context.Context, spec string, vars Vars, options Options) ([]TaskResult, error)
}
