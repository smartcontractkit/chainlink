// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsClientUpgradeHelper} from "./FunctionsClientUpgradeHelper.sol";
import {FunctionsResponse} from "../../../dev/v1_X/libraries/FunctionsResponse.sol";

/// @title Functions Client Test Harness
/// @notice Contract to expose internal functions for testing purposes
contract FunctionsClientHarness is FunctionsClientUpgradeHelper {
  constructor(address router) FunctionsClientUpgradeHelper(router) {}

  function getRouter_HARNESS() external view returns (address) {
    return address(i_functionsRouter);
  }

  function sendRequest_HARNESS(
    bytes memory data,
    uint64 subscriptionId,
    uint32 callbackGasLimit,
    bytes32 donId
  ) external returns (bytes32) {
    return super._sendRequest(data, subscriptionId, callbackGasLimit, donId);
  }
}
