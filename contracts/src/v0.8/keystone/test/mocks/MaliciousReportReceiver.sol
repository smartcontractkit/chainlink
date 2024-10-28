// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IReceiver} from "../../interfaces/IReceiver.sol";

contract MaliciousReportReceiver is IReceiver, IERC165 {
  event MessageReceived(bytes metadata, bytes[] mercuryReports);
  bytes public latestReport;

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
    // Exhaust all gas that was provided
    for (uint256 i = 0; i < 1_000_000_000; ++i) {
      bytes[] memory mercuryReports = abi.decode(rawReport, (bytes[]));
      latestReport = rawReport;
      emit MessageReceived(metadata, mercuryReports);
    }
  }

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}
