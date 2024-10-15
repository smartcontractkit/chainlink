// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MultiAggregateRateLimiter} from "../../MultiAggregateRateLimiter.sol";
import {Client} from "../../libraries/Client.sol";

contract MultiAggregateRateLimiterHelper is MultiAggregateRateLimiter {
  constructor(
    address feeQuoter,
    address[] memory authorizedCallers
  ) MultiAggregateRateLimiter(feeQuoter, authorizedCallers) {}

  function getTokenValue(
    Client.EVMTokenAmount memory tokenAmount
  ) public view returns (uint256) {
    return _getTokenValue(tokenAmount);
  }
}
