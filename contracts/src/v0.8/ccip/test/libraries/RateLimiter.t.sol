// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {RateLimiterHelper} from "../helpers/RateLimiterHelper.sol";
import {Test} from "forge-std/Test.sol";

contract RateLimiterSetup is Test {
  RateLimiterHelper internal s_helper;
  RateLimiter.Config internal s_config;

  uint256 internal constant BLOCK_TIME = 1234567890;

  function setUp() public virtual {
    s_config = RateLimiter.Config({isEnabled: true, rate: 5, capacity: 100});
    s_helper = new RateLimiterHelper(s_config);
  }
}

contract RateLimiter_constructor is RateLimiterSetup {
  function test_Constructor_Success() public view {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(s_config.capacity, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
    assertEq(BLOCK_TIME, rateLimiter.lastUpdated);
  }
}

contract RateLimiter_setTokenBucketConfig is RateLimiterSetup {
  function test_SetRateLimiterConfig_Success() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);

    s_config =
      RateLimiter.Config({isEnabled: true, rate: uint128(rateLimiter.rate * 2), capacity: rateLimiter.capacity * 8});

    vm.expectEmit();
    emit RateLimiter.ConfigChanged(s_config);

    s_helper.setTokenBucketConfig(s_config);

    rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(s_config.capacity / 8, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
    assertEq(BLOCK_TIME, rateLimiter.lastUpdated);
  }
}

contract RateLimiter_currentTokenBucketState is RateLimiterSetup {
  function test_CurrentTokenBucketState_Success() public {
    RateLimiter.TokenBucket memory bucket = s_helper.currentTokenBucketState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity, bucket.tokens);
    assertEq(s_config.isEnabled, bucket.isEnabled);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    s_config = RateLimiter.Config({isEnabled: true, rate: uint128(bucket.rate * 2), capacity: bucket.capacity * 8});

    s_helper.setTokenBucketConfig(s_config);

    bucket = s_helper.currentTokenBucketState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity / 8, bucket.tokens);
    assertEq(s_config.isEnabled, bucket.isEnabled);
    assertEq(BLOCK_TIME, bucket.lastUpdated);
  }

  function test_Refill_Success() public {
    RateLimiter.TokenBucket memory bucket = s_helper.currentTokenBucketState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity, bucket.tokens);
    assertEq(s_config.isEnabled, bucket.isEnabled);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    s_config = RateLimiter.Config({isEnabled: true, rate: uint128(bucket.rate * 2), capacity: bucket.capacity * 8});

    s_helper.setTokenBucketConfig(s_config);

    bucket = s_helper.currentTokenBucketState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity / 8, bucket.tokens);
    assertEq(s_config.isEnabled, bucket.isEnabled);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    uint256 warpTime = 4;
    vm.warp(BLOCK_TIME + warpTime);

    bucket = s_helper.currentTokenBucketState();

    assertEq(s_config.capacity / 8 + warpTime * s_config.rate, bucket.tokens);

    vm.warp(BLOCK_TIME + warpTime * 100);

    // Bucket overflow
    bucket = s_helper.currentTokenBucketState();
    assertEq(s_config.capacity, bucket.tokens);
  }
}

contract RateLimiter_consume is RateLimiterSetup {
  address internal s_token = address(100);

  function test_ConsumeAggregateValue_Success() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(s_config.capacity, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
    assertEq(BLOCK_TIME, rateLimiter.lastUpdated);

    uint256 requestTokens = 50;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(requestTokens);

    s_helper.consume(requestTokens, address(0));

    rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(s_config.capacity - requestTokens, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
    assertEq(BLOCK_TIME, rateLimiter.lastUpdated);
  }

  function test_ConsumeTokens_Success() public {
    uint256 requestTokens = 50;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(requestTokens);

    s_helper.consume(requestTokens, s_token);
  }

  function test_Refill_Success() public {
    uint256 requestTokens = 50;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(requestTokens);

    s_helper.consume(requestTokens, address(0));

    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(s_config.capacity - requestTokens, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
    assertEq(BLOCK_TIME, rateLimiter.lastUpdated);

    uint256 warpTime = 4;
    vm.warp(BLOCK_TIME + warpTime);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(requestTokens);

    s_helper.consume(requestTokens, address(0));

    rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(s_config.capacity - requestTokens * 2 + warpTime * s_config.rate, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
    assertEq(BLOCK_TIME + warpTime, rateLimiter.lastUpdated);
  }

  function test_ConsumeUnlimited_Success() public {
    s_helper.consume(0, address(0));

    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.capacity, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);

    RateLimiter.Config memory disableConfig = RateLimiter.Config({isEnabled: false, rate: 0, capacity: 0});

    s_helper.setTokenBucketConfig(disableConfig);

    uint256 requestTokens = 50;
    s_helper.consume(requestTokens, address(0));

    rateLimiter = s_helper.getRateLimiter();
    assertEq(disableConfig.capacity, rateLimiter.tokens);
    assertEq(disableConfig.isEnabled, rateLimiter.isEnabled);

    s_helper.setTokenBucketConfig(s_config);

    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 10, 0));
    s_helper.consume(requestTokens, address(0));

    rateLimiter = s_helper.getRateLimiter();
    assertEq(s_config.rate, rateLimiter.rate);
    assertEq(s_config.capacity, rateLimiter.capacity);
    assertEq(0, rateLimiter.tokens);
    assertEq(s_config.isEnabled, rateLimiter.isEnabled);
  }

  // Reverts

  function test_AggregateValueMaxCapacityExceeded_Revert() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector, rateLimiter.capacity, rateLimiter.capacity + 1
      )
    );
    s_helper.consume(rateLimiter.capacity + 1, address(0));
  }

  function test_TokenMaxCapacityExceeded_Revert() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.TokenMaxCapacityExceeded.selector, rateLimiter.capacity, rateLimiter.capacity + 1, s_token
      )
    );
    s_helper.consume(rateLimiter.capacity + 1, s_token);
  }

  function test_ConsumingMoreThanUint128_Revert() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();

    uint256 request = uint256(type(uint128).max) + 1;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.AggregateValueMaxCapacityExceeded.selector, rateLimiter.capacity, request)
    );
    s_helper.consume(request, address(0));
  }

  function test_AggregateValueRateLimitReached_Revert() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();

    uint256 overLimit = 20;
    uint256 requestTokens1 = rateLimiter.capacity / 2;
    uint256 requestTokens2 = rateLimiter.capacity / 2 + overLimit;

    uint256 waitInSeconds = overLimit / rateLimiter.rate;

    s_helper.consume(requestTokens1, address(0));

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueRateLimitReached.selector, waitInSeconds, rateLimiter.capacity - requestTokens1
      )
    );
    s_helper.consume(requestTokens2, address(0));
  }

  function test_TokenRateLimitReached_Revert() public {
    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();

    uint256 overLimit = 20;
    uint256 requestTokens1 = rateLimiter.capacity / 2;
    uint256 requestTokens2 = rateLimiter.capacity / 2 + overLimit;

    uint256 waitInSeconds = overLimit / rateLimiter.rate;

    s_helper.consume(requestTokens1, s_token);

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.TokenRateLimitReached.selector, waitInSeconds, rateLimiter.capacity - requestTokens1, s_token
      )
    );
    s_helper.consume(requestTokens2, s_token);
  }

  function test_RateLimitReachedOverConsecutiveBlocks_Revert() public {
    uint256 initBlockTime = BLOCK_TIME + 10000;
    vm.warp(initBlockTime);

    RateLimiter.TokenBucket memory rateLimiter = s_helper.getRateLimiter();

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(rateLimiter.capacity);

    s_helper.consume(rateLimiter.capacity, address(0));

    vm.warp(initBlockTime + 1);

    // Over rate limit by 1, force 1 second wait
    uint256 overLimit = 1;

    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 1, rateLimiter.rate));
    s_helper.consume(rateLimiter.rate + overLimit, address(0));
  }
}
