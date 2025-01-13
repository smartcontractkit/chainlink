// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC20Mock} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {IERC20Metadata as IERC20} from
  "../../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {BundleAggregatorProxy} from "../BundleAggregatorProxy.sol";

import {DataFeedsCache} from "../DataFeedsCache.sol";
import {IDataFeedsCache} from "../interfaces/IDataFeedsCache.sol";
import {IDecimalAggregator} from "../interfaces/IDecimalAggregator.sol";
import {BaseTest} from "./BaseTest.t.sol";

contract DataFeedsCacheTest is BaseTest {
  BundleAggregatorProxy internal dataFeedsAggregatorProxy;
  DataFeedsCacheHarness internal dataFeedsCache;

  address internal constant ILLEGAL_CALLER = address(11111); // address used as incorrect caller in tests
  address internal constant REPORT_SENDER = address(12222); // mocks keystone forwarder address

  ERC20Mock internal s_link = new ERC20Mock("LINK", "LINK", OWNER, 0);

  bytes32 internal constant workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  bytes10 internal constant workflowName = bytes10("abc");
  address internal constant workflowOwner = address(10004);
  bytes2 internal constant reportId = hex"0001";
  string[] internal descriptions = ["description"];

  uint8[][] internal decimals1By1 = new uint8[][](1);
  uint8[][] internal decimals2By1 = new uint8[][](2);
  uint8[][] internal decimals2By2 = new uint8[][](2);

  bytes internal constant METADATA = abi.encodePacked(workflowId, workflowName, workflowOwner, reportId);

  address[] internal allowedSendersList = [REPORT_SENDER, REPORT_SENDER];
  address[] internal allowedWorkflowOwnersList = [address(10004), address(10005)];
  bytes10[] internal allowedWorkflowNamesList = [bytes10("abc"), bytes10("xyz")];

  address[] internal singleProxyList = new address[](1);
  address[] internal proxyList = new address[](5);
  address[] internal newSingleProxyList = new address[](1);
  address[] internal newProxyList = new address[](5);

  bytes16[] internal singleValueId = new bytes16[](1);
  bytes16[] internal batchValueIds = new bytes16[](5);

  DataFeedsCache.WorkflowMetadata internal workflowMetadata1 = DataFeedsCache.WorkflowMetadata({
    allowedSender: allowedSendersList[0],
    allowedWorkflowOwner: allowedWorkflowOwnersList[0],
    allowedWorkflowName: allowedWorkflowNamesList[0]
  });

  DataFeedsCache.WorkflowMetadata internal workflowMetadata2 = DataFeedsCache.WorkflowMetadata({
    allowedSender: allowedSendersList[1],
    allowedWorkflowOwner: allowedWorkflowOwnersList[1],
    allowedWorkflowName: allowedWorkflowNamesList[1]
  });

  DataFeedsCache.WorkflowMetadata[] internal workflowMetadata;

  bytes internal emptyDecimalReport;
  bytes internal decimalReportlength1;
  bytes internal decimalReportlength2;
  bytes internal emptyBundleReport;
  bytes internal bundleReportlength1;
  bytes internal bundleReportlength2;
  bytes internal staleReport;
  bytes internal staleBundleReport;
  bytes32 internal constant dataId1 = hex"010e12d1e0000032000000000000000000000000000000000000000000000000";
  bytes32 internal constant dataId2 = hex"01b476d70d000232000000000000000000000000000000000000000000000000";
  bytes32 internal constant dataId3 = hex"0169bd6041000103000000000000000000000000000000000000000000000000";
  bytes32 internal constant dataId4 = hex"010e12d1e0000028000000000000000000000000000000000000000000000000";
  bytes32 internal constant dataId5 = hex"010e12d1e0000032000000000000000000000000000000000000000000000000";
  bytes16 internal constant DATA_ID_0 = bytes16(keccak256("12345"));
  bytes16 internal constant DATA_ID_1 = bytes16(keccak256("23456"));
  bytes16 internal constant DATA_ID_2 = bytes16(keccak256("34567"));
  bytes16 internal constant DATA_ID_3 = bytes16(keccak256("45678"));
  bytes16 internal constant DATA_ID_4 = bytes16(keccak256("56789"));
  bytes16 internal constant DATA_ID_5 = bytes16(keccak256("67890"));
  uint256 internal constant price1 = 123456;
  uint256 internal constant price2 = 456789;
  uint256 internal constant price3 = 789456;
  uint256 internal constant price4 = 890123;
  uint256 internal constant price5 = 654321;
  uint256 internal constant price6 = 987654;
  uint32 internal constant timestamp1 = 100;
  uint32 internal constant timestamp2 = 200;

  function setUp() public override {
    super.setUp();
    dataFeedsCache = new DataFeedsCacheHarness();
    dataFeedsCache.setFeedAdmin(OWNER, true);
    dataFeedsAggregatorProxy = new BundleAggregatorProxy(address(dataFeedsCache), OWNER);

    // reports should be encoded as calldata, which has offset and length
    emptyDecimalReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000000" // Length
    );

    // reports should be encoded as calldata, which has offset and length
    decimalReportlength1 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000001", // Length
      dataId1,
      abi.encode(timestamp1),
      abi.encode(price1)
    );

    // reports should be encoded as calldata, which has offset and length
    decimalReportlength2 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // Length
      dataId1,
      abi.encode(timestamp1),
      abi.encode(price3),
      dataId2,
      abi.encode(timestamp2),
      abi.encode(price4)
    );

    staleReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // Length
      dataId1,
      abi.encode(timestamp1 - 50), // report 1 for dataId1 is stale in this report
      abi.encode(price1),
      dataId2,
      abi.encode(timestamp2 + 50),
      abi.encode(price2)
    );

    emptyBundleReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000000", // length
      hex"0000000000000000000000000000000000000000000000000000000000000000" // offset of ReportOne
    );

    bundleReportlength1 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000001", // length
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset of ReportOne
      dataId1, // ReportOne FeedID
      abi.encode(timestamp1),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
      abi.encode(price1),
      abi.encode(price2)
    );

    bundleReportlength2 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // length
      hex"0000000000000000000000000000000000000000000000000000000000000040", // offset of ReportOne
      hex"0000000000000000000000000000000000000000000000000000000000000100", // offset of ReportTwo
      dataId1, // ReportOne FeedID
      abi.encode(timestamp1),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
      abi.encode(price3),
      abi.encode(price4),
      dataId2, // ReportTwo FeedID
      abi.encode(timestamp2),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportTwo Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportTwo Bundle
      abi.encode(price5),
      abi.encode(price6)
    );

    staleBundleReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // length
      hex"0000000000000000000000000000000000000000000000000000000000000040", // offset of ReportOne
      hex"0000000000000000000000000000000000000000000000000000000000000100", // offset of ReportTwo
      dataId1, // ReportOne FeedID
      abi.encode(timestamp1 - 50), // report is stale
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
      abi.encode(price1),
      abi.encode(price2),
      dataId2, // ReportTwo FeedID
      abi.encode(timestamp2 + 50),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportTwo Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportTwo Bundle
      abi.encode(price3),
      abi.encode(price4)
    );

    workflowMetadata.push(workflowMetadata1);
    workflowMetadata.push(workflowMetadata2);

    singleProxyList[0] = address(10002);

    proxyList[0] = address(dataFeedsAggregatorProxy);
    proxyList[1] = address(10002);
    proxyList[2] = address(10004);
    proxyList[3] = address(10005);
    proxyList[4] = address(10006);

    newSingleProxyList[0] = address(10007);

    newProxyList[0] = address(10002);
    newProxyList[1] = address(10003);
    newProxyList[2] = address(10004);
    newProxyList[3] = address(10005);
    newProxyList[4] = address(10006);

    singleValueId = new bytes16[](1);
    singleValueId[0] = bytes16(dataId1);

    batchValueIds = new bytes16[](5);
    batchValueIds[0] = bytes16(dataId1);
    batchValueIds[1] = bytes16(dataId2);
    batchValueIds[2] = bytes16(dataId3);
    batchValueIds[3] = bytes16(dataId4);
    batchValueIds[4] = bytes16(dataId5);

    decimals1By1[0] = new uint8[](1);
    decimals1By1[0][0] = 18;

    decimals2By1[0] = new uint8[](1);
    decimals2By1[0][0] = 18;
    decimals2By1[1] = new uint8[](1);
    decimals2By1[1][0] = 8;

    decimals2By2[0] = new uint8[](2);
    decimals2By2[0][0] = 6;
    decimals2By2[0][1] = 12;
    decimals2By2[1] = new uint8[](2);
    decimals2By2[1][0] = 18;
    decimals2By2[1][0] = 8;

    vm.startPrank(OWNER);
    dataFeedsCache.setFeedAdmin(OWNER, true);

    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, batchValueIds);
  }

  function test_updateDataIdMappingsForProxiesRevertInvalidLengths() public {
    address[] memory proxyList = new address[](1);
    proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](2);
    dataIdList[0] = bytes16(keccak256("12345"));
    dataIdList[1] = bytes16(keccak256("67890"));

    vm.expectRevert(DataFeedsCache.ArrayLengthMismatch.selector);

    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, dataIdList);
  }

  function test_updateDataIdMappingsForProxiesRevertUnauthorizedOwner() public {
    address[] memory proxyList = new address[](1);
    proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.stopPrank();
    vm.startPrank(ILLEGAL_CALLER);
    vm.expectRevert(
      abi.encodeWithSelector(
        DataFeedsCache.UnauthorizedCaller.selector, address(0x0000000000000000000000000000000000002B67)
      )
    );
    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, dataIdList);
  }

  function test_updateDataIdMappingsForProxiesSuccess() public {
    address[] memory proxyList = new address[](1);
    proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(proxyList[0], dataIdList[0]);

    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, dataIdList);
  }

  function test_updateDataIdMappingsForProxies_and_call_decimals() public {
    uint8 decimals = 8;

    vm.startPrank(proxyList[3]);
    uint8 decimalsAns = dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);

    decimals = 18;

    vm.startPrank(proxyList[4]);
    decimalsAns = dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);

    address[] memory newProxyList = new address[](2);
    newProxyList[0] = proxyList[3];
    newProxyList[1] = proxyList[4];

    bytes16[] memory newDataIdList = new bytes16[](2);
    newDataIdList[0] = batchValueIds[4];
    newDataIdList[1] = batchValueIds[3];

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(newProxyList[0], newDataIdList[0]);
    emit DataFeedsCache.ProxyDataIdUpdated(newProxyList[1], newDataIdList[1]);

    dataFeedsCache.updateDataIdMappingsForProxies(newProxyList, newDataIdList);

    decimals = 18;

    vm.startPrank(proxyList[3]);
    decimalsAns = dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);

    decimals = 8;

    vm.startPrank(proxyList[4]);
    decimalsAns = dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);
  }

  function test_updateDataIdMappingsForProxies_and_RevertOnWrongCaller() public {
    address[] memory proxyList = new address[](1);
    proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(proxyList[0], dataIdList[0]);

    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, dataIdList);

    uint8[] memory decimalsArr = new uint8[](1);
    decimalsArr[0] = 8;

    vm.startPrank(ILLEGAL_CALLER);
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.NoMappingForSender.selector, ILLEGAL_CALLER));

    dataFeedsCache.decimals();
  }

  function test_removeDataIdMappingsForProxiesSuccess() public {
    address[] memory proxyList = new address[](1);
    proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(proxyList[0], dataIdList[0]);

    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, dataIdList);

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdRemoved(proxyList[0], dataIdList[0]);

    dataFeedsCache.removeDataIdMappingsForProxies(proxyList);
  }

  function test_removeDataIdMappingsForProxiesSuccess_and_call_decimals() public {
    address[] memory proxyList = new address[](1);
    proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(proxyList[0], dataIdList[0]);

    dataFeedsCache.updateDataIdMappingsForProxies(proxyList, dataIdList);

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdRemoved(proxyList[0], dataIdList[0]);

    dataFeedsCache.removeDataIdMappingsForProxies(proxyList);

    uint8[] memory decimalsArr = new uint8[](1);
    decimalsArr[0] = 8;

    vm.startPrank(proxyList[0]);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.NoMappingForSender.selector, proxyList[0]));

    dataFeedsCache.decimals();
  }

  function test_supportsInterface() public view {
    assertEq(dataFeedsCache.supportsInterface(type(IDataFeedsCache).interfaceId), true);
  }

  function test_setFeedConfigsRevertEmptyConfig() public {
    // empty data ids
    bytes16[] memory dataIds = new bytes16[](0);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, workflowMetadata);

    // empty workflows
    dataIds = new bytes16[](1);
    dataIds[0] = bytes16(0);
    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, _workflowMetadata);
  }

  function test_setFeedConfigsRevertZeroDataId() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(0);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidDataId.selector));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidDataId.selector));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, workflowMetadata);
  }

  function test_setFeedConfigsRevertInvalidConfigsLengthDescriptions() public {
    // description has length of 1
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.ArrayLengthMismatch.selector));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.ArrayLengthMismatch.selector));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals2By1, workflowMetadata);
  }

  function test_setBundleFeedConfigsRevertInvalidConfigsLengthDecimals() public {
    // decimals has length of 1
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    string[] memory _descriptions = new string[](2);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.ArrayLengthMismatch.selector));
    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals1By1, workflowMetadata);
  }

  function test_setFeedConfigsRevertUnauthorizedFeedAdmin() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");
    vm.startPrank(address(123));

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.UnauthorizedCaller.selector, address(123)));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.UnauthorizedCaller.selector, address(123)));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, workflowMetadata);
  }

  function test_setFeedConfigsRevertInvalidWorkflowMetadata() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");

    // 0 address sender
    DataFeedsCache.WorkflowMetadata memory wfWithInvalidSender = DataFeedsCache.WorkflowMetadata({
      allowedSender: address(0),
      allowedWorkflowOwner: allowedWorkflowOwnersList[0],
      allowedWorkflowName: allowedWorkflowNamesList[0]
    });

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = new DataFeedsCache.WorkflowMetadata[](1);
    _workflowMetadata[0] = wfWithInvalidSender;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, _workflowMetadata);

    // 0 address owner
    DataFeedsCache.WorkflowMetadata memory wfWithInvalidOwner = DataFeedsCache.WorkflowMetadata({
      allowedSender: allowedSendersList[0],
      allowedWorkflowOwner: address(0),
      allowedWorkflowName: allowedWorkflowNamesList[0]
    });
    _workflowMetadata[0] = wfWithInvalidOwner;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, _workflowMetadata);

    // 0 address name
    DataFeedsCache.WorkflowMetadata memory wfWithInvalidName = DataFeedsCache.WorkflowMetadata({
      allowedSender: allowedSendersList[0],
      allowedWorkflowOwner: allowedWorkflowOwnersList[0],
      allowedWorkflowName: bytes10(0)
    });
    _workflowMetadata[0] = wfWithInvalidName;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidWorkflowName.selector, address(0)));
    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidWorkflowName.selector, address(0)));
    dataFeedsCache.setBundleFeedConfigs(dataIds, descriptions, decimals1By1, _workflowMetadata);
  }

  function test_setFeedConfigsSuccess() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");

    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[0],
      decimals: 0,
      description: descriptions[0],
      workflowMetadata: workflowMetadata
    });

    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);
  }

  function test_setDecimalFeedConfigs_setAgainWithClear() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadataNew = new DataFeedsCache.WorkflowMetadata[](3);
    _workflowMetadataNew[0] = workflowMetadata[1];
    _workflowMetadataNew[1] = workflowMetadata[0];
    _workflowMetadataNew[2] = workflowMetadata[1];

    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[0],
      decimals: 0,
      description: _descriptions[0],
      workflowMetadata: workflowMetadata
    });
    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[1],
      decimals: 0,
      description: _descriptions[1],
      workflowMetadata: workflowMetadata
    });

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, workflowMetadata[1].allowedSender);

    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, workflowMetadata[1].allowedSender);

    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved({dataId: dataIds[0]});
    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[0],
      decimals: 0,
      description: _descriptions[0],
      workflowMetadata: _workflowMetadataNew
    });
    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved({dataId: dataIds[1]});
    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[1],
      decimals: 0,
      description: _descriptions[1],
      workflowMetadata: _workflowMetadataNew
    });

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, _workflowMetadataNew);

    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

    assertEq(_workflowMetadata.length, 3);
    assertEq(_workflowMetadataNew[0].allowedWorkflowName, _workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadataNew[0].allowedWorkflowOwner, _workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[0].allowedSender, _workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadataNew[1].allowedWorkflowName, _workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadataNew[1].allowedWorkflowOwner, _workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[1].allowedSender, _workflowMetadata[1].allowedSender);

    assertEq(_workflowMetadataNew[2].allowedWorkflowName, _workflowMetadata[2].allowedWorkflowName);
    assertEq(_workflowMetadataNew[2].allowedWorkflowOwner, _workflowMetadata[2].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[2].allowedSender, _workflowMetadata[2].allowedSender);

    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

    assertEq(_workflowMetadata.length, 3);
    assertEq(_workflowMetadataNew[0].allowedWorkflowName, _workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadataNew[0].allowedWorkflowOwner, _workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[0].allowedSender, _workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadataNew[1].allowedWorkflowName, _workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadataNew[1].allowedWorkflowOwner, _workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[1].allowedSender, _workflowMetadata[1].allowedSender);

    assertEq(_workflowMetadataNew[2].allowedWorkflowName, _workflowMetadata[2].allowedWorkflowName);
    assertEq(_workflowMetadataNew[2].allowedWorkflowOwner, _workflowMetadata[2].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[2].allowedSender, _workflowMetadata[2].allowedSender);
  }

  function test_setBundleFeedConfigs_setAgainWithClear() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadataNew = new DataFeedsCache.WorkflowMetadata[](3);
    _workflowMetadataNew[0] = workflowMetadata[1];
    _workflowMetadataNew[1] = workflowMetadata[0];
    _workflowMetadataNew[2] = workflowMetadata[1];

    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[0],
      decimals: decimals2By1[0],
      description: _descriptions[0],
      workflowMetadata: workflowMetadata
    });
    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[1],
      decimals: decimals2By1[1],
      description: _descriptions[1],
      workflowMetadata: workflowMetadata
    });

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals2By1, workflowMetadata);

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, workflowMetadata[1].allowedSender);

    uint8[] memory decimalsArr = dataFeedsCache.getBundleDecimals(dataIds[0]);

    assertEq(decimalsArr.length, decimals2By1[0].length);
    assertEq(decimalsArr[0], decimals2By1[0][0]);

    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, workflowMetadata[1].allowedSender);

    decimalsArr = dataFeedsCache.getBundleDecimals(dataIds[1]);

    assertEq(decimalsArr.length, decimals2By1[1].length);
    assertEq(decimalsArr[0], decimals2By1[1][0]);

    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved({dataId: dataIds[0]});
    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[0],
      decimals: decimals2By2[0],
      description: _descriptions[0],
      workflowMetadata: _workflowMetadataNew
    });
    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved({dataId: dataIds[1]});
    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[1],
      decimals: decimals2By2[1],
      description: _descriptions[1],
      workflowMetadata: _workflowMetadataNew
    });

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals2By2, _workflowMetadataNew);

    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

    assertEq(_workflowMetadata.length, 3);
    assertEq(_workflowMetadataNew[0].allowedWorkflowName, _workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadataNew[0].allowedWorkflowOwner, _workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[0].allowedSender, _workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadataNew[1].allowedWorkflowName, _workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadataNew[1].allowedWorkflowOwner, _workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[1].allowedSender, _workflowMetadata[1].allowedSender);

    assertEq(_workflowMetadataNew[2].allowedWorkflowName, _workflowMetadata[2].allowedWorkflowName);
    assertEq(_workflowMetadataNew[2].allowedWorkflowOwner, _workflowMetadata[2].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[2].allowedSender, _workflowMetadata[2].allowedSender);

    decimalsArr = dataFeedsCache.getBundleDecimals(dataIds[0]);

    assertEq(decimalsArr.length, decimals2By2[0].length);
    assertEq(decimalsArr[0], decimals2By2[0][0]);
    assertEq(decimalsArr[1], decimals2By2[0][1]);

    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

    assertEq(_workflowMetadata.length, 3);
    assertEq(_workflowMetadataNew[0].allowedWorkflowName, _workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadataNew[0].allowedWorkflowOwner, _workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[0].allowedSender, _workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadataNew[1].allowedWorkflowName, _workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadataNew[1].allowedWorkflowOwner, _workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[1].allowedSender, _workflowMetadata[1].allowedSender);

    assertEq(_workflowMetadataNew[2].allowedWorkflowName, _workflowMetadata[2].allowedWorkflowName);
    assertEq(_workflowMetadataNew[2].allowedWorkflowOwner, _workflowMetadata[2].allowedWorkflowOwner);
    assertEq(_workflowMetadataNew[2].allowedSender, _workflowMetadata[2].allowedSender);

    decimalsArr = dataFeedsCache.getBundleDecimals(dataIds[1]);

    assertEq(decimalsArr.length, decimals2By2[1].length);
    assertEq(decimalsArr[0], decimals2By2[1][0]);
    assertEq(decimalsArr[1], decimals2By2[1][1]);
  }

  function test_description() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(address(dataFeedsAggregatorProxy));
    string memory description = dataFeedsCache.description();

    assertEq(descriptions[0], description);
  }

  function test_decimals() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(address(dataFeedsAggregatorProxy));
    uint8 decimals = dataFeedsCache.decimals();
    assertEq(18, decimals);
  }

  function test_bundleDecimals() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];
    uint8[][] memory _decimals = new uint8[][](1);
    _decimals[0] = new uint8[](2);
    _decimals[0][0] = 18;
    _decimals[0][1] = 8;

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, _decimals, workflowMetadata);

    vm.startPrank(address(dataFeedsAggregatorProxy));
    uint8[] memory decimals = dataFeedsCache.bundleDecimals();
    assertEq(decimals.length, 2);
    assertEq(decimals[0], 18);
    assertEq(decimals[1], 8);
  }

  function test_getFeedMetadataRevertFeedNotConfigured() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.FeedNotConfigured.selector, bytes16(0)));
    dataFeedsCache.getFeedMetadata(bytes16(0), 0, 1);
  }

  function test_getFeedMetadata() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");

    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    // limit less than the number of elements
    // first slice
    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 1);

    assertEq(_workflowMetadata.length, 1);
    assertEq(_workflowMetadata[0].allowedWorkflowName, allowedWorkflowNamesList[0]);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, allowedWorkflowOwnersList[0]);
    assertEq(_workflowMetadata[0].allowedSender, allowedSendersList[0]);

    // second slice
    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 1, 1);

    assertEq(_workflowMetadata.length, 1);
    assertEq(_workflowMetadata[0].allowedWorkflowName, allowedWorkflowNamesList[1]);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, allowedWorkflowOwnersList[1]);
    assertEq(_workflowMetadata[0].allowedSender, allowedSendersList[1]);

    // returns the full array if the maxCount is equal to the number of elements
    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, workflowMetadata.length);
    assertEq(_workflowMetadata.length, 2);

    // returns the full array if the number of elements is less than the maxCount
    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 100);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, allowedWorkflowNamesList[0]);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, allowedWorkflowOwnersList[0]);
    assertEq(_workflowMetadata[0].allowedSender, allowedSendersList[0]);

    assertEq(_workflowMetadata[1].allowedWorkflowName, allowedWorkflowNamesList[1]);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, allowedWorkflowOwnersList[1]);
    assertEq(_workflowMetadata[1].allowedSender, allowedSendersList[1]);

    // returns the full array if maxCount is 0
    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);
    assertEq(_workflowMetadata.length, 2);

    // returns empty array if the cursor is out of bounds
    _workflowMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 2, 1);
    assertEq(_workflowMetadata.length, 0);
  }

  function test_getWorkflowMetaData() public view {
    (address _workflowOwner, bytes10 _workflowName) = dataFeedsCache.getWorkflowMetaData(METADATA);

    assertEq(_workflowName, workflowName);
    assertEq(_workflowOwner, workflowOwner);
  }

  function test_getDataType() public view {
    bytes1 dataType = dataFeedsCache.getDataType(bytes16(dataId1), 7);
    assertEq(dataType, hex"32");
  }

  function testFuzzy_getDataType(bytes16 id, uint256 index) public view {
    vm.assume(index < 16);
    bytes1 expected = bytes1(uint8(id[index]));
    bytes1 result = dataFeedsCache.getDataType(id, index);
    assertEq(result, expected);
  }

  function testFuzzy_getDataTypeRevertOutOfBound(bytes16 id, uint256 index) public {
    vm.assume(index >= 16);
    vm.expectRevert();
    dataFeedsCache.getDataType(id, index);
  }

  function testFuzz_createReportHash(
    bytes16 dataId,
    address sender,
    address fuzzedWorkflowOwner,
    bytes10 fuzzedWorkflowName
  ) public view {
    bytes32 reportHash = dataFeedsCache.createReportHash(dataId, sender, fuzzedWorkflowOwner, fuzzedWorkflowName);
    bytes32 expectedReportHash = keccak256(abi.encode(dataId, sender, fuzzedWorkflowOwner, fuzzedWorkflowName));
    assertEq(reportHash, expectedReportHash);
  }

  function test_setFeedAdminRevertZeroAddress() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));

    dataFeedsCache.setFeedAdmin(address(0), true);
  }

  function testFuzz_setFeedAdminSuccess(
    address feedAdmin
  ) public {
    vm.assume(feedAdmin != address(0));
    vm.assume(feedAdmin != OWNER);
    vm.expectEmit();
    emit DataFeedsCache.FeedAdminSet(feedAdmin, true);

    dataFeedsCache.setFeedAdmin(feedAdmin, true);
  }

  function test_isFeedAdmin() public view {
    assertEq(dataFeedsCache.isFeedAdmin(OWNER), true);
    assertEq(dataFeedsCache.isFeedAdmin(address(10002)), false);
  }

  function test_removeFeedAdminSuccess() public {
    dataFeedsCache.setFeedAdmin(address(10003), true);
    vm.expectEmit();
    emit DataFeedsCache.FeedAdminSet(address(10003), false);
    dataFeedsCache.setFeedAdmin(address(10003), false);
  }

  function testFuzz_checkFeedPermissionFalse(
    bytes16 dataId,
    address sender,
    address fuzzedWorkflowOwner,
    bytes10 fuzzedWorkflowName
  ) public view {
    DataFeedsCache.WorkflowMetadata memory wfm = DataFeedsCache.WorkflowMetadata({
      allowedSender: sender,
      allowedWorkflowOwner: fuzzedWorkflowOwner,
      allowedWorkflowName: fuzzedWorkflowName
    });
    bool hasPermission = dataFeedsCache.checkFeedPermission(dataId, wfm);
    assertEq(hasPermission, false);
  }

  function testFuzz_checkFeedPermissionTrue(
    bytes16 dataId,
    address sender,
    address fuzzedWorkflowOwner,
    bytes10 fuzzedWorkflowName
  ) public {
    vm.assume(dataId != bytes16(0));
    vm.assume(sender != address(0));
    vm.assume(fuzzedWorkflowOwner != address(0));
    vm.assume(fuzzedWorkflowName != bytes10(0));

    DataFeedsCache.WorkflowMetadata memory _workflowMetadata1 = DataFeedsCache.WorkflowMetadata({
      allowedSender: sender,
      allowedWorkflowOwner: fuzzedWorkflowOwner,
      allowedWorkflowName: fuzzedWorkflowName
    });

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = new DataFeedsCache.WorkflowMetadata[](1);
    _workflowMetadata[0] = _workflowMetadata1;

    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = dataId;

    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, _workflowMetadata);

    bool hasPermission = dataFeedsCache.checkFeedPermission(dataId, _workflowMetadata[0]);
    assertEq(hasPermission, true);
  }

  function test_onReportInvalidPermission() public {
    // Invalid sender
    vm.startPrank(ILLEGAL_CALLER);

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId1),
      sender: ILLEGAL_CALLER,
      workflowOwner: workflowOwner,
      workflowName: workflowName
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId2),
      sender: ILLEGAL_CALLER,
      workflowOwner: workflowOwner,
      workflowName: workflowName
    });

    dataFeedsCache.onReport(METADATA, decimalReportlength2);

    // Data id not configured
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16("1"); // onReport will send report for dataId1 and dataId2.

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    vm.stopPrank();
    vm.startPrank(OWNER);
    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(dataId1),
      roundId: 1,
      timestamp: timestamp1,
      answer: uint224(price3)
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId2),
      sender: REPORT_SENDER,
      workflowOwner: workflowOwner,
      workflowName: workflowName
    });

    vm.stopPrank();
    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength2);

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({
      dataId: bytes16(dataId1),
      timestamp: timestamp1,
      bundle: abi.encodePacked(abi.encode(price3), abi.encode(price4))
    });

    // missing data id for bundle report
    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId2),
      sender: REPORT_SENDER,
      workflowOwner: workflowOwner,
      workflowName: workflowName
    });

    dataFeedsCache.onReport(METADATA, bundleReportlength2);
  }

  function test_onReportStaleDecimalReport() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength2);

    vm.expectEmit();
    emit DataFeedsCache.StaleDecimalReport({
      dataId: bytes16(dataId1),
      reportTimestamp: timestamp1 - 50,
      latestTimestamp: timestamp1
    });

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(dataId2),
      roundId: 2,
      timestamp: timestamp2 + 50,
      answer: uint224(price2)
    });

    dataFeedsCache.onReport(METADATA, staleReport);
  }

  function test_onReportStaleBundleReport() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals2By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, bundleReportlength2);

    vm.expectEmit();
    emit DataFeedsCache.StaleBundleReport({
      dataId: bytes16(dataId1),
      reportTimestamp: timestamp1 - 50,
      latestTimestamp: timestamp1
    });

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({
      dataId: bytes16(dataId2),
      timestamp: timestamp2 + 50,
      bundle: abi.encodePacked(abi.encode(price3), abi.encode(price4))
    });

    dataFeedsCache.onReport(METADATA, staleBundleReport);
  }

  function test_onReportRevertInvalidWorkflowName() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    // workflowName in report is 'abc'
    bytes10 invalidWorkflowName = bytes10("xyz");
    bytes memory thisMetadata = abi.encodePacked(workflowId, invalidWorkflowName, workflowOwner, reportId);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId1),
      sender: REPORT_SENDER,
      workflowOwner: workflowOwner,
      workflowName: invalidWorkflowName
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId2),
      sender: REPORT_SENDER,
      workflowOwner: workflowOwner,
      workflowName: invalidWorkflowName
    });

    dataFeedsCache.onReport(thisMetadata, decimalReportlength2);
  }

  function test_onReportRevertInvalidWorkflowOwner() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    // workFlowOwner in report is address(10004);
    address invalidWorkflowOwner = address(10005);
    bytes memory thisMetadata = abi.encodePacked(workflowId, workflowName, invalidWorkflowOwner, reportId);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId1),
      sender: REPORT_SENDER,
      workflowOwner: invalidWorkflowOwner,
      workflowName: workflowName
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(dataId2),
      sender: REPORT_SENDER,
      workflowOwner: invalidWorkflowOwner,
      workflowName: workflowName
    });

    dataFeedsCache.onReport(thisMetadata, decimalReportlength2);
  }

  function test_onReportSuccess_EmptyReport() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    vm.expectRevert();
    dataFeedsCache.onReport(METADATA, "");

    assertEq(dataFeedsCache.getLatestAnswer(dataIds[0]), int256(0));
  }

  function test_onReportSuccess_EmptyDecimalReport() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    dataFeedsCache.onReport(METADATA, emptyDecimalReport);

    assertEq(dataFeedsCache.getLatestAnswer(bytes16(dataId1)), int256(0));
    assertEq(dataFeedsCache.getLatestAnswer(bytes16(dataId2)), int256(0));
    assertEq(dataFeedsCache.getLatestAnswer(bytes16(dataId3)), int256(0));
    assertEq(dataFeedsCache.getLatestAnswer(bytes16(dataId4)), int256(0));
    assertEq(dataFeedsCache.getLatestAnswer(bytes16(dataId5)), int256(0));
  }

  function test_onReportSuccess_DecimalReportLength1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(dataId1),
      roundId: 1,
      timestamp: timestamp1,
      answer: uint224(price1)
    });

    dataFeedsCache.onReport(METADATA, decimalReportlength1);

    assertEq(dataFeedsCache.getLatestAnswer(dataIds[0]), int256(price1));
  }

  function test_onReportSuccess_DecimalReportLength2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(dataId1),
      roundId: 1,
      timestamp: timestamp1,
      answer: uint224(price3)
    });

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(dataId2),
      roundId: 1,
      timestamp: timestamp2,
      answer: uint224(price4)
    });

    dataFeedsCache.onReport(METADATA, decimalReportlength2);
  }

  function test_onReportSuccess_EmptyBundleReport() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    dataFeedsCache.onReport(METADATA, emptyBundleReport);

    assertEq(dataFeedsCache.getLatestBundle(bytes16(dataId1)), "");
    assertEq(dataFeedsCache.getLatestBundle(bytes16(dataId2)), "");
    assertEq(dataFeedsCache.getLatestBundle(bytes16(dataId3)), "");
    assertEq(dataFeedsCache.getLatestBundle(bytes16(dataId4)), "");
    assertEq(dataFeedsCache.getLatestBundle(bytes16(dataId5)), "");
  }

  function test_onReportSuccess_BundleReportLength1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals1By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    bytes memory expectedBundle =
      hex"000000000000000000000000000000000000000000000000000000000001e240000000000000000000000000000000000000000000000000000000000006f855";

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({dataId: bytes16(dataId1), timestamp: timestamp1, bundle: expectedBundle});

    dataFeedsCache.onReport(METADATA, bundleReportlength1);
  }

  function test_onReportSuccess_BundleReportLength2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals2By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    bytes memory expectedBundle1 =
      hex"00000000000000000000000000000000000000000000000000000000000c0bd000000000000000000000000000000000000000000000000000000000000d950b";

    bytes memory expectedBundle2 =
      hex"000000000000000000000000000000000000000000000000000000000009fbf100000000000000000000000000000000000000000000000000000000000f1206";

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({dataId: bytes16(dataId1), timestamp: timestamp1, bundle: expectedBundle1});

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({dataId: bytes16(dataId2), timestamp: timestamp2, bundle: expectedBundle2});

    dataFeedsCache.onReport(METADATA, bundleReportlength2);
  }

  function test_latestAnswer1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength1);

    vm.startPrank(proxyList[0]);
    int256 value = dataFeedsCache.latestAnswer();
    assertEq(value, int256(price1));
  }

  function test_latestAnswer2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength2);

    vm.startPrank(proxyList[0]);
    int256 value = dataFeedsCache.latestAnswer();
    assertEq(value, int256(price3));

    vm.startPrank(proxyList[1]);
    value = dataFeedsCache.latestAnswer();
    assertEq(value, int256(price4));
  }

  function test_getLatestAnswer1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength1);
    vm.stopPrank();

    int256 value = dataFeedsCache.getLatestAnswer(dataIds[0]);
    assertEq(value, int256(price1));
  }

  function test_getLatestAnswer2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength2);
    vm.stopPrank();

    int256 value = dataFeedsCache.getLatestAnswer(dataIds[0]);
    assertEq(value, int256(price3));

    value = dataFeedsCache.getLatestAnswer(dataIds[1]);
    assertEq(value, int256(price4));
  }

  function test_latestBundle1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals1By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, bundleReportlength1);

    vm.startPrank(proxyList[0]);
    uint256 roundId = dataFeedsCache.latestRound();

    bytes memory bundle = dataFeedsCache.latestBundle();
    uint256 timestamp = dataFeedsCache.latestBundleTimestamp();
    uint8[] memory decimals = dataFeedsCache.bundleDecimals();
    assertEq(bundle, abi.encode(price1, price2));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, price1);
    assertEq(firstBundleP2, price2);
    assertEq(timestamp, timestamp1);
    assertEq(decimals.length, decimals1By1[0].length);
    assertEq(decimals[0], decimals1By1[0][0]);
  }

  function test_latestBundle2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals2By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, bundleReportlength2);

    vm.startPrank(proxyList[0]);
    uint256 roundId = dataFeedsCache.latestRound();

    bytes memory bundle = dataFeedsCache.latestBundle();
    uint256 timestamp = dataFeedsCache.latestBundleTimestamp();
    uint8[] memory decimals = dataFeedsCache.bundleDecimals();
    assertEq(bundle, abi.encode(price3, price4));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, price3);
    assertEq(firstBundleP2, price4);
    assertEq(timestamp, timestamp1);
    assertEq(decimals.length, decimals2By1[0].length);
    assertEq(decimals[0], decimals2By1[0][0]);

    vm.startPrank(proxyList[1]);
    roundId = dataFeedsCache.latestRound();

    bundle = dataFeedsCache.latestBundle();
    timestamp = dataFeedsCache.latestBundleTimestamp();
    decimals = dataFeedsCache.bundleDecimals();
    assertEq(bundle, abi.encode(price5, price6));
    (uint256 secondBundleP1, uint256 secondBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(secondBundleP1, price5);
    assertEq(secondBundleP2, price6);
    assertEq(timestamp, timestamp2);
    assertEq(decimals.length, decimals2By1[1].length);
    assertEq(decimals[0], decimals2By1[1][0]);
  }

  function test_getLatestBundle1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals1By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, bundleReportlength1);
    vm.stopPrank();

    bytes memory bundle = dataFeedsCache.getLatestBundle(dataIds[0]);
    uint256 timestamp = dataFeedsCache.getLatestBundleTimestamp(dataIds[0]);
    uint8[] memory decimals = dataFeedsCache.getBundleDecimals(dataIds[0]);
    assertEq(bundle, abi.encode(price1, price2));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, price1);
    assertEq(firstBundleP2, price2);
    assertEq(timestamp, timestamp1);
    assertEq(decimals.length, decimals1By1[0].length);
    assertEq(decimals[0], decimals1By1[0][0]);
  }

  function test_getLatestBundle2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(dataId1);
    dataIds[1] = bytes16(dataId2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = descriptions[0];
    _descriptions[1] = descriptions[0];

    dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, decimals2By1, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, bundleReportlength2);
    vm.stopPrank();

    bytes memory bundle = dataFeedsCache.getLatestBundle(dataIds[0]);
    uint256 timestamp = dataFeedsCache.getLatestBundleTimestamp(dataIds[0]);
    uint8[] memory decimals = dataFeedsCache.getBundleDecimals(dataIds[0]);
    assertEq(bundle, abi.encode(price3, price4));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, price3);
    assertEq(firstBundleP2, price4);
    assertEq(timestamp, timestamp1);
    assertEq(decimals.length, decimals2By1[0].length);
    assertEq(decimals[0], decimals2By1[0][0]);

    bundle = dataFeedsCache.getLatestBundle(dataIds[1]);
    timestamp = dataFeedsCache.getLatestBundleTimestamp(dataIds[1]);
    decimals = dataFeedsCache.getBundleDecimals(dataIds[1]);
    assertEq(bundle, abi.encode(price5, price6));
    (uint256 secondBundleP1, uint256 secondBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(secondBundleP1, price5);
    assertEq(secondBundleP2, price6);
    assertEq(timestamp, timestamp2);
    assertEq(decimals.length, decimals2By1[1].length);
    assertEq(decimals[0], decimals2By1[1][0]);
  }

  function test_removeFeedsRevertInvalidSender() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    vm.startPrank(address(1002));

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.UnauthorizedCaller.selector, address(1002)));
    dataFeedsCache.removeFeedConfigs(dataIds);
  }

  function test_removeFeedsRevertNotConfiguredFeed() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    dataFeedsCache.setFeedAdmin(OWNER, true);

    vm.stopPrank();
    vm.startPrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.FeedNotConfigured.selector, dataIds[0]));
    dataFeedsCache.removeFeedConfigs(dataIds);
  }

  function test_removeFeedsSuccess() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);

    DataFeedsCache.WorkflowMetadata[] memory wfMetadata;

    dataFeedsCache.setDecimalFeedConfigs(dataIds, descriptions, workflowMetadata);

    wfMetadata = dataFeedsCache.getFeedMetadata(dataIds[0], 0, 2);
    assertEq(wfMetadata.length, 2);
    bool hasPermission = dataFeedsCache.checkFeedPermission(dataIds[0], wfMetadata[0]);
    assertEq(hasPermission, true);

    dataFeedsCache.setFeedAdmin(OWNER, true);

    vm.stopPrank();
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved(dataIds[0]);
    dataFeedsCache.removeFeedConfigs(dataIds);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.FeedNotConfigured.selector, dataIds[0]));
    dataFeedsCache.getFeedMetadata(dataIds[0], 0, 2);
    hasPermission = dataFeedsCache.checkFeedPermission(dataIds[0], wfMetadata[0]);
    assertEq(hasPermission, false);
  }

  function test_getDataIdForProxy() public view {
    bytes16 dataId = dataFeedsCache.getDataIdForProxy(proxyList[0]);
    assertEq(dataId, bytes16(dataId1));
  }

  function test_recoverTokensRevertUnauthorized() public {
    vm.startPrank(ILLEGAL_CALLER);

    vm.expectRevert("Only callable by owner");
    dataFeedsCache.recoverTokens(IERC20(address(s_link)), address(10008), 1 ether);
  }

  function test_recoverTokensERC20RevertNoBalance() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InsufficientBalance.selector, 0, 1));
    dataFeedsCache.recoverTokens(IERC20(address(s_link)), address(10007), 1);
  }

  function testFuzzy_recoverTokensERC20Success(
    uint256 amount
  ) public {
    vm.assume(amount > 0);
    s_link.mint(address(dataFeedsCache), amount);

    vm.expectEmit();
    emit DataFeedsCache.TokenRecovered(address(s_link), address(10008), amount);
    dataFeedsCache.recoverTokens(IERC20(address(s_link)), address(10008), amount);
    assertEq(s_link.balanceOf(address(10008)), amount);
    assertEq(s_link.balanceOf(address(dataFeedsCache)), 0);
  }

  function test_recoverTokensNativeRevertNoBalance() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InsufficientBalance.selector, 0, 1 ether));
    dataFeedsCache.recoverTokens(IERC20(address(0)), address(10007), 1 ether);
  }

  function testFuzzy_recoverTokensNativeSuccess(
    uint256 amount
  ) public {
    vm.assume(amount > 0);
    vm.deal(address(dataFeedsCache), amount);
    assertEq(address(dataFeedsCache).balance, amount);

    vm.expectEmit();
    emit DataFeedsCache.TokenRecovered(address(0), address(10007), amount);
    dataFeedsCache.recoverTokens(IERC20(address(0)), address(10007), amount);
    assertEq(address(dataFeedsCache).balance, 0);
    assertEq(address(10007).balance, amount);
  }

  function test_getLatestByFeedId() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(dataId1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = descriptions[0];

    dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    dataFeedsCache.onReport(METADATA, decimalReportlength1);

    uint256 timestamp = dataFeedsCache.getLatestTimestamp(dataIds[0]);
    assertEq(timestamp, timestamp1);

    (uint80 roundId, int256 answer, uint256 timestamp2, uint256 timestamp3, uint80 roundId2) =
      dataFeedsCache.getLatestRoundData(dataIds[0]);
    assertEq(roundId, 1);
    assertEq(roundId2, 1);
    assertEq(answer, int256(price1));
    assertEq(timestamp, timestamp2);
    assertEq(timestamp, timestamp3);

    uint8 decimals = dataFeedsCache.getDecimals(dataIds[0]);
    assertEq(decimals, 18);

    string memory description = dataFeedsCache.getDescription(dataIds[0]);
    assertEq(description, descriptions[0]);
  }
}

contract DataFeedsCacheHarness is DataFeedsCache {
  function getWorkflowMetaData(
    bytes calldata metadata
  ) public pure returns (address workflowOwner, bytes10 _workflowName) {
    return _getWorkflowMetaData(metadata);
  }

  function getDataType(bytes16 id, uint256 index) public pure returns (bytes1) {
    return _getDataType(id, index);
  }

  function createReportHash(
    bytes16 dataId,
    address sender,
    address _workflowOwner,
    bytes10 _workflowName
  ) public pure returns (bytes32) {
    return _createReportHash(dataId, sender, _workflowOwner, _workflowName);
  }
}
