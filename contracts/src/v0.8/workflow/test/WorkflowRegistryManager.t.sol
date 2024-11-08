// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../shared/access/Ownable2Step.sol";
import {WorkflowRegistry} from "../dev/WorkflowRegistry.sol";
import {WorkflowRegistryManager} from "../dev/WorkflowRegistryManager.sol";
import {Test} from "forge-std/Test.sol";

contract WorkflowRegistryManagerTest is Test {
  WorkflowRegistryManager private s_manager;
  WorkflowRegistry private s_registry;

  address private s_owner = address(1);
  address private s_unauthorizedUser = address(2);
  uint32 private s_chainID = 1;

  function setUp() public {
    vm.prank(s_owner);
    s_manager = new WorkflowRegistryManager();

    vm.prank(s_owner);
    s_registry = new WorkflowRegistry();
  }

  function testAddVersionFailsForUnauthorizedUser() public {
    vm.prank(s_unauthorizedUser);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_manager.addVersion(address(s_registry), s_chainID, false);
  }

  function testAddVersion() public {
    vm.prank(s_owner);
    s_manager.addVersion(address(s_registry), s_chainID, true);

    // Get first version using getAllVersions
    WorkflowRegistryManager.Version[] memory versions = s_manager.getAllVersions(1, 1);
    assertEq(versions[0].contractAddress, address(s_registry));
    assertEq(versions[0].chainID, s_chainID);
    assertEq(versions[0].deployedAt, block.timestamp);
  }

  function testActivateVersion() public {
    vm.startPrank(s_owner);

    // Add two versions
    s_manager.addVersion(address(s_registry), s_chainID, false);
    WorkflowRegistry newRegistry = new WorkflowRegistry();
    s_manager.addVersion(address(newRegistry), s_chainID, false);

    // Activate first version
    s_manager.activateVersion(1);
    vm.stopPrank();

    // Verify first version is active
    WorkflowRegistryManager.Version memory activeVersion = s_manager.getActiveVersion();
    assertEq(activeVersion.contractAddress, address(s_registry));
  }

  function testAddVersionFailsForInvalidContract() public {
    // Invalid version number
    InvalidContract invalidContract = new InvalidContract();
    address invalidContractAddress = address(invalidContract);
    vm.prank(s_owner);
    vm.expectRevert(
      abi.encodeWithSelector(WorkflowRegistryManager.InvalidContractType.selector, invalidContractAddress)
    );
    s_manager.addVersion(invalidContractAddress, s_chainID, false);

    // Zero address
    vm.prank(s_owner);
    vm.expectRevert(abi.encodeWithSelector(WorkflowRegistryManager.InvalidContractAddress.selector, address(0)));
    s_manager.addVersion(address(0), s_chainID, false);
  }

  function testGetAllVersionsPagination() public {
    // Add multiple versions
    vm.startPrank(s_owner);
    for (uint32 i = 0; i < 5; i++) {
      s_manager.addVersion(address(s_registry), s_chainID, false);
    }
    vm.stopPrank();

    // Test with valid start and limit
    WorkflowRegistryManager.Version[] memory versions = s_manager.getAllVersions(1, 3);
    assertEq(versions.length, 3);

    // Test with start index beyond the total versions
    versions = s_manager.getAllVersions(10, 3);
    assertEq(versions.length, 0);

    // Test with limit exceeding the total versions
    versions = s_manager.getAllVersions(1, 10);
    assertEq(versions.length, 5);

    // Test with start and limit that exactly match the total versions
    versions = s_manager.getAllVersions(1, 5);
    assertEq(versions.length, 5);

    // Test with start index at the last version
    versions = s_manager.getAllVersions(4, 1);
    assertEq(versions.length, 1);
  }

  function testTypeAndVersion() public view {
    string memory version = s_manager.typeAndVersion();
    assertEq(version, "WorkflowRegistryManager 1.0.0");
  }
}

contract InvalidContract {}
