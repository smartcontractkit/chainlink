// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_setChainRateLimiterConfigs is TokenPoolSetup {
  uint64 internal s_remoteChainSelector;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();

    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = abi.encode(address(2));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    s_remoteChainSelector = 123124;
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: s_remoteChainSelector,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: abi.encode(address(3)),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);
  }

  function testFuzz_SetChainRateLimiterConfigs_Success(uint128 capacity, uint128 rate, uint32 newTime) public {
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

    uint64[] memory chainSelectors = new uint64[](1);
    chainSelectors[0] = DEST_CHAIN_SELECTOR;

    RateLimiter.Config[] memory newOutboundConfigs = new RateLimiter.Config[](1);
    newOutboundConfigs[0] = newOutboundConfig;

    RateLimiter.Config[] memory newInboundConfigs = new RateLimiter.Config[](1);
    newInboundConfigs[0] = newInboundConfig;

    vm.expectEmit();
    emit RateLimiter.ConfigChanged(newOutboundConfig);
    vm.expectEmit();
    emit RateLimiter.ConfigChanged(newInboundConfig);
    vm.expectEmit();
    emit TokenPool.ChainConfigured(DEST_CHAIN_SELECTOR, newOutboundConfig, newInboundConfig);

    s_tokenPool.setChainRateLimiterConfigs(chainSelectors, newOutboundConfigs, newInboundConfigs);

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
    uint64[] memory chainSelectors = new uint64[](1);
    chainSelectors[0] = DEST_CHAIN_SELECTOR;

    RateLimiter.Config[] memory newOutboundConfigs = new RateLimiter.Config[](1);
    newOutboundConfigs[0] = _getOutboundRateLimiterConfig();

    RateLimiter.Config[] memory newInboundConfigs = new RateLimiter.Config[](1);
    newInboundConfigs[0] = _getInboundRateLimiterConfig();

    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));
    s_tokenPool.setChainRateLimiterConfigs(chainSelectors, newOutboundConfigs, newInboundConfigs);
  }

  function test_RevertWhen_NonExistentChain() public {
    uint64 wrongChainSelector = 9084102894;

    uint64[] memory chainSelectors = new uint64[](1);
    chainSelectors[0] = wrongChainSelector;

    RateLimiter.Config[] memory newOutboundConfigs = new RateLimiter.Config[](1);
    RateLimiter.Config[] memory newInboundConfigs = new RateLimiter.Config[](1);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, wrongChainSelector));
    s_tokenPool.setChainRateLimiterConfigs(chainSelectors, newOutboundConfigs, newInboundConfigs);
  }

  function test_RevertWhen_MismatchedArrayLengths() public {
    uint64[] memory chainSelectors = new uint64[](1);

    RateLimiter.Config[] memory newOutboundConfigs = new RateLimiter.Config[](1);
    RateLimiter.Config[] memory newInboundConfigs = new RateLimiter.Config[](2);

    // test mismatched array lengths between rate limiters
    vm.expectRevert(abi.encodeWithSelector(TokenPool.MismatchedArrayLengths.selector));
    s_tokenPool.setChainRateLimiterConfigs(chainSelectors, newOutboundConfigs, newInboundConfigs);

    newInboundConfigs = new RateLimiter.Config[](1);
    chainSelectors = new uint64[](2);

    // test mismatched array lengths between chain selectors and rate limiters
    vm.expectRevert(abi.encodeWithSelector(TokenPool.MismatchedArrayLengths.selector));
    s_tokenPool.setChainRateLimiterConfigs(chainSelectors, newOutboundConfigs, newInboundConfigs);
  }
}
