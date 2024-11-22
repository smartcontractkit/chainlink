// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract BurnMintERC20_getCCIPAdmin is BurnMintERC20Setup {
  function test_getCCIPAdmin() public view {
    assertEq(OWNER, s_burnMintERC20.getCCIPAdmin());
  }

  function test_setCCIPAdmin() public {
    address newAdmin = makeAddr("newAdmin");

    vm.expectEmit();
    emit BurnMintERC20.CCIPAdminTransferred(OWNER, newAdmin);

    s_burnMintERC20.setCCIPAdmin(newAdmin);

    assertEq(newAdmin, s_burnMintERC20.getCCIPAdmin());
  }
}
