package chainreader

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

func NewChainReader(factory Factory, family, chainId string) capabilities.ActionCapability {
	info, _ := capabilities.NewCapabilityInfo(
		fmt.Sprintf("chainreader-%s-%s", family, chainId),
		capabilities.CapabilityTypeAction,
		fmt.Sprintf("Reads from %s %s", family, chainId),
		"v0.0.1",
	)

	return &chainReader{
		Factory:        factory,
		CapabilityInfo: info,
	}
}

type Factory interface {
	ConfigObject() any
	NewChainReader(config any) (types.ChainReader, error)
}

type chainReader struct {
	Factory
	capabilities.CapabilityInfo
	workflowToReader map[string]types.ChainReader
}

var _ capabilities.ActionCapability = &chainReader{}

func (c *chainReader) Execute(ctx context.Context, workflowId string, callback chan capabilities.CapabilityResponse, inputs *values.Map) error {
	//TODO implement me
	panic("implement me")
}

func (c *chainReader) RegisterWorkflow(ctx context.Context, workflowId string, inputs *values.Map) error {
	//TODO implement me
	panic("implement me")
}

func (c *chainReader) UnregisterWorkflow(ctx context.Context, workflowId string, inputs *values.Map) error {
	//TODO implement me
	panic("implement me")
}
