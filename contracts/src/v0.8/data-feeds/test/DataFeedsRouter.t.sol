pragma solidity ^0.8.19;

import {BaseTest} from "../../shared/test/BaseTest.t.sol";
import {DataFeedsRouter} from "../dev/DataFeedsRouter.sol";
import {IDataFeedsRegistry} from "../dev/interfaces/IDataFeedsRegistry.sol";
import {IDataFeedsFeeManager} from "../dev/interfaces/IDataFeedsFeeManager.sol";

struct Feed {
  bytes32 feedId;
  string description;
  string registry;
  int256 benchmark;
  bytes report;
  uint256 observationTimestamp;
}

contract DataFeedsRouterTest is BaseTest {
  RouterHarness internal router;

  mapping(string => address) internal registries;
  MockFeeManager internal feeManager;

  Feed[] internal allFeeds;
  bytes32[] internal allFeedIds;
  mapping(string => Feed) internal feeds;

  bytes internal constant billingData = bytes("");

  event FeeManagerSet(address indexed currentFeeManager, address indexed previousFeeManager);

  function setUp() public override {
    super.setUp();
    vm.deal(OWNER, 100 ether);

    router = new RouterHarness();

    // Define test Feeds
    Feed memory feedA = Feed({
      feedId: 0x7465737431000000000000000000000000000000000000000000000000000000,
      description: "a",
      registry: "a",
      benchmark: 100,
      report: bytes("feedA"),
      observationTimestamp: 100000
    });

    Feed memory feedB = Feed({
      feedId: 0x3274657374000000000000000000000000000000000000000000000000000000,
      description: "b",
      registry: "a",
      benchmark: 200,
      report: bytes("feedB"),
      observationTimestamp: 200000
    });

    Feed memory feedC = Feed({
      feedId: 0x3362616c00000000000000000000000000000000000000000000000000000000,
      description: "c",
      registry: "b",
      benchmark: 300,
      report: bytes("feedC"),
      observationTimestamp: 300000
    });

    allFeeds.push(feedA);
    allFeeds.push(feedB);
    allFeeds.push(feedC);

    for (uint256 i; i < allFeeds.length; ++i) {
      Feed memory feed = allFeeds[i];
      feeds[feed.description] = feed;
      allFeedIds.push(feed.feedId);
    }

    // Associate test feeds with mock registries
    Feed[] memory registryAFeeds = new Feed[](2);
    registryAFeeds[0] = feeds["a"];
    registryAFeeds[1] = feeds["b"];

    Feed[] memory registryBFeeds = new Feed[](1);
    registryBFeeds[0] = feeds["c"];

    // Deploy mock registry contracts with mock data for feeds
    MockRegistry registryA = deployMockRegistry(registryAFeeds);
    MockRegistry registryB = deployMockRegistry(registryBFeeds);

    registries["a"] = address(registryA);
    registries["b"] = address(registryB);

    feeManager = new MockFeeManager();

    router.setFeeManager(address(feeManager));
  }

  // Setup feed registries in Router
  function setRouterFeedRegistries() public {
    MockRegistry registryA = MockRegistry(registries["a"]);
    bytes32[] memory registryA_feedIds = registryA.feedIds();
    router.setRegistry(registryA_feedIds, address(registryA));

    MockRegistry registryB = MockRegistry(registries["b"]);
    bytes32[] memory registryB_feedIds = registryB.feedIds();
    router.setRegistry(registryB_feedIds, address(registryB));
  }

  // Test configureFeeds for setting the Registry and FeeManager for feeds
  function test_configureFeeds() public {
    Feed[] memory feeds = allFeeds;

    bytes32[] memory feedIds = new bytes32[](feeds.length);
    string[] memory descriptions = new string[](feedIds.length);

    for (uint256 i; i < feeds.length; ++i) {
      Feed memory feed = feeds[i];
      feedIds[i] = feed.feedId;
      descriptions[i] = feed.description;
    }

    bytes32 configId = bytes32("1");
    address upkeep = address(0);
    address registryAddress = registries["a"];
    bytes32 feeConfigId = bytes32("1");

    router.configureFeeds(feedIds, descriptions, configId, upkeep, registryAddress, feeConfigId);

    for (uint256 i; i < feeds.length; ++i) {
      address registry = router.getRegistry(feedIds[i]);
      assertEq(registry, registryAddress, "registry should be set to registry A");
      uint256 feedsConfigured = feeManager.feedsConfigured(feedIds[i]);
      assertEq(feedsConfigured, 1, "fee manager should have one configuration set for each feed");
    }
  }

  function test_configureFeeds_revertIfNotOwner() public {
    Feed[] memory feeds = allFeeds;

    bytes32[] memory feedIds = new bytes32[](feeds.length);
    string[] memory descriptions = new string[](feedIds.length);

    for (uint256 i; i < feeds.length; ++i) {
      Feed memory feed = feeds[i];
      feedIds[i] = feed.feedId;
      descriptions[i] = feed.description;
    }

    bytes32 configId = bytes32("1");
    address upkeep = address(0);
    address registryAddress = registries["a"];
    bytes32 feeConfigId = bytes32("1");

    address NOT_OWNER = address(1234);
    vm.stopPrank();
    vm.startPrank(NOT_OWNER);

    vm.expectRevert(bytes("Only callable by owner"));
    router.configureFeeds(feedIds, descriptions, configId, upkeep, registryAddress, feeConfigId);
  }

  function test_setRegistry_revertIfNotOwner() public {
    bytes32[] memory feedIds = new bytes32[](1);
    feedIds[0] = feeds["a"].feedId;

    address registryAddress = registries["a"];

    address NOT_OWNER = address(1234);
    vm.stopPrank();
    vm.startPrank(NOT_OWNER);

    vm.expectRevert(bytes("Only callable by owner"));
    router.setRegistry(feedIds, registryAddress);
  }

  function test_setFeeManager_revertIfNotOwner() public {
    address NOT_OWNER = address(1234);
    vm.stopPrank();
    vm.startPrank(NOT_OWNER);

    vm.expectRevert(bytes("Only callable by owner"));
    router.setFeeManager(address(feeManager));
  }

  // Test setting registry addresses for feed ids, and retrieving them
  function test_assigningRegistries() public {
    setRouterFeedRegistries();

    // Check Feeds A and B have their registry set to Registry A
    address registryA = registries["a"];
    address feedA_Registry = router.getRegistry(feeds["a"].feedId);
    assertEq(feedA_Registry, registryA, "feed A should have registry set to registry A");
    address feedB_Registry = router.getRegistry(feeds["b"].feedId);
    assertEq(feedB_Registry, registryA, "feed B should have registry set to registry A");

    // Check Feed C has their registry set to Registry B
    address registryB = registries["b"];
    address feedC_Registry = router.getRegistry(feeds["c"].feedId);
    assertEq(feedC_Registry, registryB, "feed C should have registry set to registry B");

    // Check that getUniqueAssignedAddresses returns the correct two registries
    address[] memory uniqueRegistries = router.getUniqueRegistries(allFeedIds);
    assertEq(uniqueRegistries.length, 2, "should be two unique registries");
    assertEq(uniqueRegistries[0], registryA, "unique registries should contain registry A");
    assertEq(uniqueRegistries[1], registryB, "unique registries should contain registry B");

    // Check that getAssignedFeedIdsForAddress returns the correct feedIds and indexes for Registry A
    (bytes32[] memory registryA_assignedFeedIds, uint256[] memory registryA_feedsIdsIndexes) = router
      .getRegistryFeedIds(allFeedIds, registryA);
    assertEq(registryA_assignedFeedIds.length, 2, "registryA_assignedFeedIds should contain two assigned feed ids");
    assertEq(registryA_assignedFeedIds[0], feeds["a"].feedId, "registryA_assignedFeedIds should include feed A");
    assertEq(registryA_assignedFeedIds[1], feeds["b"].feedId, "registryA_assignedFeedIds should include feed B");

    assertEq(registryA_feedsIdsIndexes.length, 2, "registryA_feedsIdsIndexes should contain two feed id indexes");
    assertEq(registryA_feedsIdsIndexes[0], 0, "feed A index should be 0");
    assertEq(registryA_feedsIdsIndexes[1], 1, "feed B index should be 1");

    // Check that getAssignedFeedIdsForAddress returns the correct feedIds and indexes for Registry B
    (bytes32[] memory registryB_assignedFeedIds, uint256[] memory registryB_feedsIdsIndexes) = router
      .getRegistryFeedIds(allFeedIds, registryB);
    assertEq(registryB_assignedFeedIds.length, 1, "registryB_feedsIdsIndexes should contain one assigned feed id");
    assertEq(registryB_assignedFeedIds[0], feeds["c"].feedId, "registryB_assignedFeedIds should include feed C");

    assertEq(registryB_feedsIdsIndexes.length, 1, "registryB_feedsIdsIndexes should contain one feed id index");
    assertEq(registryB_feedsIdsIndexes[0], 2, "feed C index should be 2");
  }

  // Test setting fee manager addresses for feed ids, and retrieving them
  function test_assigningFeeManager() public {
    address newFeeManagerAddress = address(10001);

    vm.expectEmit();
    emit FeeManagerSet(newFeeManagerAddress, address(feeManager));

    router.setFeeManager(newFeeManagerAddress);

    address newFeeManager = router.getFeeManager();

    assertEq(newFeeManager, newFeeManagerAddress);
  }

  function test_getBenchmarks() public {
    setRouterFeedRegistries();

    // Check that benchmarks and observationTimestamps are returned, and in the same order as the provided feedIds
    (int256[] memory benchmarks, uint256[] memory observationTimestamps) = router.getBenchmarks{value: 1 ether}(
      allFeedIds,
      billingData
    );
    assertEq(benchmarks.length, 3, "should be three benchmarks, one for each feed");
    assertEq(observationTimestamps.length, 3, "should be three observationTimestamps, one for each feed");

    for (uint256 i; i < benchmarks.length; ++i) {
      Feed memory feed = allFeeds[i];
      assertEq(benchmarks[i], feed.benchmark, "benchmark should match expected value");
      assertEq(observationTimestamps[i], feed.observationTimestamp, "observationTimestamp should match expected value");
    }

    // Check that the fee manager processes a fee for each feed
    for (uint256 i; i < allFeedIds.length; ++i) {
      bytes32 feedId = allFeedIds[i];

      uint256 feesProcessed = feeManager.feesProcessed(feedId);
      assertEq(feesProcessed, 1, "fee manager should have processed one fee for each feed");
    }
  }

  function test_getReports() public {
    setRouterFeedRegistries();

    // Check that reports and observationTimestamps are returned, and in the same order as the provided feedIds
    (bytes[] memory reports, uint256[] memory observationTimestamps) = router.getReports{value: 1 ether}(
      allFeedIds,
      billingData
    );
    assertEq(reports.length, 3, "should be three reports, one for each feed");
    assertEq(observationTimestamps.length, 3, "should be three observationTimestamps, one for each feed");

    for (uint256 i; i < reports.length; ++i) {
      Feed memory feed = allFeeds[i];
      assertEq(reports[i], feed.report, "report should match expected value");
      assertEq(observationTimestamps[i], feed.observationTimestamp, "observationTimestamp should match expected value");
    }

    // Check that the fee manager processes a fee for each feed
    for (uint256 i; i < allFeedIds.length; ++i) {
      bytes32 feedId = allFeedIds[i];

      uint256 feesProcessed = feeManager.feesProcessed(feedId);
      assertEq(feesProcessed, 1, "fee manager should have processed one fee for each feed");
    }
  }

  function test_getDescriptions() public {
    setRouterFeedRegistries();

    // Check that descriptions are returned, and in the same order as the provided feedIds
    string[] memory descriptions = router.getDescriptions(allFeedIds);
    assertEq(descriptions.length, 3, "should be three descriptions, one for each feed");

    for (uint256 i; i < descriptions.length; ++i) {
      Feed memory feed = allFeeds[i];
      assertEq(descriptions[i], feed.description, "description should match expected value");
    }
  }

  function test_requestUpkeep() public {
    setRouterFeedRegistries();

    router.requestUpkeep{value: 1 ether}(allFeedIds, billingData);

    // Check that the Fee Manager processes a fee, and the Registry requests an upkeep for each feed
    for (uint256 i; i < allFeedIds.length; ++i) {
      bytes32 feedId = allFeedIds[i];

      uint256 feesProcessed = feeManager.feesProcessed(feedId);
      assertEq(feesProcessed, 1, "fee manager should have processed one fee for each feed");

      MockRegistry registry = MockRegistry(router.getRegistry(feedId));
      uint256 upkeepsRequested = registry.upkeepsRequested(feedId);
      assertEq(upkeepsRequested, 1, "registry should have requested one upkeep for each feed");
    }
  }

  function test_typeAndVersion() public {
    string memory typeAndVersion = router.typeAndVersion();
    assertEq(typeAndVersion, "DataFeedsRouter 1.0.0", "typeAndVersion should match expected value");
  }

  function deployMockRegistry(Feed[] memory _feeds) public returns (MockRegistry) {
    bytes32[] memory feedId = new bytes32[](_feeds.length);
    string[] memory descriptions = new string[](_feeds.length);
    string[] memory registryNames = new string[](_feeds.length);
    int256[] memory benchmarks = new int256[](_feeds.length);
    bytes[] memory reports = new bytes[](_feeds.length);
    uint256[] memory observationTimestamps = new uint256[](_feeds.length);

    for (uint256 i; i < _feeds.length; ++i) {
      Feed memory feed = _feeds[i];
      feedId[i] = feed.feedId;
      descriptions[i] = feed.description;
      registryNames[i] = feed.registry;
      benchmarks[i] = feed.benchmark;
      reports[i] = feed.report;
      observationTimestamps[i] = feed.observationTimestamp;
    }

    MockRegistry mockRegistry = new MockRegistry(
      feedId,
      descriptions,
      registryNames,
      benchmarks,
      reports,
      observationTimestamps
    );

    return mockRegistry;
  }
}

// Harness for accessing internal properties in DataFeedRouter for testing
contract RouterHarness is DataFeedsRouter {
  constructor() DataFeedsRouter() {}

  function getUniqueRegistries(bytes32[] calldata feedIds) external view returns (address[] memory registries) {
    return _getUniqueAssignedAddresses(feedIds, s_feedIdToRegistry);
  }

  function getRegistryFeedIds(
    bytes32[] calldata feedIds,
    address registry
  ) external view returns (bytes32[] memory assignedFeedIds, uint256[] memory feedsIdsIndexes) {
    return _getAssignedFeedIdsForAddress(feedIds, registry, s_feedIdToRegistry);
  }
}

// Mock implementation of DataFeedsRegistry for returning mock data
contract MockRegistry is IDataFeedsRegistry {
  Feed[] s_feeds;

  mapping(bytes32 => uint256) s_upkeepsRequested;

  constructor(
    bytes32[] memory feedIds,
    string[] memory descriptions,
    string[] memory registries,
    int256[] memory benchmarks,
    bytes[] memory reports,
    uint256[] memory observationTimestamps
  ) {
    for (uint256 i; i < feedIds.length; ++i) {
      Feed memory feed = Feed({
        feedId: feedIds[i],
        description: descriptions[i],
        registry: registries[i],
        benchmark: benchmarks[i],
        report: reports[i],
        observationTimestamp: observationTimestamps[i]
      });

      s_feeds.push(feed);
    }
  }

  function feedIds() public view returns (bytes32[] memory feedIds) {
    feedIds = new bytes32[](s_feeds.length);
    for (uint256 i; i < s_feeds.length; ++i) {
      feedIds[i] = s_feeds[i].feedId;
    }
    return feedIds;
  }

  function getFeeds(bytes32[] calldata feedIds) internal view returns (Feed[] memory feeds) {
    feeds = new Feed[](feedIds.length);

    for (uint256 i; i < feedIds.length; ++i) {
      bool found = false;

      for (uint256 j = 0; j < s_feeds.length; ++j) {
        Feed memory feed = s_feeds[j];
        if (feed.feedId == feedIds[i]) {
          found = true;
          feeds[i] = feed;
          break;
        }
      }

      if (!found) {
        revert("feed not found in MockRegistry mock data");
      }
    }
    return feeds;
  }

  function getBenchmarks(
    bytes32[] calldata feedIds
  ) external view returns (int256[] memory benchmarks, uint256[] memory observationTimestamps) {
    benchmarks = new int256[](feedIds.length);
    observationTimestamps = new uint256[](feedIds.length);

    Feed[] memory feeds = getFeeds(feedIds);
    for (uint256 i; i < feeds.length; ++i) {
      benchmarks[i] = feeds[i].benchmark;
      observationTimestamps[i] = feeds[i].observationTimestamp;
    }

    return (benchmarks, observationTimestamps);
  }

  function getReports(
    bytes32[] calldata feedIds
  ) external view returns (bytes[] memory reports, uint256[] memory observationTimestamps) {
    reports = new bytes[](feedIds.length);
    observationTimestamps = new uint256[](feedIds.length);

    Feed[] memory feeds = getFeeds(feedIds);
    for (uint256 i; i < feeds.length; ++i) {
      reports[i] = feeds[i].report;
      observationTimestamps[i] = feeds[i].observationTimestamp;
    }
    return (reports, observationTimestamps);
  }

  function getFeedMetadata(
    bytes32[] calldata feedIds
  )
    external
    view
    returns (
      string[] memory descriptions,
      bytes32[] memory configIds,
      uint256[] memory deviationThresholds,
      uint256[] memory stalenessSeconds,
      bool[] memory upkeepsRequested
    )
  {
    descriptions = new string[](feedIds.length);
    configIds = new bytes32[](feedIds.length);
    deviationThresholds = new uint256[](feedIds.length);
    stalenessSeconds = new uint256[](feedIds.length);
    upkeepsRequested = new bool[](feedIds.length);

    Feed[] memory feeds = getFeeds(feedIds);
    for (uint256 i; i < feeds.length; ++i) {
      descriptions[i] = feeds[i].description;
    }
    return (descriptions, configIds, deviationThresholds, stalenessSeconds, upkeepsRequested);
  }

  function setFeeds(
    bytes32[] calldata feedIds,
    string[] calldata descriptions,
    bytes32 configId,
    address upkeep
  ) external {}

  function requestUpkeep(bytes32[] calldata feedIds) external {
    for (uint256 i; i < feedIds.length; ++i) {
      s_upkeepsRequested[feedIds[i]] += 1;
    }
  }

  function upkeepsRequested(bytes32 feedId) external view returns (uint256) {
    return s_upkeepsRequested[feedId];
  }

  // ================================================================
  // │                  Un-implemented functions                    │
  // ================================================================
  modifier notImplemented() {
    revert("function not implemented on MockRegistry");
    _;
  }

  function removeFeeds(bytes32[] calldata feedIds) external notImplemented {}

  function setFeedConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata deviationThresholds,
    uint256[] calldata stalenessSeconds
  ) external notImplemented {}

  function updateDescriptions(bytes32[] calldata feedIds, string[] calldata descriptions) external notImplemented {}

  function updateFeedConfigId(bytes32[] calldata feedIds, bytes32 configId) external notImplemented {}

  function updateUpkeep(bytes32[] calldata feedIds, address upkeep) external notImplemented {}

  function performUpkeep(bytes calldata performData) external notImplemented {}

  function getFeedConfigs(
    bytes32[] calldata configIds
  ) external view notImplemented returns (uint256[] memory deviationThresholds, uint256[] memory stalenessSeconds) {}

  function getLinkAddress() external view notImplemented returns (address linkAddress) {}

  function getUpkeepFeedIds(address upkeep) external view notImplemented returns (bytes32[] memory feedIdsData) {}

  function linkAvailableForPayment() external view notImplemented returns (int256 availableBalance) {}
}

// Mock implementation of DataFeedsRegistry for returning mock data
contract MockFeeManager is IDataFeedsFeeManager {
  mapping(bytes32 => uint256) s_feesProcessed;
  mapping(bytes32 => uint256) s_feedsConfigured;

  constructor() {}

  function feesProcessed(bytes32 feedId) external view returns (uint256) {
    return s_feesProcessed[feedId];
  }

  function feedsConfigured(bytes32 feedId) external view returns (uint256) {
    return s_feedsConfigured[feedId];
  }

  function processFee(
    address sender,
    IDataFeedsFeeManager.Service service,
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable {
    // In real implementation the FeeManager can support both native and ERC20 payments,
    // but for unit testing the Router we only support native payments
    if (msg.value <= 0) {
      revert("MockFeeManager.processFee: msg.value must be greater than 0");
    }
    for (uint256 i; i < feedIds.length; ++i) {
      s_feesProcessed[feedIds[i]] += 1;
    }
  }

  function getFee(
    address sender,
    IDataFeedsFeeManager.Service service,
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external view returns (uint256 fee) {
    return feedIds.length * 0.01 ether;
  }

  function setFeedServiceFees(bytes32 configId, bytes32[] calldata feedIds) external {
    for (uint256 i; i < feedIds.length; ++i) {
      s_feedsConfigured[feedIds[i]] += 1;
    }
  }

  // ================================================================
  // │                  Un-implemented functions                    │
  // ================================================================
  modifier notImplemented() {
    revert("function not implemented on MockFeeManager");
    _;
  }

  function enableSpendingAddresses(address[] calldata spendingAddresses) external notImplemented {}

  function disableSpendingAddresses(address[] calldata spendingAddresses) external notImplemented {}

  function addFinanceAdmins(address[] calldata financeAdmins) external notImplemented {}

  function removeFinanceAdmins(address[] calldata financeAdmins) external notImplemented {}

  function addFeeTokens(
    address[] calldata tokenAddresses,
    bytes32[] calldata priceFeedIds,
    bytes32 feeTokenDiscountConfigId
  ) external notImplemented {}

  function removeFeeTokens(address[] calldata tokenAddresses) external notImplemented {}

  function setServiceFeeConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata getBenchmarkUsdFees,
    uint256[] calldata getReportUsdFees,
    uint256[] calldata requestUpkeepUsdFees
  ) external notImplemented {}

  function setFeeTokenDiscountConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata discounts
  ) external notImplemented {}

  function setFeeTokenDiscounts(bytes32 configId, address[] calldata tokenAddresses) external notImplemented {}

  function setSenderDiscountConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata discounts
  ) external notImplemented {}

  function setSenderDiscounts(bytes32 configId, address[] calldata senders) external notImplemented {}

  function withdraw(
    address[] calldata tokenAddresses,
    uint256[] calldata quantities,
    address recipientAddress
  ) external notImplemented {}
}
