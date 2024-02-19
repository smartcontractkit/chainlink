package functions

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type CoordinatorV1 struct {
	address common.Address

	client    client.Client
	logPoller logpoller.LogPoller
	lggr      logger.Logger
}

func NewCoordinatorV1(address common.Address, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) *CoordinatorV1 {
	return &CoordinatorV1{
		address:   address,
		client:    client,
		logPoller: logPoller,
		lggr:      lggr,
	}
}

func (c *CoordinatorV1) Address() common.Address {
	return c.address
}

func (c *CoordinatorV1) RegisterFilters() error {
	if (c.address == common.Address{}) {
		return nil
	}

	return c.logPoller.RegisterFilter(
		logpoller.Filter{
			Name: logpoller.FilterName("FunctionsLogPollerWrapper", c.address.String(), "-v", "1"),
			EventSigs: []common.Hash{
				functions_coordinator_1_1_0.FunctionsCoordinator110OracleRequest{}.Topic(),
				functions_coordinator_1_1_0.FunctionsCoordinator110OracleResponse{}.Topic(),
			},
			Addresses: []common.Address{c.address},
		})
}

func (c *CoordinatorV1) OracleRequestLogTopic() (common.Hash, error) {
	return functions_coordinator_1_1_0.FunctionsCoordinator110OracleRequest{}.Topic(), nil
}

func (c *CoordinatorV1) OracleResponseLogTopic() (common.Hash, error) {
	return functions_coordinator_1_1_0.FunctionsCoordinator110OracleResponse{}.Topic(), nil
}

func (c *CoordinatorV1) LogsToRequests(requestLogs []logpoller.Log) ([]evmRelayTypes.OracleRequest, error) {
	var requests []evmRelayTypes.OracleRequest

	parsingContract, err := functions_coordinator_1_1_0.NewFunctionsCoordinator110(c.address, c.client)
	if err != nil {
		return nil, errors.Errorf("LogsToRequests: creating a contract instance for NewFunctionsCoordinator110 parsing failed")
	}

	for _, log := range requestLogs {
		gethLog := log.ToGethLog()
		oracleRequest, err := parsingContract.ParseOracleRequest(gethLog)
		if err != nil {
			c.lggr.Errorw("LogsToRequests: failed to parse a request log, skipping", "err", err)
			continue
		}

		commitmentABIV1 := abi.Arguments{
			{Type: bytes32Type}, // RequestId
			{Type: addressType}, // Coordinator
			{Type: uint96Type},  // EstimatedTotalCostJuels
			{Type: addressType}, // Client
			{Type: uint64Type},  // SubscriptionId
			{Type: uint32Type},  // CallbackGasLimit
			{Type: uint72Type},  // AdminFee
			{Type: uint72Type},  // DonFee
			{Type: uint40Type},  // GasOverheadBeforeCallback
			{Type: uint40Type},  // GasOverheadAfterCallback
			{Type: uint32Type},  // TimeoutTimestamp
		}

		commitmentBytesV1, err := commitmentABIV1.Pack(
			oracleRequest.Commitment.RequestId,
			oracleRequest.Commitment.Coordinator,
			oracleRequest.Commitment.EstimatedTotalCostJuels,
			oracleRequest.Commitment.Client,
			oracleRequest.Commitment.SubscriptionId,
			oracleRequest.Commitment.CallbackGasLimit,
			oracleRequest.Commitment.AdminFee,
			oracleRequest.Commitment.DonFee,
			oracleRequest.Commitment.GasOverheadBeforeCallback,
			oracleRequest.Commitment.GasOverheadAfterCallback,
			oracleRequest.Commitment.TimeoutTimestamp,
		)
		if err != nil {
			c.lggr.Errorw("LogsToRequests: failed to pack Coordinator v1 commitment bytes, skipping", err)
		}

		OracleRequestV1 := evmRelayTypes.OracleRequest{
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
			OnchainMetadata:     commitmentBytesV1,
			CoordinatorContract: c.address,
		}

		requests = append(requests, OracleRequestV1)
	}
	return requests, nil
}

func (c *CoordinatorV1) LogsToResponses(responseLogs []logpoller.Log) ([]evmRelayTypes.OracleResponse, error) {
	var responses []evmRelayTypes.OracleResponse

	parsingContract, err := functions_coordinator_1_1_0.NewFunctionsCoordinator110(c.address, c.client)
	if err != nil {
		return nil, errors.Errorf("LogsToResponses: creating a contract instance for parsing failed")
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
