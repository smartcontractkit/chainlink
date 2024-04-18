// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Report} from "../../libraries/Report.sol";

contract KeystoneReceiver {
    event MessageReceived(bytes32 indexed workflowId, bytes32 indexed workflowExecutionId, bytes[] mercuryReports);

    error InvalidReport(bytes data);

    uint256 private constant REPORT_LENGTH = 64;

    function foo(bytes calldata rawReport) external {
        if (rawReport.length < REPORT_LENGTH) {
            revert InvalidReport(rawReport);
        }

        // decode metadata
        (bytes32 workflowId, bytes32 workflowExecutionId) = Report.getMetadata(rawReport);
        // parse actual report
        bytes[] memory mercuryReports = abi.decode(rawReport[64:], (bytes[]));
        emit MessageReceived(workflowId, workflowExecutionId, mercuryReports);
    }
}
