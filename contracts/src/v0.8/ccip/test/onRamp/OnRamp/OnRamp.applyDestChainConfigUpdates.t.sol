// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRouter} from "../../../interfaces/IRouter.sol";

import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

contract OnRamp_applyDestChainConfigUpdates is OnRampSetup {
  function test_ApplyDestChainConfigUpdates() external {
    OnRamp.DestChainConfigArgs[] memory configArgs = new OnRamp.DestChainConfigArgs[](1);
    configArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;

    // supports disabling a lane by setting a router to zero
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, 0, IRouter(address(0)), false);

    s_onRamp.applyDestChainConfigUpdates(configArgs);

    (,, address router) = s_onRamp.getDestChainConfig(DEST_CHAIN_SELECTOR);
    assertEq(address(0), router);

    // supports updating and adding lanes simultaneously
    configArgs = new OnRamp.DestChainConfigArgs[](2);
    configArgs[0] = OnRamp.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      router: s_sourceRouter,
      allowlistEnabled: false
    });
    uint64 newDestChainSelector = 99999;
    address newRouter = makeAddr("newRouter");

    configArgs[1] = OnRamp.DestChainConfigArgs({
      destChainSelector: newDestChainSelector,
      router: IRouter(newRouter),
      allowlistEnabled: false
    });

    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, 0, s_sourceRouter, false);
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(newDestChainSelector, 0, IRouter(newRouter), false);

    s_onRamp.applyDestChainConfigUpdates(configArgs);

    (,, address newGotRouter) = s_onRamp.getDestChainConfig(newDestChainSelector);
    assertEq(newRouter, newGotRouter);

    // handles empty list
    uint256 numLogs = vm.getRecordedLogs().length;
    configArgs = new OnRamp.DestChainConfigArgs[](0);
    s_onRamp.applyDestChainConfigUpdates(configArgs);
    assertEq(numLogs, vm.getRecordedLogs().length); // indicates no changes made
  }

  function test_RevertWhen_ApplyDestChainConfigUpdates_WithInvalidChainSelector() external {
    OnRamp.DestChainConfigArgs[] memory configArgs = new OnRamp.DestChainConfigArgs[](1);
    configArgs[0].destChainSelector = 0; // invalid
    vm.expectRevert(abi.encodeWithSelector(OnRamp.InvalidDestChainConfig.selector, 0));
    s_onRamp.applyDestChainConfigUpdates(configArgs);
  }
}

contract OnRamp_applyAllowlistUpdates is OnRampSetup {
  function test_applyAllowlistUpdates() public {
    OnRamp.DestChainConfigArgs[] memory configArgs = new OnRamp.DestChainConfigArgs[](2);
    configArgs[0] = OnRamp.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      router: s_sourceRouter,
      allowlistEnabled: false
    });
    configArgs[1] =
      OnRamp.DestChainConfigArgs({destChainSelector: 9999, router: IRouter(address(9999)), allowlistEnabled: false});
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, 0, s_sourceRouter, false);
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(9999, 0, IRouter(address(9999)), false);
    s_onRamp.applyDestChainConfigUpdates(configArgs);

    (uint64 sequenceNumber, bool allowlistEnabled, address router) = s_onRamp.getDestChainConfig(9999);
    assertEq(sequenceNumber, 0);
    assertEq(allowlistEnabled, false);
    assertEq(router, address(9999));

    uint64[] memory destinationChainSelectors = new uint64[](2);
    destinationChainSelectors[0] = DEST_CHAIN_SELECTOR;
    destinationChainSelectors[1] = uint64(99999);

    address[] memory addedAllowlistedSenders = new address[](4);
    addedAllowlistedSenders[0] = vm.addr(1);
    addedAllowlistedSenders[1] = vm.addr(2);
    addedAllowlistedSenders[2] = vm.addr(3);
    addedAllowlistedSenders[3] = vm.addr(4);

    vm.expectEmit();
    emit OnRamp.AllowListSendersAdded(DEST_CHAIN_SELECTOR, addedAllowlistedSenders);

    OnRamp.AllowlistConfigArgs memory allowlistConfigArgs = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: true,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: addedAllowlistedSenders,
      removedAllowlistedSenders: new address[](0)
    });

    OnRamp.AllowlistConfigArgs[] memory applyAllowlistConfigArgsItems = new OnRamp.AllowlistConfigArgs[](1);
    applyAllowlistConfigArgsItems[0] = allowlistConfigArgs;

    s_onRamp.applyAllowlistUpdates(applyAllowlistConfigArgsItems);

    (bool isActive, address[] memory gotAllowList) = s_onRamp.getAllowedSendersList(DEST_CHAIN_SELECTOR);
    assertEq(4, gotAllowList.length);
    assertEq(addedAllowlistedSenders, gotAllowList);
    assertEq(true, isActive);

    address[] memory removedAllowlistedSenders = new address[](1);
    removedAllowlistedSenders[0] = vm.addr(2);

    vm.expectEmit();
    emit OnRamp.AllowListSendersRemoved(DEST_CHAIN_SELECTOR, removedAllowlistedSenders);

    allowlistConfigArgs = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: false,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: new address[](0),
      removedAllowlistedSenders: removedAllowlistedSenders
    });

    OnRamp.AllowlistConfigArgs[] memory allowlistConfigArgsItems_2 = new OnRamp.AllowlistConfigArgs[](1);
    allowlistConfigArgsItems_2[0] = allowlistConfigArgs;

    s_onRamp.applyAllowlistUpdates(allowlistConfigArgsItems_2);
    (isActive, gotAllowList) = s_onRamp.getAllowedSendersList(DEST_CHAIN_SELECTOR);
    assertEq(3, gotAllowList.length);
    assertFalse(isActive);

    addedAllowlistedSenders = new address[](2);
    addedAllowlistedSenders[0] = vm.addr(5);
    addedAllowlistedSenders[1] = vm.addr(6);

    removedAllowlistedSenders = new address[](2);
    removedAllowlistedSenders[0] = vm.addr(1);
    removedAllowlistedSenders[1] = vm.addr(3);

    vm.expectEmit();
    emit OnRamp.AllowListSendersAdded(DEST_CHAIN_SELECTOR, addedAllowlistedSenders);
    emit OnRamp.AllowListSendersRemoved(DEST_CHAIN_SELECTOR, removedAllowlistedSenders);

    allowlistConfigArgs = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: true,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: addedAllowlistedSenders,
      removedAllowlistedSenders: removedAllowlistedSenders
    });

    OnRamp.AllowlistConfigArgs[] memory allowlistConfigArgsItems_3 = new OnRamp.AllowlistConfigArgs[](1);
    allowlistConfigArgsItems_3[0] = allowlistConfigArgs;

    s_onRamp.applyAllowlistUpdates(allowlistConfigArgsItems_3);
    (isActive, gotAllowList) = s_onRamp.getAllowedSendersList(DEST_CHAIN_SELECTOR);

    assertEq(3, gotAllowList.length);
    assertTrue(isActive);
  }

  function test_RevertWhen_applyAllowlistUpdates() public {
    OnRamp.DestChainConfigArgs[] memory configArgs = new OnRamp.DestChainConfigArgs[](2);
    configArgs[0] = OnRamp.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      router: s_sourceRouter,
      allowlistEnabled: false
    });
    configArgs[1] =
      OnRamp.DestChainConfigArgs({destChainSelector: 9999, router: IRouter(address(9999)), allowlistEnabled: false});
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(DEST_CHAIN_SELECTOR, 0, s_sourceRouter, false);
    vm.expectEmit();
    emit OnRamp.DestChainConfigSet(9999, 0, IRouter(address(9999)), false);
    s_onRamp.applyDestChainConfigUpdates(configArgs);

    uint64[] memory destinationChainSelectors = new uint64[](2);
    destinationChainSelectors[0] = DEST_CHAIN_SELECTOR;
    destinationChainSelectors[1] = uint64(99999);

    address[] memory addedAllowlistedSenders = new address[](4);
    addedAllowlistedSenders[0] = vm.addr(1);
    addedAllowlistedSenders[1] = vm.addr(2);
    addedAllowlistedSenders[2] = vm.addr(3);
    addedAllowlistedSenders[3] = vm.addr(4);

    OnRamp.AllowlistConfigArgs memory allowlistConfigArgs = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: true,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: addedAllowlistedSenders,
      removedAllowlistedSenders: new address[](0)
    });

    OnRamp.AllowlistConfigArgs[] memory applyAllowlistConfigArgsItems = new OnRamp.AllowlistConfigArgs[](1);
    applyAllowlistConfigArgsItems[0] = allowlistConfigArgs;

    vm.startPrank(STRANGER);
    vm.expectRevert(OnRamp.OnlyCallableByOwnerOrAllowlistAdmin.selector);
    s_onRamp.applyAllowlistUpdates(applyAllowlistConfigArgsItems);
    vm.stopPrank();

    applyAllowlistConfigArgsItems[0].addedAllowlistedSenders[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(OnRamp.InvalidAllowListRequest.selector, DEST_CHAIN_SELECTOR));
    vm.startPrank(OWNER);
    s_onRamp.applyAllowlistUpdates(applyAllowlistConfigArgsItems);
    vm.stopPrank();
  }

  function test_applyAllowlistUpdates_InvalidAllowListRequestDisabledAllowListWithAdds() public {
    address[] memory addedAllowlistedSenders = new address[](1);
    addedAllowlistedSenders[0] = vm.addr(1);

    OnRamp.AllowlistConfigArgs memory allowlistConfigArgs = OnRamp.AllowlistConfigArgs({
      allowlistEnabled: false,
      destChainSelector: DEST_CHAIN_SELECTOR,
      addedAllowlistedSenders: addedAllowlistedSenders,
      removedAllowlistedSenders: new address[](0)
    });
    OnRamp.AllowlistConfigArgs[] memory applyAllowlistConfigArgsItems = new OnRamp.AllowlistConfigArgs[](1);
    applyAllowlistConfigArgsItems[0] = allowlistConfigArgs;

    vm.expectRevert(abi.encodeWithSelector(OnRamp.InvalidAllowListRequest.selector, DEST_CHAIN_SELECTOR));
    s_onRamp.applyAllowlistUpdates(applyAllowlistConfigArgsItems);
  }
}
