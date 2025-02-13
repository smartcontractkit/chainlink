// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

contract CCIPReceiver_typeAndVersion is CCIPReceiverSetup {
  function test_typeAndVersion() public view {
    assertEq(s_receiver.typeAndVersion(), "CCIPReceiver 1.6.0-dev");
  }
}
