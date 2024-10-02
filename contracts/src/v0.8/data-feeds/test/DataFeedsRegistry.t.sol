// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {DataFeedsRegistry} from "../dev/DataFeedsRegistry.sol";
import {IDataStreamsFeeManager} from "../dev/interfaces/IDataStreamsFeeManager.sol";
import {IDataStreamsVerifierProxy} from "../dev/interfaces/IDataStreamsVerifierProxy.sol";
import {DataFeedsTestSetup} from "./DataFeedsTestSetup.t.sol";
import {DataFeedsBase} from "../dev/DataFeedsBase.sol";

contract DataFeedsRegistryHarness is DataFeedsRegistry {
  constructor(
    address linkAddress,
    address router,
    address verifierProxy
  ) DataFeedsRegistry(linkAddress, router, verifierProxy) {}

  function exposed_hasDuplicates(bytes32[] calldata elements) external pure returns (bool hasDuplicate) {
    return _hasDuplicates(elements);
  }
}

contract DataFeedsRegistryTest is DataFeedsTestSetup {
  using DataFeedsBase for bytes[];
  address internal constant ROUTER = address(10001);
  address internal constant UPKEEP = address(10002);
  address internal constant DS_VERIFIER_PROXY = address(10003);
  address internal constant DS_FEE_MANAGER = address(10004);
  address internal constant DS_REWARD_MANAGER = address(10005);
  address internal constant USER = address(10006);
  uint256 internal constant DEVIATION_THRESHOLD = 1e16; // 1%
  uint256 internal constant STALENESS_SECONDS = 60; // 1 minute
  uint256 internal constant LINK_FEE = 11111;

  bytes32[] internal CONFIG_IDS = [bytes32("3"), bytes32("4")];
  uint256[] internal CONFIG_DEVIATION_THRESHOLDS = [1e16, 2e16];
  uint256[] internal CONFIG_STALENESS_SECONDS = [60, 120];

  bytes32[] internal FEED_IDS = [reportStructBasic.feedId, reportStructPremium.feedId];

  string[] internal FEED_DESCRIPTIONS = ["Feed 1", "Feed 2"];

  event FeedConfigIdUpdated(bytes32 feedId, bytes32 configId);
  event FeedConfigSet(bytes32 configId, uint256 deviationThreshold, uint256 stalenessSeconds);
  event FeedDescriptionUpdated(bytes32 feedId, string description);
  event FeedRemoved(bytes32 feedId);
  event FeedSet(bytes32 feedId, string description, bytes32 configId, address upkeep);
  event FeedUpdated(bytes32 feedId, uint256 timestamp, int256 benchmark, bytes report);
  event UpkeepRequested(bytes32 feedId);
  event UpkeepUpdated(bytes32 feedId, address upkeep);

  DataFeedsRegistryHarness internal dataFeedsRegistry;

  function setUp() public override {
    DataFeedsTestSetup.setUp();

    dataFeedsRegistry = new DataFeedsRegistryHarness(address(link), ROUTER, DS_VERIFIER_PROXY);

    dataFeedsRegistry.setFeedConfigs(CONFIG_IDS, CONFIG_DEVIATION_THRESHOLDS, CONFIG_STALENESS_SECONDS);

    dataFeedsRegistry.setFeeds(FEED_IDS, FEED_DESCRIPTIONS, CONFIG_IDS[0], UPKEEP);
  }

  function test_setFeeds() public {
    bytes32[] memory feedIds = new bytes32[](2);
    feedIds[0] = bytes32("3");
    feedIds[1] = bytes32("4");

    string[] memory descriptions = new string[](2);
    descriptions[0] = "Feed 3";
    descriptions[1] = "Feed 4";

    (
      string[] memory oldDescriptions,
      bytes32[] memory oldConfigIds,
      uint256[] memory oldDeviationThresholds,
      uint256[] memory oldStalenessSeconds,
      bool[] memory oldUpkeepsRequested
    ) = dataFeedsRegistry.getFeedMetadata(feedIds);

    assertEq(oldDescriptions[0], "");
    assertEq(oldDescriptions[1], "");
    assertEq(oldConfigIds[0], bytes32(""));
    assertEq(oldConfigIds[1], bytes32(""));
    assertEq(oldDeviationThresholds[0], 0);
    assertEq(oldDeviationThresholds[1], 0);
    assertEq(oldStalenessSeconds[0], 0);
    assertEq(oldStalenessSeconds[1], 0);
    assertEq(oldUpkeepsRequested[0], false);
    assertEq(oldUpkeepsRequested[1], false);

    bytes32[] memory oldFeedIds = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(oldFeedIds[0], FEED_IDS[0]);
    assertEq(oldFeedIds[1], FEED_IDS[1]);
    assertEq(oldFeedIds.length, FEED_IDS.length);

    vm.expectEmit();
    emit FeedSet(feedIds[0], descriptions[0], CONFIG_IDS[1], UPKEEP);

    vm.expectEmit();
    emit FeedSet(feedIds[1], descriptions[1], CONFIG_IDS[1], UPKEEP);

    dataFeedsRegistry.setFeeds(feedIds, descriptions, CONFIG_IDS[1], UPKEEP);

    (
      string[] memory newDescriptions,
      bytes32[] memory newConfigIds,
      uint256[] memory newDeviationThresholds,
      uint256[] memory newStalenessSeconds,
      bool[] memory newUpkeepsRequested
    ) = dataFeedsRegistry.getFeedMetadata(feedIds);

    assertEq(newDescriptions[0], descriptions[0]);
    assertEq(newDescriptions[1], descriptions[1]);
    assertEq(newConfigIds[0], CONFIG_IDS[1]);
    assertEq(newConfigIds[1], CONFIG_IDS[1]);
    assertEq(newDeviationThresholds[0], CONFIG_DEVIATION_THRESHOLDS[1]);
    assertEq(newDeviationThresholds[1], CONFIG_DEVIATION_THRESHOLDS[1]);
    assertEq(newStalenessSeconds[0], CONFIG_STALENESS_SECONDS[1]);
    assertEq(newStalenessSeconds[1], CONFIG_STALENESS_SECONDS[1]);
    assertEq(newUpkeepsRequested[0], false);
    assertEq(newUpkeepsRequested[1], false);

    bytes32[] memory newFeedIds = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(newFeedIds[0], FEED_IDS[0]);
    assertEq(newFeedIds[1], FEED_IDS[1]);
    assertEq(newFeedIds[2], feedIds[0]);
    assertEq(newFeedIds[3], feedIds[1]);
    assertEq(newFeedIds.length, FEED_IDS.length + feedIds.length);
  }

  function test_setFeedsRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnauthorizedRouterOperation.selector));

    dataFeedsRegistry.setFeeds(FEED_IDS, FEED_DESCRIPTIONS, CONFIG_IDS[0], UPKEEP);
  }

  function test_setFeedsRevertUnequalArrayLengths() public {
    string[] memory descriptions = new string[](1);
    descriptions[0] = FEED_DESCRIPTIONS[0];

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnequalArrayLengths.selector));

    dataFeedsRegistry.setFeeds(FEED_IDS, descriptions, CONFIG_IDS[0], UPKEEP);
  }

  function test_setFeedsRevertDuplicateFeedIds() public {
    bytes32[] memory feedIds = new bytes32[](2);
    feedIds[0] = FEED_IDS[0];
    feedIds[1] = FEED_IDS[0];

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.DuplicateFeedIds.selector));

    dataFeedsRegistry.setFeeds(feedIds, FEED_DESCRIPTIONS, CONFIG_IDS[0], UPKEEP);
  }

  function test_setFeedsRevertInvalidUpkeep() public {
    address invalidUpkeep = address(0);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.InvalidUpkeep.selector));

    dataFeedsRegistry.setFeeds(FEED_IDS, FEED_DESCRIPTIONS, CONFIG_IDS[0], invalidUpkeep);
  }

  function test_setFeedsRevertFeedExists() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.FeedExists.selector, FEED_IDS[0]));

    dataFeedsRegistry.setFeeds(FEED_IDS, FEED_DESCRIPTIONS, CONFIG_IDS[0], UPKEEP);
  }

  function test_removeFeeds() public {
    bytes32[] memory feedIdsToRemove = new bytes32[](1);
    feedIdsToRemove[0] = FEED_IDS[0];

    (
      string[] memory oldDescriptions,
      bytes32[] memory oldConfigIds,
      uint256[] memory oldDeviationThresholds,
      uint256[] memory oldStalenessSeconds,
      bool[] memory oldUpkeepsRequested
    ) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(oldDescriptions[0], FEED_DESCRIPTIONS[0]);
    assertEq(oldDescriptions[1], FEED_DESCRIPTIONS[1]);
    assertEq(oldConfigIds[0], CONFIG_IDS[0]);
    assertEq(oldConfigIds[1], CONFIG_IDS[0]);
    assertEq(oldDeviationThresholds[0], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(oldDeviationThresholds[1], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(oldStalenessSeconds[0], CONFIG_STALENESS_SECONDS[0]);
    assertEq(oldStalenessSeconds[1], CONFIG_STALENESS_SECONDS[0]);
    assertEq(oldUpkeepsRequested[0], false);
    assertEq(oldUpkeepsRequested[1], false);

    bytes32[] memory oldFeedIds = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(oldFeedIds[0], FEED_IDS[0]);
    assertEq(oldFeedIds[1], FEED_IDS[1]);
    assertEq(oldFeedIds.length, FEED_IDS.length);

    vm.expectEmit();
    emit FeedRemoved(feedIdsToRemove[0]);

    dataFeedsRegistry.removeFeeds(feedIdsToRemove);

    (
      string[] memory newDescriptions,
      bytes32[] memory newConfigIds,
      uint256[] memory newDeviationThresholds,
      uint256[] memory newStalenessSeconds,
      bool[] memory newUpkeepsRequested
    ) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(newDescriptions[0], "");
    assertEq(newDescriptions[1], FEED_DESCRIPTIONS[1]);
    assertEq(newConfigIds[0], bytes32(""));
    assertEq(newConfigIds[1], CONFIG_IDS[0]);
    assertEq(newDeviationThresholds[0], 0);
    assertEq(newDeviationThresholds[1], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(newStalenessSeconds[0], 0);
    assertEq(newStalenessSeconds[1], CONFIG_STALENESS_SECONDS[0]);
    assertEq(newUpkeepsRequested[0], false);
    assertEq(newUpkeepsRequested[1], false);

    bytes32[] memory newFeedIds = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(newFeedIds[0], FEED_IDS[1]);
    assertEq(newFeedIds.length, FEED_IDS.length - feedIdsToRemove.length);
  }

  function test_removeFeedsRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert("Only callable by owner");

    dataFeedsRegistry.removeFeeds(FEED_IDS);
  }

  function test_removeFeedsRevertDuplicateFeedIds() public {
    bytes32[] memory feedIdsToRemove = new bytes32[](2);
    feedIdsToRemove[0] = FEED_IDS[0];
    feedIdsToRemove[1] = FEED_IDS[0];

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.DuplicateFeedIds.selector));

    dataFeedsRegistry.removeFeeds(feedIdsToRemove);
  }

  function test_removeFeedsRevertFeedNotConfigured() public {
    bytes32[] memory feedIdsToRemove = new bytes32[](1);
    feedIdsToRemove[0] = bytes32("3");

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.FeedNotConfigured.selector, feedIdsToRemove[0]));

    dataFeedsRegistry.removeFeeds(feedIdsToRemove);
  }

  function test_setFeedConfigs() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("5");
    configIds[1] = bytes32("6");

    (uint256[] memory oldDeviationThresholds, uint256[] memory oldStalenessSeconds) = dataFeedsRegistry.getFeedConfigs(
      configIds
    );

    assertEq(oldDeviationThresholds[0], 0);
    assertEq(oldDeviationThresholds[1], 0);
    assertEq(oldStalenessSeconds[0], 0);
    assertEq(oldStalenessSeconds[1], 0);

    uint256[] memory deviationThresholds = new uint256[](2);
    deviationThresholds[0] = 2e16;
    deviationThresholds[1] = 1e17;

    uint256[] memory stalenessSeconds = new uint256[](2);
    stalenessSeconds[0] = 30;
    stalenessSeconds[1] = 300;

    vm.expectEmit();
    emit FeedConfigSet(configIds[0], deviationThresholds[0], stalenessSeconds[0]);

    vm.expectEmit();
    emit FeedConfigSet(configIds[1], deviationThresholds[1], stalenessSeconds[1]);

    dataFeedsRegistry.setFeedConfigs(configIds, deviationThresholds, stalenessSeconds);

    (uint256[] memory newDeviationThresholds, uint256[] memory newStalenessSeconds) = dataFeedsRegistry.getFeedConfigs(
      configIds
    );

    assertEq(newDeviationThresholds[0], deviationThresholds[0]);
    assertEq(newDeviationThresholds[1], deviationThresholds[1]);
    assertEq(newStalenessSeconds[0], stalenessSeconds[0]);
    assertEq(newStalenessSeconds[1], stalenessSeconds[1]);
  }

  function test_setFeedConfigsRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert("Only callable by owner");

    dataFeedsRegistry.setFeedConfigs(CONFIG_IDS, CONFIG_DEVIATION_THRESHOLDS, CONFIG_STALENESS_SECONDS);
  }

  function test_setFeedConfigsRevertUnequalArrayLengthsDeviationThresholds() public {
    uint256[] memory deviationThresholds = new uint256[](1);
    deviationThresholds[0] = CONFIG_DEVIATION_THRESHOLDS[0];

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnequalArrayLengths.selector));

    dataFeedsRegistry.setFeedConfigs(CONFIG_IDS, deviationThresholds, CONFIG_STALENESS_SECONDS);
  }

  function test_setFeedConfigsRevertUnequalArrayLengthsStalenessSeconds() public {
    uint256[] memory stalenessSeconds = new uint256[](1);
    stalenessSeconds[0] = CONFIG_STALENESS_SECONDS[0];

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnequalArrayLengths.selector));

    dataFeedsRegistry.setFeedConfigs(CONFIG_IDS, CONFIG_DEVIATION_THRESHOLDS, stalenessSeconds);
  }

  function test_updateDescriptions() public {
    string[] memory descriptions = new string[](2);
    descriptions[0] = "New Feed 1 Description";
    descriptions[1] = "New Feed 2 Description";

    (string[] memory oldDescriptions, , , , ) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(oldDescriptions[0], FEED_DESCRIPTIONS[0]);
    assertEq(oldDescriptions[1], FEED_DESCRIPTIONS[1]);

    vm.expectEmit();
    emit FeedDescriptionUpdated(FEED_IDS[0], descriptions[0]);

    vm.expectEmit();
    emit FeedDescriptionUpdated(FEED_IDS[1], descriptions[1]);

    dataFeedsRegistry.updateDescriptions(FEED_IDS, descriptions);

    (string[] memory newDescriptions, , , , ) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(newDescriptions[0], descriptions[0]);
    assertEq(newDescriptions[1], descriptions[1]);
  }

  function test_updateDescriptionsRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert("Only callable by owner");

    dataFeedsRegistry.updateDescriptions(FEED_IDS, FEED_DESCRIPTIONS);
  }

  function test_updateDescriptionsRevertUnequalArrayLengths() public {
    bytes32[] memory feedIds = new bytes32[](1);
    feedIds[0] = FEED_IDS[0];

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnequalArrayLengths.selector));

    dataFeedsRegistry.updateDescriptions(feedIds, FEED_DESCRIPTIONS);
  }

  function test_updateFeedConfigId() public {
    (
      ,
      bytes32[] memory oldConfigIds,
      uint256[] memory oldDeviationThresholds,
      uint256[] memory oldStalenessSeconds,

    ) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(oldConfigIds[0], CONFIG_IDS[0]);
    assertEq(oldConfigIds[1], CONFIG_IDS[0]);
    assertEq(oldDeviationThresholds[0], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(oldDeviationThresholds[1], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(oldStalenessSeconds[0], CONFIG_STALENESS_SECONDS[0]);
    assertEq(oldStalenessSeconds[1], CONFIG_STALENESS_SECONDS[0]);

    vm.expectEmit();
    emit FeedConfigIdUpdated(FEED_IDS[0], CONFIG_IDS[1]);

    vm.expectEmit();
    emit FeedConfigIdUpdated(FEED_IDS[1], CONFIG_IDS[1]);

    dataFeedsRegistry.updateFeedConfigId(FEED_IDS, CONFIG_IDS[1]);

    (
      ,
      bytes32[] memory newConfigIds,
      uint256[] memory newDeviationThresholds,
      uint256[] memory newStalenessSeconds,

    ) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(newConfigIds[0], CONFIG_IDS[1]);
    assertEq(newConfigIds[1], CONFIG_IDS[1]);
    assertEq(newDeviationThresholds[0], CONFIG_DEVIATION_THRESHOLDS[1]);
    assertEq(newDeviationThresholds[1], CONFIG_DEVIATION_THRESHOLDS[1]);
    assertEq(newStalenessSeconds[0], CONFIG_STALENESS_SECONDS[1]);
    assertEq(newStalenessSeconds[1], CONFIG_STALENESS_SECONDS[1]);
  }

  function test_updateFeedConfigIdRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert("Only callable by owner");

    dataFeedsRegistry.updateFeedConfigId(FEED_IDS, CONFIG_IDS[1]);
  }

  function test_updateUpkeep() public {
    address newUpkeep = address(10005);

    bytes32[] memory oldFeedIdsPrevUpkeep = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(oldFeedIdsPrevUpkeep[0], FEED_IDS[0]);
    assertEq(oldFeedIdsPrevUpkeep[1], FEED_IDS[1]);
    assertEq(oldFeedIdsPrevUpkeep.length, FEED_IDS.length);

    bytes32[] memory oldFeedIdsNewUpkeep = dataFeedsRegistry.getUpkeepFeedIds(newUpkeep);

    assertEq(oldFeedIdsNewUpkeep.length, 0);

    vm.expectEmit();
    emit UpkeepUpdated(FEED_IDS[0], newUpkeep);

    vm.expectEmit();
    emit UpkeepUpdated(FEED_IDS[1], newUpkeep);

    dataFeedsRegistry.updateUpkeep(FEED_IDS, newUpkeep);

    bytes32[] memory newFeedIdsPrevUpkeep = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(newFeedIdsPrevUpkeep.length, 0);

    bytes32[] memory newFeedIdsNewUpkeep = dataFeedsRegistry.getUpkeepFeedIds(newUpkeep);

    assertEq(newFeedIdsNewUpkeep[0], FEED_IDS[0]);
    assertEq(newFeedIdsNewUpkeep[1], FEED_IDS[1]);
    assertEq(newFeedIdsNewUpkeep.length, FEED_IDS.length);
  }

  function test_updateUpkeepRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert("Only callable by owner");

    dataFeedsRegistry.updateUpkeep(FEED_IDS, UPKEEP);
  }

  function test_updateUpkeepRevertInvalidUpkeep() public {
    address invalidUpkeep = address(0);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.InvalidUpkeep.selector));

    dataFeedsRegistry.updateUpkeep(FEED_IDS, invalidUpkeep);
  }

  function test_requestUpkeep() public {
    (, , , , bool[] memory oldUpkeepsRequested) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(oldUpkeepsRequested[0], false);
    assertEq(oldUpkeepsRequested[1], false);

    vm.expectEmit();
    emit UpkeepRequested(FEED_IDS[0]);

    vm.expectEmit();
    emit UpkeepRequested(FEED_IDS[1]);

    vm.startPrank(ROUTER);

    dataFeedsRegistry.requestUpkeep(FEED_IDS);

    (, , , , bool[] memory newUpkeepsRequested) = dataFeedsRegistry.getFeedMetadata(FEED_IDS);

    assertEq(newUpkeepsRequested[0], true);
    assertEq(newUpkeepsRequested[1], true);
  }

  function test_requestUpkeepRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnauthorizedRouterOperation.selector));

    dataFeedsRegistry.requestUpkeep(FEED_IDS);
  }

  function test_performUpkeep() public {
    bytes[] memory reports = new bytes[](2);
    reports[0] = reportBasic;
    reports[1] = reportPremium;

    bytes[] memory reportData = reports.getReportData();

    IDataStreamsFeeManager.Asset memory verificationFee = IDataStreamsFeeManager.Asset({
      assetAddress: address(link),
      amount: LINK_FEE
    });

    vm.mockCall(
      DS_VERIFIER_PROXY,
      abi.encodeWithSelector(IDataStreamsVerifierProxy.s_feeManager.selector),
      abi.encode(DS_FEE_MANAGER)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(
        IDataStreamsFeeManager.getFeeAndReward.selector,
        address(dataFeedsRegistry),
        reportData[0],
        address(link)
      ),
      abi.encode(verificationFee, verificationFee, 0)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(
        IDataStreamsFeeManager.getFeeAndReward.selector,
        address(dataFeedsRegistry),
        reportData[1],
        address(link)
      ),
      abi.encode(verificationFee, verificationFee, 0)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(IDataStreamsFeeManager.i_rewardManager.selector),
      abi.encode(DS_REWARD_MANAGER)
    );

    vm.mockCall(
      DS_VERIFIER_PROXY,
      abi.encodeWithSelector(IDataStreamsVerifierProxy.verifyBulk.selector, reports, abi.encode(address(link))),
      abi.encode(reportData)
    );

    bytes memory performData = abi.encode(reports, "");

    vm.expectEmit();
    emit FeedUpdated(
      reportStructBasic.feedId,
      reportStructBasic.observationsTimestamp,
      reportStructBasic.benchmarkPrice,
      reportData[0]
    );

    vm.expectEmit();
    emit FeedUpdated(
      reportStructPremium.feedId,
      reportStructPremium.observationsTimestamp,
      reportStructPremium.benchmarkPrice,
      reportData[1]
    );

    vm.startPrank(UPKEEP);

    dataFeedsRegistry.performUpkeep(performData);

    vm.startPrank(ROUTER);

    (int256[] memory benchmarks, uint256[] memory observationTimestamps) = dataFeedsRegistry.getBenchmarks(FEED_IDS);

    assertEq(benchmarks[0], reportStructBasic.benchmarkPrice);
    assertEq(benchmarks[1], reportStructPremium.benchmarkPrice);
    assertEq(observationTimestamps[0], reportStructBasic.observationsTimestamp);
    assertEq(observationTimestamps[1], reportStructPremium.observationsTimestamp);
  }

  function test_getBenchmarks() public {
    bytes[] memory reports = new bytes[](2);
    reports[0] = reportBasic;
    reports[1] = reportPremium;

    bytes[] memory reportData = reports.getReportData();

    IDataStreamsFeeManager.Asset memory verificationFee = IDataStreamsFeeManager.Asset({
      assetAddress: address(link),
      amount: LINK_FEE
    });

    vm.mockCall(
      DS_VERIFIER_PROXY,
      abi.encodeWithSelector(IDataStreamsVerifierProxy.s_feeManager.selector),
      abi.encode(DS_FEE_MANAGER)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(
        IDataStreamsFeeManager.getFeeAndReward.selector,
        address(dataFeedsRegistry),
        reportData[0],
        address(link)
      ),
      abi.encode(verificationFee, verificationFee, 0)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(
        IDataStreamsFeeManager.getFeeAndReward.selector,
        address(dataFeedsRegistry),
        reportData[1],
        address(link)
      ),
      abi.encode(verificationFee, verificationFee, 0)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(IDataStreamsFeeManager.i_rewardManager.selector),
      abi.encode(DS_REWARD_MANAGER)
    );

    vm.mockCall(
      DS_VERIFIER_PROXY,
      abi.encodeWithSelector(IDataStreamsVerifierProxy.verifyBulk.selector, reports, abi.encode(address(link))),
      abi.encode(reportData)
    );

    bytes memory performData = abi.encode(reports, "");

    vm.startPrank(UPKEEP);

    dataFeedsRegistry.performUpkeep(performData);

    vm.startPrank(ROUTER);

    (int256[] memory benchmarks, uint256[] memory observationTimestamps) = dataFeedsRegistry.getBenchmarks(FEED_IDS);

    assertEq(benchmarks[0], reportStructBasic.benchmarkPrice);
    assertEq(benchmarks[1], reportStructPremium.benchmarkPrice);
    assertEq(observationTimestamps[0], reportStructBasic.observationsTimestamp);
    assertEq(observationTimestamps[1], reportStructPremium.observationsTimestamp);
  }

  function test_getBenchmarksRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnauthorizedDataFetch.selector));

    dataFeedsRegistry.getBenchmarks(FEED_IDS);
  }

  function test_getReports() public {
    bytes[] memory reports = new bytes[](2);
    reports[0] = reportBasic;
    reports[1] = reportPremium;

    bytes[] memory reportData = reports.getReportData();

    IDataStreamsFeeManager.Asset memory verificationFee = IDataStreamsFeeManager.Asset({
      assetAddress: address(link),
      amount: LINK_FEE
    });

    vm.mockCall(
      DS_VERIFIER_PROXY,
      abi.encodeWithSelector(IDataStreamsVerifierProxy.s_feeManager.selector),
      abi.encode(DS_FEE_MANAGER)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(
        IDataStreamsFeeManager.getFeeAndReward.selector,
        address(dataFeedsRegistry),
        reportData[0],
        address(link)
      ),
      abi.encode(verificationFee, verificationFee, 0)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(
        IDataStreamsFeeManager.getFeeAndReward.selector,
        address(dataFeedsRegistry),
        reportData[1],
        address(link)
      ),
      abi.encode(verificationFee, verificationFee, 0)
    );

    vm.mockCall(
      DS_FEE_MANAGER,
      abi.encodeWithSelector(IDataStreamsFeeManager.i_rewardManager.selector),
      abi.encode(DS_REWARD_MANAGER)
    );

    vm.mockCall(
      DS_VERIFIER_PROXY,
      abi.encodeWithSelector(IDataStreamsVerifierProxy.verifyBulk.selector, reports, abi.encode(address(link))),
      abi.encode(reportData)
    );

    bytes memory performData = abi.encode(reports, "");

    vm.startPrank(UPKEEP);

    dataFeedsRegistry.performUpkeep(performData);

    vm.startPrank(ROUTER);

    (bytes[] memory reportsData, uint256[] memory observationTimestamps) = dataFeedsRegistry.getReports(FEED_IDS);

    assertEq(reportsData[0], reportData[0]);
    assertEq(reportsData[1], reportData[1]);
    assertEq(observationTimestamps[0], reportStructBasic.observationsTimestamp);
    assertEq(observationTimestamps[1], reportStructPremium.observationsTimestamp);
  }

  function test_getReportsRevertUnauthorized() public {
    vm.startPrank(USER);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsRegistry.UnauthorizedDataFetch.selector));

    dataFeedsRegistry.getReports(FEED_IDS);
  }

  function test_getFeedMetadata() public {
    bytes32[] memory feedIds = new bytes32[](3);
    feedIds[0] = FEED_IDS[0];
    feedIds[1] = FEED_IDS[1];
    feedIds[2] = bytes32("3");

    (
      string[] memory descriptions,
      bytes32[] memory configIds,
      uint256[] memory deviationThresholds,
      uint256[] memory stalenessSeconds,
      bool[] memory upkeepsRequested
    ) = dataFeedsRegistry.getFeedMetadata(feedIds);

    assertEq(descriptions[0], "Feed 1");
    assertEq(descriptions[1], "Feed 2");
    assertEq(descriptions[2], "");
    assertEq(configIds[0], CONFIG_IDS[0]);
    assertEq(configIds[1], CONFIG_IDS[0]);
    assertEq(configIds[2], bytes32(""));
    assertEq(deviationThresholds[0], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(deviationThresholds[1], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(deviationThresholds[2], 0);
    assertEq(stalenessSeconds[0], CONFIG_STALENESS_SECONDS[0]);
    assertEq(stalenessSeconds[1], CONFIG_STALENESS_SECONDS[0]);
    assertEq(stalenessSeconds[2], 0);
    assertEq(upkeepsRequested[0], false);
    assertEq(upkeepsRequested[1], false);
    assertEq(upkeepsRequested[2], false);
  }

  function test_getFeedConfigs() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = CONFIG_IDS[0];
    configIds[1] = bytes32("5");

    (uint256[] memory deviationThresholds, uint256[] memory stalenessSeconds) = dataFeedsRegistry.getFeedConfigs(
      configIds
    );

    assertEq(deviationThresholds[0], CONFIG_DEVIATION_THRESHOLDS[0]);
    assertEq(deviationThresholds[1], 0);
    assertEq(stalenessSeconds[0], CONFIG_STALENESS_SECONDS[0]);
    assertEq(stalenessSeconds[1], 0);
  }

  function test_getUpkeepFeedIds() public {
    bytes32[] memory feedIds = dataFeedsRegistry.getUpkeepFeedIds(UPKEEP);

    assertEq(feedIds[0], FEED_IDS[0]);
    assertEq(feedIds[1], FEED_IDS[1]);
    assertEq(feedIds.length, FEED_IDS.length);
  }

  function test_getLinkAddress() public {
    address linkAddress = dataFeedsRegistry.getLinkAddress();

    assertEq(linkAddress, address(link));
  }

  function test_linkAvailableForPayment() public {
    uint256 LINK_MINT = 12345;
    link.mint(address(dataFeedsRegistry), LINK_MINT);

    int256 linkAvailable = dataFeedsRegistry.linkAvailableForPayment();

    assertEq(linkAvailable, int256(LINK_MINT));
  }

  function test_hasDuplicatesBytes32True() public {
    bytes32[] memory elements = new bytes32[](3);
    elements[0] = bytes32("0");
    elements[1] = bytes32("1");
    elements[2] = bytes32("1");

    bool hasDuplicates = dataFeedsRegistry.exposed_hasDuplicates(elements);

    assertEq(hasDuplicates, true);
  }

  function test_hasDuplicatesBytes32False() public {
    bytes32[] memory elements = new bytes32[](3);
    elements[0] = bytes32("0");
    elements[1] = bytes32("1");
    elements[2] = bytes32("2");

    bool hasDuplicates = dataFeedsRegistry.exposed_hasDuplicates(elements);

    assertEq(hasDuplicates, false);
  }

  function test_typeAndVersion() public {
    string memory typeAndVersion = dataFeedsRegistry.typeAndVersion();
    assertEq(typeAndVersion, "DataFeedsRegistry 1.0.0", "typeAndVersion should match expected value");
  }
}
