// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {OnRampSetup} from "../../onRamp/OnRamp/OnRampSetup.t.sol";

contract Router_constructor is OnRampSetup {
  function test_Constructor() public view {
    assertEq("Router 1.2.0", s_sourceRouter.typeAndVersion());
    assertEq(OWNER, s_sourceRouter.owner());
  }
}
