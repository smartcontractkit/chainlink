// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

/// @title Keystone Feeds Permission Handler
/// @notice This contract is designed to manage and validate permissions for accessing specific reports within a decentralized system.
/// @dev The contract uses mappings to keep track of report permissions associated with a unique report ID.
abstract contract KeystoneFeedsPermissionHandler is OwnerIsCreator {
  /// @notice Holds the details for permissions of a report
  /// @dev Workflow names and report names are stored as bytes to optimize for gas efficiency.
  struct Permission {
    address forwarder;
    bytes10 workflowName;
    bytes2 reportName;
    address workflowOwner;
    bool isAllowed;
  }

  /// @notice Event emitted when report permissions are set
  event ReportPermissionSet(bytes32 indexed reportId, Permission permission);

  /// @notice Error to be thrown when an unauthorized access attempt is made
  error Unauthorized(address forwarder, address workflowOwner, bytes10 workflowName, bytes2 reportName);

  /// @dev Mapping from a report ID to a boolean indicating whether the report is allowed or not
  mapping(bytes32 => bool) internal s_allowedReports;

  /// @notice Sets permissions for multiple reports
  /// @param permissions An array of Permission structs for which to set permissions
  /// @dev Emits a ReportPermissionSet event for each permission set
  function setReportPermissions(Permission[] memory permissions) external onlyOwner {
    for (uint256 i; i < permissions.length; i++) {
      _setReportPermission(permissions[i]);
    }
  }

  /// @dev Internal function to set a single report permission
  /// @param permission The Permission struct containing details about the permission to set
  /// @dev Emits a ReportPermissionSet event
  function _setReportPermission(Permission memory permission) internal {
    bytes32 reportId = _createReportId(
      permission.forwarder,
      permission.workflowOwner,
      permission.workflowName,
      permission.reportName
    );
    s_allowedReports[reportId] = permission.isAllowed;
    emit ReportPermissionSet(reportId, permission);
  }

  /// @dev Internal view function to validate if a report is allowed for a given set of details
  /// @param forwarder The address of the forwarder
  /// @param workflowOwner The address of the workflow owner
  /// @param workflowName The name of the workflow in bytes10
  /// @param reportName The name of the report in bytes2
  /// @dev Reverts with Unauthorized if the report is not allowed
  function _validateReportPermission(
    address forwarder,
    address workflowOwner,
    bytes10 workflowName,
    bytes2 reportName
  ) internal view {
    bytes32 reportId = _createReportId(forwarder, workflowOwner, workflowName, reportName);
    if (!s_allowedReports[reportId]) {
      revert Unauthorized(forwarder, workflowOwner, workflowName, reportName);
    }
  }

  function _createReportId(
    address forwarder,
    address workflowOwner,
    bytes10 workflowName,
    bytes2 reportName
  ) internal pure returns (bytes32) {
    return keccak256(abi.encode(forwarder, workflowOwner, workflowName, reportName));
  }
}
