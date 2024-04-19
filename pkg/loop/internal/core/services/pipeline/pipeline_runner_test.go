package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
)

type mockPipelineRunner struct {
	taskResults []core.TaskResult
	err         error
	spec        string
	vars        core.Vars
	options     core.Options
}

func (m *mockPipelineRunner) ExecuteRun(ctx context.Context, spec string, vars core.Vars, options core.Options) (core.TaskResults, error) {
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
	originalResults := []core.TaskResult{
		{
			ID: "1",
			TaskValue: core.TaskValue{
				Value: jsonserializable.JSONSerializable{
					Val:   123.123,
					Valid: true,
				},
			},
			Index: 0,
		},
		{
			ID: "2",

			TaskValue: core.TaskValue{
				Value: jsonserializable.JSONSerializable{},
				Error: errors.New("Error task"),
			},
			Index: 1,
		},
	}

	mpr := &mockPipelineRunner{taskResults: originalResults}
	srv := &RunnerServer{impl: mpr}
	client := &pipelineRunnerServiceClient{grpc: &clientAdapter{srv: srv}}

	trs, err := client.ExecuteRun(
		context.Background(),
		"my-spec",
		core.Vars{Vars: map[string]interface{}{"my-vars": true}},
		core.Options{MaxTaskDuration: 10 * time.Second},
	)
	require.NoError(t, err)
	assert.ElementsMatch(t, originalResults, trs)
}

func TestPipelineRunnerService_CallArgs(t *testing.T) {
	mpr := &mockPipelineRunner{}
	srv := &RunnerServer{impl: mpr}
	client := &pipelineRunnerServiceClient{grpc: &clientAdapter{srv: srv}}

	spec := "my-spec"
	vars := core.Vars{
		Vars: map[string]interface{}{"my-vars": true},
	}
	options := core.Options{
		MaxTaskDuration: 10 * time.Second,
	}
	_, err := client.ExecuteRun(context.Background(), spec, vars, options)
	require.NoError(t, err)
	assert.Equal(t, spec, mpr.spec)
	assert.Equal(t, vars, mpr.vars)
	assert.Equal(t, options, mpr.options)
}
