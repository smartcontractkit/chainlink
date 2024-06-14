// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsRequest} from "../../../dev/v1_0_0/libraries/FunctionsRequest.sol";
import {FunctionsClient} from "../../../dev/v1_0_0/FunctionsClient.sol";

contract FunctionsClientWithEmptyCallback is FunctionsClient {
  using FunctionsRequest for FunctionsRequest.Request;

  event SendRequestInvoked(bytes32 requestId, string sourceCode, uint64 subscriptionId);
  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  constructor(address router) FunctionsClient(router) {}

  function sendSimpleRequestWithJavaScript(
    string memory sourceCode,
    uint64 subscriptionId,
    bytes32 donId,
    uint32 callbackGasLimit
  ) public returns (bytes32 requestId) {
    FunctionsRequest.Request memory request;
    request.initializeRequestForInlineJavaScript(sourceCode);
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);
    requestId = _sendRequest(requestData, subscriptionId, callbackGasLimit, donId);
    emit SendRequestInvoked(requestId, sourceCode, subscriptionId);
  }

  function fulfillRequest(bytes32 /*requestId*/, bytes memory /*response*/, bytes memory /*err*/) internal override {
    // Do nothing
  }
}
