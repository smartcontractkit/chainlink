// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPClientSetup} from "./CCIPClientSetup.t.sol";

contract CCIPClient_usePreFundedFeeTokens is CCIPClientSetup {
  function test_usePreFundedFeeTokens() public view {
    assertEq(s_sender.usePreFundedFeeTokens(), false);
  }
}
