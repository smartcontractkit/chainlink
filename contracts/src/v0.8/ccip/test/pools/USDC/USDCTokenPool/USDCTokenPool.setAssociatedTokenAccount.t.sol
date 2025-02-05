// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../../shared/access/Ownable2Step.sol";
import {USDCTokenPoolSetup} from "./USDCTokenPoolSetup.t.sol";

contract USDCTokenPool_setAssociatedTokenAccount is USDCTokenPoolSetup {
  function test_setAssociatedTokenAccount() public {
    s_usdcTokenPool.setAssociatedTokenAccount(DEST_CHAIN_SELECTOR, bytes32("ATA"));
    assertEq(s_usdcTokenPool.getAssociatedTokenAccount(DEST_CHAIN_SELECTOR), bytes32("ATA"));
  }

  function test_setAssociatedTokenAccount_RevertsWhen_NotOwner() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_usdcTokenPool.setAssociatedTokenAccount(DEST_CHAIN_SELECTOR, bytes32("ATA"));
  }
}
