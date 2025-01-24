// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPClientSetup} from "./CCIPClientSetup.t.sol";

contract CCIPClient_typeAndVersion is CCIPClientSetup {
  function test_typeAndVersion() public view {
    assertEq(s_sender.typeAndVersion(), "CCIPClient 1.6.0-dev");
  }
}
