// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title Report Library
/// @notice A library for handling Keystone reports. Used by KeystoneForwarder
/// and recipient (end-user) contracts.
library Report {
    /// @notice Extracts the Keystone metadata from the report.
    /// @param rawReport The raw report data without the function selector.
    function getMetadata(bytes memory rawReport)
        public
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
