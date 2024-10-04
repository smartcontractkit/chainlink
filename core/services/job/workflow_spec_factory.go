package job

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

var ErrInvalidWorkflowType = errors.New("invalid workflow type")

type SDKWorkflowSpecFactory interface {
	Spec(ctx context.Context, rawSpec, config []byte) (sdk.WorkflowSpec, error)
	RawSpec(ctx context.Context, wf string) ([]byte, error)
}

type WorkflowSpecFactory map[WorkflowSpecType]SDKWorkflowSpecFactory

func (wsf WorkflowSpecFactory) Spec(
	ctx context.Context, workflow string, config []byte, tpe WorkflowSpecType) (sdk.WorkflowSpec, []byte, string, error) {
	if tpe == "" {
		tpe = DefaultSpecType
	}

	factory, ok := wsf[tpe]
	if !ok {
		return sdk.WorkflowSpec{}, nil, "", ErrInvalidWorkflowType
	}

	rawSpec, err := factory.RawSpec(ctx, workflow)
	if err != nil {
		return sdk.WorkflowSpec{}, nil, "", err
	}

	spec, err := factory.Spec(ctx, rawSpec, config)
	if err != nil {
		return sdk.WorkflowSpec{}, nil, "", err
	}

	sum := sha256.New()
	sum.Write(rawSpec)
	sum.Write(config)

	return spec, rawSpec, fmt.Sprintf("%x", sum.Sum(nil)), nil
}

func (wsf WorkflowSpecFactory) RawSpec(
	ctx context.Context, workflow string, tpe WorkflowSpecType) ([]byte, error) {
	if tpe == "" {
		tpe = DefaultSpecType
	}

	factory, ok := wsf[tpe]
	if !ok {
		return nil, ErrInvalidWorkflowType
	}

	return factory.RawSpec(ctx, workflow)
}

var workflowSpecFactory = WorkflowSpecFactory{
	YamlSpec: YAMLSpecFactory{},
}
