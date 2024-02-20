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
	address  common.Address
	abiTypes *abiTypes

	client    client.Client
	logPoller logpoller.LogPoller
	lggr      logger.Logger
}

func NewCoordinatorV2(address common.Address, abiTypes *abiTypes, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) *CoordinatorV2 {
	return &CoordinatorV2{
		address:   address,
		abiTypes:  abiTypes,
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
func (c *CoordinatorV2) OracleRequestLogTopic() (common.Hash, error) {
	return functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(), nil
}
func (c *CoordinatorV2) OracleResponseLogTopic() (common.Hash, error) {
	return functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(), nil
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

		commitmentABIV2 := abi.Arguments{
			{Type: c.abiTypes.bytes32Type}, // RequestId
			{Type: c.abiTypes.addressType}, // Coordinator
			{Type: c.abiTypes.uint96Type},  // EstimatedTotalCostJuels
			{Type: c.abiTypes.addressType}, // Client
			{Type: c.abiTypes.uint64Type},  // SubscriptionId
			{Type: c.abiTypes.uint32Type},  // CallbackGasLimit
			{Type: c.abiTypes.uint72Type},  // AdminFee
			{Type: c.abiTypes.uint72Type},  // DonFee
			{Type: c.abiTypes.uint40Type},  // GasOverheadBeforeCallback
			{Type: c.abiTypes.uint40Type},  // GasOverheadAfterCallback
			{Type: c.abiTypes.uint32Type},  // TimeoutTimestamp
			{Type: c.abiTypes.uint72Type},  // OperationFee
		}

		commitmentBytesV2, err := commitmentABIV2.Pack(
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
