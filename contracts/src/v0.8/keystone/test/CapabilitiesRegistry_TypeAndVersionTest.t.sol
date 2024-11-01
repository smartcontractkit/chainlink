// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./BaseTest.t.sol";

contract CapabilitiesRegistry_TypeAndVersionTest is BaseTest {
  function test_TypeAndVersion() public view {
    assertEq(s_CapabilitiesRegistry.typeAndVersion(), "CapabilitiesRegistry 1.1.0");
  }
}
