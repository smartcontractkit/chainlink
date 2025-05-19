// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_proposeAdministrator is TokenAdminRegistrySetup {
  function test_proposeAdministrator_module() public {
    vm.startPrank(s_registryModule);
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).pendingAdministrator, newAdmin);
    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).administrator, address(0));
    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).tokenPool, address(0));

    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
  }

  function test_proposeAdministrator_owner() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).pendingAdministrator, newAdmin);

    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
  }

  function test_proposeAdministrator_reRegisterWhileUnclaimed() public {
    address newAdmin = makeAddr("wrongAddress");
    address newToken = makeAddr("newToken");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    assertEq(s_tokenAdminRegistry.getTokenConfig(newToken).pendingAdministrator, newAdmin);

    newAdmin = makeAddr("correctAddress");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(newToken, address(0), newAdmin);

    // Ensure we can still register the correct admin while the previous admin is unclaimed.
    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);

    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
  }

  mapping(address token => address admin) internal s_AdminByToken;

  function testFuzz_proposeAdministrator_Success(address[50] memory tokens, address[50] memory admins) public {
    TokenAdminRegistry cleanTokenAdminRegistry = new TokenAdminRegistry();
    for (uint256 i = 0; i < tokens.length; i++) {
      if (admins[i] == address(0)) {
        continue;
      }
      if (cleanTokenAdminRegistry.getTokenConfig(tokens[i]).administrator != address(0)) {
        continue;
      }
      cleanTokenAdminRegistry.proposeAdministrator(tokens[i], admins[i]);
      s_AdminByToken[tokens[i]] = admins[i];
    }

    for (uint256 i = 0; i < tokens.length; i++) {
      assertEq(cleanTokenAdminRegistry.getTokenConfig(tokens[i]).pendingAdministrator, s_AdminByToken[tokens[i]]);
    }
  }

  function test_RevertWhen_proposeAdministrator_OnlyRegistryModule() public {
    address newToken = makeAddr("newToken");
    vm.stopPrank();

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.OnlyRegistryModuleOrOwner.selector, address(this)));
    s_tokenAdminRegistry.proposeAdministrator(newToken, OWNER);
  }

  function test_RevertWhen_proposeAdministrator_ZeroAddress() public {
    address newToken = makeAddr("newToken");

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.ZeroAddress.selector));
    s_tokenAdminRegistry.proposeAdministrator(newToken, address(0));
  }

  function test_RevertWhen_proposeAdministrator_AlreadyRegistered() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);
    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    changePrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.AlreadyRegistered.selector, newToken));
    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);
  }
}
