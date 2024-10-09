// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";

import {AggregateRateLimiter} from "../../AggregateRateLimiter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";

import {FeeQuoterSetup} from "../feeQuoter/FeeQuoterSetup.t.sol";
import {AggregateRateLimiterHelper} from "../helpers/AggregateRateLimiterHelper.sol";

import {stdError} from "forge-std/Test.sol";

contract AggregateTokenLimiterSetup is FeeQuoterSetup {
  AggregateRateLimiterHelper internal s_rateLimiter;
  RateLimiter.Config internal s_config;

  address internal immutable TOKEN = 0x21118E64E1fB0c487F25Dd6d3601FF6af8D32E4e;
  uint224 internal constant TOKEN_PRICE = 4e18;

  function setUp() public virtual override {
    FeeQuoterSetup.setUp();

    Internal.PriceUpdates memory priceUpdates = _getSingleTokenPriceUpdateStruct(TOKEN, TOKEN_PRICE);
    s_feeQuoter.updatePrices(priceUpdates);

    s_config = RateLimiter.Config({isEnabled: true, rate: 5, capacity: 100});
    s_rateLimiter = new AggregateRateLimiterHelper(s_config);
    s_rateLimiter.setAdmin(ADMIN);
  }
}

contract AggregateTokenLimiter_constructor is AggregateTokenLimiterSetup {
  function test_Constructor_Success() public view {
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

contract AggregateTokenLimiter_getTokenLimitAdmin is AggregateTokenLimiterSetup {
  function test_GetTokenLimitAdmin_Success() public view {
    assertEq(ADMIN, s_rateLimiter.getTokenLimitAdmin());
  }
}

contract AggregateTokenLimiter_setAdmin is AggregateTokenLimiterSetup {
  function test_Owner_Success() public {
    vm.expectEmit();
    emit AggregateRateLimiter.AdminSet(STRANGER);

    s_rateLimiter.setAdmin(STRANGER);
    assertEq(STRANGER, s_rateLimiter.getTokenLimitAdmin());
  }

  // Reverts

  function test_OnlyOwnerOrAdmin_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(RateLimiter.OnlyCallableByAdminOrOwner.selector);

    s_rateLimiter.setAdmin(STRANGER);
  }
}

contract AggregateTokenLimiter_getTokenBucket is AggregateTokenLimiterSetup {
  function test_GetTokenBucket_Success() public view {
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.capacity, bucket.tokens);
    assertEq(BLOCK_TIME, bucket.lastUpdated);
  }

  function test_Refill_Success() public {
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

  function test_TimeUnderflow_Revert() public {
    vm.warp(BLOCK_TIME - 1);

    vm.expectRevert(stdError.arithmeticError);
    s_rateLimiter.currentRateLimiterState();
  }
}

contract AggregateTokenLimiter_setRateLimiterConfig is AggregateTokenLimiterSetup {
  function test_Owner_Success() public {
    setConfig();
  }

  function test_TokenLimitAdmin_Success() public {
    vm.startPrank(ADMIN);
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
    emit RateLimiter.ConfigChanged(s_config);

    s_rateLimiter.setRateLimiterConfig(s_config);

    bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(s_config.rate, bucket.rate);
    assertEq(s_config.capacity, bucket.capacity);
    assertEq(s_config.isEnabled, bucket.isEnabled);
  }

  // Reverts

  function test_OnlyOnlyCallableByAdminOrOwner_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(RateLimiter.OnlyCallableByAdminOrOwner.selector);

    s_rateLimiter.setRateLimiterConfig(s_config);
  }
}

contract AggregateTokenLimiter_rateLimitValue is AggregateTokenLimiterSetup {
  function test_RateLimitValueSuccess_gas() public {
    vm.pauseGasMetering();
    // start from blocktime that does not equal rate limiter init timestamp
    vm.warp(BLOCK_TIME + 1);

    // 15 (tokens) * 4 (price) * 2 (number of times) > 100 (capacity)
    uint256 numberOfTokens = 15;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(value);
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
    s_rateLimiter.rateLimitValue(value);

    // Move the block time forward by 10 so the bucket refills by 10 * rate
    vm.warp(BLOCK_TIME + 1 + waitTime);

    // The bucket has filled up enough so we can take out more tokens
    s_rateLimiter.rateLimitValue(value);
    bucket = s_rateLimiter.currentRateLimiterState();
    assertEq(bucket.capacity - value + waitTime * s_config.rate - value, bucket.tokens);
    vm.resumeGasMetering();
  }

  // Reverts

  function test_AggregateValueMaxCapacityExceeded_Revert() public {
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState();

    uint256 numberOfTokens = 100;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector, bucket.capacity, (numberOfTokens * TOKEN_PRICE) / 1e18
      )
    );
    s_rateLimiter.rateLimitValue(value);
  }
}

contract AggregateTokenLimiter_getTokenValue is AggregateTokenLimiterSetup {
  function test_GetTokenValue_Success() public view {
    uint256 numberOfTokens = 10;
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: TOKEN, amount: 10});
    uint256 value = s_rateLimiter.getTokenValue(tokenAmount, IPriceRegistry(address(s_feeQuoter)));
    assertEq(value, (numberOfTokens * TOKEN_PRICE) / 1e18);
  }

  // Reverts
  function test_NoTokenPrice_Reverts() public {
    address tokenWithNoPrice = makeAddr("Token with no price");
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: tokenWithNoPrice, amount: 10});

    vm.expectRevert(abi.encodeWithSelector(AggregateRateLimiter.PriceNotFoundForToken.selector, tokenWithNoPrice));
    s_rateLimiter.getTokenValue(tokenAmount, IPriceRegistry(address(s_feeQuoter)));
  }
}
