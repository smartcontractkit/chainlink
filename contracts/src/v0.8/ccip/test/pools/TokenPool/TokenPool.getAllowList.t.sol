// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPoolWithAllowListSetup} from "./TokenPoolWithAllowListSetup.t.sol";

contract TokenPoolWithAllowList_getAllowList is TokenPoolWithAllowListSetup {
  function test_GetAllowList() public view {
    address[] memory setAddresses = s_tokenPool.getAllowList();
    assertEq(2, setAddresses.length);
    assertEq(s_allowedSenders[0], setAddresses[0]);
    assertEq(s_allowedSenders[1], setAddresses[1]);
  }
}
