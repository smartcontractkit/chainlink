// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";
import {LiquidityManager} from "../LiquidityManager.sol";

contract LiquidityManager_receive is LiquidityManagerSetup {
  function test_receive() external {
    address depositor = makeAddr("depositor");
    vm.deal(depositor, 100);
    uint256 before = address(s_liquidityManager).balance;

    vm.expectEmit();
    emit LiquidityManager.NativeDeposited(100, depositor);

    changePrank(depositor);
    payable(address(s_liquidityManager)).transfer(100);

    assertEq(address(s_liquidityManager).balance, before + 100);
  }
}
