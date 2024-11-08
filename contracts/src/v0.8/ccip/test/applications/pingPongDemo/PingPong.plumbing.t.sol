// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {PingPongDemo} from "../../../applications/PingPongDemo.sol";
import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_plumbing is PingPongDappSetup {
  function test_Fuzz_CounterPartChainSelector_Success(
    uint64 chainSelector
  ) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function test_Fuzz_CounterPartAddress_Success(
    address counterpartAddress
  ) public {
    s_pingPong.setCounterpartAddress(counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
  }

  function test_Fuzz_CounterPartAddress_Success(uint64 chainSelector, address counterpartAddress) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    s_pingPong.setCounterpart(chainSelector, counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function test_Pausing_Success() public {
    assertFalse(s_pingPong.isPaused());

    s_pingPong.setPaused(true);

    assertTrue(s_pingPong.isPaused());
  }

  function test_OutOfOrderExecution_Success() public {
    assertFalse(s_pingPong.getOutOfOrderExecution());

    vm.expectEmit();
    emit PingPongDemo.OutOfOrderExecutionChange(true);

    s_pingPong.setOutOfOrderExecution(true);

    assertTrue(s_pingPong.getOutOfOrderExecution());
  }
}
