// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_getLatestVersion is WorkflowRegistryManagerSetup {
  function test_WhenNoVersionsHaveBeenRegistered() external {
    // it should revert with NoVersionsRegistered
    vm.expectRevert(WorkflowRegistryManager.NoVersionsRegistered.selector);
    s_registryManager.getLatestVersion();
  }

  function test_WhenVersionsHaveBeenRegistered() external {
    // it should return the latest registered version details
    _deployMockRegistryAndAddVersion(true);
    WorkflowRegistryManager.Version memory version = s_registryManager.getLatestVersion();
    assertEq(version.contractAddress, address(s_mockWorkflowRegistryContract));
    assertEq(version.chainID, s_chainID);
    assertEq(version.deployedAt, s_deployedAt);
  }
}
