// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ILiquidityContainer} from "../../../../../liquiditymanager/interfaces/ILiquidityContainer.sol";

import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from "./HybridLockReleaseUSDCTokenPoolSetup.t.sol";

contract HybridLockReleaseUSDCTokenPool_TransferLiquidity is HybridLockReleaseUSDCTokenPoolSetup {
  function test_transferLiquidity_Success() public {
    // Set as the OWNER so we can provide liquidity
    vm.startPrank(OWNER);

    s_usdcTokenPool.setLiquidityProvider(DEST_CHAIN_SELECTOR, OWNER);
    s_token.approve(address(s_usdcTokenPool), type(uint256).max);

    uint256 liquidityAmount = 1e9;

    // Provide some liquidity to the pool
    s_usdcTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, liquidityAmount);

    // Set the new token pool as the rebalancer
    s_usdcTokenPool.transferOwnership(address(s_usdcTokenPoolTransferLiquidity));

    vm.expectEmit();
    emit ILiquidityContainer.LiquidityRemoved(address(s_usdcTokenPoolTransferLiquidity), liquidityAmount);

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LiquidityTransferred(
      address(s_usdcTokenPool), DEST_CHAIN_SELECTOR, liquidityAmount
    );

    s_usdcTokenPoolTransferLiquidity.transferLiquidity(address(s_usdcTokenPool), DEST_CHAIN_SELECTOR);

    assertEq(
      s_usdcTokenPool.owner(),
      address(s_usdcTokenPoolTransferLiquidity),
      "Ownership of the old pool should be transferred to the new pool"
    );

    assertEq(
      s_usdcTokenPoolTransferLiquidity.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      liquidityAmount,
      "Tokens locked for dest chain doesn't match expected amount in storage"
    );

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      0,
      "Tokens locked for dest chain in old token pool doesn't match expected amount in storage"
    );

    assertEq(
      s_token.balanceOf(address(s_usdcTokenPoolTransferLiquidity)),
      liquidityAmount,
      "Liquidity amount of tokens should be new in new pool, but aren't"
    );

    assertEq(
      s_token.balanceOf(address(s_usdcTokenPool)),
      0,
      "Liquidity amount of tokens should be zero in old pool, but aren't"
    );
  }

  function test_cannotTransferLiquidityDuringPendingMigration_Revert() public {
    // Set as the OWNER so we can provide liquidity
    vm.startPrank(OWNER);

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    s_usdcTokenPool.setLiquidityProvider(DEST_CHAIN_SELECTOR, OWNER);
    s_token.approve(address(s_usdcTokenPool), type(uint256).max);

    uint256 liquidityAmount = 1e9;

    // Provide some liquidity to the pool
    s_usdcTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, liquidityAmount);

    // Set the new token pool as the rebalancer
    s_usdcTokenPool.transferOwnership(address(s_usdcTokenPoolTransferLiquidity));

    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    vm.expectRevert(
      abi.encodeWithSelector(HybridLockReleaseUSDCTokenPool.LanePausedForCCTPMigration.selector, DEST_CHAIN_SELECTOR)
    );

    s_usdcTokenPoolTransferLiquidity.transferLiquidity(address(s_usdcTokenPool), DEST_CHAIN_SELECTOR);
  }
}
