// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_getPools is TokenAdminRegistrySetup {
  function test_getPools() public {
    address[] memory tokens = new address[](1);
    tokens[0] = s_sourceTokens[0];

    address[] memory got = s_tokenAdminRegistry.getPools(tokens);
    assertEq(got.length, 1);
    assertEq(got[0], s_sourcePoolByToken[tokens[0]]);

    got = s_tokenAdminRegistry.getPools(s_sourceTokens);
    assertEq(got.length, s_sourceTokens.length);
    for (uint256 i = 0; i < s_sourceTokens.length; i++) {
      assertEq(got[i], s_sourcePoolByToken[s_sourceTokens[i]]);
    }

    address doesNotExist = makeAddr("doesNotExist");
    tokens[0] = doesNotExist;
    got = s_tokenAdminRegistry.getPools(tokens);
    assertEq(got.length, 1);
    assertEq(got[0], address(0));
  }
}
