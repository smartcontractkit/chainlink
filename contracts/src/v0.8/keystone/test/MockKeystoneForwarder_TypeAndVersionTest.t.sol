// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockBaseTest} from "./MockKeystoneForwarderBaseTest.t.sol";

contract MockKeystoneForwarder_TypeAndVersionTest is MockBaseTest {
  function test_TypeAndVersion() public view {
    assertEq(s_mockForwarder.typeAndVersion(), "MockKeystoneForwarder 1.0.0");
  }
}