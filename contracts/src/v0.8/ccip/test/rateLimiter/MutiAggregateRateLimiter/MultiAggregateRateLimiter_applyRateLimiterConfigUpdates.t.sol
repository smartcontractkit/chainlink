// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";

import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract MultiAggregateRateLimiter_applyRateLimiterConfigUpdates is MultiAggregateRateLimiterSetup {
  function test_ZeroConfigs() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](0);

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);
  }

  function test_SingleConfig() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig1
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

  function test_SingleConfigOutbound() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: true,
      rateLimiterConfig: s_rateLimiterConfig2
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

  function test_MultipleConfigs() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](5);

    for (uint64 i; i < configUpdates.length; ++i) {
      configUpdates[i] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
        remoteChainSelector: CHAIN_SELECTOR_1 + i + 1,
        isOutboundLane: i % 2 == 0 ? false : true,
        rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 5 + i, capacity: 100 + i})
      });

      vm.expectEmit();
      emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
        configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane, configUpdates[i].rateLimiterConfig
      );
    }

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, configUpdates.length);

    for (uint256 i; i < configUpdates.length; ++i) {
      RateLimiter.TokenBucket memory bucket =
        s_rateLimiter.currentRateLimiterState(configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane);
      _assertConfigWithTokenBucketEquality(configUpdates[i].rateLimiterConfig, bucket);
      assertEq(BLOCK_TIME, bucket.lastUpdated);
    }
  }

  function test_MultipleConfigsBothLanes() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](2);

    for (uint64 i; i < configUpdates.length; ++i) {
      configUpdates[i] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
        remoteChainSelector: CHAIN_SELECTOR_1 + 1,
        isOutboundLane: i % 2 == 0 ? false : true,
        rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 5 + i, capacity: 100 + i})
      });

      vm.expectEmit();
      emit MultiAggregateRateLimiter.RateLimiterConfigUpdated(
        configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane, configUpdates[i].rateLimiterConfig
      );
    }

    vm.recordLogs();
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);

    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, configUpdates.length);

    for (uint256 i; i < configUpdates.length; ++i) {
      RateLimiter.TokenBucket memory bucket =
        s_rateLimiter.currentRateLimiterState(configUpdates[i].remoteChainSelector, configUpdates[i].isOutboundLane);
      _assertConfigWithTokenBucketEquality(configUpdates[i].rateLimiterConfig, bucket);
      assertEq(BLOCK_TIME, bucket.lastUpdated);
    }
  }

  function test_UpdateExistingConfig() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig2
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

    // Outbound lane config remains unchanged
    _assertConfigWithTokenBucketEquality(
      s_rateLimiterConfig1, s_rateLimiter.currentRateLimiterState(CHAIN_SELECTOR_1, true)
    );
  }

  function test_UpdateExistingConfigWithNoDifference() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig1
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
  function test_RevertWhen_ZeroChainSelector() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: 0,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig1
    });

    vm.expectRevert(MultiAggregateRateLimiter.ZeroChainSelectorNotAllowed.selector);
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_RevertWhen_OnlyCallableByOwner() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: s_rateLimiterConfig1
    });
    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_RevertWhen_ConfigRateMoreThanCapacity() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 100, capacity: 100})
    });

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, configUpdates[0].rateLimiterConfig)
    );
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_RevertWhen_ConfigRateZero() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: RateLimiter.Config({isEnabled: true, rate: 0, capacity: 100})
    });

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, configUpdates[0].rateLimiterConfig)
    );
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_RevertWhen_DisableConfigRateNonZero() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: RateLimiter.Config({isEnabled: false, rate: 5, capacity: 100})
    });

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.DisabledNonZeroRateLimit.selector, configUpdates[0].rateLimiterConfig)
    );
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }

  function test_RevertWhen_DiableConfigCapacityNonZero() public {
    MultiAggregateRateLimiter.RateLimiterConfigArgs[] memory configUpdates =
      new MultiAggregateRateLimiter.RateLimiterConfigArgs[](1);
    configUpdates[0] = MultiAggregateRateLimiter.RateLimiterConfigArgs({
      remoteChainSelector: CHAIN_SELECTOR_1 + 1,
      isOutboundLane: false,
      rateLimiterConfig: RateLimiter.Config({isEnabled: false, rate: 0, capacity: 100})
    });

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.DisabledNonZeroRateLimit.selector, configUpdates[0].rateLimiterConfig)
    );
    s_rateLimiter.applyRateLimiterConfigUpdates(configUpdates);
  }
}
