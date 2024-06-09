// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

contract CapabilityRegistry_TypeAndVersionTest is BaseTest {
  function test_TypeAndVersion() public view {
    assertEq(s_capabilityRegistry.typeAndVersion(), "CapabilityRegistry 1.0.0");
  }
}
