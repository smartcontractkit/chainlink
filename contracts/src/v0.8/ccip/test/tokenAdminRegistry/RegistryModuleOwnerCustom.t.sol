// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IGetCCIPAdmin} from "../../interfaces/IGetCCIPAdmin.sol";
import {IOwner} from "../../interfaces/IOwner.sol";

import {RegistryModuleOwnerCustom} from "../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {BurnMintERC677Helper} from "../helpers/BurnMintERC677Helper.sol";

contract RegistryModuleOwnerCustomSetup is TokenSetup {
  event AdministratorRegistered(address indexed token, address indexed administrator);

  RegistryModuleOwnerCustom internal s_registryModuleOwnerCustom;
  address internal s_token;

  function setUp() public virtual override {
    TokenSetup.setUp();

    s_token = address(new BurnMintERC677Helper("Test", "TST"));
    s_registryModuleOwnerCustom = new RegistryModuleOwnerCustom(address(s_tokenAdminRegistry));
    s_tokenAdminRegistry.addRegistryModule(address(s_registryModuleOwnerCustom));
  }
}

contract RegistryModuleOwnerCustom_registerAdminViaGetCCIPAdmin is RegistryModuleOwnerCustomSetup {
  function test_registerAdminViaGetCCIPAdmin_Success() public {
    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, address(0));

    address expectedOwner = IGetCCIPAdmin(s_token).getCCIPAdmin();

    vm.expectCall(s_token, abi.encodeWithSelector(IGetCCIPAdmin.getCCIPAdmin.selector), 1);
    vm.expectCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(TokenAdminRegistry.registerAdministrator.selector, s_token, expectedOwner),
      1
    );

    vm.expectEmit();
    emit AdministratorRegistered(s_token, expectedOwner);

    s_registryModuleOwnerCustom.registerAdminViaGetCCIPAdmin(s_token);

    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, OWNER);
  }

  function test_registerAdminViaGetCCIPAdmin_Revert() public {
    address expectedOwner = IGetCCIPAdmin(s_token).getCCIPAdmin();

    vm.startPrank(makeAddr("Not_expected_owner"));

    vm.expectRevert(
      abi.encodeWithSelector(RegistryModuleOwnerCustom.CanOnlySelfRegister.selector, expectedOwner, s_token)
    );

    s_registryModuleOwnerCustom.registerAdminViaGetCCIPAdmin(s_token);
  }
}

contract RegistryModuleOwnerCustom_registerAdminViaOwner is RegistryModuleOwnerCustomSetup {
  function test_registerAdminViaOwner_Success() public {
    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, address(0));

    address expectedOwner = IOwner(s_token).owner();

    vm.expectCall(s_token, abi.encodeWithSelector(IOwner.owner.selector), 1);
    vm.expectCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(TokenAdminRegistry.registerAdministrator.selector, s_token, expectedOwner),
      1
    );

    vm.expectEmit();
    emit AdministratorRegistered(s_token, expectedOwner);

    s_registryModuleOwnerCustom.registerAdminViaOwner(s_token);

    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).administrator, OWNER);
  }

  function test_registerAdminViaOwner_Revert() public {
    address expectedOwner = IOwner(s_token).owner();

    vm.startPrank(makeAddr("Not_expected_owner"));

    vm.expectRevert(
      abi.encodeWithSelector(RegistryModuleOwnerCustom.CanOnlySelfRegister.selector, expectedOwner, s_token)
    );

    s_registryModuleOwnerCustom.registerAdminViaOwner(s_token);
  }
}
