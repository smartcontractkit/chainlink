// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPSenderSetup} from "./CCIPSenderSetup.t.sol";

contract CCIPSender_typeAndVersion is CCIPSenderSetup {
  function test_typeAndVersion() public {
    assertEq(s_sender.typeAndVersion(), "CCIPSender 1.6.0-dev");
  }
}
