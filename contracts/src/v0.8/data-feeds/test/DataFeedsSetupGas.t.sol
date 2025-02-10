// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ERC20Mock} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";

import {BundleAggregatorProxy} from "../BundleAggregatorProxy.sol";
import {DataFeedsCache} from "../DataFeedsCache.sol";

import {BaseTest} from "./BaseTest.t.sol";
import {DataFeedsLegacyAggregatorProxy} from "./helpers/DataFeedsLegacyAggregatorProxy.sol";

contract DataFeedsSetupGas is BaseTest {
  struct ReceivedBundleReport {
    bytes32 dataId;
    uint32 timestamp;
    bytes bundle;
  }

  DataFeedsLegacyAggregatorProxy internal dataFeedsLegacyAggregatorProxy;
  BundleAggregatorProxy internal dataFeedsAggregatorProxy;
  DataFeedsCache internal dataFeedsCache;

  string[] internal descriptions1 = new string[](1);
  string[] internal descriptions5 = new string[](5);

  uint8[][] internal decimals1 = new uint8[][](1);
  uint8[][] internal decimals5 = new uint8[][](5);

  bytes16[] internal dataIds = new bytes16[](5);
  bytes16[] internal dataIds1Old = new bytes16[](1);
  bytes16[] internal dataIds1New = new bytes16[](1);
  bytes16[] internal dataIds5Old = new bytes16[](5);
  bytes16[] internal dataIds5New = new bytes16[](5);

  bytes16[] internal singleValueId = new bytes16[](1);
  bytes16[] internal batchValueIds = new bytes16[](5);

  bytes32[] internal paddedDataIds = new bytes32[](5);
  uint256 internal price1 = 123456;
  uint256 internal price2 = 456789;
  uint32 internal timestamp1 = 0;
  uint32 internal timestamp2 = 0;
  uint32 internal timestamp3 = 0;

  address internal reportSender = address(10002);
  string internal description = "description";
  bytes32 internal workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  bytes2 internal reportId = hex"0001";
  address[] internal senders = [reportSender, reportSender];
  address[] internal workflowOwners = [address(10004), address(10005)];
  bytes10[] internal workflowNames = [bytes10("abc"), bytes10("xyz")];

  DataFeedsCache.WorkflowMetadata internal workflowMetadata1 = DataFeedsCache.WorkflowMetadata({
    allowedSender: senders[0],
    allowedWorkflowOwner: workflowOwners[0],
    allowedWorkflowName: workflowNames[0]
  });

  DataFeedsCache.WorkflowMetadata internal workflowMetadata2 = DataFeedsCache.WorkflowMetadata({
    allowedSender: senders[1],
    allowedWorkflowOwner: workflowOwners[1],
    allowedWorkflowName: workflowNames[1]
  });
  DataFeedsCache.WorkflowMetadata[] internal workflowMetadata;

  uint256[] internal prices = [123, 456, 789, 876, 543];
  uint32[] internal timestamps = [12, 34, 56, 78, 90];
  bytes internal metadata;

  function setUp() public virtual override {
    BaseTest.setUp();

    dataFeedsCache = new DataFeedsCache();
    dataFeedsLegacyAggregatorProxy = new DataFeedsLegacyAggregatorProxy(address(dataFeedsCache));
    dataFeedsAggregatorProxy = new BundleAggregatorProxy(address(dataFeedsCache), OWNER);

    paddedDataIds = new bytes32[](10);
    paddedDataIds[0] = hex"010e12d1e0000032000000000000000000000000000000000000000000000000";
    paddedDataIds[1] = hex"010e12dde0000032000000000000000000000000000000000000000000000000";
    paddedDataIds[2] = hex"01b476d70d000232000000000000000000000000000000000000000000000000";
    paddedDataIds[3] = hex"0169bd6041000132000000000000000000000000000000000000000000000000";
    paddedDataIds[4] = hex"010e12f1e0000032000000000000000000000000000000000000000000000000";
    paddedDataIds[5] = hex"010e1ab1e0000004000000000000000000000000000000000000000000000000";
    paddedDataIds[6] = hex"0112345670000004000000000000000000000000000000000000000000000000";
    paddedDataIds[7] = hex"0198765432000004000000000000000000000000000000000000000000000000";
    paddedDataIds[8] = hex"0187654321000004000000000000000000000000000000000000000000000000";
    paddedDataIds[9] = hex"0112754834000004000000000000000000000000000000000000000000000000";

    descriptions1 = new string[](1);
    descriptions1[0] = "description0";

    descriptions5 = new string[](5);
    descriptions5[0] = "description0";
    descriptions5[1] = "description1";
    descriptions5[2] = "description2";
    descriptions5[3] = "description3";
    descriptions5[4] = "description4";

    decimals1 = new uint8[][](1);
    decimals1[0] = new uint8[](1);
    decimals1[0][0] = 18;

    decimals5 = new uint8[][](5);
    decimals5[0] = new uint8[](1);
    decimals5[0][0] = 18;
    decimals5[1] = new uint8[](2);
    decimals5[1][0] = 18;
    decimals5[1][1] = 0;
    decimals5[2] = new uint8[](1);
    decimals5[2][0] = 18;
    decimals5[3] = new uint8[](3);
    decimals5[3][0] = 18;
    decimals5[3][1] = 8;
    decimals5[3][2] = 1;
    decimals5[4] = new uint8[](1);
    decimals5[4][0] = 18;

    dataIds = new bytes16[](10);
    dataIds[0] = bytes16(paddedDataIds[0]);
    dataIds[1] = bytes16(paddedDataIds[1]);
    dataIds[2] = bytes16(paddedDataIds[2]);
    dataIds[3] = bytes16(paddedDataIds[3]);
    dataIds[4] = bytes16(paddedDataIds[4]);
    dataIds[5] = bytes16(paddedDataIds[5]);
    dataIds[6] = bytes16(paddedDataIds[6]);
    dataIds[7] = bytes16(paddedDataIds[7]);
    dataIds[8] = bytes16(paddedDataIds[8]);
    dataIds[9] = bytes16(paddedDataIds[9]);

    dataIds1Old[0] = dataIds[0];

    dataIds1New[0] = dataIds[5];

    dataIds5Old[0] = dataIds[0];
    dataIds5Old[1] = dataIds[1];
    dataIds5Old[2] = dataIds[2];
    dataIds5Old[3] = dataIds[3];
    dataIds5Old[4] = dataIds[4];

    dataIds5New[0] = dataIds[5];
    dataIds5New[1] = dataIds[6];
    dataIds5New[2] = dataIds[7];
    dataIds5New[3] = dataIds[8];
    dataIds5New[4] = dataIds[9];

    singleValueId = new bytes16[](1);
    singleValueId[0] = dataIds[0];

    batchValueIds = new bytes16[](5);
    batchValueIds[0] = dataIds[0];
    batchValueIds[1] = dataIds[1];
    batchValueIds[2] = dataIds[2];
    batchValueIds[3] = dataIds[3];
    batchValueIds[4] = dataIds[4];

    metadata = abi.encodePacked(workflowId, workflowNames[0], workflowOwners[0], reportId);

    workflowMetadata.push(workflowMetadata1);
    workflowMetadata.push(workflowMetadata2);

    dataFeedsCache.setFeedAdmin(OWNER, true);

    dataFeedsCache.setDecimalFeedConfigs(dataIds5Old, descriptions5, workflowMetadata);
    dataFeedsCache.setBundleFeedConfigs(dataIds5Old, descriptions5, decimals5, workflowMetadata);

    vm.stopPrank();
    vm.startPrank(reportSender);

    dataFeedsCache.onReport(
      metadata,
      abi.encodePacked(
        hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
        hex"0000000000000000000000000000000000000000000000000000000000000003", // Length
        paddedDataIds[0],
        abi.encode(timestamps[0]),
        abi.encode(prices[0]),
        paddedDataIds[1],
        abi.encode(timestamps[1]),
        abi.encode(prices[1]),
        paddedDataIds[2],
        abi.encode(timestamps[2]),
        abi.encode(prices[2])
      )
    );

    dataFeedsCache.onReport(
      metadata,
      abi.encodePacked(
        hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
        hex"0000000000000000000000000000000000000000000000000000000000000002", // length
        hex"0000000000000000000000000000000000000000000000000000000000000040", // offset of ReportOne
        hex"0000000000000000000000000000000000000000000000000000000000000100", // offset of ReportTwo
        paddedDataIds[0], // ReportOne FeedID
        abi.encode(timestamps[0]),
        hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
        hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
        abi.encode(prices[0]),
        abi.encode(prices[1]),
        paddedDataIds[1], // ReportTwo FeedID
        abi.encode(timestamps[1]),
        hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportTwo Bundle
        hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportTwo Bundle
        abi.encode(prices[2]),
        abi.encode(prices[3])
      )
    );

    vm.stopPrank();
  }
}
