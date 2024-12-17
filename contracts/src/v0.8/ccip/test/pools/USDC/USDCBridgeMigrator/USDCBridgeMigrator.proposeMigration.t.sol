// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from "./USDCBridgeMigratorSetup.t.sol";

contract USDCBridgeMigrator_proposeMigration is HybridLockReleaseUSDCTokenPoolSetup {
  function test_RevertWhen_ChainNotUsingLockRelease() public {
    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.InvalidChainSelector.selector));

    vm.startPrank(OWNER);

    s_usdcTokenPool.proposeCCTPMigration(0x98765);
  }
}
