// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsClient} from "./FunctionsClient.sol";
import {Functions} from "./Functions.sol";
import {ConfirmedOwner} from "../../ConfirmedOwner.sol";

/**
 * @title Chainlink Functions example client contract implementation
 */
contract FunctionsClientExample is FunctionsClient, ConfirmedOwner {
  using Functions for Functions.Request;

  uint32 public constant MAX_CALLBACK_GAS = 70_000;

  bytes32 public lastRequestId;
  bytes32 public lastResponse;
  bytes32 public lastError;
  uint32 public lastResponseLength;
  uint32 public lastErrorLength;

  error UnexpectedRequestID(bytes32 requestId);

  constructor(address oracle) FunctionsClient(oracle) ConfirmedOwner(msg.sender) {}

  /**
   * @notice Send a simple request
   * @param source JavaScript source code
   * @param secrets Encrypted secrets payload
   * @param args List of arguments accessible from within the source code
   * @param subscriptionId Billing ID
   */
  function SendRequest(
    string calldata source,
    bytes calldata secrets,
    string[] calldata args,
    uint64 subscriptionId
  ) external onlyOwner {
    Functions.Request memory req;
    req.initializeRequestForInlineJavaScript(source);
    if (secrets.length > 0) req.addInlineSecrets(secrets);
    if (args.length > 0) req.addArgs(args);
    lastRequestId = sendRequest(req, subscriptionId, MAX_CALLBACK_GAS);
  }

  /**
   * @notice Store latest result/error
   * @param requestId The request ID, returned by sendRequest()
   * @param response Aggregated response from the user code
   * @param err Aggregated error from the user code or from the execution pipeline
   * Either response or error parameter will be set, but never both
   */
  function fulfillRequest(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) internal override {
    if (lastRequestId != requestId) {
      revert UnexpectedRequestID(requestId);
    }
    // Save only the first 32 bytes of reponse/error to always fit within MAX_CALLBACK_GAS
    lastResponse = bytesToBytes32(response);
    lastResponseLength = uint32(response.length);
    lastError = bytesToBytes32(err);
    lastErrorLength = uint32(err.length);
  }

  function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
    bytes32 out;
    uint256 maxLen = 32;
    if (b.length < 32) {
      maxLen = b.length;
    }
    for (uint256 i = 0; i < maxLen; i++) {
      out |= bytes32(b[i]) >> (i * 8);
    }
    return out;
  }
}
