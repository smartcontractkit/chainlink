// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./OCR2DRClient.sol";
import "./OCR2DR.sol";
import "../../ConfirmedOwner.sol";

/**
 * @title OCR2DR example client contract implementation
 */
contract OCR2DRClientExample is OCR2DRClient, ConfirmedOwner {
  using OCR2DR for OCR2DR.Request;

  bytes32 public lastRequestId;
  bytes public lastResponse;
  bytes public lastError;

  error UnexpectedRequestID(bytes32 requestId);

  constructor(address oracle) OCR2DRClient(oracle) ConfirmedOwner(msg.sender) {}

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
    OCR2DR.Request memory req;
    req.initializeRequestForInlineJavaScript(source);
    if (secrets.length > 0) req.addInlineSecrets(secrets);
    if (args.length > 0) req.addArgs(args);
    lastRequestId = sendRequest(req, subscriptionId, 40_000, tx.gasprice);
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
    lastResponse = response;
    lastError = err;
  }
}
