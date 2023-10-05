// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouter} from "../../dev/v1_0_0/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/v1_0_0/FunctionsSubscriptions.sol";
import {FunctionsRequest} from "../../dev/v1_0_0/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_0_0/libraries/FunctionsResponse.sol";

import {FunctionsSubscriptionSetup} from "./Setup.t.sol";

/// @notice #constructor
contract FunctionsClient_Constructor {

}

/// @notice #_sendRequest
contract FunctionsClient__SendRequest is FunctionsSubscriptionSetup {
  // TODO: make contract internal function helper

  function test__SendRequest_RevertIfInvalidCallbackGasLimit() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

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

/// @notice #fulfillRequest
contract FunctionsClient_FulfillRequest {

}

/// @notice #handleOracleFulfillment
contract FunctionsClient_HandleOracleFulfillment {

}
