// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";

contract ReportCodec {
  event ExecuteReportDecoded(Internal.ExecutionReportSingleChain[] report);
  event CommitReportDecoded(OffRamp.CommitReport report);

  function decodeExecuteReport(bytes memory report) public pure returns (Internal.ExecutionReportSingleChain[] memory) {
    return abi.decode(report, (Internal.ExecutionReportSingleChain[]));
  }

  function decodeCommitReport(bytes memory report) public pure returns (OffRamp.CommitReport memory) {
    return abi.decode(report, (OffRamp.CommitReport));
  }
}
