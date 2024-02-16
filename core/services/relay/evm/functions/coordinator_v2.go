package functions

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type CoordinatorV2 struct {
	address common.Address

	client    client.Client
	logPoller logpoller.LogPoller
	lggr      logger.Logger
}

func NewCoordinatorV2(address common.Address, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) *CoordinatorV2 {
	return &CoordinatorV2{
		address:   address,
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

	uint32Type, errType1 := abi.NewType("uint32", "uint32", nil)
	uint40Type, errType2 := abi.NewType("uint40", "uint40", nil)
	uint64Type, errType3 := abi.NewType("uint64", "uint64", nil)
	uint72Type, errType4 := abi.NewType("uint72", "uint72", nil)
	uint96Type, errType5 := abi.NewType("uint96", "uint96", nil)
	addressType, errType6 := abi.NewType("address", "address", nil)
	bytes32Type, errType7 := abi.NewType("bytes32", "bytes32", nil)

	if errType1 != nil || errType2 != nil || errType3 != nil || errType4 != nil || errType5 != nil || errType6 != nil || errType7 != nil {
		c.lggr.Errorw("LogsToRequests: failed to initialize types", "errType1", errType1,
			"errType2", errType2, "errType3", errType3, "errType4", errType4, "errType5", errType5, "errType6", errType6, "errType7", errType7,
		)
	}

	parsingContract, err := functions_coordinator.NewFunctionsCoordinator(c.address, c.client)
	if err != nil {
		return nil, errors.Errorf("LogsToRequests: creating a contract instance for NewFunctionsCoordinator parsing failed")
	}

	for _, log := range requestLogs {
		gethLog := log.ToGethLog()
		oracleRequest, err := parsingContract.ParseOracleRequest(gethLog)
		if err != nil {
			c.lggr.Errorw("LogsToRequests: failed to parse a request log, skipping", "err", err)
			continue
		}

		commitmentABIV2 := abi.Arguments{
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
			{Type: uint72Type},  // OperationFee
		}

		commitmentBytesV2, err := commitmentABIV2.Pack(
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
			oracleRequest.Commitment.OperationFee,
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
