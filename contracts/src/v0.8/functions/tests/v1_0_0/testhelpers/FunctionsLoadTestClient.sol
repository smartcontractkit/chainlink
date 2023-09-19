// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsClient} from "../../../dev/v1_0_0/FunctionsClient.sol";
import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";
import {FunctionsRequest} from "../../../dev/v1_0_0/libraries/FunctionsRequest.sol";

/**
 * @title Chainlink Functions load test client implementation
 */
contract FunctionsLoadTestClient is FunctionsClient, ConfirmedOwner {
  using FunctionsRequest for FunctionsRequest.Request;

  uint32 public constant MAX_CALLBACK_GAS = 250_000;

  bytes32 public lastRequestID;
  bytes public lastResponse;
  bytes public lastError;
  uint32 public totalRequests;
  uint32 public totalEmptyResponses;
  uint32 public totalSucceededResponses;
  uint32 public totalFailedResponses;

  constructor(address router) FunctionsClient(router) ConfirmedOwner(msg.sender) {}

  /**
   * @notice Send a simple request
   * @param times Number of times to send the request
   * @param source JavaScript source code
   * @param encryptedSecretsReferences Encrypted secrets payload
   * @param args List of arguments accessible from within the source code
   * @param subscriptionId Billing ID
   * @param donId DON ID
   */
  function sendRequest(
    uint32 times,
    string calldata source,
    bytes calldata encryptedSecretsReferences,
    string[] calldata args,
    uint64 subscriptionId,
    bytes32 donId
  ) external onlyOwner {
    FunctionsRequest.Request memory req;
    req.initializeRequestForInlineJavaScript(source);
    if (encryptedSecretsReferences.length > 0) req.addSecretsReference(encryptedSecretsReferences);
    if (args.length > 0) req.setArgs(args);
    uint i = 0;
    for (i = 0; i < times; i++) {
      lastRequestID = _sendRequest(req.encodeCBOR(), subscriptionId, MAX_CALLBACK_GAS, donId);
      totalRequests += 1;
    }
  }

  /**
   * @notice Same as sendRequest but for DONHosted secrets
   * @param times Number of times to send the request
   * @param source JavaScript source code
   * @param slotId DON hosted secrets slot ID
   * @param slotVersion DON hosted secrets slot version
   * @param args List of arguments accessible from within the source code
   * @param subscriptionId Billing ID
   * @param donId DON ID
   */
  function sendRequestWithDONHostedSecrets(
    uint32 times,
    string calldata source,
    uint8 slotId,
    uint64 slotVersion,
    string[] calldata args,
    uint64 subscriptionId,
    bytes32 donId
  ) public onlyOwner {
    FunctionsRequest.Request memory req;
    req.initializeRequestForInlineJavaScript(source);
    req.addDONHostedSecrets(slotId, slotVersion);
    if (args.length > 0) req.setArgs(args);
    uint i = 0;
    for (i = 0; i < times; i++) {
      lastRequestID = _sendRequest(req.encodeCBOR(), subscriptionId, MAX_CALLBACK_GAS, donId);
      totalRequests += 1;
    }
  }

  /**
   * @notice Sends a Chainlink Functions request that has already been CBOR encoded
   * @param times Number of times to send the request
   * @param cborEncodedRequest The CBOR encoded bytes data for a Functions request
   * @param subscriptionId The subscription ID that will be charged to service the request
   * @param donId DON ID
   */
  function sendEncodedRequest(
    uint32 times,
    bytes memory cborEncodedRequest,
    uint64 subscriptionId,
    bytes32 donId
  ) public onlyOwner {
    uint i = 0;
    for (i = 0; i < times; i++) {
      lastRequestID = _sendRequest(cborEncodedRequest, subscriptionId, MAX_CALLBACK_GAS, donId);
      totalRequests += 1;
    }
  }

  function resetStats() external onlyOwner {
    lastRequestID = "";
    lastResponse = "";
    lastError = "";
    totalRequests = 0;
    totalSucceededResponses = 0;
    totalFailedResponses = 0;
    totalEmptyResponses = 0;
  }

  function getStats()
    public
    view
    onlyOwner
    returns (bytes32, bytes memory, bytes memory, uint32, uint32, uint32, uint32)
  {
    return (
      lastRequestID,
      lastResponse,
      lastError,
      totalRequests,
      totalSucceededResponses,
      totalFailedResponses,
      totalEmptyResponses
    );
  }

  /**
   * @notice Store latest result/error
   * @param requestId The request ID, returned by sendRequest()
   * @param response Aggregated response from the user code
   * @param err Aggregated error from the user code or from the execution pipeline
   * Either response or error parameter will be set, but never both
   */
  function fulfillRequest(bytes32 requestId, bytes memory response, bytes memory err) internal override {
    lastRequestID = requestId;
    lastResponse = response;
    lastError = err;
    if (response.length == 0) {
      totalEmptyResponses += 1;
    }
    if (err.length != 0) {
      totalFailedResponses += 1;
    }
    if (response.length != 0 && err.length == 0) {
      totalSucceededResponses += 1;
    }
  }
}
