package test

import (
	"context"
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.PipelineRunnerService = (*StaticPipelineRunnerService)(nil)

type StaticPipelineRunnerService struct{}

func (pr *StaticPipelineRunnerService) ExecuteRun(ctx context.Context, s string, v types.Vars, o types.Options) (types.TaskResults, error) {
	if s != spec {
		return nil, fmt.Errorf("expected %s but got %s", spec, s)
	}
	if !reflect.DeepEqual(v, vars) {
		return nil, fmt.Errorf("expected %+v but got %+v", vars, v)
	}
	if !reflect.DeepEqual(o, options) {
		return nil, fmt.Errorf("expected %+v but got %+v", options, o)
	}
	return taskResults, nil
}
