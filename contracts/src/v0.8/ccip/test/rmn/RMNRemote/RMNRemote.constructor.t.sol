// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_constructor is RMNRemoteSetup {
  function test_constructor() public view {
    assertEq(s_rmnRemote.getLocalChainSelector(), 1);
  }
}
