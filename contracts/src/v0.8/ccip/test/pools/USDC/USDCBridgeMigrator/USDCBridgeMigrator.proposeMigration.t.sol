// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../../../shared/token/ERC20/IBurnMintERC20.sol";

import {BurnMintERC20} from "../../../../../shared/token/ERC20/BurnMintERC20.sol";
import {Router} from "../../../../Router.sol";

import {TokenPool} from "../../../../pools/TokenPool.sol";
import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from "../USDCTokenPoolSetup.t.sol";

contract USDCBridgeMigrator_proposeMigration is HybridLockReleaseUSDCTokenPoolSetup {
  function test_ChainNotUsingLockRelease_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.InvalidChainSelector.selector));

    vm.startPrank(OWNER);

    s_usdcTokenPool.proposeCCTPMigration(0x98765);
  }
}
