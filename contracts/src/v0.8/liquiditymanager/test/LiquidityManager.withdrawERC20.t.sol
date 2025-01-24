// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";
import {LiquidityManager} from "../LiquidityManager.sol";

contract LiquidityManager_withdrawERC20 is LiquidityManagerSetup {
  function test_withdrawERC20() external {
    uint256 amount = 100;
    deal(address(s_otherToken), address(s_liquidityManager), amount);

    assertEq(s_otherToken.balanceOf(address(1)), 0);
    assertEq(s_otherToken.balanceOf(address(s_liquidityManager)), amount);

    vm.startPrank(FINANCE);
    s_liquidityManager.withdrawERC20(address(s_otherToken), amount, address(1));

    assertEq(s_otherToken.balanceOf(address(1)), amount);
    assertEq(s_otherToken.balanceOf(address(s_liquidityManager)), 0);
  }

  function test_withdrawERC20_RevertWhen_InvalidCondition() external {
    uint256 amount = 100;
    deal(address(s_otherToken), address(s_liquidityManager), amount);

    vm.startPrank(STRANGER);
    vm.expectRevert(LiquidityManager.OnlyFinanceRole.selector);

    s_liquidityManager.withdrawERC20(address(s_otherToken), amount, address(1));
  }
}
