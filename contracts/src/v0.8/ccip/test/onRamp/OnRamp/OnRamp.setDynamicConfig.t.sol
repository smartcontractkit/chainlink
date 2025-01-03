// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

contract OnRamp_setDynamicConfig is OnRampSetup {
  function test_setDynamicConfig() public {
    OnRamp.StaticConfig memory staticConfig = s_onRamp.getStaticConfig();
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(23423),
      reentrancyGuardEntered: false,
      messageInterceptor: makeAddr("messageInterceptor"),
      feeAggregator: FEE_AGGREGATOR,
      allowlistAdmin: address(0)
    });

    vm.expectEmit();
    emit OnRamp.ConfigSet(staticConfig, newConfig);

    s_onRamp.setDynamicConfig(newConfig);

    OnRamp.DynamicConfig memory gotDynamicConfig = s_onRamp.getDynamicConfig();
    assertEq(newConfig.feeQuoter, gotDynamicConfig.feeQuoter);
  }

  // Reverts

  function test_RevertWhen_setDynamicConfig_InvalidConfigFeeQuoterEqAddressZero() public {
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(0),
      reentrancyGuardEntered: false,
      feeAggregator: FEE_AGGREGATOR,
      messageInterceptor: makeAddr("messageInterceptor"),
      allowlistAdmin: address(0)
    });

    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_RevertWhen_setDynamicConfig_InvalidConfigInvalidConfig() public {
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(23423),
      reentrancyGuardEntered: false,
      messageInterceptor: address(0),
      feeAggregator: FEE_AGGREGATOR,
      allowlistAdmin: address(0)
    });

    // Invalid price reg reverts.
    newConfig.feeQuoter = address(0);
    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_RevertWhen_setDynamicConfig_InvalidConfigFeeAggregatorEqAddressZero() public {
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(23423),
      reentrancyGuardEntered: false,
      messageInterceptor: address(0),
      feeAggregator: address(0),
      allowlistAdmin: address(0)
    });

    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }

  function test_RevertWhen_setDynamicConfig_InvalidConfigOnlyOwner() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_onRamp.setDynamicConfig(_generateDynamicOnRampConfig(address(2)));
  }

  function test_RevertWhen_setDynamicConfig_InvalidConfigReentrancyGuardEnteredEqTrue() public {
    OnRamp.DynamicConfig memory newConfig = OnRamp.DynamicConfig({
      feeQuoter: address(23423),
      reentrancyGuardEntered: true,
      messageInterceptor: makeAddr("messageInterceptor"),
      feeAggregator: FEE_AGGREGATOR,
      allowlistAdmin: address(0)
    });

    vm.expectRevert(OnRamp.InvalidConfig.selector);
    s_onRamp.setDynamicConfig(newConfig);
  }
}
