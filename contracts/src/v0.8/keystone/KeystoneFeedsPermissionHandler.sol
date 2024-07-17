// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

abstract contract KeystoneFeedsPermissionHandler {
    struct Permission {
        address forwarder;
        address workflowOwner;
        bytes10 workflowName;
        bytes2 reportName;
        bool isAllowed;
    }

    event ReportPermissionSet(bytes32 indexed reportId, Permission config);
    error Unauthorized(address fowarder, address workflowOwner, bytes10 workflowName, bytes2 reportName);

    mapping(bytes32 => bool) internal s_allowedReports;

    function setReportPermissions(
        Permission[] memory permissions
    ) external {
        for (uint i = 0; i < permissions.length; i++) {
            _setReportPermission(permissions[i]);
        }
    }

    function _setReportPermission(Permission memory permission) internal {
        bytes32 reportId = keccak256(abi.encode(
            permission.forwarder,
            permission.workflowOwner,
            permission.workflowName,
            permission.reportName
        ));
        s_allowedReports[reportId] = permission.isAllowed;
        emit ReportPermissionSet(reportId, permission);
    }

    function _validateReportPermission(
        address forwarder,
        address workflowOwner,
        bytes10 workflowName,
        bytes2 reportName
    ) internal view {
        bytes32 reportId = keccak256(abi.encode(
            forwarder,
            workflowOwner,
            workflowName,
            reportName
        ));
        if (!s_allowedReports[reportId]) {
            revert Unauthorized(
                forwarder,
                workflowOwner,
                workflowName,
                reportName
            );
        }
    }
}
