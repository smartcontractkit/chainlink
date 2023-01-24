// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/functions/FunctionsClient.sol";

contract FunctionsClientTestHelper is FunctionsClient {
  using Functions for Functions.Request;

  event SendRequestInvoked(bytes32 requestId, string sourceCode, uint64 subscriptionId);
  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  bool private s_revertFulfillRequest;
  bool private s_doInvalidOperation;

  constructor(address oracle) FunctionsClient(oracle) {}

  function sendSimpleRequestWithJavaScript(string memory sourceCode, uint64 subscriptionId)
    public
    returns (bytes32 requestId)
  {
    Functions.Request memory request;
    request.initializeRequestForInlineJavaScript(sourceCode);
    requestId = sendRequest(request, subscriptionId, 20_000);
    emit SendRequestInvoked(requestId, sourceCode, subscriptionId);
  }

  function estimateJuelCost(
    string memory sourceCode,
    uint64 subscriptionId,
    uint256 gasCost
  ) public view returns (uint96) {
    Functions.Request memory request;
    request.initializeRequestForInlineJavaScript(sourceCode);
    return estimateCost(request, subscriptionId, 20_000, gasCost);
  }

  function fulfillRequest(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) internal override {
    if (s_revertFulfillRequest) {
      revert("asked to revert");
    }
    if (s_doInvalidOperation) {
      uint256 x = 1;
      uint256 y = 0;
      x = x / y;
    }
    emit FulfillRequestInvoked(requestId, response, err);
  }

  function setRevertFulfillRequest(bool on) external {
    s_revertFulfillRequest = on;
  }

  function setDoInvalidOperation(bool on) external {
    s_doInvalidOperation = on;
  }
}
