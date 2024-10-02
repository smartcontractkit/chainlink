// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {DataFeedsBase} from "../dev/DataFeedsBase.sol";
import {DataFeedsTestSetup} from "./DataFeedsTestSetup.t.sol";

contract DataFeedsBaseTest is DataFeedsTestSetup {
  function setUp() public virtual override {
    DataFeedsTestSetup.setUp();
  }

  function test_getFeedIds() public {
    bytes[] memory reports = new bytes[](3);
    reports[0] = reportBlockPremium;
    reports[1] = reportBasic;
    reports[2] = reportPremium;

    bytes[] memory reportData = DataFeedsBase.getReportData(reports);

    bytes32[] memory feedIds = DataFeedsBase.getFeedIds(reportData);

    assertEq(feedIds[0], reportStructBlockPremium.feedId);
    assertEq(feedIds[1], reportStructBasic.feedId);
    assertEq(feedIds[2], reportStructPremium.feedId);
  }

  function test_getBenchmarksAndTimestamps() public {
    bytes[] memory reports = new bytes[](3);
    reports[0] = reportBlockPremium;
    reports[1] = reportBasic;
    reports[2] = reportPremium;

    bytes[] memory reportData = DataFeedsBase.getReportData(reports);

    (int256[] memory benchmarksData, uint256[] memory timestampsData) = DataFeedsBase.getBenchmarksAndTimestamps(
      reportData
    );

    assertEq(benchmarksData[0], reportStructBlockPremium.benchmarkPrice);
    assertEq(benchmarksData[1], reportStructBasic.benchmarkPrice);
    assertEq(benchmarksData[2], reportStructPremium.benchmarkPrice);
    assertEq(timestampsData[0], reportStructBlockPremium.observationsTimestamp);
    assertEq(timestampsData[1], reportStructBasic.observationsTimestamp);
    assertEq(timestampsData[2], reportStructPremium.observationsTimestamp);
  }
}
