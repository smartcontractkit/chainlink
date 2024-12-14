// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IGetCCIPAdmin} from "../../../interfaces/IGetCCIPAdmin.sol";

import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";

import {RegistryModuleOwnerCustomSetup} from "./RegistryModuleOwnerCustomSetup.t.sol";

contract RegistryModuleOwnerCustom_registerAdminViaGetCCIPAdmin is RegistryModuleOwnerCustomSetup {
  function test_registerAdminViaGetCCIPAdmin() public {
    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, address(0));

    address expectedOwner = IGetCCIPAdmin(s_token).getCCIPAdmin();

    vm.expectCall(s_token, abi.encodeWithSelector(IGetCCIPAdmin.getCCIPAdmin.selector), 1);
    vm.expectCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(TokenAdminRegistry.proposeAdministrator.selector, s_token, expectedOwner),
      1
    );

    vm.expectEmit();
    emit RegistryModuleOwnerCustom.AdministratorRegistered(s_token, expectedOwner);

    s_registryModuleOwnerCustom.registerAdminViaGetCCIPAdmin(s_token);

    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).pendingAdministrator, OWNER);
  }

  function test_RevertWhen_registerAdminViaGetCCIPAdmin() public {
    address expectedOwner = IGetCCIPAdmin(s_token).getCCIPAdmin();

    vm.startPrank(makeAddr("Not_expected_owner"));

    vm.expectRevert(
      abi.encodeWithSelector(RegistryModuleOwnerCustom.CanOnlySelfRegister.selector, expectedOwner, s_token)
    );

    s_registryModuleOwnerCustom.registerAdminViaGetCCIPAdmin(s_token);
  }
}
