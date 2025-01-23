// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IBridgeAdapter} from "../interfaces/IBridge.sol";
import {ILiquidityManager} from "../interfaces/ILiquidityManager.sol";

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";
import {LiquidityManager} from "../LiquidityManager.sol";

contract LiquidityManager_setCrossChainRebalancer is LiquidityManagerSetup {
  function test_setCrossChainRebalancer() external {
    address newRebalancer = address(23892423);
    uint64 remoteChainSelector = 12301293;

    uint64[] memory supportedChains = s_liquidityManager.getSupportedDestChains();
    assertEq(supportedChains.length, 0);

    LiquidityManager.CrossChainRebalancerArgs[] memory args = new LiquidityManager.CrossChainRebalancerArgs[](1);
    args[0] = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: newRebalancer,
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: remoteChainSelector,
      enabled: true
    });

    vm.expectEmit();
    emit LiquidityManager.CrossChainRebalancerSet(
      remoteChainSelector,
      args[0].localBridge,
      args[0].remoteToken,
      newRebalancer,
      args[0].enabled
    );

    s_liquidityManager.setCrossChainRebalancers(args);

    assertEq(s_liquidityManager.getCrossChainRebalancer(remoteChainSelector).remoteRebalancer, newRebalancer);

    LiquidityManager.CrossChainRebalancerArgs[] memory crossChainRebalancers = s_liquidityManager
      .getAllCrossChainRebalancers();
    assertEq(crossChainRebalancers.length, 1);
    assertEq(crossChainRebalancers[0].remoteRebalancer, args[0].remoteRebalancer);
    assertEq(address(crossChainRebalancers[0].localBridge), address(args[0].localBridge));
    assertEq(crossChainRebalancers[0].remoteToken, args[0].remoteToken);
    assertEq(crossChainRebalancers[0].remoteChainSelector, args[0].remoteChainSelector);
    assertEq(crossChainRebalancers[0].enabled, args[0].enabled);

    supportedChains = s_liquidityManager.getSupportedDestChains();
    assertEq(supportedChains.length, 1);
    assertEq(supportedChains[0], remoteChainSelector);

    address anotherRebalancer = address(123);
    args[0].remoteRebalancer = anotherRebalancer;

    vm.expectEmit();
    emit LiquidityManager.CrossChainRebalancerSet(
      remoteChainSelector,
      args[0].localBridge,
      args[0].remoteToken,
      anotherRebalancer,
      args[0].enabled
    );

    s_liquidityManager.setCrossChainRebalancer(args[0]);

    assertEq(s_liquidityManager.getCrossChainRebalancer(remoteChainSelector).remoteRebalancer, anotherRebalancer);

    supportedChains = s_liquidityManager.getSupportedDestChains();
    assertEq(supportedChains.length, 1);
    assertEq(supportedChains[0], remoteChainSelector);
  }

  function test_setCrossChainRebalancer_RevertWhen_ZeroChainSelector() external {
    LiquidityManager.CrossChainRebalancerArgs memory arg = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(9),
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: 0,
      enabled: true
    });

    vm.expectRevert(LiquidityManager.ZeroChainSelector.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);
  }

  function test_setCrossChainRebalancer_RevertWhen_ZeroAddressRemoteRebalancer() external {
    LiquidityManager.CrossChainRebalancerArgs memory arg = ILiquidityManager.CrossChainRebalancerArgs({
      remoteRebalancer: address(0),
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: 123,
      enabled: true
    });

    vm.expectRevert(LiquidityManager.ZeroAddress.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);

    arg.remoteRebalancer = address(9);
    arg.localBridge = IBridgeAdapter(address(0));

    vm.expectRevert(LiquidityManager.ZeroAddress.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);

    arg.localBridge = s_bridgeAdapter;
    arg.remoteToken = address(0);

    vm.expectRevert(LiquidityManager.ZeroAddress.selector);

    s_liquidityManager.setCrossChainRebalancer(arg);
  }

  function test_setCrossChainRebalancer_RevertWhen_CallerNotOwner() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    // Test the entrypoint that takes a list
    s_liquidityManager.setCrossChainRebalancers(new LiquidityManager.CrossChainRebalancerArgs[](0));

    vm.expectRevert("Only callable by owner");

    // Test the entrypoint that takes a single item
    s_liquidityManager.setCrossChainRebalancer(
      ILiquidityManager.CrossChainRebalancerArgs({
        remoteRebalancer: address(9),
        localBridge: s_bridgeAdapter,
        remoteToken: address(190490124908),
        remoteChainSelector: 124,
        enabled: true
      })
    );
  }
}
