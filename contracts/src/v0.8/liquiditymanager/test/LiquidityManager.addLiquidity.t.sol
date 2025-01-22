// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import "./LiquidityManagerSetup.t.sol";
import {LiquidityManager} from "../LiquidityManager.sol";

contract LiquidityManager_addLiquidity is LiquidityManagerSetup {
  function test_addLiquidity() external {
    address caller = STRANGER;
    changePrank(caller);

    uint256 amount = 12345679;
    deal(address(s_l1Token), caller, amount);

    s_l1Token.approve(address(s_liquidityManager), amount);

    vm.expectEmit();
    emit LiquidityManager.LiquidityAddedToContainer(caller, amount);

    s_liquidityManager.addLiquidity(amount);

    assertEq(s_l1Token.balanceOf(address(s_lockReleaseTokenPool)), amount);
  }
}
