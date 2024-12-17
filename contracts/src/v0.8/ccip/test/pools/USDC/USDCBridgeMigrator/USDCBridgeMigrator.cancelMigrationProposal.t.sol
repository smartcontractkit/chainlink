// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from "./USDCBridgeMigratorSetup.t.sol";

contract USDCBridgeMigrator_cancelMigrationProposal is HybridLockReleaseUSDCTokenPoolSetup {
  function test_cancelExistingCCTPMigrationProposal() public {
    vm.startPrank(OWNER);

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationProposed(DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    assertEq(
      s_usdcTokenPool.getCurrentProposedCCTPChainMigration(),
      DEST_CHAIN_SELECTOR,
      "migration proposal should exist, but doesn't"
    );

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationCancelled(DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.cancelExistingCCTPMigrationProposal();

    assertEq(
      s_usdcTokenPool.getCurrentProposedCCTPChainMigration(),
      0,
      "migration proposal exists, but shouldn't after being cancelled"
    );

    vm.expectRevert(USDCBridgeMigrator.NoMigrationProposalPending.selector);
    s_usdcTokenPool.cancelExistingCCTPMigrationProposal();
  }

  function test_RevertWhen_cannotCancelANonExistentMigrationProposal() public {
    vm.expectRevert(USDCBridgeMigrator.NoMigrationProposalPending.selector);

    // Proposal to migrate doesn't exist, and so the chain selector is zero, and therefore should revert
    s_usdcTokenPool.cancelExistingCCTPMigrationProposal();
  }
}
