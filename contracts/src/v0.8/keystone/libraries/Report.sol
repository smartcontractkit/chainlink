// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title Report Library
/// @notice A library for handling Keystone reports. Used by KeystoneForwarder
/// and recipient (end-user) contracts.
library Report {
  /// @notice This error is returned when the report is malformed. We expect
  /// Keystone reports to be at least 64 bytes long (the length of the metadata).
  /// @param report the data that was received
  error InvalidReport(bytes report);

  /// @notice Extracts the Keystone metadata from the report.
  /// @param report The raw report data without the function selector.
  function getMetadata(bytes memory report) public pure returns (bytes32 workflowId, bytes32 workflowExecutionId) {
    if (report.length < 64) {
      revert InvalidReport(report);
    }

    assembly ("memory-safe") {
      // skip first 32 bytes, contains length of the report
      workflowId := mload(add(report, 32))
      workflowExecutionId := mload(add(report, 64))
    }

    return (workflowId, workflowExecutionId);
  }
}
