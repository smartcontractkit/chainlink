// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract KeystoneReceiver {
    event MessageReceived(bytes32 indexed workflowId, bytes32 indexed workflowExecutionId, bytes[] mercuryReports);

    error InvalidReport(bytes data);

    uint256 private constant REPORT_LENGTH = 64;

    function foo(bytes calldata rawReport) external {
        if (rawReport.length < REPORT_LENGTH) {
            revert InvalidReport(rawReport);
        }

        // decode metadata
        (bytes32 workflowId, bytes32 workflowExecutionId) = _splitReport(rawReport);
        // parse actual report
        bytes[] memory mercuryReports = abi.decode(rawReport[64:], (bytes[]));
        emit MessageReceived(workflowId, workflowExecutionId, mercuryReports);
    }

    function _splitReport(bytes memory rawReport)
        internal
        pure
        returns (bytes32 workflowId, bytes32 workflowExecutionId)
    {
        assembly {
            // skip first 32 bytes, contains length of the report
            workflowId := mload(add(rawReport, 32))
            workflowExecutionId := mload(add(rawReport, 64))
        }

        return (workflowId, workflowExecutionId);
    }
}
