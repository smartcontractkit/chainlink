// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../dev/WorkflowRegistryManager.sol";
import {Test} from "forge-std/Test.sol";

contract WorkflowRegistryManagerSetup is Test {
  WorkflowRegistryManager internal s_registryManager;
  address internal s_owner;
  address internal s_nonOwner;
  address internal s_contractAddress;
  uint64 internal s_chainID;
  uint32 internal s_versionNumber;
  uint32 internal s_deployedAt;

  function setUp() public virtual {
    s_owner = makeAddr("owner");
    s_nonOwner = makeAddr("nonOwner");
    s_contractAddress = makeAddr("contractAddress");
    s_chainID = 1;
    s_versionNumber = 1;
    s_deployedAt = uint32(block.timestamp);

    // Deploy the WorkflowRegistryManager contract
    vm.prank(s_owner);
    s_registryManager = new WorkflowRegistryManager();
  }
}
