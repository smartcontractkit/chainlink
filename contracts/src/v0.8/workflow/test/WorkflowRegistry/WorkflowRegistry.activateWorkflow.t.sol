// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_activateWorkflow is WorkflowRegistrySetup {
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Lock the registry as the owner.
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Attempt to activate the workflow now after the registry is locked.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.activateWorkflow(s_validWorkflowKey);
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
    s_registry.activateWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheWorkflowIsAlreadyActive() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Attempt to activate the workflow.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyInDesiredStatus.selector);
    s_registry.activateWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsPaused
  function test_RevertWhen_TheDonIDIsNotAllowed() external {
    // Register a paused workflow first.
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

    // Remove the DON from the allowed DONs list.
    _removeDONFromAllowedDONs(s_allowedDonID);

    // Attempt to activate the workflow.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.DONNotAllowed.selector, s_allowedDonID));
    s_registry.activateWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsPaused whenTheDonIDIsAllowed
  function test_RevertWhen_TheCallerIsNotAnAuthorizedAddress() external {
    // Register a paused workflow first.
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

    // Remove the address from the authorized addresses list.
    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);

    // Attempt to activate the workflow.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.AddressNotAuthorized.selector, s_authorizedAddress));
    s_registry.activateWorkflow(s_validWorkflowKey);
  }

  // whenTheRegistryIsNotLocked whenTheCallerIsTheWorkflowOwner whenTheWorkflowIsPaused whenTheDonIDIsAllowed
  function test_WhenTheCallerIsAnAuthorizedAddress() external {
    // Register a paused workflow first.
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

    // It should emit {WorkflowActivatedV1} when the workflow is activated.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowActivatedV1(
      s_validWorkflowID, s_authorizedAddress, s_allowedDonID, s_validWorkflowName
    );

    // Activate the workflow.
    vm.prank(s_authorizedAddress);
    s_registry.activateWorkflow(s_validWorkflowKey);

    // Check that the workflow is active.
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);
  }
}
