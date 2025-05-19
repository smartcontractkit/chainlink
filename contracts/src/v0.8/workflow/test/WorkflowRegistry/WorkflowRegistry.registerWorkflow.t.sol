// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";

contract WorkflowRegistry_registerWorkflow is WorkflowRegistrySetup {
  function test_RevertWhen_TheCallerIsNotAnAuthorizedAddress() external {
    vm.prank(s_unauthorizedAddress);

    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.AddressNotAuthorized.selector, s_unauthorizedAddress));
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Lock the registry as the owner
    vm.startPrank(s_owner);
    s_registry.lockRegistry();

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
    vm.stopPrank();
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked
  function test_RevertWhen_TheDonIDIsNotAllowed() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.DONNotAllowed.selector, s_disallowedDonID));
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_disallowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheWorkflowNameIsEmpty() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(WorkflowRegistry.WorkflowNameRequired.selector);
    s_registry.registerWorkflow(
      "",
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheWorkflowNameIsTooLong() external {
    vm.prank(s_authorizedAddress);

    // Ensure the expected error encoding matches the actual error
    vm.expectRevert(
      abi.encodeWithSelector(WorkflowRegistry.WorkflowNameTooLong.selector, bytes(s_invalidWorkflowName).length, 64)
    );
    s_registry.registerWorkflow(
      s_invalidWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheBinaryURLIsEmpty() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(WorkflowRegistry.BinaryURLRequired.selector);
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      "",
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheBinaryURLIsTooLong() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.URLTooLong.selector, bytes(s_invalidURL).length, 200));
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_invalidURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheConfigURLIsTooLong() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.URLTooLong.selector, bytes(s_invalidURL).length, 200));
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_invalidURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheSecretsURLIsTooLong() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistry.URLTooLong.selector, bytes(s_invalidURL).length, 200));
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_invalidURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheWorkflowIDIsInvalid() external {
    vm.prank(s_authorizedAddress);

    vm.expectRevert(WorkflowRegistry.InvalidWorkflowID.selector);
    s_registry.registerWorkflow(
      s_validWorkflowName,
      bytes32(0),
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheWorkflowIDIsAlreadyInUsedByAnotherWorkflow() external {
    vm.startPrank(s_authorizedAddress);

    // Register a valid workflow first
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    vm.expectRevert(WorkflowRegistry.WorkflowIDAlreadyExists.selector);
    s_registry.registerWorkflow(
      "ValidWorkflow2",
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    vm.stopPrank();
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_RevertWhen_TheWorkflowNameIsAlreadyUsedByTheOwner() external {
    vm.startPrank(s_authorizedAddress);

    // Register a valid workflow first
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    // Register the same workflow again
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyRegistered.selector);
    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    vm.stopPrank();
  }

  // whenTheCallerIsAnAuthorizedAddress whenTheRegistryIsNotLocked whenTheDonIDIsAllowed
  function test_WhenTheWorkflowInputsAreAllValid() external {
    vm.startPrank(s_authorizedAddress);

    // it should emit {WorkflowRegisteredV1}
    vm.expectEmit();
    emit WorkflowRegistry.WorkflowRegisteredV1(
      s_validWorkflowID,
      s_authorizedAddress,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validWorkflowName,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    s_registry.registerWorkflow(
      s_validWorkflowName,
      s_validWorkflowID,
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_validBinaryURL,
      s_validConfigURL,
      s_validSecretsURL
    );

    // it should store the new workflow in s_workflows
    WorkflowRegistry.WorkflowMetadata memory workflow =
      s_registry.getWorkflowMetadata(s_authorizedAddress, s_validWorkflowName);
    assertEq(workflow.owner, s_authorizedAddress);
    assertEq(workflow.donID, s_allowedDonID);
    assertEq(workflow.workflowName, s_validWorkflowName);
    assertEq(workflow.workflowID, s_validWorkflowID);
    assertEq(workflow.binaryURL, s_validBinaryURL);
    assertEq(workflow.configURL, s_validConfigURL);
    assertEq(workflow.secretsURL, s_validSecretsURL);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // it should add the workflow key to s_ownerWorkflowKeys
    WorkflowRegistry.WorkflowMetadata[] memory workflows =
      s_registry.getWorkflowMetadataListByOwner(s_authorizedAddress, 0, 1);
    assertEq(workflows[0].owner, s_authorizedAddress);
    assertEq(workflows[0].donID, s_allowedDonID);
    assertEq(workflows[0].workflowName, s_validWorkflowName);
    assertEq(workflows[0].workflowID, s_validWorkflowID);
    assertEq(workflows[0].binaryURL, s_validBinaryURL);
    assertEq(workflows[0].configURL, s_validConfigURL);
    assertEq(workflows[0].secretsURL, s_validSecretsURL);
    assertTrue(workflows[0].status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // it should add the workflow key to s_donWorkflowKeys
    workflows = s_registry.getWorkflowMetadataListByDON(s_allowedDonID, 0, 1);
    assertEq(workflows[0].owner, s_authorizedAddress);
    assertEq(workflows[0].donID, s_allowedDonID);
    assertEq(workflows[0].workflowName, s_validWorkflowName);
    assertEq(workflows[0].workflowID, s_validWorkflowID);
    assertEq(workflows[0].binaryURL, s_validBinaryURL);
    assertEq(workflows[0].configURL, s_validConfigURL);
    assertEq(workflows[0].secretsURL, s_validSecretsURL);
    assertTrue(workflows[0].status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // it should add the url + key to s_secretsHashToWorkflows when the secretsURL is not empty
    vm.expectEmit(true, true, false, true);
    emit WorkflowRegistry.WorkflowForceUpdateSecretsRequestedV1(
      s_authorizedAddress, keccak256(abi.encodePacked(s_authorizedAddress, s_validSecretsURL)), s_validWorkflowName
    );

    // Call the function that should emit the event
    s_registry.requestForceUpdateSecrets(s_validSecretsURL);

    vm.stopPrank();
  }
}
