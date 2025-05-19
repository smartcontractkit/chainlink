// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_getPool is TokenAdminRegistrySetup {
  function test_getPool() public view {
    address got = s_tokenAdminRegistry.getPool(s_sourceTokens[0]);
    assertEq(got, s_sourcePoolByToken[s_sourceTokens[0]]);
  }
}
