// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsSubscriptions} from "../../../dev/v1_X/FunctionsSubscriptions.sol";

/// @title Functions Subscriptions Test Harness
/// @notice Contract to expose internal functions for testing purposes
contract FunctionsSubscriptionsHarness is FunctionsSubscriptions {
  constructor(address link) FunctionsSubscriptions(link) {}

  function markRequestInFlight_HARNESS(address client, uint64 subscriptionId, uint96 estimatedTotalCostJuels) external {
    return super._markRequestInFlight(client, subscriptionId, estimatedTotalCostJuels);
  }

  function pay_HARNESS(
    uint64 subscriptionId,
    uint96 estimatedTotalCostJuels,
    address client,
    uint96 adminFee,
    uint96 juelsPerGas,
    uint96 gasUsed,
    uint96 costWithoutCallbackJuels
  ) external returns (Receipt memory) {
    return
      super._pay(
        subscriptionId,
        estimatedTotalCostJuels,
        client,
        adminFee,
        juelsPerGas,
        gasUsed,
        costWithoutCallbackJuels
      );
  }

  function isExistingSubscription_HARNESS(uint64 subscriptionId) external view {
    return super._isExistingSubscription(subscriptionId);
  }

  function isAllowedConsumer_HARNESS(address client, uint64 subscriptionId) external view {
    return super._isAllowedConsumer(client, subscriptionId);
  }

  // Overrides
  function _getMaxConsumers() internal view override returns (uint16) {}

  function _getSubscriptionDepositDetails() internal override returns (uint16, uint72) {}

  function _onlySenderThatAcceptedToS() internal override {}

  function _onlyRouterOwner() internal override {}

  function _whenNotPaused() internal override {}
}
