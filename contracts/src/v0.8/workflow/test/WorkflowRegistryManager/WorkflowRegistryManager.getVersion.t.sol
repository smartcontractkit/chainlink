// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";
// import {MockWorkflowRegistryContract} from "../../mocks/MockWorkflowRegistryContract.sol";
import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";

contract WorkflowRegistryManager_getVersion is WorkflowRegistryManagerSetup {
  function test_WhenVersionNumberIsNotRegistered() external {
    // it should revert with VersionNotRegistered
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.VersionNotRegistered.selector, 1));
    s_registryManager.getVersion(1);
  }

  function test_WhenVersionNumberIsRegistered() external {
    // it should return the correct version details
    _deployMockRegistryAndAddVersion(true);
    WorkflowRegistryManager.Version memory version = s_registryManager.getVersion(1);
    assertEq(version.contractAddress, address(s_mockWorkflowRegistryContract));
    assertEq(version.chainID, s_chainID);
    assertEq(version.deployedAt, s_deployedAt);
  }
}
