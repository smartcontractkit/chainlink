// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {WorkflowRegistryManager} from "../../WorkflowRegistryManager.sol";

import {MockWorkflowRegistryContract} from "../../mocks/MockWorkflowRegistryContract.sol";
import {WorkflowRegistryManagerSetup} from "./WorkflowRegistryManagerSetup.t.sol";

contract WorkflowRegistryManager_getAllVersions is WorkflowRegistryManagerSetup {
  MockWorkflowRegistryContract internal s_mockContract1;
  MockWorkflowRegistryContract internal s_mockContract2;
  MockWorkflowRegistryContract internal s_mockContract3;

  function setUp() public override {
    super.setUp();
    // Add 3 versions
    s_mockContract1 = new MockWorkflowRegistryContract();
    s_mockContract2 = new MockWorkflowRegistryContract();
    s_mockContract3 = new MockWorkflowRegistryContract();

    vm.startPrank(s_owner);
    s_registryManager.addVersion(address(s_mockContract1), s_chainID, s_deployedAt, true);
    s_registryManager.addVersion(address(s_mockContract2), s_chainID, s_deployedAt, false);
    s_registryManager.addVersion(address(s_mockContract3), s_chainID, s_deployedAt, true);
    vm.stopPrank();
  }

  function test_WhenRequestingWithInvalidStartIndex() external view {
    // It should return an empty array.
    WorkflowRegistryManager.Version[] memory versions = s_registryManager.getAllVersions(10, 1);
    assertEq(versions.length, 0);
  }

  function test_WhenRequestingWithValidStartIndexAndLimitWithinBounds() external view {
    // It should return the correct versions based on pagination.
    WorkflowRegistryManager.Version[] memory versions = s_registryManager.getAllVersions(1, 2);
    assertEq(versions.length, 2);
    assertEq(versions[0].contractAddress, address(s_mockContract1));
    assertEq(versions[1].contractAddress, address(s_mockContract2));
  }

  function test_WhenLimitExceedsMaximumPaginationLimit() external view {
    // it should return results up to MAX_PAGINATION_LIMIT
    WorkflowRegistryManager.Version[] memory versions = s_registryManager.getAllVersions(1, 200);
    assertEq(versions.length, 3);
    assertEq(versions[0].contractAddress, address(s_mockContract1));
    assertEq(versions[1].contractAddress, address(s_mockContract2));
    assertEq(versions[2].contractAddress, address(s_mockContract3));
  }
}
