// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_provideSiloedLiquidity is SiloedLockReleaseTokenPoolSetup {
  address public UNAUTHORIZED_ADDRESS = address(0xdeadbeef);

  function setUp() public override {
    super.setUp();

    s_siloedLockReleaseTokenPool.setSiloRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
  }

  function test_UnsiloedChain() public {
    uint256 amount = 1e24;

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityAdded(DEST_CHAIN_SELECTOR, OWNER, amount);

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(DEST_CHAIN_SELECTOR, amount);

    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), amount);

    // Since the funds for the destination chain are not siloed,
    // the locked token amount should not be increased
    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(DEST_CHAIN_SELECTOR), amount);
    assertEq(s_siloedLockReleaseTokenPool.getUnsiloedLiquidity(), amount);
  }

  function test_SiloedChain() public {
    uint256 amount = 1e24;

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityAdded(SILOED_CHAIN_SELECTOR, OWNER, amount);

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, amount);

    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), amount);
    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), amount);

    // Since the funds for the destination chain are not siloed, the locked token amount should not be increased
    assertEq(s_siloedLockReleaseTokenPool.getUnsiloedLiquidity(), 0);
  }

  // Reverts

  function test_RevertWhen_UnauthorizedForSiloedChain() public {
    vm.startPrank(UNAUTHORIZED_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, UNAUTHORIZED_ADDRESS));

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, 1);
  }

  function test_RevertWhen_UnauthorizedForUnsiloedChain() public {
    vm.startPrank(UNAUTHORIZED_ADDRESS);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, UNAUTHORIZED_ADDRESS));

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(DEST_CHAIN_SELECTOR, 1);
  }

  function test_RevertWhen_LiquidityAmountCannotBeZero() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.LiquidityAmountCannotBeZero.selector));

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, 0);
  }

  function test_RevertWhen_InvalidChainSelector_Zero() public {
    vm.expectRevert(abi.encodeWithSelector(SiloedLockReleaseTokenPool.InvalidChainSelector.selector, 0));

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(0, 1);
  }
}
