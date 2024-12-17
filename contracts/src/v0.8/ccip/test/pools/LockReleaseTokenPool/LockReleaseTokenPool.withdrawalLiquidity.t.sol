// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";

import {TokenPool} from "../../../pools/TokenPool.sol";
import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

contract LockReleaseTokenPool_withdrawalLiquidity is LockReleaseTokenPoolSetup {
  function testFuzz_WithdrawalLiquidity_Success(
    uint256 amount
  ) public {
    uint256 balancePre = s_token.balanceOf(OWNER);
    s_token.approve(address(s_lockReleaseTokenPool), amount);
    s_lockReleaseTokenPool.provideLiquidity(amount);

    s_lockReleaseTokenPool.withdrawLiquidity(amount);

    assertEq(s_token.balanceOf(OWNER), balancePre);
  }

  // Reverts
  function test_RevertWhen_Unauthorized() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, STRANGER));

    s_lockReleaseTokenPool.withdrawLiquidity(1);
  }

  function test_RevertWhen_InsufficientLiquidity() public {
    uint256 maxUint256 = 2 ** 256 - 1;
    s_token.approve(address(s_lockReleaseTokenPool), maxUint256);
    s_lockReleaseTokenPool.provideLiquidity(maxUint256);

    vm.startPrank(address(s_lockReleaseTokenPool));
    s_token.transfer(OWNER, maxUint256);
    vm.startPrank(OWNER);

    vm.expectRevert(LockReleaseTokenPool.InsufficientLiquidity.selector);
    s_lockReleaseTokenPool.withdrawLiquidity(1);
  }
}
