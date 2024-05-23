// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title IReceiver - receives keystone reports
interface IReceiver {
  function onReport(
    bytes32 workflowId,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName,
    bytes calldata rawReport
  ) external;
}
