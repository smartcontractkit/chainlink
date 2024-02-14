// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Client} from "../../libraries/Client.sol";
import {AggregateRateLimiterHelper} from "../helpers/AggregateRateLimiterHelper.sol";
import {AggregateRateLimiter} from "../../AggregateRateLimiter.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";

import {BaseTest, stdError} from "../BaseTest.t.sol";

contract AggregateTokenLimiterSetup is BaseTest, PriceRegistrySetup {
  AggregateRateLimiterHelper internal s_rateLimiter;
  RateLimiter.Config internal s_config;

  address internal immutable TOKEN = 0x21118E64E1fB0c487F25Dd6d3601FF6af8D32E4e;
  uint224 internal constant TOKEN_PRICE = 4e18;

  function setUp() public virtual override(BaseTest, PriceRegistrySetup) {
    BaseTest.setUp();
    PriceRegistrySetup.setUp();

    Internal.PriceUpdates memory priceUpdates = getSingleTokenPriceUpdateStruct(TOKEN, TOKEN_PRICE);
    s_priceRegistry.updatePrices(priceUpdates);

    s_config = RateLimiter.Config({isEnabled: true, rate: 5, capacity: 100});
    s_rateLimiter = new AggregateRateLimiterHelper(s_config);
    s_rateLimiter.setAdmin(ADMIN);
  }
}

/// @notice #constructor
contract AggregateTokenLimiter_constructor is AggregateTokenLimiterSetup {
  function testConstructorSuccess() public {
    assertEq(ADMIN, s_rateLimiter.getTokenLimitAdmin());
    assertEq(OWNER, s_rateLimiter.owner());

    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity, bucket.tokens);
    assertEq(s_config.isEnabled, bucket.isEnabled);
    assertEq(BLOCK_TIME, bucket.lastUpdated);
  }
}

/// @notice #getTokenLimitAdmin
contract AggregateTokenLimiter_getTokenLimitAdmin is AggregateTokenLimiterSetup {
  function testGetTokenLimitAdminSuccess() public {
    assertEq(ADMIN, s_rateLimiter.getTokenLimitAdmin());
  }
}

/// @notice #setAdmin
contract AggregateTokenLimiter_setAdmin is AggregateTokenLimiterSetup {
  event AdminSet(address newAdmin);

  function testOwnerSuccess() public {
    vm.expectEmit();
    emit AdminSet(STRANGER);

    s_rateLimiter.setAdmin(STRANGER);
    assertEq(STRANGER, s_rateLimiter.getTokenLimitAdmin());
  }

  // Reverts

  function testOnlyOwnerOrAdminReverts() public {
    changePrank(STRANGER);
    vm.expectRevert(RateLimiter.OnlyCallableByAdminOrOwner.selector);

    s_rateLimiter.setAdmin(STRANGER);
  }
}

/// @notice #getTokenBucket
contract AggregateTokenLimiter_getTokenBucket is AggregateTokenLimiterSetup {
  function testGetTokenBucketSuccess() public {
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity, bucket.tokens);
    assertEq(BLOCK_TIME, bucket.lastUpdated);
  }

  function testRefillSuccess() public {
    s_config.capacity = s_config.capacity * 2;
    s_rateLimiter.setRateLimiterConfig(s_config);

    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();

    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity / 2, bucket.tokens);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    uint256 warpTime = 4;
    vm.warp(BLOCK_TIME + warpTime);

    bucket = s_rateLimiter.currentRateLimiterState();

    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity / 2 + warpTime * s_config.rate, bucket.tokens);
    assertEq(BLOCK_TIME + warpTime, bucket.lastUpdated);

    vm.warp(BLOCK_TIME + warpTime * 100);

    // Bucket overflow
    bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.capacity, bucket.tokens);
  }

  // Reverts

  function testTimeUnderflowReverts() public {
    vm.warp(BLOCK_TIME - 1);

    vm.expectRevert(stdError.arithmeticError);
    s_rateLimiter.currentRateLimiterState();
  }
}

/// @notice #setRateLimiterConfig
contract AggregateTokenLimiter_setRateLimiterConfig is AggregateTokenLimiterSetup {
  event ConfigChanged(RateLimiter.Config config);

  function testOwnerSuccess() public {
    setConfig();
  }

  function testTokenLimitAdminSuccess() public {
    changePrank(ADMIN);
    setConfig();
  }

  function setConfig() private {
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);

    if (bucket.isEnabled) {
      s_config = RateLimiter.Config({isEnabled: false, rate: 0, capacity: 0});
    } else {
      s_config = RateLimiter.Config({isEnabled: true, rate: 100, capacity: 200});
    }

    vm.expectEmit();
    emit ConfigChanged(s_config);

    s_rateLimiter.setRateLimiterConfig(s_config);

    bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.isEnabled, bucket.isEnabled);
  }

  // Reverts

  function testOnlyOnlyCallableByAdminOrOwnerReverts() public {
    changePrank(STRANGER);

    vm.expectRevert(RateLimiter.OnlyCallableByAdminOrOwner.selector);

    s_rateLimiter.setRateLimiterConfig(s_config);
  }
}

/// @notice #_rateLimitValue
contract AggregateTokenLimiter__rateLimitValue is AggregateTokenLimiterSetup {
  event TokensConsumed(uint256 tokens);

  function testRateLimitValueSuccess_gas() public {
    vm.pauseGasMetering();
    // start from blocktime that does not equal rate limiter init timestamp
    vm.warp(BLOCK_TIME + 1);

    // 15 (tokens) * 4 (price) * 2 (number of times) > 100 (capacity)
    uint256 numberOfTokens = 15;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0].token = TOKEN;
    tokenAmounts[0].amount = numberOfTokens;

    vm.expectEmit();
    emit TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(tokenAmounts, s_priceRegistry);
    vm.pauseGasMetering();

    // Get the updated bucket status
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();
    // Assert the proper value has been taken out of the bucket
    assertEq(bucket.capacity - value, bucket.tokens);

    // Since value * 2 > bucket.capacity we cannot take it out twice.
    // Expect a revert when we try, with a wait time.
    uint256 waitTime = 4;
    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, waitTime, bucket.tokens)
    );
    s_rateLimiter.rateLimitValue(tokenAmounts, s_priceRegistry);

    // Move the block time forward by 10 so the bucket refills by 10 * rate
    vm.warp(BLOCK_TIME + 1 + waitTime);

    // The bucket has filled up enough so we can take out more tokens
    s_rateLimiter.rateLimitValue(tokenAmounts, s_priceRegistry);
    bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(bucket.capacity - value + waitTime * s_config.rate - value, bucket.tokens);
    vm.resumeGasMetering();
  }

  // Reverts

  function testUnknownTokenReverts() public {
    vm.expectRevert(abi.encodeWithSelector(AggregateRateLimiter.PriceNotFoundForToken.selector, address(0)));
    s_rateLimiter.rateLimitValue(new Client.EVMTokenAmount[](1), s_priceRegistry);
  }

  function testAggregateValueMaxCapacityExceededReverts() public {
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0].token = TOKEN;
    tokenAmounts[0].amount = 100;

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector,
        bucket.capacity,
        (tokenAmounts[0].amount * TOKEN_PRICE) / 1e18
      )
    );
    s_rateLimiter.rateLimitValue(tokenAmounts, s_priceRegistry);
  }
}
