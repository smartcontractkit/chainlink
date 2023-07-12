// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {Functions} from "./Functions.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsClient} from "./interfaces/IFunctionsClient.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";

/**
 * @title The Chainlink Functions client contract
 * @notice Contract writers can inherit this contract in order to create Chainlink Functions requests
 */
abstract contract FunctionsClient is IFunctionsClient {
  IFunctionsRouter private s_router;
  mapping(bytes32 => address) internal s_pendingRequests; /* requestId => fulfillment sender */

  event RequestSent(bytes32 indexed id);
  event RequestFulfilled(bytes32 indexed id);

  error SenderIsNotRegistry();
  error RequestIsAlreadyPending();
  error RequestIsNotPending();

  constructor(address router) {
    setRouter(router);
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
  ) internal returns (bytes32) {
    bytes memory requestData = Functions.encodeCBOR(req);
    bytes32 requestId = _sendRequestBytes(requestData, subscriptionId, callbackGasLimit, donId);
    return requestId;
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
  ) internal returns (bytes32) {
    bytes32 requestId = s_router.sendRequest(
      subscriptionId,
      data,
      Functions.REQUEST_DATA_VERSION,
      callbackGasLimit,
      donId
    );
    s_pendingRequests[requestId] = s_router.getRoute(donId);
    emit RequestSent(requestId);
    return requestId;
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
  function handleOracleFulfillment(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) external override recordChainlinkFulfillment(requestId) {
    fulfillRequest(requestId, response, err);
  }

  /**
   * @notice Sets the stored router address
   * @param router The address of Functions router contract
   */
  function setRouter(address router) internal {
    s_router = IFunctionsRouter(router);
  }

  /**
   * @notice Gets the stored address of the router contract
   * @return The address of the router contract
   */
  function getRouter() internal view returns (address) {
    return address(s_router);
  }

  /**
   * @dev Reverts if the sender is not the oracle that serviced the request.
   * Emits RequestFulfilled event.
   * @param requestId The request ID for fulfillment
   */
  modifier recordChainlinkFulfillment(bytes32 requestId) {
    if (msg.sender != s_pendingRequests[requestId]) {
      revert SenderIsNotRegistry();
    }
    delete s_pendingRequests[requestId];
    emit RequestFulfilled(requestId);
    _;
  }

  /**
   * @dev Reverts if the request is already pending
   * @param requestId The request ID for fulfillment
   */
  modifier notPendingRequest(bytes32 requestId) {
    if (s_pendingRequests[requestId] != address(0)) {
      revert RequestIsAlreadyPending();
    }
    _;
  }
}
