// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";
import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";

contract BurnMintERC20grantMintAndBurnRoles is BurnMintERC20Setup {
  function test_GrantMintAndBurnRoles_Success() public {
    assertFalse(s_burnMintERC20.isMinter(STRANGER));
    assertFalse(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit BurnMintERC20.MintAccessGranted(STRANGER);
    vm.expectEmit();
    emit BurnMintERC20.BurnAccessGranted(STRANGER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.isMinter(STRANGER));
    assertTrue(s_burnMintERC20.isBurner(STRANGER));
  }
}