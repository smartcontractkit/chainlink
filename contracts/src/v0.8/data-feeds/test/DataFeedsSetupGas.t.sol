// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {BundleAggregatorProxy} from "../BundleAggregatorProxy.sol";
import {DataFeedsCache} from "../DataFeedsCache.sol";

import {BaseTest} from "./BaseTest.t.sol";
import {DataFeedsLegacyAggregatorProxy} from "./helpers/DataFeedsLegacyAggregatorProxy.sol";

// solhint-disable-next-line max-states-count
contract DataFeedsSetupGas is BaseTest {
  struct ReceivedBundleReport {
    bytes32 dataId;
    uint32 timestamp;
    bytes bundle;
  }

  DataFeedsLegacyAggregatorProxy internal s_dataFeedsLegacyAggregatorProxy;
  BundleAggregatorProxy internal s_dataFeedsAggregatorProxy;
  DataFeedsCache internal s_dataFeedsCache;

  string[] internal s_descriptions1 = new string[](1);
  string[] internal s_descriptions5 = new string[](5);

  uint8[][] internal s_decimals1 = new uint8[][](1);
  uint8[][] internal s_decimals5 = new uint8[][](5);

  bytes16[] internal s_dataIds = new bytes16[](5);
  bytes16[] internal s_dataIds1Old = new bytes16[](1);
  bytes16[] internal s_dataIds1New = new bytes16[](1);
  bytes16[] internal s_dataIds5Old = new bytes16[](5);
  bytes16[] internal s_dataIds5New = new bytes16[](5);

  bytes16[] internal s_singleValueId = new bytes16[](1);
  bytes16[] internal s_batchValueIds = new bytes16[](5);

  bytes32[] internal s_paddedDataIds = new bytes32[](5);
  uint256 internal s_price1 = 123456;
  uint256 internal s_price2 = 456789;
  uint32 internal s_timestamp1 = 0;
  uint32 internal s_timestamp2 = 0;
  uint32 internal s_timestamp3 = 0;

  address internal s_reportSender = address(10002);
  string internal s_description = "description";
  bytes32 internal s_workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  bytes2 internal s_reportId = hex"0001";
  address[] internal s_senders = [s_reportSender, s_reportSender];
  address[] internal s_workflowOwners = [address(10004), address(10005)];
  bytes10[] internal s_workflowNames = [bytes10("abc"), bytes10("xyz")];

  DataFeedsCache.WorkflowMetadata internal s_workflowMetadata1 = DataFeedsCache.WorkflowMetadata({
    allowedSender: s_senders[0],
    allowedWorkflowOwner: s_workflowOwners[0],
    allowedWorkflowName: s_workflowNames[0]
  });

  DataFeedsCache.WorkflowMetadata internal s_workflowMetadata2 = DataFeedsCache.WorkflowMetadata({
    allowedSender: s_senders[1],
    allowedWorkflowOwner: s_workflowOwners[1],
    allowedWorkflowName: s_workflowNames[1]
  });
  DataFeedsCache.WorkflowMetadata[] internal s_workflowMetadata;

  uint256[] internal s_prices = [123, 456, 789, 876, 543];
  uint32[] internal s_timestamps = [12, 34, 56, 78, 90];
  bytes internal s_metadata;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_dataFeedsCache = new DataFeedsCache();
    s_dataFeedsLegacyAggregatorProxy = new DataFeedsLegacyAggregatorProxy(address(s_dataFeedsCache));
    s_dataFeedsAggregatorProxy = new BundleAggregatorProxy(address(s_dataFeedsCache), OWNER);

    s_paddedDataIds = new bytes32[](10);
    s_paddedDataIds[0] = hex"010e12d1e0000032000000000000000000000000000000000000000000000000";
    s_paddedDataIds[1] = hex"010e12dde0000032000000000000000000000000000000000000000000000000";
    s_paddedDataIds[2] = hex"01b476d70d000232000000000000000000000000000000000000000000000000";
    s_paddedDataIds[3] = hex"0169bd6041000132000000000000000000000000000000000000000000000000";
    s_paddedDataIds[4] = hex"010e12f1e0000032000000000000000000000000000000000000000000000000";
    s_paddedDataIds[5] = hex"010e1ab1e0000004000000000000000000000000000000000000000000000000";
    s_paddedDataIds[6] = hex"0112345670000004000000000000000000000000000000000000000000000000";
    s_paddedDataIds[7] = hex"0198765432000004000000000000000000000000000000000000000000000000";
    s_paddedDataIds[8] = hex"0187654321000004000000000000000000000000000000000000000000000000";
    s_paddedDataIds[9] = hex"0112754834000004000000000000000000000000000000000000000000000000";

    s_descriptions1 = new string[](1);
    s_descriptions1[0] = "description0";

    s_descriptions5 = new string[](5);
    s_descriptions5[0] = "description0";
    s_descriptions5[1] = "description1";
    s_descriptions5[2] = "description2";
    s_descriptions5[3] = "description3";
    s_descriptions5[4] = "description4";

    s_decimals1 = new uint8[][](1);
    s_decimals1[0] = new uint8[](1);
    s_decimals1[0][0] = 18;

    s_decimals5 = new uint8[][](5);
    s_decimals5[0] = new uint8[](1);
    s_decimals5[0][0] = 18;
    s_decimals5[1] = new uint8[](2);
    s_decimals5[1][0] = 18;
    s_decimals5[1][1] = 0;
    s_decimals5[2] = new uint8[](1);
    s_decimals5[2][0] = 18;
    s_decimals5[3] = new uint8[](3);
    s_decimals5[3][0] = 18;
    s_decimals5[3][1] = 8;
    s_decimals5[3][2] = 1;
    s_decimals5[4] = new uint8[](1);
    s_decimals5[4][0] = 18;

    s_dataIds = new bytes16[](10);
    s_dataIds[0] = bytes16(s_paddedDataIds[0]);
    s_dataIds[1] = bytes16(s_paddedDataIds[1]);
    s_dataIds[2] = bytes16(s_paddedDataIds[2]);
    s_dataIds[3] = bytes16(s_paddedDataIds[3]);
    s_dataIds[4] = bytes16(s_paddedDataIds[4]);
    s_dataIds[5] = bytes16(s_paddedDataIds[5]);
    s_dataIds[6] = bytes16(s_paddedDataIds[6]);
    s_dataIds[7] = bytes16(s_paddedDataIds[7]);
    s_dataIds[8] = bytes16(s_paddedDataIds[8]);
    s_dataIds[9] = bytes16(s_paddedDataIds[9]);

    s_dataIds1Old[0] = s_dataIds[0];

    s_dataIds1New[0] = s_dataIds[5];

    s_dataIds5Old[0] = s_dataIds[0];
    s_dataIds5Old[1] = s_dataIds[1];
    s_dataIds5Old[2] = s_dataIds[2];
    s_dataIds5Old[3] = s_dataIds[3];
    s_dataIds5Old[4] = s_dataIds[4];

    s_dataIds5New[0] = s_dataIds[5];
    s_dataIds5New[1] = s_dataIds[6];
    s_dataIds5New[2] = s_dataIds[7];
    s_dataIds5New[3] = s_dataIds[8];
    s_dataIds5New[4] = s_dataIds[9];

    s_singleValueId = new bytes16[](1);
    s_singleValueId[0] = s_dataIds[0];

    s_batchValueIds = new bytes16[](5);
    s_batchValueIds[0] = s_dataIds[0];
    s_batchValueIds[1] = s_dataIds[1];
    s_batchValueIds[2] = s_dataIds[2];
    s_batchValueIds[3] = s_dataIds[3];
    s_batchValueIds[4] = s_dataIds[4];

    s_metadata = abi.encodePacked(s_workflowId, s_workflowNames[0], s_workflowOwners[0], s_reportId);

    s_workflowMetadata.push(s_workflowMetadata1);
    s_workflowMetadata.push(s_workflowMetadata2);

    s_dataFeedsCache.setFeedAdmin(OWNER, true);

    s_dataFeedsCache.setDecimalFeedConfigs(s_dataIds5Old, s_descriptions5, s_workflowMetadata);
    s_dataFeedsCache.setBundleFeedConfigs(s_dataIds5Old, s_descriptions5, s_decimals5, s_workflowMetadata);

    vm.stopPrank();
    vm.startPrank(s_reportSender);

    s_dataFeedsCache.onReport(
      s_metadata,
      abi.encodePacked(
        hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
        hex"0000000000000000000000000000000000000000000000000000000000000003", // Length
        s_paddedDataIds[0],
        abi.encode(s_timestamps[0]),
        abi.encode(s_prices[0]),
        s_paddedDataIds[1],
        abi.encode(s_timestamps[1]),
        abi.encode(s_prices[1]),
        s_paddedDataIds[2],
        abi.encode(s_timestamps[2]),
        abi.encode(s_prices[2])
      )
    );

    s_dataFeedsCache.onReport(
      s_metadata,
      abi.encodePacked(
        hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
        hex"0000000000000000000000000000000000000000000000000000000000000002", // length
        hex"0000000000000000000000000000000000000000000000000000000000000040", // offset of ReportOne
        hex"0000000000000000000000000000000000000000000000000000000000000100", // offset of ReportTwo
        s_paddedDataIds[0], // ReportOne FeedID
        abi.encode(s_timestamps[0]),
        hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
        hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
        abi.encode(s_prices[0]),
        abi.encode(s_prices[1]),
        s_paddedDataIds[1], // ReportTwo FeedID
        abi.encode(s_timestamps[1]),
        hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportTwo Bundle
        hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportTwo Bundle
        abi.encode(s_prices[2]),
        abi.encode(s_prices[3])
      )
    );

    vm.stopPrank();
  }
}
