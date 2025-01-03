// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";

import {TokenPool} from "../../../pools/TokenPool.sol";
import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

contract LockReleaseTokenPool_provideLiquidity is LockReleaseTokenPoolSetup {
  function testFuzz_ProvideLiquidity_Success(
    uint256 amount
  ) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPool), amount);

    s_lockReleaseTokenPool.provideLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre - amount);
    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), amount);
  }

  // Reverts

  function test_RevertWhen_Unauthorized() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPool.provideLiquidity(1);
  }

  function testFuzz_ExceedsAllowance(
    uint256 amount
  ) public {
    vm.assume(amount > 0);
    vm.expectRevert("ERC20: insufficient allowance");
    s_lockReleaseTokenPool.provideLiquidity(amount);
  }

  function test_RevertWhen_LiquidityNotAccepted() public {
    s_lockReleaseTokenPool = new LockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), false, address(s_sourceRouter)
    );

    vm.expectRevert(LockReleaseTokenPool.LiquidityNotAccepted.selector);
    s_lockReleaseTokenPool.provideLiquidity(1);
  }
}
