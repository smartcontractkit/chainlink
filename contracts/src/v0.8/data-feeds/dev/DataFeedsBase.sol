// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice DataFeedsBase
/// Utility library that handles certain operations for feeds.
library DataFeedsBase {
  bytes32 private constant SCHEMA_MASK = 0xffff000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant BLOCK_PREMIUM_SCHEMA = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant BASIC_SCHEMA = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant PREMIUM_SCHEMA = 0x0003000000000000000000000000000000000000000000000000000000000000;

  /// @notice https://github.com/smartcontractkit/chainlink/blob/develop/core/services/relay/evm/mercury/v1/types/types.go
  struct BlockPremiumSchema {
    bytes32 feedId; /// The feed ID the report has data for
    uint32 observationsTimestamp; /// The time the median value was observed on
    int192 benchmarkPrice; /// The median value agreed in an OCR round
    int192 bid; /// The best bid value agreed in an OCR round
    int192 ask; /// The best ask value agreed in an OCR round
    uint64 currentBlockNum; /// The upper bound of the block range the median value was observed within
    bytes32 currentBlockHash; /// The blockhash for the upper bound of block range (ensures correct blockchain)
    uint64 validFromBlockNum; /// The lower bound of the block range the median value was observed within
    uint64 currentBlockTimestamp; /// Upper timestamp for validity of report
  }

  /// @notice https://github.com/smartcontractkit/chainlink/blob/develop/core/services/relay/evm/mercury/v2/types/types.go
  struct BasicSchema {
    bytes32 feedId; /// The feed ID the report has data for
    uint32 validFromTimestamp; /// Lower timestamp for validity of report
    uint32 observationsTimestamp; /// The time the median value was observed on
    uint192 nativeFee; /// Base ETH/WETH fee to verify report
    uint192 linkFee; /// Base LINK fee to verify report
    uint32 expiresAt; /// Upper timestamp for validity of report
    int192 benchmarkPrice; /// The median value agreed in an OCR round
  }

  /// @notice https://github.com/smartcontractkit/chainlink/blob/develop/core/services/relay/evm/mercury/v3/types/types.go
  struct PremiumSchema {
    bytes32 feedId; /// The feed ID the report has data for
    uint32 validFromTimestamp; /// Lower timestamp for validity of report
    uint32 observationsTimestamp; /// The time the median value was observed on
    uint192 nativeFee; /// Base ETH/WETH fee to verify report
    uint192 linkFee; /// Base LINK fee to verify report
    uint32 expiresAt; /// Upper timestamp for validity of report
    int192 benchmarkPrice; /// The median value agreed in an OCR round
    int192 bid; /// The best bid value agreed in an OCR round
    int192 ask; /// The best ask value agreed in an OCR round
  }

  /// @notice Given a report, get the inner report data
  /// @param reports Data Streams reports
  /// @return data
  function getReportData(bytes[] memory reports) internal pure returns (bytes[] memory data) {
    bytes[] memory reportData = new bytes[](reports.length);

    for (uint256 i; i < reports.length; i++) {
      (, bytes memory thisReport, , , ) = abi.decode(reports[i], (bytes32[3], bytes, bytes32[], bytes32[], bytes32));
      reportData[i] = thisReport;
    }

    return reportData;
  }

  /// @notice Given a report, parse out the feed ID
  /// @param reportData Data Streams reports
  /// @return feedIdsData
  function getFeedIds(bytes[] memory reportData) internal pure returns (bytes32[] memory feedIdsData) {
    bytes32[] memory feedIds = new bytes32[](reportData.length);

    for (uint256 i; i < reportData.length; i++) {
      feedIds[i] = bytes32(reportData[i]);
    }

    return feedIds;
  }

  /// @notice Given a report, parse out the benchmark value (price, median, or other)
  /// and timestamp from within.
  /// @param reportData Data Streams report
  /// @return benchmarksData
  /// @return timestampsData
  function getBenchmarksAndTimestamps(
    bytes[] memory reportData
  ) internal pure returns (int256[] memory benchmarksData, uint256[] memory timestampsData) {
    int256[] memory benchmarks = new int256[](reportData.length);
    uint256[] memory timestamps = new uint256[](reportData.length);

    for (uint256 i; i < reportData.length; i++) {
      bytes memory thisReport = reportData[i];

      bytes32 reportSchema = SCHEMA_MASK & bytes32(thisReport);

      /// Check for Basic Schema
      if (reportSchema == BASIC_SCHEMA) {
        BasicSchema memory decodedReport = abi.decode(thisReport, (BasicSchema));

        benchmarks[i] = int256(decodedReport.benchmarkPrice);
        timestamps[i] = uint256(decodedReport.observationsTimestamp);

        /// Check for Premium Schema
      } else if (reportSchema == PREMIUM_SCHEMA) {
        PremiumSchema memory decodedReport = abi.decode(thisReport, (PremiumSchema));

        benchmarks[i] = int256(decodedReport.benchmarkPrice);
        timestamps[i] = uint256(decodedReport.observationsTimestamp);

        /// Default to Block Premium Schema
      } else {
        BlockPremiumSchema memory decodedReport = abi.decode(thisReport, (BlockPremiumSchema));

        benchmarks[i] = int256(decodedReport.benchmarkPrice);
        timestamps[i] = uint256(decodedReport.observationsTimestamp);
      }
    }
    return (benchmarks, timestamps);
  }
}
