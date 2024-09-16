// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {KeystoneFeedsPermissionHandler} from "../KeystoneFeedsPermissionHandler.sol";
import "./helpers/KeystoneFeedsPermissionHandlerHelper.sol";

contract KeystoneFeedsPermissionHandlerTestSetup is Test {
    address constant FORWARDER_1 = address(0x1);
    address constant FORWARDER_2 = address(0x2);
    address constant WORKFLOW_OWNER_1 = address(0x3);
    address constant WORKFLOW_OWNER_2 = address(0x4);
    bytes10 constant WORKFLOW_NAME_1 = "workflow1";
    bytes10 constant WORKFLOW_NAME_2 = "workflow2";
    address constant UNKNOWN_ADDRESS = address(0x5);
    KeystoneFeedsPermissionHandlerHelper permissionHandlerHelper;

    function _createPermissions() internal pure returns (KeystoneFeedsPermissionHandler.Permission[] memory) {
        KeystoneFeedsPermissionHandler.Permission[] memory permissions = new KeystoneFeedsPermissionHandler.Permission[](2);
        permissions[0] = KeystoneFeedsPermissionHandler.Permission(
            FORWARDER_1,
            WORKFLOW_NAME_1,
            WORKFLOW_OWNER_1,
            true
        );
        permissions[1] = KeystoneFeedsPermissionHandler.Permission(
            FORWARDER_2,
            WORKFLOW_NAME_2,
            WORKFLOW_OWNER_2,
            true
        );
        return permissions;
    }

    function setUp() public {
        permissionHandlerHelper = new KeystoneFeedsPermissionHandlerHelper();
    }
}

contract KeystoneFeeds_setReportPermissions is KeystoneFeedsPermissionHandlerTestSetup {
    function test__setReportPermissions() public {
//create permissions
        KeystoneFeedsPermissionHandler.Permission[] memory permissions = _createPermissions();
        bytes32[] memory reportIds = new bytes32[](permissions.length);
        reportIds[0] = permissionHandlerHelper.createReportId(permissions[0]);
        reportIds[1] = permissionHandlerHelper.createReportId(permissions[1]);

//expected events
        for (uint256 i; i < permissions.length; i++) {
            vm.expectEmit();
            emit KeystoneFeedsPermissionHandler.ReportPermissionSet(reportIds[i], permissions[i]);
        }

//set permissions
        permissionHandlerHelper.setReportPermissions(permissions);

//assert permissions
        for (uint256 i; i < permissions.length; i++) {
            vm.assertEq(permissionHandlerHelper.getAllowedReports(reportIds[i]), permissions[i].isAllowed);
        }
    }

    function test__setReportPermissions_InvalidOwner_Reverts() public {
//create permissions
        KeystoneFeedsPermissionHandler.Permission[] memory permissions = _createPermissions();
        bytes32[] memory reportIds = new bytes32[](permissions.length);
        reportIds[0] = permissionHandlerHelper.createReportId(permissions[0]);
        reportIds[1] = permissionHandlerHelper.createReportId(permissions[1]);

//set permissions
        vm.startPrank(UNKNOWN_ADDRESS);
        vm.expectRevert("Only callable by owner");
        permissionHandlerHelper.setReportPermissions(permissions);
        vm.stopPrank();
    }
}

contract KeystoneFeeds_validatePermissions is KeystoneFeedsPermissionHandlerTestSetup {
    function test__validatePermissions() public {
//create permissions
        KeystoneFeedsPermissionHandler.Permission[] memory permissions = _createPermissions();
        bytes32[] memory reportIds = new bytes32[](permissions.length);
        reportIds[0] = permissionHandlerHelper.createReportId(permissions[0]);
        reportIds[1] = permissionHandlerHelper.createReportId(permissions[1]);

//set permissions
        permissionHandlerHelper.setReportPermissions(permissions);

//validate permissions - if this does not revert it means the permissions are valid
        vm.prank(FORWARDER_1);
        permissionHandlerHelper.validateReportPermission(
            permissions[0].workflowOwner,
            permissions[0].workflowName
        );

        vm.prank(FORWARDER_2);
        permissionHandlerHelper.validateReportPermission(
            permissions[1].workflowOwner,
            permissions[1].workflowName
        );
    }

    function test__validatePermissions_InvalidForwarder_Reverts() public {
//create permissions
        KeystoneFeedsPermissionHandler.Permission[] memory permissions = _createPermissions();
        bytes32[] memory reportIds = new bytes32[](permissions.length);
        reportIds[0] = permissionHandlerHelper.createReportId(permissions[0]);
        reportIds[1] = permissionHandlerHelper.createReportId(permissions[1]);

//set permissions
        permissionHandlerHelper.setReportPermissions(permissions);

//trying out invalid permissions by changing the forwarder to FORWARDER_2 for Permission : 1
        vm.expectRevert(
            abi.encodeWithSelector(
                KeystoneFeedsPermissionHandler.ReportForwarderUnauthorized.selector,
                UNKNOWN_ADDRESS,
                permissions[0].workflowOwner,
                permissions[0].workflowName
            )
        );
        vm.startPrank(UNKNOWN_ADDRESS);
        permissionHandlerHelper.validateReportPermission(
            permissions[0].workflowOwner,
            permissions[0].workflowName
        );
        vm.stopPrank();
    }
}