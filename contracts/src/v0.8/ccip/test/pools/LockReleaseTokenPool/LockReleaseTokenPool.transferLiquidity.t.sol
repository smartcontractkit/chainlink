// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

contract LockReleaseTokenPool_transferLiquidity is LockReleaseTokenPoolSetup {
  LockReleaseTokenPool internal s_oldLockReleaseTokenPool;
  uint256 internal s_amount = 100000;

  function setUp() public virtual override {
    super.setUp();

    s_oldLockReleaseTokenPool = new LockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), true, address(s_sourceRouter)
    );

    deal(address(s_token), address(s_oldLockReleaseTokenPool), s_amount);
  }

  function test_transferLiquidity() public {
    uint256 balancePre = s_token.balanceOf(address(s_lockReleaseTokenPool));

    s_oldLockReleaseTokenPool.setRebalancer(address(s_lockReleaseTokenPool));

    vm.expectEmit();
    emit LockReleaseTokenPool.LiquidityTransferred(address(s_oldLockReleaseTokenPool), s_amount);

    s_lockReleaseTokenPool.transferLiquidity(address(s_oldLockReleaseTokenPool), s_amount);

    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), balancePre + s_amount);
  }

  function test_RevertWhen_transferLiquidity_transferTooMuch() public {
    uint256 balancePre = s_token.balanceOf(address(s_lockReleaseTokenPool));

    s_oldLockReleaseTokenPool.setRebalancer(address(s_lockReleaseTokenPool));

    vm.expectRevert(LockReleaseTokenPool.InsufficientLiquidity.selector);
    s_lockReleaseTokenPool.transferLiquidity(address(s_oldLockReleaseTokenPool), s_amount + 1);

    assertEq(s_token.balanceOf(address(s_lockReleaseTokenPool)), balancePre);
  }
}
