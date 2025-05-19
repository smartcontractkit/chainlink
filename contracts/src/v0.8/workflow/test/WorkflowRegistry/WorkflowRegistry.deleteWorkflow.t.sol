// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_deleteWorkflow is WorkflowRegistrySetup {
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Lock the registry as the owner.
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Attempt to delete the workflow now after the registry is locked.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.deleteWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked
  function test_RevertWhen_TheCallerIsNotTheWorkflowOwner() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Add the previously unauthorized address to the authorized addresses list.
    _addAddressToAuthorizedAddresses(s_unauthorizedAddress);

    // Update the workflow now as the new authorized user.
    vm.prank(s_unauthorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.CallerIsNotWorkflowOwner.selector, s_unauthorizedAddress));
    s_registry.deleteWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheCallerIsNotAnAuthorizedAddress() external {
    // Register the workflow first as an authorized address.
    _registerValidWorkflow();

    // Remove the address from the authorized addresses list.
    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);

    // Delete the workflow now after the workflow owner is no longer an authorized address.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.AddressNotAuthorized.selector, s_authorizedAddress));
    s_registry.deleteWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner
  function test_WhenTheCallerIsAnAuthorizedAddress_AndTheDonIDIsAllowed() external {
    // Register the workflow.
    _registerValidWorkflow();

    // Check that the workflow exists.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertEq(workflow.workflowName, s_validWorkflowName);

    // It should emit {WorkflowDeletedV1} when the workflow is deleted.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowDeletedV1(s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName);

    // Delete the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.deleteWorkflow(s_validWorkflowKey);

    // Check that the workflow was deleted.
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner
  function test_WhenTheCallerIsAnAuthorizedAddress_AndTheDonIDIsNotAllowed() external {
    // Register the workflow.
    _registerValidWorkflow();

    // Check that the workflow exists.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertEq(workflow.workflowName, s_validWorkflowName);

    // Remove the DON from the allowed DONs list.
    _removeDONFromAllowedDONs(s_allowedDonID);

    // It should emit {WorkflowDeletedV1} when the workflow is deleted.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowDeletedV1(s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName);

    // Delete the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.deleteWorkflow(s_validWorkflowKey);

    // Check that the workflow was deleted.
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
  }
}
