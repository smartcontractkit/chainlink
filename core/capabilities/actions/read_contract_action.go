package actions

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ReadContractConfig struct {
	ChainId uint64 `json:"chainId"`
	Network string `json:"network"`
}

type RequestConfig struct {
	ContractReaderConfig evmtypes.ChainReaderConfig `json:"contractReaderConfig"`
}

type Input struct {
	ContractName    string         `json:"contractName"`
	ContractAddress common.Address `json:"contractAddress"`
	ConfidenceLevel string         `json:"confidenceLevel"`
}

type ReadContractAction struct {
	capabilities.CapabilityInfo
	capabilities.Validator[RequestConfig, Input, capabilities.TriggerResponse]

	lggr logger.Logger

	relayer core.Relayer
}

func NewReadContractAction(lggr logger.Logger, config ReadContractConfig, relayer core.Relayer) *ReadContractAction {
	id := fmt.Sprintf("read-contract-%s-%d@1.0.0", config.Network, config.ChainId)

	info := capabilities.MustNewCapabilityInfo(
		id,
		capabilities.CapabilityTypeAction,
		"Read Contract Action.  Supports reading from a contract.",
	)

	return &ReadContractAction{
		CapabilityInfo: info,
		lggr:           lggr,
		relayer:        relayer,
	}
}

func (r ReadContractAction) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {

	// going to need a contract reader instance per config

	request.
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
