// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/ocr2dr/OCR2DRClient.sol";

contract OCR2DRClientTestHelper is OCR2DRClient {
  using OCR2DR for OCR2DR.Request;
  using OCR2DR for OCR2DR.HttpQuery;
  using OCR2DR for OCR2DR.HttpHeader;

  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  constructor(address oracle) OCR2DRClient(oracle) {}

  function sendSimpleRequestWithJavaScript(string memory sourceCode, uint256 subscriptionId) public returns (bytes32) {
    OCR2DR.Request memory request;
    request.initializeRequestForInlineJavaScript(sourceCode);
    return sendRequest(request, subscriptionId);
  }

  function fulfillRequest(
    bytes32 requestId,
    bytes memory response,
    bytes memory err
  ) internal override {
    emit FulfillRequestInvoked(requestId, response, err);
  }
}
