// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import "./LiquidityManagerSetup.t.sol";

contract LiquidityManager_removeLiquidity is LiquidityManagerSetup {
  function test_removeLiquiditySuccess() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    vm.expectEmit();
    emit LiquidityManager.LiquidityRemovedFromContainer(FINANCE, amount);

    changePrank(FINANCE);
    s_liquidityManager.removeLiquidity(amount);

    assertEq(s_l1Token.balanceOf(address(s_liquidityManager)), 0);
  }

  function test_removeLiquidity_RevertWhen_InsufficientLiquidity() external {
    uint256 balance = 923;
    uint256 requested = balance + 1;

    deal(address(s_l1Token), address(s_lockReleaseTokenPool), balance);

    vm.expectRevert(abi.encodeWithSelector(LiquidityManager.InsufficientLiquidity.selector, requested, balance, 0));

    changePrank(FINANCE);
    s_liquidityManager.removeLiquidity(requested);
  }

  function test_removeLiquidity_RevertWhen_NotFinanceRole() external {
    vm.stopPrank();

    vm.expectRevert(LiquidityManager.OnlyFinanceRole.selector);

    s_liquidityManager.removeLiquidity(123);
  }
}
