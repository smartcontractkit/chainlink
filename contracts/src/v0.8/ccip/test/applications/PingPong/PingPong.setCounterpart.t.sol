// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_setCounterpart is PingPongDappSetup {
  function testFuzz_CounterPartAddress_Success(uint64 chainSelector, address counterpartAddress) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    s_pingPong.setCounterpart(chainSelector, counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }
}
