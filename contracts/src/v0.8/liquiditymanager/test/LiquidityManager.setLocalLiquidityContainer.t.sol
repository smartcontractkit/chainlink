// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";
import {MockLockReleaseTokenPool} from "./mocks/MockTokenPool.sol";
import {LiquidityManager} from "../LiquidityManager.sol";

contract LiquidityManager_setLocalLiquidityContainer is LiquidityManagerSetup {
  function test_setLocalLiquidityContainer() external {
    MockLockReleaseTokenPool newPool = new MockLockReleaseTokenPool(s_l1Token);

    vm.expectEmit();
    emit LiquidityManager.LiquidityContainerSet(address(newPool));

    s_liquidityManager.setLocalLiquidityContainer(newPool);

    assertEq(s_liquidityManager.getLocalLiquidityContainer(), address(newPool));
  }

  function test_setLocalLiquidityContainer_RevertWhen_CallerNotOwner() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    s_liquidityManager.setLocalLiquidityContainer(MockLockReleaseTokenPool(address(1)));
  }

  function test_setLocalLiquidityContainer_RevertWhen_CalledWithTheZeroAddress() external {
    vm.expectRevert(LiquidityManager.ZeroAddress.selector);
    s_liquidityManager.setLocalLiquidityContainer(MockLockReleaseTokenPool(address(0)));
  }
}
