// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_typeAndVersion is PingPongDappSetup {
  function test_typeAndVersion() public view {
    assertEq(s_pingPong.typeAndVersion(), "PingPongDemo 1.6.0-dev");
  }
}
