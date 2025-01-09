// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";

contract SiloedLockReleaseTokenPool_withdrawLiqudity is SiloedLockReleaseTokenPoolSetup {
  function setUp() public override {
    super.setUp();

    s_siloedLockReleaseTokenPool.setRebalancer(SILOED_CHAIN_SELECTOR, OWNER);
  }

  function test_withdrawLiquidity_SiloedFunds() public {
    uint256 amount = 1e24;

    uint256 balanceBefore = s_token.balanceOf(OWNER);

    // Provide the Liquidity first
    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityRemoved(SILOED_CHAIN_SELECTOR, OWNER, amount);

    // Remove the Liquidity
    s_siloedLockReleaseTokenPool.withdrawLiquidity(SILOED_CHAIN_SELECTOR, amount);

    assertEq(s_token.balanceOf(OWNER), balanceBefore);
    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), 0);
  }

  function test_withdrawLiquidity_UnsiloedFunds_Success() public {
    uint256 amount = 1e24;

    uint256 balanceBefore = s_token.balanceOf(OWNER);

    // Provide the Liquidity first
    s_siloedLockReleaseTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, amount);

    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), amount);

    vm.expectEmit();
    emit SiloedLockReleaseTokenPool.LiquidityRemoved(DEST_CHAIN_SELECTOR, OWNER, amount);

    // Remove the Liquidity
    s_siloedLockReleaseTokenPool.withdrawLiquidity(DEST_CHAIN_SELECTOR, amount);

    assertEq(s_token.balanceOf(OWNER), balanceBefore);
    assertEq(s_token.balanceOf(address(s_siloedLockReleaseTokenPool)), 0);
    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), 0);
  }

  // Reverts

  function test_withdrawLiquidity_RevertWhen_SiloedFunds_NotEnoughLiquidity() public {
    uint256 amount = 1e24;

    s_siloedLockReleaseTokenPool.provideLiquidity(SILOED_CHAIN_SELECTOR, amount);

    // Call should revert due to underflow error due to trying to burn more tokens than are locked via CCIP.
    vm.expectRevert(abi.encodeWithSignature("Panic(uint256)", 0x11));

    s_siloedLockReleaseTokenPool.withdrawLiquidity(SILOED_CHAIN_SELECTOR, amount + 1);
  }

  function test_withdrawLiquidity_RevertWhen_UnsiloedFunds_NotEnoughLiquidity() public {
    uint256 amount = 1e24;

    s_siloedLockReleaseTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, amount);

    // Call should revert due to underflow error due to trying to burn more tokens than are locked via CCIP.
    vm.expectRevert(abi.encodeWithSignature("Panic(uint256)", 0x11));

    // Test withdrawing funds from unsiloed liquidity pool but underflow
    s_siloedLockReleaseTokenPool.withdrawLiquidity(DEST_CHAIN_SELECTOR, amount + 1);
  }
}
