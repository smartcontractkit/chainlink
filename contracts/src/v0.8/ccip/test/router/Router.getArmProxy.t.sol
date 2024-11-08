// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {RouterSetup} from "./RouterSetup.t.sol";

contract Router_getArmProxy is RouterSetup {
  function test_getArmProxy() public view {
    assertEq(s_sourceRouter.getArmProxy(), address(s_mockRMN));
  }
}
