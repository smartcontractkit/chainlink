// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "../../interfaces/IReceiver.sol";

contract Receiver is IReceiver {
  error Unauthorized(
    address forwarderAddress,
    bytes32 workflowId,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName
  );

  event MessageReceived(
    bytes32 workflowId,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName,
    bytes rawReport
  );

  bytes32 internal s_allowedWorkflowReport;
  bytes32 internal s_allowedOwnerReport;

  function _getWorkflowReportHash(
    address forwarderAddress,
    bytes32 workflowId,
    bytes32 reportName
  ) internal pure returns (bytes32) {
    return keccak256(abi.encode(forwarderAddress, workflowId, reportName));
  }

  function setAllowedWorkflowReport(address forwarderAddress, bytes32 workflowId, bytes32 reportName) external {
    s_allowedWorkflowReport = _getWorkflowReportHash(forwarderAddress, workflowId, reportName);
  }

  function _getOwnerReportHash(
    address forwarderAddress,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName
  ) internal pure returns (bytes32) {
    return keccak256(abi.encode(forwarderAddress, workflowOwner, workflowName, reportName));
  }

  function setAllowedOwnerReport(
    address forwarderAddress,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName
  ) external {
    s_allowedOwnerReport = _getOwnerReportHash(forwarderAddress, workflowOwner, workflowName, reportName);
  }

  function onReport(
    bytes32 workflowId,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName,
    bytes calldata rawReport
  ) external {
    // if (
    //   _getWorkflowReportHash(msg.sender, workflowId, reportName) != s_allowedWorkflowReport &&
    //   _getOwnerReportHash(msg.sender, workflowOwner, workflowName, reportName) != s_allowedOwnerReport
    // ) {
    //   revert Unauthorized(msg.sender, workflowId, workflowOwner, workflowName, reportName);
    // }

    emit MessageReceived(workflowId, workflowOwner, workflowName, reportName, rawReport);
  }
}
