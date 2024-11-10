// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Ownable2Step} from "../../shared/access/Ownable2Step.sol";
import {DONAccessControl} from "../dev/DONAccessControl.sol";
import {WorkflowRegistry} from "../dev/WorkflowRegistry.sol";
import {Test} from "forge-std/Test.sol";
import {Vm} from "forge-std/Vm.sol";

contract WorkflowRegistryTest is Test {
  WorkflowRegistry private s_registry;

  address private s_owner = address(1);
  address private s_unauthorizedUser = address(2);
  address private s_authorizedUser = address(3);
  bytes32 private s_workflowID1 = keccak256(abi.encodePacked("workflow1"));
  string private s_workflowName1 = "s_workflowName1";
  bytes32 private s_workflowID2 = keccak256(abi.encodePacked("workflow2"));
  string private s_workflowName2 = "s_workflowName2";
  bytes32 private s_newWorkflowID = keccak256(abi.encodePacked("workflow_new"));
  uint32 private s_donID = 1;
  string private s_testBinaryURL = "binaryURL";
  string private s_testConfigURL = "configURL";
  string private s_testSecretsURL = "secretsURL";

  function setUp() public {
    vm.prank(s_owner);
    s_registry = new WorkflowRegistry();
  }

  function _allowAccessAndRegisterWorkflow(
    address workflowOwner,
    string memory workflowName,
    bytes32 workflowID,
    WorkflowRegistry.WorkflowStatus initialStatus
  ) internal {
    _setupAuthorizedUser(workflowOwner);
    _setupAllowedDON(s_donID);

    // authorized user registers workflow
    vm.prank(workflowOwner);
    s_registry.registerWorkflow(
      workflowName,
      workflowID,
      s_donID,
      initialStatus,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );
  }

  function _setupAuthorizedUser(address workflowOwner) internal {
    // s_owner adds a single authorized address capable of registering workflows
    address[] memory authorizedUsers = new address[](1);
    authorizedUsers[0] = workflowOwner;
    vm.prank(s_owner);
    s_registry.updateDONPermissions(s_donID, authorizedUsers, true);
  }

  function _setupAllowedDON(uint32 _donID) internal {
    // s_owner adds a single DON ID allowed for registering workflows
    uint32[] memory allowedDONs = new uint32[](1);
    allowedDONs[0] = _donID;
    vm.prank(s_owner);
    s_registry.updateAllowedDONs(allowedDONs, true);
  }

  function testLockRegistry() public {
    // Ensure only the s_owner can lock the s_registry
    vm.prank(s_authorizedUser);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registry.lockRegistry();

    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Lock the s_registry as the s_owner
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Test all state-changing functions revert when s_registry is locked
    vm.startPrank(s_authorizedUser);

    // Test registerWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.registerWorkflow(
      s_workflowName2,
      s_workflowID2,
      s_donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );

    // Test updateWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.updateWorkflow(s_workflowName1, s_newWorkflowID, "newBinaryURL", "newConfigURL", "newSecretsURL");

    // Test pauseWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.pauseWorkflow(s_workflowName1);

    // Test activateWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.activateWorkflow(s_workflowName1);

    // Test deleteWorkflow
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.deleteWorkflow(s_workflowName1);

    // Test requestForceUpdateSecrets
    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    s_registry.requestForceUpdateSecrets(s_testSecretsURL);

    vm.stopPrank();

    // Test s_owner functions still revert when s_registry is locked
    vm.startPrank(s_owner);

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    address[] memory addresses = new address[](1);
    addresses[0] = s_authorizedUser;
    s_registry.updateDONPermissions(s_donID, addresses, true);

    vm.expectRevert(WorkflowRegistry.RegistryLocked.selector);
    uint32[] memory dons = new uint32[](1);
    dons[0] = s_donID;
    s_registry.updateAllowedDONs(dons, true);

    vm.stopPrank();
  }

  function testUnlockRegistry() public {
    // Check that the s_registry is initially unlocked
    bool isLocked = s_registry.isRegistryLocked();
    assertFalse(isLocked, "Registry should start off as unlocked");

    // Lock the s_registry first
    vm.prank(s_owner);
    s_registry.lockRegistry();

    // Ensure only the s_owner can unlock the s_registry
    vm.prank(s_authorizedUser);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_registry.unlockRegistry();

    // Unlock the s_registry as the s_owner
    vm.prank(s_owner);
    s_registry.unlockRegistry();

    // Perform an action that requires the s_registry to be unlocked
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Verify the workflow was registered successfully
    WorkflowRegistry.WorkflowMetadata memory workflow = s_registry.getWorkflowMetadata(
      s_authorizedUser,
      s_workflowName1
    );
    assertEq(workflow.workflowID, s_workflowID1);
  }

  function testRegisterWorkflowFailsForNotAuthorizedAddressOrForNotAllowedDONId() public {
    // s_owner of the contract is not allowed to register workflows without first setting permissions
    vm.prank(s_owner);
    vm.expectRevert(abi.encodeWithSelector(DONAccessControl.DONNotAllowed.selector, s_donID));
    s_registry.registerWorkflow(
      s_workflowName1,
      s_workflowID1,
      s_donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );

    // s_owner adds a single authorized address capable of registering workflows
    address[] memory authorizedUsers = new address[](1);
    authorizedUsers[0] = s_authorizedUser;
    vm.prank(s_owner);
    s_registry.updateDONPermissions(s_donID, authorizedUsers, true);

    // authorized address is still not able to register because DON ID is not allowed
    vm.prank(s_authorizedUser);
    vm.expectRevert(abi.encodeWithSelector(DONAccessControl.DONNotAllowed.selector, s_donID));
    s_registry.registerWorkflow(
      s_workflowName1,
      s_workflowID1,
      s_donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );

    // s_owner adds a single DON ID allowed for registering workflows
    uint32[] memory allowedDONs = new uint32[](1);
    allowedDONs[0] = s_donID;
    vm.prank(s_owner);
    s_registry.updateAllowedDONs(allowedDONs, true);

    // authorized address is finally able to register workflow
    vm.prank(s_authorizedUser);
    s_registry.registerWorkflow(
      s_workflowName1,
      s_workflowID1,
      s_donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );

    // sanity check by retrieving the workflow metadata
    WorkflowRegistry.WorkflowMetadata memory workflow = s_registry.getWorkflowMetadata(
      s_authorizedUser,
      s_workflowName1
    );
    assertEq(workflow.workflowID, s_workflowID1);
    assertEq(workflow.workflowName, s_workflowName1);
    assertEq(workflow.owner, s_authorizedUser);
    assertEq(workflow.binaryURL, s_testBinaryURL);
    assertEq(workflow.configURL, s_testConfigURL);
    assertEq(workflow.secretsURL, s_testSecretsURL);
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // any other unauthorized address still gets the unauthorized error
    vm.prank(s_unauthorizedUser);
    vm.expectRevert(
      abi.encodeWithSelector(DONAccessControl.AddressNotAuthorized.selector, s_donID, s_unauthorizedUser)
    );
    s_registry.registerWorkflow(
      s_workflowName1,
      s_workflowID1,
      s_donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );
  }

  function testUpdateWorkflow() public {
    // create a new workflow
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // authorized user tries to update the workflow by using the same workflow ID as before
    vm.prank(s_authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowIDNotUpdated.selector);
    s_registry.updateWorkflow(s_workflowName1, s_workflowID1, "newBinaryURL", "newConfigURL", "newSecretsURL");

    // now the authorizer user sets the new workflow ID
    vm.prank(s_authorizedUser);
    s_registry.updateWorkflow(s_workflowName1, s_newWorkflowID, "newBinaryURL", "newConfigURL", "newSecretsURL");

    // sanity check by retrieving the workflow metadata to make sure parameters are updated
    WorkflowRegistry.WorkflowMetadata memory workflow = s_registry.getWorkflowMetadata(
      s_authorizedUser,
      s_workflowName1
    );
    assertEq(workflow.workflowID, s_newWorkflowID);
    assertEq(workflow.workflowName, s_workflowName1);
    assertEq(workflow.owner, s_authorizedUser);
    assertEq(workflow.binaryURL, "newBinaryURL");
    assertEq(workflow.configURL, "newConfigURL");
    assertEq(workflow.secretsURL, "newSecretsURL");
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);
  }

  function testPauseWorkflow() public {
    // create a new workflow
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // authorized user pauses the workflow
    vm.prank(s_authorizedUser);
    s_registry.pauseWorkflow(s_workflowName1);

    // sanity check the workflow status update
    WorkflowRegistry.WorkflowMetadata memory workflow = s_registry.getWorkflowMetadata(
      s_authorizedUser,
      s_workflowName1
    );
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.PAUSED);

    // authorized user is not able to pause the workflow twice in a row
    vm.prank(s_authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyInDesiredStatus.selector);
    s_registry.pauseWorkflow(s_workflowName1);
  }

  function testActivateWorkflow() public {
    // create a new workflow
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.PAUSED
    );

    // authorized user activates the workflow
    vm.prank(s_authorizedUser);
    s_registry.activateWorkflow(s_workflowName1);

    // sanity check the workflow status update
    WorkflowRegistry.WorkflowMetadata memory workflow = s_registry.getWorkflowMetadata(
      s_authorizedUser,
      s_workflowName1
    );
    assertTrue(workflow.status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    // authorized user is not able to activate the workflow twice in a row
    vm.prank(s_authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowAlreadyInDesiredStatus.selector);
    s_registry.activateWorkflow(s_workflowName1);
  }

  function testNonWorkflowOwnerUserCannotUpdateWorkflow() public {
    // create a new workflow for one authorized user
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // add a new authorized user capable of registering workflows
    address anotherAuthorizedUser = address(567);
    address[] memory authorizedUsers = new address[](1);
    authorizedUsers[0] = anotherAuthorizedUser;
    vm.prank(s_owner);
    s_registry.updateDONPermissions(s_donID, authorizedUsers, true);

    // new authorized user is not able to update another user's workflow (same workflow name)
    vm.prank(anotherAuthorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.updateWorkflow(s_workflowName1, s_newWorkflowID, "newBinaryURL", "newConfigURL", "newSecretsURL");
  }

  function testRequestForceUpdateSecrets() public {
    // Register two workflows with the same secretsURL for the authorized user
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    vm.prank(s_authorizedUser);
    s_registry.registerWorkflow(
      s_workflowName2,
      s_workflowID2,
      s_donID,
      WorkflowRegistry.WorkflowStatus.ACTIVE,
      s_testBinaryURL,
      s_testConfigURL,
      s_testSecretsURL
    );

    // Attempt force update secrets from an unauthorized user
    vm.prank(s_unauthorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.requestForceUpdateSecrets(s_testSecretsURL);

    // Start recording the logs to later check the event content
    vm.recordLogs();

    // Authorized user requests force update secrets
    vm.prank(s_authorizedUser);
    s_registry.requestForceUpdateSecrets(s_testSecretsURL);

    // Verify the events emitted with correct details
    Vm.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries.length, 2); // Expecting two separate events for each individual workflow

    for (uint256 i = 0; i < entries.length; i++) {
      assertEq(entries[i].topics[0], keccak256("WorkflowForceUpdateSecretsRequestedV1(string,address,string)"));

      // Compare the hash of the expected string with the topic
      bytes32 expectedSecretsURLHash = keccak256(abi.encodePacked(s_testSecretsURL));
      assertEq(entries[i].topics[1], expectedSecretsURLHash);

      // Decode the indexed address
      address decodedOwner = abi.decode(abi.encodePacked(entries[i].topics[2]), (address));
      assertEq(decodedOwner, s_authorizedUser);

      // Decode the non-indexed workflow name
      string memory decodedWorkflowName = abi.decode(entries[i].data, (string));

      // Assert the values
      if (i == 0) {
        assertEq(decodedWorkflowName, s_workflowName1);
      } else {
        assertEq(decodedWorkflowName, s_workflowName2);
      }
    }
  }

  function testDeleteWorkflow() public {
    // Create a new workflow
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Unauthorized user should not be able to delete the workflow
    vm.prank(s_unauthorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.deleteWorkflow(s_workflowName1);

    // Authorized user deletes the workflow
    vm.prank(s_authorizedUser);
    s_registry.deleteWorkflow(s_workflowName1);

    // Sanity check to verify that the workflow has been deleted
    vm.prank(s_authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.getWorkflowMetadata(s_authorizedUser, s_workflowName1);

    // Authorized user should not be able to delete a non-existing workflow
    vm.prank(s_authorizedUser);
    vm.expectRevert(WorkflowRegistry.WorkflowDoesNotExist.selector);
    s_registry.deleteWorkflow(s_workflowName1);
  }

  function testGetAllAllowedDONs() public {
    // Add allowed DON IDs
    uint32[] memory allowedDONs = new uint32[](3);
    allowedDONs[0] = 1;
    allowedDONs[1] = 2;
    allowedDONs[2] = 3;
    vm.prank(s_owner);

    for (uint32 i = 0; i < allowedDONs.length; i++) {
      vm.expectEmit(true, true, false, false);
      emit DONAccessControl.AllowedDONUpdatedV1(allowedDONs[i], true);
    }

    s_registry.updateAllowedDONs(allowedDONs, true);

    // Verify the allowed DONs list
    uint32[] memory fetchedDONs = s_registry.getAllAllowedDONs();
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
    vm.prank(s_owner);

    for (uint32 i = 0; i < authorizedAddresses.length; i++) {
      vm.expectEmit(true, true, false, false);
      emit DONAccessControl.DONPermissionUpdatedV1(s_donID, authorizedAddresses[i], true);
    }

    s_registry.updateDONPermissions(s_donID, authorizedAddresses, true);

    // Verify the authorized addresses list
    address[] memory permissionedAddresses = s_registry.getAllAuthorizedAddressesByDON(s_donID, 0, 10);
    assertEq(permissionedAddresses.length, authorizedAddresses.length);
    for (uint256 i = 0; i < authorizedAddresses.length; i++) {
      assertEq(permissionedAddresses[i], authorizedAddresses[i]);
    }
  }

  function testGetWorkflowMetadataListByOwner() public {
    // Register multiple workflows for the same s_owner
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName2,
      s_workflowID2,
      WorkflowRegistry.WorkflowStatus.PAUSED
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      "workflow3",
      keccak256(abi.encodePacked("workflow3")),
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      "workflow4",
      keccak256(abi.encodePacked("workflow4")),
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      "workflow5",
      keccak256(abi.encodePacked("workflow5")),
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Retrieve the list of workflows for the s_owner
    WorkflowRegistry.WorkflowMetadata[] memory workflows = s_registry.getWorkflowMetadataListByOwner(
      s_authorizedUser,
      0,
      10
    );

    // Verify the individual workflows are retrieved correctly
    assertEq(workflows[0].workflowID, s_workflowID1);
    assertEq(workflows[0].workflowName, s_workflowName1);
    assertEq(workflows[0].owner, s_authorizedUser);
    assertEq(workflows[0].binaryURL, s_testBinaryURL);
    assertEq(workflows[0].configURL, s_testConfigURL);
    assertEq(workflows[0].secretsURL, s_testSecretsURL);
    assertTrue(workflows[0].status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    assertEq(workflows[1].workflowID, s_workflowID2);
    assertEq(workflows[1].workflowName, s_workflowName2);
    assertEq(workflows[1].owner, s_authorizedUser);
    assertEq(workflows[1].binaryURL, s_testBinaryURL);
    assertEq(workflows[1].configURL, s_testConfigURL);
    assertEq(workflows[1].secretsURL, s_testSecretsURL);
    assertTrue(workflows[1].status == WorkflowRegistry.WorkflowStatus.PAUSED);

    // Pagination: Get first page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory firstPage = s_registry.getWorkflowMetadataListByOwner(
      s_authorizedUser,
      0,
      2
    );
    assertEq(firstPage.length, 2);
    assertEq(firstPage[0].workflowName, s_workflowName1);
    assertEq(firstPage[1].workflowName, s_workflowName2);

    // Pagination: Get second page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory secondPage = s_registry.getWorkflowMetadataListByOwner(
      s_authorizedUser,
      2,
      2
    );
    assertEq(secondPage.length, 2);
    assertEq(secondPage[0].workflowName, "workflow3");
    assertEq(secondPage[1].workflowName, "workflow4");

    // Pagination: Get last page (1 item)
    WorkflowRegistry.WorkflowMetadata[] memory lastPage = s_registry.getWorkflowMetadataListByOwner(
      s_authorizedUser,
      4,
      2
    );
    assertEq(lastPage.length, 1);
    assertEq(lastPage[0].workflowName, "workflow5");

    // Pagination: Request page beyond available items
    WorkflowRegistry.WorkflowMetadata[] memory emptyPage = s_registry.getWorkflowMetadataListByOwner(
      s_authorizedUser,
      6,
      2
    );
    assertEq(emptyPage.length, 0);

    // Pagination: Request all items at once
    WorkflowRegistry.WorkflowMetadata[] memory allItems = s_registry.getWorkflowMetadataListByOwner(
      s_authorizedUser,
      0,
      10
    );
    assertEq(allItems.length, 5);
  }

  function testGetWorkflowMetadataListByDON() public {
    // Register multiple workflows for the same DON ID
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName1,
      s_workflowID1,
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      s_workflowName2,
      s_workflowID2,
      WorkflowRegistry.WorkflowStatus.PAUSED
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      "workflow3",
      keccak256(abi.encodePacked("workflow3")),
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      "workflow4",
      keccak256(abi.encodePacked("workflow4")),
      WorkflowRegistry.WorkflowStatus.PAUSED
    );
    _allowAccessAndRegisterWorkflow(
      s_authorizedUser,
      "workflow5",
      keccak256(abi.encodePacked("workflow5")),
      WorkflowRegistry.WorkflowStatus.ACTIVE
    );

    // Retrieve the list of workflows for the DON ID
    WorkflowRegistry.WorkflowMetadata[] memory workflows = s_registry.getWorkflowMetadataListByDON(s_donID, 0, 10);

    // Verify the individual workflows are retrieved correctly
    assertEq(workflows[0].workflowID, s_workflowID1);
    assertEq(workflows[0].workflowName, s_workflowName1);
    assertEq(workflows[0].owner, s_authorizedUser);
    assertEq(workflows[0].binaryURL, s_testBinaryURL);
    assertEq(workflows[0].configURL, s_testConfigURL);
    assertEq(workflows[0].secretsURL, s_testSecretsURL);
    assertTrue(workflows[0].status == WorkflowRegistry.WorkflowStatus.ACTIVE);

    assertEq(workflows[1].workflowID, s_workflowID2);
    assertEq(workflows[1].workflowName, s_workflowName2);
    assertEq(workflows[1].owner, s_authorizedUser);
    assertEq(workflows[1].binaryURL, s_testBinaryURL);
    assertEq(workflows[1].configURL, s_testConfigURL);
    assertEq(workflows[1].secretsURL, s_testSecretsURL);
    assertTrue(workflows[1].status == WorkflowRegistry.WorkflowStatus.PAUSED);

    // Pagination: Get first page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory firstPage = s_registry.getWorkflowMetadataListByDON(s_donID, 0, 2);
    assertEq(firstPage.length, 2);
    assertEq(firstPage[0].workflowName, s_workflowName1);
    assertEq(firstPage[1].workflowName, s_workflowName2);

    // Pagination: Get second page (2 items)
    WorkflowRegistry.WorkflowMetadata[] memory secondPage = s_registry.getWorkflowMetadataListByDON(s_donID, 2, 2);
    assertEq(secondPage.length, 2);
    assertEq(secondPage[0].workflowName, "workflow3");
    assertEq(secondPage[1].workflowName, "workflow4");

    // Pagination: Get last page (1 item)
    WorkflowRegistry.WorkflowMetadata[] memory lastPage = s_registry.getWorkflowMetadataListByDON(s_donID, 4, 2);
    assertEq(lastPage.length, 1);
    assertEq(lastPage[0].workflowName, "workflow5");

    // Pagination: Request page beyond available items
    WorkflowRegistry.WorkflowMetadata[] memory emptyPage = s_registry.getWorkflowMetadataListByDON(s_donID, 6, 2);
    assertEq(emptyPage.length, 0);

    // Pagination: Request all items at once
    WorkflowRegistry.WorkflowMetadata[] memory allItems = s_registry.getWorkflowMetadataListByDON(s_donID, 0, 10);
    assertEq(allItems.length, 5);

    // Request from non-existent DON ID
    uint32 nonExistentDonID = 999;
    WorkflowRegistry.WorkflowMetadata[] memory emptyDON = s_registry.getWorkflowMetadataListByDON(
      nonExistentDonID,
      0,
      10
    );
    assertEq(emptyDON.length, 0);
  }
}
