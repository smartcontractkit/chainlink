// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Functions} from "./Functions.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsClient} from "./interfaces/IFunctionsClient.sol";

/**
 * @title The Chainlink Functions client contract
 * @notice Contract writers can inherit this contract in order to create Chainlink Functions requests
 */
abstract contract FunctionsClient is IFunctionsClient {
  IFunctionsRouter internal immutable s_router;

  event RequestSent(bytes32 indexed id);
  event RequestFulfilled(bytes32 indexed id);

  error OnlyRouterCanFufill();

  constructor(address router) {
    s_router = IFunctionsRouter(router);
  }

  /**
   * @notice Sends a Chainlink Functions request to the stored oracle address
   * @param req The initialized Functions.Request
   * @param subscriptionId The subscription ID
   * @param callbackGasLimit gas limit for the fulfillment callback
   * @return requestId The generated request ID
   */
  function _sendRequest(
    Functions.Request memory req,
    uint64 subscriptionId,
    uint32 callbackGasLimit,
    bytes32 donId
  ) internal returns (bytes32 requestId) {
    bytes memory requestData = Functions.encodeCBOR(req);
    requestId = _sendRequestBytes(requestData, subscriptionId, callbackGasLimit, donId);
  }

  /**
   * @notice Sends a Chainlink Functions request to the stored oracle address
   * @param data The initialized Functions request data
   * @param subscriptionId The subscription ID
   * @param callbackGasLimit gas limit for the fulfillment callback
   * @return requestId The generated request ID
   */
  function _sendRequestBytes(
    bytes memory data,
    uint64 subscriptionId,
    uint32 callbackGasLimit,
    bytes32 donId
  ) internal returns (bytes32 requestId) {
    requestId = s_router.sendRequest(subscriptionId, data, Functions.REQUEST_DATA_VERSION, callbackGasLimit, donId);
    emit RequestSent(requestId);
  }

  /**
   * @notice User defined function to handle a response
   * @param requestId The request ID, returned by sendRequest()
   * @param response Aggregated response from the user code
   * @param err Aggregated error from the user code or from the execution pipeline
   * Either response or error parameter will be set, but never both
   */
  function fulfillRequest(bytes32 requestId, bytes memory response, bytes memory err) internal virtual;

  /**
   * @inheritdoc IFunctionsClient
   */
  function handleOracleFulfillment(bytes32 requestId, bytes memory response, bytes memory err) external override {
    if (msg.sender != address(s_router)) {
      revert OnlyRouterCanFufill();
    }
    fulfillRequest(requestId, response, err);
  }
}
