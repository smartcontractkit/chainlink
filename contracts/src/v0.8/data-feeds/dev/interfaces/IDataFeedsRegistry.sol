// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice IDataFeedsRegistry
/// Responsible for storing the latest benchmark value and report
/// for every feed in the Data Feeds service. Also manages the upkeep
/// permissions for each feed. Automation upkeeps will fetch reports
/// that are sent to the registry, which will verify the received
/// reports against the Data Streams VerifierProxy contract (payment
/// required in LINK, hence inheriting ILinkAvailable). This
/// registry should also hold relevant metadata associated with the
/// feeds, such as the canonical feed name.

interface IDataFeedsRegistry {
  /// @notice For each feed in feedIds, add the feed config and associated upkeep.
  /// @param feedIds List of feed IDs
  /// @param descriptions List of feed descriptions/names
  /// @param configId Feed config level to set for all feeds in feedIds list
  /// @param upkeep Upkeep given permission to submit reports for the list of feeds
  function setFeeds(
    bytes32[] calldata feedIds,
    string[] calldata descriptions,
    bytes32 configId,
    address upkeep
  ) external;

  /// @notice Remove the List of feeds and all upkeeps for each.
  /// @param feedIds List of feed IDs
  function removeFeeds(bytes32[] calldata feedIds) external;

  /// @notice Set the parameters for a feed config level
  /// @param configIds Unique IDs for the config parameters
  /// @param deviationThresholds Deviation Thresholds as a percentage of the benchmark value (10^18 = 100%)
  /// @param stalenessSeconds Heartbeat Thresholds in seconds
  function setFeedConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata deviationThresholds,
    uint256[] calldata stalenessSeconds
  ) external;

  /// @notice Update the description associated with each feed.
  /// @param feedIds List of feed IDs
  /// @param descriptions List of feed descriptions/names
  function updateDescriptions(bytes32[] calldata feedIds, string[] calldata descriptions) external;

  /// @notice Update the feed config level associated with each feed.
  /// @param feedIds List of feed IDs
  /// @param configId Feed config level to set for all feeds in feedIds list
  function updateFeedConfigId(bytes32[] calldata feedIds, bytes32 configId) external;

  /// @notice Update the upkeep associated with a set of feeds.
  /// @param feedIds List of feed IDs
  /// @param upkeep Upkeep to give permission to submit reports for the List of feeds
  function updateUpkeep(bytes32[] calldata feedIds, address upkeep) external;

  /// @notice Allow a user to request upkeeps for a set of feeds.
  /// Will only execute for the upkeeps that are permissioned for a given feed.
  /// @param feedIds List of feed IDs
  function requestUpkeep(bytes32[] calldata feedIds) external;

  /// @notice Perform upkeep to update the latest reports for a set of feeds. Can
  /// only be called for the set of feeds the caller (upkeep) is permissioned for.
  /// @param performData // Encoded data passed in from the performUpkeep method on an upkeep
  function performUpkeep(bytes calldata performData) external;

  /// @notice Get the latest benchmark for a set of feeds. For price feeds this will be the
  /// median/benchmark price, and for other feeds it will be the "main" data point of the
  /// report. Other report data should be fetched via getReports and decoded with the
  /// associated DataFeedsUtils.
  /// @param feedIds List of feed IDs
  /// @return benchmarksData
  /// @return observationTimestampsData
  function getBenchmarks(
    bytes32[] calldata feedIds
  ) external view returns (int256[] memory benchmarksData, uint256[] memory observationTimestampsData);

  /// @notice Get latest reports for a set of feeds.
  /// @param feedIds List of feed IDs
  /// @return reportsData
  /// @return observationTimestampsData
  function getReports(
    bytes32[] calldata feedIds
  ) external view returns (bytes[] memory reportsData, uint256[] memory observationTimestampsData);

  /// @notice Get metadata associated with a set of feeds.
  /// @param feedIds List of feed IDs
  /// @return descriptionsData List of feed descriptions/names
  /// @return configIdsData List of feed config IDs (call getFeedConfigs to return threshold values)
  /// @return deviationThresholdsData Deviation Thresholds as a percentage of the benchmark value (10^18 = 100%)
  /// @return stalenessSecondsData Heartbeat Thresholds in seconds
  /// @return upkeepsRequestedData List of T/F flags whether an upkeep is requested or not
  function getFeedMetadata(
    bytes32[] calldata feedIds
  )
    external
    view
    returns (
      string[] memory descriptionsData,
      bytes32[] memory configIdsData,
      uint256[] memory deviationThresholdsData,
      uint256[] memory stalenessSecondsData,
      bool[] memory upkeepsRequestedData
    );

  /// @notice Get the data for feed config levels
  /// @param configIds Unique IDs for the config parameters
  /// @return deviationThresholdsData Deviation Thresholds as a percentage of the benchmark value (10^18 = 100%)
  /// @return stalenessSecondsData Heartbeat Thresholds in seconds
  function getFeedConfigs(
    bytes32[] calldata configIds
  ) external view returns (uint256[] memory deviationThresholdsData, uint256[] memory stalenessSecondsData);

  /// @notice Return LINK token ERC-20 address
  /// @return linkAddress address
  function getLinkAddress() external view returns (address linkAddress);

  /// @notice Get the feeds associated with an upkeep.
  /// @param upkeep Upkeep
  /// @return feedIdsData Feed IDs associated with the given upkeep
  function getUpkeepFeedIds(address upkeep) external view returns (bytes32[] memory feedIdsData);
}
