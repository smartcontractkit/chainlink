// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsClient} from "../../../dev/1_0_0/FunctionsClient.sol";
import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";
import {FunctionsRequest} from "../../../dev/1_0_0/libraries/FunctionsRequest.sol";

/**
 * @title Chainlink Functions load test client implementation
 */
contract FunctionsLoadTestClient is FunctionsClient, ConfirmedOwner {
  using FunctionsRequest for FunctionsRequest.Request;

  uint32 public constant MAX_CALLBACK_GAS = 70_000;

  bytes32 public lastRequestID;
  bytes32 public lastResponse;
  bytes32 public lastError;
  uint32 public totalRequests;
  uint32 public totalEmptyResponses;
  uint32 public totalSucceededResponses;
  uint32 public totalFailedResponses;

  constructor(address router) FunctionsClient(router) ConfirmedOwner(msg.sender) {}

  /**
   * @notice Send a simple request
   * @param source JavaScript source code
   * @param encryptedSecretsReferences Encrypted secrets payload
   * @param args List of arguments accessible from within the source code
   * @param subscriptionId Billing ID
   */
  function sendRequest(
    string calldata source,
    bytes calldata encryptedSecretsReferences,
    string[] calldata args,
    uint64 subscriptionId,
    bytes32 jobId
  ) external onlyOwner {
    FunctionsRequest.Request memory req;
    req.initializeRequestForInlineJavaScript(source);
    if (encryptedSecretsReferences.length > 0) req.addSecretsReference(encryptedSecretsReferences);
    if (args.length > 0) req.setArgs(args);
    lastRequestID = _sendRequest(req.encodeCBOR(), subscriptionId, MAX_CALLBACK_GAS, jobId);
    totalRequests += 1;
  }

  function resetCounters() external onlyOwner {
    totalRequests = 0;
    totalSucceededResponses = 0;
    totalFailedResponses = 0;
    totalEmptyResponses = 0;
  }

  function getStats() public view onlyOwner returns (uint32, uint32, uint32, uint32) {
    return (totalRequests, totalSucceededResponses, totalFailedResponses, totalEmptyResponses);
  }

  /**
   * @notice Store latest result/error
   * @param requestId The request ID, returned by sendRequest()
   * @param response Aggregated response from the user code
   * @param err Aggregated error from the user code or from the execution pipeline
   * Either response or error parameter will be set, but never both
   */
  function fulfillRequest(bytes32 requestId, bytes memory response, bytes memory err) internal override {
    // Save only the first 32 bytes of response/error to always fit within MAX_CALLBACK_GAS
    lastRequestID = requestId;
    lastResponse = bytesToBytes32(response);
    lastError = bytesToBytes32(err);
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

  function bytesToBytes32(bytes memory b) private pure returns (bytes32 out) {
    uint256 maxLen = 32;
    if (b.length < 32) {
      maxLen = b.length;
    }
    for (uint256 i = 0; i < maxLen; ++i) {
      out |= bytes32(b[i]) >> (i * 8);
    }
    return out;
  }
}
