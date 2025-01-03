// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";

contract ReportCodec {
  event ExecuteReportDecoded(Internal.ExecutionReport[] report);
  event CommitReportDecoded(OffRamp.CommitReport report);

  function decodeExecuteReport(
    bytes memory report
  ) public pure returns (Internal.ExecutionReport[] memory) {
    return abi.decode(report, (Internal.ExecutionReport[]));
  }

  function decodeCommitReport(
    bytes memory report
  ) public pure returns (OffRamp.CommitReport memory) {
    return abi.decode(report, (OffRamp.CommitReport));
  }
}
