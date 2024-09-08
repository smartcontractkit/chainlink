// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {RateLimiter} from "../../libraries/RateLimiter.sol";

contract RateLimiterHelper {
  using RateLimiter for RateLimiter.TokenBucket;

  RateLimiter.TokenBucket internal s_rateLimiter;

  constructor(RateLimiter.Config memory config) {
    s_rateLimiter = RateLimiter.TokenBucket({
      rate: config.rate,
      capacity: config.capacity,
      tokens: config.capacity,
      lastUpdated: uint32(block.timestamp),
      isEnabled: config.isEnabled
    });
  }

  function consume(uint256 requestTokens, address tokenAddress) external {
    s_rateLimiter._consume(requestTokens, tokenAddress);
  }

  function currentTokenBucketState() external view returns (RateLimiter.TokenBucket memory) {
    return s_rateLimiter._currentTokenBucketState();
  }

  function setTokenBucketConfig(RateLimiter.Config memory config) external {
    s_rateLimiter._setTokenBucketConfig(config);
  }

  function getRateLimiter() external view returns (RateLimiter.TokenBucket memory) {
    return s_rateLimiter;
  }
}
