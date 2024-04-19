package test

import (
	"context"
	"fmt"
	"reflect"
	"time"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
)

const pipleineSpec = `
answer [type=sum values=<[ $(val), 2 ]>]
answer;
`

var PipelineRunner = staticPipelineRunnerService{
	staticPipelineRunnerConfig: staticPipelineRunnerConfig{
		spec: pipleineSpec,
		vars: core.Vars{
			Vars: map[string]interface{}{"foo": "baz"},
		},
		options: core.Options{
			MaxTaskDuration: 10 * time.Second,
		},
		taskResults: core.TaskResults([]core.TaskResult{
			{
				TaskValue: core.TaskValue{
					Value: jsonserializable.JSONSerializable{
						Val:   "hello",
						Valid: true,
					},
				},
				Index: 0,
			},
		}),
	},
}

var _ testtypes.PipelineEvaluator = (*staticPipelineRunnerService)(nil)

type staticPipelineRunnerConfig struct {
	spec        string
	vars        core.Vars
	options     core.Options
	taskResults core.TaskResults
}

type staticPipelineRunnerService struct {
	staticPipelineRunnerConfig
}

func (pr staticPipelineRunnerService) ExecuteRun(ctx context.Context, s string, v core.Vars, o core.Options) (core.TaskResults, error) {
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

func (pr staticPipelineRunnerService) Evaluate(ctx context.Context, other core.PipelineRunnerService) error {
	tr, err := pr.ExecuteRun(ctx, pr.spec, pr.vars, pr.options)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}
	if !reflect.DeepEqual(tr, pr.taskResults) {
		return fmt.Errorf("expected TaskResults %+v but got %+v", pr.taskResults, tr)
	}
	return nil
}
