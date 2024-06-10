// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./KeystoneForwarderBaseTest.t.sol";

contract KeystoneForwarder_TypeAndVersionTest is BaseTest {
  function test_TypeAndVersion() public view {
    assertEq(s_forwarder.typeAndVersion(), "KeystoneForwarder 1.0.0");
  }
}
