// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";
import {LiquidityManager} from "../LiquidityManager.sol";

contract LiquidityManager_setFinanceRole is LiquidityManagerSetup {
  function test_setFinanceRole() external {
    vm.expectEmit();
    address newFinanceRole = makeAddr("newFinanceRole");
    assertEq(s_liquidityManager.getFinanceRole(), FINANCE);
    emit LiquidityManager.FinanceRoleSet(newFinanceRole);

    s_liquidityManager.setFinanceRole(newFinanceRole);

    assertEq(s_liquidityManager.getFinanceRole(), newFinanceRole);
  }

  function test_setFinanceRole_RevertWhen_CallerNotOwner() external {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");

    s_liquidityManager.setFinanceRole(address(1));
  }
}
