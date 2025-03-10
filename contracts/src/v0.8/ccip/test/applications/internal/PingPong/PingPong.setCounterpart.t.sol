// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_setCounterpart is PingPongDappSetup {
  function testFuzz_CounterPartAddress_Success(uint64 chainSelector, address counterpartAddress) public {
    // s_pingPong.setCounterpartChainSelector(chainSelector);
    vm.assume(counterpartAddress != address(0));
    vm.assume(chainSelector != 0);

    s_pingPong.setCounterpart(chainSelector, counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function test_setCounterpart_RevertWhen_InvalidZeroAddress() public {
    vm.expectRevert(CCIPBase.ZeroAddressNotAllowed.selector);
    s_pingPong.setCounterpart(0, address(1));
  }

  function test_setCounterpart_RevertWhen_InvalidZeroChainSelector() public {
    vm.expectRevert(CCIPBase.ZeroAddressNotAllowed.selector);
    s_pingPong.setCounterpart(0, address(1));
  }
}
