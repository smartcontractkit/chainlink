// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {MockContract} from "../../mocks/MockContract.sol";
import {MockWorkflowRegistryContract} from "../../mocks/MockWorkflowRegistryContract.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_addVersion is WorkflowRegistryManagerSetup {
  function test_RevertWhen_TheCallerIsNotTheOwner() external {
    // Deploy the MockWorkflowRegistryContract contract
    vm.prank(s_owner);
    MockWorkflowRegistryContract mockContract = new MockWorkflowRegistryContract();

    // Add it as a non owner
    vm.prank(s_stranger);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registryManager.addVersion(address(mockContract), s_chainID, s_deployedAt, true);
  }

  // whenTheCallerIsTheOwner
  function test_RevertWhen_TheContractAddressIsInvalid() external {
    // Add a 0 address to the WorkflowRegistryManager
    vm.prank(s_owner);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.InvalidContractAddress.selector, address(0)));
    s_registryManager.addVersion(address(0), s_chainID, s_deployedAt, true);
  }

  // whenTheCallerIsTheOwner whenTheContractAddressIsValid
  function test_RevertWhen_TheContractIsAlreadyRegistered() external {
    // Deploy a MockWorkflowRegistryContract contract
    MockWorkflowRegistryContract mockWfrContract = new MockWorkflowRegistryContract();

    // Add it to the WorkflowRegistryManager
    vm.prank(s_owner);
    s_registryManager.addVersion(address(mockWfrContract), s_chainID, s_deployedAt, true);

    // Try to add it again
    vm.prank(s_owner);
    vm.expectRevert(
      abi.encodeWithSelector(
        WorkflowRegistryManager.ContractAlreadyRegistered.selector, address(mockWfrContract), s_chainID
      )
    );
    s_registryManager.addVersion(address(mockWfrContract), s_chainID, s_deployedAt, true);
  }

  // whenTheCallerIsTheOwner whenTheContractAddressIsValid
  function test_RevertWhen_TheContractTypeIsInvalid() external {
    // Deploy a random contract
    MockContract mockContract = new MockContract();

    // Add it to the WorkflowRegistryManager
    vm.prank(s_owner);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.InvalidContractType.selector, address(mockContract)));
    s_registryManager.addVersion(address(mockContract), s_chainID, s_deployedAt, true);
  }

  // whenTheCallerIsTheOwner whenTheContractAddressIsValid whenTheContractTypeIsValid
  function test_WhenAutoActivateIsTrue() external {
    // Get the latest version number, which should revert.
    vm.expectRevert(WorkflowRegistryManager.NoVersionsRegistered.selector);
    s_registryManager.getLatestVersionNumber();

    // Deploy a MockWorkflowRegistryContract contract
    MockWorkflowRegistryContract mockWfrContract = new MockWorkflowRegistryContract();

    // Expect both VersionAdded and VersionActivated events to be emitted.
    vm.expectEmit(true, true, false, true);
    emit WorkflowRegistryManager.VersionAdded(address(mockWfrContract), s_chainID, s_deployedAt, 1);
    emit WorkflowRegistryManager.VersionActivated(address(mockWfrContract), s_chainID, 1);

    // Add the MockWorkflowRegistryContract to the WorkflowRegistryManager
    vm.prank(s_owner);
    s_registryManager.addVersion(address(mockWfrContract), s_chainID, s_deployedAt, true);

    // Get the latest version number again, which should be 1.
    uint32 versionNumber = s_registryManager.getLatestVersionNumber();
    assertEq(versionNumber, 1);

    // Get the latest active version number, which should also be 1.
    uint32 activeVersionNumber = s_registryManager.getActiveVersionNumber();
    assertEq(activeVersionNumber, 1);
  }

  // whenTheCallerIsTheOwner whenTheContractAddressIsValid whenTheContractTypeIsValid
  function test_WhenAutoActivateIsFalse() external {
    // Get the latest version number, which should revert.
    vm.expectRevert(WorkflowRegistryManager.NoVersionsRegistered.selector);
    s_registryManager.getLatestVersionNumber();

    // Deploy a MockWorkflowRegistryContract contract
    MockWorkflowRegistryContract mockWfrContract = new MockWorkflowRegistryContract();

    vm.expectEmit(true, true, false, true);
    emit WorkflowRegistryManager.VersionAdded(address(mockWfrContract), s_chainID, s_deployedAt, 1);

    // Add the MockWorkflowRegistryContract to the WorkflowRegistryManager
    vm.prank(s_owner);
    s_registryManager.addVersion(address(mockWfrContract), s_chainID, s_deployedAt, false);

    // Get the latest version number again, which should be 1.
    uint32 versionNumber = s_registryManager.getLatestVersionNumber();
    assertEq(versionNumber, 1);

    // Get the latest active version number, which should revert.
    vm.expectRevert(WorkflowRegistryManager.NoActiveVersionAvailable.selector);
    s_registryManager.getActiveVersionNumber();
  }
}
