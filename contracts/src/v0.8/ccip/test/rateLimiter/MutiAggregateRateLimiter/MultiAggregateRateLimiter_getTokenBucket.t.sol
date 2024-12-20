// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";

import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";
import {stdError} from "forge-std/Test.sol";

contract MultiAggregateRateLimiter_getTokenBucket is MultiAggregateRateLimiterSetup {
  function test_GetTokenBucket() public view {
    RateLimiter.TokenBucket memory bucketInbound = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    _assertConfigWithTokenBucketEquality(s_rateLimiterConfig1, bucketInbound);
    assertEq(BLOCK_TIME, bucketInbound.lastUpdated);

    RateLimiter.TokenBucket memory bucketOutbound = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    _assertConfigWithTokenBucketEquality(s_rateLimiterConfig1, bucketOutbound);
    assertEq(BLOCK_TIME, bucketOutbound.lastUpdated);
  }

  function test_Refill() public {
    s_rateLimiterConfig1.capacity = s_rateLimiterConfig1.capacity * 2;

    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig1
    });

    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    assertEq(s_rateLimiterConfig1.rate, bucket.rate);
    assertEq(s_rateLimiterConfig1.capacity, bucket.capacity);
    assertEq(s_rateLimiterConfig1.capacity / 2, bucket.tokens);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    uint256 warpTime = 4;
    vm.warp(BLOCK_TIME + warpTime);

    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    assertEq(s_rateLimiterConfig1.rate, bucket.rate);
    assertEq(s_rateLimiterConfig1.capacity, bucket.capacity);
    assertEq(s_rateLimiterConfig1.capacity / 2 + warpTime * s_rateLimiterConfig1.rate, bucket.tokens);
    assertEq(BLOCK_TIME + warpTime, bucket.lastUpdated);

    vm.warp(BLOCK_TIME + warpTime * 100);

    // Bucket overflow
    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(s_rateLimiterConfig1.capacity, bucket.tokens);
  }

  // Reverts

  function test_RevertWhen_TimeUnderflow() public {
    vm.warp(BLOCK_TIME - 1);

    vm.expectRevert(stdError.arithmeticError);
    s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
  }
}
