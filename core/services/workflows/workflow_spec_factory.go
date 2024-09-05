package workflows

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var ErrInvalidWorkflowType = errors.New("invalid workflow type")

type SDKWorkflowSpecFactory interface {
	GetSpec(rawSpec, config []byte) (sdk.WorkflowSpec, error)
	GetRawSpec(wf string) ([]byte, error)
}

type WorkflowSpecFactory map[job.WorkflowSpecType]SDKWorkflowSpecFactory

func (wsf WorkflowSpecFactory) ToSpec(workflow string, config []byte, tpe job.WorkflowSpecType) (sdk.WorkflowSpec, string, error) {
	if tpe == "" {
		tpe = job.DefaultSpecType
	}

	factory, ok := wsf[tpe]
	if !ok {
		return sdk.WorkflowSpec{}, "", ErrInvalidWorkflowType
	}

	rawSpec, err := factory.GetRawSpec(workflow)
	if err != nil {
		return sdk.WorkflowSpec{}, "", err
	}

	spec, err := factory.GetSpec(rawSpec, config)
	if err != nil {
		return sdk.WorkflowSpec{}, "", err
	}

	sum := sha256.New()
	sum.Write(rawSpec)
	sum.Write(config)

	return spec, fmt.Sprintf("%x", sum.Sum(nil)), nil
}

var workflowSpecFactory = WorkflowSpecFactory{
	job.YamlSpec: YAMLSpecFactory{},
}
