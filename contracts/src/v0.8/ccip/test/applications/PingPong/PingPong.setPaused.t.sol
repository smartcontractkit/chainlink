// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_setPaused is PingPongDappSetup {
  function test_Pausing() public {
    assertFalse(s_pingPong.isPaused());

    s_pingPong.setPaused(true);

    assertTrue(s_pingPong.isPaused());
  }
}
