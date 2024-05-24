// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsClient} from "../../dev/v1_X/FunctionsClient.sol";
import {FunctionsRouter} from "../../dev/v1_X/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/v1_X/FunctionsSubscriptions.sol";
import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";

import {FunctionsClientSetup, FunctionsSubscriptionSetup, FunctionsClientRequestSetup} from "./Setup.t.sol";

/// @notice #constructor
contract FunctionsClient_Constructor is FunctionsClientSetup {
  function test_Constructor_Success() public {
    assertEq(address(s_functionsRouter), s_functionsClient.getRouter_HARNESS());
  }
}

/// @notice #_sendRequest
contract FunctionsClient__SendRequest is FunctionsSubscriptionSetup {
  function test__SendRequest_RevertIfInvalidCallbackGasLimit() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest._initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest._encodeCBOR(request);

    uint8 MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;
    bytes32 subscriptionFlags = s_functionsRouter.getFlags(s_subscriptionId);
    uint8 callbackGasLimitsIndexSelector = uint8(subscriptionFlags[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    uint32[] memory _maxCallbackGasLimits = config.maxCallbackGasLimits;
    uint32 maxCallbackGasLimit = _maxCallbackGasLimits[callbackGasLimitsIndexSelector];

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.GasLimitTooBig.selector, maxCallbackGasLimit));
    s_functionsClient.sendRequestBytes(requestData, s_subscriptionId, 500_000, s_donId);
  }
}

/// @notice #handleOracleFulfillment
contract FunctionsClient_HandleOracleFulfillment is FunctionsClientRequestSetup {
  function test_HandleOracleFulfillment_RevertIfNotRouter() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsClient.OnlyRouterCanFulfill.selector);
    s_functionsClient.handleOracleFulfillment(s_requests[1].requestId, new bytes(0), new bytes(0));
  }

  event RequestFulfilled(bytes32 indexed id);
  event ResponseReceived(bytes32 indexed requestId, bytes result, bytes err);

  function test_HandleOracleFulfillment_Success() public {
    // Send as Router
    vm.stopPrank();
    vm.startPrank(address(s_functionsRouter));

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ResponseReceived(s_requests[1].requestId, new bytes(0), new bytes(0));
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestFulfilled(s_requests[1].requestId);

    s_functionsClient.handleOracleFulfillment(s_requests[1].requestId, new bytes(0), new bytes(0));
  }
}
