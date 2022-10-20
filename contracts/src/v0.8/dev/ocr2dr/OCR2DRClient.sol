// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./OCR2DR.sol";
import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";

/**
 * @title The OCR2DR client contract
 * @notice Contract writers can inherit this contract in order to create on-demand OCR requests
 */
abstract contract OCR2DRClient is OCR2DRClientInterface {
  OCR2DROracleInterface private s_oracle;
  mapping(bytes32 => address) private s_pendingRequests;

  event RequestSent(bytes32 indexed id);
  event RequestFulfilled(bytes32 indexed id);

  error SenderIsNotOracle();
  error RequestIsAlreadyPending();
  error RequestIsNotPending();

  constructor(address oracle) {
    setOracle(oracle);
  }

  /// @inheritdoc OCR2DRClientInterface
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_oracle.getDONPublicKey();
  }

  /**
   * @notice Sends OCR2DR request to the stored oracle address
   * @param req The initialized OCR2DR.Request
   * @param subscriptionId The subscription ID
   * @return requestId The generated request ID
   */
  function sendRequest(OCR2DR.Request memory req, uint256 subscriptionId) internal returns (bytes32) {
    bytes32 requestId = s_oracle.sendRequest(subscriptionId, OCR2DR.encodeCBOR(req));
    s_pendingRequests[requestId] = address(s_oracle);
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

  /// @inheritdoc OCR2DRClientInterface
  function handleOracleFulfillment(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) external override recordChainlinkFulfillment(requestId) {
    fulfillRequest(requestId, response, err);
  }

  /**
   * @notice Sets the stored oracle address
   * @param oracleAddress The address of OCR2DR oracle contract
   */
  function setOracle(address oracleAddress) internal {
    s_oracle = OCR2DROracleInterface(oracleAddress);
  }

  /**
   * @notice Gets the stored address of the oracle contract
   * @return The address of the oracle contract
   */
  function getChainlinkOracleAddress() internal view returns (address) {
    return address(s_oracle);
  }

  /**
   * @notice Allows for a request which was created on another contract to be fulfilled
   * on this contract
   * @param oracleAddress The address of the oracle contract that will fulfill the request
   * @param requestId The request ID used for the response
   */
  function addExternalRequest(address oracleAddress, bytes32 requestId) internal notPendingRequest(requestId) {
    s_pendingRequests[requestId] = oracleAddress;
  }

  /**
   * @dev Reverts if the sender is not the oracle of the request.
   * Emits RequestFulfilled event.
   * @param requestId The request ID for fulfillment
   */
  modifier recordChainlinkFulfillment(bytes32 requestId) {
    if (msg.sender != s_pendingRequests[requestId]) {
      revert SenderIsNotOracle();
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
