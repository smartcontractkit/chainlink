// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_getCapabilityConfiguration is CCIPHomeTestSetup {
  function test_getCapabilityConfiguration() public view {
    bytes memory config = s_ccipHome.getCapabilityConfiguration(DEFAULT_DON_ID);
    assertEq(config.length, 0);
  }
}
