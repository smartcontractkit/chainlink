// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";
import {MockWorkflowRegistryContract} from "../../mocks/MockWorkflowRegistryContract.sol";
import {Test} from "forge-std/Test.sol";

contract WorkflowRegistryManagerSetup is Test {
  WorkflowRegistryManager internal s_registryManager;
  MockWorkflowRegistryContract internal s_mockWorkflowRegistryContract;
  address internal s_owner;
  address internal s_stranger;
  address internal s_invalidContractAddress;
  uint64 internal s_chainID;
  uint32 internal s_deployedAt;

  function setUp() public virtual {
    s_owner = makeAddr("owner");
    s_stranger = makeAddr("nonOwner");
    s_invalidContractAddress = makeAddr("contractAddress");
    s_chainID = 1;
    s_deployedAt = uint32(block.timestamp);

    // Deploy the WorkflowRegistryManager contract
    vm.prank(s_owner);
    s_registryManager = new WorkflowRegistryManager();
  }

  // Helper function to deploy the MockWorkflowRegistryContract and add it to the WorkflowRegistryManager
  function _deployMockRegistryAndAddVersion(
    bool activate
  ) internal {
    // Deploy the MockWorkflowRegistryContract contract
    s_mockWorkflowRegistryContract = new MockWorkflowRegistryContract();

    // Add the MockWorkflowRegistryContract to the WorkflowRegistryManager
    vm.prank(s_owner);
    s_registryManager.addVersion(address(s_mockWorkflowRegistryContract), s_chainID, s_deployedAt, activate);
  }
}
