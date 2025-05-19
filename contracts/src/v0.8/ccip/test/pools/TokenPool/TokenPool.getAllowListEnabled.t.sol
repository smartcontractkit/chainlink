// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPoolWithAllowListSetup} from "./TokenPoolWithAllowListSetup.t.sol";

contract TokenPoolWithAllowList_getAllowListEnabled is TokenPoolWithAllowListSetup {
  function test_GetAllowListEnabled() public view {
    assertTrue(s_tokenPool.getAllowListEnabled());
  }
}
