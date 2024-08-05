// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IReceiver} from "../../interfaces/IReceiver.sol";

contract Receiver is IReceiver {
  event MessageReceived(bytes metadata, bytes[] mercuryReports);
  bytes public latestReport;

  constructor() {}

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
    latestReport = rawReport;

    // parse actual report
    bytes[] memory mercuryReports = abi.decode(rawReport, (bytes[]));
    emit MessageReceived(metadata, mercuryReports);
  }
}
