// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_getLatestVersionNumber is WorkflowRegistryManagerSetup {
  function test_WhenNoVersionsHaveBeenRegistered() external {
    vm.expectRevert(WorkflowRegistryManager.NoVersionsRegistered.selector);
    s_registryManager.getLatestVersionNumber();
  }

  function test_WhenVersionsHaveBeenRegistered() external {
    _deployMockRegistryAndAddVersion(true);
    uint32 latestVersionNumber = s_registryManager.getLatestVersionNumber();
    assertEq(latestVersionNumber, 1);
  }
}
