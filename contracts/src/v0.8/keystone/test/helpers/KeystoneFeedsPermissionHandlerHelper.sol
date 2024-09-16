// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "../../KeystoneFeedsPermissionHandler.sol";

contract KeystoneFeedsPermissionHandlerHelper is KeystoneFeedsPermissionHandler {
    function validateReportPermission(address workflowOwner, bytes10 workflowName) external view {
        _validateReportPermission(msg.sender, workflowOwner, workflowName);
    }

    function createReportId(Permission memory permission) external pure returns (bytes32) {
        return
            _createReportId(permission.forwarder, permission.workflowOwner, permission.workflowName);
    }

    function getAllowedReports(bytes32 reportId) external view returns (bool) {
        return s_allowedReports[reportId];
    }
}