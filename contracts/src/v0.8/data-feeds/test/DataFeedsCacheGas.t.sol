// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {DataFeedsSetupGas} from "./DataFeedsSetupGas.t.sol";

contract DataFeedsCacheGasTest is DataFeedsSetupGas {
  address[] internal singleProxyList = new address[](1);
  address[] internal proxyList = new address[](5);
  address[] internal newSingleProxyList = new address[](1);
  address[] internal newProxyList = new address[](5);

  bytes priceReportBytes1 = abi.encodePacked(
    hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
    hex"0000000000000000000000000000000000000000000000000000000000000001", // Length
    hex"010e12d1e0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(prices[0])
  );
  bytes priceReportBytes5 = abi.encodePacked(
    hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
    hex"0000000000000000000000000000000000000000000000000000000000000005", // Length
    hex"010e12d1e0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(prices[0]),
    hex"010e12dde0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(prices[1]),
    hex"01b476d70d000232000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(prices[2]),
    hex"0169bd6041000132000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(prices[3]),
    hex"010e12f1e0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(prices[4])
  );

  function setUp() public virtual override {
    DataFeedsSetupGas.setUp();

    singleProxyList[0] = address(10002);

    proxyList[0] = address(10002);
    proxyList[1] = address(dataFeedsLegacyAggregatorProxy);
    proxyList[2] = address(dataFeedsAggregatorProxy);
    proxyList[3] = address(10005);
    proxyList[4] = address(10006);

    newSingleProxyList[0] = address(10007);

    newProxyList[0] = address(10002);
    newProxyList[1] = address(10003);
    newProxyList[2] = address(10004);
    newProxyList[3] = address(10005);
    newProxyList[4] = address(10006);

    vm.startPrank(OWNER);
    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, batchValueIds);
  }

  function test_write_setDecimalFeedConfigs_1_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_1_gas");
    dataFeedsCache.setDecimalFeedConfigs(dataIds1New, descriptions1, workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_1_gas");
  }

  function test_write_setDecimalFeedConfigs_5_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_5_gas");
    dataFeedsCache.setDecimalFeedConfigs(dataIds5New, descriptions5, workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_5_gas");
  }

  function test_write_setDecimalFeedConfigs_with_delete_1_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_1_gas");
    dataFeedsCache.setDecimalFeedConfigs(dataIds1Old, descriptions1, workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_1_gas");
  }

  function test_write_setDecimalFeedConfigs_with_delete_5_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_5_gas");
    dataFeedsCache.setDecimalFeedConfigs(dataIds5Old, descriptions5, workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_5_gas");
  }

  function test_write_setBundleFeedConfigs_1_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_1_gas");
    dataFeedsCache.setBundleFeedConfigs(dataIds1New, descriptions1, decimals1, workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_1_gas");
  }

  function test_write_setBundleFeedConfigs_5_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_5_gas");
    dataFeedsCache.setBundleFeedConfigs(dataIds5New, descriptions5, decimals5, workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_5_gas");
  }

  function test_write_setBundleFeedConfigs_with_delete_1_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_with_delete_1_gas");
    dataFeedsCache.setBundleFeedConfigs(dataIds1Old, descriptions1, decimals1, workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_with_delete_1_gas");
  }

  function test_write_setBundleFeedConfigs_with_delete_5_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_with_delete_5_gas");
    dataFeedsCache.setBundleFeedConfigs(dataIds5Old, descriptions5, decimals5, workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_with_delete_5_gas");
  }

  function test_write_removeFeedConfigs_1_gas() public {
    vm.startSnapshotGas("test_write_removeFeedConfigs_1_gas");
    dataFeedsCache.removeFeedConfigs(dataIds1Old);
    vm.stopSnapshotGas("test_write_removeFeedConfigs_1_gas");
  }

  function test_write_removeFeedConfigs_5_gas() public {
    vm.startSnapshotGas("test_write_removeFeedConfigs_5_gas");
    dataFeedsCache.removeFeedConfigs(dataIds5Old);
    vm.stopSnapshotGas("test_write_removeFeedConfigs_5_gas");
  }

  function test_write_onReport_prices_1_gas() public {
    vm.startSnapshotGas("test_write_onReport_prices_1_gas");
    vm.startPrank(reportSender);
    dataFeedsCache.onReport(metadata, priceReportBytes1);
    vm.stopSnapshotGas("test_write_onReport_prices_1_gas");
  }

  function test_write_onReport_prices_5_gas() public {
    vm.startSnapshotGas("test_write_onReport_prices_5_gas");
    vm.startPrank(reportSender);
    dataFeedsCache.onReport(metadata, priceReportBytes5);
    vm.stopSnapshotGas("test_write_onReport_prices_5_gas");
  }

  function test_updateDataIdMappingsForProxies1feed_gas() public {
    vm.startSnapshotGas("test_updateDataIdMappingsForProxies1feed_gas");
    dataFeedsCache.updateDataIdMappingsForProxies(newSingleProxyList, singleValueId);
    vm.stopSnapshotGas("test_updateDataIdMappingsForProxies1feed_gas");
  }

  function test_updateDataIdMappingsForProxies5feeds_gas() public {
    vm.startSnapshotGas("test_updateDataIdMappingsForProxies5feeds_gas");
    dataFeedsCache.updateDataIdMappingsForProxies(newProxyList, batchValueIds);
    vm.stopSnapshotGas("test_updateDataIdMappingsForProxies5feeds_gas");
  }

  function test_removeDataIdMappingsForProxies1feed_gas() public {
    vm.startSnapshotGas("test_removeDataIdMappingsForProxies1feed_gas");
    dataFeedsCache.removeDataIdMappingsForProxies(singleProxyList);
    vm.stopSnapshotGas("test_removeDataIdMappingsForProxies1feed_gas");
  }

  function test_removeDataIdMappingsForProxies5feeds_gas() public {
    vm.startSnapshotGas("test_removeDataIdMappingsForProxies5feeds_gas");
    dataFeedsCache.removeDataIdMappingsForProxies(proxyList);
    vm.stopSnapshotGas("test_removeDataIdMappingsForProxies5feeds_gas");
  }

  /// AggregatorInterface

  function test_latestAnswer_proxy_gas() public {
    vm.startSnapshotGas("test_latestAnswer_proxy_gas");
    dataFeedsLegacyAggregatorProxy.latestAnswer();
    vm.stopSnapshotGas("test_latestAnswer_proxy_gas");
  }

  function test_latestTimestamp_proxy_gas() public {
    vm.startSnapshotGas("test_latestTimestamp_proxy_gas");
    dataFeedsLegacyAggregatorProxy.latestTimestamp();
    vm.stopSnapshotGas("test_latestTimestamp_proxy_gas");
  }

  function test_latestRound_proxy_gas() public {
    vm.startSnapshotGas("test_latestRound_proxy_gas");
    dataFeedsLegacyAggregatorProxy.latestRound();
    vm.stopSnapshotGas("test_latestRound_proxy_gas");
  }

  function test_getAnswer_proxy_gas() public {
    vm.startSnapshotGas("test_getAnswer_proxy_gas");
    dataFeedsLegacyAggregatorProxy.getAnswer(18446744073709551617);
    vm.stopSnapshotGas("test_getAnswer_proxy_gas");
  }

  function test_getTimestamp_proxy_gas() public {
    vm.startSnapshotGas("test_getTimestamp_proxy_gas");
    dataFeedsLegacyAggregatorProxy.getTimestamp(18446744073709551617);
    vm.stopSnapshotGas("test_getTimestamp_proxy_gas");
  }

  /// AggregatorV3Interface

  function test_decimals_proxy_gas() public {
    vm.startSnapshotGas("test_decimals_proxy_gas");
    dataFeedsLegacyAggregatorProxy.decimals();
    vm.stopSnapshotGas("test_decimals_proxy_gas");
  }

  function test_description_proxy_gas() public {
    vm.startSnapshotGas("test_description_proxy_gas");
    dataFeedsLegacyAggregatorProxy.description();
    vm.stopSnapshotGas("test_description_proxy_gas");
  }

  function test_getRoundData_proxy_gas() public {
    vm.startSnapshotGas("test_getRoundData_proxy_gas");
    dataFeedsLegacyAggregatorProxy.getRoundData(uint80(18446744073709551617));
    vm.stopSnapshotGas("test_getRoundData_proxy_gas");
  }

  function test_latestRoundData_proxy_gas() public {
    vm.startSnapshotGas("test_latestRoundData_proxy_gas");
    dataFeedsLegacyAggregatorProxy.latestRoundData();
    vm.stopSnapshotGas("test_latestRoundData_proxy_gas");
  }

  /// BundleAggregatorInterface
  function test_bundleDecimals_proxy_gas() public {
    vm.startSnapshotGas("test_bundleDecimals_proxy_gas");
    dataFeedsAggregatorProxy.bundleDecimals();
    vm.stopSnapshotGas("test_bundleDecimals_proxy_gas");
  }

  function test_latestBundle_proxy_gas() public {
    vm.startSnapshotGas("test_latestBundle_proxy_gas");
    dataFeedsAggregatorProxy.latestBundle();
    vm.stopSnapshotGas("test_latestBundle_proxy_gas");
  }

  function test_latestBundleTimestamp_proxy_gas() public {
    vm.startSnapshotGas("test_latestBundleTimestamp_proxy_gas");
    dataFeedsAggregatorProxy.latestBundleTimestamp();
    vm.stopSnapshotGas("test_latestBundleTimestamp_proxy_gas");
  }
}
