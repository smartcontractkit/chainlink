// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../shared/access/Ownable2Step.sol";
import {DONAccessControl} from "../dev/DONAccessControl.sol";
import {WorkflowRegistry} from "../dev/WorkflowRegistry.sol";
import {Test} from "forge-std/Test.sol";
import {Vm} from "forge-std/Vm.sol";

contract WorkflowRegistryTest is Test {
  WorkflowRegistry private registry;

  address private owner = address(1);
  address private unauthorizedUser = address(2);
  address private authorizedUser = address(3);
  bytes32 private workflowID1 = keccak256(abi.encodePacked("workflow1"));
  string private workflowName1 = "workflowName1";
  bytes32 private workflowID2 = keccak256(abi.encodePacked("workflow2"));
  string private workflowName2 = "workflowName2";
  bytes32 private newWorkflowID = keccak256(abi.encodePacked("workflow_new"));
  uint32 private donID = 1;
  string private testBinaryURL = "binaryURL";
  string private testConfigURL = "configURL";
  string private testSecretsURL = "secretsURL";

  function setUp() public {
    vm.prank(owner);
    registry = new WorkflowRegistry();
  }

  function _allowAccessAndRegisterWorkflow(
    address workflowOwner,
    string memory workflowName,
    bytes32 workflowID,
    WorkflowRegistry.WorkflowStatus initialStatus
  ) internal {
    _setupAuthorizedUser(workflowOwner);
    _setupAllowedDON(donID);

    // authorized user registers workflow
    vm.prank(workflowOwner);
    registry.registerWorkflow(
      workflowName, workflowID, donID, initialStatus, testBinaryURL, testConfigURL, testSecretsURL
    );
  }

  function _setupAuthorizedUser(address workflowOwner) internal {
    // owner adds a single authorized address capable of registering workflows
    address[] memory authorizedUsers = new address[](1);
    authorizedUsers[0] = workflowOwner;
    vm.prank(owner);
    registry.updateDONPermissions(donID, authorizedUsers, true);
  }

  function _setupAllowedDON(uint32 _donID) internal {
    // owner adds a single DON ID allowed for registering workflows
    uint32[] memory allowedDONs = new uint32[](1);
    allowedDONs[0] = _donID;
    vm.prank(owner);
    registry.updateAllowedDONs(allowedDONs, true);
  }

  function testLockRegistry() public {
    // Ensure only the owner can lock the registry
    vm.prank(authorizedUser);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    registry.lockRegistry();

    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    // Lock the registry as the owner
    vm.prank(owner);
    registry.lockRegistry();

    // Test all state-changing functions revert when registry is locked
    vm.startPrank(authorizedUser);

    // Test registerWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    registry.registerWorkflow(
      workflowName2,
      workflowID2,
      donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      testBinaryURL,
      testConfigURL,
      testSecretsURL
    );

    // Test updateWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    registry.updateWorkflow(workflowName1, newWorkflowID, "newBinaryURL", "newConfigURL", "newSecretsURL");

    // Test pauseWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    registry.pauseWorkflow(workflowName1);

    // Test activateWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    registry.activateWorkflow(workflowName1);

    // Test deleteWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    registry.deleteWorkflow(workflowName1);

    // Test requestForceUpdateSecrets
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    registry.requestForceUpdateSecrets(testSecretsURL);

    vm.stopPrank();

    // Test owner functions still revert when registry is locked
    vm.startPrank(owner);

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    address[] memory addresses = new address[](1);
    addresses[0] = authorizedUser;
    registry.updateDONPermissions(donID, addresses, true);

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    uint32[] memory dons = new uint32[](1);
    dons[0] = donID;
    registry.updateAllowedDONs(dons, true);

    vm.stopPrank();
  }

  function testUnlockRegistry() public {
    // Check that the registry is initially unlocked
    bool isLocked = registry.isRegistryLocked();
    assertFalse(isLocked, "Registry should start off as unlocked");

    // Lock the registry first
    vm.prank(owner);
    registry.lockRegistry();

    // Ensure only the owner can unlock the registry
    vm.prank(authorizedUser);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    registry.unlockRegistry();

    // Unlock the registry as the owner
    vm.prank(owner);
    registry.unlockRegistry();

    // Perform an action that requires the registry to be unlocked
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    // Verify the workflow was registered successfully
    WorkflowRegistry.WorkflowMetadata memory workflow = registry.getWorkflowMetadata(authorizedUser, workflowName1);
    assertEq(workflow.workflowID, workflowID1);
  }

  function testRegisterWorkflowFailsForNotAuthorizedAddressOrForNotAllowedDONId() public {
    // owner of the contract is not allowed to register workflows without first setting permissions
    vm.prank(owner);
    vm.expectRevert(abi.encodeWithSelector(DONAccessControl.DONNotAllowed.selector, donID));
    registry.registerWorkflow(
      workflowName1,
      workflowID1,
      donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      testBinaryURL,
      testConfigURL,
      testSecretsURL
    );

    // owner adds a single authorized address capable of registering workflows
    address[] memory authorizedUsers = new address[](1);
    authorizedUsers[0] = authorizedUser;
    vm.prank(owner);
    registry.updateDONPermissions(donID, authorizedUsers, true);

    // authorized address is still not able to register because DON ID is not allowed
    vm.prank(authorizedUser);
    vm.expectRevert(abi.encodeWithSelector(DONAccessControl.DONNotAllowed.selector, donID));
    registry.registerWorkflow(
      workflowName1,
      workflowID1,
      donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      testBinaryURL,
      testConfigURL,
      testSecretsURL
    );

    // owner adds a single DON ID allowed for registering workflows
    uint32[] memory allowedDONs = new uint32[](1);
    allowedDONs[0] = donID;
    vm.prank(owner);
    registry.updateAllowedDONs(allowedDONs, true);

    // authorized address is finally able to register workflow
    vm.prank(authorizedUser);
    registry.registerWorkflow(
      workflowName1,
      workflowID1,
      donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      testBinaryURL,
      testConfigURL,
      testSecretsURL
    );

    // sanity check by retrieving the workflow metadata
    WorkflowRegistry.WorkflowMetadata memory workflow = registry.getWorkflowMetadata(authorizedUser, workflowName1);
    assertEq(workflow.workflowID, workflowID1);
    assertEq(workflow.workflowName, workflowName1);
    assertEq(workflow.owner, authorizedUser);
    assertEq(workflow.binaryURL, testBinaryURL);
    assertEq(workflow.configURL, testConfigURL);
    assertEq(workflow.secretsURL, testSecretsURL);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // any other unauthorized address still gets the unauthorized error
    vm.prank(unauthorizedUser);
    vm.expectRevert(abi.encodeWithSelector(DONAccessControl.AddressNotAuthorized.selector, donID, unauthorizedUser));
    registry.registerWorkflow(
      workflowName1,
      workflowID1,
      donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      testBinaryURL,
      testConfigURL,
      testSecretsURL
    );
  }

  function testUpdateWorkflow() public {
    // create a new workflow
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    // authorized user tries to update the workflow by using the same workflow ID as before
    vm.prank(authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowIDNotUpdated.selector);
    registry.updateWorkflow(workflowName1, workflowID1, "newBinaryURL", "newConfigURL", "newSecretsURL");

    // now the authorizer user sets the new workflow ID
    vm.prank(authorizedUser);
    registry.updateWorkflow(workflowName1, newWorkflowID, "newBinaryURL", "newConfigURL", "newSecretsURL");

    // sanity check by retrieving the workflow metadata to make sure parameters are updated
    WorkflowRegistry.WorkflowMetadata memory workflow = registry.getWorkflowMetadata(authorizedUser, workflowName1);
    assertEq(workflow.workflowID, newWorkflowID);
    assertEq(workflow.workflowName, workflowName1);
    assertEq(workflow.owner, authorizedUser);
    assertEq(workflow.binaryURL, "newBinaryURL");
    assertEq(workflow.configURL, "newConfigURL");
    assertEq(workflow.secretsURL, "newSecretsURL");
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);
  }

  function testPauseWorkflow() public {
    // create a new workflow
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    // authorized user pauses the workflow
    vm.prank(authorizedUser);
    registry.pauseWorkflow(workflowName1);

    // sanity check the workflow status update
    WorkflowRegistry.WorkflowMetadata memory workflow = registry.getWorkflowMetadata(authorizedUser, workflowName1);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.PAUSED);

    // authorized user is not able to pause the workflow twice in a row
    vm.prank(authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyInDesiredStatus.selector);
    registry.pauseWorkflow(workflowName1);
  }

  function testActivateWorkflow() public {
    // create a new workflow
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.PAUSED);

    // authorized user activates the workflow
    vm.prank(authorizedUser);
    registry.activateWorkflow(workflowName1);

    // sanity check the workflow status update
    WorkflowRegistry.WorkflowMetadata memory workflow = registry.getWorkflowMetadata(authorizedUser, workflowName1);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // authorized user is not able to activate the workflow twice in a row
    vm.prank(authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyInDesiredStatus.selector);
    registry.activateWorkflow(workflowName1);
  }

  function testNonWorkflowOwnerUserCannotUpdateWorkflow() public {
    // create a new workflow for one authorized user
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    // add a new authorized user capable of registering workflows
    address anotherAuthorizedUser = address(567);
    address[] memory authorizedUsers = new address[](1);
    authorizedUsers[0] = anotherAuthorizedUser;
    vm.prank(owner);
    registry.updateDONPermissions(donID, authorizedUsers, true);

    // new authorized user is not able to update another user's workflow (same workflow name)
    vm.prank(anotherAuthorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    registry.updateWorkflow(workflowName1, newWorkflowID, "newBinaryURL", "newConfigURL", "newSecretsURL");
  }

  function testRequestForceUpdateSecrets() public {
    // Register two workflows with the same secretsURL for the authorized user
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    vm.prank(authorizedUser);
    registry.registerWorkflow(
      workflowName2,
      workflowID2,
      donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      testBinaryURL,
      testConfigURL,
      testSecretsURL
    );

    // Attempt force update secrets from an unauthorized user
    vm.prank(unauthorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    registry.requestForceUpdateSecrets(testSecretsURL);

    // Start recording the logs to later check the event content
    vm.recordLogs();

    // Authorized user requests force update secrets
    vm.prank(authorizedUser);
    registry.requestForceUpdateSecrets(testSecretsURL);

    // Verify the events emitted with correct details
    Vm.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries.length, 2); // Expecting two separate events for each individual workflow

    for (uint256 i = 0; i < entries.length; i++) {
      assertEq(entries[i].topics[0], keccak256("WorkflowForceUpdateSecretsRequestedV1(string,address,string)"));

      // Compare the hash of the expected string with the topic
      bytes32 expectedSecretsURLHash = keccak256(abi.encodePacked(testSecretsURL));
      assertEq(entries[i].topics[1], expectedSecretsURLHash);

      // Decode the indexed address
      address decodedOwner = abi.decode(abi.encodePacked(entries[i].topics[2]), (address));
      assertEq(decodedOwner, authorizedUser);

      // Decode the non-indexed workflow name
      string memory decodedWorkflowName = abi.decode(entries[i].data, (string));

      // Assert the values
      if (i == 0) {
        assertEq(decodedWorkflowName, workflowName1);
      } else {
        assertEq(decodedWorkflowName, workflowName2);
      }
    }
  }

  function testDeleteWorkflow() public {
    // Create a new workflow
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);

    // Unauthorized user should not be able to delete the workflow
    vm.prank(unauthorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    registry.deleteWorkflow(workflowName1);

    // Authorized user deletes the workflow
    vm.prank(authorizedUser);
    registry.deleteWorkflow(workflowName1);

    // Sanity check to verify that the workflow has been deleted
    vm.prank(authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    registry.getWorkflowMetadata(authorizedUser, workflowName1);

    // Authorized user should not be able to delete a non-existing workflow
    vm.prank(authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    registry.deleteWorkflow(workflowName1);
  }

  function testGetAllAllowedDONs() public {
    // Add allowed DON IDs
    uint32[] memory allowedDONs = new uint32[](3);
    allowedDONs[0] = 1;
    allowedDONs[1] = 2;
    allowedDONs[2] = 3;
    vm.prank(owner);

    for (uint32 i = 0; i < allowedDONs.length; i++) {
      vm.expectEmit(true, true, false, false);
      emit DONAccessControl.AllowedDONUpdatedV1(allowedDONs[i], true);
    }

    registry.updateAllowedDONs(allowedDONs, true);

    // Verify the allowed DONs list
    uint32[] memory fetchedDONs = registry.getAllAllowedDONs();
    assertEq(fetchedDONs.length, allowedDONs.length);
    for (uint256 i = 0; i < allowedDONs.length; i++) {
      assertEq(fetchedDONs[i], allowedDONs[i]);
    }
  }

  function testGetAllAuthorizedAddressesByDON() public {
    // Add authorized addresses
    address[] memory authorizedAddresses = new address[](3);
    authorizedAddresses[0] = address(4);
    authorizedAddresses[1] = address(5);
    authorizedAddresses[2] = address(6);
    vm.prank(owner);

    for (uint32 i = 0; i < authorizedAddresses.length; i++) {
      vm.expectEmit(true, true, false, false);
      emit DONAccessControl.DONPermissionUpdatedV1(donID, authorizedAddresses[i], true);
    }

    registry.updateDONPermissions(donID, authorizedAddresses, true);

    // Verify the authorized addresses list
    address[] memory permissionedAddresses = registry.getAllAuthorizedAddressesByDON(donID, 0, 10);
    assertEq(permissionedAddresses.length, authorizedAddresses.length);
    for (uint256 i = 0; i < authorizedAddresses.length; i++) {
      assertEq(permissionedAddresses[i], authorizedAddresses[i]);
    }
  }

  function testGetWorkflowMetadataListByOwner() public {
    // Register multiple workflows for the same owner
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName2, workflowID2, WorkflowRegistry.WorkflowStatus.PAUSED);
    _allowAccessAndRegisterWorkflow(
      authorizedUser, "workflow3", keccak256(abi.encodePacked("workflow3")), WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      authorizedUser, "workflow4", keccak256(abi.encodePacked("workflow4")), WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      authorizedUser, "workflow5", keccak256(abi.encodePacked("workflow5")), WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Retrieve the list of workflows for the owner
    WorkflowRegistry.WorkflowMetadata[] memory workflows =
      registry.getWorkflowMetadataListByOwner(authorizedUser, 0, 10);

    // Verify the individual workflows are retrieved correctly
    assertEq(workflows[0].workflowID, workflowID1);
    assertEq(workflows[0].workflowName, workflowName1);
    assertEq(workflows[0].owner, authorizedUser);
    assertEq(workflows[0].binaryURL, testBinaryURL);
    assertEq(workflows[0].configURL, testConfigURL);
    assertEq(workflows[0].secretsURL, testSecretsURL);
    assertTrue(workflows[0].status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    assertEq(workflows[1].workflowID, workflowID2);
    assertEq(workflows[1].workflowName, workflowName2);
    assertEq(workflows[1].owner, authorizedUser);
    assertEq(workflows[1].binaryURL, testBinaryURL);
    assertEq(workflows[1].configURL, testConfigURL);
    assertEq(workflows[1].secretsURL, testSecretsURL);
    assertTrue(workflows[1].status == WorkflowRegistry.WorkflowStatus.PAUSED);

    // Pagination: Get first page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory firstPage = registry.getWorkflowMetadataListByOwner(authorizedUser, 0, 2);
    assertEq(firstPage.length, 2);
    assertEq(firstPage[0].workflowName, workflowName1);
    assertEq(firstPage[1].workflowName, workflowName2);

    // Pagination: Get second page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory secondPage =
      registry.getWorkflowMetadataListByOwner(authorizedUser, 2, 2);
    assertEq(secondPage.length, 2);
    assertEq(secondPage[0].workflowName, "workflow3");
    assertEq(secondPage[1].workflowName, "workflow4");

    // Pagination: Get last page (1 item)
    WorkflowRegistry.WorkflowMetadata[] memory lastPage = registry.getWorkflowMetadataListByOwner(authorizedUser, 4, 2);
    assertEq(lastPage.length, 1);
    assertEq(lastPage[0].workflowName, "workflow5");

    // Pagination: Request page beyond available items
    WorkflowRegistry.WorkflowMetadata[] memory emptyPage = registry.getWorkflowMetadataListByOwner(authorizedUser, 6, 2);
    assertEq(emptyPage.length, 0);

    // Pagination: Request all items at once
    WorkflowRegistry.WorkflowMetadata[] memory allItems = registry.getWorkflowMetadataListByOwner(authorizedUser, 0, 10);
    assertEq(allItems.length, 5);
  }

  function testGetWorkflowMetadataListByDON() public {
    // Register multiple workflows for the same DON ID
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName1, workflowID1, WorkflowRegistry.WorkflowStatus.ACTIVE);
    _allowAccessAndRegisterWorkflow(authorizedUser, workflowName2, workflowID2, WorkflowRegistry.WorkflowStatus.PAUSED);
    _allowAccessAndRegisterWorkflow(
      authorizedUser, "workflow3", keccak256(abi.encodePacked("workflow3")), WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      authorizedUser, "workflow4", keccak256(abi.encodePacked("workflow4")), WorkflowRegistry.WorkflowStatus.PAUSED
    );
    _allowAccessAndRegisterWorkflow(
      authorizedUser, "workflow5", keccak256(abi.encodePacked("workflow5")), WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Retrieve the list of workflows for the DON ID
    WorkflowRegistry.WorkflowMetadata[] memory workflows = registry.getWorkflowMetadataListByDON(donID, 0, 10);

    // Verify the individual workflows are retrieved correctly
    assertEq(workflows[0].workflowID, workflowID1);
    assertEq(workflows[0].workflowName, workflowName1);
    assertEq(workflows[0].owner, authorizedUser);
    assertEq(workflows[0].binaryURL, testBinaryURL);
    assertEq(workflows[0].configURL, testConfigURL);
    assertEq(workflows[0].secretsURL, testSecretsURL);
    assertTrue(workflows[0].status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    assertEq(workflows[1].workflowID, workflowID2);
    assertEq(workflows[1].workflowName, workflowName2);
    assertEq(workflows[1].owner, authorizedUser);
    assertEq(workflows[1].binaryURL, testBinaryURL);
    assertEq(workflows[1].configURL, testConfigURL);
    assertEq(workflows[1].secretsURL, testSecretsURL);
    assertTrue(workflows[1].status == WorkflowRegistry.WorkflowStatus.PAUSED);

    // Pagination: Get first page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory firstPage = registry.getWorkflowMetadataListByDON(donID, 0, 2);
    assertEq(firstPage.length, 2);
    assertEq(firstPage[0].workflowName, workflowName1);
    assertEq(firstPage[1].workflowName, workflowName2);

    // Pagination: Get second page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory secondPage = registry.getWorkflowMetadataListByDON(donID, 2, 2);
    assertEq(secondPage.length, 2);
    assertEq(secondPage[0].workflowName, "workflow3");
    assertEq(secondPage[1].workflowName, "workflow4");

    // Pagination: Get last page (1 item)
    WorkflowRegistry.WorkflowMetadata[] memory lastPage = registry.getWorkflowMetadataListByDON(donID, 4, 2);
    assertEq(lastPage.length, 1);
    assertEq(lastPage[0].workflowName, "workflow5");

    // Pagination: Request page beyond available items
    WorkflowRegistry.WorkflowMetadata[] memory emptyPage = registry.getWorkflowMetadataListByDON(donID, 6, 2);
    assertEq(emptyPage.length, 0);

    // Pagination: Request all items at once
    WorkflowRegistry.WorkflowMetadata[] memory allItems = registry.getWorkflowMetadataListByDON(donID, 0, 10);
    assertEq(allItems.length, 5);

    // Request from non-existent DON ID
    uint32 nonExistentDonID = 999;
    WorkflowRegistry.WorkflowMetadata[] memory emptyDON = registry.getWorkflowMetadataListByDON(nonExistentDonID, 0, 10);
    assertEq(emptyDON.length, 0);
  }
}
