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
	ctx context.Context, workflow string, config []byte, tpe WorkflowSpecType) (sdk.WorkflowSpec, string, error) {
	if tpe == "" {
		tpe = DefaultSpecType
	}

	factory, ok := wsf[tpe]
	if !ok {
		return sdk.WorkflowSpec{}, "", ErrInvalidWorkflowType
	}

	rawSpec, err := factory.RawSpec(ctx, workflow)
	if err != nil {
		return sdk.WorkflowSpec{}, "", err
	}

	spec, err := factory.Spec(ctx, rawSpec, config)
	if err != nil {
		return sdk.WorkflowSpec{}, "", err
	}

	sum := sha256.New()
	sum.Write(rawSpec)
	sum.Write(config)

	return spec, fmt.Sprintf("%x", sum.Sum(nil)), nil
}

var workflowSpecFactory = WorkflowSpecFactory{
	YamlSpec: YAMLSpecFactory{},
}
