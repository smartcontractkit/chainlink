package core

import (
	"context"
	"fmt"
	"reflect"
	"time"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const pipleineSpec = `
answer [type=sum values=<[ $(val), 2 ]>]
answer;
`

var PipelineRunner = staticPipelineRunnerService{
	staticPipelineRunnerConfig: staticPipelineRunnerConfig{
		spec: pipleineSpec,
		vars: types.Vars{
			Vars: map[string]interface{}{"foo": "baz"},
		},
		options: types.Options{
			MaxTaskDuration: 10 * time.Second,
		},
		taskResults: types.TaskResults([]types.TaskResult{
			{
				TaskValue: types.TaskValue{
					Value: "hello",
				},
				Index: 0,
			},
		}),
	},
}

var _ testtypes.PipelineEvaluator = (*staticPipelineRunnerService)(nil)

type staticPipelineRunnerConfig struct {
	spec        string
	vars        types.Vars
	options     types.Options
	taskResults types.TaskResults
}

type staticPipelineRunnerService struct {
	staticPipelineRunnerConfig
}

func (pr staticPipelineRunnerService) ExecuteRun(ctx context.Context, s string, v types.Vars, o types.Options) (types.TaskResults, error) {
	if s != pr.spec {
		return nil, fmt.Errorf("expected %s but got %s", pr.spec, s)
	}
	if !reflect.DeepEqual(v, pr.vars) {
		return nil, fmt.Errorf("expected %+v but got %+v", pr.vars, v)
	}
	if !reflect.DeepEqual(o, pr.options) {
		return nil, fmt.Errorf("expected %+v but got %+v", pr.options, o)
	}
	return pr.taskResults, nil
}

func (pr staticPipelineRunnerService) Evaluate(ctx context.Context, other types.PipelineRunnerService) error {
	tr, err := pr.ExecuteRun(ctx, pr.spec, pr.vars, pr.options)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}
	if !reflect.DeepEqual(tr, pr.taskResults) {
		return fmt.Errorf("expected TaskResults %+v but got %+v", pr.taskResults, tr)
	}
	return nil
}
