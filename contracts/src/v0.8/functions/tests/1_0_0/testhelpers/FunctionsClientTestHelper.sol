// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsClient, Functions} from "../../../dev/1_0_0/FunctionsClient.sol";
import {ITermsOfServiceAllowList} from "../../../dev/1_0_0/accessControl/interfaces/ITermsOfServiceAllowList.sol";

contract FunctionsClientTestHelper is FunctionsClient {
  using Functions for Functions.Request;

  event SendRequestInvoked(bytes32 requestId, string sourceCode, uint64 subscriptionId);
  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  bool private s_revertFulfillRequest;
  bool private s_doInvalidOperation;

  constructor(address router) FunctionsClient(router) {}

  function sendSimpleRequestWithJavaScript(
    string memory sourceCode,
    uint64 subscriptionId,
    bytes32 donId
  ) public returns (bytes32 requestId) {
    Functions.Request memory request;
    uint32 callbackGasLimit = 20_000;
    request.initializeRequestForInlineJavaScript(sourceCode);
    requestId = _sendRequest(request, subscriptionId, callbackGasLimit, donId);
    emit SendRequestInvoked(requestId, sourceCode, subscriptionId);
  }

  function sendRequestProposed(
    string memory sourceCode,
    uint64 subscriptionId,
    bytes32 donId
  ) public returns (bytes32 requestId) {
    Functions.Request memory request;
    uint32 callbackGasLimit = 20_000;
    request.initializeRequestForInlineJavaScript(sourceCode);
    bytes memory requestData = Functions.encodeCBOR(request);
    requestId = bytes32(
      s_router.validateProposedContracts(
        donId,
        abi.encode(subscriptionId, requestData, Functions.REQUEST_DATA_VERSION, callbackGasLimit)
      )
    );
    emit RequestSent(requestId);
    emit SendRequestInvoked(requestId, sourceCode, subscriptionId);
  }

  function acceptTermsOfService(address acceptor, address recipient, bytes calldata proof) external {
    bytes32 allowListId = s_router.getAllowListId();
    ITermsOfServiceAllowList allowList = ITermsOfServiceAllowList(s_router.getContractById(allowListId));
    allowList.acceptTermsOfService(acceptor, recipient, proof);
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
    emit FulfillRequestInvoked(requestId, response, err);
  }

  function setRevertFulfillRequest(bool on) external {
    s_revertFulfillRequest = on;
  }

  function setDoInvalidOperation(bool on) external {
    s_doInvalidOperation = on;
  }
}
