// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IGetCCIPAdmin} from "../../interfaces/IGetCCIPAdmin.sol";
import {IOwner} from "../../interfaces/IOwner.sol";

import {RegistryModuleOwnerCustom} from "../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {BurnMintERC677Helper} from "../helpers/BurnMintERC677Helper.sol";

import {Test} from "forge-std/Test.sol";

contract RegistryModuleOwnerCustomSetup is Test {
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;

  RegistryModuleOwnerCustom internal s_registryModuleOwnerCustom;
  TokenAdminRegistry internal s_tokenAdminRegistry;
  address internal s_token;

  function setUp() public virtual {
    vm.startPrank(OWNER);

    s_tokenAdminRegistry = new TokenAdminRegistry();
    s_token = address(new BurnMintERC677Helper("Test", "TST"));
    s_registryModuleOwnerCustom = new RegistryModuleOwnerCustom(address(s_tokenAdminRegistry));
    s_tokenAdminRegistry.addRegistryModule(address(s_registryModuleOwnerCustom));
  }
}

contract RegistryModuleOwnerCustom_constructor is RegistryModuleOwnerCustomSetup {
  function test_constructor_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(RegistryModuleOwnerCustom.AddressZero.selector));

    new RegistryModuleOwnerCustom(address(0));
  }
}

contract RegistryModuleOwnerCustom_registerAdminViaGetCCIPAdmin is RegistryModuleOwnerCustomSetup {
  function test_registerAdminViaGetCCIPAdmin_Success() public {
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
      abi.encodeWithSelector(TokenAdminRegistry.proposeAdministrator.selector, s_token, expectedOwner),
      1
    );

    vm.expectEmit();
    emit RegistryModuleOwnerCustom.AdministratorRegistered(s_token, expectedOwner);

    s_registryModuleOwnerCustom.registerAdminViaOwner(s_token);

    assertEq(s_tokenAdminRegistry.getTokenConfig(s_token).pendingAdministrator, OWNER);
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
