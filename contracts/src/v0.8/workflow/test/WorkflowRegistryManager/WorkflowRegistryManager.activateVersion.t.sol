// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract WorkflowRegistryManager_activateVersion is WorkflowRegistryManagerSetup {
  function test_RevertWhen_TheCallerIsNotTheOwner() external {
    // it should revert
    vm.prank(s_stranger);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registryManager.activateVersion(2);
  }

  // whenTheCallerIsTheOwner
  function test_RevertWhen_TheVersionNumberDoesNotExist() external {
    // it should revert
    vm.prank(s_owner);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.VersionNotRegistered.selector, 5));
    s_registryManager.activateVersion(5);
  }

  // whenTheCallerIsTheOwner whenTheVersionNumberExists
  function test_RevertWhen_TheVersionNumberIsAlreadyActive() external {
    // Deploy a mock registry and add that to the registry manager
    _deployMockRegistryAndAddVersion(true);

    // Get the latest version number
    uint32 versionNumber = s_registryManager.getLatestVersionNumber();

    // Activate the same version
    vm.prank(s_owner);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.VersionAlreadyActive.selector, versionNumber));
    s_registryManager.activateVersion(versionNumber);
  }

  function test_WhenTheVersionNumberIsNotActive_AndWhenThereAreNoActiveVersions() external {
    // Deploy a mock registry and add but not activate it.
    _deployMockRegistryAndAddVersion(false);

    // Get the latest version number, which should be 1.
    uint32 versionNumber = s_registryManager.getLatestVersionNumber();
    assertEq(versionNumber, 1);

    // Start recording logs to check that VersionDeactivated is not emitted.
    vm.recordLogs();

    // Since there are no existing active versions, this should only activate the version and emit VersionActivated
    vm.expectEmit(true, true, false, true);
    emit WorkflowRegistryManager.VersionActivated(address(s_mockWorkflowRegistryContract), s_chainID, versionNumber);

    // Activate the version
    vm.prank(s_owner);
    s_registryManager.activateVersion(versionNumber);

    // Retrieve the recorded logs.
    Vm.Log[] memory entries = vm.getRecordedLogs();

    // Event signature hash for WorkflowForceUpdateSecretsRequestedV1.
    bytes32 eventSignature = keccak256("VersionDeactivated(string,address,string)");

    // Iterate through the logs to ensure VersionDeactivatedV1 was not emitted.
    bool deactivateEventEmitted = false;
    for (uint256 i = 0; i < entries.length; ++i) {
      if (entries[i].topics[0] == eventSignature) {
        deactivateEventEmitted = true;
        break;
      }
    }

    // Assert that the event was not emitted.
    assertFalse(deactivateEventEmitted);

    // Deploy another mock registry.
    _deployMockRegistryAndAddVersion(false);

    // Get the latest version number, which should now be 2.
    uint32 newVersionNumber = s_registryManager.getLatestVersionNumber();
    assertEq(newVersionNumber, 2);

    // It should now emit both VersionActivated and VersionDeactivated.
    vm.expectEmit(true, true, false, true);
    emit WorkflowRegistryManager.VersionActivated(address(s_mockWorkflowRegistryContract), s_chainID, newVersionNumber);
    emit WorkflowRegistryManager.VersionDeactivated(address(s_mockWorkflowRegistryContract), s_chainID, versionNumber);

    // Activate the version
    vm.prank(s_owner);
    s_registryManager.activateVersion(newVersionNumber);
  }
}
