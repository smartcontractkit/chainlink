// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_getWorkflowMetadata is WorkflowRegistrySetup {
  function test_WhenTheWorkflowExistsWithTheOwnerAndName() external {
    _registerValidWorkflow();

    WorkflowRegistry.WorkflowMetadata memory metadata =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);

    assertEq(metadata.workflowName, s_validWorkflowName);
    assertEq(metadata.workflowID, s_validWorkflowID);
    assertEq(metadata.binaryURL, s_validBinaryURL);
    assertEq(metadata.configURL, s_validConfigURL);
    assertEq(metadata.secretsURL, s_validSecretsURL);
  }

  function test_WhenTheWorkflowDoesNotExist() external {
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.getWorkflowMetadata(s_authorizedAddress, "RandomWorkflowName");
  }

  function test_WhenTheRegistryIsLocked() external {
    // Register a workflow
    _registerValidWorkflow();

    // Lock the registry
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // It should behave the same as when the registry is not locked
    vm.prank(s_stranger);
    WorkflowRegistry.WorkflowMetadata memory metadata =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);

    assertEq(metadata.workflowName, s_validWorkflowName);
    assertEq(metadata.workflowID, s_validWorkflowID);
    assertEq(metadata.binaryURL, s_validBinaryURL);
    assertEq(metadata.configURL, s_validConfigURL);
    assertEq(metadata.secretsURL, s_validSecretsURL);
  }
}
