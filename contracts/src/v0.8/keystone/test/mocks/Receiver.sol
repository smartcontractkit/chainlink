// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "../../interfaces/IReceiver.sol";

contract Receiver is IReceiver {
  event MessageReceived(bytes metadata, bytes[] mercuryReports);

  constructor() {}

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
    // parse actual report
    bytes[] memory mercuryReports = abi.decode(rawReport, (bytes[]));
    emit MessageReceived(metadata, mercuryReports);
  }
}
