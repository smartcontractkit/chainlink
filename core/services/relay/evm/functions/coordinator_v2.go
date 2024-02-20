package functions

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type CoordinatorV2 struct {
	address common.Address
	abiArgs abi.Arguments

	client    client.Client
	logPoller logpoller.LogPoller
	lggr      logger.Logger
}

func NewCoordinatorV2(address common.Address, abiArgs abi.Arguments, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) *CoordinatorV2 {
	return &CoordinatorV2{
		address:   address,
		abiArgs:   abiArgs,
		client:    client,
		logPoller: logPoller,
		lggr:      lggr,
	}
}

func (c *CoordinatorV2) Address() common.Address {
	return c.address
}

func (c *CoordinatorV2) RegisterFilters() error {
	if (c.address == common.Address{}) {
		return nil
	}

	return c.logPoller.RegisterFilter(
		logpoller.Filter{
			Name: logpoller.FilterName("FunctionsLogPollerWrapper", c.address.String(), "-v", "2"),
			EventSigs: []common.Hash{
				functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(),
				functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(),
			},
			Addresses: []common.Address{c.address},
		})
}
func (c *CoordinatorV2) OracleRequestLogTopic() common.Hash {
	return functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic()
}
func (c *CoordinatorV2) OracleResponseLogTopic() common.Hash {
	return functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic()
}
func (c *CoordinatorV2) LogsToRequests(requestLogs []logpoller.Log) ([]evmRelayTypes.OracleRequest, error) {
	var requests []evmRelayTypes.OracleRequest

	parsingContract, err := functions_coordinator.NewFunctionsCoordinator(c.address, c.client)
	if err != nil {
		return nil, fmt.Errorf("LogsToRequests: creating a contract instance for NewFunctionsCoordinator parsing failed: %w", err)
	}

	for _, log := range requestLogs {
		gethLog := log.ToGethLog()
		oracleRequest, err := parsingContract.ParseOracleRequest(gethLog)
		if err != nil {
			c.lggr.Errorw("LogsToRequests: failed to parse a request log, skipping", "err", err)
			continue
		}

		commitmentBytesV2, err := c.abiArgs.Pack(
			oracleRequest.Commitment.RequestId,
			oracleRequest.Commitment.Coordinator,
			oracleRequest.Commitment.EstimatedTotalCostJuels,
			oracleRequest.Commitment.Client,
			oracleRequest.Commitment.SubscriptionId,
			oracleRequest.Commitment.CallbackGasLimit,
			oracleRequest.Commitment.AdminFeeJuels,
			oracleRequest.Commitment.DonFeeJuels,
			oracleRequest.Commitment.GasOverheadBeforeCallback,
			oracleRequest.Commitment.GasOverheadAfterCallback,
			oracleRequest.Commitment.TimeoutTimestamp,
			oracleRequest.Commitment.OperationFeeJuels,
		)
		if err != nil {
			c.lggr.Errorw("LogsToRequests: failed to pack Coordinator v2 commitment bytes, skipping", err)
		}

		OracleRequestV2 := evmRelayTypes.OracleRequest{
			RequestId:           oracleRequest.RequestId,
			RequestingContract:  oracleRequest.RequestingContract,
			RequestInitiator:    oracleRequest.RequestInitiator,
			SubscriptionId:      oracleRequest.SubscriptionId,
			SubscriptionOwner:   oracleRequest.SubscriptionOwner,
			Data:                oracleRequest.Data,
			DataVersion:         oracleRequest.DataVersion,
			Flags:               oracleRequest.Flags,
			CallbackGasLimit:    oracleRequest.CallbackGasLimit,
			TxHash:              oracleRequest.Raw.TxHash,
			OnchainMetadata:     commitmentBytesV2,
			CoordinatorContract: c.address,
		}

		requests = append(requests, OracleRequestV2)
	}
	return requests, nil
}

func (c *CoordinatorV2) LogsToResponses(responseLogs []logpoller.Log) ([]evmRelayTypes.OracleResponse, error) {
	var responses []evmRelayTypes.OracleResponse

	parsingContract, err := functions_coordinator.NewFunctionsCoordinator(c.address, c.client)
	if err != nil {
		return nil, fmt.Errorf("LogsToResponses: creating a contract instance for parsing failed: %w", err)
	}
	for _, log := range responseLogs {
		gethLog := log.ToGethLog()
		oracleResponse, err := parsingContract.ParseOracleResponse(gethLog)
		if err != nil {
			c.lggr.Errorw("LogsToResponses: failed to parse a response log, skipping")
			continue
		}
		responses = append(responses, evmRelayTypes.OracleResponse{
			RequestId: oracleResponse.RequestId,
		})
	}
	return responses, nil
}
