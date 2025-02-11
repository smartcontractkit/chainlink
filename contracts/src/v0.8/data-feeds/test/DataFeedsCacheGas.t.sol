// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {DataFeedsSetupGas} from "./DataFeedsSetupGas.t.sol";

contract DataFeedsCacheGasTest is DataFeedsSetupGas {
  address[] internal s_singleProxyList = new address[](1);
  address[] internal s_proxyList = new address[](5);
  address[] internal s_newSingleProxyList = new address[](1);
  address[] internal s_newProxyList = new address[](5);

  bytes internal s_priceReportBytes1 = abi.encodePacked(
    hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
    hex"0000000000000000000000000000000000000000000000000000000000000001", // Length
    hex"010e12d1e0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(s_prices[0])
  );
  bytes internal s_priceReportBytes5 = abi.encodePacked(
    hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
    hex"0000000000000000000000000000000000000000000000000000000000000005", // Length
    hex"010e12d1e0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(s_prices[0]),
    hex"010e12dde0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(s_prices[1]),
    hex"01b476d70d000232000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(s_prices[2]),
    hex"0169bd6041000132000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(s_prices[3]),
    hex"010e12f1e0000032000000000000000000000000000000000000000000000000",
    abi.encode(100), // Timestamp
    abi.encode(s_prices[4])
  );

  function setUp() public virtual override {
    DataFeedsSetupGas.setUp();

    s_singleProxyList[0] = address(10002);

    s_proxyList[0] = address(10002);
    s_proxyList[1] = address(s_dataFeedsLegacyAggregatorProxy);
    s_proxyList[2] = address(s_dataFeedsAggregatorProxy);
    s_proxyList[3] = address(10005);
    s_proxyList[4] = address(10006);

    s_newSingleProxyList[0] = address(10007);

    s_newProxyList[0] = address(10002);
    s_newProxyList[1] = address(10003);
    s_newProxyList[2] = address(10004);
    s_newProxyList[3] = address(10005);
    s_newProxyList[4] = address(10006);

    vm.startPrank(OWNER);
    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, s_batchValueIds);
  }

  function test_write_setDecimalFeedConfigs_1_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_1_gas");
    s_dataFeedsCache.setDecimalFeedConfigs(s_dataIds1New, s_descriptions1, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_1_gas");
  }

  function test_write_setDecimalFeedConfigs_5_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_5_gas");
    s_dataFeedsCache.setDecimalFeedConfigs(s_dataIds5New, s_descriptions5, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_5_gas");
  }

  function test_write_setDecimalFeedConfigs_with_delete_1_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_1_gas");
    s_dataFeedsCache.setDecimalFeedConfigs(s_dataIds1Old, s_descriptions1, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_1_gas");
  }

  function test_write_setDecimalFeedConfigs_with_delete_5_gas() public {
    vm.startSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_5_gas");
    s_dataFeedsCache.setDecimalFeedConfigs(s_dataIds5Old, s_descriptions5, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setDecimalFeedConfigs_with_delete_5_gas");
  }

  function test_write_setBundleFeedConfigs_1_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_1_gas");
    s_dataFeedsCache.setBundleFeedConfigs(s_dataIds1New, s_descriptions1, s_decimals1, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_1_gas");
  }

  function test_write_setBundleFeedConfigs_5_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_5_gas");
    s_dataFeedsCache.setBundleFeedConfigs(s_dataIds5New, s_descriptions5, s_decimals5, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_5_gas");
  }

  function test_write_setBundleFeedConfigs_with_delete_1_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_with_delete_1_gas");
    s_dataFeedsCache.setBundleFeedConfigs(s_dataIds1Old, s_descriptions1, s_decimals1, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_with_delete_1_gas");
  }

  function test_write_setBundleFeedConfigs_with_delete_5_gas() public {
    vm.startSnapshotGas("test_write_setBundleFeedConfigs_with_delete_5_gas");
    s_dataFeedsCache.setBundleFeedConfigs(s_dataIds5Old, s_descriptions5, s_decimals5, s_workflowMetadata);
    vm.stopSnapshotGas("test_write_setBundleFeedConfigs_with_delete_5_gas");
  }

  function test_write_removeFeedConfigs_1_gas() public {
    vm.startSnapshotGas("test_write_removeFeedConfigs_1_gas");
    s_dataFeedsCache.removeFeedConfigs(s_dataIds1Old);
    vm.stopSnapshotGas("test_write_removeFeedConfigs_1_gas");
  }

  function test_write_removeFeedConfigs_5_gas() public {
    vm.startSnapshotGas("test_write_removeFeedConfigs_5_gas");
    s_dataFeedsCache.removeFeedConfigs(s_dataIds5Old);
    vm.stopSnapshotGas("test_write_removeFeedConfigs_5_gas");
  }

  function test_write_onReport_prices_1_gas() public {
    vm.startSnapshotGas("test_write_onReport_prices_1_gas");
    vm.startPrank(s_reportSender);
    s_dataFeedsCache.onReport(s_metadata, s_priceReportBytes1);
    vm.stopSnapshotGas("test_write_onReport_prices_1_gas");
  }

  function test_write_onReport_prices_5_gas() public {
    vm.startSnapshotGas("test_write_onReport_prices_5_gas");
    vm.startPrank(s_reportSender);
    s_dataFeedsCache.onReport(s_metadata, s_priceReportBytes5);
    vm.stopSnapshotGas("test_write_onReport_prices_5_gas");
  }

  function test_updateDataIdMappingsForProxies1feed_gas() public {
    vm.startSnapshotGas("test_updateDataIdMappingsForProxies1feed_gas");
    s_dataFeedsCache.updateDataIdMappingsForProxies(s_newSingleProxyList, s_singleValueId);
    vm.stopSnapshotGas("test_updateDataIdMappingsForProxies1feed_gas");
  }

  function test_updateDataIdMappingsForProxies5feeds_gas() public {
    vm.startSnapshotGas("test_updateDataIdMappingsForProxies5feeds_gas");
    s_dataFeedsCache.updateDataIdMappingsForProxies(s_newProxyList, s_batchValueIds);
    vm.stopSnapshotGas("test_updateDataIdMappingsForProxies5feeds_gas");
  }

  function test_removeDataIdMappingsForProxies1feed_gas() public {
    vm.startSnapshotGas("test_removeDataIdMappingsForProxies1feed_gas");
    s_dataFeedsCache.removeDataIdMappingsForProxies(s_singleProxyList);
    vm.stopSnapshotGas("test_removeDataIdMappingsForProxies1feed_gas");
  }

  function test_removeDataIdMappingsForProxies5feeds_gas() public {
    vm.startSnapshotGas("test_removeDataIdMappingsForProxies5feeds_gas");
    s_dataFeedsCache.removeDataIdMappingsForProxies(s_proxyList);
    vm.stopSnapshotGas("test_removeDataIdMappingsForProxies5feeds_gas");
  }

  /// AggregatorInterface

  function test_latestAnswer_proxy_gas() public {
    vm.startSnapshotGas("test_latestAnswer_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.latestAnswer();
    vm.stopSnapshotGas("test_latestAnswer_proxy_gas");
  }

  function test_latestTimestamp_proxy_gas() public {
    vm.startSnapshotGas("test_latestTimestamp_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.latestTimestamp();
    vm.stopSnapshotGas("test_latestTimestamp_proxy_gas");
  }

  function test_latestRound_proxy_gas() public {
    vm.startSnapshotGas("test_latestRound_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.latestRound();
    vm.stopSnapshotGas("test_latestRound_proxy_gas");
  }

  function test_getAnswer_proxy_gas() public {
    vm.startSnapshotGas("test_getAnswer_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.getAnswer(18446744073709551617);
    vm.stopSnapshotGas("test_getAnswer_proxy_gas");
  }

  function test_getTimestamp_proxy_gas() public {
    vm.startSnapshotGas("test_getTimestamp_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.getTimestamp(18446744073709551617);
    vm.stopSnapshotGas("test_getTimestamp_proxy_gas");
  }

  /// AggregatorV3Interface

  function test_decimals_proxy_gas() public {
    vm.startSnapshotGas("test_decimals_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.decimals();
    vm.stopSnapshotGas("test_decimals_proxy_gas");
  }

  function test_description_proxy_gas() public {
    vm.startSnapshotGas("test_description_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.description();
    vm.stopSnapshotGas("test_description_proxy_gas");
  }

  function test_getRoundData_proxy_gas() public {
    vm.startSnapshotGas("test_getRoundData_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.getRoundData(uint80(18446744073709551617));
    vm.stopSnapshotGas("test_getRoundData_proxy_gas");
  }

  function test_latestRoundData_proxy_gas() public {
    vm.startSnapshotGas("test_latestRoundData_proxy_gas");
    s_dataFeedsLegacyAggregatorProxy.latestRoundData();
    vm.stopSnapshotGas("test_latestRoundData_proxy_gas");
  }

  /// BundleAggregatorInterface
  function test_bundleDecimals_proxy_gas() public {
    vm.startSnapshotGas("test_bundleDecimals_proxy_gas");
    s_dataFeedsAggregatorProxy.bundleDecimals();
    vm.stopSnapshotGas("test_bundleDecimals_proxy_gas");
  }

  function test_latestBundle_proxy_gas() public {
    vm.startSnapshotGas("test_latestBundle_proxy_gas");
    s_dataFeedsAggregatorProxy.latestBundle();
    vm.stopSnapshotGas("test_latestBundle_proxy_gas");
  }

  function test_latestBundleTimestamp_proxy_gas() public {
    vm.startSnapshotGas("test_latestBundleTimestamp_proxy_gas");
    s_dataFeedsAggregatorProxy.latestBundleTimestamp();
    vm.stopSnapshotGas("test_latestBundleTimestamp_proxy_gas");
  }
}
