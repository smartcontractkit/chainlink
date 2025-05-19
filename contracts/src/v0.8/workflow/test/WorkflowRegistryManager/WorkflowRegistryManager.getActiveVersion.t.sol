// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_getActiveVersion is WorkflowRegistryManagerSetup {
  function test_WhenNoActiveVersionIsAvailable() external {
    vm.expectRevert(WorkflowRegistryManager.NoActiveVersionAvailable.selector);
    s_registryManager.getActiveVersion();
  }

  function test_WhenAnActiveVersionExists() external {
    _deployMockRegistryAndAddVersion(true);
    WorkflowRegistryManager.Version memory activeVersion = s_registryManager.getActiveVersion();
    assertEq(activeVersion.contractAddress, address(s_mockWorkflowRegistryContract));
    assertEq(activeVersion.chainID, s_chainID);
    assertEq(activeVersion.deployedAt, s_deployedAt);
  }
}
