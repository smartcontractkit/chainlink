// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMMultiOffRamp} from "../../offRamp/EVM2EVMMultiOffRamp.sol";

contract ReportCodec {
  event ExecuteReportDecoded(Internal.ExecutionReportSingleChain[] report);
  event CommitReportDecoded(EVM2EVMMultiOffRamp.CommitReport report);

  function decodeExecuteReport(bytes memory report) public pure returns (Internal.ExecutionReportSingleChain[] memory) {
    return abi.decode(report, (Internal.ExecutionReportSingleChain[]));
  }

  function decodeCommitReport(bytes memory report) public pure returns (EVM2EVMMultiOffRamp.CommitReport memory) {
    return abi.decode(report, (EVM2EVMMultiOffRamp.CommitReport));
  }
}
