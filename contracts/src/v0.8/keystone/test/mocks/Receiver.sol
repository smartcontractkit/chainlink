// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
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

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}
