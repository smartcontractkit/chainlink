// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/access/IAccessControl.sol";

contract BurnMintERC20_grantMintAndBurnRoles is BurnMintERC20Setup {
  function test_GrantMintAndBurnRoles() public {
    assertFalse(s_burnMintERC20.hasRole(s_burnMintERC20.MINTER_ROLE(), STRANGER));
    assertFalse(s_burnMintERC20.hasRole(s_burnMintERC20.BURNER_ROLE(), STRANGER));

    vm.expectEmit();
    emit IAccessControl.RoleGranted(s_burnMintERC20.MINTER_ROLE(), STRANGER, OWNER);
    vm.expectEmit();
    emit IAccessControl.RoleGranted(s_burnMintERC20.BURNER_ROLE(), STRANGER, OWNER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.hasRole(s_burnMintERC20.MINTER_ROLE(), STRANGER));
    assertTrue(s_burnMintERC20.hasRole(s_burnMintERC20.BURNER_ROLE(), STRANGER));
  }
}
