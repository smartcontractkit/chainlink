// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../../AggregateRateLimiter.sol";

contract AggregateRateLimiterHelper is AggregateRateLimiter {
  constructor(RateLimiter.Config memory config) AggregateRateLimiter(config) {}

  function rateLimitValue(Client.EVMTokenAmount[] memory tokenAmounts, IPriceRegistry priceRegistry) public {
    _rateLimitValue(tokenAmounts, priceRegistry);
  }
}
