// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_setCounterpartChainSelector is PingPongDappSetup {
  function testFuzz_CounterPartChainSelector_Success(
    uint64 chainSelector
  ) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }
}
