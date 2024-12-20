// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_grantMintAndBurnRoles is BurnMintERC20Setup {
  function test_GrantMintAndBurnRoles() public {
    assertFalse(s_burnMintERC20.isMinter(STRANGER));
    assertFalse(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.MintAccessGranted(STRANGER);
    vm.expectEmit();
    emit FactoryBurnMintERC20.BurnAccessGranted(STRANGER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.isMinter(STRANGER));
    assertTrue(s_burnMintERC20.isBurner(STRANGER));
  }
}
