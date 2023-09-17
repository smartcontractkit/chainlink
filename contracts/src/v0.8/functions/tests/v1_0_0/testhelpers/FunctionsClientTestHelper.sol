// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {ITermsOfServiceAllowList} from "../../../dev/v1_0_0/accessControl/interfaces/ITermsOfServiceAllowList.sol";
import {IFunctionsSubscriptions} from "../../../dev/v1_0_0/interfaces/IFunctionsSubscriptions.sol";

import {FunctionsRequest} from "../../../dev/v1_0_0/libraries/FunctionsRequest.sol";
import {FunctionsClient} from "../../../dev/v1_0_0/FunctionsClient.sol";

contract FunctionsClientTestHelper is FunctionsClient {
  using FunctionsRequest for FunctionsRequest.Request;

  event SendRequestInvoked(bytes32 requestId, string sourceCode, uint64 subscriptionId);
  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  bool private s_revertFulfillRequest;
  bool private s_doInvalidOperation;
  bool private s_doInvalidReentrantOperation;
  bool private s_doValidReentrantOperation;

  uint64 private s_subscriptionId;
  bytes32 private s_donId;

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

  function sendRequestProposed(
    string memory sourceCode,
    uint64 subscriptionId,
    bytes32 donId
  ) public returns (bytes32 requestId) {
    FunctionsRequest.Request memory request;
    uint32 callbackGasLimit = 20_000;
    request.initializeRequestForInlineJavaScript(sourceCode);
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);
    requestId = i_router.sendRequestToProposed(
      subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      donId
    );
    emit RequestSent(requestId);
    emit SendRequestInvoked(requestId, sourceCode, subscriptionId);
  }

  function acceptTermsOfService(address acceptor, address recipient, bytes32 r, bytes32 s, uint8 v) external {
    bytes32 allowListId = i_router.getAllowListId();
    ITermsOfServiceAllowList allowList = ITermsOfServiceAllowList(i_router.getContractById(allowListId));
    allowList.acceptTermsOfService(acceptor, recipient, r, s, v);
  }

  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external {
    IFunctionsSubscriptions(address(i_router)).acceptSubscriptionOwnerTransfer(subscriptionId);
  }

  function fulfillRequest(bytes32 requestId, bytes memory response, bytes memory err) internal override {
    if (s_revertFulfillRequest) {
      revert("asked to revert");
    }
    if (s_doInvalidOperation) {
      uint256 x = 1;
      uint256 y = 0;
      x = x / y;
    }
    if (s_doValidReentrantOperation) {
      sendSimpleRequestWithJavaScript("somedata", s_subscriptionId, s_donId, 20_000);
    }
    if (s_doInvalidReentrantOperation) {
      IFunctionsSubscriptions(address(i_router)).cancelSubscription(s_subscriptionId, msg.sender);
    }
    emit FulfillRequestInvoked(requestId, response, err);
  }

  function setRevertFulfillRequest(bool on) external {
    s_revertFulfillRequest = on;
  }

  function setDoInvalidOperation(bool on) external {
    s_doInvalidOperation = on;
  }

  function setDoInvalidReentrantOperation(bool on, uint64 subscriptionId) external {
    s_doInvalidReentrantOperation = on;
    s_subscriptionId = subscriptionId;
  }

  function setDoValidReentrantOperation(bool on, uint64 subscriptionId, bytes32 donId) external {
    s_doValidReentrantOperation = on;
    s_subscriptionId = subscriptionId;
    s_donId = donId;
  }
}
