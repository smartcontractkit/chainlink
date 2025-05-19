// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_getAllConfiguredTokens is TokenAdminRegistrySetup {
  function testFuzz_getAllConfiguredTokens_Success(
    uint8 numberOfTokens
  ) public {
    TokenAdminRegistry cleanTokenAdminRegistry = new TokenAdminRegistry();
    for (uint160 i = 0; i < numberOfTokens; ++i) {
      cleanTokenAdminRegistry.proposeAdministrator(address(i), address(i + 1000));
    }

    uint160 count = 0;
    for (uint160 start = 0; start < numberOfTokens; start += count++) {
      address[] memory got = cleanTokenAdminRegistry.getAllConfiguredTokens(uint64(start), uint64(count));
      if (start + count > numberOfTokens) {
        assertEq(got.length, numberOfTokens - start);
      } else {
        assertEq(got.length, count);
      }

      for (uint160 j = 0; j < got.length; ++j) {
        assertEq(got[j], address(j + start));
      }
    }
  }

  function test_getAllConfiguredTokens_outOfBounds() public view {
    address[] memory tokens = s_tokenAdminRegistry.getAllConfiguredTokens(type(uint64).max, 10);
    assertEq(tokens.length, 0);
  }
}
