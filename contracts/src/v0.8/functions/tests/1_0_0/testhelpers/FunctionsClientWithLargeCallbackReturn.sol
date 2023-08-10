// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsRequest} from "../../../dev/1_0_0/libraries/FunctionsRequest.sol";
import {FunctionsClient} from "../../../dev/1_0_0/FunctionsClient.sol";

contract FunctionsClientWithLargeCallbackReturn is FunctionsClient {
  using FunctionsRequest for FunctionsRequest.Request;

  event SendRequestInvoked(bytes32 requestId, string sourceCode, uint64 subscriptionId);
  event FulfillRequestInvoked(bytes32 requestId, bytes response, bytes err);

  error AnExtremelyLargeErrorThatHasALotOfDataInItAndAnExtremelyLongNameToGoWithTheLargeAmountOfData(
    bytes lotsOfData,
    bytes moreData,
    string andAString
  );

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
    revert AnExtremelyLargeErrorThatHasALotOfDataInItAndAnExtremelyLongNameToGoWithTheLargeAmountOfData(
      abi.encode(
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur."
      ),
      abi.encode(
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur."
      ),
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur."
    );
  }
}
