// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice IDataFeedsRouter
/// Routes requests for feed data to the correct registry, including
/// benchmark values and full reports as separate endpoints. Fees for
/// for requests are delegated to a FeeManager. Users may also request
/// upkeep for a given feed via a set of upkeeps through this contract.
interface IDataFeedsRouter {
  /// @notice Get the latest benchmark for a set of feeds. For price feeds this will be the
  /// median/benchmark price, and for other feeds it will be the "main" data point of the
  /// report. Other report data should be fetched via getReports and decoded with the
  /// associated DataFeedsUtils.
  /// @param feedIds List of feed IDs
  /// @param billingData Encoded data for additional flexibility
  /// @return benchmarks
  /// @return observationTimestamps
  function getBenchmarks(
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable returns (int256[] memory benchmarks, uint256[] memory observationTimestamps);

  /// @notice Get the latest benchmark for a set of feeds. Only callable by authorized
  /// non-payable users.
  /// @param feedIds List of feed IDs
  /// @return benchmarks
  /// @return observationTimestamps
  function getBenchmarksNonbillable(
    bytes32[] calldata feedIds
  ) external view returns (int256[] memory benchmarks, uint256[] memory observationTimestamps);

  /// @notice Get latest reports for a set of feeds.
  /// @param feedIds List of feed IDs
  /// @param billingData Encoded data for additional flexibility
  /// @return reports
  /// @return observationTimestamps
  function getReports(
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable returns (bytes[] memory reports, uint256[] memory observationTimestamps);

  /// @notice Get latest reports for a set of feeds. Only callable by authorized
  /// non-payable users.
  /// @param feedIds List of feed IDs
  /// @return reports
  /// @return observationTimestamps
  function getReportsNonbillable(
    bytes32[] calldata feedIds
  ) external view returns (bytes[] memory reports, uint256[] memory observationTimestamps);

  /// @notice Get metadata associated with a set of feeds.
  /// @param feedIds List of feed IDs
  /// @return descriptions List of feed descriptions/names
  function getDescriptions(bytes32[] calldata feedIds) external view returns (string[] memory descriptions);

  /// @notice Allow a user to request upkeeps for a set of feeds on a set of upkeeps.
  /// @param feedIds List of feed IDs
  /// @param billingData Encoded data for additional flexibility
  function requestUpkeep(bytes32[] calldata feedIds, bytes calldata billingData) external payable;

  /// @notice For each feed in feedIds, fully configure the feed
  /// @param feedIds List of feed IDs
  /// @param descriptions List of feed descriptions/names
  /// @param configId Feed config level to set for all feeds in feedIds list (
  /// specifies deviation and heartbeat thresholds)
  /// @param upkeep Upkeep given permission to submit reports for the list of feeds
  /// @param registryAddress Registry that will manage these feeds
  /// @param feeConfigId Fee config for all feeds in the list (specifies USD fee
  /// for each user request type)
  function configureFeeds(
    bytes32[] calldata feedIds,
    string[] calldata descriptions,
    bytes32 configId,
    address upkeep,
    address registryAddress,
    bytes32 feeConfigId
  ) external;

  /// @notice For each feed in feedIds, update the Registry address.
  /// @param feedIds List of feed IDs
  /// @param registryAddress Registry that will manage these feeds
  function setRegistry(bytes32[] calldata feedIds, address registryAddress) external;

  /// @notice Get the Registry address for a given feed
  /// @param feedId Feed ID
  /// @return registryAddress Registry that manages this feed
  function getRegistry(bytes32 feedId) external view returns (address registryAddress);

  /// @notice Set the FeeManager contract address.
  /// @param feeManagerAddress FeeManager that will manage billing
  function setFeeManager(address feeManagerAddress) external;

  /// @notice Fetch the FeeManager contract address. This can be used to call
  /// functions to determine the expected fee associated with certain operations.
  /// @return feeManagerAddress FeeManager that manages billing
  function getFeeManager() external view returns (address feeManagerAddress);

  /// @notice Add user to nonpayable list
  /// @param user User address
  function addNonbillableUser(address user) external;

  /// @notice Remove user from nonpayable list
  /// @param user User address
  function removeNonbillableUser(address user) external;
}
