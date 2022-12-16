// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/ocr2dr/OCR2DRClient.sol";

contract OCR2DRClientTestHelper is OCR2DRClient {
  using OCR2DR for OCR2DR.Request;

  event SendRequestInvoked(bytes32 requestId, string sourceCode, uint64 subscriptionId);
  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  bool private s_revertFulfillRequest;
  bool private s_doInvalidOperation;

  constructor(address oracle) OCR2DRClient(oracle) {}

  function sendSimpleRequestWithJavaScript(string memory sourceCode, uint64 subscriptionId)
    public
    returns (bytes32 requestId)
  {
    OCR2DR.Request memory request;
    request.initializeRequestForInlineJavaScript(sourceCode);
    requestId = sendRequest(request, subscriptionId, 20_000, tx.gasprice);
    emit SendRequestInvoked(requestId, sourceCode, subscriptionId);
  }

  function estimateJuelCost(string memory sourceCode, uint64 subscriptionId) public view returns (uint96) {
    OCR2DR.Request memory request;
    request.initializeRequestForInlineJavaScript(sourceCode);
    return estimateCost(request, subscriptionId, 20_000, 4_388_265);
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
