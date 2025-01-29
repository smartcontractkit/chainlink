// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_updateWorkflow is WorkflowRegistrySetup {
  bytes32 private s_newValidWorkflowID = keccak256("newValidWorkflowID");
  string private s_newValidSecretsURL = "https://example.com/new-secrets";
  string private s_newValidConfigURL = "https://example.com/new-config";
  string private s_newValidBinaryURL = "https://example.com/new-binary";

  function test_RevertWhen_TheCallerIsNotAnAuthorizedAddress() external {
    // Register the workflow first as an authorized address.
    _registerValidWorkflow();

    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);

    // Update the workflow now after the workflow owner is no longer an authorized address.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.AddressNotAuthorized.selector, s_authorizedAddress));
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Lock the registry as the owner.
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Update the workflow now after the registry is locked.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked
  function test_RevertWhen_TheDonIDIsNotAllowed() external {
    // Register a workflow first.
    _registerValidWorkflow();

    _removeDONFromAllowedDONs(s_allowedDonID);

    // Update the workflow now after the DON is no longer allowed.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.DONNotAllowed.selector, s_allowedDonID));
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheCallerIsNotTheWorkflowOwner() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Add the previously unauthorized address to the authorized addresses list.
    _addAddressToAuthorizedAddresses(s_unauthorizedAddress);

    // Update the workflow now as the new authorized user.
    vm.prank(s_unauthorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.CallerIsNotWorkflowOwner.selector, s_unauthorizedAddress));
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_AnExistingWorkflowIsNotFoundWithTheGivenWorkflowName() external {
    // Update a workflow with a non-existent workflow name
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.updateWorkflow(
      "nonExistentWorkflow", s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_NoneOfTheURLsAreUpdated() external {
    // Register a workflow first
    _registerValidWorkflow();

    // Update the workflow with no changes to any URLs
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.WorkflowContentNotUpdated.selector);
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheBinaryURLIsEmpty() external {
    // Register a workflow first
    _registerValidWorkflow();

    // Update the workflow with a binary URL that is empty
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.BinaryURLRequired.selector);
    s_registry.updateWorkflow(s_validWorkflowKey, s_newValidWorkflowID, "", s_validConfigURL, s_validSecretsURL);
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheBinaryURLIsTooLong() external {
    // Register a workflow first
    _registerValidWorkflow();

    // Update the workflow with a binary URL that is too long
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.URLTooLong.selector, bytes(s_invalidURL).length, 200));
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_invalidURL, s_validConfigURL, s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheConfigURLIsTooLong() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Update the workflow with a config URL that is too long.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.URLTooLong.selector, bytes(s_invalidURL).length, 200));
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_invalidURL, s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheSecretsURLIsTooLong() external {
    // Register a workflow first
    _registerValidWorkflow();

    // Update the workflow with a secrets URL that is too long.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.URLTooLong.selector, bytes(s_invalidURL).length, 200));
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_invalidURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheWorkflowIDIsInvalid() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Update the workflow with an invalid workflow ID.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.InvalidWorkflowID.selector);
    s_registry.updateWorkflow(s_validWorkflowKey, bytes32(0), s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL);
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_RevertWhen_TheWorkflowIDIsAlreadyInUsedByAnotherWorkflow() external {
    // Register a workflow first
    _registerValidWorkflow();

    // Register another workflow with another workflow ID
    vm.startPrank(s_authorizedAddress);
    s_registry.registerWorkflow(
      "ValidWorkflow2",
      s_newValidWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    // Update the workflow with a workflow ID that is already in use by another workflow.
    vm.expectRevert(WorkflowRegistry.WorkflowIDAlreadyExists.selector);
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_validBinaryURL, s_validConfigURL, s_newValidSecretsURL
    );

    vm.stopPrank();
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed whenTheCallerIsTheWorkflowOwner
  function test_WhenTheWorkflowInputsAreAllValid() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Update the workflow.
    // It should emit {WorkflowUpdatedV1}.
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowUpdatedV1(
      s_validWorkflowID,
      s_authorizedAddress,
      s_allowedDonID,
      s_newValidWorkflowID,
      s_validWorkflowName,
      s_newValidBinaryURL,
      s_newValidConfigURL,
      s_newValidSecretsURL
    );

    vm.startPrank(s_authorizedAddress);
    s_registry.updateWorkflow(
      s_validWorkflowKey, s_newValidWorkflowID, s_newValidBinaryURL, s_newValidConfigURL, s_newValidSecretsURL
    );

    // It should update the workflow in s_workflows with the new values
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertEq(workflow.owner, s_authorizedAddress);
    assertEq(workflow.donID, s_allowedDonID);
    assertEq(workflow.workflowName, s_validWorkflowName);
    assertEq(workflow.workflowID, s_newValidWorkflowID);
    assertEq(workflow.binaryURL, s_newValidBinaryURL);
    assertEq(workflow.configURL, s_newValidConfigURL);
    assertEq(workflow.secretsURL, s_newValidSecretsURL);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // It should add the url + key to s_secretsHashToWorkflows when the secretsURL is not empty
    vm.expectEmit(true, true, false, true);
    emit WorkflowRegistry.WorkflowForceUpdateSecretsRequestedV1(
      s_authorizedAddress, keccak256(abi.encodePacked(s_authorizedAddress, s_newValidSecretsURL)), s_validWorkflowName
    );

    // Call the function that should emit the event.
    s_registry.requestForceUpdateSecrets(s_newValidSecretsURL);
    vm.stopPrank();
  }
}
