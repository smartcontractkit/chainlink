// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_getCCIPAdmin is BurnMintERC20Setup {
  function test_getCCIPAdmin() public view {
    assertEq(s_alice, s_burnMintERC20.getCCIPAdmin());
  }

  function test_setCCIPAdmin() public {
    address newAdmin = makeAddr("newAdmin");

    vm.expectEmit();
    emit FactoryBurnMintERC20.CCIPAdminTransferred(s_alice, newAdmin);

    s_burnMintERC20.setCCIPAdmin(newAdmin);

    assertEq(newAdmin, s_burnMintERC20.getCCIPAdmin());
  }
}
