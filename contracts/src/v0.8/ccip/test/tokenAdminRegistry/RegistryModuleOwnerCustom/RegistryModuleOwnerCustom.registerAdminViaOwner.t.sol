// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IOwner} from "../../../interfaces/IOwner.sol";

import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";

import {RegistryModuleOwnerCustomSetup} from "./RegistryModuleOwnerCustomSetup.t.sol";

contract RegistryModuleOwnerCustom_registerAdminViaOwner is RegistryModuleOwnerCustomSetup {
  function test_registerAdminViaOwner() public {
    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, address(0));

    address expectedOwner = IOwner(s_token).owner();

    vm.expectCall(s_token, abi.encodeWithSelector(IOwner.owner.selector), 1);
    vm.expectCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(TokenAdminRegistry.proposeAdministrator.selector, s_token, expectedOwner),
      1
    );

    vm.expectEmit();
    emit RegistryModuleOwnerCustom.AdministratorRegistered(s_token, expectedOwner);

    s_registryModuleOwnerCustom.registerAdminViaOwner(s_token);

    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).pendingAdministrator, OWNER);
  }

  function test_RevertWhen_registerAdminViaOwner() public {
    address expectedOwner = IOwner(s_token).owner();

    vm.startPrank(makeAddr("Not_expected_owner"));

    vm.expectRevert(
      abi.encodeWithSelector(RegistryModuleOwnerCustom.CanOnlySelfRegister.selector, expectedOwner, s_token)
    );

    s_registryModuleOwnerCustom.registerAdminViaOwner(s_token);
  }
}
