// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_pauseWorkflow is WorkflowRegistrySetup {
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Lock the registry as the owner.
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Attempt to pause the workflow now after the registry is locked.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.pauseWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked
  function test_RevertWhen_TheCallerIsNotTheWorkflowOwner() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Attempt to pause the workflow from a different address.
    vm.prank(s_unauthorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.CallerIsNotWorkflowOwner.selector, s_unauthorizedAddress));
    s_registry.pauseWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheWorkflowIsAlreadyPaused() external {
    // Register a paused workflow.
    vm.prank(s_authorizedAddress);
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.PAUSED,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    // Attempt to pause the workflow.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyInDesiredStatus.selector);
    s_registry.pauseWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsActive
  function test_WhenTheDonIDIsNotAllowed_AndTheCallerIsAnAuthorizedAddress() external {
    // Register a workflow first.
    _registerValidWorkflow();

    _removeDONFromAllowedDONs(s_allowedDonID);

    // It should emit {WorkflowPausedV1} when the workflow is paused.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowPausedV1(s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName);

    // Pause the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.pauseWorkflow(s_validWorkflowKey);

    // Check that the workflow is paused.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.PAUSED);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsActive
  function test_WhenTheDonIDIsNotAllowed_AndTheCallerIsAnUnauthorizedAddress() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Remove the allowed DON ID and the authorized address.
    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);
    _removeDONFromAllowedDONs(s_allowedDonID);

    // It should emit {WorkflowPausedV1} when the workflow is paused.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowPausedV1(s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName);

    // Pause the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.pauseWorkflow(s_validWorkflowKey);

    // Check that the workflow is paused.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.PAUSED);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsActive
  function test_WhenTheDonIDIsAllowed_AndTheCallerIsAnUnauthorizedAddress() external {
    // Register a workflow first.
    _registerValidWorkflow();

    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);

    // It should emit {WorkflowPausedV1} when the workflow is paused.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowPausedV1(s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName);

    // Pause the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.pauseWorkflow(s_validWorkflowKey);

    // Check that the workflow is paused.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.PAUSED);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsActive
  function test_WhenTheDonIDIsAllowed_AndTheCallerIsAnAuthorizedAddress() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // It should emit {WorkflowPausedV1} when the workflow is paused.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowPausedV1(s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName);

    // Pause the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.pauseWorkflow(s_validWorkflowKey);

    // Check that the workflow is paused.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.PAUSED);
  }
}
