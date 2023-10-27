package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type mockPipelineRunner struct {
	taskResults []types.TaskResult
	err         error
	spec        string
	vars        types.Vars
	options     types.Options
}

func (m *mockPipelineRunner) ExecuteRun(ctx context.Context, spec string, vars types.Vars, options types.Options) ([]types.TaskResult, error) {
	m.spec, m.vars, m.options = spec, vars, options
	return m.taskResults, m.err
}

type clientAdapter struct {
	srv pb.PipelineRunnerServiceServer
}

func (c *clientAdapter) ExecuteRun(ctx context.Context, in *pb.RunRequest, opts ...grpc.CallOption) (*pb.RunResponse, error) {
	return c.srv.ExecuteRun(ctx, in)
}

func TestPipelineRunnerService(t *testing.T) {
	originalResults := []types.TaskResult{
		{
			ID:    "1",
			Value: float64(123),
			Index: 0,
		},
		{
			ID:    "2",
			Error: errors.New("Error task"),
			Index: 1,
		},
	}

	mpr := &mockPipelineRunner{taskResults: originalResults}
	srv := &pipelineRunnerServiceServer{impl: mpr}
	client := &pipelineRunnerServiceClient{grpc: &clientAdapter{srv: srv}}

	trs, err := client.ExecuteRun(
		context.Background(),
		"my-spec",
		types.Vars{Vars: map[string]interface{}{"my-vars": true}},
		types.Options{MaxTaskDuration: time.Duration(10 * time.Second)},
	)
	require.NoError(t, err)
	assert.ElementsMatch(t, originalResults, trs)
}

func TestPipelineRunnerService_CallArgs(t *testing.T) {
	mpr := &mockPipelineRunner{}
	srv := &pipelineRunnerServiceServer{impl: mpr}
	client := &pipelineRunnerServiceClient{grpc: &clientAdapter{srv: srv}}

	spec := "my-spec"
	vars := types.Vars{
		Vars: map[string]interface{}{"my-vars": true},
	}
	options := types.Options{
		MaxTaskDuration: time.Duration(10 * time.Second),
	}
	_, err := client.ExecuteRun(context.Background(), spec, vars, options)
	require.NoError(t, err)
	assert.Equal(t, spec, mpr.spec)
	assert.Equal(t, vars, mpr.vars)
	assert.Equal(t, options, mpr.options)
}
