// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./Functions.sol";
import "../interfaces/FunctionsClientInterface.sol";
import "../interfaces/FunctionsDONInterface.sol";

/**
 * @title The Chainlink Functions client contract
 * @notice Contract writers can inherit this contract in order to create Chainlink Functions requests
 */
abstract contract FunctionsClient is FunctionsClientInterface {
  FunctionsDONInterface private s_don;
  mapping(bytes32 => address) private s_pendingRequests;

  event RequestSent(bytes32 indexed id);
  event RequestFulfilled(bytes32 indexed id);

  error SenderIsNotRegistry();
  error RequestIsAlreadyPending();
  error RequestIsNotPending();

  constructor(address don) {
    setDON(don);
  }

  /**
   * @inheritdoc FunctionsClientInterface
   */
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_don.getDONPublicKey();
  }

  /**
   * @notice Estimate the total cost that will be charged to a subscription to make a request: gas re-imbursement, plus DON fee, plus Registry fee
   * @param req The initialized Functions.Request
   * @param subscriptionId The subscription ID
   * @param gasLimit gas limit for the fulfillment callback
   * @return billedCost Cost in Juels (1e18) of LINK
   */
  function estimateCost(
    Functions.Request memory req,
    uint64 subscriptionId,
    uint32 gasLimit,
    uint256 gasPrice
  ) public view returns (uint96) {
    return s_don.estimateCost(subscriptionId, Functions.encodeCBOR(req), gasLimit, gasPrice);
  }

  /**
   * @notice Sends a Chainlink Functions request to the stored DON address
   * @param req The initialized Functions.Request
   * @param subscriptionId The subscription ID
   * @param gasLimit gas limit for the fulfillment callback
   * @return requestId The generated request ID
   */
  function sendRequest(
    Functions.Request memory req,
    uint64 subscriptionId,
    uint32 gasLimit,
    uint256 gasPrice
  ) internal returns (bytes32) {
    bytes32 requestId = s_don.sendRequest(subscriptionId, Functions.encodeCBOR(req), gasLimit, gasPrice);
    s_pendingRequests[requestId] = s_don.getRegistry();
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
  function fulfillRequest(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) internal virtual;

  /**
   * @inheritdoc FunctionsClientInterface
   */
  function handleDONFulfillment(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) external override recordChainlinkFulfillment(requestId) {
    fulfillRequest(requestId, response, err);
  }

  /**
   * @notice Sets the stored DON address
   * @param don The address of Functions DON contract
   */
  function setDON(address don) internal {
    s_don = FunctionsDONInterface(don);
  }

  /**
   * @notice Gets the stored address of the DON contract
   * @return The address of the DON contract
   */
  function getChainlinkDONAddress() internal view returns (address) {
    return address(s_don);
  }

  /**
   * @notice Allows for a request which was created on another contract to be fulfilled
   * on this contract
   * @param donAddress The address of the DON contract that will fulfill the request
   * @param requestId The request ID used for the response
   */
  function addExternalRequest(address donAddress, bytes32 requestId) internal notPendingRequest(requestId) {
    s_pendingRequests[requestId] = donAddress;
  }

  /**
   * @dev Reverts if the sender is not the DON that serviced the request.
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
