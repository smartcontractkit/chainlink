// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ZKSyncFunctionsRouter} from "../../../v1_3_0_zksync/ZKSyncFunctionsRouter.sol";

import {FunctionsRouter} from "../../../v1_0_0/FunctionsRouter.sol";

/// @title ZKSync Functions Router Test Harness
/// @notice Contract to expose internal functions for testing purposes
contract ZKSyncFunctionsRouterHarness is ZKSyncFunctionsRouter {
  constructor(address linkToken, FunctionsRouter.Config memory config) ZKSyncFunctionsRouter(linkToken, config) {}

  function exposed_callback(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint32 callbackGasLimit,
    address client
  ) public returns (CallbackResult memory) {
    // simply call the internal `_callback` method
    return super._callback(requestId, response, err, callbackGasLimit, client);
  }
}
