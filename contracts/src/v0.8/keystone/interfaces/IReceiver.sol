// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title IReceiver - receives keystone reports
interface IReceiver {
  /// @notice Handles incoming keystone reports.
  /// @dev If this function call reverts, it can be retried with a higher gas
  /// limit. The receiver is responsible for discarding stale reports.
  /// @param metadata Report's metadata.
  /// @param report Workflow report.
  function onReport(bytes calldata metadata, bytes calldata report) external;
}
