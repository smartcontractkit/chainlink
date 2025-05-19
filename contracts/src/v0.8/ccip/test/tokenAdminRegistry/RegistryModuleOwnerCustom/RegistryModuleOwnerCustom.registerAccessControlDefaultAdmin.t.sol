// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";

import {AccessControl} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/AccessControl.sol";

import {RegistryModuleOwnerCustomSetup} from "./RegistryModuleOwnerCustomSetup.t.sol";

contract AccessController is AccessControl {
  constructor(
    address admin
  ) {
    _grantRole(DEFAULT_ADMIN_ROLE, admin);
  }
}

contract RegistryModuleOwnerCustom_registerAccessControlDefaultAdmin is RegistryModuleOwnerCustomSetup {
  function setUp() public override {
    super.setUp();

    s_token = address(new AccessController(OWNER));
  }

  function test_registerAccessControlDefaultAdmin() public {
    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, address(0));

    bytes32 defaultAdminRole = AccessController(s_token).DEFAULT_ADMIN_ROLE();

    vm.expectCall(address(s_token), abi.encodeWithSelector(AccessControl.hasRole.selector, defaultAdminRole, OWNER), 1);
    vm.expectCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(TokenAdminRegistry.proposeAdministrator.selector, s_token, OWNER),
      1
    );

    vm.expectEmit();
    emit RegistryModuleOwnerCustom.AdministratorRegistered(s_token, OWNER);

    s_registryModuleOwnerCustom.registerAccessControlDefaultAdmin(s_token);

    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).pendingAdministrator, OWNER);
  }

  function test_RevertWhen_registerAccessControlDefaultAdmin() public {
    bytes32 defaultAdminRole = AccessController(s_token).DEFAULT_ADMIN_ROLE();

    address wrongSender = makeAddr("Not_expected_owner");
    vm.startPrank(wrongSender);

    vm.expectRevert(
      abi.encodeWithSelector(
        RegistryModuleOwnerCustom.RequiredRoleNotFound.selector, wrongSender, defaultAdminRole, s_token
      )
    );

    s_registryModuleOwnerCustom.registerAccessControlDefaultAdmin(s_token);
  }
}
