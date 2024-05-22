// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "../../interfaces/IReceiver.sol";

contract Receiver is IReceiver {
  event MessageReceived(bytes32 indexed workflowId, address indexed workflowOwner, bytes[] mercuryReports);

  constructor() {}

  function onReport(bytes32 workflowId, address workflowOwner, bytes calldata rawReport) external {
    // parse actual report
    bytes[] memory mercuryReports = abi.decode(rawReport, (bytes[]));
    emit MessageReceived(workflowId, workflowOwner, mercuryReports);
  }
}
