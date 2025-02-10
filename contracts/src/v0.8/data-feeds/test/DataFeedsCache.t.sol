// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ERC20Mock} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {IERC20Metadata as IERC20} from
  "../../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {BundleAggregatorProxy} from "../BundleAggregatorProxy.sol";

import {DataFeedsCache} from "../DataFeedsCache.sol";
import {IDataFeedsCache} from "../interfaces/IDataFeedsCache.sol";
import {BaseTest} from "./BaseTest.t.sol";

// solhint-disable-next-line max-states-count
contract DataFeedsCacheTest is BaseTest {
  BundleAggregatorProxy internal s_dataFeedsAggregatorProxy;
  DataFeedsCacheHarness internal s_dataFeedsCache;

  address internal constant ILLEGAL_CALLER = address(11111); // address used as incorrect caller in tests
  address internal constant REPORT_SENDER = address(12222); // mocks keystone forwarder address

  ERC20Mock internal s_link = new ERC20Mock("LINK", "LINK", OWNER, 0);

  bytes32 internal constant WORKFLOWID = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  bytes10 internal constant WORKFLOWNAME = bytes10("abc");
  address internal constant WORKFLOWOWNER = address(10004);
  bytes2 internal constant REPORTID = hex"0001";
  string[] internal s_descriptions = ["description"];

  uint8[][] internal s_decimals1By1 = new uint8[][](1);
  uint8[][] internal s_decimals2By1 = new uint8[][](2);
  uint8[][] internal s_decimals2By2 = new uint8[][](2);

  bytes internal constant METADATA = abi.encodePacked(WORKFLOWID, WORKFLOWNAME, WORKFLOWOWNER, REPORTID);

  address[] internal s_allowedSendersList = [REPORT_SENDER, REPORT_SENDER];
  address[] internal s_allowedWorkflowOwnersList = [address(10004), address(10005)];
  bytes10[] internal s_allowedWorkflowNamesList = [bytes10("abc"), bytes10("xyz")];

  address[] internal s_singleProxyList = new address[](1);
  address[] internal s_proxyList = new address[](5);
  address[] internal s_newSingleProxyList = new address[](1);
  address[] internal s_newProxyList = new address[](5);

  bytes16[] internal s_singleValueId = new bytes16[](1);
  bytes16[] internal s_batchValueIds = new bytes16[](5);

  DataFeedsCache.WorkflowMetadata internal s_workflowMetadata1 = DataFeedsCache.WorkflowMetadata({
    allowedSender: s_allowedSendersList[0],
    allowedWorkflowOwner: s_allowedWorkflowOwnersList[0],
    allowedWorkflowName: s_allowedWorkflowNamesList[0]
  });

  DataFeedsCache.WorkflowMetadata internal s_workflowMetadata2 = DataFeedsCache.WorkflowMetadata({
    allowedSender: s_allowedSendersList[1],
    allowedWorkflowOwner: s_allowedWorkflowOwnersList[1],
    allowedWorkflowName: s_allowedWorkflowNamesList[1]
  });

  DataFeedsCache.WorkflowMetadata[] internal s_workflowMetadata;

  bytes internal s_emptyDecimalReport;
  bytes internal s_decimalReportlength1;
  bytes internal s_decimalReportlength2;
  bytes internal s_emptyBundleReport;
  bytes internal s_bundleReportlength1;
  bytes internal s_bundleReportlength2;
  bytes internal s_staleReport;
  bytes internal s_staleBundleReport;
  bytes32 internal constant DATAID1 = hex"010e12d1e0000032000000000000000000000000000000000000000000000000";
  bytes32 internal constant DATAID2 = hex"01b476d70d000232000000000000000000000000000000000000000000000000";
  bytes32 internal constant DATAID3 = hex"0169bd6041000103000000000000000000000000000000000000000000000000";
  bytes32 internal constant DATAID4 = hex"010e12d1e0000028000000000000000000000000000000000000000000000000";
  bytes32 internal constant DATAID5 = hex"010e12d1e0000032000000000000000000000000000000000000000000000000";
  bytes16 internal constant DATA_ID_0 = bytes16(keccak256("12345"));
  bytes16 internal constant DATA_ID_1 = bytes16(keccak256("23456"));
  bytes16 internal constant DATA_ID_2 = bytes16(keccak256("34567"));
  bytes16 internal constant DATA_ID_3 = bytes16(keccak256("45678"));
  bytes16 internal constant DATA_ID_4 = bytes16(keccak256("56789"));
  bytes16 internal constant DATA_ID_5 = bytes16(keccak256("67890"));
  uint256 internal constant PRICE1 = 123456;
  uint256 internal constant PRICE2 = 456789;
  uint256 internal constant PRICE3 = 789456;
  uint256 internal constant PRICE4 = 890123;
  uint256 internal constant PRICE5 = 654321;
  uint256 internal constant PRICE6 = 987654;
  uint32 internal constant TIMESTAMP1 = 100;
  uint32 internal constant TIMESTAMP2 = 200;

  function setUp() public override {
    super.setUp();
    s_dataFeedsCache = new DataFeedsCacheHarness();
    s_dataFeedsCache.setFeedAdmin(OWNER, true);
    s_dataFeedsAggregatorProxy = new BundleAggregatorProxy(address(s_dataFeedsCache), OWNER);

    // reports should be encoded as calldata, which has offset and length
    s_emptyDecimalReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000000" // Length
    );

    // reports should be encoded as calldata, which has offset and length
    s_decimalReportlength1 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000001", // Length
      DATAID1,
      abi.encode(TIMESTAMP1),
      abi.encode(PRICE1)
    );

    // reports should be encoded as calldata, which has offset and length
    s_decimalReportlength2 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // Length
      DATAID1,
      abi.encode(TIMESTAMP1),
      abi.encode(PRICE3),
      DATAID2,
      abi.encode(TIMESTAMP2),
      abi.encode(PRICE4)
    );

    s_staleReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // Offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // Length
      DATAID1,
      abi.encode(TIMESTAMP1 - 50), // report 1 for DATAID1 is stale in this report
      abi.encode(PRICE1),
      DATAID2,
      abi.encode(TIMESTAMP2 + 50),
      abi.encode(PRICE2)
    );

    s_emptyBundleReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000000", // length
      hex"0000000000000000000000000000000000000000000000000000000000000000" // offset of ReportOne
    );

    s_bundleReportlength1 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000001", // length
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset of ReportOne
      DATAID1, // ReportOne FeedID
      abi.encode(TIMESTAMP1),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
      abi.encode(PRICE1),
      abi.encode(PRICE2)
    );

    s_bundleReportlength2 = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // length
      hex"0000000000000000000000000000000000000000000000000000000000000040", // offset of ReportOne
      hex"0000000000000000000000000000000000000000000000000000000000000100", // offset of ReportTwo
      DATAID1, // ReportOne FeedID
      abi.encode(TIMESTAMP1),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
      abi.encode(PRICE3),
      abi.encode(PRICE4),
      DATAID2, // ReportTwo FeedID
      abi.encode(TIMESTAMP2),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportTwo Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportTwo Bundle
      abi.encode(PRICE5),
      abi.encode(PRICE6)
    );

    s_staleBundleReport = abi.encodePacked(
      hex"0000000000000000000000000000000000000000000000000000000000000020", // offset
      hex"0000000000000000000000000000000000000000000000000000000000000002", // length
      hex"0000000000000000000000000000000000000000000000000000000000000040", // offset of ReportOne
      hex"0000000000000000000000000000000000000000000000000000000000000100", // offset of ReportTwo
      DATAID1, // ReportOne FeedID
      abi.encode(TIMESTAMP1 - 50), // report is stale
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportOne Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportOne Bundle
      abi.encode(PRICE1),
      abi.encode(PRICE2),
      DATAID2, // ReportTwo FeedID
      abi.encode(TIMESTAMP2 + 50),
      hex"0000000000000000000000000000000000000000000000000000000000000060", // offset of ReportTwo Bundle
      hex"0000000000000000000000000000000000000000000000000000000000000040", // length of ReportTwo Bundle
      abi.encode(PRICE3),
      abi.encode(PRICE4)
    );

    s_workflowMetadata.push(s_workflowMetadata1);
    s_workflowMetadata.push(s_workflowMetadata2);

    s_singleProxyList[0] = address(10002);

    s_proxyList[0] = address(s_dataFeedsAggregatorProxy);
    s_proxyList[1] = address(10002);
    s_proxyList[2] = address(10004);
    s_proxyList[3] = address(10005);
    s_proxyList[4] = address(10006);

    s_newSingleProxyList[0] = address(10007);

    s_newProxyList[0] = address(10002);
    s_newProxyList[1] = address(10003);
    s_newProxyList[2] = address(10004);
    s_newProxyList[3] = address(10005);
    s_newProxyList[4] = address(10006);

    s_singleValueId = new bytes16[](1);
    s_singleValueId[0] = bytes16(DATAID1);

    s_batchValueIds = new bytes16[](5);
    s_batchValueIds[0] = bytes16(DATAID1);
    s_batchValueIds[1] = bytes16(DATAID2);
    s_batchValueIds[2] = bytes16(DATAID3);
    s_batchValueIds[3] = bytes16(DATAID4);
    s_batchValueIds[4] = bytes16(DATAID5);

    s_decimals1By1[0] = new uint8[](1);
    s_decimals1By1[0][0] = 18;

    s_decimals2By1[0] = new uint8[](1);
    s_decimals2By1[0][0] = 18;
    s_decimals2By1[1] = new uint8[](1);
    s_decimals2By1[1][0] = 8;

    s_decimals2By2[0] = new uint8[](2);
    s_decimals2By2[0][0] = 6;
    s_decimals2By2[0][1] = 12;
    s_decimals2By2[1] = new uint8[](2);
    s_decimals2By2[1][0] = 18;
    s_decimals2By2[1][0] = 8;

    vm.startPrank(OWNER);
    s_dataFeedsCache.setFeedAdmin(OWNER, true);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, s_batchValueIds);
  }

  function test_updateDataIdMappingsForProxiesRevertInvalidLengths() public {
    address[] memory s_proxyList = new address[](1);
    s_proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](2);
    dataIdList[0] = bytes16(keccak256("12345"));
    dataIdList[1] = bytes16(keccak256("67890"));

    vm.expectRevert(DataFeedsCache.ArrayLengthMismatch.selector);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, dataIdList);
  }

  function test_updateDataIdMappingsForProxiesRevertUnauthorizedOwner() public {
    address[] memory s_proxyList = new address[](1);
    s_proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.stopPrank();
    vm.startPrank(ILLEGAL_CALLER);
    vm.expectRevert(
      abi.encodeWithSelector(
        DataFeedsCache.UnauthorizedCaller.selector, address(0x0000000000000000000000000000000000002B67)
      )
    );
    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, dataIdList);
  }

  function test_updateDataIdMappingsForProxiesSuccess() public {
    address[] memory s_proxyList = new address[](1);
    s_proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(s_proxyList[0], dataIdList[0]);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, dataIdList);
  }

  function test_updateDataIdMappingsForProxies_and_call_decimals() public {
    uint8 decimals = 8;

    vm.startPrank(s_proxyList[3]);
    uint8 decimalsAns = s_dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);

    decimals = 18;

    vm.startPrank(s_proxyList[4]);
    decimalsAns = s_dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);

    address[] memory s_newProxyList = new address[](2);
    s_newProxyList[0] = s_proxyList[3];
    s_newProxyList[1] = s_proxyList[4];

    bytes16[] memory newDataIdList = new bytes16[](2);
    newDataIdList[0] = s_batchValueIds[4];
    newDataIdList[1] = s_batchValueIds[3];

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(s_newProxyList[0], newDataIdList[0]);
    emit DataFeedsCache.ProxyDataIdUpdated(s_newProxyList[1], newDataIdList[1]);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_newProxyList, newDataIdList);

    decimals = 18;

    vm.startPrank(s_proxyList[3]);
    decimalsAns = s_dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);

    decimals = 8;

    vm.startPrank(s_proxyList[4]);
    decimalsAns = s_dataFeedsCache.decimals();

    assertEq(decimalsAns, decimals);
  }

  function test_updateDataIdMappingsForProxies_and_RevertOnWrongCaller() public {
    address[] memory s_proxyList = new address[](1);
    s_proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(s_proxyList[0], dataIdList[0]);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, dataIdList);

    uint8[] memory decimalsArr = new uint8[](1);
    decimalsArr[0] = 8;

    vm.startPrank(ILLEGAL_CALLER);
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.NoMappingForSender.selector, ILLEGAL_CALLER));

    s_dataFeedsCache.decimals();
  }

  function test_removeDataIdMappingsForProxiesSuccess() public {
    address[] memory s_proxyList = new address[](1);
    s_proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(s_proxyList[0], dataIdList[0]);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, dataIdList);

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdRemoved(s_proxyList[0], dataIdList[0]);

    s_dataFeedsCache.removeDataIdMappingsForProxies(s_proxyList);
  }

  function test_removeDataIdMappingsForProxiesSuccess_and_call_decimals() public {
    address[] memory s_proxyList = new address[](1);
    s_proxyList[0] = address(10002);

    bytes16[] memory dataIdList = new bytes16[](1);
    dataIdList[0] = bytes16(keccak256("12345"));

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdUpdated(s_proxyList[0], dataIdList[0]);

    s_dataFeedsCache.updateDataIdMappingsForProxies(s_proxyList, dataIdList);

    vm.expectEmit();
    emit DataFeedsCache.ProxyDataIdRemoved(s_proxyList[0], dataIdList[0]);

    s_dataFeedsCache.removeDataIdMappingsForProxies(s_proxyList);

    uint8[] memory decimalsArr = new uint8[](1);
    decimalsArr[0] = 8;

    vm.startPrank(s_proxyList[0]);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.NoMappingForSender.selector, s_proxyList[0]));

    s_dataFeedsCache.decimals();
  }

  function test_supportsInterface() public view {
    assertEq(s_dataFeedsCache.supportsInterface(type(IDataFeedsCache).interfaceId), true);
  }

  function test_setFeedConfigsRevertEmptyConfig() public {
    // empty data ids
    bytes16[] memory dataIds = new bytes16[](0);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, s_workflowMetadata);

    // empty workflows
    dataIds = new bytes16[](1);
    dataIds[0] = bytes16(0);
    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.EmptyConfig.selector));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, _workflowMetadata);
  }

  function test_setFeedConfigsRevertZeroDataId() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(0);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidDataId.selector));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidDataId.selector));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, s_workflowMetadata);
  }

  function test_setFeedConfigsRevertInvalidConfigsLengthDescriptions() public {
    // description has length of 1
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.ArrayLengthMismatch.selector));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.ArrayLengthMismatch.selector));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals2By1, s_workflowMetadata);
  }

  function test_setBundleFeedConfigsRevertInvalidConfigsLengthDecimals() public {
    // decimals has length of 1
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    string[] memory _descriptions = new string[](2);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.ArrayLengthMismatch.selector));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals1By1, s_workflowMetadata);
  }

  function test_setFeedConfigsRevertUnauthorizedFeedAdmin() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");
    vm.startPrank(address(123));

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.UnauthorizedCaller.selector, address(123)));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.UnauthorizedCaller.selector, address(123)));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, s_workflowMetadata);
  }

  function test_setFeedConfigsRevertInvalidWorkflowMetadata() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");

    // 0 address sender
    DataFeedsCache.WorkflowMetadata memory wfWithInvalidSender = DataFeedsCache.WorkflowMetadata({
      allowedSender: address(0),
      allowedWorkflowOwner: s_allowedWorkflowOwnersList[0],
      allowedWorkflowName: s_allowedWorkflowNamesList[0]
    });

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = new DataFeedsCache.WorkflowMetadata[](1);
    _workflowMetadata[0] = wfWithInvalidSender;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, _workflowMetadata);

    // 0 address owner
    DataFeedsCache.WorkflowMetadata memory wfWithInvalidOwner = DataFeedsCache.WorkflowMetadata({
      allowedSender: s_allowedSendersList[0],
      allowedWorkflowOwner: address(0),
      allowedWorkflowName: s_allowedWorkflowNamesList[0]
    });
    _workflowMetadata[0] = wfWithInvalidOwner;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, _workflowMetadata);

    // 0 address name
    DataFeedsCache.WorkflowMetadata memory wfWithInvalidName = DataFeedsCache.WorkflowMetadata({
      allowedSender: s_allowedSendersList[0],
      allowedWorkflowOwner: s_allowedWorkflowOwnersList[0],
      allowedWorkflowName: bytes10(0)
    });
    _workflowMetadata[0] = wfWithInvalidName;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidWorkflowName.selector, address(0)));
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, _workflowMetadata);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidWorkflowName.selector, address(0)));
    s_dataFeedsCache.setBundleFeedConfigs(dataIds, s_descriptions, s_decimals1By1, _workflowMetadata);
  }

  function test_setFeedConfigsSuccess() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");

    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[0],
      decimals: 0,
      description: s_descriptions[0],
      workflowMetadata: s_workflowMetadata
    });

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);
  }

  function test_setDecimalFeedConfigs_setAgainWithClear() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16("1");
    dataIds[1] = bytes16("2");

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadataNew = new DataFeedsCache.WorkflowMetadata[](3);
    _workflowMetadataNew[0] = s_workflowMetadata[1];
    _workflowMetadataNew[1] = s_workflowMetadata[0];
    _workflowMetadataNew[2] = s_workflowMetadata[1];

    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[0],
      decimals: 0,
      description: _descriptions[0],
      workflowMetadata: s_workflowMetadata
    });
    vm.expectEmit();
    emit DataFeedsCache.DecimalFeedConfigSet({
      dataId: dataIds[1],
      decimals: 0,
      description: _descriptions[1],
      workflowMetadata: s_workflowMetadata
    });

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, s_workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, s_workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, s_workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, s_workflowMetadata[1].allowedSender);

    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, s_workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, s_workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, s_workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, s_workflowMetadata[1].allowedSender);

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

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, _workflowMetadataNew);

    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

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

    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

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
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadataNew = new DataFeedsCache.WorkflowMetadata[](3);
    _workflowMetadataNew[0] = s_workflowMetadata[1];
    _workflowMetadataNew[1] = s_workflowMetadata[0];
    _workflowMetadataNew[2] = s_workflowMetadata[1];

    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[0],
      decimals: s_decimals2By1[0],
      description: _descriptions[0],
      workflowMetadata: s_workflowMetadata
    });
    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[1],
      decimals: s_decimals2By1[1],
      description: _descriptions[1],
      workflowMetadata: s_workflowMetadata
    });

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals2By1, s_workflowMetadata);

    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, s_workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, s_workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, s_workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, s_workflowMetadata[1].allowedSender);

    uint8[] memory decimalsArr = s_dataFeedsCache.getBundleDecimals(dataIds[0]);

    assertEq(decimalsArr.length, s_decimals2By1[0].length);
    assertEq(decimalsArr[0], s_decimals2By1[0][0]);

    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_workflowMetadata[0].allowedWorkflowName);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_workflowMetadata[0].allowedWorkflowOwner);
    assertEq(_workflowMetadata[0].allowedSender, s_workflowMetadata[0].allowedSender);

    assertEq(_workflowMetadata[1].allowedWorkflowName, s_workflowMetadata[1].allowedWorkflowName);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, s_workflowMetadata[1].allowedWorkflowOwner);
    assertEq(_workflowMetadata[1].allowedSender, s_workflowMetadata[1].allowedSender);

    decimalsArr = s_dataFeedsCache.getBundleDecimals(dataIds[1]);

    assertEq(decimalsArr.length, s_decimals2By1[1].length);
    assertEq(decimalsArr[0], s_decimals2By1[1][0]);

    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved({dataId: dataIds[0]});
    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[0],
      decimals: s_decimals2By2[0],
      description: _descriptions[0],
      workflowMetadata: _workflowMetadataNew
    });
    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved({dataId: dataIds[1]});
    vm.expectEmit();
    emit DataFeedsCache.BundleFeedConfigSet({
      dataId: dataIds[1],
      decimals: s_decimals2By2[1],
      description: _descriptions[1],
      workflowMetadata: _workflowMetadataNew
    });

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals2By2, _workflowMetadataNew);

    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);

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

    decimalsArr = s_dataFeedsCache.getBundleDecimals(dataIds[0]);

    assertEq(decimalsArr.length, s_decimals2By2[0].length);
    assertEq(decimalsArr[0], s_decimals2By2[0][0]);
    assertEq(decimalsArr[1], s_decimals2By2[0][1]);

    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[1], 0, 0);

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

    decimalsArr = s_dataFeedsCache.getBundleDecimals(dataIds[1]);

    assertEq(decimalsArr.length, s_decimals2By2[1].length);
    assertEq(decimalsArr[0], s_decimals2By2[1][0]);
    assertEq(decimalsArr[1], s_decimals2By2[1][1]);
  }

  function test_description() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(address(s_dataFeedsAggregatorProxy));
    string memory description = s_dataFeedsCache.description();

    assertEq(s_descriptions[0], description);
  }

  function test_decimals() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(address(s_dataFeedsAggregatorProxy));
    uint8 decimals = s_dataFeedsCache.decimals();
    assertEq(18, decimals);
  }

  function test_bundleDecimals() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];
    uint8[][] memory _decimals = new uint8[][](1);
    _decimals[0] = new uint8[](2);
    _decimals[0][0] = 18;
    _decimals[0][1] = 8;

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, _decimals, s_workflowMetadata);

    vm.startPrank(address(s_dataFeedsAggregatorProxy));
    uint8[] memory decimals = s_dataFeedsCache.bundleDecimals();
    assertEq(decimals.length, 2);
    assertEq(decimals[0], 18);
    assertEq(decimals[1], 8);
  }

  function test_getFeedMetadataRevertFeedNotConfigured() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.FeedNotConfigured.selector, bytes16(0)));
    s_dataFeedsCache.getFeedMetadata(bytes16(0), 0, 1);
  }

  function test_getFeedMetadata() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16("1");

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    // limit less than the number of elements
    // first slice
    DataFeedsCache.WorkflowMetadata[] memory _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 1);

    assertEq(_workflowMetadata.length, 1);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_allowedWorkflowNamesList[0]);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_allowedWorkflowOwnersList[0]);
    assertEq(_workflowMetadata[0].allowedSender, s_allowedSendersList[0]);

    // second slice
    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 1, 1);

    assertEq(_workflowMetadata.length, 1);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_allowedWorkflowNamesList[1]);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_allowedWorkflowOwnersList[1]);
    assertEq(_workflowMetadata[0].allowedSender, s_allowedSendersList[1]);

    // returns the full array if the maxCount is equal to the number of elements
    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, s_workflowMetadata.length);
    assertEq(_workflowMetadata.length, 2);

    // returns the full array if the number of elements is less than the maxCount
    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 100);

    assertEq(_workflowMetadata.length, 2);
    assertEq(_workflowMetadata[0].allowedWorkflowName, s_allowedWorkflowNamesList[0]);
    assertEq(_workflowMetadata[0].allowedWorkflowOwner, s_allowedWorkflowOwnersList[0]);
    assertEq(_workflowMetadata[0].allowedSender, s_allowedSendersList[0]);

    assertEq(_workflowMetadata[1].allowedWorkflowName, s_allowedWorkflowNamesList[1]);
    assertEq(_workflowMetadata[1].allowedWorkflowOwner, s_allowedWorkflowOwnersList[1]);
    assertEq(_workflowMetadata[1].allowedSender, s_allowedSendersList[1]);

    // returns the full array if maxCount is 0
    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 0);
    assertEq(_workflowMetadata.length, 2);

    // returns empty array if the cursor is out of bounds
    _workflowMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 2, 1);
    assertEq(_workflowMetadata.length, 0);
  }

  function test_getWorkflowMetaData() public view {
    (address _workflowOwner, bytes10 _workflowName) = s_dataFeedsCache.getWorkflowMetaData(METADATA);

    assertEq(_workflowName, WORKFLOWNAME);
    assertEq(_workflowOwner, WORKFLOWOWNER);
  }

  function test_getDataType() public view {
    bytes1 dataType = s_dataFeedsCache.getDataType(bytes16(DATAID1), 7);
    assertEq(dataType, hex"32");
  }

  function testFuzzy_getDataType(bytes16 id, uint256 index) public view {
    vm.assume(index < 16);
    bytes1 expected = bytes1(uint8(id[index]));
    bytes1 result = s_dataFeedsCache.getDataType(id, index);
    assertEq(result, expected);
  }

  function testFuzzy_getDataTypeRevertOutOfBound(bytes16 id, uint256 index) public {
    vm.assume(index >= 16);
    vm.expectRevert();
    s_dataFeedsCache.getDataType(id, index);
  }

  function testFuzz_createReportHash(
    bytes16 dataId,
    address sender,
    address fuzzedWorkflowOwner,
    bytes10 fuzzedWorkflowName
  ) public view {
    bytes32 reportHash = s_dataFeedsCache.createReportHash(dataId, sender, fuzzedWorkflowOwner, fuzzedWorkflowName);
    bytes32 expectedReportHash = keccak256(abi.encode(dataId, sender, fuzzedWorkflowOwner, fuzzedWorkflowName));
    assertEq(reportHash, expectedReportHash);
  }

  function test_setFeedAdminRevertZeroAddress() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InvalidAddress.selector, address(0)));

    s_dataFeedsCache.setFeedAdmin(address(0), true);
  }

  function testFuzz_setFeedAdminSuccess(
    address feedAdmin
  ) public {
    vm.assume(feedAdmin != address(0));
    vm.assume(feedAdmin != OWNER);
    vm.expectEmit();
    emit DataFeedsCache.FeedAdminSet(feedAdmin, true);

    s_dataFeedsCache.setFeedAdmin(feedAdmin, true);
  }

  function test_isFeedAdmin() public view {
    assertEq(s_dataFeedsCache.isFeedAdmin(OWNER), true);
    assertEq(s_dataFeedsCache.isFeedAdmin(address(10002)), false);
  }

  function test_removeFeedAdminSuccess() public {
    s_dataFeedsCache.setFeedAdmin(address(10003), true);
    vm.expectEmit();
    emit DataFeedsCache.FeedAdminSet(address(10003), false);
    s_dataFeedsCache.setFeedAdmin(address(10003), false);
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
    bool hasPermission = s_dataFeedsCache.checkFeedPermission(dataId, wfm);
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

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, _workflowMetadata);

    bool hasPermission = s_dataFeedsCache.checkFeedPermission(dataId, _workflowMetadata[0]);
    assertEq(hasPermission, true);
  }

  function test_onReportInvalidPermission() public {
    // Invalid sender
    vm.startPrank(ILLEGAL_CALLER);

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID1),
      sender: ILLEGAL_CALLER,
      workflowOwner: WORKFLOWOWNER,
      workflowName: WORKFLOWNAME
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID2),
      sender: ILLEGAL_CALLER,
      workflowOwner: WORKFLOWOWNER,
      workflowName: WORKFLOWNAME
    });

    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength2);

    // Data id not configured
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16("1"); // onReport will send report for DATAID1 and DATAID2.

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    vm.stopPrank();
    vm.startPrank(OWNER);
    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(DATAID1),
      roundId: 1,
      timestamp: TIMESTAMP1,
      answer: uint224(PRICE3)
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID2),
      sender: REPORT_SENDER,
      workflowOwner: WORKFLOWOWNER,
      workflowName: WORKFLOWNAME
    });

    vm.stopPrank();
    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength2);

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({
      dataId: bytes16(DATAID1),
      timestamp: TIMESTAMP1,
      bundle: abi.encodePacked(abi.encode(PRICE3), abi.encode(PRICE4))
    });

    // missing data id for bundle report
    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID2),
      sender: REPORT_SENDER,
      workflowOwner: WORKFLOWOWNER,
      workflowName: WORKFLOWNAME
    });

    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength2);
  }

  function test_onReportStaleDecimalReport() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength2);

    vm.expectEmit();
    emit DataFeedsCache.StaleDecimalReport({
      dataId: bytes16(DATAID1),
      reportTimestamp: TIMESTAMP1 - 50,
      latestTimestamp: TIMESTAMP1
    });

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(DATAID2),
      roundId: 2,
      timestamp: TIMESTAMP2 + 50,
      answer: uint224(PRICE2)
    });

    s_dataFeedsCache.onReport(METADATA, s_staleReport);
  }

  function test_onReportStaleBundleReport() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals2By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength2);

    vm.expectEmit();
    emit DataFeedsCache.StaleBundleReport({
      dataId: bytes16(DATAID1),
      reportTimestamp: TIMESTAMP1 - 50,
      latestTimestamp: TIMESTAMP1
    });

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({
      dataId: bytes16(DATAID2),
      timestamp: TIMESTAMP2 + 50,
      bundle: abi.encodePacked(abi.encode(PRICE3), abi.encode(PRICE4))
    });

    s_dataFeedsCache.onReport(METADATA, s_staleBundleReport);
  }

  function test_onReportRevertInvalidWorkflowName() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    // workflowName in report is 'abc'
    bytes10 invalidWorkflowName = bytes10("xyz");
    bytes memory thisMetadata = abi.encodePacked(WORKFLOWID, invalidWorkflowName, WORKFLOWOWNER, REPORTID);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID1),
      sender: REPORT_SENDER,
      workflowOwner: WORKFLOWOWNER,
      workflowName: invalidWorkflowName
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID2),
      sender: REPORT_SENDER,
      workflowOwner: WORKFLOWOWNER,
      workflowName: invalidWorkflowName
    });

    s_dataFeedsCache.onReport(thisMetadata, s_decimalReportlength2);
  }

  function test_onReportRevertInvalidWorkflowOwner() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    // workFlowOwner in report is address(10004);
    address invalidWorkflowOwner = address(10005);
    bytes memory thisMetadata = abi.encodePacked(WORKFLOWID, WORKFLOWNAME, invalidWorkflowOwner, REPORTID);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID1),
      sender: REPORT_SENDER,
      workflowOwner: invalidWorkflowOwner,
      workflowName: WORKFLOWNAME
    });

    vm.expectEmit();
    emit DataFeedsCache.InvalidUpdatePermission({
      dataId: bytes16(DATAID2),
      sender: REPORT_SENDER,
      workflowOwner: invalidWorkflowOwner,
      workflowName: WORKFLOWNAME
    });

    s_dataFeedsCache.onReport(thisMetadata, s_decimalReportlength2);
  }

  function test_onReportSuccess_EmptyReport() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    vm.expectRevert();
    s_dataFeedsCache.onReport(METADATA, "");

    assertEq(s_dataFeedsCache.getLatestAnswer(dataIds[0]), int256(0));
  }

  function test_onReportSuccess_EmptyDecimalReport() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    s_dataFeedsCache.onReport(METADATA, s_emptyDecimalReport);

    assertEq(s_dataFeedsCache.getLatestAnswer(bytes16(DATAID1)), int256(0));
    assertEq(s_dataFeedsCache.getLatestAnswer(bytes16(DATAID2)), int256(0));
    assertEq(s_dataFeedsCache.getLatestAnswer(bytes16(DATAID3)), int256(0));
    assertEq(s_dataFeedsCache.getLatestAnswer(bytes16(DATAID4)), int256(0));
    assertEq(s_dataFeedsCache.getLatestAnswer(bytes16(DATAID5)), int256(0));
  }

  function test_onReportSuccess_DecimalReportLength1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(DATAID1),
      roundId: 1,
      timestamp: TIMESTAMP1,
      answer: uint224(PRICE1)
    });

    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength1);

    assertEq(s_dataFeedsCache.getLatestAnswer(dataIds[0]), int256(PRICE1));
  }

  function test_onReportSuccess_DecimalReportLength2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);

    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(DATAID1),
      roundId: 1,
      timestamp: TIMESTAMP1,
      answer: uint224(PRICE3)
    });

    vm.expectEmit();
    emit DataFeedsCache.DecimalReportUpdated({
      dataId: bytes16(DATAID2),
      roundId: 1,
      timestamp: TIMESTAMP2,
      answer: uint224(PRICE4)
    });

    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength2);
  }

  function test_onReportSuccess_EmptyBundleReport() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    s_dataFeedsCache.onReport(METADATA, s_emptyBundleReport);

    assertEq(s_dataFeedsCache.getLatestBundle(bytes16(DATAID1)), "");
    assertEq(s_dataFeedsCache.getLatestBundle(bytes16(DATAID2)), "");
    assertEq(s_dataFeedsCache.getLatestBundle(bytes16(DATAID3)), "");
    assertEq(s_dataFeedsCache.getLatestBundle(bytes16(DATAID4)), "");
    assertEq(s_dataFeedsCache.getLatestBundle(bytes16(DATAID5)), "");
  }

  function test_onReportSuccess_BundleReportLength1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals1By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    bytes memory expectedBundle =
      hex"000000000000000000000000000000000000000000000000000000000001e240000000000000000000000000000000000000000000000000000000000006f855";

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({dataId: bytes16(DATAID1), timestamp: TIMESTAMP1, bundle: expectedBundle});

    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength1);
  }

  function test_onReportSuccess_BundleReportLength2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals2By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);

    bytes memory expectedBundle1 =
      hex"00000000000000000000000000000000000000000000000000000000000c0bd000000000000000000000000000000000000000000000000000000000000d950b";

    bytes memory expectedBundle2 =
      hex"000000000000000000000000000000000000000000000000000000000009fbf100000000000000000000000000000000000000000000000000000000000f1206";

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({dataId: bytes16(DATAID1), timestamp: TIMESTAMP1, bundle: expectedBundle1});

    vm.expectEmit();
    emit DataFeedsCache.BundleReportUpdated({dataId: bytes16(DATAID2), timestamp: TIMESTAMP2, bundle: expectedBundle2});

    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength2);
  }

  function test_latestAnswer1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength1);

    vm.startPrank(s_proxyList[0]);
    int256 value = s_dataFeedsCache.latestAnswer();
    assertEq(value, int256(PRICE1));
  }

  function test_latestAnswer2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength2);

    vm.startPrank(s_proxyList[0]);
    int256 value = s_dataFeedsCache.latestAnswer();
    assertEq(value, int256(PRICE3));

    vm.startPrank(s_proxyList[1]);
    value = s_dataFeedsCache.latestAnswer();
    assertEq(value, int256(PRICE4));
  }

  function test_getLatestAnswer1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength1);
    vm.stopPrank();

    int256 value = s_dataFeedsCache.getLatestAnswer(dataIds[0]);
    assertEq(value, int256(PRICE1));
  }

  function test_getLatestAnswer2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength2);
    vm.stopPrank();

    int256 value = s_dataFeedsCache.getLatestAnswer(dataIds[0]);
    assertEq(value, int256(PRICE3));

    value = s_dataFeedsCache.getLatestAnswer(dataIds[1]);
    assertEq(value, int256(PRICE4));
  }

  function test_latestBundle1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals1By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength1);

    vm.startPrank(s_proxyList[0]);
    uint256 roundId = s_dataFeedsCache.latestRound();
    assertEq(roundId, 0);

    bytes memory bundle = s_dataFeedsCache.latestBundle();
    uint256 timestamp = s_dataFeedsCache.latestBundleTimestamp();
    uint8[] memory decimals = s_dataFeedsCache.bundleDecimals();
    assertEq(bundle, abi.encode(PRICE1, PRICE2));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, PRICE1);
    assertEq(firstBundleP2, PRICE2);
    assertEq(timestamp, TIMESTAMP1);
    assertEq(decimals.length, s_decimals1By1[0].length);
    assertEq(decimals[0], s_decimals1By1[0][0]);
  }

  function test_latestBundle2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals2By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength2);

    vm.startPrank(s_proxyList[0]);
    uint256 roundId = s_dataFeedsCache.latestRound();

    bytes memory bundle = s_dataFeedsCache.latestBundle();
    uint256 timestamp = s_dataFeedsCache.latestBundleTimestamp();
    uint8[] memory decimals = s_dataFeedsCache.bundleDecimals();
    assertEq(bundle, abi.encode(PRICE3, PRICE4));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, PRICE3);
    assertEq(firstBundleP2, PRICE4);
    assertEq(timestamp, TIMESTAMP1);
    assertEq(decimals.length, s_decimals2By1[0].length);
    assertEq(decimals[0], s_decimals2By1[0][0]);

    vm.startPrank(s_proxyList[1]);
    roundId = s_dataFeedsCache.latestRound();

    bundle = s_dataFeedsCache.latestBundle();
    timestamp = s_dataFeedsCache.latestBundleTimestamp();
    decimals = s_dataFeedsCache.bundleDecimals();
    assertEq(bundle, abi.encode(PRICE5, PRICE6));
    (uint256 secondBundleP1, uint256 secondBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(secondBundleP1, PRICE5);
    assertEq(secondBundleP2, PRICE6);
    assertEq(timestamp, TIMESTAMP2);
    assertEq(decimals.length, s_decimals2By1[1].length);
    assertEq(decimals[0], s_decimals2By1[1][0]);
  }

  function test_getLatestBundle1() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals1By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength1);
    vm.stopPrank();

    bytes memory bundle = s_dataFeedsCache.getLatestBundle(dataIds[0]);
    uint256 timestamp = s_dataFeedsCache.getLatestBundleTimestamp(dataIds[0]);
    uint8[] memory decimals = s_dataFeedsCache.getBundleDecimals(dataIds[0]);
    assertEq(bundle, abi.encode(PRICE1, PRICE2));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, PRICE1);
    assertEq(firstBundleP2, PRICE2);
    assertEq(timestamp, TIMESTAMP1);
    assertEq(decimals.length, s_decimals1By1[0].length);
    assertEq(decimals[0], s_decimals1By1[0][0]);
  }

  function test_getLatestBundle2() public {
    bytes16[] memory dataIds = new bytes16[](2);
    dataIds[0] = bytes16(DATAID1);
    dataIds[1] = bytes16(DATAID2);
    string[] memory _descriptions = new string[](2);
    _descriptions[0] = s_descriptions[0];
    _descriptions[1] = s_descriptions[0];

    s_dataFeedsCache.setBundleFeedConfigs(dataIds, _descriptions, s_decimals2By1, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_bundleReportlength2);
    vm.stopPrank();

    bytes memory bundle = s_dataFeedsCache.getLatestBundle(dataIds[0]);
    uint256 timestamp = s_dataFeedsCache.getLatestBundleTimestamp(dataIds[0]);
    uint8[] memory decimals = s_dataFeedsCache.getBundleDecimals(dataIds[0]);
    assertEq(bundle, abi.encode(PRICE3, PRICE4));
    (uint256 firstBundleP1, uint256 firstBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(firstBundleP1, PRICE3);
    assertEq(firstBundleP2, PRICE4);
    assertEq(timestamp, TIMESTAMP1);
    assertEq(decimals.length, s_decimals2By1[0].length);
    assertEq(decimals[0], s_decimals2By1[0][0]);

    bundle = s_dataFeedsCache.getLatestBundle(dataIds[1]);
    timestamp = s_dataFeedsCache.getLatestBundleTimestamp(dataIds[1]);
    decimals = s_dataFeedsCache.getBundleDecimals(dataIds[1]);
    assertEq(bundle, abi.encode(PRICE5, PRICE6));
    (uint256 secondBundleP1, uint256 secondBundleP2) = abi.decode(bundle, (uint256, uint256));
    assertEq(secondBundleP1, PRICE5);
    assertEq(secondBundleP2, PRICE6);
    assertEq(timestamp, TIMESTAMP2);
    assertEq(decimals.length, s_decimals2By1[1].length);
    assertEq(decimals[0], s_decimals2By1[1][0]);
  }

  function test_removeFeedsRevertInvalidSender() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    vm.startPrank(address(1002));

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.UnauthorizedCaller.selector, address(1002)));
    s_dataFeedsCache.removeFeedConfigs(dataIds);
  }

  function test_removeFeedsRevertNotConfiguredFeed() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    s_dataFeedsCache.setFeedAdmin(OWNER, true);

    vm.stopPrank();
    vm.startPrank(OWNER);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.FeedNotConfigured.selector, dataIds[0]));
    s_dataFeedsCache.removeFeedConfigs(dataIds);
  }

  function test_removeFeedsSuccess() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);

    DataFeedsCache.WorkflowMetadata[] memory wfMetadata;

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, s_descriptions, s_workflowMetadata);

    wfMetadata = s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 2);
    assertEq(wfMetadata.length, 2);
    bool hasPermission = s_dataFeedsCache.checkFeedPermission(dataIds[0], wfMetadata[0]);
    assertEq(hasPermission, true);

    s_dataFeedsCache.setFeedAdmin(OWNER, true);

    vm.stopPrank();
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit DataFeedsCache.FeedConfigRemoved(dataIds[0]);
    s_dataFeedsCache.removeFeedConfigs(dataIds);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.FeedNotConfigured.selector, dataIds[0]));
    s_dataFeedsCache.getFeedMetadata(dataIds[0], 0, 2);
    hasPermission = s_dataFeedsCache.checkFeedPermission(dataIds[0], wfMetadata[0]);
    assertEq(hasPermission, false);
  }

  function test_getDataIdForProxy() public view {
    bytes16 dataId = s_dataFeedsCache.getDataIdForProxy(s_proxyList[0]);
    assertEq(dataId, bytes16(DATAID1));
  }

  function test_recoverTokensRevertUnauthorized() public {
    vm.startPrank(ILLEGAL_CALLER);

    vm.expectRevert("Only callable by owner");
    s_dataFeedsCache.recoverTokens(IERC20(address(s_link)), address(10008), 1 ether);
  }

  function test_recoverTokensERC20RevertNoBalance() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InsufficientBalance.selector, 0, 1));
    s_dataFeedsCache.recoverTokens(IERC20(address(s_link)), address(10007), 1);
  }

  function testFuzzy_recoverTokensERC20Success(
    uint256 amount
  ) public {
    vm.assume(amount > 0);
    s_link.mint(address(s_dataFeedsCache), amount);

    vm.expectEmit();
    emit DataFeedsCache.TokenRecovered(address(s_link), address(10008), amount);
    s_dataFeedsCache.recoverTokens(IERC20(address(s_link)), address(10008), amount);
    assertEq(s_link.balanceOf(address(10008)), amount);
    assertEq(s_link.balanceOf(address(s_dataFeedsCache)), 0);
  }

  function test_recoverTokensNativeRevertNoBalance() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsCache.InsufficientBalance.selector, 0, 1 ether));
    s_dataFeedsCache.recoverTokens(IERC20(address(0)), address(10007), 1 ether);
  }

  function testFuzzy_recoverTokensNativeSuccess(
    uint256 amount
  ) public {
    vm.assume(amount > 0);
    vm.deal(address(s_dataFeedsCache), amount);
    assertEq(address(s_dataFeedsCache).balance, amount);

    vm.expectEmit();
    emit DataFeedsCache.TokenRecovered(address(0), address(10007), amount);
    s_dataFeedsCache.recoverTokens(IERC20(address(0)), address(10007), amount);
    assertEq(address(s_dataFeedsCache).balance, 0);
    assertEq(address(10007).balance, amount);
  }

  function test_getLatestByFeedId() public {
    bytes16[] memory dataIds = new bytes16[](1);
    dataIds[0] = bytes16(DATAID1);
    string[] memory _descriptions = new string[](1);
    _descriptions[0] = s_descriptions[0];

    s_dataFeedsCache.setDecimalFeedConfigs(dataIds, _descriptions, s_workflowMetadata);

    vm.startPrank(REPORT_SENDER);
    s_dataFeedsCache.onReport(METADATA, s_decimalReportlength1);

    uint256 timestamp = s_dataFeedsCache.getLatestTimestamp(dataIds[0]);
    assertEq(timestamp, TIMESTAMP1);

    (uint80 roundId, int256 answer, uint256 TIMESTAMP2, uint256 timestamp3, uint80 roundId2) =
      s_dataFeedsCache.getLatestRoundData(dataIds[0]);
    assertEq(roundId, 1);
    assertEq(roundId2, 1);
    assertEq(answer, int256(PRICE1));
    assertEq(timestamp, TIMESTAMP2);
    assertEq(timestamp, timestamp3);

    uint8 decimals = s_dataFeedsCache.getDecimals(dataIds[0]);
    assertEq(decimals, 18);

    string memory description = s_dataFeedsCache.getDescription(dataIds[0]);
    assertEq(description, s_descriptions[0]);
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
