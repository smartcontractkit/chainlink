// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_isAdministrator is TokenAdminRegistrySetup {
  function test_isAdministrator() public {
    address newAdmin = makeAddr("newAdmin");
    address newToken = makeAddr("newToken");
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, OWNER));

    s_tokenAdminRegistry.proposeAdministrator(newToken, newAdmin);
    changePrank(newAdmin);
    s_tokenAdminRegistry.acceptAdminRole(newToken);

    assertTrue(s_tokenAdminRegistry.isAdministrator(newToken, newAdmin));
    assertFalse(s_tokenAdminRegistry.isAdministrator(newToken, OWNER));
  }
}
