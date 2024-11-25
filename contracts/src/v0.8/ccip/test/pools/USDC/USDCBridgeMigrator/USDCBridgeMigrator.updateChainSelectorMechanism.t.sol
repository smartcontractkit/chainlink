// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";

import {USDCBridgeMigrator_BurnLockedUSDC} from "./USDCBridgeMigrator.burnLockedUSDC.t.sol";

contract USDCBridgeMigrator_updateChainSelectorMechanism is USDCBridgeMigrator_BurnLockedUSDC {
  function test_RevertWhen_cannotRevertChainMechanism_afterMigration() public {
    test_lockOrBurn_then_BurnInCCTPMigration();

    vm.startPrank(OWNER);

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    vm.expectRevert(
      abi.encodeWithSelector(
        HybridLockReleaseUSDCTokenPool.TokenLockingNotAllowedAfterMigration.selector, DEST_CHAIN_SELECTOR
      )
    );

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);
  }
}
