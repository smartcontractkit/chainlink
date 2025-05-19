// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_transferAdminRole is TokenAdminRegistrySetup {
  function test_transferAdminRole() public {
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
  }

  function test_RevertWhen_transferAdminRole_OnlyAdministrator() public {
    vm.stopPrank();

    vm.expectRevert(
      abi.encodeWithSelector(TokenAdminRegistry.OnlyAdministrator.selector, address(this), s_sourceTokens[0])
    );
    s_tokenAdminRegistry.transferAdminRole(s_sourceTokens[0], makeAddr("newAdmin"));
  }
}
