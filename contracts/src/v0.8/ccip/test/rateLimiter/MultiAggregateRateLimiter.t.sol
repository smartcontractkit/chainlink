// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Vm} from "forge-std/Vm.sol";

import {MultiAggregateRateLimiter} from "../../MultiAggregateRateLimiter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {MultiAggregateRateLimiterHelper} from "../helpers/MultiAggregateRateLimiterHelper.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";

import {stdError} from "forge-std/Test.sol";

contract MultiAggregateRateLimiterSetup is BaseTest, PriceRegistrySetup {
  MultiAggregateRateLimiterHelper internal s_rateLimiter;

  address internal immutable TOKEN = 0x21118E64E1fB0c487F25Dd6d3601FF6af8D32E4e;
  uint224 internal constant TOKEN_PRICE = 4e18;

  uint64 internal constant CHAIN_SELECTOR_1 = 5009297550715157269;
  uint64 internal constant CHAIN_SELECTOR_2 = 4949039107694359620;

  RateLimiter.Config internal RATE_LIMITER_CONFIG_1 = RateLimiter.Config({isEnabled: true, rate: 5, capacity: 100});
  RateLimiter.Config internal RATE_LIMITER_CONFIG_2 = RateLimiter.Config({isEnabled: true, rate: 10, capacity: 200});

  address internal immutable MOCK_OFFRAMP = address(1111);
  address internal immutable MOCK_ONRAMP = address(1112);

  address[] internal s_authorizedCallers;

  function setUp() public virtual override(BaseTest, PriceRegistrySetup) {
    BaseTest.setUp();
    PriceRegistrySetup.setUp();

    Internal.PriceUpdates memory priceUpdates = getSingleTokenPriceUpdateStruct(TOKEN, TOKEN_PRICE);
    s_priceRegistry.updatePrices(priceUpdates);

    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](3);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    configUpdates[1] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_2,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });
    configUpdates[2] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });

    s_authorizedCallers = new address[](2);
    s_authorizedCallers[0] = MOCK_OFFRAMP;
    s_authorizedCallers[1] = MOCK_ONRAMP;

    s_rateLimiter = new MultiAggregateRateLimiterHelper(configUpdates, address(s_priceRegistry), s_authorizedCallers);
  }

  function _assertConfigWithTokenBucketEquality(
    RateLimiter.Config memory config,
    RateLimiter.TokenBucket memory tokenBucket
  ) internal pure {
    assertEq(config.rate, tokenBucket.rate);
    assertEq(config.capacity, tokenBucket.capacity);
    assertEq(config.capacity, tokenBucket.tokens);
    assertEq(config.isEnabled, tokenBucket.isEnabled);
  }

  function _assertTokenBucketEquality(
    RateLimiter.TokenBucket memory tokenBucketA,
    RateLimiter.TokenBucket memory tokenBucketB
  ) internal pure {
    assertEq(tokenBucketA.rate, tokenBucketB.rate);
    assertEq(tokenBucketA.capacity, tokenBucketB.capacity);
    assertEq(tokenBucketA.tokens, tokenBucketB.tokens);
    assertEq(tokenBucketA.isEnabled, tokenBucketB.isEnabled);
  }
}

contract MultiAggregateRateLimiter_constructor is MultiAggregateRateLimiterSetup {
  function test_ConstructorNoAuthorizedCallers_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](0);
    address[] memory authorizedCallers = new address[](0);

    vm.recordLogs();
    s_rateLimiter = new MultiAggregateRateLimiterHelper(configUpdates, address(s_priceRegistry), authorizedCallers);

    // PriceRegistrySet
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_priceRegistry), s_rateLimiter.getPriceRegistry());
  }

  function test_ConstructorNoConfigs_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](0);
    address[] memory authorizedCallers = new address[](2);
    authorizedCallers[0] = MOCK_OFFRAMP;
    authorizedCallers[1] = MOCK_ONRAMP;

    vm.recordLogs();
    s_rateLimiter = new MultiAggregateRateLimiterHelper(configUpdates, address(s_priceRegistry), authorizedCallers);

    // PriceRegistrySet + 2 authorized caller sets
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 3);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_priceRegistry), s_rateLimiter.getPriceRegistry());
  }

  function test_Constructor_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](3);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    configUpdates[1] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_2,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });
    configUpdates[2] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });

    address[] memory authorizedCallers = new address[](2);
    authorizedCallers[0] = MOCK_OFFRAMP;
    authorizedCallers[1] = MOCK_ONRAMP;

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(CHAIN_SELECTOR_1, false, RATE_LIMITER_CONFIG_1);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(CHAIN_SELECTOR_2, false, RATE_LIMITER_CONFIG_2);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(CHAIN_SELECTOR_1, true, RATE_LIMITER_CONFIG_2);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.PriceRegistrySet(address(s_priceRegistry));

    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(MOCK_OFFRAMP);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(MOCK_ONRAMP);

    s_rateLimiter = new MultiAggregateRateLimiterHelper(configUpdates, address(s_priceRegistry), authorizedCallers);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_priceRegistry), s_rateLimiter.getPriceRegistry());

    RateLimiter.TokenBucket memory bucketSrcChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_1, bucketSrcChain1);
    assertEq(BLOCK_TIME, bucketSrcChain1.lastUpdated);

    RateLimiter.TokenBucket memory bucketSrcChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_2, bucketSrcChain2);
    assertEq(BLOCK_TIME, bucketSrcChain2.lastUpdated);

    RateLimiter.TokenBucket memory bucketSrcChainOutgoing =
      s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_2, bucketSrcChainOutgoing);
    assertEq(BLOCK_TIME, bucketSrcChainOutgoing.lastUpdated);
  }
}

contract MultiAggregateRateLimiter_setPriceRegistry is MultiAggregateRateLimiterSetup {
  function test_Owner_Success() public {
    address newAddress = address(42);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.PriceRegistrySet(newAddress);

    s_rateLimiter.setPriceRegistry(newAddress);
    assertEq(newAddress, s_rateLimiter.getPriceRegistry());
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(bytes("Only callable by owner"));

    s_rateLimiter.setPriceRegistry(STRANGER);
  }

  function test_ZeroAddress_Revert() public {
    vm.expectRevert(MultiAggregateRateLimiter.ZeroAddressNotAllowed.selector);
    s_rateLimiter.setPriceRegistry(address(0));
  }
}

contract MultiAggregateRateLimiter_setAuthorizedCallers is MultiAggregateRateLimiterSetup {
  function test_OnlyAdd_Success() public {
    address[] memory addedCallers = new address[](2);
    addedCallers[0] = address(42);
    addedCallers[1] = address(43);

    address[] memory removedCallers = new address[](0);

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), s_authorizedCallers);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(addedCallers[0]);
    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(addedCallers[1]);

    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    address[] memory expectedCallers = new address[](4);
    expectedCallers[0] = s_authorizedCallers[0];
    expectedCallers[1] = s_authorizedCallers[1];
    expectedCallers[2] = addedCallers[0];
    expectedCallers[3] = addedCallers[1];

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_OnlyRemove_Success() public {
    address[] memory addedCallers = new address[](0);

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = s_authorizedCallers[0];

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), s_authorizedCallers);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerRemoved(removedCallers[0]);

    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    address[] memory expectedCallers = new address[](1);
    expectedCallers[0] = s_authorizedCallers[1];

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_AddAndRemove_Success() public {
    address[] memory addedCallers = new address[](2);
    addedCallers[0] = address(42);
    addedCallers[1] = address(43);

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = s_authorizedCallers[0];

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), s_authorizedCallers);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(addedCallers[0]);
    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(addedCallers[1]);
    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerRemoved(removedCallers[0]);

    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    // Order of the set changes on removal
    address[] memory expectedCallers = new address[](3);
    expectedCallers[0] = addedCallers[1];
    expectedCallers[1] = s_authorizedCallers[1];
    expectedCallers[2] = addedCallers[0];

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_AddThenRemove_Success() public {
    address[] memory addedCallers = new address[](1);
    addedCallers[0] = address(42);

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = address(42);

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), s_authorizedCallers);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerAdded(addedCallers[0]);
    vm.expectEmit();
    emit MultiAggregateRateLimiter.AuthorizedCallerRemoved(addedCallers[0]);

    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), s_authorizedCallers);
  }

  function test_SkipRemove_Success() public {
    address[] memory addedCallers = new address[](0);

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = address(42);

    vm.recordLogs();
    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    assertEq(s_rateLimiter.getAllAuthorizedCallers(), s_authorizedCallers);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(bytes("Only callable by owner"));

    address[] memory addedCallers = new address[](2);
    addedCallers[0] = address(42);
    addedCallers[1] = address(43);

    address[] memory removedCallers = new address[](0);

    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );
  }

  function test_ZeroAddressAdd_Revert() public {
    address[] memory addedCallers = new address[](1);
    addedCallers[0] = address(0);
    address[] memory removedCallers = new address[](0);

    vm.expectRevert(MultiAggregateRateLimiter.ZeroAddressNotAllowed.selector);
    s_rateLimiter.applyAuthorizedCallerUpdates(
      MultiAggregateRateLimiter.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );
  }
}

contract MultiAggregateRateLimiter_getTokenBucket is MultiAggregateRateLimiterSetup {
  function test_GetTokenBucket_Success() public view {
    RateLimiter.TokenBucket memory bucketIncoming = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_1, bucketIncoming);
    assertEq(BLOCK_TIME, bucketIncoming.lastUpdated);

    RateLimiter.TokenBucket memory bucketOutgoing = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);
    _assertConfigWithTokenBucketEquality(RATE_LIMITER_CONFIG_2, bucketOutgoing);
    assertEq(BLOCK_TIME, bucketOutgoing.lastUpdated);
  }

  function test_Refill_Success() public {
    RATE_LIMITER_CONFIG_1.capacity = RATE_LIMITER_CONFIG_1.capacity * 2;

    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    assertEq(RATE_LIMITER_CONFIG_1.rate, bucket.rate);
    assertEq(RATE_LIMITER_CONFIG_1.capacity, bucket.capacity);
    assertEq(RATE_LIMITER_CONFIG_1.capacity / 2, bucket.tokens);
    assertEq(BLOCK_TIME, bucket.lastUpdated);

    uint256 warpTime = 4;
    vm.warp(BLOCK_TIME + warpTime);

    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    assertEq(RATE_LIMITER_CONFIG_1.rate, bucket.rate);
    assertEq(RATE_LIMITER_CONFIG_1.capacity, bucket.capacity);
    assertEq(RATE_LIMITER_CONFIG_1.capacity / 2 + warpTime * RATE_LIMITER_CONFIG_1.rate, bucket.tokens);
    assertEq(BLOCK_TIME + warpTime, bucket.lastUpdated);

    vm.warp(BLOCK_TIME + warpTime * 100);

    // Bucket overflow
    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(RATE_LIMITER_CONFIG_1.capacity, bucket.tokens);
  }

  // Reverts

  function test_TimeUnderflow_Revert() public {
    vm.warp(BLOCK_TIME - 1);

    vm.expectRevert(stdError.arithmeticError);
    s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
  }
}

contract MultiAggregateRateLimiter_applyRateLimiterConfigUpdates is MultiAggregateRateLimiterSetup {
  function test_ZeroConfigs_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](0);

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_SingleConfig_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, false, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    RateLimiter.TokenBucket memory bucket1 =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);
    _assertConfigWithTokenBucketEquality(configUpdates[0].rateLimiterConfig, bucket1);
    assertEq(BLOCK_TIME, bucket1.lastUpdated);
  }

  function test_SingleConfigOutgoing_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutgoingLane: true,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, true, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    RateLimiter.TokenBucket memory bucket1 =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, true);
    _assertConfigWithTokenBucketEquality(configUpdates[0].rateLimiterConfig, bucket1);
    assertEq(BLOCK_TIME, bucket1.lastUpdated);
  }

  function test_MultipleConfigs_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](5);

    for (uint64 i; i < configUpdates.length; ++i) {
      configUpdates[i] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
        remoteChainSelector: CHAIN_SELECTOR_1 + i + 1,
        isOutgoingLane: i % 2 == 0 ? false : true,
        rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 5 + i, capacity: 100 + i})
      });

      vm.expectEmit();
      emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
        configUpdates[i].remoteChainSelector, configUpdates[i].isOutgoingLane, configUpdates[i].rateLimiterConfig
      );
    }

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, configUpdates.length);

    for (uint256 i; i < configUpdates.length; ++i) {
      RateLimiter.TokenBucket memory bucket =
        s_rateLimiter.currentRateLimiterState(configUpdates[i].remoteChainSelector, configUpdates[i].isOutgoingLane);
      _assertConfigWithTokenBucketEquality(configUpdates[i].rateLimiterConfig, bucket);
      assertEq(BLOCK_TIME, bucket.lastUpdated);
    }
  }

  function test_MultipleConfigsBothLanes_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](2);

    for (uint64 i; i < configUpdates.length; ++i) {
      configUpdates[i] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
        remoteChainSelector: CHAIN_SELECTOR_1 + 1,
        isOutgoingLane: i % 2 == 0 ? false : true,
        rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 5 + i, capacity: 100 + i})
      });

      vm.expectEmit();
      emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
        configUpdates[i].remoteChainSelector, configUpdates[i].isOutgoingLane, configUpdates[i].rateLimiterConfig
      );
    }

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, configUpdates.length);

    for (uint256 i; i < configUpdates.length; ++i) {
      RateLimiter.TokenBucket memory bucket =
        s_rateLimiter.currentRateLimiterState(configUpdates[i].remoteChainSelector, configUpdates[i].isOutgoingLane);
      _assertConfigWithTokenBucketEquality(configUpdates[i].rateLimiterConfig, bucket);
      assertEq(BLOCK_TIME, bucket.lastUpdated);
    }
  }

  function test_UpdateExistingConfig_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_2
    });

    RateLimiter.TokenBucket memory bucket1 =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);

    // Capacity equals tokens
    assertEq(bucket1.capacity, bucket1.tokens);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, false, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    vm.warp(BLOCK_TIME + 1);
    bucket1 = s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);
    assertEq(BLOCK_TIME + 1, bucket1.lastUpdated);

    // Tokens < capacity since capacity doubled
    assertTrue(bucket1.capacity != bucket1.tokens);

    // Outgoing lane config remains unchanged
    _assertConfigWithTokenBucketEquality(
      RATE_LIMITER_CONFIG_2, s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true)
    );
  }

  function test_UpdateExistingConfigWithNoDifference_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    RateLimiter.TokenBucket memory bucketPreUpdate =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
      configUpdates[0].remoteChainSelector, false, configUpdates[0].rateLimiterConfig
    );

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    vm.warp(BLOCK_TIME + 1);
    RateLimiter.TokenBucket memory bucketPostUpdate =
      s_rateLimiter.currentRateLimiterState(configUpdates[0].remoteChainSelector, false);
    _assertTokenBucketEquality(bucketPreUpdate, bucketPostUpdate);
    assertEq(BLOCK_TIME + 1, bucketPostUpdate.lastUpdated);
  }

  // Reverts
  function test_ZeroChainSelector_Revert() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: 0,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });

    vm.expectRevert(MultiAggregateRateLimiter.ZeroChainSelectorNotAllowed.selector);
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_OnlyCallableByOwner_Revert() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    vm.startPrank(STRANGER);

    vm.expectRevert(bytes("Only callable by owner"));
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }
}

contract MultiAggregateRateLimiter__rateLimitValue is MultiAggregateRateLimiterSetup {
  function test_RateLimitValue_Success_gas() public {
    vm.pauseGasMetering();
    // start from blocktime that does not equal rate limiter init timestamp
    vm.warp(BLOCK_TIME + 1);

    // 15 (tokens) * 4 (price) * 2 (number of times) > 100 (capacity)
    uint256 numberOfTokens = 15;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, false, value);
    vm.pauseGasMetering();

    // Get the updated bucket status
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    // Assert the proper value has been taken out of the bucket
    assertEq(bucket.capacity - value, bucket.tokens);

    // Since value * 2 > bucket.capacity we cannot take it out twice.
    // Expect a revert when we try, with a wait time.
    uint256 waitTime = 4;
    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, waitTime, bucket.tokens)
    );
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, false, value);

    // Move the block time forward by 10 so the bucket refills by 10 * rate
    vm.warp(BLOCK_TIME + 1 + waitTime);

    // The bucket has filled up enough so we can take out more tokens
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, false, value);
    bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucket.capacity - value + waitTime * RATE_LIMITER_CONFIG_1.rate - value, bucket.tokens);
    vm.resumeGasMetering();
  }

  function test_RateLimitValueDifferentChainSelectors_Success() public {
    vm.pauseGasMetering();
    // start from blocktime that does not equal rate limiter init timestamp
    vm.warp(BLOCK_TIME + 1);

    // 15 (tokens) * 4 (price) * 2 (number of times) > 100 (capacity)
    uint256 numberOfTokens = 15;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, false, value);
    vm.pauseGasMetering();

    // Get the updated bucket status
    RateLimiter.TokenBucket memory bucket1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    RateLimiter.TokenBucket memory bucket2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);

    // Assert the proper value has been taken out of the bucket
    assertEq(bucket1.capacity - value, bucket1.tokens);
    // CHAIN_SELECTOR_2 should remain unchanged
    assertEq(bucket2.capacity, bucket2.tokens);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_2, false, value);
    vm.pauseGasMetering();

    bucket1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    bucket2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);

    assertEq(bucket2.capacity - value, bucket2.tokens);
    // CHAIN_SELECTOR_1 should remain unchanged
    assertEq(bucket1.capacity - value, bucket1.tokens);
  }

  function test_RateLimitValueDifferentLanes_Success() public {
    vm.pauseGasMetering();
    // start from blocktime that does not equal rate limiter init timestamp
    vm.warp(BLOCK_TIME + 1);

    // 15 (tokens) * 4 (price) * 2 (number of times) > 100 (capacity)
    uint256 numberOfTokens = 15;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, false, value);
    vm.pauseGasMetering();

    // Get the updated bucket status
    RateLimiter.TokenBucket memory bucket1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    RateLimiter.TokenBucket memory bucket2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);

    // Assert the proper value has been taken out of the bucket
    assertEq(bucket1.capacity - value, bucket1.tokens);
    // Outgoing lane should remain unchanged
    assertEq(bucket2.capacity, bucket2.tokens);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(value);

    vm.resumeGasMetering();
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, true, value);
    vm.pauseGasMetering();

    bucket1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    bucket2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true);

    assertEq(bucket2.capacity - value, bucket2.tokens);
    // Incoming lane should remain unchanged
    assertEq(bucket1.capacity - value, bucket1.tokens);
  }

  // Reverts

  function test_AggregateValueMaxCapacityExceeded_Revert() public {
    RateLimiter.TokenBucket memory bucket = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);

    uint256 numberOfTokens = 100;
    uint256 value = (numberOfTokens * TOKEN_PRICE) / 1e18;

    vm.expectRevert(
      abi.encodeWithSelector(
        RateLimiter.AggregateValueMaxCapacityExceeded.selector, bucket.capacity, (numberOfTokens * TOKEN_PRICE) / 1e18
      )
    );
    s_rateLimiter.rateLimitValue(CHAIN_SELECTOR_1, false, value);
  }
}

contract MultiAggregateRateLimiter__getTokenValue is MultiAggregateRateLimiterSetup {
  function test_GetTokenValue_Success() public view {
    uint256 numberOfTokens = 10;
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: TOKEN, amount: 10});
    uint256 value = s_rateLimiter.getTokenValue(tokenAmount);
    assertEq(value, (numberOfTokens * TOKEN_PRICE) / 1e18);
  }

  // Reverts
  function test_NoTokenPrice_Reverts() public {
    address tokenWithNoPrice = makeAddr("Token with no price");
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: tokenWithNoPrice, amount: 10});

    vm.expectRevert(abi.encodeWithSelector(MultiAggregateRateLimiter.PriceNotFoundForToken.selector, tokenWithNoPrice));
    s_rateLimiter.getTokenValue(tokenAmount);
  }
}

contract MultiAggregateRateLimiter_updateRateLimitTokens is MultiAggregateRateLimiterSetup {
  function setUp() public virtual override {
    super.setUp();

    // Clear rate limit tokens state
    MultiAggregateRateLimiter.RateLimitToken[] memory remove =
      new MultiAggregateRateLimiter.RateLimitToken[](s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      remove[i] =
        MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[i], localToken: s_destTokens[i]});
    }
    s_rateLimiter.updateRateLimitTokens(remove, new MultiAggregateRateLimiter.RateLimitToken[](0));
  }

  function test_UpdateRateLimitTokens_Success() public {
    MultiAggregateRateLimiter.RateLimitToken[] memory adds = new MultiAggregateRateLimiter.RateLimitToken[](2);
    adds[0] = MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[0], localToken: s_destTokens[0]});
    adds[1] = MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[1], localToken: s_destTokens[1]});

    for (uint256 i = 0; i < adds.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitAdded(adds[i].remoteToken, adds[i].localToken);
    }

    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.RateLimitToken[](0), adds);

    (address[] memory sourceTokens, address[] memory destTokens) = s_rateLimiter.getAllRateLimitTokens();

    for (uint256 i = 0; i < adds.length; ++i) {
      assertEq(adds[i].remoteToken, sourceTokens[i]);
      assertEq(adds[i].localToken, destTokens[i]);
    }
  }

  function test_UpdateRateLimitTokens_AddsAndRemoves_Success() public {
    MultiAggregateRateLimiter.RateLimitToken[] memory adds = new MultiAggregateRateLimiter.RateLimitToken[](2);
    adds[0] = MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[0], localToken: s_destTokens[0]});
    adds[1] = MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[1], localToken: s_destTokens[1]});

    MultiAggregateRateLimiter.RateLimitToken[] memory removes = new MultiAggregateRateLimiter.RateLimitToken[](1);
    removes[0] = adds[0];

    for (uint256 i = 0; i < adds.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitAdded(adds[i].remoteToken, adds[i].localToken);
    }

    s_rateLimiter.updateRateLimitTokens(removes, adds);

    for (uint256 i = 0; i < removes.length; ++i) {
      vm.expectEmit();
      emit MultiAggregateRateLimiter.TokenAggregateRateLimitRemoved(removes[i].remoteToken, removes[i].localToken);
    }

    s_rateLimiter.updateRateLimitTokens(removes, new MultiAggregateRateLimiter.RateLimitToken[](0));

    (address[] memory sourceTokens, address[] memory destTokens) = s_rateLimiter.getAllRateLimitTokens();

    assertEq(1, sourceTokens.length);
    assertEq(adds[1].remoteToken, sourceTokens[0]);

    assertEq(1, destTokens.length);
    assertEq(adds[1].localToken, destTokens[0]);
  }

  // Reverts

  function test_ZeroSourceToken_Revert() public {
    MultiAggregateRateLimiter.RateLimitToken[] memory adds = new MultiAggregateRateLimiter.RateLimitToken[](1);
    adds[0] = MultiAggregateRateLimiter.RateLimitToken({remoteToken: address(0), localToken: s_destTokens[0]});

    vm.expectRevert(MultiAggregateRateLimiter.ZeroAddressNotAllowed.selector);
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.RateLimitToken[](0), adds);
  }

  function test_ZeroDestToken_Revert() public {
    MultiAggregateRateLimiter.RateLimitToken[] memory adds = new MultiAggregateRateLimiter.RateLimitToken[](1);
    adds[0] = MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_destTokens[0], localToken: address(0)});

    vm.expectRevert(MultiAggregateRateLimiter.ZeroAddressNotAllowed.selector);
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.RateLimitToken[](0), adds);
  }

  function test_NonOwner_Revert() public {
    MultiAggregateRateLimiter.RateLimitToken[] memory addsAndRemoves = new MultiAggregateRateLimiter.RateLimitToken[](4);

    vm.startPrank(STRANGER);

    vm.expectRevert(bytes("Only callable by owner"));
    s_rateLimiter.updateRateLimitTokens(addsAndRemoves, addsAndRemoves);
  }
}

contract MultiAggregateRateLimiter_onIncomingMessage is MultiAggregateRateLimiterSetup {
  address internal immutable MOCK_RECEIVER = address(1113);

  function setUp() public virtual override {
    super.setUp();

    MultiAggregateRateLimiter.RateLimitToken[] memory tokensToAdd =
      new MultiAggregateRateLimiter.RateLimitToken[](s_sourceTokens.length);
    for (uint224 i = 0; i < s_sourceTokens.length; ++i) {
      tokensToAdd[i] =
        MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[i], localToken: s_destTokens[i]});

      Internal.PriceUpdates memory priceUpdates =
        getSingleTokenPriceUpdateStruct(s_destTokens[i], TOKEN_PRICE * (i + 1));
      s_priceRegistry.updatePrices(priceUpdates);
    }
    s_rateLimiter.updateRateLimitTokens(new MultiAggregateRateLimiter.RateLimitToken[](0), tokensToAdd);
  }

  function test_ValidateMessageWithNoTokens_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    vm.recordLogs();
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessageNoTokens(CHAIN_SELECTOR_1));

    // No consumed rate limit events
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_ValidateMessageWithTokens_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 3});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    // 3 tokens * TOKEN_PRICE + 1 token * (2 * TOKEN_PRICE)
    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  function test_ValidateMessageWithDisabledRateLimitToken_Success() public {
    MultiAggregateRateLimiter.RateLimitToken[] memory tokensToRemove = new MultiAggregateRateLimiter.RateLimitToken[](1);
    tokensToRemove[0] =
      MultiAggregateRateLimiter.RateLimitToken({remoteToken: s_sourceTokens[1], localToken: s_destTokens[1]});
    s_rateLimiter.updateRateLimitTokens(tokensToRemove, new MultiAggregateRateLimiter.RateLimitToken[](0));

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 5});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    vm.startPrank(MOCK_OFFRAMP);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed((5 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  function test_ValidateMessageWithRateLimitDisabled_Success() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutgoingLane: false,
      rateLimiterConfig: RATE_LIMITER_CONFIG_1
    });
    configUpdates[0].rateLimiterConfig.isEnabled = false;

    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 1000});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 50});

    vm.startPrank(MOCK_OFFRAMP);
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // No consumed rate limit events
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_ValidateMessageWithTokensOnDifferentChains_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 2});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 1});

    // 2 tokens * (TOKEN_PRICE) + 1 token * (2 * TOKEN_PRICE)
    uint256 totalValue = (4 * TOKEN_PRICE) / 1e18;

    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Chain 1 changed
    RateLimiter.TokenBucket memory bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 unchanged
    RateLimiter.TokenBucket memory bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    assertEq(bucketChain2.capacity, bucketChain2.tokens);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed((4 * TOKEN_PRICE) / 1e18);

    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_2, tokenAmounts));

    // Chain 1 unchanged
    bucketChain1 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, false);
    assertEq(bucketChain1.capacity - totalValue, bucketChain1.tokens);

    // Chain 2 changed
    bucketChain2 = s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_2, false);
    assertEq(bucketChain2.capacity - totalValue, bucketChain2.tokens);
  }

  function test_ValidateMessageWithRateLimitReset_Success() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 20});

    // Remaining capacity: 100 -> 20
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Cannot fit 80 rate limit value (need to wait at least 12 blocks, current capacity is 20)
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 12, 20));
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Remaining capacity: 20 -> 35 (need to wait 9 more blocks)
    vm.warp(BLOCK_TIME + 3);
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueRateLimitReached.selector, 9, 35));
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));

    // Remaining capacity: 35 -> 80 (can fit exactly 80)
    vm.warp(BLOCK_TIME + 12);
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  // Reverts

  function test_ValidateMessageWithRateLimitExceeded_Revert() public {
    vm.startPrank(MOCK_OFFRAMP);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: s_destTokens[0], amount: 80});
    tokenAmounts[1] = Client.EVMTokenAmount({token: s_destTokens[1], amount: 30});

    uint256 totalValue = (80 * TOKEN_PRICE + 2 * (30 * TOKEN_PRICE)) / 1e18;
    vm.expectRevert(abi.encodeWithSelector(RateLimiter.AggregateValueMaxCapacityExceeded.selector, 100, totalValue));
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessage(CHAIN_SELECTOR_1, tokenAmounts));
  }

  function test_ValidateMessageFromUnauthorizedCaller_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(MultiAggregateRateLimiter.UnauthorizedCaller.selector, STRANGER));
    s_rateLimiter.onIncomingMessage(_generateAny2EVMMessageNoTokens(CHAIN_SELECTOR_1));
  }

  function _generateAny2EVMMessageNoTokens(uint64 sourceChainSelector)
    internal
    pure
    returns (Client.Any2EVMMessage memory)
  {
    return _generateAny2EVMMessage(sourceChainSelector, new Client.EVMTokenAmount[](0));
  }

  function _generateAny2EVMMessage(
    uint64 sourceChainSelector,
    Client.EVMTokenAmount[] memory tokenAmounts
  ) internal pure returns (Client.Any2EVMMessage memory) {
    return Client.Any2EVMMessage({
      messageId: keccak256(bytes("messageId")),
      sourceChainSelector: sourceChainSelector,
      sender: abi.encode(OWNER),
      data: abi.encode(0),
      destTokenAmounts: tokenAmounts
    });
  }
}
