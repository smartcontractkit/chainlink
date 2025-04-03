// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IGetCCIPAdmin} from "../../../../../ccip/interfaces/IGetCCIPAdmin.sol";
import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";

import {ERC20UpgradableBaseTest, IERC20UpgradeableBase} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_roles is ERC20UpgradableBaseTest {
  function should_GrantMintAndBurnRoles(address implementation, bytes32 MINTER_ROLE, bytes32 BURNER_ROLE) public {
    assertFalse(IAccessControl(implementation).hasRole(MINTER_ROLE, STRANGER));
    assertFalse(IAccessControl(implementation).hasRole(BURNER_ROLE, STRANGER));

    changePrank(DEFAULT_ADMIN);

    vm.expectEmit();
    emit IAccessControl.RoleGranted(MINTER_ROLE, STRANGER, DEFAULT_ADMIN);
    vm.expectEmit();
    emit IAccessControl.RoleGranted(BURNER_ROLE, STRANGER, DEFAULT_ADMIN);

    IERC20UpgradeableBase(implementation).grantMintAndBurnRoles(STRANGER);

    assertTrue(IAccessControl(implementation).hasRole(MINTER_ROLE, STRANGER));
    assertTrue(IAccessControl(implementation).hasRole(BURNER_ROLE, STRANGER));
  }

  function should_GetCCIPAdmin(address implementation) public view {
    assertEq(IGetCCIPAdmin(implementation).getCCIPAdmin(), DEFAULT_ADMIN);
  }

  function should_SetCCIPAdmin(address implementation) public {
    changePrank(DEFAULT_ADMIN);

    vm.expectEmit();
    emit IERC20UpgradeableBase.CCIPAdminTransferred(DEFAULT_ADMIN, STRANGER);

    IERC20UpgradeableBase(implementation).setCCIPAdmin(STRANGER);

    assertEq(IGetCCIPAdmin(implementation).getCCIPAdmin(), STRANGER);
  }
}
