// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {DataFeedsBase} from "./DataFeedsBase.sol";
import {IDataFeedsRegistry} from "./interfaces/IDataFeedsRegistry.sol";
import {IDataStreamsFeeManager} from "./interfaces/IDataStreamsFeeManager.sol";
import {IDataStreamsVerifierProxy} from "./interfaces/IDataStreamsVerifierProxy.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";
import {ILinkAvailable} from "../../automation/upkeeps/LinkAvailableBalanceMonitor.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";

contract DataFeedsRegistry is ConfirmedOwner, IDataFeedsRegistry, ILinkAvailable, TypeAndVersionInterface {
  using EnumerableSet for EnumerableSet.Bytes32Set;
  using DataFeedsBase for bytes[];

  string public constant override typeAndVersion = "DataFeedsRegistry 1.0.0";

  uint256 immutable MIN_GAS_FOR_PERFORM = 200_000; // TODO determine more accurate number

  struct Feed {
    string description; // Description of the feed
    bytes32 configId; // Identifier for the config parameters object
    address upkeep; // Upkeep associated with the feed
    bool upkeepRequested; // T/F to trigger manual upkeep regardless of thresholds
    int256 benchmark; // The latest benchmark. For price feeds this will be
    // the median/benchmark price, and for other feeds
    // it will be the "main" data point of the report
    bytes report; // The latest report for the feed
    uint256 observationTimestamp; // Observation timestamp of the latest report
  }

  struct Config {
    uint256 deviationThreshold; // Deviation threshold as a percentage of the benchmark value (10^18 = 100%)
    uint256 stalenessSeconds; // Heartbeat threshold in seconds
  }

  mapping(bytes32 feedId => Feed) private s_feeds;
  mapping(bytes32 configId => Config) private s_configs;
  mapping(address upkeep => EnumerableSet.Bytes32Set feedIds) private s_upkeepFeedIdSet;

  event FeedConfigIdUpdated(bytes32 feedId, bytes32 configId);
  event FeedConfigSet(bytes32 configId, uint256 deviationThreshold, uint256 stalenessSeconds);
  event FeedDescriptionUpdated(bytes32 feedId, string description);
  event FeedRemoved(bytes32 feedId);
  event FeedSet(bytes32 feedId, string description, bytes32 configId, address upkeep);
  event FeedUpdated(bytes32 feedId, uint256 timestamp, int256 benchmark, bytes report);
  event StaleReport(bytes32 feedId, uint256 latestTimestamp, uint256 reportTimestamp);
  event UnauthorizedUpkeep(bytes32 feedId, address upkeep, address sender);
  event UpkeepRequested(bytes32 feedId);
  event UpkeepUpdated(bytes32 feedId, address upkeep);

  error DuplicateFeedIds();
  error FeedExists(bytes32 feedId);
  error FeedNotConfigured(bytes32 feedId);
  error InvalidUpkeep();
  error UnauthorizedDataFetch();
  error UnauthorizedRouterOperation();
  error UnequalArrayLengths();

  address public immutable i_linkAddress;
  address public immutable i_router;
  address public immutable i_verifierProxy;

  modifier authorizeRouterOperation() {
    if (msg.sender != owner() && msg.sender != i_router) {
      revert UnauthorizedRouterOperation();
    }

    _;
  }

  modifier authorizeDataFetch(bytes32[] calldata feedIds) {
    bool isAuthorizedCaller = false;
    if (msg.sender == owner() || msg.sender == i_router || msg.sender == address(0)) {
      isAuthorizedCaller = true;
    } else {
      isAuthorizedCaller = true;

      EnumerableSet.Bytes32Set storage upkeepFeedIds = s_upkeepFeedIdSet[msg.sender];
      for (uint256 i; i < feedIds.length; i++) {
        if (!upkeepFeedIds.contains(feedIds[i])) {
          isAuthorizedCaller = false;
          break;
        }
      }
    }

    if (!isAuthorizedCaller) revert UnauthorizedDataFetch();

    _;
  }

  constructor(address linkAddress, address router, address verifierProxy) ConfirmedOwner(msg.sender) {
    i_linkAddress = linkAddress;
    i_router = router;
    i_verifierProxy = verifierProxy;
  }

  /// @inheritdoc IDataFeedsRegistry
  function setFeeds(
    bytes32[] calldata feedIds,
    string[] calldata descriptions,
    bytes32 configId,
    address upkeep
  ) external authorizeRouterOperation {
    if (_hasDuplicates(feedIds)) revert DuplicateFeedIds();

    if (feedIds.length != descriptions.length) revert UnequalArrayLengths();

    if (upkeep == address(0)) revert InvalidUpkeep();

    for (uint256 i; i < feedIds.length; i++) {
      if (s_feeds[feedIds[i]].upkeep != address(0)) revert FeedExists(feedIds[i]);

      s_feeds[feedIds[i]] = Feed({
        description: descriptions[i],
        configId: configId,
        upkeep: upkeep,
        upkeepRequested: false,
        benchmark: 0,
        report: hex"",
        observationTimestamp: 0
      });

      s_upkeepFeedIdSet[upkeep].add(feedIds[i]);

      emit FeedSet(feedIds[i], descriptions[i], configId, upkeep);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function removeFeeds(bytes32[] calldata feedIds) external onlyOwner {
    if (_hasDuplicates(feedIds)) revert DuplicateFeedIds();

    for (uint256 i; i < feedIds.length; i++) {
      if (s_feeds[feedIds[i]].upkeep == address(0)) revert FeedNotConfigured(feedIds[i]);

      s_upkeepFeedIdSet[s_feeds[feedIds[i]].upkeep].remove(feedIds[i]);

      delete s_feeds[feedIds[i]];

      emit FeedRemoved(feedIds[i]);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function setFeedConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata deviationThresholds,
    uint256[] calldata stalenessSeconds
  ) external onlyOwner {
    if (configIds.length != deviationThresholds.length || configIds.length != stalenessSeconds.length)
      revert UnequalArrayLengths();

    for (uint256 i; i < configIds.length; i++) {
      s_configs[configIds[i]] = Config({
        deviationThreshold: deviationThresholds[i],
        stalenessSeconds: stalenessSeconds[i]
      });

      emit FeedConfigSet(configIds[i], deviationThresholds[i], stalenessSeconds[i]);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function updateDescriptions(bytes32[] calldata feedIds, string[] calldata descriptions) external onlyOwner {
    if (feedIds.length != descriptions.length) revert UnequalArrayLengths();

    for (uint256 i; i < feedIds.length; i++) {
      s_feeds[feedIds[i]].description = descriptions[i];

      emit FeedDescriptionUpdated(feedIds[i], descriptions[i]);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function updateFeedConfigId(bytes32[] calldata feedIds, bytes32 configId) external onlyOwner {
    for (uint256 i; i < feedIds.length; i++) {
      s_feeds[feedIds[i]].configId = configId;

      emit FeedConfigIdUpdated(feedIds[i], configId);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function updateUpkeep(bytes32[] calldata feedIds, address upkeep) external onlyOwner {
    if (upkeep == address(0)) revert InvalidUpkeep();

    for (uint256 i; i < feedIds.length; i++) {
      address prevUpkeep = s_feeds[feedIds[i]].upkeep;
      s_upkeepFeedIdSet[prevUpkeep].remove(feedIds[i]);
      s_upkeepFeedIdSet[upkeep].add(feedIds[i]);

      s_feeds[feedIds[i]].upkeep = upkeep;

      emit UpkeepUpdated(feedIds[i], upkeep);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function requestUpkeep(bytes32[] calldata feedIds) external authorizeRouterOperation {
    for (uint256 i; i < feedIds.length; i++) {
      s_feeds[feedIds[i]].upkeepRequested = true;

      emit UpkeepRequested(feedIds[i]);
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function performUpkeep(bytes calldata performData) external {
    (bytes[] memory reports, ) = abi.decode(performData, (bytes[], bytes));

    // TODO add check for heartbeat threshold from last update (acts as a cooldown)

    bytes[] memory reportData = reports.getReportData();

    address feeManager = IDataStreamsVerifierProxy(i_verifierProxy).s_feeManager();

    uint256 totalFee = 0;

    for (uint256 i; i < reports.length; i++) {
      (IDataStreamsFeeManager.Asset memory fee, , ) = IDataStreamsFeeManager(feeManager).getFeeAndReward(
        address(this),
        reportData[i],
        i_linkAddress
      );

      totalFee += fee.amount;
    }

    address rewardManager = IDataStreamsFeeManager(feeManager).i_rewardManager();

    IERC20(i_linkAddress).approve(rewardManager, totalFee);

    // TODO filter out unauthorized upkeeps and stalereports before verifying reports
    IDataStreamsVerifierProxy(i_verifierProxy).verifyBulk(reports, abi.encode(i_linkAddress));

    bytes32[] memory feedIds = reportData.getFeedIds();
    (int256[] memory benchmarks, uint256[] memory timestamps) = reportData.getBenchmarksAndTimestamps();

    for (uint256 i; i < reports.length; i++) {
      bytes32 feedId = feedIds[i];
      if (s_feeds[feedId].upkeep != msg.sender) {
        emit UnauthorizedUpkeep(feedId, s_feeds[feedId].upkeep, msg.sender);
        continue;
      }

      if (s_feeds[feedId].observationTimestamp >= timestamps[i]) {
        emit StaleReport(feedId, s_feeds[feedId].observationTimestamp, timestamps[i]);
        continue;
      }

      s_feeds[feedId].benchmark = benchmarks[i];
      s_feeds[feedId].observationTimestamp = timestamps[i];
      s_feeds[feedId].report = reportData[i];
      s_feeds[feedId].upkeepRequested = false;

      emit FeedUpdated(
        feedIds[i],
        s_feeds[feedId].observationTimestamp,
        s_feeds[feedId].benchmark,
        s_feeds[feedId].report
      );

      if (gasleft() < MIN_GAS_FOR_PERFORM) {
        return;
      }
    }
  }

  /// @inheritdoc IDataFeedsRegistry
  function getBenchmarks(
    bytes32[] calldata feedIds
  )
    external
    view
    authorizeDataFetch(feedIds)
    returns (int256[] memory benchmarksData, uint256[] memory observationTimestampsData)
  {
    int256[] memory benchmarks = new int256[](feedIds.length);
    uint256[] memory observationTimestamps = new uint256[](feedIds.length);

    for (uint256 i; i < feedIds.length; i++) {
      benchmarks[i] = s_feeds[feedIds[i]].benchmark;
      observationTimestamps[i] = s_feeds[feedIds[i]].observationTimestamp;
    }

    return (benchmarks, observationTimestamps);
  }

  /// @inheritdoc IDataFeedsRegistry
  function getReports(
    bytes32[] calldata feedIds
  )
    external
    view
    authorizeDataFetch(feedIds)
    returns (bytes[] memory reportsData, uint256[] memory observationTimestampsData)
  {
    bytes[] memory reports = new bytes[](feedIds.length);
    uint256[] memory observationTimestamps = new uint256[](feedIds.length);

    for (uint256 i; i < feedIds.length; i++) {
      reports[i] = s_feeds[feedIds[i]].report;
      observationTimestamps[i] = s_feeds[feedIds[i]].observationTimestamp;
    }

    return (reports, observationTimestamps);
  }

  /// @inheritdoc IDataFeedsRegistry
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
    )
  {
    string[] memory descriptions = new string[](feedIds.length);
    bytes32[] memory configIds = new bytes32[](feedIds.length);
    uint256[] memory deviationThresholds = new uint256[](feedIds.length);
    uint256[] memory stalenessSeconds = new uint256[](feedIds.length);
    bool[] memory upkeepsRequested = new bool[](feedIds.length);

    for (uint256 i; i < feedIds.length; i++) {
      descriptions[i] = s_feeds[feedIds[i]].description;
      configIds[i] = s_feeds[feedIds[i]].configId;
      deviationThresholds[i] = s_configs[configIds[i]].deviationThreshold;
      stalenessSeconds[i] = s_configs[configIds[i]].stalenessSeconds;
      upkeepsRequested[i] = s_feeds[feedIds[i]].upkeepRequested;
    }

    return (descriptions, configIds, deviationThresholds, stalenessSeconds, upkeepsRequested);
  }

  /// @inheritdoc IDataFeedsRegistry
  function getFeedConfigs(
    bytes32[] calldata configIds
  ) external view returns (uint256[] memory deviationThresholdsData, uint256[] memory stalenessSecondsData) {
    uint256[] memory deviationThresholds = new uint256[](configIds.length);
    uint256[] memory stalenessSeconds = new uint256[](configIds.length);

    for (uint256 i; i < configIds.length; i++) {
      deviationThresholds[i] = s_configs[configIds[i]].deviationThreshold;
      stalenessSeconds[i] = s_configs[configIds[i]].stalenessSeconds;
    }

    return (deviationThresholds, stalenessSeconds);
  }

  /// @inheritdoc IDataFeedsRegistry
  function getUpkeepFeedIds(address upkeep) external view returns (bytes32[] memory feedIdsData) {
    uint256 setLength = s_upkeepFeedIdSet[upkeep].length();
    bytes32[] memory feedIds = new bytes32[](setLength);

    for (uint256 i; i < setLength; i++) {
      feedIds[i] = s_upkeepFeedIdSet[upkeep].at(i);
    }
    return feedIds;
  }

  /// @inheritdoc IDataFeedsRegistry
  function getLinkAddress() external view returns (address linkAddress) {
    return i_linkAddress;
  }

  function _hasDuplicates(bytes32[] calldata elements) internal pure returns (bool hasDuplicate) {
    for (uint256 i; i < elements.length; ) {
      for (uint256 j = i + 1; j < elements.length; ) {
        if (elements[i] == elements[j]) {
          return true;
        }
        unchecked {
          ++j;
        }
      }
      unchecked {
        ++i;
      }
    }
    return false;
  }

  /// @inheritdoc ILinkAvailable
  function linkAvailableForPayment() external view returns (int256 availableBalance) {
    return int256(IERC20(i_linkAddress).balanceOf(address(this)));
  }

  // TODO add withdraw method for ERC-20s and native
}
