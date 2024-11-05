// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";
import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";

contract BurnMintERC20getCCIPAdmin is BurnMintERC20Setup {
  function test_getCCIPAdmin_Success() public view {
    assertEq(OWNER, s_burnMintERC20.getCCIPAdmin());
  }

  function test_setCCIPAdmin_Success() public {
    address newAdmin = makeAddr("newAdmin");

    vm.expectEmit();
    emit BurnMintERC20.CCIPAdminTransferred(OWNER, newAdmin);

    s_burnMintERC20.setCCIPAdmin(newAdmin);

    assertEq(newAdmin, s_burnMintERC20.getCCIPAdmin());
  }
}
