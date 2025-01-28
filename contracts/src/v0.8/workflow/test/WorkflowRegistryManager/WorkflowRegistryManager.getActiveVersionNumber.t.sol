// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_getActiveVersionNumber is WorkflowRegistryManagerSetup {
  function test_WhenNoActiveVersionIsAvailable() external {
    vm.expectRevert(WorkflowRegistryManager.NoActiveVersionAvailable.selector);
    s_registryManager.getActiveVersionNumber();
  }

  function test_WhenAnActiveVersionExists() external {
    _deployMockRegistryAndAddVersion(true);
    uint32 activeVersionNumber = s_registryManager.getActiveVersionNumber();
    assertEq(activeVersionNumber, 1);
  }
}
