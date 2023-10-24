// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRouter} from "../../../dev/v1_X/FunctionsRouter.sol";

/// @title Functions Router Test Harness
/// @notice Contract to expose internal functions for testing purposes
contract FunctionsRouterHarness is FunctionsRouter {
  constructor(address linkToken, Config memory config) FunctionsRouter(linkToken, config) {}

  function getMaxConsumers_HARNESS() external view returns (uint16) {
    return super._getMaxConsumers();
  }

  function getSubscriptionDepositDetails_HARNESS() external view returns (uint16, uint72) {
    return super._getSubscriptionDepositDetails();
  }

  function whenNotPaused_HARNESS() external view {
    return super._whenNotPaused();
  }

  function onlyRouterOwner_HARNESS() external view {
    return super._onlyRouterOwner();
  }

  function onlySenderThatAcceptedToS_HARNESS() external view {
    return super._onlySenderThatAcceptedToS();
  }
}
