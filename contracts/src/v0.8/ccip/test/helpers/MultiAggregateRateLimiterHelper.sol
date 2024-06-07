// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MultiAggregateRateLimiter} from "../../MultiAggregateRateLimiter.sol";
import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";
import {Client} from "../../libraries/Client.sol";

contract MultiAggregateRateLimiterHelper is MultiAggregateRateLimiter {
  constructor(
    address priceRegistry,
    address[] memory authorizedCallers
  ) MultiAggregateRateLimiter(priceRegistry, authorizedCallers) {}

  function getTokenValue(Client.EVMTokenAmount memory tokenAmount) public view returns (uint256) {
    return _getTokenValue(tokenAmount);
  }
}
