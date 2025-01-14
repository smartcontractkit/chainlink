// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_setChainRateLimiterConfig is TokenPoolSetup {
  function testFuzz_SetChainRateLimiterConfig_Success(uint128 capacity, uint128 rate, uint32 newTime) public {
    // Cap the lower bound to 4 so 4/2 is still >= 2
    vm.assume(capacity >= 4);
    // Cap the lower bound to 2 so 2/2 is still >= 1
    rate = uint128(bound(rate, 2, capacity - 2));
    // Bucket updates only work on increasing time
    newTime = uint32(bound(newTime, block.timestamp + 1, type(uint32).max));
    vm.warp(newTime);

    uint256 oldOutboundTokens = s_tokenPool.getCurrentOutboundRateLimiterState(DEST_CHAIN_SELECTOR).tokens;
    uint256 oldInboundTokens = s_tokenPool.getCurrentInboundRateLimiterState(DEST_CHAIN_SELECTOR).tokens;

    RateLimiter.Config memory newOutboundConfig = RateLimiter.Config({isEnabled: true, capacity: capacity, rate: rate});
    RateLimiter.Config memory newInboundConfig =
      RateLimiter.Config({isEnabled: true, capacity: capacity / 2, rate: rate / 2});

    vm.expectEmit();
    emit RateLimiter.ConfigChanged(newOutboundConfig);
    vm.expectEmit();
    emit RateLimiter.ConfigChanged(newInboundConfig);
    vm.expectEmit();
    emit TokenPool.ChainConfigured(DEST_CHAIN_SELECTOR, newOutboundConfig, newInboundConfig);

    s_tokenPool.setChainRateLimiterConfig(DEST_CHAIN_SELECTOR, newOutboundConfig, newInboundConfig);

    uint256 expectedTokens = RateLimiter._min(newOutboundConfig.capacity, oldOutboundTokens);

    RateLimiter.TokenBucket memory bucket = s_tokenPool.getCurrentOutboundRateLimiterState(DEST_CHAIN_SELECTOR);
    assertEq(bucket.capacity, newOutboundConfig.capacity);
    assertEq(bucket.rate, newOutboundConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);

    expectedTokens = RateLimiter._min(newInboundConfig.capacity, oldInboundTokens);

    bucket = s_tokenPool.getCurrentInboundRateLimiterState(DEST_CHAIN_SELECTOR);
    assertEq(bucket.capacity, newInboundConfig.capacity);
    assertEq(bucket.rate, newInboundConfig.rate);
    assertEq(bucket.tokens, expectedTokens);
    assertEq(bucket.lastUpdated, newTime);
  }

  // Reverts

  function test_RevertWhen_OnlyOwnerOrRateLimitAdmin() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));
    s_tokenPool.setChainRateLimiterConfig(
      DEST_CHAIN_SELECTOR, _getOutboundRateLimiterConfig(), _getInboundRateLimiterConfig()
    );
  }

  function test_RevertWhen_NonExistentChain() public {
    uint64 wrongChainSelector = 9084102894;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, wrongChainSelector));
    s_tokenPool.setChainRateLimiterConfig(
      wrongChainSelector, _getOutboundRateLimiterConfig(), _getInboundRateLimiterConfig()
    );
  }
}
