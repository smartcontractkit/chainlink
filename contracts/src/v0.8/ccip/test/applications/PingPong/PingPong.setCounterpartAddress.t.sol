// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_setCounterpartAddress is PingPongDappSetup {
  function testFuzz_CounterPartAddress_Success(
    address counterpartAddress
  ) public {
    s_pingPong.setCounterpartAddress(counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
  }
}
