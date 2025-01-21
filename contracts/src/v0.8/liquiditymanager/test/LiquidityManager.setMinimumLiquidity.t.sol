// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {LiquidityManagerSetup} from "./LiquidityManagerSetup.t.sol";

contract LiquidityManager_setMinimumLiquidity is LiquidityManagerSetup {
    event MinimumLiquiditySet(uint256 oldBalance, uint256 newBalance);

    function test_setMinimumLiquidity() external {
        vm.expectEmit();
        emit MinimumLiquiditySet(uint256(0), uint256(1000));
        s_liquidityManager.setMinimumLiquidity(1000);
        assertEq(s_liquidityManager.getMinimumLiquidity(), uint256(1000));
    }

    function test_setMinimumLiquidity_RevertWhen_CallerNotOwner() external {
        vm.stopPrank();
        vm.expectRevert("Only callable by owner");
        s_liquidityManager.setMinimumLiquidity(uint256(1000));
    }
}
