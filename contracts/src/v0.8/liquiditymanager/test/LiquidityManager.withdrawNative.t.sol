// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LiquidityManager} from "../LiquidityManager.sol";
import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";

contract LiquidityManager_withdrawNative is LiquidityManagerSetup {
  address private receiver = makeAddr("receiver");

  function setUp() public override {
    super.setUp();
    vm.deal(address(s_liquidityManager), 1);
  }

  function test_withdrawNative() external {
    assertEq(receiver.balance, 0);

    vm.expectEmit();
    emit LiquidityManager.NativeWithdrawn(1, receiver);

    changePrank(FINANCE);
    s_liquidityManager.withdrawNative(1, payable(receiver));

    assertEq(receiver.balance, 1);
  }

  function test_withdrawNative_RevertWhen_NotFinanceRole() external {
    vm.stopPrank();

    vm.expectRevert(LiquidityManager.OnlyFinanceRole.selector);

    s_liquidityManager.withdrawNative(1, payable(receiver));
  }
}
