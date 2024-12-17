// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDemo} from "../../../applications/PingPongDemo.sol";
import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_setOutOfOrderExecution is PingPongDappSetup {
  function test_OutOfOrderExecution() public {
    assertFalse(s_pingPong.getOutOfOrderExecution());

    vm.expectEmit();
    emit PingPongDemo.OutOfOrderExecutionChange(true);

    s_pingPong.setOutOfOrderExecution(true);

    assertTrue(s_pingPong.getOutOfOrderExecution());
  }
}
