package pipeline

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
)

var _ core.PipelineRunnerService = (*pipelineRunnerServiceClient)(nil)

type pipelineRunnerServiceClient struct {
	*net.BrokerExt
	grpc pb.PipelineRunnerServiceClient
}

func NewRunnerClient(cc grpc.ClientConnInterface) *pipelineRunnerServiceClient {
	return &pipelineRunnerServiceClient{grpc: pb.NewPipelineRunnerServiceClient(cc)}
}

func (p pipelineRunnerServiceClient) ExecuteRun(ctx context.Context, spec string, vars core.Vars, options core.Options) (core.TaskResults, error) {
	varsStruct, err := structpb.NewStruct(vars.Vars)
	if err != nil {
		return nil, err
	}

	rr := pb.RunRequest{
		Spec: spec,
		Vars: varsStruct,
		Options: &pb.Options{
			MaxTaskDuration: durationpb.New(options.MaxTaskDuration),
		},
	}

	executeRunResult, err := p.grpc.ExecuteRun(ctx, &rr)
	if err != nil {
		return nil, err
	}

	trs := make([]core.TaskResult, len(executeRunResult.Results))
	for i, trr := range executeRunResult.Results {
		var err error
		if trr.HasError {
			err = errors.New(trr.Error)
		}

		js := jsonserializable.JSONSerializable{}
		err2 := js.UnmarshalJSON(trr.Value)
		if err2 != nil {
			return nil, err2
		}
		trs[i] = core.TaskResult{
			ID:   trr.Id,
			Type: trr.Type,
			TaskValue: core.TaskValue{
				Value:      js,
				Error:      err,
				IsTerminal: trr.IsTerminal,
			},
			Index: int(trr.Index),
		}
	}

	return trs, nil
}

var _ pb.PipelineRunnerServiceServer = (*RunnerServer)(nil)

type RunnerServer struct {
	pb.UnimplementedPipelineRunnerServiceServer
	*net.BrokerExt

	impl core.PipelineRunnerService
}

func NewRunnerServer(impl core.PipelineRunnerService) *RunnerServer {
	return &RunnerServer{impl: impl}
}

func (p *RunnerServer) ExecuteRun(ctx context.Context, rr *pb.RunRequest) (*pb.RunResponse, error) {
	vars := core.Vars{
		Vars: rr.Vars.AsMap(),
	}
	options := core.Options{
		MaxTaskDuration: rr.Options.MaxTaskDuration.AsDuration(),
	}
	trs, err := p.impl.ExecuteRun(ctx, rr.Spec, vars, options)
	if err != nil {
		return nil, err
	}

	taskResults := make([]*pb.TaskResult, len(trs))
	for i, trr := range trs {
		v, err := trr.Value.MarshalJSON()
		if err != nil {
			return nil, err
		}

		hasError := trr.Error != nil
		errs := ""
		if hasError {
			errs = trr.Error.Error()
		}
		taskResults[i] = &pb.TaskResult{
			Id:         trr.ID,
			Type:       trr.Type,
			Value:      v,
			Error:      errs,
			HasError:   hasError,
			IsTerminal: trr.IsTerminal,
			Index:      int32(trr.Index),
		}
	}

	return &pb.RunResponse{
		Results: taskResults,
	}, nil
}
