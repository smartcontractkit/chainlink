package actions

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ReadContractAction struct {
	capabilities.CapabilityInfo

	lggr           logger.Logger
	contractReader commontypes.ContractReader
}

func NewReadContractAction(lggr logger.Logger, id string, cr commontypes.ContractReader) *ReadContractAction {
	info := capabilities.MustNewCapabilityInfo(
		id,
		capabilities.CapabilityTypeAction,
		"Read Contract Action.  Supports reading from a contract.",
	)

	return &ReadContractAction{
		CapabilityInfo: info,
		lggr:           lggr,
		contractReader: cr,
	}
}

func (r ReadContractAction) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	r.contractReader.Bind(request.Config)

	//TODO implement me
	panic("implement me")
}

func (r ReadContractAction) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	// Do Nothing
	return nil
}

func (r ReadContractAction) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	// Do Nothing
	return nil
}
