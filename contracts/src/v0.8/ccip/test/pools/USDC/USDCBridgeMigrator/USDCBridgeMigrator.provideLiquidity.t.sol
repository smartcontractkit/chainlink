// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPool} from "../../../../pools/TokenPool.sol";
import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCBridgeMigrator_BurnLockedUSDC} from "./USDCBridgeMigrator.burnLockedUSDC.t.sol";

contract USDCBridgeMigrator_provideLiquidity is USDCBridgeMigrator_BurnLockedUSDC {
  function test_RevertWhen_cannotModifyLiquidityWithoutPermissions() public {
    address randomAddr = makeAddr("RANDOM");

    vm.startPrank(randomAddr);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, randomAddr));

    // Revert because there's insufficient permissions for the DEST_CHAIN_SELECTOR to provide liquidity
    s_usdcTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, 1e6);
  }

  function test_RevertWhen_cannotProvideLiquidity_AfterMigration() public {
    test_lockOrBurn_then_BurnInCCTPMigration();

    vm.startPrank(OWNER);

    vm.expectRevert(
      abi.encodeWithSelector(
        HybridLockReleaseUSDCTokenPool.TokenLockingNotAllowedAfterMigration.selector, DEST_CHAIN_SELECTOR
      )
    );

    s_usdcTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, 1e6);
  }

  function test_RevertWhen_cannotProvideLiquidityWhenMigrationProposalPending() public {
    vm.startPrank(OWNER);

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    vm.expectRevert(
      abi.encodeWithSelector(HybridLockReleaseUSDCTokenPool.LanePausedForCCTPMigration.selector, DEST_CHAIN_SELECTOR)
    );
    s_usdcTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, 1e6);
  }
}
