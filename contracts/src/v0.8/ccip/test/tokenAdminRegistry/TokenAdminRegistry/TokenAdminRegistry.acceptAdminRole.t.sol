// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_acceptAdminRole is TokenAdminRegistrySetup {
  function test_acceptAdminRole() public {
    address token = s_sourceTokens[0];

    address currentAdmin = s_tokenAdminRegistry.getTokenConfig(token).administrator;
    address newAdmin = makeAddr("newAdmin");

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferRequested(token, currentAdmin, newAdmin);

    s_tokenAdminRegistry.transferAdminRole(token, newAdmin);

    TokenAdminRegistry.TokenConfig memory config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, newAdmin);
    assertEq(config.administrator, currentAdmin);

    vm.startPrank(newAdmin);

    vm.expectEmit();
    emit TokenAdminRegistry.AdministratorTransferred(token, newAdmin);

    s_tokenAdminRegistry.acceptAdminRole(token);

    config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, address(0));
    assertEq(config.administrator, newAdmin);
  }

  function test_RevertWhen_acceptAdminRole_OnlyPendingAdministrator() public {
    address token = s_sourceTokens[0];
    address currentAdmin = s_tokenAdminRegistry.getTokenConfig(token).administrator;
    address newAdmin = makeAddr("newAdmin");

    s_tokenAdminRegistry.transferAdminRole(token, newAdmin);

    TokenAdminRegistry.TokenConfig memory config = s_tokenAdminRegistry.getTokenConfig(token);

    // Assert only the pending admin updates, without affecting the pending admin.
    assertEq(config.pendingAdministrator, newAdmin);
    assertEq(config.administrator, currentAdmin);

    address notNewAdmin = makeAddr("notNewAdmin");
    vm.startPrank(notNewAdmin);

    vm.expectRevert(abi.encodeWithSelector(TokenAdminRegistry.OnlyPendingAdministrator.selector, notNewAdmin, token));
    s_tokenAdminRegistry.acceptAdminRole(token);
  }
}
