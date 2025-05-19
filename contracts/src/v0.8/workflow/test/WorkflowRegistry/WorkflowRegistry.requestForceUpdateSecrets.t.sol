// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {WorkflowRegistry} from "../../WorkflowRegistry.sol";
import {WorkflowRegistrySetup} from "./WorkflowRegistrySetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract WorkflowRegistry_requestForceUpdateSecrets is WorkflowRegistrySetup {
  function test_RevertWhen_TheRegistryIsLocked() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Lock the registry as the owner.
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Attempt to request force update secrets now after the registry is locked.
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.requestForceUpdateSecrets(s_validSecretsURL);
  }

  // whenTheRegistryIsNotLocked
  function test_RevertWhen_TheCallerDoesNotOwnAnyWorkflowsWithTheSecretsURL() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Call the requestForceUpdateSecrets function now on a random URL
    vm.prank(s_authorizedAddress);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.requestForceUpdateSecrets(s_validBinaryURL);
  }

  // whenTheRegistryIsNotLocked whenTheCallerOwnsWorkflowsWithTheSecretsURL
  function test_WhenTheCallerIsNotAnAuthorizedAddress() external {
    // Register a workflow first.
    _registerValidWorkflow();

    _removeAddressFromAuthorizedAddresses(s_authorizedAddress);

    // Start recording logs
    vm.recordLogs();

    vm.prank(s_authorizedAddress);
    s_registry.requestForceUpdateSecrets(s_validSecretsURL);

    // Retrieve the recorded logs.
    Vm.Log[] memory entries = vm.getRecordedLogs();

    // Event signature hash for WorkflowForceUpdateSecretsRequestedV1.
    bytes32 eventSignature = keccak256("WorkflowForceUpdateSecretsRequestedV1(string,address,string)");

    // Iterate through the logs to ensure WorkflowForceUpdateSecretsRequestedV1 was not emitted.
    bool eventEmitted = false;
    for (uint256 i = 0; i < entries.length; ++i) {
      if (entries[i].topics[0] == eventSignature) {
        eventEmitted = true;
        break;
      }
    }
    // Assert that the event was not emitted
    assertFalse(eventEmitted);
  }

  // whenTheRegistryIsNotLocked whenTheCallerOwnsWorkflowsWithTheSecretsURL
  function test_WhenTheCallerIsAnAuthorizedAddress_AndTheWorkflowIsNotInAnAllowedDON() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Start recording logs
    vm.recordLogs();

    _removeDONFromAllowedDONs(s_allowedDonID);

    // Call the requestForceUpdateSecrets function now after the don is removed.
    vm.prank(s_authorizedAddress);
    s_registry.requestForceUpdateSecrets(s_validSecretsURL);

    // Retrieve the recorded logs
    Vm.Log[] memory entries = vm.getRecordedLogs();

    // Event signature hash for WorkflowForceUpdateSecretsRequestedV1
    bytes32 eventSignature = keccak256("WorkflowForceUpdateSecretsRequestedV1(string,address,string)");

    // Iterate through the logs to ensure WorkflowForceUpdateSecretsRequestedV1 was not emitted
    bool eventEmitted = false;
    for (uint256 i = 0; i < entries.length; ++i) {
      if (entries[i].topics[0] == eventSignature) {
        eventEmitted = true;
        break;
      }
    }
    // Assert that the event was not emitted
    assertFalse(eventEmitted);
  }

  // whenTheRegistryIsNotLocked whenTheCallerOwnsWorkflowsWithTheSecretsURL
  function test_WhenTheCallerIsAnAuthorizedAddress_AndTheWorkflowIsInAnAllowedDON() external {
    // Register a workflow first.
    _registerValidWorkflow();

    // Register another workflow with the same owner but different secrets URL.
    vm.prank(s_authorizedAddress);
    s_registry.registerWorkflow(
      "ValidWorkflow2",
      keccak256("validWorkflow2"),
      s_allowedDonID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      "https://example.com/valid-binary2",
      s_validConfigURL,
      s_validSecretsURL
    );

    // Start recording logs
    vm.recordLogs();

    vm.prank(s_authorizedAddress);
    s_registry.requestForceUpdateSecrets(s_validSecretsURL);
    // Verify the event emitted with correct details
    Vm.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries.length, 2);

    bytes32 eventSignature = keccak256("WorkflowForceUpdateSecretsRequestedV1(address,bytes32,string)");

    // Check the first event
    assertEq(entries[0].topics[0], eventSignature);
    // Verify owner (indexed)
    address decodedAddress = abi.decode(abi.encodePacked(entries[0].topics[1]), (address));
    assertEq(decodedAddress, s_authorizedAddress);
    // Decode non-indexed parameters (secretsURLHash and workflowName)
    (bytes32 decodedSecretsURLHash, string memory decodedWorkflowName) = abi.decode(entries[0].data, (bytes32, string));
    // Verify the decoded values
    bytes32 expectedSecretsURLHash = keccak256(abi.encodePacked(s_authorizedAddress, s_validSecretsURL));
    assertEq(decodedSecretsURLHash, expectedSecretsURLHash);
    assertEq(decodedWorkflowName, s_validWorkflowName);

    // // Check the second event
    assertEq(entries[1].topics[0], eventSignature);
    // Verify owner (indexed)
    address decodedAddress2 = abi.decode(abi.encodePacked(entries[1].topics[1]), (address));
    assertEq(decodedAddress2, s_authorizedAddress);
    // Decode non-indexed parameters (secretsURLHash and workflowName)
    (bytes32 decodedSecretsURLHash2, string memory decodedWorkflowName2) =
      abi.decode(entries[1].data, (bytes32, string));
    // Verify the decoded values
    bytes32 expectedSecretsURLHash2 = keccak256(abi.encodePacked(s_authorizedAddress, s_validSecretsURL));
    assertEq(decodedSecretsURLHash2, expectedSecretsURLHash2);
    assertEq(decodedWorkflowName2, "ValidWorkflow2");
  }
}
